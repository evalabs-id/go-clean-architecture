package queries

import _ "embed"

var (
	//go:embed users/select.sql
	selectUsers string
	//go:embed users/get_by_id.sql
	getUserByID string
	//go:embed users/insert.sql
	insertUser string
	//go:embed users/update.sql
	updateUser string
)
