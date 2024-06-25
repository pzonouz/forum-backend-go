package middlewares

import (
	"log"
	"net/http"
)

func RoleGuard(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("access")
			if err != nil {
				log.Print(err.Error())
			}

			log.Print(cookie)

			next.ServeHTTP(w, r)
		})
}
