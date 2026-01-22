package router

import (
	"net/http"
	"schoolREST/internal/api/handlers"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.RootHandler)

	mux.HandleFunc("GET /teachers/", handlers.GetTeachersHandler)
	mux.HandleFunc("POST /teachers/", handlers.AddTeachersHandler)
	mux.HandleFunc("PATCH /teachers/", handlers.PatchTeachersHandler)
	mux.HandleFunc("DELETE /teachers/", handlers.DeleteMultipleTeachersHandler)

	mux.HandleFunc("GET /teachers/{id}", handlers.GetOneTeacherHandler)
	mux.HandleFunc("PUT /teachers/{id}", handlers.UpdateTeachersHandler)
	mux.HandleFunc("PATCH /teachers/{id}", handlers.PatchOneTeacherHandler)
	mux.HandleFunc("DELETE /teachers/{id}", handlers.DeleteTeacherHandler)

	mux.HandleFunc("/students/", handlers.StudentsHandler)

	mux.HandleFunc("/execs/", handlers.ExecsHandler)

	return mux
}
