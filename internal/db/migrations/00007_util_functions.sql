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

-- +goose Down
-- +goose StatementBegin
drop function if exists set_updated_at_trigger_func();
-- +goose StatementEnd
