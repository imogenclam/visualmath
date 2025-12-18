package storage

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "log"
)

func InitSQLite() *sql.DB {
    db, err := sql.Open("sqlite3", "./visualmath.db")
    if err != nil {
        log.Fatal(err)
    }
    
    // Создаем таблицы
    queries := []string{
        `CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            login TEXT UNIQUE NOT NULL,
            password_hash TEXT NOT NULL,
            full_name TEXT NOT NULL,
            user_type TEXT NOT NULL CHECK (user_type IN ('student', 'teacher', 'admin')),
            group_number TEXT,
            email TEXT UNIQUE NOT NULL,
            email_verified BOOLEAN DEFAULT FALSE,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )`,
        
        `CREATE TABLE IF NOT EXISTS courses (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT UNIQUE NOT NULL
        )`,
        
        `INSERT OR IGNORE INTO courses (name) VALUES 
            ('Математический анализ'),
            ('Линейная алгебра и аналитическая геометрия'),
            ('Дискретная математика'),
            ('Экономика')`,
    }
    
    for _, query := range queries {
        _, err := db.Exec(query)
        if err != nil {
            log.Printf("Warning: %v", err)
        }
    }
    
    return db
}