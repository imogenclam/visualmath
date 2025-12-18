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
    
    // ============ ДОБАВЬТЕ ЭТИ МАРШРУТЫ ============
    // Маршруты модулей
    r.Get("/modules", moduleHandler.ListModules)                    // Список модулей
    r.Get("/modules/create", moduleHandler.CreateModulePage)        // Страница создания
    r.Get("/modules/view/{id}", moduleHandler.ViewModulePage)       // Просмотр модуля
    r.Get("/modules/edit/{id}", moduleHandler.EditModulePage)       // Редактирование модуля
    
    // API endpoints для модулей
    r.Get("/api/modules/list", moduleHandler.ListModulesAPI)        // API: список модулей
    r.Post("/api/modules", moduleHandler.CreateModule)              // API: создание модуля
    r.Get("/api/modules/{id}", moduleHandler.GetModule)             // API: получить модуль
    r.Put("/api/modules/{id}", moduleHandler.UpdateModule)          // API: обновить модуль
    r.Delete("/api/modules/{id}", moduleHandler.DeleteModule)       // API: удалить модуль
    // ============ КОНЕЦ ДОБАВЛЕНИЯ ============
    
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
    
    html := `
    <!DOCTYPE html>
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
    </html>
    `
    
    fmt.Fprintf(w, html)
}

// loginPageHandler обрабатывает страницу входа
func loginPageHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    
    html := `
    <!DOCTYPE html>
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
    </html>
    `
    
    fmt.Fprintf(w, html)
}

// registerPageHandler обрабатывает страницу регистрации
func registerPageHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    
    html := `
    <!DOCTYPE html>
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
            .form-row {
                display: grid;
                grid-template-columns: 1fr 1fr;
                gap: 15px;
                margin-bottom: 15px;
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
                <h1>📝 Регистрация</h1>
                <p>Создайте аккаунт VisualMath</p>
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
                <div class="form-row">
                    <div class="form-group">
                        <label for="login">Логин *</label>
                        <input type="text" id="login" name="login" placeholder="Придумайте логин" required>
                    </div>
                    <div class="form-group">
                        <label for="email">Email *</label>
                        <input type="email" id="email" name="email" placeholder="Ваш email" required>
                    </div>
                </div>
                
                <div class="form-row">
                    <div class="form-group">
                        <label for="full_name">ФИО *</label>
                        <input type="text" id="full_name" name="full_name" placeholder="Иванов Иван Иванович" required>
                    </div>
                    <div class="form-group">
                        <label for="user_type">Тип пользователя *</label>
                        <select id="user_type" name="user_type" required>
                            <option value="">Выберите тип</option>
                            <option value="student">Студент</option>
                            <option value="teacher">Преподаватель</option>
                        </select>
                    </div>
                </div>
                
                <div class="form-group">
                    <label for="password">Пароль *</label>
                    <input type="password" id="password" name="password" placeholder="Придумайте пароль" required>
                    <div class="info-text">Минимум 8 символов</div>
                </div>
                
                <div class="form-group" id="groupField" style="display: none;">
                    <label for="group_number">Номер группы *</label>
                    <input type="text" id="group_number" name="group_number" placeholder="Например: ИУ6-32Б">
                    <div class="info-text">Только для студентов</div>
                </div>
                
                <button type="submit" class="submit-btn">Создать аккаунт</button>
            </form>
            
            <div class="form-links">
                <a href="/login">Уже есть аккаунт? Войти</a>
                <a href="/">На главную страницу</a>
            </div>
        </div>

        <script>
            // Показываем поле группы только для студентов
            document.getElementById('user_type').addEventListener('change', function() {
                const groupField = document.getElementById('groupField');
                const groupInput = document.getElementById('group_number');
                
                if (this.value === 'student') {
                    groupField.style.display = 'block';
                    groupInput.required = true;
                } else {
                    groupField.style.display = 'none';
                    groupInput.required = false;
                }
            });
            
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
    </html>
    `
    
    fmt.Fprintf(w, html)
}

// dashboardHandler показывает личный кабинет
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    
    html := `
    <!DOCTYPE html>
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
    </html>
    `
    
    fmt.Fprintf(w, html)
}

// testHandler для проверки работы
func testHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    html := `
    <!DOCTYPE html>
    <html>
    <head><title>Тест</title></head>
    <body style="padding: 50px; text-align: center;">
        <h1 style="color: green;">✅ Тест пройден!</h1>
        <p>Сервер работает правильно</p>
        <p><a href="/">Вернуться на главную</a></p>
    </body>
    </html>
    `
    fmt.Fprintf(w, html)
}