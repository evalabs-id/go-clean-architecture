package queries

import _ "embed"

var (
	//go:embed users/select.sql
	SelectUsers string
	//go:embed users/get_by_id.sql
	GetUserByID string
	//go:embed users/get_by_email.sql
	GetUserByEmail string
	//go:embed users/insert.sql
	InsertUser string
	//go:embed users/update.sql
	UpdateUser string
	//go:embed users/check_email.sql
	CountCheckEmail string
)
