package handlers

import (
	"encoding/json"
	"fmt"
	//"html/template"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ModuleHandler struct {
	// Позже добавим БД
}

// ListModules показывает список всех модулей
func (h *ModuleHandler) ListModules(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Библиотека модулей - VisualMath</title>
    <link rel="stylesheet" href="/static/css/style.css">
    <style>
        .modules-container {
            max-width: 1200px;
            margin: 30px auto;
            padding: 0 20px;
        }
        .modules-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 30px;
            padding-bottom: 20px;
            border-bottom: 2px solid #eee;
        }
        .modules-header h1 {
            color: #2c3e50;
            margin: 0;
        }
        .create-btn {
            background: #2ecc71;
            color: white;
            border: none;
            padding: 12px 24px;
            border-radius: 6px;
            font-size: 16px;
            cursor: pointer;
            text-decoration: none;
            display: inline-block;
        }
        .create-btn:hover {
            background: #27ae60;
        }
        .search-filter {
            display: flex;
            gap: 15px;
            margin-bottom: 30px;
            padding: 20px;
            background: #f8f9fa;
            border-radius: 8px;
        }
        .search-box {
            flex: 1;
            padding: 12px;
            border: 1px solid #ddd;
            border-radius: 6px;
            font-size: 16px;
        }
        .filter-select {
            padding: 12px;
            border: 1px solid #ddd;
            border-radius: 6px;
            font-size: 16px;
            background: white;
            min-width: 200px;
        }
        .modules-grid {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
            gap: 25px;
        }
        .module-card {
            background: white;
            border-radius: 10px;
            overflow: hidden;
            box-shadow: 0 4px 12px rgba(0,0,0,0.1);
            transition: transform 0.3s, box-shadow 0.3s;
            border: 1px solid #eaeaea;
        }
        .module-card:hover {
            transform: translateY(-5px);
            box-shadow: 0 8px 20px rgba(0,0,0,0.15);
        }
        .module-header {
            padding: 20px;
            background: #f8f9fa;
            border-bottom: 1px solid #eee;
        }
        .module-title {
            margin: 0 0 10px 0;
            color: #2c3e50;
            font-size: 18px;
        }
        .module-body {
            padding: 20px;
        }
        .module-meta {
            display: flex;
            justify-content: space-between;
            color: #7f8c8d;
            font-size: 14px;
            margin-bottom: 15px;
        }
        .module-type {
            display: inline-block;
            padding: 4px 12px;
            border-radius: 20px;
            font-size: 12px;
            font-weight: 500;
            margin-bottom: 15px;
        }
        .type-text { background: #d4edda; color: #155724; }
        .type-visual { background: #d1ecf1; color: #0c5460; }
        .type-question { background: #fff3cd; color: #856404; }
        .type-test { background: #f8d7da; color: #721c24; }
        .module-actions {
            padding: 15px 20px;
            background: #f8f9fa;
            border-top: 1px solid #eee;
            display: flex;
            gap: 10px;
        }
        .action-btn {
            flex: 1;
            padding: 8px 12px;
            border: 1px solid #ddd;
            border-radius: 4px;
            background: white;
            cursor: pointer;
            font-size: 14px;
            text-align: center;
            text-decoration: none;
            color: #2c3e50;
        }
        .action-btn:hover {
            background: #f8f9fa;
        }
        .action-btn.primary {
            background: #3498db;
            color: white;
            border-color: #3498db;
        }
        .action-btn.primary:hover {
            background: #2980b9;
        }
        .empty-state {
            text-align: center;
            padding: 60px 20px;
            color: #7f8c8d;
        }
        .empty-state h3 {
            color: #95a5a6;
            margin-bottom: 10px;
        }
    </style>
</head>
<body>
    <div class="modules-container">
        <div class="modules-header">
            <h1>📚 Библиотека модулей</h1>
            <a href="/modules/create" class="create-btn">➕ Создать новый модуль</a>
        </div>
        
        <div class="search-filter">
            <input type="text" class="search-box" placeholder="Поиск модулей..." id="searchInput">
            <select class="filter-select" id="courseFilter">
                <option value="">Все предметы</option>
                <option value="Математический анализ">Математический анализ</option>
                <option value="Линейная алгебра">Линейная алгебра</option>
                <option value="Дискретная математика">Дискретная математика</option>
                <option value="Экономика">Экономика</option>
            </select>
            <select class="filter-select" id="typeFilter">
                <option value="">Все типы</option>
                <option value="text">Текстовый</option>
                <option value="visual">Визуальный</option>
                <option value="question">Вопросник</option>
                <option value="test">Проверочный</option>
            </select>
            <button class="create-btn" onclick="filterModules()">Применить</button>
        </div>
        
        <div class="modules-grid" id="modulesGrid">
            <!-- Модули будут загружены через JavaScript -->
            <div class="empty-state">
                <h3>📭 Нет модулей</h3>
                <p>Создайте свой первый модуль</p>
            </div>
        </div>
    </div>
    
    <script>
        // Загружаем модули
        async function loadModules() {
            try {
                const response = await fetch('/api/modules/list');
                const modules = await response.json();
                displayModules(modules);
            } catch (error) {
                console.error('Error loading modules:', error);
            }
        }
        
        function displayModules(modules) {
            const grid = document.getElementById('modulesGrid');
            
            if (!modules || modules.length === 0) {
                grid.innerHTML = '<div class="empty-state"><h3>📭 Нет модулей</h3><p>Создайте свой первый модуль</p><a href="/modules/create" class="create-btn" style="margin-top: 15px;">Создать модуль</a></div>';
                return;
            }
            
            let html = '';
            
            modules.forEach(module => {
                const typeLabels = {
                    'text': { name: 'Текстовый', class: 'type-text' },
                    'visual': { name: 'Визуальный', class: 'type-visual' },
                    'question': { name: 'Вопросник', class: 'type-question' },
                    'test': { name: 'Проверочный', class: 'type-test' }
                };
                
                const typeInfo = typeLabels[module.type] || { name: module.type, class: '' };
                
                html += '<div class="module-card">' +
                        '<div class="module-header">' +
                        '<h3 class="module-title">' + module.title + '</h3>' +
                        '<span class="module-type ' + typeInfo.class + '">' + typeInfo.name + '</span>' +
                        '</div>' +
                        '<div class="module-body">' +
                        '<div class="module-meta">' +
                        '<span>📚 ' + module.course + '</span>' +
                        '<span>👤 ' + module.author + '</span>' +
                        '</div>' +
                        '<p>' + (module.description || 'Описание отсутствует') + '</p>' +
                        '</div>' +
                        '<div class="module-actions">' +
                        '<a href="/modules/view/' + module.id + '" class="action-btn primary">Открыть</a>' +
                        '<a href="/modules/edit/' + module.id + '" class="action-btn">Редактировать</a>' +
                        '</div>' +
                        '</div>';
            });
            
            grid.innerHTML = html;
        }
        
        function filterModules() {
            const search = document.getElementById('searchInput').value.toLowerCase();
            const course = document.getElementById('courseFilter').value;
            const type = document.getElementById('typeFilter').value;
            
            // Здесь будет фильтрация
            console.log('Filtering:', { search, course, type });
            // В реальном приложении здесь будет запрос к API
        }
        
        // Загружаем модули при загрузке страницы
        window.addEventListener('DOMContentLoaded', loadModules);
    </script>
</body>
</html>`

	fmt.Fprintf(w, html)
}

// CreateModulePage показывает страницу создания модуля
func (h *ModuleHandler) CreateModulePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Создать модуль - VisualMath</title>
    <link rel="stylesheet" href="/static/css/style.css">
    <style>
        .create-container {
            max-width: 800px;
            margin: 30px auto;
            padding: 0 20px;
        }
        .create-header {
            margin-bottom: 30px;
            text-align: center;
        }
        .create-header h1 {
            color: #2c3e50;
            margin-bottom: 10px;
        }
        .create-header p {
            color: #7f8c8d;
            font-size: 16px;
        }
        .create-form {
            background: white;
            padding: 40px;
            border-radius: 12px;
            box-shadow: 0 8px 25px rgba(0,0,0,0.1);
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
            min-height: 150px;
            resize: vertical;
            font-family: monospace;
        }
        .module-type-selector {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
            gap: 15px;
            margin: 20px 0;
        }
        .type-option {
            padding: 25px 15px;
            border: 2px solid #e0e0e0;
            border-radius: 10px;
            text-align: center;
            cursor: pointer;
            transition: all 0.3s;
            background: white;
        }
        .type-option:hover {
            border-color: #3498db;
            background: #f8f9fa;
            transform: translateY(-2px);
        }
        .type-option.selected {
            border-color: #2ecc71;
            background: #f0f9f4;
            box-shadow: 0 4px 12px rgba(46, 204, 113, 0.2);
        }
        .type-option h4 {
            margin: 0 0 10px 0;
            color: #2c3e50;
            font-size: 16px;
        }
        .type-option p {
            margin: 0;
            color: #7f8c8d;
            font-size: 14px;
            line-height: 1.4;
        }
        .type-icon {
            font-size: 24px;
            margin-bottom: 10px;
            display: block;
        }
        .form-actions {
            display: flex;
            gap: 15px;
            margin-top: 40px;
            padding-top: 25px;
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
            transition: background 0.3s;
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
            transition: background 0.3s;
        }
        .cancel-btn:hover {
            background: #7f8c8d;
        }
        .content-section {
            display: none;
            animation: fadeIn 0.3s;
        }
        .content-section.active {
            display: block;
        }
        @keyframes fadeIn {
            from { opacity: 0; transform: translateY(10px); }
            to { opacity: 1; transform: translateY(0); }
        }
        .image-upload {
            border: 2px dashed #ddd;
            border-radius: 8px;
            padding: 40px;
            text-align: center;
            margin: 15px 0;
            cursor: pointer;
            transition: border-color 0.3s;
        }
        .image-upload:hover {
            border-color: #3498db;
        }
        .image-upload.dragover {
            border-color: #2ecc71;
            background: #f0f9f4;
        }
        .image-preview {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
            gap: 15px;
            margin-top: 20px;
        }
        .preview-item {
            position: relative;
            border-radius: 6px;
            overflow: hidden;
            border: 1px solid #ddd;
        }
        .preview-item img {
            width: 100%;
            height: 120px;
            object-fit: cover;
        }
        .remove-btn {
            position: absolute;
            top: 5px;
            right: 5px;
            background: #e74c3c;
            color: white;
            border: none;
            border-radius: 50%;
            width: 24px;
            height: 24px;
            cursor: pointer;
            font-size: 12px;
        }
        .latex-hint {
            background: #f8f9fa;
            padding: 15px;
            border-radius: 6px;
            margin-top: 10px;
            font-size: 14px;
            color: #7f8c8d;
        }
        .latex-hint code {
            background: #e9ecef;
            padding: 2px 6px;
            border-radius: 3px;
            font-family: monospace;
        }
    </style>
</head>
<body>
    <div class="create-container">
        <div class="create-header">
            <h1>➕ Создать новый модуль</h1>
            <p>Заполните форму для создания учебного модуля</p>
        </div>
        
        <form id="createModuleForm" class="create-form">
            <!-- Основная информация -->
            <div class="form-group">
                <label for="moduleTitle">Название модуля *</label>
                <input type="text" id="moduleTitle" name="title" required 
                       placeholder="Введите название модуля (например: 'Производные функций')">
            </div>
            
            <div class="form-group">
                <label for="moduleCourse">Предмет *</label>
                <select id="moduleCourse" name="course" required>
                    <option value="">Выберите предмет</option>
                    <option value="Математический анализ">Математический анализ</option>
                    <option value="Линейная алгебра и аналитическая геометрия">Линейная алгебра и аналитическая геометрия</option>
                    <option value="Дискретная математика">Дискретная математика</option>
                    <option value="Экономика">Экономика</option>
                    <option value="Физика">Физика</option>
                    <option value="Дифференциальные уравнения">Дифференциальные уравнения</option>
                    <option value="Литература">Литература</option>
                </select>
            </div>
            
            <div class="form-group">
                <label for="moduleDescription">Краткое описание</label>
                <textarea id="moduleDescription" name="description" 
                          placeholder="Опишите содержание модуля (необязательно)"></textarea>
            </div>
            
            <!-- Выбор типа модуля -->
            <div class="form-group">
                <label>Тип модуля *</label>
                <div class="module-type-selector">
                    <div class="type-option" data-type="text" onclick="selectModuleType('text')">
                        <span class="type-icon">📝</span>
                        <h4>Текстовый модуль</h4>
                        <p>Текст с формулами LaTeX, возможность добавления изображений</p>
                    </div>
                    <div class="type-option" data-type="visual" onclick="selectModuleType('visual')">
                        <span class="type-icon">🎨</span>
                        <h4>Визуальный модуль</h4>
                        <p>Графики, диаграммы, интерактивные визуализации</p>
                    </div>
                    <div class="type-option" data-type="question" onclick="selectModuleType('question')">
                        <span class="type-icon">❓</span>
                        <h4>Вопросник</h4>
                        <p>Вопросы с вариантами ответов, тесты</p>
                    </div>
                    <div class="type-option" data-type="test" onclick="selectModuleType('test')">
                        <span class="type-icon">📋</span>
                        <h4>Проверочный блок</h4>
                        <p>Контрольные работы, экзаменационные задания</p>
                    </div>
                </div>
                <input type="hidden" id="moduleType" name="type" required>
            </div>
            
            <!-- Контент в зависимости от типа -->
            <div id="contentArea">
                <!-- Динамически меняется -->
            </div>
            
            <!-- Кнопки -->
            <div class="form-actions">
                <button type="submit" class="submit-btn">Создать модуль</button>
                <a href="/modules" class="cancel-btn">Отмена</a>
            </div>
        </form>
        
        <!-- Сообщения -->
        <div id="message" style="margin-top: 20px; padding: 15px; border-radius: 8px; display: none;"></div>
    </div>
    
    <script>
        let selectedType = '';
        let uploadedImages = [];
        
        // Выбор типа модуля
        function selectModuleType(type) {
            selectedType = type;
            document.getElementById('moduleType').value = type;
            
            // Снимаем выделение со всех
            document.querySelectorAll('.type-option').forEach(opt => {
                opt.classList.remove('selected');
            });
            
            // Выделяем выбранный
            document.querySelector('[data-type="' + type + '"]').classList.add('selected');
            
            // Показываем соответствующий контент
            updateContentArea(type);
        }
        
        // Обновление области контента
        function updateContentArea(type) {
            const contentArea = document.getElementById('contentArea');
            let html = '';
            
            switch(type) {
                case 'text':
                    html = '<div class="form-group">' +
                           '<label for="contentText">Текст модуля *</label>' +
                           '<textarea id="contentText" name="content" rows="15" placeholder="Введите текст лекции..." required></textarea>' +
                           '<div class="latex-hint">' +
                           '<strong>Подсказка по LaTeX:</strong><br>' +
                           '• Формулы в строке: <code>$E = mc^2$</code><br>' +
                           '• Отдельные формулы: <code>$$\\int_a^b f(x)dx$$</code><br>' +
                           '• Греческие буквы: <code>$\\alpha, \\beta, \\gamma$</code>' +
                           '</div>' +
                           '</div>' +
                           '<div class="form-group">' +
                           '<label>Добавить иллюстрации</label>' +
                           '<div class="image-upload" id="imageUpload" onclick="document.getElementById(\'imageInput\').click()" ondragover="handleDragOver(event)" ondragleave="handleDragLeave(event)" ondrop="handleDrop(event)">' +
                           '<p>📁 Перетащите изображения сюда или нажмите для выбора</p>' +
                           '<p style="font-size: 14px; color: #95a5a6; margin-top: 10px;">Поддерживаются: JPG, PNG, GIF (макс. 5MB)</p>' +
                           '</div>' +
                           '<input type="file" id="imageInput" multiple accept="image/*" style="display: none;" onchange="handleImageSelect(event)">' +
                           '<div class="image-preview" id="imagePreview"></div>' +
                           '</div>';
                    break;
                    
                case 'visual':
                    html = '<div class="form-group">' +
                           '<label for="visualFile">Файл визуализации *</label>' +
                           '<input type="file" id="visualFile" name="visual_file" accept=".json,.xml,.svg,.png,.jpg" required>' +
                           '<p style="color: #7f8c8d; font-size: 14px; margin-top: 5px;">Поддерживаемые форматы: JSON, XML, SVG, PNG, JPG</p>' +
                           '</div>' +
                           '<div class="form-group">' +
                           '<label for="visualConfig">Конфигурация визуализации</label>' +
                           '<textarea id="visualConfig" name="config" rows="8" placeholder=\'{\n  "width": 800,\n  "height": 600,\n  "interactive": true,\n  "animation": false,\n  "controls": ["zoom", "pan"]\n}\'></textarea>' +
                           '<p style="color: #7f8c8d; font-size: 14px; margin-top: 5px;">Укажите параметры визуализации в формате JSON (опционально)</p>' +
                           '</div>';
                    break;
                    
                case 'question':
                    html = '<div class="form-group">' +
                           '<label for="questions">Вопросы и ответы *</label>' +
                           '<textarea id="questions" name="questions" rows="12" required placeholder=\'[\n  {\n    "question": "Что такое производная функции?",\n    "answers": [\n      "Скорость изменения функции",\n      "Площадь под графиком",\n      "Корень уравнения",\n      "Предел функции"\n    ],\n    "correct": 0,\n    "explanation": "Производная показывает скорость изменения функции в точке"\n  },\n  {\n    "question": "Чему равна производная константы?",\n    "answers": ["0", "1", "Сама константа", "Не существует"],\n    "correct": 0\n  }\n]\'></textarea>' +
                           '<p style="color: #7f8c8d; font-size: 14px; margin-top: 5px;">Формат: JSON массив объектов. Каждый вопрос должен содержать:<br>• <code>"question"</code> - текст вопроса<br>• <code>"answers"</code> - массив вариантов ответов<br>• <code>"correct"</code> - индекс правильного ответа (0, 1, 2...)<br>• <code>"explanation"</code> - объяснение (опционально)</p>' +
                           '</div>';
                    break;
                    
                case 'test':
                    html = '<div class="form-group">' +
                           '<label for="testConfig">Конфигурация теста *</label>' +
                           '<textarea id="testConfig" name="test_config" rows="10" required placeholder=\'{\n  "time_limit": 60,\n  "questions_count": 10,\n  "passing_score": 70,\n  "shuffle_questions": true,\n  "shuffle_answers": true,\n  "show_results": true,\n  "allow_retake": false,\n  "questions": [\n    {"id": 1, "points": 2},\n    {"id": 2, "points": 3}\n  ]\n}\'></textarea>' +
                           '<p style="color: #7f8c8d; font-size: 14px; margin-top: 5px;">Укажите параметры теста в JSON формате</p>' +
                           '</div>' +
                           '<div class="form-group">' +
                           '<label for="testSource">Источник вопросов</label>' +
                           '<select id="testSource" name="test_source">' +
                           '<option value="new">Создать новые вопросы</option>' +
                           '<option value="existing">Использовать существующие</option>' +
                           '<option value="mixed">Смешанный (новые + существующие)</option>' +
                           '</select>' +
                           '</div>';
                    break;
            }
            
            contentArea.innerHTML = html;
            
            // Инициализируем загрузку изображений для текстового модуля
            if (type === 'text') {
                initImageUpload();
            }
        }
        
        // Загрузка изображений
        function initImageUpload() {
            const imageInput = document.getElementById('imageInput');
            const imageUpload = document.getElementById('imageUpload');
            
            if (imageInput && imageUpload) {
                imageInput.onchange = handleImageSelect;
            }
        }
        
        function handleImageSelect(event) {
            const files = event.target.files;
            handleImages(files);
        }
        
        function handleDragOver(event) {
            event.preventDefault();
            event.currentTarget.classList.add('dragover');
        }
        
        function handleDragLeave(event) {
            event.currentTarget.classList.remove('dragover');
        }
        
        function handleDrop(event) {
            event.preventDefault();
            event.currentTarget.classList.remove('dragover');
            const files = event.dataTransfer.files;
            handleImages(files);
        }
        
        function handleImages(files) {
            const preview = document.getElementById('imagePreview');
            
            for (let file of files) {
                if (file.type.startsWith('image/')) {
                    const reader = new FileReader();
                    reader.onload = function(e) {
                        uploadedImages.push({
                            name: file.name,
                            data: e.target.result,
                            type: file.type
                        });
                        updateImagePreview();
                    };
                    reader.readAsDataURL(file);
                }
            }
        }
        
        function updateImagePreview() {
            const preview = document.getElementById('imagePreview');
            preview.innerHTML = '';
            
            uploadedImages.forEach((image, index) => {
                const item = document.createElement('div');
                item.className = 'preview-item';
                item.innerHTML = '<img src="' + image.data + '" alt="' + image.name + '">' +
                               '<button class="remove-btn" onclick="removeImage(' + index + ')">×</button>';
                preview.appendChild(item);
            });
        }
        
        function removeImage(index) {
            uploadedImages.splice(index, 1);
            updateImagePreview();
        }
        
        // Обработка формы
        document.getElementById('createModuleForm').addEventListener('submit', async function(e) {
            e.preventDefault();
            
            const messageDiv = document.getElementById('message');
            messageDiv.style.display = 'none';
            
            // Собираем данные
            const formData = {
                title: document.getElementById('moduleTitle').value,
                course: document.getElementById('moduleCourse').value,
                description: document.getElementById('moduleDescription').value,
                type: document.getElementById('moduleType').value,
                content: getModuleContent()
            };
            
            // Валидация
            if (!formData.title || !formData.course || !formData.type) {
                showMessage('Заполните все обязательные поля', 'error');
                return;
            }
            
            if (selectedType === 'text' && !formData.content.text) {
                showMessage('Введите текст модуля', 'error');
                return;
            }
            
            try {
                const response = await fetch('/api/modules', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(formData)
                });
                
                const result = await response.json();
                
                if (response.ok) {
                    showMessage('✅ Модуль успешно создан! Перенаправление...', 'success');
                    
                    setTimeout(() => {
                        window.location.href = '/modules';
                    }, 2000);
                    
                } else {
                    showMessage('❌ Ошибка: ' + (result.message || 'Не удалось создать модуль'), 'error');
                }
            } catch (error) {
                showMessage('❌ Ошибка сети: ' + error.message, 'error');
            }
        });
        
        function getModuleContent() {
            switch(selectedType) {
                case 'text':
                    return {
                        text: document.getElementById('contentText').value,
                        images: uploadedImages
                    };
                case 'question':
                    try {
                        return JSON.parse(document.getElementById('questions').value);
                    } catch {
                        return { error: 'Неверный JSON формат' };
                    }
                case 'test':
                    try {
                        return JSON.parse(document.getElementById('testConfig').value);
                    } catch {
                        return { error: 'Неверный JSON формат' };
                    }
                case 'visual':
                    // Для визуального модуля нужно обработать файл
                    return {
                        file: 'visual-file',
                        config: document.getElementById('visualConfig').value
                    };
                default:
                    return {};
            }
        }
        
        function showMessage(text, type) {
            const messageDiv = document.getElementById('message');
            messageDiv.textContent = text;
            messageDiv.style.background = type === 'success' ? '#d4edda' : '#f8d7da';
            messageDiv.style.color = type === 'success' ? '#155724' : '#721c24';
            messageDiv.style.display = 'block';
            messageDiv.style.border = '1px solid ' + (type === 'success' ? '#c3e6cb' : '#f5c6cb');
        }
        
        // Выбираем текстовый модуль по умолчанию
        window.addEventListener('DOMContentLoaded', function() {
            selectModuleType('text');
        });
    </script>
</body>
</html>`

	fmt.Fprintf(w, html)
}

// CreateModule обрабатывает создание модуля
func (h *ModuleHandler) CreateModule(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Title       string      `json:"title"`
		Course      string      `json:"course"`
		Description string      `json:"description"`
		Type        string      `json:"type"`
		Content     interface{} `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Валидация
	if request.Title == "" || request.Course == "" || request.Type == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Здесь будет сохранение в БД
	// Пока просто возвращаем успех

	response := map[string]interface{}{
		"success":  true,
		"message":  "Module created successfully",
		"module": map[string]interface{}{
			"id":          1,
			"title":       request.Title,
			"course":      request.Course,
			"description": request.Description,
			"type":        request.Type,
			"content":     request.Content,
			"created_at":  "2024-01-01 12:00:00",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetModule возвращает информацию о модуле
func (h *ModuleHandler) GetModule(w http.ResponseWriter, r *http.Request) {
	moduleID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(moduleID)
	if err != nil {
		http.Error(w, "Invalid module ID", http.StatusBadRequest)
		return
	}

	// Здесь будет получение из БД
	// Пока возвращаем тестовые данные
	module := map[string]interface{}{
		"id":          id,
		"title":       "Пример модуля",
		"course":      "Математический анализ",
		"description": "Описание модуля",
		"type":        "text",
		"content": map[string]interface{}{
			"text": "Текст модуля с формулами $E=mc^2$",
		},
		"author":      "Иванов И.И.",
		"created_at":  "2024-01-01 12:00:00",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(module)
}

// UpdateModule обновляет модуль
func (h *ModuleHandler) UpdateModule(w http.ResponseWriter, r *http.Request) {
	moduleID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(moduleID)
	if err != nil {
		http.Error(w, "Invalid module ID", http.StatusBadRequest)
		return
	}

	var request struct {
		Title       string      `json:"title"`
		Course      string      `json:"course"`
		Description string      `json:"description"`
		Type        string      `json:"type"`
		Content     interface{} `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Здесь будет обновление в БД

	response := map[string]interface{}{
		"success":  true,
		"message":  "Module updated successfully",
		"module_id": id,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteModule удаляет модуль
func (h *ModuleHandler) DeleteModule(w http.ResponseWriter, r *http.Request) {
	moduleID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(moduleID)
	if err != nil {
		http.Error(w, "Invalid module ID", http.StatusBadRequest)
		return
	}

	// Здесь будет удаление из БД

	response := map[string]interface{}{
		"success":  true,
		"message":  "Module deleted successfully",
		"module_id": id,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListModulesAPI возвращает список модулей для API
func (h *ModuleHandler) ListModulesAPI(w http.ResponseWriter, r *http.Request) {
	// Здесь будет получение из БД
	// Пока возвращаем тестовые данные
	modules := []map[string]interface{}{
		{
			"id":          1,
			"title":       "Производные функций",
			"course":      "Математический анализ",
			"description": "Основы дифференцирования",
			"type":        "text",
			"author":      "Иванов И.И.",
			"created_at":  "2024-01-01 12:00:00",
		},
		{
			"id":          2,
			"title":       "Матрицы и определители",
			"course":      "Линейная алгебра",
			"description": "Работа с матрицами",
			"type":        "visual",
			"author":      "Петров П.П.",
			"created_at":  "2024-01-02 14:30:00",
		},
		{
			"id":          3,
			"title":       "Тест по теории вероятностей",
			"course":      "Теория вероятностей",
			"description": "Контрольный тест",
			"type":        "test",
			"author":      "Сидоров С.С.",
			"created_at":  "2024-01-03 10:15:00",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(modules)
}

// ViewModulePage показывает страницу просмотра модуля
func (h *ModuleHandler) ViewModulePage(w http.ResponseWriter, r *http.Request) {
	moduleID := chi.URLParam(r, "id")
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Просмотр модуля - VisualMath</title>
    <link rel="stylesheet" href="/static/css/style.css">
    <style>
        .view-container {
            max-width: 1000px;
            margin: 30px auto;
            padding: 0 20px;
        }
        .view-header {
            margin-bottom: 30px;
            padding-bottom: 20px;
            border-bottom: 2px solid #eee;
        }
        .view-header h1 {
            color: #2c3e50;
            margin-bottom: 10px;
        }
        .module-info {
            display: flex;
            gap: 20px;
            margin-bottom: 20px;
            color: #7f8c8d;
        }
        .module-content {
            background: white;
            padding: 40px;
            border-radius: 12px;
            box-shadow: 0 4px 12px rgba(0,0,0,0.1);
            margin-bottom: 30px;
        }
        .module-actions {
            display: flex;
            gap: 15px;
            margin-top: 30px;
        }
        .btn {
            padding: 12px 24px;
            border: none;
            border-radius: 6px;
            font-size: 16px;
            cursor: pointer;
            text-decoration: none;
            display: inline-block;
        }
        .btn-edit {
            background: #3498db;
            color: white;
        }
        .btn-back {
            background: #95a5a6;
            color: white;
        }
        .btn-delete {
            background: #e74c3c;
            color: white;
        }
        .latex-content {
            font-family: "Times New Roman", Times, serif;
            font-size: 18px;
            line-height: 1.6;
        }
    </style>
</head>
<body>
    <div class="view-container">
        <div class="view-header">
            <h1>Просмотр модуля #` + moduleID + `</h1>
            <div class="module-info">
                <span>📚 Математический анализ</span>
                <span>👤 Иванов И.И.</span>
                <span>📅 2024-01-01</span>
            </div>
        </div>
        
        <div class="module-content">
            <h2>Производные функций</h2>
            <div class="latex-content">
                <p>Производная функции $f(x)$ в точке $x_0$ определяется как предел:</p>
                <p>$$f'(x_0) = \lim_{\Delta x \to 0} \frac{f(x_0 + \Delta x) - f(x_0)}{\Delta x}$$</p>
                <p>Это основное понятие дифференциального исчисления.</p>
            </div>
        </div>
        
        <div class="module-actions">
            <a href="/modules/edit/` + moduleID + `" class="btn btn-edit">Редактировать</a>
            <a href="/modules" class="btn btn-back">Назад к списку</a>
            <button class="btn btn-delete" onclick="deleteModule(` + moduleID + `)">Удалить</button>
        </div>
    </div>
    
    <script>
        function deleteModule(id) {
            if (confirm('Вы уверены, что хотите удалить этот модуль?')) {
                fetch('/api/modules/' + id, {
                    method: 'DELETE'
                })
                .then(response => response.json())
                .then(data => {
                    if (data.success) {
                        alert('Модуль удален');
                        window.location.href = '/modules';
                    } else {
                        alert('Ошибка при удалении: ' + data.message);
                    }
                });
            }
        }
    </script>
</body>
</html>`

	fmt.Fprintf(w, html)
}

// EditModulePage показывает страницу редактирования модуля
func (h *ModuleHandler) EditModulePage(w http.ResponseWriter, r *http.Request) {
	moduleID := chi.URLParam(r, "id")
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Редактирование модуля - VisualMath</title>
    <link rel="stylesheet" href="/static/css/style.css">
    <style>
        .edit-container {
            max-width: 800px;
            margin: 30px auto;
            padding: 0 20px;
        }
        .edit-header {
            margin-bottom: 30px;
            text-align: center;
        }
        .edit-header h1 {
            color: #2c3e50;
            margin-bottom: 10px;
        }
        .edit-form {
            background: white;
            padding: 40px;
            border-radius: 12px;
            box-shadow: 0 8px 25px rgba(0,0,0,0.1);
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
        }
        .form-actions {
            display: flex;
            gap: 15px;
            margin-top: 40px;
            padding-top: 25px;
            border-top: 1px solid #eee;
        }
        .submit-btn {
            flex: 1;
            background: #3498db;
            color: white;
            border: none;
            padding: 16px;
            border-radius: 8px;
            font-size: 16px;
            cursor: pointer;
        }
        .cancel-btn {
            flex: 1;
            background: #95a5a6;
            color: white;
            border: none;
            padding: 16px;
            border-radius: 8px;
            font-size: 16px;
            cursor: pointer;
            text-decoration: none;
            text-align: center;
        }
    </style>
</head>
<body>
    <div class="edit-container">
        <div class="edit-header">
            <h1>✏️ Редактирование модуля #` + moduleID + `</h1>
            <p>Внесите изменения в учебный модуль</p>
        </div>
        
        <form id="editModuleForm" class="edit-form">
            <div class="form-group">
                <label for="editTitle">Название модуля *</label>
                <input type="text" id="editTitle" name="title" value="Производные функций" required>
            </div>
            
            <div class="form-group">
                <label for="editCourse">Предмет *</label>
                <select id="editCourse" name="course" required>
                    <option value="Математический анализ" selected>Математический анализ</option>
                    <option value="Линейная алгебра">Линейная алгебра</option>
                    <option value="Дискретная математика">Дискретная математика</option>
                </select>
            </div>
            
            <div class="form-group">
                <label for="editDescription">Краткое описание</label>
                <textarea id="editDescription" name="description">Основы дифференцирования</textarea>
            </div>
            
            <div class="form-group">
                <label for="editContent">Содержание модуля *</label>
                <textarea id="editContent" name="content" rows="15" required>Производная функции $f(x)$ в точке $x_0$ определяется как предел:

$$f'(x_0) = \lim_{\Delta x \to 0} \frac{f(x_0 + \Delta x) - f(x_0)}{\Delta x}$$

Это основное понятие дифференциального исчисления.</textarea>
            </div>
            
            <div class="form-actions">
                <button type="submit" class="submit-btn">Сохранить изменения</button>
                <a href="/modules/view/` + moduleID + `" class="cancel-btn">Отмена</a>
            </div>
        </form>
        
        <div id="message" style="margin-top: 20px; padding: 15px; border-radius: 8px; display: none;"></div>
    </div>
    
    <script>
        document.getElementById('editModuleForm').addEventListener('submit', async function(e) {
            e.preventDefault();
            
            const formData = {
                title: document.getElementById('editTitle').value,
                course: document.getElementById('editCourse').value,
                description: document.getElementById('editDescription').value,
                content: document.getElementById('editContent').value
            };
            
            try {
                const response = await fetch('/api/modules/` + moduleID + `', {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(formData)
                });
                
                const result = await response.json();
                
                if (response.ok) {
                    document.getElementById('message').textContent = '✅ Изменения сохранены!';
                    document.getElementById('message').style.background = '#d4edda';
                    document.getElementById('message').style.color = '#155724';
                    document.getElementById('message').style.display = 'block';
                } else {
                    document.getElementById('message').textContent = '❌ Ошибка: ' + (result.message || 'Не удалось сохранить');
                    document.getElementById('message').style.background = '#f8d7da';
                    document.getElementById('message').style.color = '#721c24';
                    document.getElementById('message').style.display = 'block';
                }
            } catch (error) {
                document.getElementById('message').textContent = '❌ Ошибка сети: ' + error.message;
                document.getElementById('message').style.background = '#f8d7da';
                document.getElementById('message').style.color = '#721c24';
                document.getElementById('message').style.display = 'block';
            }
        });
    </script>
</body>
</html>`

	fmt.Fprintf(w, html)
}