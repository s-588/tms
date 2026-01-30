-- +goose Up
-- +goose StatementBegin
create function set_updated_at_trigger_func()
returns trigger as $$
begin
    NEW.updated_at = now();
    return NEW;
end;
$$ language plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
create function set_deleted_at_trigger_func()
returns trigger as $$
begin
    NEW.deleted_at = now();
    return NEW;
end;
$$ language plpgsql;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
delete function set_updated_at_trigger_func;
-- +goose StatementEnd

-- +goose StatementBegin
delete function set_deleted_at_trigger_func;
-- +goose StatementEnd
