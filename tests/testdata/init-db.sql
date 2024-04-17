CREATE TABLE IF NOT EXISTS customers (
    id varchar(255),
    document_id varchar(255),
    document_type int,
    is_anonymous boolean,
    password varchar(255),
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (id)
);

INSERT INTO customers (id, document_id, document_type, is_anonymous, password, created_at, updated_at)
    VALUES ('1', '548.644.620-97', 1, false, '$2a$10$2jAz4kMPG1VgrUpxlq/KR.lsJMSFuG7ghpB1qdH5DL9cSqZOKcsFu', NOW(), NOW());