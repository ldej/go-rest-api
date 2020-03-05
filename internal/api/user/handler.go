package user

import (
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/gorilla/mux"
	config "github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/ldej/go-rest-example/internal/api"
	"github.com/ldej/go-rest-example/pkg/jwt"
	"github.com/ldej/go-rest-example/pkg/middleware"
	"github.com/ldej/go-rest-example/pkg/server"
	"github.com/ldej/go-rest-example/pkg/store"
)

type Handler struct {
	logger      *zap.SugaredLogger
	router      *mux.Router
	jwtKey      []byte
	userService *Service
}

func Setup(router *mux.Router, userService *Service) {
	h := Handler{
		logger:      zap.S().With("package", "user"),
		router:      router,
		userService: userService,
	}
	h.router.HandleFunc("", h.Register).Methods(http.MethodPost)
	h.router.HandleFunc("/{uid}", middleware.RequireUser(h.GetUser)).Methods(http.MethodGet)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	c := r.Context()

	var registerUserRequest api.RegisterUserRequest
	if err := render.DecodeJSON(r.Body, &registerUserRequest); err != nil {
		render.Render(w, r, server.ErrInvalidRequest(err))
		return
	}

	user, err := h.userService.Register(c, registerUserRequest)
	if err != nil {
		render.Render(w, r, server.ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, api.UserResponse{
		UID:          user.UID,
		Name:         user.Name,
		EmailAddress: user.EmailAddress,
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	c := r.Context()

	var loginRequest api.LoginRequest
	err := render.Bind(r, &loginRequest)
	if err != nil {
		render.Render(w, r, server.ErrInvalidRequest(err))
		return
	}

	user, err := h.userService.Login(c, loginRequest)
	if errors.Is(err, store.ErrNotFound) {
		render.Render(w, r, server.ErrNotFound)
		return
	} else if err != nil {
		render.Render(w, r, server.ErrInvalidRequest(err))
		return
	}

	render.JSON(w, r, api.UserResponse{
		UID:          user.UID,
		Name:         user.Name,
		EmailAddress: user.EmailAddress,
	})
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request, authUser *api.AuthUser) {
	c := r.Context()

	userUID := mux.Vars(r)["uid"]
	h.logger.Infow("User requesting user",
		"requester UID", authUser.UID,
		"requested UID", userUID,
	)

	user, err := h.userService.GetByUID(c, userUID)
	if errors.Is(err, store.ErrNotFound) {
		render.Render(w, r, server.ErrNotFound)
		return
	} else if err != nil {
		render.Render(w, r, server.ErrInvalidRequest(err))
		return
	}

	jwtKey := []byte(config.GetString("server.jwt.key"))
	isSecure := config.GetBool("server.tls")
	domain := net.JoinHostPort(config.GetString("server.host"), config.GetString("server.port"))
	tokenAge := config.GetInt("server.jwt.token_age")
	{
		expirationTime := time.Now().Add(time.Duration(tokenAge) * time.Minute)
		token, err := jwt.CreateToken(jwtKey, user.UID, "", expirationTime)
		if err != nil {
			render.Render(w, r, server.ErrInvalidRequest(err))
			return
		}
		http.SetCookie(w, &http.Cookie{Name: "token", Value: token, Expires: expirationTime, Domain: domain, Secure: isSecure})
	}

	render.JSON(w, r, api.UserResponse{
		UID:          user.UID,
		Name:         user.Name,
		EmailAddress: user.EmailAddress,
	})
}
