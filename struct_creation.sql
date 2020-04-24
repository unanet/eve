WITH models AS (
    WITH data AS (
        SELECT
            replace(initcap(table_name::text), '_', '') table_name,
            replace(initcap(column_name::text), '_', '') column_name,
            CASE data_type
                WHEN 'timestamp without time zone' THEN 'time.Time'
                WHEN 'timestamp with time zone' THEN 'time.Time'
                WHEN 'boolean' THEN 'bool'
                -- add your own type converters as needed or it will default to 'string'
                ELSE 'string'
                END AS type_info,
            '`db:"' || column_name ||'"`' AS annotation
        FROM information_schema.columns
        WHERE table_schema IN ('eve-api', 'public')
        ORDER BY table_schema, table_name, ordinal_position
    )
    SELECT table_name, STRING_AGG(E'\t' || column_name || E'\t' || type_info || E'\t' || annotation, E'\n') fields
    FROM data
    GROUP BY table_name
)
SELECT 'type ' || table_name || E' struct {\n' || fields || E'\n}' models
FROM models ORDER BY 1