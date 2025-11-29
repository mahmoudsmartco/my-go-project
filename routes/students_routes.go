package routes

import (
	"app2_http_api_database/auth"
	v1 "app2_http_api_database/handler/v1"
	"net/http"
)

// RegisterStudentsRoutes يربط كل endpoint متعلق بالطلاب
func RegisterStudentsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			v1.GetStudents(w, r)
		case http.MethodPost:
			auth.JWTMiddleware(http.HandlerFunc(v1.CreateStudent)).ServeHTTP(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/students/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			v1.GetStudent(w, r)
		case http.MethodPut:
			auth.JWTMiddleware(http.HandlerFunc(v1.UpdateStudentHandler)).ServeHTTP(w, r)
		case http.MethodDelete:
			auth.JWTMiddleware(http.HandlerFunc(v1.DeleteStudentHandler)).ServeHTTP(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
