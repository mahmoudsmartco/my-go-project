package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// InitDB يقوم بإنشاء الاتصال بقاعدة البيانات
func InitDB() {
	// بيانات الاتصال (بدلها بالقيم المناسبة عندك)
	dbUser := "root"
	dbPass := "root"
	dbHost := "127.0.0.1"
	dbPort := "3306"
	dbName := "school"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPass, dbHost, dbPort, dbName)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("❌ فشل الاتصال بقاعدة البيانات: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("❌ لا يمكن الاتصال بـ MySQL: %v", err)
	}

	fmt.Println("✅ تم الاتصال بقاعدة بيانات MySQL بنجاح")
}

// CloseDB لإغلاق الاتصال عند إيقاف السيرفر
func CloseDB() {
	if DB != nil {
		DB.Close()
		fmt.Println("🟡 تم إغلاق اتصال MySQL")
	}
}
