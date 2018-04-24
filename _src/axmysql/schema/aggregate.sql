CREATE TABLE aggregate
(
    aggregate_id VARBINARY(255) PRIMARY KEY,
    revision     BIGINT UNSIGNED NOT NULL
)
ROW_FORMAT=COMPRESSED;

CREATE TABLE aggregate_event
(
    message_id     VARBINARY(255) PRIMARY KEY,
    aggregate_id   VARBINARY(255) NOT NULL,
    revision       BIGINT UNSIGNED NOT NULL,
    correlation_id VARBINARY(255) NOT NULL,
    causation_id   VARBINARY(255),
    content_type   VARBINARY(255) NOT NULL,
    body           LONGBLOB,

    UNIQUE (aggregate_id, revision),
    INDEX (correlation_id),
    INDEX (causation_id),
)
ROW_FORMAT=COMPRESSED;
