-- +goose Up
-- +goose StatementBegin
CREATE TRIGGER trigger_employees_updated_at
    BEFORE UPDATE ON employees
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at_trigger_func();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER trigger_employees_updated_at ON employees;
-- +goose StatementEnd