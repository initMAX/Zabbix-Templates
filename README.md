<!-- *********************************************************************************************************************************** -->
<!-- *** HEADER ************************************************************************************************************************ -->
<!-- *********************************************************************************************************************************** -->
<div align="center">
    <a href="http://www.initmax.com"><img src="./.readme/logo/initMAX_banner.png" alt="initMAX Logo"></a>
    <h3>
        <span>
            Honesty, diligence and MAXimum knowledge of our products is our standard.
        </span>
    </h3>
    <h3>
        <a href="https://www.linkedin.com/company/initmax/">
            <img alt="Static Badge" src="./.readme/logo/linkedin.png">
        </a>&nbsp;&nbsp;&nbsp;
        <a href="https://www.youtube.com/@initmax1">
            <img alt="Static Badge" src="./.readme/logo/youtube.png">
        </a>&nbsp;&nbsp;&nbsp;
        <a href="https://www.facebook.com/initmax">
            <img alt="Static Badge" src="./.readme/logo/facebook.png">
        </a>&nbsp;&nbsp;&nbsp;
        <a href="https://www.instagram.com/initmax/">
            <img alt="Static Badge" src="./.readme/logo/instagram.png">
        </a>&nbsp;&nbsp;&nbsp;
        <a href="https://x.com/initmax">
            <img alt="Static Badge" src="./.readme/logo/x.png">
        </a>&nbsp;&nbsp;&nbsp;
        <a href="https://github.com/initmax">
            <img alt="Static Badge" src="./.readme/logo/github.png">
        </a>
    </h3>
</div>
<br>
<br>

---
---

<div align="center">
    <h1>
        Zabbix Templates <!-- !!! change version !!! -->
    </h1>
</div>
<br>
<br>

<!-- *********************************************************************************************************************************** -->
<!-- *** BODY ************************************************************************************************************************** -->
<!-- *********************************************************************************************************************************** -->
## Templates

This repository contains the following public Zabbix templates:

- [AI on the edge device HTTP](./Meter_ESP32_by_HTTP)
  - This template integrates jomjol’s AI-on-the-Edge ESP32 meter project with Zabbix over HTTP.

- [AI on the edge device MQTT](./Meter_ESP32_by_MQTT)
  - This template integrates jomjol’s AI-on-the-Edge ESP32 meter project with Zabbix over MQTT via Zabbix Agent 2.

- [BambuLAB printers by Zabbix Agent 2](./BambuLab)
  - This template monitors BambuLab 3D printers via MQTT using a Zabbix Agent 2 plugin.

- [CVE_2024_3094](./CVE_2024_3094)
  - Checks host for vulnerability CVE-2024-3094.

- [CVE_zabbix_threats](./CVE_zabbix_threats)
  - Evaluates Zabbix CVEs for the effective version and creates detailed alerts.

- [CVE_zabbix_threats_check (ARCHIVED)](./CVE_zabbix_threats_check)
  - Checks current Zabbix server version for CVEs and creates alerts. (ARCHIVED → [CVE_zabbix_threats](./CVE_zabbix_threats))

- [CVE_zabbix_threats_check_controlled_by_macro (ARCHIVED)](./CVE_zabbix_threats_check_controlled_by_macro)
  - Evaluates CVEs for different Zabbix server versions using a macro-defined version. (ARCHIVED → [CVE_zabbix_threats](./CVE_zabbix_threats))

- [EHEIM_Digital](./EHEIM_Digital)
  - Template for monitoring EHEIM Digital aquarium devices (chillers, feeders, filters, LED controllers) via their HTTP API.

- [Fail2ban](./Fail2ban)
  - Template for Fail2ban with automatic jail discovery, metrics (counts + IP list), and ready‑made triggers/dashboard.

- [Firebird](./FirebirdDB)
  - This template monitors Firebird databases (v3+)

- [FreeIPA Client](./FreeIPA_Client)
  - Native-first monitoring for FreeIPA clients, including DNS resolver validation and response time checks.

- [FreeIPA Server](./FreeIPA_Server)
  - This template monitors FreeIPA Server health checks via ipa-healthcheck JSON output.

- [FS Switch S3270](./FS_Switch_S3270)
  - SNMP monitoring of FS S3270 switches - interfaces, CPU, memory, PSU and temperature.

- [FS Switch S3400](./FS_Switch_S3400)
  - SNMP monitoring of FS S3400 switches - interfaces, CPU, memory, PoE ports and PSU.

- [FS Switch S3900](./FS_Switch_S3900)
  - SNMP monitoring of FS S3900 switches - interfaces, CPU, memory, fans and PSU.

- [FS Switch S5810](./FS_Switch_S5810)
  - SNMP monitoring of FS S5810 switches - interfaces, CPU, memory, fans, PSU and temperature.

- [Logstash](./Logstash)
  - Monitors Logstash instances via the Monitoring API (HTTP) — JVM, pipelines, events, and node health.

- [Media_HTML_Template](./Media_HTML_Template)
  - Reusable HTML email media type with severity-aware colors, light/dark mode support, and smart action buttons.

- [Media_WhatsApp](./Media_WhatsApp)
  - Integrates Zabbix notifications with WhatsApp Cloud API using webhook.

- [Multiple_Website_certificate_by_Zabbix_agent_2](./Multiple_Website_certificate_by_Zabbix_agent_2)
  - Monitors multiple TLS/SSL certificates using Zabbix agent 2.

- [Net_Planet_Switch](./Net_Planet_Switch)
  - Monitors Planet switches.

- [Not_Supported_Items](./Not_Supported_Items)
  - Counts Unsupported and No Longer Discovered Items on Host.

- [Open_Files](./Open_Files)
  - Checks OpenFiles settings for Zabbix server processes.

- [Proxmox Backup Server](./Proxmox_Backup_Server)
  - Agentless monitoring of Proxmox Backup Server (PBS) via the HTTP API - host metrics, automatic datastore discovery, garbage collection and backup freshness.

- [PrusaLink 3D Printer](./PrusaLink_3D_Printer)
  - Base template for monitoring Prusa 3D printers via PrusaLink HTTP API.

- [SELinux_by_Zabbix_Agent_2](./SELinux_by_Zabbix_Agent_2)
  - Enables monitoring of basic SELinux properties

- [SMSEagle (SMS & voice media type)](./Media_SMSEagle)
  - SMS and voice-call (TTS) notification media types for SMSEagle gateways, with two-way acknowledge by SMS.

- [SMSEagle by HTTP](./SMSEagle_by_HTTP)
  - Monitors SMSEagle SMS gateway devices via HTTP API v2 — modem status, signal strength, message counters, firmware, support, HA, services, and temperature sensors.

- [Solax_Solar_power_plant](./Solax_Solar_power_plant)
  - Template for monitoring Solax X3 - Hybrid using modbus protocol.

- [SSL_Labs_Certificate_Grade](./SSL_Labs_Certificate_Grade)
  - Zabbix template that checks public SSL/TLS configuration via the Qualys SSL Labs API.

- [Synology_Active_Backup_for_Business](./Synology_Active_Backup_for_Business)
  - Monitors Synology Active Backup for Business via SSH + SQLite — task and device discovery, backup status, and failure alerts.

- [Synology_NAS](./Synology_NAS)
  - Monitors Synology DiskStation via SNMP walk with dependent items and LLD.

- [TCP_UDP_Sockets](./TCP_UDP_Sockets)
  - Creates items for standard TCP/UDP socket checks.

- [Ubiquiti UniFi UAP Access Point](./Ubiquiti_UniFi_UAP_Access_Point)
  - SNMP monitoring of Ubiquiti UniFi access points - interfaces, CPU, memory, radios and clients.

- [Ubiquiti UniFi USW Pro Switch](./Ubiquiti_UniFi_USW_Pro_Switch)
  - SNMP monitoring of Ubiquiti UniFi USW Pro switches - interfaces, CPU, memory, fans and temperature.

- [Ubiquiti UniFi USW Switch](./Ubiquiti_UniFi_USW_Switch)
  - SNMP monitoring of Ubiquiti UniFi USW switches - interfaces, CPU and memory.

- [Windows Defender](./Windows_Defender_by_Zabbix_Agent)
  - Monitoring Microsoft Windows Defender via Zabbix Agent - protection states, signatures, scans and threat events from the event log.

- [Zabbix_DB_tables_size](./Zabbix_DB_tables_size)
  - Checks and monitors the sizes of Zabbix database tables. 

- [Zabbix_Server_version](./Zabbix_Server_version)
  - Checks Zabbix server version and triggers if not up to date.


<!-- *********************************************************************************************************************************** -->
<!-- *** FOOTER ************************************************************************************************************************ -->
<!-- *********************************************************************************************************************************** -->
<br>
<br>

---
---
<div align="center">
    <h4>
        <a href="https://www.initmax.com/" style="display: inline-flex; align-items: center; gap: 5px; color: inherit; text-decoration: none;">
            <img alt="Static Badge" src="./.readme/logo/initMAX_symbol_cerne.svg" height="24">
            initMAX.com
        </a>&nbsp;&nbsp;&nbsp;
        <a href="tel:+420800244442" style="display: inline-flex; align-items: center; gap: 5px; color: inherit; text-decoration: none;">
            <img alt="Static Badge" src="./.readme/logo/phone.svg" height="24">
            +420800244442
        </a>&nbsp;&nbsp;&nbsp;
        <a href="mailto:info@initmax.com" style="display: inline-flex; align-items: center; gap: 5px; color: inherit; text-decoration: none;">
            <img alt="Static Badge" src="./.readme/logo/mail.svg" height="24">
            info@initmax.com
        </a>
        <br><br>
        <a href="https://www.linkedin.com/company/initmax/">
            <img alt="Static Badge" src="./.readme/logo/linkedin.png">
        </a>&nbsp;
        <a href="https://www.youtube.com/@initmax1">
            <img alt="Static Badge" src="./.readme/logo/youtube.png">
        </a>&nbsp;
        <a href="https://www.facebook.com/initmax">
            <img alt="Static Badge" src="./.readme/logo/facebook.png">
        </a>&nbsp;
        <a href="https://www.instagram.com/initmax/">
            <img alt="Static Badge" src="./.readme/logo/instagram.png">
        </a>&nbsp;
        <a href="https://x.com/initmax">
            <img alt="Static Badge" src="./.readme/logo/x.png">
        </a>&nbsp;
        <a href="https://github.com/initmax">
            <img alt="Static Badge" src="./.readme/logo/github.png">
        </a><br><br><br>
        <a><img src="./.readme/logo/zabbix-premium-partner.png" alt="Zabbix premium partner" width="80"></a>&nbsp;&nbsp;&nbsp;
        <a><img src="./.readme/logo/zabbix-certified-trainer.png" alt="Zabbix certified trainer" width="80"></a>
        <br><br><br>
        <a>
            <img src="./.readme/logo/agplv3.png" width="100">
        </a>
    </h4>
</div>
