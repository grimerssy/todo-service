CREATE TABLE users (
    id SERIAL NOT NULL,
    first_name VARCHAR(64) NOT NULL,
    last_name VARCHAR(64) NOT NULL,
    email VARCHAR(64) NOT NULL,
    username VARCHAR(64) NOT NULL,
    password VARCHAR(64) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_users_id PRIMARY KEY (id)
);
CREATE INDEX idx_users_username ON users (username);

CREATE TABLE todos (
    id SERIAL NOT NULL,
    title VARCHAR(64) NOT NULL,
    description VARCHAR(512) NULL,
    completed BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_todos_id PRIMARY KEY (id)
);

CREATE TABLE users_todos (
    user_id INTEGER NOT NULL,
    todo_id INTEGER NOT NULL,
    CONSTRAINT pk_users_todos PRIMARY KEY (user_id, todo_id),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_todo_id FOREIGN KEY (todo_id) REFERENCES todos(id) ON DELETE CASCADE
);

CREATE FUNCTION set_updated_at()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_updated_at
    BEFORE UPDATE
    OF first_name, last_name, email, username, password
    ON users
    FOR EACH ROW
EXECUTE PROCEDURE set_updated_at();

CREATE TRIGGER todos_updated_at
    BEFORE UPDATE
    OF title, description, completed
    ON todos
    FOR EACH ROW
EXECUTE PROCEDURE set_updated_at();
