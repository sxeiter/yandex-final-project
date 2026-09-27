package api

import "net/http"

func Init() {
	http.HandleFunc("/api/signin", SignInHandler)
	http.HandleFunc("/api/nextdate", NextDateHandler)

	http.HandleFunc("/api/task", AuthMiddleware(taskHandler))
	http.HandleFunc("/api/tasks", AuthMiddleware(tasksHandler))
	http.HandleFunc("/api/task/done", AuthMiddleware(doneTaskHandler))
}
