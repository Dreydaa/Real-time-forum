package database

// Création des tables dans la database

const (
	CreateTables = `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		firstname TEXT NOT NULL,
		lastname TEXT NOT NULL,
		gender TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);
	
	CREATE TABLE IF NOT EXISTS sessions (
		session_uuid TEXT NOT NULL PRIMARY KEY UNIQUE,
		user_id INTEGER NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		category TEXT NOT NULL,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		date TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		date TEXT NOT NULL,
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sender_id INTEGER NOT NULL,
		receiver_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		date TEXT NOT NULL,
		FOREIGN KEY (sender_id) REFERENCES users(id),
		FOREIGN KEY (receiver_id) REFERENCES users(id)
	);
	
	CREATE TABLE IF NOT EXISTS chats (
		id_one INTEGER NOT NULL,
		id_two INTEGER NOT NULL,
		time INTEGER NOT NULL,
		FOREIGN KEY (id_one) REFERENCES users(id),
		FOREIGN KEY (id_two) REFERENCES users(id)
	);
	`
)
