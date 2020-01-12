package forum

import "time"

type User struct {
	About    string `json:"about"`
	Email    string `json:"email"`
	FullName string `json:"fullname"`
	NickName string `json:"nickname"`
}

type Forum struct {
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	User    string `json:"user"`
	Posts   int    `json:"posts"`
	Threads int    `json:"threads"`
}

type Thread struct {
	Author  string    `json:"author"`
	Created time.Time `json:"created"`
	Forum   string    `json:"forum"`
	Id      int       `json:"id"`
	Message string    `json:"message"`
	Slug    string    `json:"string"`
	Title   string    `json:"title"`
	Votes   int       `json:"votes"`
}

type Post struct {
	Author   string    `json:"author"`
	Created  time.Time `json:"created"`
	Forum    string    `json:"forum"`
	Id       int       `json:"id"`
	IsEdited bool      `json:"isEdited"`
	Message  string    `json:"message"`
	Parent   int       `json:"parent"`
	Thread   int       `json:"thread"`
}

type Message struct {
	Message string `json:"message"`
}


type ErrorMessage struct {
	Message string `json:"message"`
}

type FullPost struct {
	Author User   `json:"author"`
	Forum  Forum  `json:"forum"`
	Post   Post   `json:"post"`
	Thread Thread `json:"thread"`
}

type Status struct {
	Post int `json:"post"`
	Thread int `json:"thread"`
	User	int `json:"user"`
	Forum int `json:"forum"`
}
