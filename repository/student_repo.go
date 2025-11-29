package repository

import (
	"app2_http_api_database/config"
	"app2_http_api_database/model"
)

func GetAllStudents() ([]model.Student, error) {
	rows, err := config.DB.Query("SELECT id, name, age FROM students")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []model.Student
	for rows.Next() {
		var s model.Student
		if err := rows.Scan(&s.ID, &s.Name, &s.Age); err != nil {
			return nil, err
		}
		students = append(students, s)
	}
	return students, nil
}

func CreateStudent(student model.Student) (int64, error) {
	result, err := config.DB.Exec("INSERT INTO students (name, age) VALUES (?, ?)", student.Name, student.Age)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// جلب طالب واحد حسب ID
func GetStudentByID(id int) (*model.Student, error) {
	row := config.DB.QueryRow("SELECT id, name, age FROM students WHERE id = ?", id)
	var s model.Student
	err := row.Scan(&s.ID, &s.Name, &s.Age)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// تحديث طالب حسب ID
func UpdateStudent(id int, student model.Student) error {
	_, err := config.DB.Exec("UPDATE students SET name=?, age=? WHERE id=?", student.Name, student.Age, id)
	return err
}

// حذف طالب حسب ID
func DeleteStudent(id int) error {
	_, err := config.DB.Exec("DELETE FROM students WHERE id=?", id)
	return err
}

func GetStudentsWithFilter(offset, limit int, minAge int) ([]model.Student, error) {
	query := "SELECT id, name, age FROM students WHERE age >= ? LIMIT ? OFFSET ?"
	rows, err := config.DB.Query(query, minAge, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []model.Student
	for rows.Next() {
		var s model.Student
		if err := rows.Scan(&s.ID, &s.Name, &s.Age); err != nil {
			return nil, err
		}
		students = append(students, s)
	}
	return students, nil
}
