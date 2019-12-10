package output

import "github.com/BarniBl/DB-HW/internal/input"

type ErrorMessage struct {
	Message string `json:"message"`
}

type FullPost struct {
	Author input.User   `json:"author"`
	Forum  input.Forum  `json:"forum"`
	Post   input.Post   `json:"post"`
	Thread input.Thread `json:"thread"`
}
