package middleware

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	config "github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/ldej/go-rest-example/internal/api"
	"github.com/ldej/go-rest-example/pkg/jwt"
	"github.com/ldej/go-rest-example/pkg/server"
)

func RequireUser(next func(http.ResponseWriter, *http.Request, *api.AuthUser)) http.HandlerFunc {
	logger := zap.S().With("package", "middleware.auth")

	jwtKey := []byte(config.GetString("server.jwt_key"))
	return func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie("token")
		if errors.Is(err, http.ErrNoCookie) {
			logger.Infow("cookie 'token' not present")
			render.Render(w, r, server.ErrUnauthorized(err))
			return
		} else if err != nil {
			logger.Warnw("something went wrong", "error", err.Error())
			render.Render(w, r, server.ErrInvalidRequest(err))
			return
		}

		claims, err := jwt.ValidateToken(jwtKey, tokenCookie.Value)
		if err != nil {
			logger.Infow("unauthorized", "error", err.Error())
			render.Render(w, r, server.ErrUnauthorized(err))
			return
		}

		if claims.UserUID == "" {
			// No user so redirect to login
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		user := &api.AuthUser{
			UID:  claims.UserUID,
			Role: claims.Role,
		}
		next(w, r, user)
	}
}
