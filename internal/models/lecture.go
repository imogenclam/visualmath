package models

import (
	"time"
	"encoding/json"
)

type Lecture struct {
	ID          int             `json:"id"`
	Title       string          `json:"title"`
	CourseID    int             `json:"course_id"`
	CourseName  string          `json:"course_name"`
	AuthorID    int             `json:"author_id"`
	AuthorName  string          `json:"author_name"`
	Description string          `json:"description"`
	Modules     []LectureModule `json:"modules"`
	CreatedAt   time.Time       `json:"created_at"`
	Published   bool            `json:"published"`
	AllowBack   bool            `json:"allow_back"` // можно ли возвращаться к пройденным
}

type LectureModule struct {
	ID        int             `json:"id"`
	LectureID int             `json:"lecture_id"`
	ModuleID  int             `json:"module_id"`
	Order     int             `json:"order"`
	Title     string          `json:"title"`
	Type      string          `json:"type"`
	Module    json.RawMessage `json:"module"`
}

// LectureRequest для создания/обновления лекции
type LectureRequest struct {
	Title       string   `json:"title"`
	CourseID    int      `json:"course_id"`
	CourseName  string   `json:"course_name"`
	Description string   `json:"description"`
	ModuleIDs   []int    `json:"module_ids"` // ID модулей в порядке
	Published   bool     `json:"published"`
	AllowBack   bool     `json:"allow_back"`
}

// StudentProgress для отслеживания прогресса
type StudentProgress struct {
	ID         int       `json:"id"`
	StudentID  int       `json:"student_id"`
	LectureID  int       `json:"lecture_id"`
	ModuleID   int       `json:"module_id"`
	Completed  bool      `json:"completed"`
	Score      float64   `json:"score"`
	StartedAt  time.Time `json:"started_at"`
	CompletedAt time.Time `json:"completed_at"`
}
