package main

import (
	"log"
	"net/http"

	"github.com/TakedaB/akademi-api/internal/handler"
	"github.com/TakedaB/akademi-api/internal/middleware"
	"github.com/TakedaB/akademi-api/internal/repository"
	"github.com/TakedaB/akademi-api/internal/service"
)

func main() {
	db := repository.NewPostgresRepository()
	defer db.Close()

	studentRepo := repository.NewStudentRepository(db)
	studentService := service.NewStudentService(studentRepo)
	studentHandler := handler.NewStudentHandler(studentService)

	userRepo := repository.NewUserRepository(db)
	teacherRepo := repository.NewTeacherRepository(db)
	teacherService := service.NewTeacherService(teacherRepo, userRepo)
	teacherHandler := handler.NewTeacherHandler(teacherService)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthCheckHandler)

	mux.HandleFunc("POST /students", studentHandler.Create)
	mux.HandleFunc("GET /students", studentHandler.FindAll)
	mux.HandleFunc("GET /students/{id}", studentHandler.FindByID)
	mux.HandleFunc("PUT /students/{id}", studentHandler.Update)
	mux.HandleFunc("DELETE /students/{id}", studentHandler.Delete)

	mux.HandleFunc("POST /teachers", teacherHandler.Create)
	mux.HandleFunc("GET /teachers", teacherHandler.FindAll)
	mux.HandleFunc("GET /teachers/{id}", teacherHandler.FindByID)
	mux.HandleFunc("DELETE /teachers/{id}", teacherHandler.Delete)

	log.Println("servidor rodando na porta 8080")
	if err := http.ListenAndServe(":8080", middleware.CORS(mux)); err != nil {
		log.Fatal(err)
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}
