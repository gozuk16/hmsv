-- name: InsertCpuOriginalData :exec
INSERT INTO cpu_original (
	timestamp,
	cpu_usage
) VALUES (
	?, ?
);

-- name: InsertCpuDownsampledData :exec
INSERT INTO cpu_downsampled (
	timestamp,
	avg_cpu_usage,
	max_cpu_usage
) VALUES (
	?, ?, ?
);

-- name: SelectCpuDownsamplingData :many
SELECT
	CAST(strftime('%Y-%m-%d %H:%M:%S', datetime((strftime('%s', timestamp) / @duration) * @duration, 'unixepoch')) AS TEXT) as dstimestamp,
	CAST(AVG(cpu_usage) AS REAL) as ave_cpu_usage,
	CAST(MAX(cpu_usage) AS REAL) as max_cpu_usage
FROM
	cpu_original
GROUP BY dstimestamp
ORDER BY dstimestamp
;
