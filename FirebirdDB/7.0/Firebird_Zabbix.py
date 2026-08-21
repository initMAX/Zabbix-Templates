#!/usr/bin/env python3
"""
Firebird → Zabbix monitoring helper
----------------------------------
• Collects many MON$-based statistics (incl. page I/O, mem usage, txn gap, …)
• Works on Firebird 3.x; for columns that exist only in 4/5 it returns 0
• Can either print the value (for testing) or push it with zabbix_sender
"""

import argparse
import subprocess
import sys
from typing import Any, Dict

import firebirdsql


class FirebirdZabbix:
    def __init__(
        self,
        host: str,
        port: int,
        database: str,
        username: str,
        password: str,
        charset: str = "UTF8",
    ):
        self.connection_params = {
            "host": host,
            "port": port,
            "database": database,
            "user": username,
            "password": password,
            "charset": charset,
        }

    # ------------------------------------------------------------------ DB helpers
    def query(self, sql: str) -> Any:
        """Run SQL and return the first column of the first row."""
        with firebirdsql.connect(**self.connection_params) as conn:
            cur = conn.cursor()
            cur.execute(sql)
            row = cur.fetchone()
            return row[0] if row else None

    def safe_query(self, sql: str, default: int | float = 0):
        """
        Execute SQL and return its first value.

        • If the server reports “column unknown”, “token unknown”, or any
          GDS status -206 / -104, we return *default* (0) instead.
        • Any other error is re-raised so we still notice real problems.
        """
        try:
            return self.query(sql)

        except firebirdsql.DatabaseError as exc:
            msg = str(exc).lower()

            # 1) look for GDS status codes (if present)
            codes = getattr(exc, "gds_codes", []) or []
            if any(code in (-206, -104) for code in codes):
                return default

            # 2) fall back to text pattern matching (older driver builds)
            if (
                "column unknown" in msg
                or "token unknown" in msg
                or "dynamic sql error" in msg
            ):
                return default

            raise  # anything else = real failure

    # ---------------------------------------------------------------- Zabbix sender
    @staticmethod
    def zbx_send(host: str, key: str, value: Any, sender: str, agent_conf: str) -> int:
        args = [sender, "-c", agent_conf, "-s", host, "-k", key, "-o", str(value)]
        proc = subprocess.run(args, capture_output=True, text=True)
        if proc.returncode != 0:
            print(f"Zabbix sender error: {proc.stderr.strip()}", file=sys.stderr)
        return proc.returncode

    # --------------------------------------------------------------------- Main API
    def run_check(self, check: str) -> Any:
        """Return the value for *check* using the queries dict below."""
        queries: Dict[str, str] = {
            # -------------------- Database state
            "database_status": "select mon$shutdown_mode from mon$database",
            "database_readonly": "select mon$read_only from mon$database",
            "database_backup": "select mon$backup_state from mon$database",
            # -------------------- Memory
            "memory_usage": (
                "select mon$memory_used from mon$memory_usage where mon$stat_group = 0"
            ),
            "memory_alloc": (
                "select mon$memory_allocated from mon$memory_usage where mon$stat_group = 0"
            ),
            # -------------------- Page I/O
            "iostat_reads": "select mon$page_reads from mon$io_stats where mon$stat_group = 0",
            "iostat_writes": "select mon$page_writes from mon$io_stats where mon$stat_group = 0",
            "iostat_fetches": "select mon$page_fetches from mon$io_stats where mon$stat_group = 0",
            "iostat_marks": "select mon$page_marks from mon$io_stats where mon$stat_group = 0",
            # -------------------- Attachments
            "attach_idle": "select count(*) from mon$attachments where mon$state = 0",
            "attach_active": "select count(*) from mon$attachments where mon$state = 1",
            # -------------------- Record reads
            "record_sequences": (
                "select mon$record_seq_reads from mon$record_stats where mon$stat_group = 0"
            ),
            "record_indexes": (
                "select mon$record_idx_reads from mon$record_stats where mon$stat_group = 0"
            ),
            # -------------------- SQL statements
            "query_idle": "select count(*) from mon$statements where mon$state = 0",
            "query_active": "select count(*) from mon$statements where mon$state = 1",
            "query_stalled": "select count(*) from mon$statements where mon$state = 2",
            # -------------------- Transactions (basic)
            "trans_idle": "select count(*) from mon$transactions where mon$state = 0",
            "trans_active": "select count(*) from mon$transactions where mon$state = 1",
            # ==================== NEW METRICS ==============================
            # Txn gap (NEXT_TXN − OLDEST_ACTIVE)
            "txn_gap": """
                select mon$next_transaction - mon$oldest_active
                from   mon$database
            """,
            # Sweep
            "sweep_interval": "select mon$sweep_interval from mon$database",
            # FB-3 fallback: return 0 when column is missing
            "sweep_active": "select mon$sweep_in_progress from mon$database",
            # DB settings
            "forced_writes": "select mon$forced_writes from mon$database",
            "reserve_space": "select mon$reserve_space from mon$database",
            # Size
            "pages_alloc": "select mon$pages from mon$database",
            "db_size_bytes": "select mon$pages * mon$page_size from mon$database",
            # Cache hit ratio (%)
            "cache_hit_ratio": """
                select cast(
                    100.0 * (mon$page_fetches - mon$page_reads)
/ nullif(mon$page_fetches,0)
                    as numeric(15,2)
                )
                from mon$io_stats
                where mon$stat_group = 0
            """,
            # Encryption & replication (FB 4+)  → 0 on FB 3
            "crypt_state": "select coalesce(mon$crypt_state,0) from mon$database",
            "replica_mode": "select coalesce(mon$replica_mode,0) from mon$database",
            # Lock conflicts (FB 5+) → 0 on FB 3/4
            "lock_conflicts": "select coalesce(mon$lock_conflicts,0) from mon$database",
            # Long-running tx (> 300 s); (timestamp difference is in *days*)
            "trans_long_running": """
                select count(*)
                from   mon$transactions
                where  (current_timestamp - mon$timestamp) * 86400 > 300
            """,
        }

        if check not in queries:
            print(f"Unknown check: {check}", file=sys.stderr)
            sys.exit(1)

        return self.safe_query(queries[check])

    # ------------------------------------------------------------------


def parse_args():
    p = argparse.ArgumentParser(description="Firebird Zabbix Monitoring")
    p.add_argument("-H", "--host", required=True, help="Firebird server host/IP")
    p.add_argument("-P", "--port", type=int, default=3050, help="Firebird port (default 3050)")
    p.add_argument("-d", "--database", required=True, help="Database path or alias")
    p.add_argument("-u", "--user", default="sysdba", help="DB user (default SYSDBA)")
    p.add_argument("-p", "--password", default="masterkey", help="DB password")
    p.add_argument("--charset", default="UTF8", help="Connection charset (default UTF8)")
    p.add_argument("-c", "--check", required=True, help="Metric name (key)")
    # Zabbix sender options
    p.add_argument("--zabbix-sender", default="zabbix_sender", help="Path to zabbix_sender")
    p.add_argument("--zabbix-conf", default="/etc/zabbix/zabbix_agentd.conf", help="Agent conf")
    p.add_argument("--zabbix-host", help="Host name as defined in Zabbix")
    return p.parse_args()


def main():
    args = parse_args()

    fb = FirebirdZabbix(
        host=args.host,
        port=args.port,
        database=args.database,
        username=args.user,
        password=args.password,
        charset=args.charset,
    )

    value = fb.run_check(args.check)

    if args.zabbix_host:
        exit_code = fb.zbx_send(
            host=args.zabbix_host,
            key=args.check,
            value=value,
            sender=args.zabbix_sender,
            agent_conf=args.zabbix_conf,
        )
        sys.exit(exit_code)
    else:
        # Local test run: just print the number
        print(value)


if __name__ == "__main__":
    main()