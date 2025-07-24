SELECT id, email, name, password, is_active, created_at, updated_at
FROM users
WHERE email = $1 AND is_active = true;
