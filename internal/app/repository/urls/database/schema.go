package database

const (
	CreateTable = `CREATE TABLE IF NOT EXISTS urls 
(
			id SERIAL PRIMARY KEY NOT NULL,
			original_url TEXT NOT NULL,
			short_url TEXT NOT NULL,
			user_id VARCHAR(250)
);

CREATE UNIQUE INDEX IF NOT EXISTS urls_short_id_user_id_uniq_idx ON urls (short_id, user_id);
ALTER TABLE urls DROP CONSTRAINT IF EXISTS urls_short_url_uniq_for_user;
ALTER TABLE urls ADD CONSTRAINT urls_short_url_uniq_for_user UNIQUE USING INDEX urls_short_id_user_id_uniq_idx;
CREATE INDEX IF NOT EXISTS urls_user_id_idx ON urls(user_id);
`
	InsertStatement = `INSERT INTO urls (original_url, short_id, user_id) VALUES($1, $2, $3) "+
			"ON CONFLICT ON CONSTRAINT urls_short_url_uniq_for_user DO "+
			"UPDATE SET original_url = EXCLUDED.original_url`
)
