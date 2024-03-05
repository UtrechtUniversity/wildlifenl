package wildlifenl

import (
	"net/http"
	"strings"
)

func getBearerToken(r *http.Request) string {
	if bearer, ok := r.Header["Authorization"]; ok {
		if len(bearer) != 1 {
			return ""
		}
		parts := strings.Split(bearer[0], " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return ""
		}
		return parts[1]
	}
	return ""
}

func writeResponseJSON(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
