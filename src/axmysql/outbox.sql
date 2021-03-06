--
-- This file defines the SQL schema used by OutboxRepository
--

--
-- outbox stores the timestamp at which an outbox was created.
--
-- The presence of a row in this table indicates that the message has already
-- been handled, even if the outbox is now empty.
--
CREATE TABLE IF NOT EXISTS outbox (
    message_id  VARBINARY(255) NOT NULL,
    insert_time TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),

    PRIMARY KEY (message_id),
    INDEX (insert_time)
);

--
-- outbox_message stores the messages within a single outbox.
--
CREATE TABLE IF NOT EXISTS outbox_message (
    message_id     VARBINARY(255) NOT NULL,
    causation_id   VARBINARY(255) NOT NULL, -- outbox.message_id
    correlation_id VARBINARY(255) NOT NULL,
    time           VARBINARY(255) NOT NULL,
    content_type   VARBINARY(255) NOT NULL,
    body           BLOB NOT NULL,
    operation      INTEGER NOT NULL,
    destination    VARBINARY(255) NOT NULL,

    PRIMARY KEY (message_id),
    INDEX (causation_id),
    INDEX (correlation_id)
) ROW_FORMAT=COMPRESSED;
