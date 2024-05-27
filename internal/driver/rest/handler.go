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

	users, err := h.UserService.ListUser(r.Context(), req)
	if err != nil {
		fmt.Println("error get users", err)
		ResponseErrorBadRequest(w, err.Error())
		return
	}

	ResponseSuccess(w, *users)
}

func (h *Handler) HandleRequestProjectUrl(w http.ResponseWriter, r *http.Request) {
	req := model.RequestUploadUrlInput{
		UserID:   r.Context().Value("auth_id").(int64),
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
	req.UserID = 3

	response, err := h.ProjectService.RequestUploadUrl(r.Context(), req)
	if err != nil {
		fmt.Println("error request upload url", err)
		ResponseErrorBadRequest(w, err.Error())
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

	req.UserID = r.Context().Value("auth_id").(int64)
	response, err := h.ProjectService.SubmitProject(r.Context(), req)
	if err != nil {
		fmt.Println("error submit project", err)
		ResponseErrorBadRequest(w, err.Error())
		return
	}

	ResponseSuccess(w, response)
}

func (h *Handler) HandleProjectReview(w http.ResponseWriter, r *http.Request) {
	var req model.ReviewProjectByAdminInput
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println("error decoding request", err)
		ResponseErrorBadRequest(w, "invalid request")
		return
	}

	req.UserID = r.Context().Value("auth_id").(int64)
	err = h.ProjectService.ReviewProjectByAdmin(r.Context(), req)
	if err != nil {
		return
	}

	ResponseSuccess(w, "")
}

func (h *Handler) HandleProjects(w http.ResponseWriter, r *http.Request) {
	var req model.ListProjectInput
	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		limit = 10
	}
	req.Limit = limit

	if r.Context().Value("auth_id") != "" {
		req.UserID = r.Context().Value("auth_id").(int64)
	}

	projects, err := h.ProjectService.ListProject(r.Context(), req)
	if err != nil {
		fmt.Println("error get projects", err)
		ResponseErrorBadRequest(w, err.Error())
		return
	}

	ResponseSuccess(w, projects)
}

func (h *Handler) HandleProjectDetail(w http.ResponseWriter, r *http.Request) {
	var req model.GetProjectByIdInput
	projectIDStr := r.URL.Query().Get("id")
	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		fmt.Println("error parsing project ID", err)
		ResponseErrorBadRequest(w, "invalid project ID")
		return
	}
	req.ProjectId = projectID

	// Call the project service to get the project detail
	project, err := h.ProjectService.GetProjectById(r.Context(), req)
	if err != nil {
		fmt.Println("error getting project detail", err)
		ResponseErrorBadRequest(w, err.Error())
		return
	}

	ResponseSuccess(w, project)
}

func (h *Handler) HandleDonateProject(w http.ResponseWriter, r *http.Request) {
	// TODO: implement donate to project
	var req model.DonateToProjectInput
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println("error decoding request", err)
		ResponseErrorBadRequest(w, "invalid request")
		return
	}

	// Call the project service to donate to the project
	req.UserID = r.Context().Value("auth_id").(int64)
	err = h.ProjectService.DonateToProject(r.Context(), req)
	if err != nil {
		fmt.Println("error donating to project", err)
		ResponseErrorBadRequest(w, err.Error())
		return
	}

	ResponseSuccess(w, "Donation successful")
}

func (h *Handler) HandleProjectDonations(w http.ResponseWriter, r *http.Request) {
	var req model.ListProjectDonationInput
	req.ProjectId, _ = strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	req.LastKey = r.URL.Query().Get("last_key")
	req.Limit, _ = strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)

	// Call the project service to get the project donations
	donations, err := h.ProjectService.ListDonationByProjectId(r.Context(), req)
	if err != nil {
		fmt.Println("error getting project donations", err)
		ResponseErrorBadRequest(w, err.Error())
		return
	}

	ResponseSuccess(w, donations)
}
