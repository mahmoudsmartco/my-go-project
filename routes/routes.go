package routes

import (
	"app2_http_api_database/middleware"
	"net/http"
	"time"
)

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// -----------------------------
	// Students Routes + RateLimit
	// -----------------------------
	// هنفترض إن RegisterStudentsRoutes بتسجل route اسمه "/students"
	// لذلك هنلفه داخل RateLimit

	studentMux := http.NewServeMux()
	RegisterStudentsRoutes(studentMux) // routes الأصلية للطلاب

	// نحدد limit → طلب كل 3 ثواني
	rateLimitedStudents := middleware.RateLimit(studentMux, 3*time.Second)

	// نربطه على نفس المسار
	mux.Handle("/students", rateLimitedStudents)
	mux.Handle("/students/", rateLimitedStudents)

	// -----------------------------
	// Protected Routes (بدون RateLimit)
	// -----------------------------
	RegisterProtectedRoutes(mux)

	// -----------------------------
	// Swagger Routes (بدون RateLimit)
	// -----------------------------
	RegisterSwaggerRoutes(mux)

	return mux
}
