CREATE TABLE documents (
    doc_id VARCHAR(255) PRIMARY KEY,
    type VARCHAR(50) NOT NULL,
    content BYTEA NOT NULL,
    sender_id VARCHAR(255),
    receiver_id UUID NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE
);


CREATE TABLE IF NOT EXISTS outbox (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    scheduled_at TIMESTAMP WITH TIME ZONE NOT NULL,
    metadata BYTEA,
    payload BYTEA NOT NULL,
    times_attempted INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_outbox_created_at ON outbox (created_at);
CREATE INDEX IF NOT EXISTS idx_outbox_scheduled_at ON outbox (scheduled_at);
