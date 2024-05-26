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
	UserServicce   user.Service
}

func (c Config) Validate() error {
	if c.Port == "" {
		return fmt.Errorf("port is required")
	}

	if c.ProjectService == nil {
		return fmt.Errorf("project service is required")
	}

	if c.UserServicce == nil {
		return fmt.Errorf("user service is required")
	}

	return nil
}

func StartApp(c Config) error {
	if err := c.Validate(); err != nil {
		return err
	}

	app := http.NewServeMux()

	app.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(" (╯°□°）╯︵ ┻━┻ "))
	})

	app.HandleFunc("POST /users/login", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Login"))
	})

	fmt.Println("Starting app on port", c.Port)
	http.ListenAndServe(":"+c.Port, app)

	return nil
}
