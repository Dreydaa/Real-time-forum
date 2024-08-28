package structure

type Post struct {
	ID       int    `json:"id"`
	User_id  int    `json:"user_id"`
	Title    string `json:"title"`
	Category string `json:"category"`
	Date     string `json:"date"`
	Content  string `json:"content"`
}

type Comment struct {
	ID      int    `json:"id"`
	Post_id int    `json:"post_id"`
	User_id int    `json:"user_id"`
	Comment string `json:"comment"`
	Date    string `json:"date"`
	Content string `json:"content"`
}

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Gender    string `json:"gender"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Message struct {
	ID          int    `json:"id"`
	Sender_id   int    `json:"sender_id"`
	Receiver_id int    `json:"receiver_id"`
	Content     string `json:"content"`
	Date        string `json:"date"`
	Msg_type    string `json:"msg_type"`
	UserID      int    `json:"user_id"`
}

type Login struct {
	Data     string `json:"emailUsername"`
	Password string `json:"password"`
}

type Chat struct {
	User_one int
	User_two int
	Time     int
}

type OnlineUsers struct {
	UserIds  []int  `json:"user_ids"`
	Msg_type string `json:"msg_type"`
}

type Resp struct {
	Msg string `json:"msg"`
}

type Session struct {
	Session_uuid string
	User_id      int
}
