package request

type IssueRequestTo struct {
	UserID  int64  `json:"userId" validate:"required"` // ID автора
	Title   string `json:"title" validate:"required,min=2,max=64"`
	Content string `json:"content" validate:"required,min=4,max=2048"`
}
