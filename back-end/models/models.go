package models

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Password    string `json:"password"`
	Avatar      string `json:"avatar"`
	DateOfBirth string `json:"birthday"`
	Gender      string `json:"gender"`
	About_Me    string `json:"aboutme"`
	Country     string `json:"country"`
	Privacy     string `json:"privacy"`
	Created_at  string `json:"createad_at"`
	/*Followers   []Followers `json:"followers"`
	Following   []Followers `json:"following"`*/
}

type Followers struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	First_user  int    `json:"first_user"`
	Second_user int    `json:"second_user"`
	Status      string `json:"status"`
}

// new struct because followers struct is used to send follow request but table has no FirstName & LastName, using the struct below to send data to front-end.
type FollowersInfo struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	First_user  int    `json:"first_user"`
	Second_user int    `json:"second_user"`
	Status      string `json:"status"`
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
}

type Chatroom struct {
	ID          uint `json:"id" gorm:"primary_key"`
	First_user  int  `json:"first_user"`
	Chatroom_id int  `json:"chatroom_id"`
}

type Message struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Message     string `json:"message"`
	User_id     int    `json:"user_id"`
	Chatroom_id int    `json:"chatroom_id"`
	Timestamp   string `json:"timestamp"`
}

type Post struct {
	ID         uint   `json:"id" gorm:"primary_key"`
	Header     string `json:"header"`
	Text       string `json:"text"`
	User_id    int    `json:"user_id"`
	Created_at string `json:"created_at"`
	Image      string `json:"image"`
	Privacy    int    `json:"privacy"`
}

type Grouppost struct {
	ID         uint   `json:"id" gorm:"primary_key"`
	Header     string `json:"header"`
	Text       string `json:"text"`
	User_id    int    `json:"user_id"`
	Created_at string `json:"created_at"`
	Image      string `json:"image"`
	Group_id   int    `json:"group_id"`
}

type Comment struct {
	ID         uint   `json:"id" gorm:"primary_key"`
	Text       string `json:"text"`
	User_id    int    `json:"user_id"`
	Post_id    int    `json:"post_id"`
	Created_at string `json:"created_at"`
	Image      string `json:"image"`
}

type Group struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Creator     int    `json:"creator"`
	Created_at  string `json:"created_at"`
	Avatar      string `json:"avatar"`
}

type InitialData struct {
	UserData  User       `json:"userData"`
	Messages  []Message  `json:"messages"`
	Chatrooms []Chatroom `json:"chatrooms"`
}

type MessageRequest struct {
	Chatroom_id int `json:"chatroomID"`
	MessageID   int `json:"messageID"`
	Amount      int `json:"amount"`
}

type Submit struct {
	Type           string `json:"type"`
	Creator_id     int    `json:"creator_id"`
	Target_user_id int    `json:"target_user_id"`
	Header         string `json:"header"`
	Content        string `json:"content"`
	Group_id       int    `json:"group_id"`
	Master_post    int    `json:"master_post"`
}

type CustomPrivacy struct {
	ID      uint `json:"id" gorm:"primary_key"`
	Post_id int  `json:"post_id"`
	User_id int  `json:"user_id"`
}

type Feed struct {
	PostsAmount string `json:"post_amount"`
	PostID      string `json:"post_ID"`
}

type GroupRequest struct {
	ID           uint   `json:"id" gorm:"primary_key"`
	Group_id     int    `json:"group_id"`
	User_id      int    `json:"user_id"`
	Pending_type string `json:"pending_type"`
	Status       string `json:"status"`
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Id      int    `json:"id"`
}
