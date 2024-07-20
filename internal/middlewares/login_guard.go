package middlewares

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"forum-backend-go/internal/utils"
)

func LoginGuard(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		access, err := r.Cookie("forum_access")

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)

			return
		}

		token, err := jwt.ParseWithClaims(access.Value, &utils.MyClaims{}, func(_ *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)

			return
		}

		claims := token.Claims.(*utils.MyClaims)

		if claims.Expired < time.Now().Unix() {
			w.WriteHeader(http.StatusUnauthorized)

			return
		}

		f(w, r)
	}
}
