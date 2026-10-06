package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"golang.zabbix.com/sdk/errs"
	"golang.zabbix.com/sdk/plugin"
	"golang.zabbix.com/sdk/plugin/container"
)

const (
	pluginName          = "BambuLab"
	metricKey           = "bambulab.get"
	defaultTimeoutSec   = 30
	defaultMinPrintKeys = 20
)

// Compile-time check that we implement plugin.Exporter.
var _ plugin.Exporter = (*bambuPlugin)(nil)

type bambuPlugin struct {
	plugin.Base
}

func main() {
	p := &bambuPlugin{}

	if err := plugin.RegisterMetrics(
		p,
		pluginName,
		metricKey,
		"Return first MQTT message from a topic over TLS with disabled certificate verification.",
	); err != nil {
		panic(errs.Wrap(err, "failed to register metrics"))
	}

	h, err := container.NewHandler(pluginName)
	if err != nil {
		panic(errs.Wrap(err, "failed to create handler"))
	}

	// Route plugin logs into the agent log.
	p.Logger = h

	if err := h.Execute(); err != nil {
		panic(errs.Wrap(err, "failed to execute plugin handler"))
	}
}

func (p *bambuPlugin) Export(key string, params []string, ctx plugin.ContextProvider) (any, error) {
	if key != metricKey {
		return nil, errs.Errorf("unknown item key %q", key)
	}

	if len(params) < 2 {
		return nil, errs.New(`usage: bambulab.get["tls://host:port","topic","user","password","min_print_keys","pushall"]`)
	}
	if len(params) > 6 {
		return nil, errs.New(`usage: bambulab.get["tls://host:port","topic","user","password","min_print_keys","pushall"]`)
	}

	rawBrokerURL := params[0]
	topic := params[1]
	if topic == "" {
		return nil, errs.New("topic must not be empty")
	}

	var user, pass string
	if len(params) > 2 {
		user = params[2]
	}
	if len(params) > 3 {
		pass = params[3]
	}
	minPrintKeys := defaultMinPrintKeys
	if len(params) > 4 && params[4] != "" {
		n, err := strconv.Atoi(params[4])
		if err != nil || n < 0 {
			return nil, errs.New("min_print_keys must be a non-negative integer")
		}
		minPrintKeys = n
	}
	pushAll := true
	if len(params) > 5 && params[5] != "" {
		v, err := strconv.ParseBool(params[5])
		if err != nil {
			return nil, errs.New("pushall must be a boolean (true/false/1/0)")
		}
		pushAll = v
	}

	brokerURL, u, isTLS, err := normalizeBrokerURL(rawBrokerURL)
	if err != nil {
		return nil, err
	}

	opts := mqtt.NewClientOptions().
		AddBroker(brokerURL).
		SetClientID(fmt.Sprintf("zbx-bambu-%d", time.Now().UnixNano())).
		SetAutoReconnect(false).
		SetConnectRetry(false).
		SetCleanSession(true)

	if user != "" {
		opts.SetUsername(user)
		opts.SetPassword(pass)
	}

	// This is the main point of the plugin: TLS with disabled certificate verification.
	if isTLS {
		opts.SetTLSConfig(&tls.Config{
			InsecureSkipVerify: true, //nolint:gosec // explicitly requested (Bambu uses mismatching/self-signed certs)
			ServerName:         u.Hostname(),
			MinVersion:         tls.VersionTLS12,
		})
	}

	timeoutSec := ctx.Timeout()
	if timeoutSec <= 0 {
		timeoutSec = defaultTimeoutSec
	}
	timeout := time.Duration(timeoutSec) * time.Second
	deadline := time.Now().Add(timeout)

	remaining := func() time.Duration {
		d := time.Until(deadline)
		if d < 0 {
			return 0
		}
		return d
	}

	msgCh := make(chan mqtt.Message, 1)
	errCh := make(chan error, 1)

	opts.SetConnectionLostHandler(func(_ mqtt.Client, e error) {
		select {
		case errCh <- e:
		default:
		}
	})

	client := mqtt.NewClient(opts)

	tok := client.Connect()
	if !tok.WaitTimeout(remaining()) {
		return nil, errs.New("timeout while connecting to MQTT broker")
	}
	if err := tok.Error(); err != nil {
		return nil, errs.Wrap(err, "failed to connect to MQTT broker")
	}
	defer client.Disconnect(250)

	tok = client.Subscribe(topic, 0, func(_ mqtt.Client, m mqtt.Message) {
		select {
		case msgCh <- m:
		default:
		}
	})
	if !tok.WaitTimeout(remaining()) {
		return nil, errs.New("timeout while subscribing to MQTT topic")
	}
	if err := tok.Error(); err != nil {
		return nil, errs.Wrap(err, "failed to subscribe to MQTT topic")
	}

	if pushAll {
		requestTopic, err := requestTopicFromReportTopic(topic)
		if err != nil {
			return nil, err
		}

		tok = client.Publish(
			requestTopic,
			0,
			false,
			[]byte(`{"pushing":{"sequence_id":1,"command":"pushall"},"user_id":"1234567890"}`),
		)
		if !tok.WaitTimeout(remaining()) {
			return nil, errs.New("timeout while publishing pushall")
		}
		if err := tok.Error(); err != nil {
			return nil, errs.Wrap(err, "failed to publish pushall")
		}
	}

	waitTimeout := remaining()
	if waitTimeout <= 0 {
		return nil, errs.New("timeout waiting for MQTT message")
	}

	waitTimer := time.NewTimer(waitTimeout)
	defer waitTimer.Stop()

	for {
		select {
		case <-waitTimer.C:
			return nil, errs.New("timeout waiting for MQTT message")
		case e := <-errCh:
			return nil, errs.Wrap(e, "connection lost")
		case msg := <-msgCh:
			payload := msg.Payload()
			if shouldAcceptPayload(payload, minPrintKeys) {
				return string(payload), nil
			}
		}
	}
}

func normalizeBrokerURL(raw string) (normalized string, parsed *url.URL, isTLS bool, err error) {
	if !strings.Contains(raw, "://") {
		raw = "tcp://" + raw
	}

	u, err := url.Parse(raw)
	if err != nil {
		return "", nil, false, errs.Wrap(err, "invalid broker url")
	}

	if u.Port() != "" && u.Hostname() == "" {
		return "", nil, false, errs.New("broker host is required")
	}

	if len(u.Query()) > 0 {
		return "", nil, false, errs.New("broker url must not contain query parameters")
	}

	switch strings.ToLower(u.Scheme) {
	case "tcp", "ws", "wss", "ssl", "tls", "mqtts":
	default:
		return "", nil, false, errs.Errorf("unsupported broker scheme %q", u.Scheme)
	}

	isTLS = strings.EqualFold(u.Scheme, "ssl") ||
		strings.EqualFold(u.Scheme, "tls") ||
		strings.EqualFold(u.Scheme, "mqtts") ||
		strings.EqualFold(u.Scheme, "wss")

	// Default ports.
	if u.Port() == "" {
		if isTLS {
			u.Host += ":8883"
		} else {
			u.Host += ":1883"
		}
	}

	// Paho uses "ssl://" for MQTT over TLS. Accept "tls://" and "mqtts://" as aliases.
	normalizedScheme := strings.ToLower(u.Scheme)
	if normalizedScheme == "tls" || normalizedScheme == "mqtts" {
		normalizedScheme = "ssl"
	}

	u2 := *u
	u2.Scheme = normalizedScheme

	return u2.String(), u, isTLS, nil
}

func requestTopicFromReportTopic(reportTopic string) (string, error) {
	parts := strings.Split(reportTopic, "/")
	if len(parts) < 3 || parts[0] != "device" || parts[2] != "report" {
		return "", errs.New("pushall requires topic starting with device/<SN>/report")
	}
	if parts[1] == "" {
		return "", errs.New("pushall requires non-empty serial in topic device/<SN>/report")
	}
	return fmt.Sprintf("device/%s/request", parts[1]), nil
}

func shouldAcceptPayload(payload []byte, minPrintKeys int) bool {
	if minPrintKeys <= 0 {
		return true
	}

	var root map[string]any
	if err := json.Unmarshal(payload, &root); err != nil {
		return false
	}

	printObj, ok := root["print"].(map[string]any)
	if !ok {
		return false
	}

	if len(printObj) >= minPrintKeys {
		return true
	}

	for _, key := range []string{
		"ams",
		"device",
		"net",
		"job",
		"upgrade_state",
		"lights_report",
		"ipcam",
		"online",
		"xcam",
	} {
		if _, ok := printObj[key]; ok {
			return true
		}
	}

	return false
}
