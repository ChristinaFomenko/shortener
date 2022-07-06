package urls

const (
	CreateTable = `CREATE TABLE IF NOT EXISTS urls 
(
			id SERIAL PRIMARY KEY NOT NULL,
			original_url TEXT NOT NULL,
			short_url TEXT NOT NULL,
			user_id VARCHAR(250)
)`
	InsertStatement = `INSERT INTO urls(original_url, user_id) VALUES($1, $2) RETURNING id`
)
