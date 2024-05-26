package main

import (
	"github.com/isdzulqor/donation-hub/internal/driven/auth/jwt"
	"github.com/isdzulqor/donation-hub/internal/driven/storage/mysql/projectstorage"
	"github.com/isdzulqor/donation-hub/internal/driven/storage/mysql/userstorage"
	"github.com/isdzulqor/donation-hub/internal/driven/storage/s3/projectfilestorage"
	"log"
	"net/http"
)

func main() {
	c := InitContainer()
	_ = userstorage.New(c)
	_ = projectstorage.New(c)
	_ = projectfilestorage.New(c)
	_ = jwt.New(c)

	log.Println("Starting server on :8180")
	http.ListenAndServe(":8180", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	}))
}
