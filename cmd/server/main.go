package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"visualmath/internal/handlers"
)

func main() {
	fmt.Println("🚀 Запускаем VisualMath сервер...")

	godotenv.Load()

	// Создаем обработчик модулей
	moduleHandler := &handlers.ModuleHandler{}
	lectureHandler := &handlers.LectureHandler{}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Handle("/static/*", http.StripPrefix("/static/",
		http.FileServer(http.Dir("web/static"))))

	// Публичные маршруты
	r.Get("/", homeHandler)
	r.Get("/login", loginPageHandler)
	r.Get("/register", registerPageHandler)
	r.Get("/dashboard", dashboardHandler)
	r.Get("/test", testHandler)

	// Маршруты модулей
	r.Get("/modules", moduleHandler.ListModules)              // Список модулей
	r.Get("/modules/create", moduleHandler.CreateModulePage)  // Страница создания
	r.Get("/modules/view/{id}", moduleHandler.ViewModulePage) // Просмотр модуля
	r.Get("/modules/edit/{id}", moduleHandler.EditModulePage) // Редактирование модуля

	// API endpoints для модулей
	r.Get("/api/modules/list", moduleHandler.ListModulesAPI)  // API: список модулей
	r.Post("/api/modules", moduleHandler.CreateModule)        // API: создание модуля
	r.Get("/api/modules/{id}", moduleHandler.GetModule)       // API: получить модуль
	r.Put("/api/modules/{id}", moduleHandler.UpdateModule)    // API: обновить модуль
	r.Delete("/api/modules/{id}", moduleHandler.DeleteModule) // API: удалить модуль

	// Маршруты лекций
	r.Get("/lectures", lecturesPageHandler)              // страница списка лекций
	r.Get("/lectures/create", createLecturePageHandler)  // страница создания лекции
	r.Get("/lectures/edit/{id}", editLecturePageHandler) // страница редактирования лекции
	r.Get("/lectures/view/{id}", viewLecturePageHandler) // страница просмотра лекции

	// API endpoints для лекций
	r.Get("/api/lectures", lectureHandler.ListLectures)
	r.Post("/api/lectures", lectureHandler.CreateLecture)
	r.Get("/api/lectures/{id}", lectureHandler.GetLecture)
	r.Put("/api/lectures/{id}", lectureHandler.UpdateLecture)
	r.Delete("/api/lectures/{id}", lectureHandler.DeleteLecture)
	r.Get("/api/modules/available", lectureHandler.GetAvailableModules)
	r.Post("/api/lectures/start", lectureHandler.StartLecture)
	r.Post("/api/lectures/complete", lectureHandler.CompleteModule)
	r.Get("/api/lectures/progress", lectureHandler.GetStudentProgress)

	// API заглушки
	r.Post("/api/register", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Register endpoint"}`)
	})

	r.Post("/api/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Login endpoint"}`)
	})

	port := "8080"
	fmt.Printf("✅ Сервер запущен на http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

// homeHandler обрабатывает главную страницу
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>VisualMath - Главная</title>
    <link rel="stylesheet" href="/static/css/style.css">
    <style>
        .oauth-buttons {
            display: flex;
            flex-direction: column;
            align-items: center;
            gap: 12px;
            margin: 25px 0;
        }
        .oauth-btn {
            display: flex;
            align-items: center;
            justify-content: center;
            width: 280px;
            padding: 14px 20px;
            border-radius: 8px;
            text-decoration: none;
            font-weight: 500;
            font-size: 16px;
            transition: all 0.3s ease;
            border: 1px solid #ddd;
            color: white;
        }
        .oauth-btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 12px rgba(0,0,0,0.15);
        }
        .oauth-btn.vk {
            background: #4a76a8;
        }
        .oauth-btn.vk:hover {
            background: #3a6398;
        }
        .oauth-btn.google {
            background: #db4437;
        }
        .oauth-btn.google:hover {
            background: #c23327;
        }
        .oauth-btn .icon {
            margin-right: 12px;
            font-size: 20px;
        }
        .divider {
            display: flex;
            align-items: center;
            margin: 20px 0;
            width: 100%;
            max-width: 280px;
        }
        .divider::before,
        .divider::after {
            content: '';
            flex: 1;
            height: 1px;
            background: #ddd;
        }
        .divider span {
            padding: 0 15px;
            color: #777;
            font-size: 14px;
        }
    </style>
</head>
<body>
    <div style="text-align: center; padding: 50px;">
        <h1 style="color: #2c3e50; margin-bottom: 10px;">🎓 VisualMath Platform</h1>
        <p style="color: #7f8c8d; margin-bottom: 40px;">Платформа для интерактивного изучения математики</p>
        
        <!-- Кнопки быстрого входа через OAuth -->
        <div class="oauth-buttons">
            <a href="/auth/vk" class="oauth-btn vk">
                <span class="icon">VK</span>
                Войти через ВКонтакте
            </a>
            <a href="/auth/google" class="oauth-btn google">
                <span class="icon">G</span>
                Войти через Google
            </a>
        </div>
        
        <div class="divider">
            <span>или используйте почту</span>
        </div>
        
        <!-- Стандартные кнопки -->
        <div style="margin: 30px 0;">
            <a href="/login" style="display: inline-block; padding: 14px 30px; background: #3498db; 
               color: white; text-decoration: none; border-radius: 8px; margin: 10px; font-weight: 500;">
               🔑 Войти в аккаунт
            </a>
            <a href="/register" style="display: inline-block; padding: 14px 30px; background: #2ecc71; 
               color: white; text-decoration: none; border-radius: 8px; margin: 10px; font-weight: 500;">
               📝 Создать аккаунт
            </a>
        </div>
        
        <!-- Простой информационный блок -->
        <div style="margin-top: 40px; padding-top: 20px; border-top: 1px solid #e0e0e0;">
            <p style="margin: 5px 0;">
                <a href="/test">Тестовая страница</a> | 
                <a href="/static/css/style.css">CSS файл</a> |
                <a href="/dashboard">Личный кабинет</a>
            </p>
        </div>
    </div>
    
    <script>
        // Автоматическое определение OAuth провайдеров
        window.addEventListener('DOMContentLoaded', function() {
            // Проверяем параметры URL
            const urlParams = new URLSearchParams(window.location.search);
            const oauthSuccess = urlParams.get('oauth_success');
            
            if (oauthSuccess === 'true') {
                alert('✅ OAuth авторизация успешна! Пожалуйста, войдите снова для подтверждения.');
            }
            
            // Проверяем, авторизован ли пользователь
            const token = localStorage.getItem('token');
            if (token) {
                setTimeout(() => {
                    window.location.href = '/dashboard';
                }, 1000);
            }
        });
    </script>
</body>
</html>`

	fmt.Fprintf(w, html)
}

// loginPageHandler обрабатывает страницу входа
func loginPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := `<!DOCTYPE html>
<html>
<head>
    <title>Вход - VisualMath</title>
    <link rel="stylesheet" href="/static/css/style.css">
    <style>
        .login-container {
            max-width: 420px;
            margin: 60px auto;
            padding: 40px;
            background: white;
            border-radius: 12px;
            box-shadow: 0 8px 25px rgba(0,0,0,0.1);
        }
        .login-header {
            text-align: center;
            margin-bottom: 30px;
        }
        .login-header h1 {
            color: #2c3e50;
            margin-bottom: 8px;
            font-size: 28px;
        }
        .login-header p {
            color: #7f8c8d;
            font-size: 16px;
        }
        .oauth-buttons {
            display: flex;
            flex-direction: column;
            gap: 12px;
            margin-bottom: 25px;
        }
        .oauth-btn {
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 14px;
            border-radius: 8px;
            text-decoration: none;
            font-weight: 500;
            font-size: 15px;
            transition: all 0.3s ease;
            border: 1px solid #ddd;
        }
        .oauth-btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 12px rgba(0,0,0,0.15);
        }
        .oauth-btn.vk {
            background: #4a76a8;
            color: white;
        }
        .oauth-btn.vk:hover {
            background: #3a6398;
        }
        .oauth-btn.google {
            background: #fff;
            color: #444;
            border: 1px solid #ddd;
        }
        .oauth-btn.google:hover {
            background: #f8f9fa;
            border-color: #ccc;
        }
        .oauth-btn .icon {
            margin-right: 12px;
            font-weight: bold;
            font-size: 16px;
        }
        .divider {
            display: flex;
            align-items: center;
            margin: 25px 0;
        }
        .divider::before,
        .divider::after {
            content: '';
            flex: 1;
            height: 1px;
            background: #eee;
        }
        .divider span {
            padding: 0 15px;
            color: #95a5a6;
            font-size: 14px;
            text-transform: uppercase;
        }
        .form-group {
            margin-bottom: 20px;
        }
        .form-group label {
            display: block;
            margin-bottom: 8px;
            color: #2c3e50;
            font-weight: 500;
            font-size: 14px;
        }
        .form-group input {
            width: 100%;
            padding: 14px 16px;
            border: 1px solid #ddd;
            border-radius: 8px;
            font-size: 16px;
            transition: border-color 0.3s;
        }
        .form-group input:focus {
            outline: none;
            border-color: #3498db;
            box-shadow: 0 0 0 3px rgba(52, 152, 219, 0.1);
        }
        .submit-btn {
            width: 100%;
            padding: 15px;
            background: #3498db;
            color: white;
            border: none;
            border-radius: 8px;
            font-size: 16px;
            font-weight: 500;
            cursor: pointer;
            transition: background 0.3s;
        }
        .submit-btn:hover {
            background: #2980b9;
        }
        .form-links {
            text-align: center;
            margin-top: 25px;
            padding-top: 20px;
            border-top: 1px solid #eee;
        }
        .form-links a {
            color: #3498db;
            text-decoration: none;
            margin: 0 10px;
        }
        .form-links a:hover {
            text-decoration: underline;
        }
        .message {
            padding: 12px 16px;
            border-radius: 8px;
            margin-bottom: 20px;
            display: none;
        }
        .message.success {
            background: #d4edda;
            color: #155724;
            border: 1px solid #c3e6cb;
        }
        .message.error {
            background: #f8d7da;
            color: #721c24;
            border: 1px solid #f5c6cb;
        }
    </style>
</head>
<body style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); min-height: 100vh; padding: 20px;">
    <div class="login-container">
        <div class="login-header">
            <h1>🔐 Вход в систему</h1>
            <p>Войдите в свой аккаунт VisualMath</p>
        </div>
        
        <!-- Сообщения -->
        <div id="message" class="message"></div>
        
        <!-- OAuth кнопки -->
        <div class="oauth-buttons">
            <a href="/auth/vk" class="oauth-btn vk">
                <span class="icon">VK</span>
                Войти через ВКонтакте
            </a>
            <a href="/auth/google" class="oauth-btn google">
                <span class="icon">G</span>
                Войти через Google
            </a>
        </div>
        
        <div class="divider">
            <span>или</span>
        </div>
        
        <!-- Форма стандартного входа -->
        <form id="loginForm">
            <div class="form-group">
                <label for="login">Логин или Email:</label>
                <input type="text" id="login" name="login" placeholder="Введите логин или email" required>
            </div>
            
            <div class="form-group">
                <label for="password">Пароль:</label>
                <input type="password" id="password" name="password" placeholder="Введите пароль" required>
            </div>
            
            <button type="submit" class="submit-btn">Войти в аккаунт</button>
        </form>
        
        <div class="form-links">
            <a href="/register">Нет аккаунта? Зарегистрироваться</a>
            <a href="/">На главную страницу</a>
        </div>
    </div>

    <script>
        // Обработка формы входа
        document.getElementById('loginForm').addEventListener('submit', async function(e) {
            e.preventDefault();
            
            const messageDiv = document.getElementById('message');
            messageDiv.style.display = 'none';
            
            const formData = {
                login: document.getElementById('login').value,
                password: document.getElementById('password').value
            };
            
            try {
                const response = await fetch('/api/login', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(formData)
                });
                
                const result = await response.json();
                
                if (response.ok) {
                    messageDiv.className = 'message success';
                    messageDiv.textContent = '✅ Вход выполнен успешно! Перенаправление...';
                    messageDiv.style.display = 'block';
                    
                    // Сохраняем токен
                    localStorage.setItem('token', result.token);
                    localStorage.setItem('user', JSON.stringify(result.user));
                    
                    // Перенаправляем на dashboard
                    setTimeout(() => {
                        window.location.href = '/dashboard';
                    }, 1500);
                    
                } else {
                    messageDiv.className = 'message error';
                    messageDiv.textContent = '❌ ' + (result.message || 'Неверный логин или пароль');
                    messageDiv.style.display = 'block';
                }
            } catch (error) {
                messageDiv.className = 'message error';
                messageDiv.textContent = '❌ Ошибка сети: ' + error.message;
                messageDiv.style.display = 'block';
            }
        });
        
        // Проверяем, авторизован ли пользователь
        window.addEventListener('DOMContentLoaded', function() {
            const token = localStorage.getItem('token');
            if (token) {
                const messageDiv = document.getElementById('message');
                messageDiv.className = 'message success';
                messageDiv.textContent = '✅ Вы уже авторизованы. Перенаправление...';
                messageDiv.style.display = 'block';
                
                setTimeout(() => {
                    window.location.href = '/dashboard';
                }, 1000);
            }
        });
    </script>
</body>
</html>`

	fmt.Fprintf(w, html)
}

// registerPageHandler обрабатывает страницу регистрации
func registerPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := `<!DOCTYPE html>
<html>
<head>
    <title>Регистрация - VisualMath</title>
    <link rel="stylesheet" href="/static/css/style.css">
    <style>
        .register-container {
            max-width: 450px;
            margin: 40px auto;
            padding: 40px;
            background: white;
            border-radius: 12px;
            box-shadow: 0 8px 25px rgba(0,0,0,0.1);
        }
        .register-header {
            text-align: center;
            margin-bottom: 30px;
        }
        .register-header h1 {
            color: #2c3e50;
            margin-bottom: 8px;
            font-size: 28px;
        }
        .register-header p {
            color: #7f8c8d;
            font-size: 16px;
        }
        .oauth-section {
            margin-bottom: 25px;
            text-align: center;
        }
        .oauth-title {
            color: #95a5a6;
            font-size: 14px;
            margin-bottom: 15px;
            text-transform: uppercase;
            letter-spacing: 1px;
        }
        .oauth-buttons {
            display: flex;
            gap: 15px;
            justify-content: center;
        }
        .oauth-btn {
            display: flex;
            align-items: center;
            padding: 12px 20px;
            border-radius: 8px;
            text-decoration: none;
            font-weight: 500;
            font-size: 14px;
            transition: all 0.3s ease;
            border: 1px solid #ddd;
        }
        .oauth-btn.vk {
            background: #4a76a8;
            color: white;
        }
        .oauth-btn.vk:hover {
            background: #3a6398;
            transform: translateY(-2px);
        }
        .oauth-btn.google {
            background: #fff;
            color: #444;
            border: 1px solid #ddd;
        }
        .oauth-btn.google:hover {
            background: #f8f9fa;
            border-color: #ccc;
            transform: translateY(-2px);
        }
        .oauth-btn .icon {
            margin-right: 8px;
            font-weight: bold;
        }
        .divider {
            display: flex;
            align-items: center;
            margin: 25px 0;
        }
        .divider::before,
        .divider::after {
            content: '';
            flex: 1;
            height: 1px;
            background: #eee;
        }
        .divider span {
            padding: 0 15px;
            color: #95a5a6;
            font-size: 14px;
        }
        .form-group {
            margin-bottom: 20px;
        }
        .form-group label {
            display: block;
            margin-bottom: 8px;
            color: #2c3e50;
            font-weight: 500;
            font-size: 14px;
        }
        .form-group input,
        .form-group select {
            width: 100%;
            padding: 14px 16px;
            border: 1px solid #ddd;
            border-radius: 8px;
            font-size: 16px;
            transition: border-color 0.3s;
        }
        .form-group input:focus,
        .form-group select:focus {
            outline: none;
            border-color: #2ecc71;
            box-shadow: 0 0 0 3px rgba(46, 204, 113, 0.1);
        }
        .submit-btn {
            width: 100%;
            padding: 15px;
            background: #2ecc71;
            color: white;
            border: none;
            border-radius: 8px;
            font-size: 16px;
            font-weight: 500;
            cursor: pointer;
            transition: background 0.3s;
            margin-top: 10px;
        }
        .submit-btn:hover {
            background: #27ae60;
        }
        .form-links {
            text-align: center;
            margin-top: 25px;
            padding-top: 20px;
            border-top: 1px solid #eee;
        }
        .form-links a {
            color: #3498db;
            text-decoration: none;
            margin: 0 10px;
        }
        .form-links a:hover {
            text-decoration: underline;
        }
        .message {
            padding: 12px 16px;
            border-radius: 8px;
            margin-bottom: 20px;
            display: none;
        }
        .message.success {
            background: #d4edda;
            color: #155724;
            border: 1px solid #c3e6cb;
        }
        .message.error {
            background: #f8d7da;
            color: #721c24;
            border: 1px solid #f5c6cb;
        }
        .info-text {
            font-size: 13px;
            color: #95a5a6;
            margin-top: 5px;
        }
    </style>
</head>
<body style="background: linear-gradient(135deg, #1a2980 0%, #26d0ce 100%); min-height: 100vh; padding: 20px;">
    <div class="register-container">
        <div class="register-header">
            <h1>📝 Регистрация студента</h1>
            <p>Создайте аккаунт VisualMath для студентов</p>
        </div>
        
        <!-- Сообщения -->
        <div id="message" class="message"></div>
        
        <!-- Быстрая регистрация через OAuth -->
        <div class="oauth-section">
            <div class="oauth-title">Быстрая регистрация</div>
            <div class="oauth-buttons">
                <a href="/auth/vk" class="oauth-btn vk">
                    <span class="icon">VK</span>
                    VK
                </a>
                <a href="/auth/google" class="oauth-btn google">
                    <span class="icon">G</span>
                    Google
                </a>
            </div>
        </div>
        
        <div class="divider">
            <span>или через email</span>
        </div>
        
        <!-- Форма стандартной регистрации -->
        <form id="registerForm">
            <div class="form-group">
                <label for="login">Логин *</label>
                <input type="text" id="login" name="login" placeholder="Придумайте логин" required>
            </div>
            
            <div class="form-group">
                <label for="email">Email *</label>
                <input type="email" id="email" name="email" placeholder="Ваш email" required>
            </div>
            
            <div class="form-group">
                <label for="full_name">ФИО *</label>
                <input type="text" id="full_name" name="full_name" placeholder="Иванов Иван Иванович" required>
            </div>
            
            <!-- Скрытое поле для типа пользователя -->
            <input type="hidden" id="user_type" name="user_type" value="student">
            
            <div class="form-group">
                <label for="password">Пароль *</label>
                <input type="password" id="password" name="password" placeholder="Придумайте пароль" required>
                <div class="info-text">Минимум 8 символов</div>
            </div>
            
            <div class="form-group">
                <label for="group_number">Номер группы *</label>
                <input type="text" id="group_number" name="group_number" placeholder="Например: ИУ6-32Б" required>
                <div class="info-text">Укажите вашу учебную группу</div>
            </div>
            
            <button type="submit" class="submit-btn">Создать аккаунт</button>
        </form>
        
        <div class="form-links">
            <a href="/login">Уже есть аккаунт? Войти</a>
            <a href="/">На главную страницу</a>
        </div>
    </div>

    <script>
        // Обработка формы регистрации
        document.getElementById('registerForm').addEventListener('submit', async function(e) {
            e.preventDefault();
            
            const messageDiv = document.getElementById('message');
            messageDiv.style.display = 'none';
            
            const formData = {
                login: document.getElementById('login').value,
                password: document.getElementById('password').value,
                full_name: document.getElementById('full_name').value,
                email: document.getElementById('email').value,
                user_type: document.getElementById('user_type').value,
                group_number: document.getElementById('group_number').value
            };
            
            // Проверка пароля
            if (formData.password.length < 8) {
                messageDiv.className = 'message error';
                messageDiv.textContent = '❌ Пароль должен содержать минимум 8 символов';
                messageDiv.style.display = 'block';
                return;
            }
            
            try {
                const response = await fetch('/api/register', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(formData)
                });
                
                const result = await response.json();
                
                if (response.ok) {
                    messageDiv.className = 'message success';
                    messageDiv.textContent = '✅ Регистрация успешна! Перенаправление на страницу входа...';
                    messageDiv.style.display = 'block';
                    
                    // Очищаем форму
                    document.getElementById('registerForm').reset();
                    
                    // Перенаправляем на страницу входа
                    setTimeout(() => {
                        window.location.href = '/login';
                    }, 2000);
                    
                } else {
                    messageDiv.className = 'message error';
                    messageDiv.textContent = '❌ ' + (result.message || 'Ошибка регистрации');
                    messageDiv.style.display = 'block';
                }
            } catch (error) {
                messageDiv.className = 'message error';
                messageDiv.textContent = '❌ Ошибка сети: ' + error.message;
                messageDiv.style.display = 'block';
            }
        });
    </script>
</body>
</html>`

	fmt.Fprintf(w, html)
}

// dashboardHandler показывает личный кабинет
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := `<!DOCTYPE html>
<html>
<head>
    <title>Личный кабинет - VisualMath</title>
    <style>
        .dashboard-container {
            display: flex;
            min-height: 100vh;
            background: #f8f9fa;
        }
        .sidebar {
            width: 280px;
            background: #2c3e50;
            color: white;
            padding: 20px;
        }
        .user-info {
            background: #34495e;
            padding: 20px;
            border-radius: 8px;
            margin-bottom: 30px;
        }
        .menu-section {
            margin-bottom: 30px;
        }
        .menu-section h3 {
            color: #95a5a6;
            font-size: 14px;
            text-transform: uppercase;
            margin-bottom: 15px;
            padding-bottom: 5px;
            border-bottom: 1px solid #4a6572;
        }
        .menu-section ul {
            list-style: none;
            padding: 0;
            margin: 0;
        }
        .menu-section li {
            margin-bottom: 8px;
        }
        .menu-section a {
            display: block;
            color: #ecf0f1;
            text-decoration: none;
            padding: 10px 15px;
            border-radius: 6px;
            transition: background 0.3s;
        }
        .menu-section a:hover {
            background: #4a6572;
            text-decoration: none;
        }
        .main-content {
            flex: 1;
            padding: 40px;
        }
        .welcome-card {
            background: white;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
            margin-bottom: 30px;
        }
    </style>
</head>
<body>
    <div class="dashboard-container">
        <aside class="sidebar">
            <div class="user-info">
                <h3>👤 Тестовый пользователь</h3>
                <p>Тип: Преподаватель</p>
                <a href="/" style="color: #e74c3c; margin-top: 15px; display: inline-block;">🚪 Выйти</a>
            </div>
            <div class="menu-section">
                <h3>👨‍🏫 Преподаватель</h3>
                <ul>
                    <li><a href="#">📚 Библиотека модулей</a></li>
                    <li><a href="#">📖 Библиотека лекций</a></li>
                    <li><a href="#">🚀 Начать лекцию</a></li>
                    <li><a href="#">➕ Создать модуль</a></li>
                    <li><a href="#">➕ Создать лекцию</a></li>
                </ul>
            </div>
        </aside>
        
        <main class="main-content">
            <div class="welcome-card">
                <h1>👋 Добро пожаловать!</h1>
                <p>Личный кабинет VisualMath. Выберите действие в меню слева.</p>
            </div>
        </main>
    </div>
</body>
</html>`

	fmt.Fprintf(w, html)
}

// testHandler для проверки работы
func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := `<!DOCTYPE html>
<html>
<head><title>Тест</title></head>
<body style="padding: 50px; text-align: center;">
    <h1 style="color: green;">✅ Тест пройден!</h1>
    <p>Сервер работает правильно</p>
    <p><a href="/">Вернуться на главную</a></p>
</body>
</html>`
	fmt.Fprintf(w, html)
}

// lecturesPageHandler показывает список лекций
func lecturesPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := `<!DOCTYPE html>
<html>
<head>
    <title>Лекции - VisualMath</title>
    <link rel="stylesheet" href="/static/css/style.css">
    <style>
        .lectures-container { max-width: 1200px; margin: 30px auto; padding: 0 20px; }
        .page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 30px; }
        .create-btn { background: #2ecc71; color: white; padding: 12px 24px; border-radius: 6px; text-decoration: none; }
        .lectures-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); gap: 25px; }
        .lecture-card { background: white; border-radius: 10px; padding: 20px; box-shadow: 0 4px 12px rgba(0,0,0,0.1); }
    </style>
</head>
<body>
    <div class="lectures-container">
        <div class="page-header">
            <h1>📖 Лекции</h1>
            <a href="/lectures/create" class="create-btn">➕ Создать лекцию</a>
        </div>
        <div class="lectures-grid">
            <div class="lecture-card">
                <h3>Введение в математический анализ</h3>
                <p>Математический анализ • 5 модулей</p>
                <a href="/lectures/view/1">Открыть →</a>
            </div>
            <div class="lecture-card">
                <h3>Линейная алгебра для начинающих</h3>
                <p>Линейная алгебра • 4 модуля</p>
                <a href="/lectures/view/2">Открыть →</a>
            </div>
        </div>
    </div>
</body>
</html>`

	fmt.Fprintf(w, html)
}

// createLecturePageHandler показывает страницу создания лекции
func createLecturePageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := `<!DOCTYPE html>
<html>
<head>
    <title>Создать лекцию - VisualMath</title>
    <link rel="stylesheet" href="/static/css/style.css">
    <!-- MathJax для предпросмотра -->
    <script src="https://polyfill.io/v3/polyfill.min.js?features=es6"></script>
    <script id="MathJax-script" async src="https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-mml-chtml.js"></script>
    <script>
        MathJax = {
            tex: {
                inlineMath: [['$', '$'], ['\\(', '\\)']],
                displayMath: [['$$', '$$'], ['\\[', '\\]']]
            },
            svg: {
                fontCache: 'global'
            }
        };
    </script>
    <style>
        .create-container {
            max-width: 1200px;
            margin: 30px auto;
            padding: 0 20px;
        }
        .page-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 30px;
            padding-bottom: 20px;
            border-bottom: 2px solid #eee;
        }
        .page-header h1 {
            color: #2c3e50;
            margin: 0;
        }
        .two-column {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 30px;
            margin-top: 20px;
        }
        .form-section {
            background: white;
            padding: 30px;
            border-radius: 12px;
            box-shadow: 0 4px 12px rgba(0,0,0,0.1);
        }
        .modules-section {
            background: white;
            padding: 30px;
            border-radius: 12px;
            box-shadow: 0 4px 12px rgba(0,0,0,0.1);
        }
        .form-group {
            margin-bottom: 25px;
        }
        .form-group label {
            display: block;
            margin-bottom: 8px;
            font-weight: 500;
            color: #2c3e50;
            font-size: 16px;
        }
        .form-group input,
        .form-group select,
        .form-group textarea {
            width: 100%;
            padding: 14px 16px;
            border: 1px solid #ddd;
            border-radius: 8px;
            font-size: 16px;
            transition: border-color 0.3s;
        }
        .form-group input:focus,
        .form-group select:focus,
        .form-group textarea:focus {
            outline: none;
            border-color: #3498db;
            box-shadow: 0 0 0 3px rgba(52, 152, 219, 0.1);
        }
        .form-group textarea {
            min-height: 100px;
            resize: vertical;
        }
        .modules-list {
            border: 2px dashed #e0e0e0;
            border-radius: 10px;
            padding: 20px;
            min-height: 300px;
            margin-bottom: 20px;
            background: #fafafa;
        }
        .module-item {
            background: white;
            border: 1px solid #e0e0e0;
            border-radius: 8px;
            padding: 15px;
            margin-bottom: 10px;
            display: flex;
            justify-content: space-between;
            align-items: center;
            cursor: move;
        }
        .module-item:hover {
            border-color: #3498db;
            box-shadow: 0 2px 8px rgba(52, 152, 219, 0.2);
        }
        .module-info {
            flex: 1;
        }
        .module-title {
            font-weight: 500;
            color: #2c3e50;
            margin-bottom: 5px;
        }
        .module-meta {
            font-size: 13px;
            color: #7f8c8d;
        }
        .module-actions {
            display: flex;
            gap: 8px;
        }
        .action-btn {
            background: none;
            border: 1px solid #ddd;
            border-radius: 4px;
            width: 32px;
            height: 32px;
            cursor: pointer;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 14px;
        }
        .action-btn:hover {
            background: #f8f9fa;
        }
        .action-btn.up:hover {
            border-color: #2ecc71;
            color: #2ecc71;
        }
        .action-btn.down:hover {
            border-color: #3498db;
            color: #3498db;
        }
        .action-btn.remove:hover {
            border-color: #e74c3c;
            color: #e74c3c;
        }
        .available-modules {
            max-height: 400px;
            overflow-y: auto;
            border: 1px solid #eee;
            border-radius: 8px;
            padding: 15px;
        }
        .available-module {
            padding: 12px 15px;
            border-bottom: 1px solid #eee;
            cursor: pointer;
            transition: background 0.2s;
        }
        .available-module:hover {
            background: #f8f9fa;
        }
        .available-module:last-child {
            border-bottom: none;
        }
        .add-module-btn {
            width: 100%;
            padding: 12px;
            background: #2ecc71;
            color: white;
            border: none;
            border-radius: 8px;
            font-size: 16px;
            cursor: pointer;
            margin-top: 15px;
            display: flex;
            align-items: center;
            justify-content: center;
            gap: 8px;
        }
        .add-module-btn:hover {
            background: #27ae60;
        }
        .empty-state {
            text-align: center;
            padding: 40px 20px;
            color: #95a5a6;
        }
        .checkbox-group {
            display: flex;
            align-items: center;
            gap: 10px;
            margin-top: 10px;
        }
        .checkbox-group input {
            width: auto;
        }
        .form-actions {
            display: flex;
            gap: 15px;
            margin-top: 40px;
            padding-top: 30px;
            border-top: 1px solid #eee;
        }
        .submit-btn {
            flex: 1;
            background: #2ecc71;
            color: white;
            border: none;
            padding: 16px;
            border-radius: 8px;
            font-size: 16px;
            font-weight: 500;
            cursor: pointer;
        }
        .submit-btn:hover {
            background: #27ae60;
        }
        .cancel-btn {
            flex: 1;
            background: #95a5a6;
            color: white;
            border: none;
            padding: 16px;
            border-radius: 8px;
            font-size: 16px;
            font-weight: 500;
            cursor: pointer;
            text-decoration: none;
            text-align: center;
        }
        .cancel-btn:hover {
            background: #7f8c8d;
        }
        .message {
            padding: 12px 16px;
            border-radius: 8px;
            margin-bottom: 20px;
            display: none;
        }
        .message.success {
            background: #d4edda;
            color: #155724;
            border: 1px solid #c3e6cb;
        }
        .message.error {
            background: #f8d7da;
            color: #721c24;
            border: 1px solid #f5c6cb;
        }
        .search-box {
            width: 100%;
            padding: 12px 16px;
            border: 1px solid #ddd;
            border-radius: 8px;
            margin-bottom: 15px;
            font-size: 16px;
        }
        .filter-buttons {
            display: flex;
            gap: 10px;
            margin-bottom: 15px;
        }
        .filter-btn {
            padding: 8px 16px;
            border: 1px solid #ddd;
            background: white;
            border-radius: 6px;
            cursor: pointer;
        }
        .filter-btn.active {
            background: #3498db;
            color: white;
            border-color: #3498db;
        }
    </style>
</head>
<body>
    <div class="create-container">
        <div class="page-header">
            <h1>➕ Создать новую лекцию</h1>
            <a href="/lectures" class="cancel-btn" style="width: auto; flex: none;">← Назад</a>
        </div>
        
        <div id="message" class="message"></div>
        
        <div class="two-column">
            <!-- Левая колонка: Форма лекции -->
            <div class="form-section">
                <h2>📝 Информация о лекции</h2>
                
                <form id="lectureForm">
                    <div class="form-group">
                        <label for="lectureTitle">Название лекции *</label>
                        <input type="text" id="lectureTitle" name="title" required 
                               placeholder="Например: 'Введение в математический анализ'">
                    </div>
                    
                    <div class="form-group">
                        <label for="lectureCourse">Предмет *</label>
                        <select id="lectureCourse" name="course" required>
                            <option value="">Выберите предмет</option>
                            <option value="Математический анализ">Математический анализ</option>
                            <option value="Линейная алгебра">Линейная алгебра</option>
                            <option value="Дискретная математика">Дискретная математика</option>
                            <option value="Теория вероятностей">Теория вероятностей</option>
                            <option value="Дифференциальные уравнения">Дифференциальные уравнения</option>
                            <option value="Экономика">Экономика</option>
                            <option value="Физика">Физика</option>
                        </select>
                    </div>
                    
                    <div class="form-group">
                        <label for="lectureDescription">Описание лекции</label>
                        <textarea id="lectureDescription" name="description" 
                                  placeholder="Опишите содержание лекции, цели обучения..."></textarea>
                    </div>
                    
                    <div class="form-group">
                        <label>Настройки лекции</label>
                        <div class="checkbox-group">
                            <input type="checkbox" id="allowBack" name="allow_back" checked>
                            <label for="allowBack">Разрешить студентам возвращаться к пройденным модулям</label>
                        </div>
                        <div class="checkbox-group">
                            <input type="checkbox" id="published" name="published" checked>
                            <label for="published">Опубликовать лекцию сразу после создания</label>
                        </div>
                    </div>
                </form>
                
                <div class="form-actions">
                    <button type="button" class="submit-btn" onclick="saveLecture()">💾 Сохранить лекцию</button>
                    <a href="/lectures" class="cancel-btn">Отмена</a>
                </div>
            </div>
            
            <!-- Правая колонка: Модули -->
            <div class="modules-section">
                <h2>📦 Состав лекции</h2>
                <p style="color: #7f8c8d; margin-bottom: 20px;">Добавьте модули и расположите их в нужном порядке</p>
                
                <div class="modules-list" id="modulesList">
                    <div class="empty-state" id="emptyState">
                        <h3>📭 Нет модулей</h3>
                        <p>Добавьте модули из библиотеки справа</p>
                    </div>
                </div>
                
                <h3>📚 Доступные модули</h3>
                <input type="text" class="search-box" id="moduleSearch" placeholder="Поиск модулей..." 
                       oninput="filterModules()">
                
                <div class="filter-buttons">
                    <button class="filter-btn active" onclick="setFilter('all')">Все</button>
                    <button class="filter-btn" onclick="setFilter('text')">📝 Текст</button>
                    <button class="filter-btn" onclick="setFilter('visual')">🎨 Визуал</button>
                    <button class="filter-btn" onclick="setFilter('question')">❓ Вопросы</button>
                    <button class="filter-btn" onclick="setFilter('test')">📋 Тесты</button>
                </div>
                
                <div class="available-modules" id="availableModules">
                    <!-- Модули загружаются через JavaScript -->
                    <div class="empty-state">
                        <h3>📭 Нет модулей</h3>
                        <p>Загрузка модулей...</p>
                    </div>
                </div>
            </div>
        </div>
    </div>
    
    <script>
        let selectedModules = [];
        let allModules = [];
        let currentFilter = 'all';
        
        // Загрузка доступных модулей
        async function loadAvailableModules() {
            try {
                const response = await fetch('/api/modules/available');
                allModules = await response.json();
                displayAvailableModules();
            } catch (error) {
                console.error('Error loading modules:', error);
            }
        }
        
        // Отображение доступных модулей
        function displayAvailableModules() {
            const container = document.getElementById('availableModules');
            const filtered = filterModuleList(allModules);
            
            if (filtered.length === 0) {
                container.innerHTML = '<div class="empty-state"><h3>📭 Нет модулей</h3><p>Создайте модули в библиотеке</p></div>';
                return;
            }
            
            let html = '';
            filtered.forEach(module => {
                const typeIcons = {
                    'text': '📝',
                    'visual': '🎨', 
                    'question': '❓',
                    'test': '📋'
                };
                
                html += '<div class="available-module" onclick="addModule(' + module.id + ')">' +
                    '<div class="module-info">' +
                    '<div class="module-title">' + module.title + '</div>' +
                    '<div class="module-meta">' +
                    (typeIcons[module.type] || '📄') + ' ' + module.type + ' • ' +
                    '📚 ' + module.course + ' • ' +
                    '👤 ' + module.author +
                    '</div>' +
                    '</div>' +
                    '</div>';
            });
            
            container.innerHTML = html;
        }
        
        // Фильтрация модулей
        function filterModuleList(modules) {
            let filtered = modules;
            
            // Поиск
            const searchTerm = document.getElementById('moduleSearch').value.toLowerCase();
            if (searchTerm) {
                filtered = filtered.filter(m => 
                    m.title.toLowerCase().includes(searchTerm) || 
                    m.description.toLowerCase().includes(searchTerm)
                );
            }
            
            // Фильтр по типу
            if (currentFilter !== 'all') {
                filtered = filtered.filter(m => m.type === currentFilter);
            }
            
            return filtered;
        }
        
        function setFilter(filter) {
            currentFilter = filter;
            document.querySelectorAll('.filter-btn').forEach(btn => {
                btn.classList.remove('active');
            });
            event.target.classList.add('active');
            displayAvailableModules();
        }
        
        function filterModules() {
            displayAvailableModules();
        }
        
        // Добавление модуля в лекцию
        function addModule(moduleId) {
            const module = allModules.find(m => m.id === moduleId);
            if (!module) return;
            
            // Проверяем, не добавлен ли уже
            if (selectedModules.some(m => m.id === moduleId)) {
                showMessage('Этот модуль уже добавлен в лекцию', 'error');
                return;
            }
            
            selectedModules.push({
                id: module.id,
                title: module.title,
                type: module.type,
                course: module.course,
                author: module.author
            });
            
            updateModulesList();
            showMessage('Модуль "' + module.title + '" добавлен', 'success');
        }
        
        // Обновление списка выбранных модулей
        function updateModulesList() {
            const container = document.getElementById('modulesList');
            const emptyState = document.getElementById('emptyState');
            
            if (selectedModules.length === 0) {
                container.innerHTML = '<div class="empty-state" id="emptyState"><h3>📭 Нет модулей</h3><p>Добавьте модули из библиотеки справа</p></div>';
                return;
            }
            
            let html = '';
            selectedModules.forEach((module, index) => {
                const typeIcons = {
                    'text': '📝',
                    'visual': '🎨',
                    'question': '❓',
                    'test': '📋'
                };
                
                html += '<div class="module-item" data-index="' + index + '">' +
                    '<div class="module-info">' +
                    '<div class="module-title">' + module.title + '</div>' +
                    '<div class="module-meta">' +
                    (typeIcons[module.type] || '📄') + ' ' + module.type + ' • ' +
                    '📚 ' + module.course + ' • ' +
                    '👤 ' + module.author +
                    '</div>' +
                    '</div>' +
                    '<div class="module-actions">' +
                    '<button class="action-btn up" onclick="moveModuleUp(' + index + ')" ' + (index === 0 ? 'disabled' : '') + '>↑</button>' +
                    '<button class="action-btn down" onclick="moveModuleDown(' + index + ')" ' + (index === selectedModules.length - 1 ? 'disabled' : '') + '>↓</button>' +
                    '<button class="action-btn remove" onclick="removeModule(' + index + ')">×</button>' +
                    '</div>' +
                    '</div>';
            });
            
            container.innerHTML = html;
        }
        
        // Перемещение модулей
        function moveModuleUp(index) {
            if (index <= 0) return;
            
            const temp = selectedModules[index];
            selectedModules[index] = selectedModules[index - 1];
            selectedModules[index - 1] = temp;
            
            updateModulesList();
        }
        
        function moveModuleDown(index) {
            if (index >= selectedModules.length - 1) return;
            
            const temp = selectedModules[index];
            selectedModules[index] = selectedModules[index + 1];
            selectedModules[index + 1] = temp;
            
            updateModulesList();
        }
        
        function removeModule(index) {
            if (confirm('Удалить этот модуль из лекции?')) {
                selectedModules.splice(index, 1);
                updateModulesList();
            }
        }
        
        // Сохранение лекции
        async function saveLecture() {
            const title = document.getElementById('lectureTitle').value;
            const course = document.getElementById('lectureCourse').value;
            const description = document.getElementById('lectureDescription').value;
            const allowBack = document.getElementById('allowBack').checked;
            const published = document.getElementById('published').checked;
            
            // Валидация
            if (!title || !course) {
                showMessage('Заполните обязательные поля: название и предмет', 'error');
                return;
            }
            
            if (selectedModules.length === 0) {
                showMessage('Добавьте хотя бы один модуль в лекцию', 'error');
                return;
            }
            
            const moduleIds = selectedModules.map(m => m.id);
            
            const lectureData = {
                title: title,
                course_name: course,
                description: description,
                module_ids: moduleIds,
                allow_back: allowBack,
                published: published
            };
            
            try {
                const response = await fetch('/api/lectures', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(lectureData)
                });
                
                const result = await response.json();
                
                if (response.ok) {
                    showMessage('✅ Лекция успешно создана! Перенаправление...', 'success');
                    
                    setTimeout(() => {
                        window.location.href = '/lectures/view/' + result.lecture.id;
                    }, 2000);
                    
                } else {
                    showMessage('❌ Ошибка: ' + (result.message || 'Не удалось создать лекцию'), 'error');
                }
            } catch (error) {
                showMessage('❌ Ошибка сети: ' + error.message, 'error');
            }
        }
        
        function showMessage(text, type) {
            const messageDiv = document.getElementById('message');
            messageDiv.textContent = text;
            messageDiv.className = 'message message-' + type;
            messageDiv.style.display = 'block';
            
            setTimeout(() => {
                messageDiv.style.display = 'none';
            }, 5000);
        }
        
        // Загружаем модули при старте
        window.addEventListener('DOMContentLoaded', loadAvailableModules);
    </script>
</body>
</html>`

	fmt.Fprintf(w, html)
}

// editLecturePageHandler показывает страницу редактирования лекции
func editLecturePageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := `<!DOCTYPE html>
<html>
<head>
    <title>Редактировать лекцию - VisualMath</title>
    <link rel="stylesheet" href="/static/css/style.css">
    <!-- MathJax для предпросмотра -->
    <script src="https://polyfill.io/v3/polyfill.min.js?features=es6"></script>
    <script id="MathJax-script" async src="https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-mml-chtml.js"></script>
    <script>
        MathJax = {
            tex: {
                inlineMath: [['$', '$'], ['\\(', '\\)']],
                displayMath: [['$$', '$$'], ['\\[', '\\]']]
            },
            svg: {
                fontCache: 'global'
            }
        };
    </script>
    <style>
        /* Стили такие же как в createLecturePageHandler */
        .create-container { max-width: 1200px; margin: 30px auto; padding: 0 20px; }
        .page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 30px; }
        .two-column { display: grid; grid-template-columns: 1fr 1fr; gap: 30px; margin-top: 20px; }
        .form-section, .modules-section { background: white; padding: 30px; border-radius: 12px; box-shadow: 0 4px 12px rgba(0,0,0,0.1); }
        .form-group { margin-bottom: 25px; }
        .form-group label { display: block; margin-bottom: 8px; font-weight: 500; color: #2c3e50; font-size: 16px; }
        .form-group input, .form-group select, .form-group textarea { width: 100%; padding: 14px 16px; border: 1px solid #ddd; border-radius: 8px; font-size: 16px; }
        .modules-list { border: 2px dashed #e0e0e0; border-radius: 10px; padding: 20px; min-height: 300px; margin-bottom: 20px; background: #fafafa; }
        .module-item { background: white; border: 1px solid #e0e0e0; border-radius: 8px; padding: 15px; margin-bottom: 10px; display: flex; justify-content: space-between; align-items: center; }
        .available-modules { max-height: 400px; overflow-y: auto; border: 1px solid #eee; border-radius: 8px; padding: 15px; }
        .available-module { padding: 12px 15px; border-bottom: 1px solid #eee; cursor: pointer; }
        .empty-state { text-align: center; padding: 40px 20px; color: #95a5a6; }
        .checkbox-group { display: flex; align-items: center; gap: 10px; margin-top: 10px; }
        .form-actions { display: flex; gap: 15px; margin-top: 40px; padding-top: 30px; border-top: 1px solid #eee; }
        .submit-btn { flex: 1; background: #3498db; color: white; border: none; padding: 16px; border-radius: 8px; font-size: 16px; cursor: pointer; }
        .cancel-btn { flex: 1; background: #95a5a6; color: white; border: none; padding: 16px; border-radius: 8px; font-size: 16px; cursor: pointer; text-decoration: none; text-align: center; }
        .message { padding: 12px 16px; border-radius: 8px; margin-bottom: 20px; display: none; }
        .message.success { background: #d4edda; color: #155724; }
        .message.error { background: #f8d7da; color: #721c24; }
        .search-box { width: 100%; padding: 12px 16px; border: 1px solid #ddd; border-radius: 8px; margin-bottom: 15px; font-size: 16px; }
        .filter-buttons { display: flex; gap: 10px; margin-bottom: 15px; }
        .filter-btn { padding: 8px 16px; border: 1px solid #ddd; background: white; border-radius: 6px; cursor: pointer; }
        .filter-btn.active { background: #3498db; color: white; border-color: #3498db; }
    </style>
</head>
<body>
    <div class="create-container">
        <div class="page-header">
            <h1>✏️ Редактировать лекцию</h1>
            <a href="/lectures/view/1" class="cancel-btn" style="width: auto; flex: none;">← Назад к просмотру</a>
        </div>
        
        <div id="message" class="message"></div>
        
        <div class="two-column">
            <!-- Левая колонка: Форма лекции -->
            <div class="form-section">
                <h2>📝 Информация о лекции</h2>
                
                <form id="lectureForm">
                    <div class="form-group">
                        <label for="lectureTitle">Название лекции *</label>
                        <input type="text" id="lectureTitle" name="title" required 
                               value="Введение в математический анализ">
                    </div>
                    
                    <div class="form-group">
                        <label for="lectureCourse">Предмет *</label>
                        <select id="lectureCourse" name="course" required>
                            <option value="Математический анализ" selected>Математический анализ</option>
                            <option value="Линейная алгебра">Линейная алгебра</option>
                            <option value="Дискретная математика">Дискретная математика</option>
                        </select>
                    </div>
                    
                    <div class="form-group">
                        <label for="lectureDescription">Описание лекции</label>
                        <textarea id="lectureDescription" name="description">Базовые понятия математического анализа: пределы, производные, интегралы.</textarea>
                    </div>
                    
                    <div class="form-group">
                        <label>Настройки лекции</label>
                        <div class="checkbox-group">
                            <input type="checkbox" id="allowBack" name="allow_back" checked>
                            <label for="allowBack">Разрешить студентам возвращаться к пройденным модулям</label>
                        </div>
                        <div class="checkbox-group">
                            <input type="checkbox" id="published" name="published" checked>
                            <label for="published">Опубликовать лекцию</label>
                        </div>
                    </div>
                </form>
                
                <div class="form-actions">
                    <button type="button" class="submit-btn" onclick="updateLecture()">💾 Сохранить изменения</button>
                    <a href="/lectures/view/1" class="cancel-btn">Отмена</a>
                    <button type="button" class="cancel-btn" style="background: #e74c3c;" onclick="deleteLecture()">🗑️ Удалить лекцию</button>
                </div>
            </div>
            
            <!-- Правая колонка: Модули -->
            <div class="modules-section">
                <h2>📦 Состав лекции</h2>
                <p style="color: #7f8c8d; margin-bottom: 20px;">Текущие модули в лекции (5 модулей)</p>
                
                <div class="modules-list" id="modulesList">
                    <!-- Модули предзаполняются через JavaScript -->
                </div>
                
                <h3>📚 Доступные модули</h3>
                <input type="text" class="search-box" id="moduleSearch" placeholder="Поиск модулей..." 
                       oninput="filterModules()">
                
                <div class="filter-buttons">
                    <button class="filter-btn active" onclick="setFilter('all')">Все</button>
                    <button class="filter-btn" onclick="setFilter('text')">📝 Текст</button>
                    <button class="filter-btn" onclick="setFilter('visual')">🎨 Визуал</button>
                    <button class="filter-btn" onclick="setFilter('question')">❓ Вопросы</button>
                    <button class="filter-btn" onclick="setFilter('test')">📋 Тесты</button>
                </div>
                
                <div class="available-modules" id="availableModules">
                    <!-- Модули загружаются через JavaScript -->
                </div>
            </div>
        </div>
    </div>
    
    <script>
        let selectedModules = [
            {id: 1, title: "Понятие предела функции", type: "text", course: "Математический анализ", author: "Иванов И.И."},
            {id: 2, title: "Производная функции", type: "text", course: "Математический анализ", author: "Иванов И.И."},
            {id: 3, title: "Проверка понимания производных", type: "question", course: "Математический анализ", author: "Иванов И.И."},
            {id: 4, title: "Графическое представление производной", type: "visual", course: "Математический анализ", author: "Иванов И.И."},
            {id: 5, title: "Итоговый тест по теме", type: "test", course: "Математический анализ", author: "Иванов И.И."}
        ];
        
        let allModules = [];
        let currentFilter = 'all';
        
        // Инициализация
        window.addEventListener('DOMContentLoaded', function() {
            updateModulesList();
            loadAvailableModules();
        });
        
        // Обновление списка модулей
        function updateModulesList() {
            const container = document.getElementById('modulesList');
            
            if (selectedModules.length === 0) {
                container.innerHTML = '<div class="empty-state"><h3>📭 Нет модулей</h3><p>Добавьте модули из библиотеки</p></div>';
                return;
            }
            
            let html = '';
            selectedModules.forEach((module, index) => {
                const typeIcons = {'text': '📝', 'visual': '🎨', 'question': '❓', 'test': '📋'};
                
                html += '<div class="module-item" data-index="' + index + '">' +
                    '<div class="module-info">' +
                    '<div class="module-title">' + module.title + '</div>' +
                    '<div class="module-meta">' +
                    (typeIcons[module.type] || '📄') + ' ' + module.type + ' • ' +
                    '📚 ' + module.course + ' • ' +
                    '👤 ' + module.author +
                    '</div>' +
                    '</div>' +
                    '<div class="module-actions">' +
                    '<button class="action-btn up" onclick="moveModuleUp(' + index + ')" ' + (index === 0 ? 'disabled' : '') + '>↑</button>' +
                    '<button class="action-btn down" onclick="moveModuleDown(' + index + ')" ' + (index === selectedModules.length - 1 ? 'disabled' : '') + '>↓</button>' +
                    '<button class="action-btn remove" onclick="removeModule(' + index + ')">×</button>' +
                    '</div>' +
                    '</div>';
            });
            
            container.innerHTML = html;
        }
        
        // Обновление лекции
        async function updateLecture() {
            const title = document.getElementById('lectureTitle').value;
            const course = document.getElementById('lectureCourse').value;
            const description = document.getElementById('lectureDescription').value;
            const allowBack = document.getElementById('allowBack').checked;
            const published = document.getElementById('published').checked;
            
            if (!title || !course) {
                showMessage('Заполните обязательные поля', 'error');
                return;
            }
            
            if (selectedModules.length === 0) {
                showMessage('Добавьте хотя бы один модуль', 'error');
                return;
            }
            
            const moduleIds = selectedModules.map(m => m.id);
            const lectureId = 1; // В реальности из URL
            
            const lectureData = {
                title: title,
                course_name: course,
                description: description,
                module_ids: moduleIds,
                allow_back: allowBack,
                published: published
            };
            
            try {
                const response = await fetch('/api/lectures/' + lectureId, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(lectureData)
                });
                
                const result = await response.json();
                
                if (response.ok) {
                    showMessage('✅ Изменения сохранены!', 'success');
                } else {
                    showMessage('❌ Ошибка: ' + (result.message || 'Не удалось сохранить'), 'error');
                }
            } catch (error) {
                showMessage('❌ Ошибка сети: ' + error.message, 'error');
            }
        }
        
        function deleteLecture() {
            if (confirm('Вы уверены, что хотите удалить эту лекцию? Это действие нельзя отменить.')) {
                const lectureId = 1;
                
                fetch('/api/lectures/' + lectureId, {
                    method: 'DELETE'
                })
                .then(response => response.json())
                .then(result => {
                    if (result.success) {
                        alert('Лекция удалена');
                        window.location.href = '/lectures';
                    } else {
                        alert('Ошибка: ' + (result.message || 'Не удалось удалить'));
                    }
                });
            }
        }
        
        function showMessage(text, type) {
            const messageDiv = document.getElementById('message');
            messageDiv.textContent = text;
            messageDiv.className = 'message message-' + type;
            messageDiv.style.display = 'block';
            
            setTimeout(() => {
                messageDiv.style.display = 'none';
            }, 5000);
        }
        
        // Остальные функции из createLecturePageHandler
        async function loadAvailableModules() {
            try {
                const response = await fetch('/api/modules/available');
                allModules = await response.json();
                displayAvailableModules();
            } catch (error) {
                console.error('Error loading modules:', error);
            }
        }
        
        function displayAvailableModules() {
            const container = document.getElementById('availableModules');
            const filtered = allModules.filter(m => !selectedModules.some(sm => sm.id === m.id));
            
            if (filtered.length === 0) {
                container.innerHTML = '<div class="empty-state"><h3>📭 Нет доступных модулей</h3><p>Все модули уже в лекции</p></div>';
                return;
            }
            
            let html = '';
            filtered.forEach(module => {
                const typeIcons = {'text': '📝', 'visual': '🎨', 'question': '❓', 'test': '📋'};
                
                html += '<div class="available-module" onclick="addModule(' + module.id + ')">' +
                    '<div class="module-info">' +
                    '<div class="module-title">' + module.title + '</div>' +
                    '<div class="module-meta">' +
                    (typeIcons[module.type] || '📄') + ' ' + module.type + ' • ' +
                    '📚 ' + module.course + ' • ' +
                    '👤 ' + module.author +
                    '</div>' +
                    '</div>' +
                    '</div>';
            });
            
            container.innerHTML = html;
        }
        
        function addModule(moduleId) {
            const module = allModules.find(m => m.id === moduleId);
            if (!module) return;
            
            selectedModules.push({
                id: module.id,
                title: module.title,
                type: module.type,
                course: module.course,
                author: module.author
            });
            
            updateModulesList();
            displayAvailableModules();
            showMessage('Модуль "' + module.title + '" добавлен', 'success');
        }
        
        function moveModuleUp(index) {
            if (index <= 0) return;
            [selectedModules[index], selectedModules[index-1]] = [selectedModules[index-1], selectedModules[index]];
            updateModulesList();
        }
        
        function moveModuleDown(index) {
            if (index >= selectedModules.length - 1) return;
            [selectedModules[index], selectedModules[index+1]] = [selectedModules[index+1], selectedModules[index]];
            updateModulesList();
        }
        
        function removeModule(index) {
            if (confirm('Удалить этот модуль из лекции?')) {
                selectedModules.splice(index, 1);
                updateModulesList();
                displayAvailableModules();
            }
        }
        
        function setFilter(filter) {
            currentFilter = filter;
            document.querySelectorAll('.filter-btn').forEach(btn => {
                btn.classList.remove('active');
            });
            event.target.classList.add('active');
            displayAvailableModules();
        }
        
        function filterModules() {
            displayAvailableModules();
        }
    </script>
</body>
</html>`

	fmt.Fprintf(w, html)
}

// viewLecturePageHandler показывает лекцию как единый документ для преподавателя
func viewLecturePageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := `<!DOCTYPE html>
<html>
<head>
    <title>Лекция: Введение в матанализ - VisualMath</title>
    <link rel="stylesheet" href="/static/css/style.css">
    <!-- MathJax для LaTeX -->
    <script src="https://polyfill.io/v3/polyfill.min.js?features=es6"></script>
    <script id="MathJax-script" async src="https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-mml-chtml.js"></script>
    <script>
        MathJax = {
            tex: {
                inlineMath: [['$', '$'], ['\\(', '\\)']],
                displayMath: [['$$', '$$'], ['\\[', '\\]']]
            },
            svg: {
                fontCache: 'global'
            }
        };
    </script>
    <style>
        .lecture-container {
            max-width: 1000px;
            margin: 0 auto;
            padding: 20px;
            background: white;
            box-shadow: 0 0 20px rgba(0,0,0,0.1);
            min-height: 100vh;
        }
        .lecture-header {
            text-align: center;
            padding-bottom: 30px;
            border-bottom: 2px solid #eee;
            margin-bottom: 40px;
        }
        .lecture-header h1 {
            color: #2c3e50;
            margin-bottom: 10px;
            font-size: 32px;
        }
        .lecture-meta {
            display: flex;
            justify-content: center;
            gap: 30px;
            color: #7f8c8d;
            margin-top: 15px;
        }
        .lecture-actions {
            display: flex;
            gap: 15px;
            justify-content: center;
            margin-top: 25px;
        }
        .btn {
            padding: 10px 20px;
            border: none;
            border-radius: 6px;
            cursor: pointer;
            font-size: 14px;
            text-decoration: none;
            display: inline-flex;
            align-items: center;
            gap: 8px;
        }
        .btn-edit {
            background: #3498db;
            color: white;
        }
        .btn-back {
            background: #95a5a6;
            color: white;
        }
        .btn-share {
            background: #2ecc71;
            color: white;
        }
        .module-container {
            margin-bottom: 50px;
            border-left: 4px solid #3498db;
            padding-left: 20px;
            position: relative;
        }
        .module-number {
            position: absolute;
            left: -15px;
            top: -10px;
            background: #3498db;
            color: white;
            width: 30px;
            height: 30px;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            font-weight: bold;
        }
        .module-header {
            display: flex;
            justify-content: space-between;
            align-items: flex-start;
            margin-bottom: 20px;
        }
        .module-title {
            font-size: 22px;
            color: #2c3e50;
            margin: 0;
        }
        .module-type {
            padding: 4px 12px;
            border-radius: 20px;
            font-size: 12px;
            font-weight: 500;
        }
        .type-text { background: #d4edda; color: #155724; }
        .type-visual { background: #d1ecf1; color: #0c5460; }
        .type-question { background: #fff3cd; color: #856404; }
        .type-test { background: #f8d7da; color: #721c24; }
        .module-content {
            background: #f8f9fa;
            padding: 25px;
            border-radius: 10px;
            margin-top: 15px;
        }
        .latex-content {
            font-family: "Times New Roman", Times, serif;
            font-size: 18px;
            line-height: 1.6;
        }
        .latex-content p {
            margin-bottom: 20px;
        }
        .question-block {
            background: white;
            border: 1px solid #ddd;
            border-radius: 8px;
            padding: 20px;
            margin: 15px 0;
        }
        .question-text {
            font-weight: 500;
            margin-bottom: 15px;
            color: #2c3e50;
        }
        .answer-option {
            padding: 10px 15px;
            margin: 8px 0;
            background: #f8f9fa;
            border-radius: 6px;
            border-left: 3px solid #3498db;
        }
        .correct-answer {
            border-left-color: #2ecc71;
            background: #d4edda;
        }
        .test-config {
            background: #f8f9fa;
            padding: 15px;
            border-radius: 8px;
            margin-top: 15px;
            font-family: monospace;
            font-size: 14px;
        }
        .image-placeholder {
            background: #e9ecef;
            border: 2px dashed #dee2e6;
            border-radius: 8px;
            padding: 40px;
            text-align: center;
            color: #6c757d;
            margin: 15px 0;
        }
        .lecture-footer {
            margin-top: 50px;
            padding-top: 30px;
            border-top: 1px solid #eee;
            text-align: center;
            color: #7f8c8d;
        }
        .navigation {
            display: flex;
            justify-content: space-between;
            margin-top: 40px;
            padding: 20px 0;
            border-top: 1px solid #eee;
        }
        .nav-btn {
            padding: 10px 20px;
            background: #f8f9fa;
            border: 1px solid #ddd;
            border-radius: 6px;
            cursor: pointer;
            text-decoration: none;
            color: #2c3e50;
        }
        .nav-btn:hover {
            background: #e9ecef;
        }
    </style>
</head>
<body style="background: #f8f9fa;">
    <div class="lecture-container">
        <!-- Заголовок лекции -->
        <div class="lecture-header">
            <h1>📚 Введение в математический анализ</h1>
            <p style="color: #7f8c8d; max-width: 800px; margin: 0 auto;">
                Базовые понятия математического анализа: пределы, производные, интегралы. 
                Лекция состоит из 5 модулей, расположенных в логической последовательности.
            </p>
            <div class="lecture-meta">
                <span>📚 Предмет: <strong>Математический анализ</strong></span>
                <span>👤 Автор: <strong>Иванов И.И.</strong></span>
                <span>📅 Создана: <strong>10 января 2024</strong></span>
                <span>📊 Модулей: <strong>5</strong></span>
            </div>
            <div class="lecture-actions">
                <a href="/lectures/edit/1" class="btn btn-edit">✏️ Редактировать лекцию</a>
                <a href="/lectures" class="btn btn-back">← Назад к списку</a>
                <button class="btn btn-share" onclick="shareLecture()">🔗 Поделиться</button>
            </div>
        </div>
        
        <!-- Модуль 1: Пределы -->
        <div class="module-container">
            <div class="module-number">1</div>
            <div class="module-header">
                <h2 class="module-title">Понятие предела функции</h2>
                <span class="module-type type-text">📝 Текстовый модуль</span>
            </div>
            <div class="module-content">
                <div class="latex-content">
                    <p>Предел функции — одно из основных понятий математического анализа.</p>
                    <p><strong>Определение:</strong> Число $A$ называется пределом функции $f(x)$ в точке $x_0$, если для любого $\epsilon > 0$ существует $\delta > 0$ такое, что для всех $x \neq x_0$, удовлетворяющих условию $|x - x_0| < \delta$, выполняется неравенство $|f(x) - A| < \epsilon$.</p>
                    <p>Записывается это так:</p>
                    <p>$$\lim_{x \to x_0} f(x) = A$$</p>
                    <p><strong>Пример 1:</strong> Найти предел:</p>
                    <p>$$\lim_{x \to 2} (3x + 1) = 3 \cdot 2 + 1 = 7$$</p>
                    <p><strong>Пример 2:</strong> Более сложный предел:</p>
                    <p>$$\lim_{x \to 0} \frac{\sin x}{x} = 1$$</p>
                    <div class="image-placeholder">
                        📈 Здесь будет график функции $\frac{\sin x}{x}$ при $x \to 0$
                    </div>
                </div>
            </div>
        </div>
        
        <!-- Модуль 2: Производные -->
        <div class="module-container">
            <div class="module-number">2</div>
            <div class="module-header">
                <h2 class="module-title">Производная функции</h2>
                <span class="module-type type-text">📝 Текстовый модуль</span>
            </div>
            <div class="module-content">
                <div class="latex-content">
                    <p>Производная функции $f(x)$ в точке $x_0$ определяется как предел отношения приращения функции к приращению аргумента:</p>
                    <p>$$f'(x_0) = \lim_{\Delta x \to 0} \frac{f(x_0 + \Delta x) - f(x_0)}{\Delta x}$$</p>
                    <p><strong>Геометрический смысл:</strong> Производная равна тангенсу угла наклона касательной к графику функции.</p>
                    <p><strong>Основные правила дифференцирования:</strong></p>
                    <ul>
                        <li>Производная константы: $(c)' = 0$</li>
                        <li>Производная суммы: $(f + g)' = f' + g'$</li>
                        <li>Производная произведения: $(fg)' = f'g + fg'$</li>
                        <li>Производная частного: $\left(\frac{f}{g}\right)' = \frac{f'g - fg'}{g^2}$</li>
                        <li>Производная сложной функции: $(f(g(x)))' = f'(g(x)) \cdot g'(x)$</li>
                    </ul>
                    <p><strong>Примеры:</strong></p>
                    <p>1. $f(x) = x^3 \Rightarrow f'(x) = 3x^2$</p>
                    <p>2. $f(x) = \sin x \Rightarrow f'(x) = \cos x$</p>
                    <p>3. $f(x) = e^x \Rightarrow f'(x) = e^x$</p>
                </div>
            </div>
        </div>
        
        <!-- Модуль 3: Вопросник -->
        <div class="module-container">
            <div class="module-number">3</div>
            <div class="module-header">
                <h2 class="module-title">Проверка понимания производных</h2>
                <span class="module-type type-question">❓ Вопросник</span>
            </div>
            <div class="module-content">
                <div class="question-block">
                    <div class="question-text">1. Чему равна производная функции $f(x) = 5x^2 + 3x - 7$?</div>
                    <div class="answer-option">A) $10x + 3$</div>
                    <div class="answer-option correct-answer">B) $10x + 3$</div>
                    <div class="answer-option">C) $5x + 3$</div>
                    <div class="answer-option">D) $10x$</div>
                </div>
                
                <div class="question-block">
                    <div class="question-text">2. Что показывает производная функции в точке?</div>
                    <div class="answer-option">A) Площадь под графиком</div>
                    <div class="answer-option correct-answer">B) Скорость изменения функции</div>
                    <div class="answer-option">C) Максимальное значение</div>
                    <div class="answer-option">D) Корень уравнения</div>
                </div>
                
                <div class="question-block">
                    <div class="question-text">3. Чему равна производная константы $c$?</div>
                    <div class="answer-option correct-answer">A) 0</div>
                    <div class="answer-option">B) 1</div>
                    <div class="answer-option">C) $c$</div>
                    <div class="answer-option">D) Не существует</div>
                </div>
            </div>
        </div>
        
        <!-- Модуль 4: Визуализация -->
        <div class="module-container">
            <div class="module-number">4</div>
            <div class="module-header">
                <h2 class="module-title">Графическое представление производной</h2>
                <span class="module-type type-visual">🎨 Визуализация</span>
            </div>
            <div class="module-content">
                <div class="image-placeholder" style="text-align: center;">
                    <div style="font-size: 20px; margin-bottom: 10px;">📊 График функции и её производной</div>
                    <p>В полной версии здесь будет интерактивный график, где можно:</p>
                    <ul style="text-align: left; display: inline-block;">
                        <li>Передвигать точку на графике $f(x) = x^2$</li>
                        <li>Видеть касательную в выбранной точке</li>
                        <li>Наблюдать значение производной $f'(x) = 2x$</li>
                        <li>Сравнивать графики функции и её производной</li>
                    </ul>
                    <div style="margin-top: 20px; padding: 15px; background: #e9ecef; border-radius: 6px;">
                        <strong>Пример визуализации:</strong><br>
                        Для $f(x) = x^2$:<br>
                        - В точке $x=1$: $f(1)=1$, $f'(1)=2$<br>
                        - В точке $x=2$: $f(2)=4$, $f'(2)=4$<br>
                        - В точке $x=3$: $f(3)=9$, $f'(3)=6$
                    </div>
                </div>
            </div>
        </div>
        
        <!-- Модуль 5: Тест -->
        <div class="module-container">
            <div class="module-number">5</div>
            <div class="module-header">
                <h2 class="module-title">Итоговый тест по теме</h2>
                <span class="module-type type-test">📋 Проверочный тест</span>
            </div>
            <div class="module-content">
                <div class="test-config">
                    <strong>Настройки теста:</strong><br>
                    - Время на выполнение: 30 минут<br>
                    - Количество вопросов: 10<br>
                    - Проходной балл: 70%<br>
                    - Разрешен возврат к вопросам: Да<br>
                    - Перемешивание вопросов: Да
                </div>
                
                <div class="question-block">
                    <div class="question-text">Пример тестового вопроса: Найдите производную $f(x) = \ln(x^2 + 1)$</div>
                    <div class="answer-option">A) $\frac{1}{x^2 + 1}$</div>
                    <div class="answer-option correct-answer">B) $\frac{2x}{x^2 + 1}$</div>
                    <div class="answer-option">C) $\frac{2}{x^2 + 1}$</div>
                    <div class="answer-option">D) $\frac{x}{x^2 + 1}$</div>
                </div>
                
                <p style="color: #7f8c8d; font-style: italic;">
                    * В реальном тесте будет 10 вопросов разных типов: выбор ответа, заполнение пропусков, 
                    сопоставление, вычисление производных.
                </p>
            </div>
        </div>
        
        <!-- Навигация -->
        <div class="navigation">
            <a href="/modules" class="nav-btn">📚 Посмотреть другие модули</a>
            <a href="/lectures/create" class="nav-btn">➕ Создать свою лекцию</a>
        </div>
        
        <!-- Подвал -->
        <div class="lecture-footer">
            <p>Лекция "Введение в математический анализ" • Автор: Иванов И.И. • VisualMath Platform © 2024</p>
            <p style="font-size: 14px; margin-top: 10px;">
                <a href="#" style="color: #3498db;">Экспорт в PDF</a> • 
                <a href="#" style="color: #3498db;">Поделиться со студентами</a> • 
                <a href="#" style="color: #3498db;">Статистика прохождений</a>
            </p>
        </div>
    </div>
    
    <script>
        function shareLecture() {
            const url = window.location.href;
            navigator.clipboard.writeText(url)
                .then(() => alert('Ссылка на лекцию скопирована в буфер обмена!'))
                .catch(err => console.error('Ошибка копирования:', err));
        }
        
        // Автоматическое обновление MathJax после загрузки
        window.addEventListener('DOMContentLoaded', function() {
            if (window.MathJax) {
                MathJax.typesetPromise();
            }
        });
    </script>
</body>
</html>`

	fmt.Fprintf(w, html)
}
