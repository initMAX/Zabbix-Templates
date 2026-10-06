#!/usr/bin/env python3

import os
import json
import sys

def discover_databases(path):
    databases = []
    for root, dirs, files in os.walk(path):
        for file in files:
            if file.lower().endswith('.fdb'):
                fullpath = os.path.join(root, file)
                dbname = os.path.splitext(file)[0]
                databases.append({
                    "{#FIREBIRD.DB.PATH}": fullpath,
                    "{#FIREBIRD.DB.NAME}": dbname
                })
    return databases

if __name__ == "__main__":
    db_folder = sys.argv[1] if len(sys.argv) > 1 else "/var/lib/firebird/3.0/data"
    dbs = discover_databases(db_folder)
    print(json.dumps(dbs, indent=4))