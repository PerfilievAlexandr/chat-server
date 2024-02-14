-- +goose Up
-- +goose StatementBegin
create table messages
(
    id         uuid primary key,
    text       text      not null,
    owner      text      not null,
    chat_id    uuid REFERENCES chat (id),
    created_at timestamp not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table messages;
-- +goose StatementEnd
