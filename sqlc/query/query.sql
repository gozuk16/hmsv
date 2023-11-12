-- name: InsertCpuOriginalData :exec
INSERT INTO cpu_original (
	timestamp
	,cpu_usage
) VALUES (
	?, ?
);

-- name: InsertCpuDownsampledData :exec
INSERT INTO cpu_downsampled (
	timestamp
	,avg_cpu_usage
	,max_cpu_usage
) VALUES (
	?, ?, ?
);

-- name: SelectCpuDownsamplingData :many
SELECT
	strftime('%Y-%m-%d %H:%M:%S', datetime((strftime('%s', timestamp) / @duration) * @duration, 'unixepoch')) AS dstimestamp
	,AVG(cpu_usage) as ave_cpu_usage
	,MAX(cpu_usage) as max_cpu_usage
FROM
	cpu_original
GROUP BY dstimestamp
ORDER BY dstimestamp
;
