package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/isdzulqor/donation-hub/internal/core/model"
	"github.com/isdzulqor/donation-hub/internal/core/service/project"
	"github.com/isdzulqor/donation-hub/internal/core/service/user"
)

func NewHandler(projectService project.Service, userService user.Service) *Handler {
	return &Handler{
		ProjectService: projectService,
		UserService:    userService,
	}
}

type Handler struct {
	ProjectService project.Service
	UserService    user.Service
}

func (h *Handler) DefaultHandler(w http.ResponseWriter, r *http.Request) {
	ResponseSuccess(w, "Donation Hub App")
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {

	var req model.UserLoginInput
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println("error decoding request", err)
		ResponseErrorBadRequest(w, "invalid request")
		return
	}

	response, err := h.UserService.Login(r.Context(), req)
	if err != nil {
		fmt.Println("error login", err)
		ResponseErrorInvalidCredentials(w, err.Error())
		return
	}

	ResponseSuccess(w, response)
}

func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req model.UserRegisterInput
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println("error decoding request", err)
		ResponseErrorBadRequest(w, "invalid request")
		return
	}

	response, err := h.UserService.Register(r.Context(), req)
	if err != nil {
		fmt.Println("error register", err)
		ResponseErrorBadRequest(w, err.Error())
		return
	}

	ResponseSuccess(w, response)
}

func (h *Handler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	var req model.ListUserInput
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1
	}
	req.Page = page

	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		limit = 10
	}
	req.Limit = limit

	req.Role = r.URL.Query().Get("role")
	req.Role = "donor"

	users, err := h.UserService.ListUser(r.Context(), req)
	if err != nil {
		fmt.Println("error get users", err)
		ResponseErrorBadRequest(w, err.Error())
		return
	}

	ResponseSuccess(w, users)
}

func (h *Handler) HandleRequestProjectUrl(w http.ResponseWriter, r *http.Request) {
	req := model.RequestUploadUrlInput{
		UserID:   1, // todo: get user id from context
		MimeType: r.URL.Query().Get("mime_type"),
	}

	fileSizeStr := r.URL.Query().Get("file_size")
	fileSize, err := strconv.ParseInt(fileSizeStr, 10, 64)
	if err != nil {
		fmt.Println("error parsing file size", err)
		ResponseErrorBadRequest(w, "invalid file size")
		return
	}

	req.FileSize = fileSize

	response, err := h.ProjectService.RequestUploadUrl(r.Context(), req)
	if err != nil {
		return
	}

	ResponseSuccess(w, response)
}

func (h *Handler) HandleSubmitProject(w http.ResponseWriter, r *http.Request) {
	var req model.SubmitProjectInput
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println("error decoding request", err)
		ResponseErrorBadRequest(w, "invalid request")
		return
	}

	req.UserID = 4 // todo: get user id from context

	response, err := h.ProjectService.SubmitProject(r.Context(), req)
	if err != nil {
		fmt.Println("error submit project", err)
		ResponseErrorBadRequest(w, err.Error())
		return
	}

	ResponseSuccess(w, response)
}

func (h *Handler) HandleProjectReview(w http.ResponseWriter, r *http.Request) {
	// TODO: implement project review
}

func (h *Handler) HandleProjects(w http.ResponseWriter, r *http.Request) {
	// TODO: implement get projects
}

func (h *Handler) HandleProjectDetail(w http.ResponseWriter, r *http.Request) {
	// TODO: implement get project detail
}

func (h *Handler) HandleDonateProject(w http.ResponseWriter, r *http.Request) {
	// TODO: implement donate to project
}

func (h *Handler) HandleProjectDonations(w http.ResponseWriter, r *http.Request) {
	// TODO: implement get project donations
}
