package project

import (
	"context"
	"errors"
	"fmt"
	"github.com/isdzulqor/donation-hub/internal/utill/role"
	"log"

	"github.com/isdzulqor/donation-hub/internal/core/model"
	"github.com/isdzulqor/donation-hub/internal/core/service/user"
	_type "github.com/isdzulqor/donation-hub/internal/core/type"
)

type Storage struct {
	storage         DataStorage
	fileStorage     FileStorage
	userDataStorage user.DataStorage
}

type Service interface {
	RequestUploadUrl(context.Context, model.RequestUploadUrlInput) (*model.RequestUploadUrlOutput, error)
	SubmitProject(context.Context, model.SubmitProjectInput) (*model.SubmitProjectOutput, error)
	ReviewProjectByAdmin(context.Context, model.ReviewProjectByAdminInput) error
	ListProject(context.Context, model.ListProjectInput) (*model.ListProjectOutput, error)
	GetProjectById(context.Context, model.GetProjectByIdInput) (*model.GetProjectByIdOutput, error)
	DonateToProject(context.Context, model.DonateToProjectInput) error
	ListDonationByProjectId(context.Context, model.ListProjectDonationInput) (*model.ListProjectDonationOutput, error)
}

func New(storage DataStorage, fileStorage FileStorage, userDataStorage user.DataStorage) Service {
	return &Storage{
		storage:         storage,
		fileStorage:     fileStorage,
		userDataStorage: userDataStorage,
	}
}

func (s *Storage) RequestUploadUrl(ctx context.Context, input model.RequestUploadUrlInput) (*model.RequestUploadUrlOutput, error) {
	// validate user, make sure role is valid
	ok, err := s.userDataStorage.UserHasRole(ctx, input.UserID, _type.ROLE_REQUESTER)
	if !ok || err != nil {
		return nil, errors.New("ERR_FORBIDDEN_ACCESS")
	}

	// validate size
	if input.FileSize > 1048576 {
		return nil, errors.New("filesize can't greater than 1MB")
	}

	if input.FileSize <= 0 {
		return nil, errors.New("filesize must greater than 0Kb")
	}

	// validate mimetype
	if input.MimeType != "image/jpeg" && input.MimeType != "image/png" {
		return nil, errors.New("mimetype must be image/jpeg or image/png")
	}

	r, err := s.fileStorage.RequestUploadUrl(input.MimeType, input.FileSize)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to request upload url, err: %s", err.Error()))
	}

	return &model.RequestUploadUrlOutput{
		MimeType:  input.MimeType,
		FileSize:  input.FileSize,
		URL:       r.Url,
		ExpiresAt: r.ExpiresAt,
	}, nil
}

func (s *Storage) SubmitProject(ctx context.Context, input model.SubmitProjectInput) (*model.SubmitProjectOutput, error) {
	input.Validate()
	// validate user, make sure role is valid
	ok, err := s.userDataStorage.UserHasRole(ctx, input.UserID, _type.ROLE_REQUESTER)
	if !ok || err != nil {
		return nil, errors.New("ERR_FORBIDDEN_ACCESS")
	}

	ok, err = s.storage.HasName(ctx, input.Title)
	if ok || err != nil {
		log.Println(err)
		return nil, errors.New("project name already exist")
	}

	// save to database
	projectId, err := s.storage.Submit(ctx, input)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to submit project, err: %s", err.Error()))
	}

	return &model.SubmitProjectOutput{
		ID:           projectId,
		Title:        input.Title,
		Description:  input.Description,
		ImageURLs:    input.ImageURLs,
		DueAt:        input.DueAt,
		TargetAmount: input.TargetAmount,
		Currency:     input.Currency,
	}, nil
}

func (s *Storage) ReviewProjectByAdmin(ctx context.Context, input model.ReviewProjectByAdminInput) error {
	// validate user, make sure role is valid
	ok, err := s.userDataStorage.UserHasRole(ctx, input.UserID, _type.ROLE_ADMIN)
	if !ok || err != nil {
		return errors.New("ERR_FORBIDDEN_ACCESS")
	}

	if input.Status != _type.PROJECT_APPROVED && input.Status != _type.PROJECT_REJECTED {
		return errors.New("status must be approved or rejected")
	}

	// jika status project bukan need_review maka tidak bisa review
	p, err := s.storage.GetProjectById(ctx, model.GetProjectByIdInput{ProjectId: input.ProjectId})
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get project, err: %s", err.Error()))
	}
	if p.Status != _type.PROJECT_NEED_REVIEW {
		return errors.New("ERR_PROJECT_NOT_NEED_REVIEW")
	}

	err = s.storage.ReviewByAdmin(ctx, input)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to review project, err: %s", err.Error()))
	}

	return nil
}

func (s *Storage) ListProject(ctx context.Context, input model.ListProjectInput) (*model.ListProjectOutput, error) {
	input.IsAdmin = role.HasRoleFromContext(ctx, _type.ROLE_ADMIN)
	// make sure user has role admin if status need_review
	if input.Status == _type.PROJECT_NEED_REVIEW {
		if !input.IsAdmin {
			return nil, errors.New("ERR_FORBIDDEN_ACCESS")
		}
	}

	output, err := s.storage.ListProject(ctx, input)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to list project, err: %s", err.Error()))
	}

	return &model.ListProjectOutput{
		Projects: output.Projects,
		LastKey:  output.LastKey,
	}, nil
}

func (s *Storage) GetProjectById(ctx context.Context, input model.GetProjectByIdInput) (*model.GetProjectByIdOutput, error) {
	output, err := s.storage.GetProjectById(ctx, input)

	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (s *Storage) DonateToProject(ctx context.Context, input model.DonateToProjectInput) error {
	// validate donations
	if err := input.Validate(); err != nil {
		return errors.New(fmt.Sprintf("failed to validate donation, err: %s", err.Error()))
	}

	// make sure user has role donor
	ok, err := s.userDataStorage.UserHasRole(ctx, input.UserID, _type.ROLE_DONOR)
	if !ok || err != nil {
		return errors.New("ERR_FORBIDDEN_ACCESS")
	}

	p, err := s.storage.GetProjectById(ctx, model.GetProjectByIdInput{ProjectId: input.ProjectId})

	// jika status project masih need_review, maka tidak bisa donate
	if p.Status == _type.PROJECT_NEED_REVIEW {
		return errors.New("ERR_PROJECT_NEED_REVIEW")
	}

	if err != nil {
		return errors.New(fmt.Sprintf("failed to get project, err: %s", err.Error()))
	}

	if float64(input.Amount) > p.TargetAmount {
		return errors.New("ERR_TOO_MUCH_DONATION")
	}

	err = s.storage.DonateToProject(ctx, input)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to donate to project, err: %s", err.Error()))
	}

	return nil
}

func (s *Storage) ListDonationByProjectId(ctx context.Context, input model.ListProjectDonationInput) (*model.ListProjectDonationOutput, error) {

	p, err := s.storage.GetProjectById(ctx, model.GetProjectByIdInput{ProjectId: input.ProjectId})

	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to get project, err: %s", err.Error()))
	}

	// jika status project masih need_review, maka tidak bisa donate
	if p.Status == _type.PROJECT_NEED_REVIEW {
		return nil, errors.New("ERR_PROJECT_NEED_REVIEW")
	}

	output, err := s.storage.ListDonationByProjectId(ctx, input)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to list donation, err: %s", err.Error()))
	}

	return &output, nil
}
