package rest

import (
	"fmt"
	"net/http"

	"github.com/isdzulqor/donation-hub/internal/core/service/project"
	"github.com/isdzulqor/donation-hub/internal/core/service/user"
)

type Config struct {
	Port           string
	ProjectService project.Service
	UserService    user.Service
}

func (c Config) Validate() error {
	if c.Port == "" {
		return fmt.Errorf("port is required")
	}

	if c.ProjectService == nil {
		return fmt.Errorf("project service is required")
	}

	if c.UserService == nil {
		return fmt.Errorf("user service is required")
	}

	return nil
}

func StartApp(c Config) error {
	if err := c.Validate(); err != nil {
		return err
	}

	app := http.NewServeMux()

	handler := NewHandler(c.ProjectService, c.UserService)

	app.HandleFunc("GET /", handler.DefaultHandler)
	app.HandleFunc("POST /users/register", handler.HandleRegister)
	app.HandleFunc("POST /users/login", handler.HandleLogin)
	app.HandleFunc("GET /users", handler.HandleUsers)
	app.HandleFunc("GET /projects/upload", handler.HandleRequestProjectUrl)
	app.HandleFunc("POST /projects", handler.HandleSubmitProject)
	app.HandleFunc("PUT /projects/{id}/review", handler.HandleProjectReview)
	app.HandleFunc("GET /projects", handler.HandleProjects)
	app.HandleFunc("GET /projects/{id}", handler.HandleProjectDetail)
	app.HandleFunc("POST /projects/{id}/donations", handler.HandleDonateProject)
	app.HandleFunc("GET /projects/{id}/donations", handler.HandleProjectDonations)

	fmt.Println("Starting app on port", c.Port)
	http.ListenAndServe(":"+c.Port, app)

	return nil
}