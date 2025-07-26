INSERT INTO users (email, name, password, is_active, created_at, updated_at)
VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING id, email, name, password, is_active, created_at, updated_at;
