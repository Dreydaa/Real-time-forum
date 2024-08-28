package database

const (
	AddUser     = `INSERT INTO users(username, firstname, lastname, gender, email, password) VALUES (?,?,?,?,?,?)`
	AddPost     = `INSERT INTO posts(user_id, category, title, content, date) VALUES (?,?,?,?,?)`
	AddComment  = `INSERT INTO comments(post_id, user_id, content, date) VALUES (?,?,?,?)`
	AddMessage  = `INSERT INTO messages(sender_id, receiver_id, content, date) VALUES (?,?,?,?)`
	AddSessions = `INSERT INTO sessions(session_uuid, user_id) VALUES (?,?)`
	AddChat     = `INSERT INTO chats(id_one, id_two, time) VALUES (?,?,?)`
)

const (
	GetUserByID          = `SELECT * FROM users WHERE id = ?`
	GetUserByUsername    = `SELECT * FROM users WHERE username = ?`
	GetUserByEmail       = `SELECT * FROM users WHERE email = ?`
	GetAllUser           = `SELECT * FROM users ORDER BY username ASC`
	GetPostById          = `SELECT * FROM posts Where id = ? ORDER BY id DESC`
	GetAllPost           = `SELECT * FROM posts ORDER id DESC`
	GetAllPostByCategory = `SELECT * FROM posts WHERE category = ? ORDER BY id DESC`
	GetAllPostByUser     = `SELECT * FROM posts WHERE user_id = ? ORDER BY id DESC`
	GetCommentById       = `SELECT * FROM comments WHERE id = ?`
	GetAllPostComment    = `SELECT * FROM comments WHERE post_id = ?`
	GetAllUserComment    = `SELECT * FROM comments WHERE user_id = ?`
	GetMessage           = `SELECT * FROM messages WHERE id = ?`
	GetAllChatMessage    = `SELECT * FROM messages WHERE ((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)) AND ( id <= ? ) ORDER BY id DESC LIMIT 10`
	GetLastMessage       = `SELECT * FROM messages WHERE ((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)) ORDER BY id DESC LIMIT 1`
	GetSessionUser       = `SELECT users * FROM seesions INNER JOIN users ON sessions.users_id = users.id WHERE sessions.session_uuid = ?`
	GetUserChats         = `SELECT * FROM chats WHERE id_one = ? OR id_two = ? ORDER BY time DESC`
	GetChatBetween       = `SELECT * FROM chats Where id_one = ? AND id_two = ? OR id_one = ? AND id_two = ?`
)

const (
	RemoveCookie = `DELETE FROM sessions WHERE user_id = ?`
)

const (
	UpdateChat = `UPDATE chats SET time = ? WHERE id_one = ? AND id_two = ?`
)
