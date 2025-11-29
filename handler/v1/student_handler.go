package v1

import (
	"app2_http_api_database/cache"
	"app2_http_api_database/model"
	"app2_http_api_database/repository"
	"app2_http_api_database/service/rabbitmq"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// مدة التخزين في الكاش (مثلاً 30 ثانية)
const cacheTTL = 30 * time.Second

// @Summary Get all students
// @Tags students
// @Produce json
// @Success 200 {array} model.Student
// @Router / [get]
func GetStudents(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	cacheKey := "students:all"

	// قراءة Query Params
	offsetStr := r.URL.Query().Get("offset")
	limitStr := r.URL.Query().Get("limit")
	minAgeStr := r.URL.Query().Get("minAge")

	offset, _ := strconv.Atoi(offsetStr)
	limit, _ := strconv.Atoi(limitStr)
	minAge, _ := strconv.Atoi(minAgeStr)

	if limit == 0 {
		limit = 100 // Default limit
	}

	cached, err := cache.Rdb.Get(ctx, cacheKey).Result()
	if err == nil && cached != "" {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cached))
		return
	}

	students, err := repository.GetStudentsWithFilter(offset, limit, minAge)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(students)
	cache.Rdb.Set(ctx, cacheKey, data, cacheTTL)

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// @Summary Add a new student
// @Tags students
// @Accept json
// @Produce json
// @Param student body model.Student true "Student Data"
// @Success 200 {object} model.Student
// @Router / [post]
func CreateStudent(w http.ResponseWriter, r *http.Request) {
	var s model.Student
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := repository.CreateStudent(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.ID = int(id)

	// ✅ بعد إضافة طالب جديد، نمسح الكاش القديم
	cache.Rdb.Del(context.Background(), "students:all")

	// نشر حدث إلى RabbitMQ (لو DefaultPublisher مُهيأ)
	if rabbitmq.DefaultPublisher != nil {
		evt := rabbitmq.StudentCreatedEvent{
			ID:    s.ID,
			Name:  s.Name,
			Email: s.Age,
			When:  time.Now().Unix(),
		}
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		err := rabbitmq.DefaultPublisher.PublishStudentCreated(ctx, evt, "students.created")
		if err != nil {
			// لا نفرّط في فشل النشر، فقط نسجل اللوق
			log.Printf("warning: publish student.created failed: %v", err)
		}
	}

	json.NewEncoder(w).Encode(s)
}

// GET /students/{id}
// @Summary Get student by ID
// @Tags students
// @Produce json
// @Param id path int true "Student ID"
// @Success 200 {object} model.Student
// @Failure 404 {string} string "Student not found"
// @Router /{id} [get]
func GetStudent(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/students/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	cacheKey := fmt.Sprintf("student:%d", id)
	ctx := context.Background()

	// حاول نجيب من Redis
	cached, err := cache.Rdb.Get(ctx, cacheKey).Result()
	if err == nil && cached != "" {
		fmt.Println("📦 Fetched from Redis cache")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cached))
		return
	}

	// جلب من MySQL
	student, err := repository.GetStudentByID(id)
	if err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	data, _ := json.Marshal(student)
	cache.Rdb.Set(ctx, cacheKey, data, cacheTTL)

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// PUT /students/{id}
// @Summary Update student by ID
// @Tags students
// @Accept json
// @Produce json
// @Param id path int true "Student ID"
// @Param student body model.Student true "Student Data"
// @Success 200 {object} model.Student
// @Failure 404 {string} string "Student not found"
// @Router /{id} [put]
func UpdateStudentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/students/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var s model.Student
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = repository.UpdateStudent(id, s)
	if err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	// مسح الكاش القديم
	cache.Rdb.Del(context.Background(), fmt.Sprintf("student:%d", id))
	cache.Rdb.Del(context.Background(), "students:all")

	s.ID = id
	json.NewEncoder(w).Encode(s)
}

// DELETE /students/{id}
// @Summary Delete student by ID
// @Tags students
// @Produce json
// @Param id path int true "Student ID"
// @Success 200 {string} string "Student deleted successfully"
// @Failure 404 {string} string "Student not found"
// @Router /{id} [delete]
func DeleteStudentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/students/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = repository.DeleteStudent(id)
	if err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	cache.Rdb.Del(context.Background(), fmt.Sprintf("student:%d", id))
	cache.Rdb.Del(context.Background(), "students:all")

	w.Write([]byte("Student deleted successfully"))
}
