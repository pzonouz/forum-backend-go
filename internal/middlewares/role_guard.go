package middlewares

import (
	"net/http"
	"strings"

	"forum-backend-go/internal/utils"
)

func AdminRoleGuard(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role := utils.GetUserRoleFromRequest(r, w)

		if strings.Compare(role, "admin") != 0 {
			http.Error(w, "not admin user", http.StatusUnauthorized)

			return
		}

		f(w, r)
	}
}
