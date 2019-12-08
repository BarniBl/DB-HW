package input

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
