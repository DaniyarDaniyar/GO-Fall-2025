migrate -path internal/db/migrations -database "postgres://postgres@localhost:5432/expense_tracker?sslmode=disable" up     
20251006180205/u create_user_table (14.804ms)
20251006180302/u create_category_table (22.69ms)
20251006180333/u create_expense_table (35.7141ms)

migrate -path internal/db/migrations -database "postgres://postgres@localhost:5432/expense_tracker?sslmode=disable" version
20251006180333