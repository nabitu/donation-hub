package model

import (
	"errors"
	"github.com/isdzulqor/donation-hub/internal/utill/validator"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	AppPort                 string
	DBDriverName            string
	DBDataSource            string
	AwsDefaultRegion        string
	AwsAccessKey            string
	AwsSecretAccessKey      string
	AwsEndpoint             string
	AwsUsePathStyleEndpoint bool
	AwsS3Bucket             string
	TokenSecretKey          string
	TokenIssuer             string
}

type Connection struct {
	S3Client *s3.Client
	DB       *sqlx.DB
}

type Container struct {
	Connection *Connection
	Config     *Config
}

type UserRegisterInput struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=20"`
	Role     string `json:"role" validate:"required"`
}

func (u UserRegisterInput) Validate() error {
	return validator.Validate().Struct(u)
}

type UserRegisterOutput struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserLoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginOutput struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
}

type ListUserInput struct {
	Limit int64  `json:"limit"`
	Page  int64  `json:"page"`
	Role  string `json:"role"`
}

// UserStorage raw data user from database
type UserStorage struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"` // for api, please don't query password
	Roles    string `json:"roles"`
}

type User struct {
	ID       int64    `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Password string   `json:"password,omitempty"` // for api, please don't query password
	Roles    []string `json:"roles"`
}

type ListUserRole struct {
	Role string `json:"role"`
}

type ListUserMeta struct {
	Page       int64 `json:"page"`
	TotalPages int64 `json:"total_pages"`
}

type ListUserOutput struct {
	Users      []User       `json:"users"`
	Pagination ListUserMeta `json:"pagination"`
}

type RequestUploadUrlInput struct {
	UserID   int64  `json:"user_id"` // user auth id from jwt or other
	MimeType string `json:"mime_type"`
	FileSize int64  `json:"file_size"`
}

type RequestUploadUrlStorage struct {
	Url       string `json:"url"`
	ExpiresAt int64  `json:"expires_at"`
}

type RequestUploadUrlOutput struct {
	MimeType  string `json:"mime_type"`
	FileSize  int64  `json:"file_size"`
	URL       string `json:"url"`
	ExpiresAt int64  `json:"expires_at"`
}

type SubmitProjectInput struct {
	UserID       int64    `json:"user_id"` // user auth id from jwt or other
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	ImageURLs    []string `json:"image_urls"`
	DueAt        int64    `json:"due_at"`
	TargetAmount int64    `json:"target_amount"`
	Currency     string   `json:"currency"`
}

type SubmitProjectOutput struct {
	ID           int64    `json:"id"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	ImageURLs    []string `json:"image_urls"`
	DueAt        int64    `json:"due_at"`
	TargetAmount int64    `json:"target_amount"`
	Currency     string   `json:"currency"`
}

type ReviewProjectByAdminInput struct {
	UserID    int64  `json:"user_id"` // user auth id from jwt or other
	ProjectId int64  `json:"project_id"`
	Status    string `json:"status"`
}

type GetProjectByIdInput struct {
	ProjectId int64 `json:"project_id"`
}

type ListProjectInput struct {
	UserID  int64  `json:"user_id"` // user auth id from jwt or other
	Status  string `json:"status"`
	Limit   int64  `json:"limit"`
	StartTs int64  `json:"start_ts"` // jangan lupa, ini nanti Unix timestamp
	EndTs   int64  `json:"end_ts"`   // jangan lupa, ini nanti Unix timestamp
	LastKey string `json:"last_key"`
}

type Requester struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Project struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	ImageURLs    []string  `json:"image_urls"`
	DueAt        int64     `json:"due_at"`
	TargetAmount float64   `json:"target_amount"`
	Currency     string    `json:"currency"`
	Status       string    `json:"status"`
	Requester    Requester `json:"requester"`
}

type GetProjectByIdOutput struct {
	ID               int64     `json:"id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	ImageURLs        []string  `json:"image_urls"`
	DueAt            int64     `json:"due_at"`
	TargetAmount     float64   `json:"target_amount"`
	CollectionAmount float64   `json:"collection_amount"`
	Currency         string    `json:"currency"`
	Status           string    `json:"status"`
	Requester        Requester `json:"requester"`
}

type ListProjectOutput struct {
	Projects []Project `json:"projects"`
	LastKey  string    `json:"last_key"`
}

type DonateToProjectInput struct {
	UserID    int64  `json:"user_id"` // user auth id from jwt or other
	ProjectId int64  `json:"project_id"`
	Amount    int64  `json:"amount"`
	Currency  string `json:"currency"`
	Message   string `json:"message"`
}

func (d DonateToProjectInput) Validate() error {
	if d.Amount <= 0 {
		return errors.New("invalid amount")
	}

	if d.Currency == "" {
		return errors.New("currency is required")
	}

	return nil
}

type ListProjectDonationInput struct {
	ProjectId int64  `json:"project_id"`
	Limit     int64  `json:"limit"`
	LastKey   string `json:"last_key"`
}

type Donor struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type Donation struct {
	ID            int64  `json:"id"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
	Message       string `json:"message"`
	Donor         Donor  `json:"donor"`
	CreatedAt     int64  `json:"created_at" db:"created_at"`
	DonorId       int64  `json:"donor_id" db:"donor_id"`
	DonorUsername string `json:"donor_username" db:"donor_username"`
}

type ListProjectDonationOutput struct {
	Donations []Donation `json:"donations"`
	LastKey   string     `json:"last_key"`
}

type AuthPayload struct {
	UserID   int64
	Username string
	Email    string
	Role     []string
}
