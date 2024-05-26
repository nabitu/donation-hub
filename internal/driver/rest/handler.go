package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	// TODO: implement register
}

func (h *Handler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	// TODO: implement get users
}

func (h *Handler) HandleRequestProjectUrl(w http.ResponseWriter, r *http.Request) {
	// TODO: implement request project URL
}

func (h *Handler) HandleSubmitProject(w http.ResponseWriter, r *http.Request) {
	// TODO: implement submit project
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
