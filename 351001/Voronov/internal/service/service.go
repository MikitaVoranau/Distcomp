package service

import (
	"Voronov/internal/model"
	"Voronov/internal/transport/dto/request"
	"Voronov/internal/transport/dto/response"
)

type UserService interface {
	FindByID(id int64) (*response.UserResponseTo, error)
	FindAll() ([]*response.UserResponseTo, error)
	Create(req *request.UserRequestTo) (*response.UserResponseTo, error)
	Update(id int64, req *request.UserRequestTo) (*response.UserResponseTo, error)
	Delete(id int64) error
}

type IssueService interface {
	FindByID(id int64) (*response.IssueResponseTo, error)
	FindAll() ([]*response.IssueResponseTo, error)
	Create(req *request.IssueRequestTo) (*response.IssueResponseTo, error)
	Update(id int64, req *request.IssueRequestTo) (*response.IssueResponseTo, error)
	Delete(id int64) error
	FindByUserID(userID int64) (*response.UserResponseTo, error)
	FindByIssueID(issueID int64) ([]*response.LabelResponseTo, []*response.ReactionResponseTo, error)
	SearchIssues(labelNames []string, labelIDs []int64, userLogin, title, content string) ([]*response.IssueResponseTo, error)
}

type LabelService interface {
	FindByID(id int64) (*response.LabelResponseTo, error)
	FindAll() ([]*response.LabelResponseTo, error)
	Create(req *request.LabelRequestTo) (*response.LabelResponseTo, error)
	Update(id int64, req *request.LabelRequestTo) (*response.LabelResponseTo, error)
	Delete(id int64) error
}

type ReactionService interface {
	FindByID(id int64) (*response.ReactionResponseTo, error)
	FindAll() ([]*response.ReactionResponseTo, error)
	Create(req *request.ReactionRequestTo) (*response.ReactionResponseTo, error)
	Update(id int64, req *request.ReactionRequestTo) (*response.ReactionResponseTo, error)
	Delete(id int64) error
}

type Mapper interface {
	ToUserResponse(m *model.User) *response.UserResponseTo
	ToUserModel(req *request.UserRequestTo) *model.User
	ToIssueResponse(m *model.Issue) *response.IssueResponseTo
	ToIssueModel(req *request.IssueRequestTo) *model.Issue
	ToLabelResponse(m *model.Label) *response.LabelResponseTo
	ToLabelModel(req *request.LabelRequestTo) *model.Label
	ToReactionResponse(m *model.Reaction) *response.ReactionResponseTo
	ToReactionModel(req *request.ReactionRequestTo) *model.Reaction
}
