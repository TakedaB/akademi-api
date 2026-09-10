package main

import (
	"log"
	"net/http"

	"github.com/TakedaB/akademi-api/internal/handler"
	"github.com/TakedaB/akademi-api/internal/middleware"
	"github.com/TakedaB/akademi-api/internal/repository"
	"github.com/TakedaB/akademi-api/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("aviso: .env não encontrado,usando variáveis de ambiente do sistema")
	}

	db := repository.NewPostgresRepository()
	defer db.Close()

	studentRepo := repository.NewStudentRepository(db)
	userRepo := repository.NewUserRepository(db)
	studentService := service.NewStudentService(studentRepo, userRepo)
	studentHandler := handler.NewStudentHandler(studentService)

	teacherRepo := repository.NewTeacherRepository(db)
	teacherService := service.NewTeacherService(teacherRepo, userRepo)
	teacherHandler := handler.NewTeacherHandler(teacherService)

	financeRepo := repository.NewFinanceRepository(db)
	financeService := service.NewFinanceService(financeRepo)
	financeHandler := handler.NewFinanceHandler(financeService)

	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthCheckHandler)
	mux.HandleFunc("POST /login", authHandler.Login)
	mux.HandleFunc("GET /me", middleware.RequireAuth(authHandler.Me))

	staffOnly := middleware.RequireRole("diretoria", "financeiro")
	allRoles := middleware.RequireRole("diretoria", "financeiro", "professor", "aluno")
	directoriaOnly := middleware.RequireRole("diretoria")

	mux.HandleFunc("POST /students", middleware.RequireAuth(staffOnly(studentHandler.Create)))
	mux.HandleFunc("GET /students", middleware.RequireAuth(middleware.RequireRole("diretoria", "financeiro", "professor")(studentHandler.FindAll)))
	mux.HandleFunc("GET /students/{id}", middleware.RequireAuth(middleware.RequireRole("diretoria", "financeiro", "professor")(studentHandler.FindByID)))
	mux.HandleFunc("PUT /students/{id}", middleware.RequireAuth(staffOnly(studentHandler.Update)))
	mux.HandleFunc("DELETE /students/{id}", middleware.RequireAuth(staffOnly(studentHandler.Delete)))

	mux.HandleFunc("POST /teachers", middleware.RequireAuth(directoriaOnly(teacherHandler.Create)))
	mux.HandleFunc("GET /teachers", middleware.RequireAuth(allRoles(teacherHandler.FindAll)))
	mux.HandleFunc("GET /teachers/{id}", middleware.RequireAuth(allRoles(teacherHandler.FindByID)))
	mux.HandleFunc("DELETE /teachers/{id}", middleware.RequireAuth(directoriaOnly(teacherHandler.Delete)))

	mux.HandleFunc("POST /finance", middleware.RequireAuth(staffOnly(financeHandler.Create)))
	mux.HandleFunc("GET /finance", middleware.RequireAuth(staffOnly(financeHandler.FindAll)))
	mux.HandleFunc("GET /students/{studentId}/finance", middleware.RequireAuth(staffOnly(financeHandler.FindByStudentID)))
	mux.HandleFunc("PATCH /finance/{id}/status", middleware.RequireAuth(staffOnly(financeHandler.UpdateStatus)))
	mux.HandleFunc("DELETE /finance/{id}", middleware.RequireAuth(staffOnly(financeHandler.Delete)))

	log.Println("servidor rodando na porta 8080")
	if err := http.ListenAndServe(":8080", middleware.CORS(mux)); err != nil {
		log.Fatal(err)
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}
