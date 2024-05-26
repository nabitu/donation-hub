package main

import (
	"fmt"

	"github.com/isdzulqor/donation-hub/internal/core/service/project"
	"github.com/isdzulqor/donation-hub/internal/core/service/user"
	"github.com/isdzulqor/donation-hub/internal/driven/auth/jwt"
	"github.com/isdzulqor/donation-hub/internal/driven/storage/mysql/projectstorage"
	"github.com/isdzulqor/donation-hub/internal/driven/storage/mysql/userstorage"
	"github.com/isdzulqor/donation-hub/internal/driven/storage/s3/projectfilestorage"
	"github.com/isdzulqor/donation-hub/internal/driver/rest"
)

func main() {
	c := InitContainer()
	userStorage := userstorage.New(c)
	projectStorage := projectstorage.New(c)
	projectFileStorage := projectfilestorage.New(c)
	jwt := jwt.New(c)

	userService := user.New(userStorage, jwt)
	projectService := project.New(projectStorage, projectFileStorage, userStorage)

	err := rest.StartApp(rest.Config{
		Port:           c.Config.AppPort,
		ProjectService: projectService,
		UserService:    userService,
	})
	if err != nil {
		fmt.Errorf("failed to start app, err: %s", err.Error())
	}
}
