package main

import (
	"github.com/ldej/go-rest-example/internal/api/user"
	"github.com/ldej/go-rest-example/pkg/encrypt"
	"github.com/ldej/go-rest-example/pkg/server"
	"github.com/ldej/go-rest-example/pkg/store/postgres"
	"github.com/ldej/go-rest-example/pkg/uid"
)

func main() {
	s, err := server.New()
	if err != nil {
		logger.Fatalf("whelp %v", err)
		return
	}

	db, err := postgres.New()
	if err != nil {
		logger.Fatalf("whelp %v", err)
		return
	}
	uidGenerator := uid.NewGenerator()
	encryptor := encrypt.Encryptor{}

	// Hook up the services
	userService := user.NewService(db, uidGenerator, encryptor)
	userRouter := s.Router().PathPrefix("/users").Subrouter()
	user.Setup(userRouter, userService)

	err = s.ListenAndServe()
	if err != nil {
		logger.Fatalf("whelp %v", err)
		return
	}
}
