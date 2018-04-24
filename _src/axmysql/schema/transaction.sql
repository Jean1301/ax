CREATE TABLE transaction
(
    causation_id VARBINARY(255) PRIMARY KEY,
    revision     BIGINT UNSIGNED NOT NULL,
    state        VARBINARY(255) NOT NULL
)
ROW_FORMAT=COMPRESSED;

CREATE TABLE outbox_message
(
    message_id     VARBINARY(255) PRIMARY KEY,
    correlation_id VARBINARY(255) NOT NULL,
    causation_id   VARBINARY(255) NOT NULL,
    content_type   VARBINARY(255) NOT NULL,
    body           LONGBLOB,

    INDEX (correlation_id),
    INDEX (causation_id)
)
ROW_FORMAT=COMPRESSED;
