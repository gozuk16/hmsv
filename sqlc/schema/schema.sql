-- CPUデータ（オリジナル）
CREATE TABLE IF NOT EXISTS cpu_original (
    id             INTEGER PRIMARY KEY
    ,timestamp     TEXT    NOT NULL
    ,cpu_usage     REAL    NOT NULL
);

-- CPUデータ（ダウンサンプリングされたデータ）
CREATE TABLE IF NOT EXISTS cpu_downsampled (
    id             INTEGER PRIMARY KEY
    ,timestamp     TEXT    NOT NULL
    ,avg_cpu_usage REAL    NOT NULL
    ,max_cpu_usage REAL    NOT NULL
);

