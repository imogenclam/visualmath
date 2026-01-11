package main

import (
	"fmt"
	"net/http"
)

func main() {
    fmt.Println("✅ Сервер запускается...")
    
    // Самый простой сервер
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, `
        <html>
        <head><title>VisualMath</title></head>
        <body>
            <h1>🎉 Ура! Сервер работает!</h1>
            <p>Если вы это видите, значит все настроено правильно!</p>
            <a href="/test">Тестовая страница</a>
        </body>
        </html>
        `)
    })
    
    http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Тестовая страница работает!")
    })
    
    fmt.Println("🌐 Откройте в браузере: http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}