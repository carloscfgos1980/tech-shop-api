-- name: CreateEmployee :one
INSERT INTO employees (id, created_at, updated_at, email, password, role)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetEmployeeByEmail :one
SELECT * FROM employees
WHERE email = $1;

-- name: GetAdminById :one
SELECT * FROM employees
WHERE id = $1;