CREATE TABLE IF NOT EXISTS yba (
    version VARCHAR(255) PRIMARY KEY,
    type VARCHAR(255),
    architecture VARCHAR(255),
    platform VARCHAR(255),
    commit VARCHAR(255),
    branch VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS ybdb (
    version VARCHAR(255) PRIMARY KEY,
    type VARCHAR(255),
    architecture VARCHAR(255),
    platform VARCHAR(255),
    download_url VARCHAR(255),
    commit VARCHAR(255),
    branch VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS yba_ybdb_compatibility (
    yba_version VARCHAR(255),
    ybdb_version VARCHAR(255),
    PRIMARY KEY (yba_version, ybdb_version),
    FOREIGN KEY (yba_version) REFERENCES yba(version),
    FOREIGN KEY (ybdb_version) REFERENCES ybdb(version)
);