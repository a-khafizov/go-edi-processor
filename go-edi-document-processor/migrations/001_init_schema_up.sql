create table if not exists documents (
    id uuid primary key default gen_random_uuid(),
    external_id varchar(255) not null unique,
    document_type varchar(50) not null check (document_type in ('xml', 'pdf', 'json')),
    content_base64 text not null,
    sender_id varchar(255) not null,
    receiver_id varchar(255) not null,
    metadata jsonb,
    status varchar(20) not null default 'pending' check (status in ('pending', 'received', 'processed', 'failed')),
    processing_time timestamp with time zone,
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone default now()
);

create table if not exists outbox_messages (
    id bigserial primary key,
    topic varchar(255) not null,
    message_base64 text not null,
    key_base64 text,
    headers jsonb,
    created_at timestamp with time zone default now(),
    processed_at timestamp with time zone,
    delay interval default '0 seconds'
);

-- indexes for performance
create index if not exists idx_documents_external_id on documents(external_id);
create index if not exists idx_documents_status on documents(status);
create index if not exists idx_documents_sender_id on documents(sender_id);
create index if not exists idx_documents_receiver_id on documents(receiver_id);
create index if not exists idx_outbox_messages_topic on outbox_messages(topic);
create index if not exists idx_outbox_messages_created_at on outbox_messages(created_at);
create index if not exists idx_outbox_messages_processed_at on outbox_messages(processed_at);
