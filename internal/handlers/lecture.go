package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"
   // "time"
)

type LectureHandler struct {
    DB *sql.DB
}

func (h *LectureHandler) GetLectures(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode([]string{"Lecture 1", "Lecture 2"})
}

func (h *LectureHandler) CreateLecturePage(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")
    w.Write([]byte("<h1>Create Lecture - coming soon</h1>"))
}

func (h *LectureHandler) CreateLecture(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Lecture created"})
}