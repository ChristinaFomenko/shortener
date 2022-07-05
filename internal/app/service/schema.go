package service

const CreateTable = `CREATE TABLE IF NOT EXISTS urls 
(
			id SERIAL PRIMARY KEY NOT NULL,
			hash VARCHAR(250),
			original_url TEXT NOT NULL,
			user_id VARCHAR(250)
)`
