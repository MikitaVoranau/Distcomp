package service

import (
	"Voronov/internal/errors"
	"Voronov/internal/model"
	"Voronov/internal/repository"
	"Voronov/internal/transport/dto/request"
	"Voronov/internal/transport/dto/response"
	"time"
)

type UserServiceImpl struct {
	repo   repository.CRUDRepository[model.User]
	mapper Mapper
}

func NewUserService(repo repository.CRUDRepository[model.User], mapper Mapper) UserService {
	return &UserServiceImpl{repo: repo, mapper: mapper}
}

func (s *UserServiceImpl) FindByID(id int64) (*response.UserResponseTo, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToUserResponse(user), nil
}

func (s *UserServiceImpl) FindAll() ([]*response.UserResponseTo, error) {
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	result := make([]*response.UserResponseTo, 0, len(users))
	for _, u := range users {
		result = append(result, s.mapper.ToUserResponse(u))
	}
	return result, nil
}

func (s *UserServiceImpl) Create(req *request.UserRequestTo) (*response.UserResponseTo, error) {
	if req.Login == "" || req.Password == "" || req.Firstname == "" || req.Lastname == "" {
		return nil, errors.ErrBadRequest
	}
	user := s.mapper.ToUserModel(req)
	created, err := s.repo.Create(user)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToUserResponse(created), nil
}

func (s *UserServiceImpl) Update(id int64, req *request.UserRequestTo) (*response.UserResponseTo, error) {
	if req.Login == "" || req.Password == "" || req.Firstname == "" || req.Lastname == "" {
		return nil, errors.ErrBadRequest
	}
	user := s.mapper.ToUserModel(req)
	user.ID = id
	updated, err := s.repo.Update(id, user)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToUserResponse(updated), nil
}

func (s *UserServiceImpl) Delete(id int64) error {
	return s.repo.Delete(id)
}

type IssueServiceImpl struct {
	issueRepo      repository.CRUDRepository[model.Issue]
	userRepo       repository.CRUDRepository[model.User]
	labelRepo      repository.CRUDRepository[model.Label]
	reactionRepo   repository.CRUDRepository[model.Reaction]
	issueLabelRepo repository.CRUDRepository[model.IssueLabel]
	mapper         Mapper
}

func NewIssueService(
	issueRepo repository.CRUDRepository[model.Issue],
	userRepo repository.CRUDRepository[model.User],
	labelRepo repository.CRUDRepository[model.Label],
	reactionRepo repository.CRUDRepository[model.Reaction],
	issueLabelRepo repository.CRUDRepository[model.IssueLabel],
	mapper Mapper,
) IssueService {
	return &IssueServiceImpl{
		issueRepo:      issueRepo,
		userRepo:       userRepo,
		labelRepo:      labelRepo,
		reactionRepo:   reactionRepo,
		issueLabelRepo: issueLabelRepo,
		mapper:         mapper,
	}
}

func (s *IssueServiceImpl) FindByID(id int64) (*response.IssueResponseTo, error) {
	issue, err := s.issueRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user, err := s.userRepo.FindByID(issue.UserID); err == nil {
		issue.User = user
	}
	return s.mapper.ToIssueResponse(issue), nil
}

func (s *IssueServiceImpl) FindAll() ([]*response.IssueResponseTo, error) {
	issues, err := s.issueRepo.FindAll()
	if err != nil {
		return nil, err
	}
	result := make([]*response.IssueResponseTo, 0, len(issues))
	for _, i := range issues {
		if user, err := s.userRepo.FindByID(i.UserID); err == nil {
			i.User = user
		}
		result = append(result, s.mapper.ToIssueResponse(i))
	}
	return result, nil
}

func (s *IssueServiceImpl) Create(req *request.IssueRequestTo) (*response.IssueResponseTo, error) {
	if req.UserID == 0 || req.Title == "" || req.Content == "" {
		return nil, errors.ErrBadRequest
	}
	if _, err := s.userRepo.FindByID(req.UserID); err != nil {
		return nil, errors.ErrBadRequest
	}
	issue := s.mapper.ToIssueModel(req)
	issue.Created = time.Now()
	issue.Modified = time.Now()
	created, err := s.issueRepo.Create(issue)
	if err != nil {
		return nil, err
	}
	if user, err := s.userRepo.FindByID(created.UserID); err == nil {
		created.User = user
	}
	return s.mapper.ToIssueResponse(created), nil
}

func (s *IssueServiceImpl) Update(id int64, req *request.IssueRequestTo) (*response.IssueResponseTo, error) {
	if req.UserID == 0 || req.Title == "" || req.Content == "" {
		return nil, errors.ErrBadRequest
	}
	issue := s.mapper.ToIssueModel(req)
	issue.ID = id
	issue.Modified = time.Now()
	updated, err := s.issueRepo.Update(id, issue)
	if err != nil {
		return nil, err
	}
	if user, err := s.userRepo.FindByID(updated.UserID); err == nil {
		updated.User = user
	}
	return s.mapper.ToIssueResponse(updated), nil
}

func (s *IssueServiceImpl) Delete(id int64) error {
	return s.issueRepo.Delete(id)
}

func (s *IssueServiceImpl) FindByUserID(userID int64) (*response.UserResponseTo, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToUserResponse(user), nil
}

func (s *IssueServiceImpl) FindByIssueID(issueID int64) ([]*response.LabelResponseTo, []*response.ReactionResponseTo, error) {
	if _, err := s.issueRepo.FindByID(issueID); err != nil {
		return nil, nil, err
	}
	labels, err := s.labelRepo.FindAll()
	if err != nil {
		return nil, nil, err
	}
	labelResults := make([]*response.LabelResponseTo, 0)
	for _, l := range labels {
		labelResults = append(labelResults, s.mapper.ToLabelResponse(l))
	}
	reactions, err := s.reactionRepo.FindAll()
	if err != nil {
		return nil, nil, err
	}
	reactionResults := make([]*response.ReactionResponseTo, 0)
	for _, r := range reactions {
		if r.IssueID == issueID {
			reactionResults = append(reactionResults, s.mapper.ToReactionResponse(r))
		}
	}
	return labelResults, reactionResults, nil
}

func (s *IssueServiceImpl) SearchIssues(labelNames []string, labelIDs []int64, userLogin, title, content string) ([]*response.IssueResponseTo, error) {
	issues, err := s.issueRepo.FindAll()
	if err != nil {
		return nil, err
	}
	results := make([]*response.IssueResponseTo, 0)
	for _, issue := range issues {
		if title != "" && issue.Title != title {
			continue
		}
		if content != "" && issue.Content != content {
			continue
		}
		if userLogin != "" {
			user, err := s.userRepo.FindByID(issue.UserID)
			if err != nil || user.Login != userLogin {
				continue
			}
		}
		if user, err := s.userRepo.FindByID(issue.UserID); err == nil {
			issue.User = user
		}
		results = append(results, s.mapper.ToIssueResponse(issue))
	}
	return results, nil
}

type LabelServiceImpl struct {
	repo   repository.CRUDRepository[model.Label]
	mapper Mapper
}

func NewLabelService(repo repository.CRUDRepository[model.Label], mapper Mapper) LabelService {
	return &LabelServiceImpl{repo: repo, mapper: mapper}
}

func (s *LabelServiceImpl) FindByID(id int64) (*response.LabelResponseTo, error) {
	label, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToLabelResponse(label), nil
}

func (s *LabelServiceImpl) FindAll() ([]*response.LabelResponseTo, error) {
	labels, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	result := make([]*response.LabelResponseTo, 0, len(labels))
	for _, l := range labels {
		result = append(result, s.mapper.ToLabelResponse(l))
	}
	return result, nil
}

func (s *LabelServiceImpl) Create(req *request.LabelRequestTo) (*response.LabelResponseTo, error) {
	if req.Name == "" {
		return nil, errors.ErrBadRequest
	}
	label := s.mapper.ToLabelModel(req)
	created, err := s.repo.Create(label)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToLabelResponse(created), nil
}

func (s *LabelServiceImpl) Update(id int64, req *request.LabelRequestTo) (*response.LabelResponseTo, error) {
	if req.Name == "" {
		return nil, errors.ErrBadRequest
	}
	label := s.mapper.ToLabelModel(req)
	label.ID = id
	updated, err := s.repo.Update(id, label)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToLabelResponse(updated), nil
}

func (s *LabelServiceImpl) Delete(id int64) error {
	return s.repo.Delete(id)
}

type ReactionServiceImpl struct {
	repo   repository.CRUDRepository[model.Reaction]
	mapper Mapper
}

func NewReactionService(repo repository.CRUDRepository[model.Reaction], mapper Mapper) ReactionService {
	return &ReactionServiceImpl{repo: repo, mapper: mapper}
}

func (s *ReactionServiceImpl) FindByID(id int64) (*response.ReactionResponseTo, error) {
	reaction, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToReactionResponse(reaction), nil
}

func (s *ReactionServiceImpl) FindAll() ([]*response.ReactionResponseTo, error) {
	reactions, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	result := make([]*response.ReactionResponseTo, 0, len(reactions))
	for _, r := range reactions {
		result = append(result, s.mapper.ToReactionResponse(r))
	}
	return result, nil
}

func (s *ReactionServiceImpl) Create(req *request.ReactionRequestTo) (*response.ReactionResponseTo, error) {
	if req.IssueID == 0 || req.Content == "" {
		return nil, errors.ErrBadRequest
	}
	reaction := s.mapper.ToReactionModel(req)
	created, err := s.repo.Create(reaction)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToReactionResponse(created), nil
}

func (s *ReactionServiceImpl) Update(id int64, req *request.ReactionRequestTo) (*response.ReactionResponseTo, error) {
	if req.IssueID == 0 || req.Content == "" {
		return nil, errors.ErrBadRequest
	}
	reaction := s.mapper.ToReactionModel(req)
	reaction.ID = id
	updated, err := s.repo.Update(id, reaction)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToReactionResponse(updated), nil
}

func (s *ReactionServiceImpl) Delete(id int64) error {
	return s.repo.Delete(id)
}
