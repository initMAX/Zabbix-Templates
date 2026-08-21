SELECT
    t.tablename,
    CASE
        WHEN h.table_name IS NOT NULL THEN 'true'
        ELSE 'false'
    END AS is_hypertable,
    CASE
        WHEN h.table_name IS NOT NULL THEN
            (SELECT COALESCE(SUM(total_bytes), 0)
             FROM chunks_detailed_size('public.' || t.tablename))
        ELSE
            pg_total_relation_size('public.' || t.tablename)
    END AS size_bytes,
    CASE
        WHEN h.table_name IS NOT NULL THEN
            (SELECT COALESCE(SUM(cds.total_bytes), 0)
             FROM chunks_detailed_size('public.' || t.tablename) cds
             JOIN timescaledb_information.chunks c
                 ON c.chunk_schema = cds.chunk_schema
                AND c.chunk_name = cds.chunk_name
             WHERE c.is_compressed)
        ELSE NULL
    END AS compressed_bytes,
    CASE
        WHEN h.table_name IS NOT NULL THEN
            (SELECT COALESCE(SUM(cds.total_bytes), 0)
             FROM chunks_detailed_size('public.' || t.tablename) cds
             JOIN timescaledb_information.chunks c
                 ON c.chunk_schema = cds.chunk_schema
                AND c.chunk_name = cds.chunk_name
             WHERE NOT c.is_compressed)
        ELSE NULL
    END AS uncompressed_bytes,
    CASE
        WHEN h.table_name IS NOT NULL THEN
            (SELECT COUNT(*)
             FROM timescaledb_information.chunks c
             WHERE c.hypertable_schema = 'public'
               AND c.hypertable_name = t.tablename)
        ELSE NULL
    END AS chunks_total,
    CASE
        WHEN h.table_name IS NOT NULL THEN
            (SELECT COUNT(*)
             FROM timescaledb_information.chunks c
             WHERE c.hypertable_schema = 'public'
               AND c.hypertable_name = t.tablename
               AND c.is_compressed)
        ELSE NULL
    END AS chunks_compressed,
    CASE
        WHEN h.table_name IS NOT NULL THEN
            (SELECT COUNT(*)
             FROM timescaledb_information.chunks c
             WHERE c.hypertable_schema = 'public'
               AND c.hypertable_name = t.tablename
               AND NOT c.is_compressed)
        ELSE NULL
    END AS chunks_uncompressed
FROM pg_tables t
LEFT JOIN _timescaledb_catalog.hypertable h
    ON t.tablename = h.table_name AND h.num_dimensions = 1
WHERE t.schemaname = 'public'
ORDER BY
    CASE
        WHEN h.table_name IS NOT NULL THEN
            (SELECT COALESCE(SUM(total_bytes), 0)
             FROM chunks_detailed_size('public.' || t.tablename))
        ELSE
            pg_total_relation_size('public.' || t.tablename)
    END DESC;
