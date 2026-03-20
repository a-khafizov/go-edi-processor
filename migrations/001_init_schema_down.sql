drop index if exists idx_documents_external_id;
drop index if exists idx_documents_status;
drop index if exists idx_documents_sender_id;
drop index if exists idx_documents_receiver_id;
drop index if exists idx_outbox_messages_topic;
drop index if exists idx_outbox_messages_created_at;
drop index if exists idx_outbox_messages_processed_at;

drop table if exists documents;
drop table if exists outbox_messages;
