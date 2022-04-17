DROP TRIGGER todos_updated_at ON todos;
DROP TRIGGER users_updated_at ON users;

DROP FUNCTION set_updated_at();

DROP TABLE users_todos;
DROP TABLE todos;
DROP TABLE users;
