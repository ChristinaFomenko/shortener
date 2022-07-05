package service

const CreateTable = `create table if not exists urls (
			id serial primary key not null,
			hash varchar(25),
			original_url text,
			user_id varchar(250),
		)`
