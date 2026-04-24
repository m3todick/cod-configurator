# ЦОД Конфигуратор
**Администрация Константиновского района Ростовской области**

Информационная система для подбора, конфигурирования и расчёта стоимости оборудования центра обработки данных.

---

## Стек технологий

| Слой | Технологии |
|------|-----------|
| Backend | Go 1.21+, `net/http` |
| База данных | PostgreSQL (Supabase) |
| Frontend | HTML5, CSS3, Vanilla JavaScript |
| Аутентификация | Cookie-сессии (bcrypt) |
| Шрифты | Playfair Display + Raleway (Google Fonts) |

---

## Быстрый старт

### 1. Настройте переменные окружения

Скопируйте `.env.example` в `.env` и заполните значения:

```bash
cp .env.example .env
```

`.env`:
```env
SUPABASE_DB_URL=postgres://user:password@host:5432/dbname?sslmode=require
PORT=8080
ROLE_CHANGE_PASSWORD=your_secret_role_password
```

> `ROLE_CHANGE_PASSWORD` — специальный пароль, который требуется при повышении или понижении роли пользователя через интерфейс.

### 2. Запустите сервер

```bash
go run .
```

Или соберите бинарник:

```bash
go build -o cod-server .
./cod-server
```

Сервер доступен по адресу: **http://localhost:8080**

---

## Демо-доступ

| Роль | Email | Пароль |
|------|-------|--------|
| Администратор | admin@konst-adm.ru | admin123 |
| Пользователь | user@konst-adm.ru | user123 |

---

## Структура проекта

```
cod-configurator/
│
├── main.go                      # Точка входа — загрузка .env, подключение к БД, запуск сервера
├── server.go                    # HTTP-сервер, маршрутизация, вспомогательные хелперы
├── models.go                    # Структуры данных (User, Component, Configuration и др.)
├── store.go                     # Слой работы с PostgreSQL (CRUD-операции)
├── crypto.go                    # Хеширование паролей (bcrypt) и управление сессиями
├── dotenv.go                    # Загрузка переменных окружения из файла .env
│
├── handlers_auth.go             # API: /api/auth/login, /api/auth/logout, /api/auth/me, /api/auth/register
├── handlers_components.go       # API: /api/components (CRUD оборудования, только admin)
├── handlers_configurations.go   # API: /api/configurations (сохранение, экспорт, импорт)
├── handlers_calculator.go       # API: /api/calculator/estimate (расчёт стоимости)
├── handlers_users.go            # API: /api/users, /api/profile (управление пользователями)
│
├── go.mod                       # Go-модуль и зависимости
├── go.sum                       # Контрольные суммы зависимостей
├── .env.example                 # Шаблон переменных окружения
├── run.sh                       # Скрипт запуска
│
└── web/
    ├── templates/               # HTML-страницы (рендерятся как статика)
    │   ├── index.html           # Главная страница
    │   ├── login.html           # Вход в систему
    │   ├── register.html        # Регистрация нового аккаунта
    │   ├── cabinet.html         # Личный кабинет (конфигурации, профиль, пользователи)
    │   ├── configurator.html    # Конфигуратор оборудования
    │   ├── calculator.html      # Калькулятор стоимости
    │   ├── admin.html           # Панель администратора
    │   ├── terms.html           # Условия использования
    │   └── privacy.html         # Политика конфиденциальности
    └── static/
        ├── css/
        │   └── main.css         # Главная таблица стилей (CSS-переменные, компоненты, адаптив)
        ├── js/
        │   └── common.js        # Общие JS-утилиты (http-клиент, toast, модалки, навбар)
        └── img/
            ├── server1.svg      # Иконки оборудования (SVG)
            ├── server2.svg
            ├── storage.svg
            ├── network.svg
            ├── cooling.svg
            ├── security.svg
            └── ups.svg
```

---

## API Эндпоинты

### Аутентификация
| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/api/auth/login` | Вход в систему |
| POST | `/api/auth/logout` | Выход |
| GET  | `/api/auth/me` | Текущий пользователь |
| POST | `/api/auth/register` | Регистрация нового пользователя |

### Оборудование
| Метод | URL | Описание |
|-------|-----|----------|
| GET    | `/api/components` | Список оборудования |
| POST   | `/api/components` | Добавить позицию (admin) |
| PUT    | `/api/components/:id` | Обновить позицию (admin) |
| DELETE | `/api/components/:id` | Удалить позицию (admin) |

### Конфигурации
| Метод | URL | Описание |
|-------|-----|----------|
| GET    | `/api/configurations` | Список конфигураций пользователя |
| POST   | `/api/configurations` | Создать конфигурацию |
| PUT    | `/api/configurations/:id` | Обновить конфигурацию |
| DELETE | `/api/configurations/:id` | Удалить конфигурацию |
| GET    | `/api/configurations/export/:id` | Экспорт в JSON |
| POST   | `/api/configurations/import` | Импорт из JSON |

### Прочее
| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/api/calculator/estimate` | Расчёт стоимости |
| PUT  | `/api/profile` | Обновление профиля / смена пароля |
| GET  | `/api/users` | Список пользователей (admin) |
| PUT  | `/api/users/:id` | Редактировать пользователя (admin) |
| DELETE | `/api/users/:id` | Удалить пользователя (admin) |

---

## Страницы

| URL | Описание |
|-----|----------|
| `/` | Главная страница |
| `/login` | Вход |
| `/register` | Регистрация |
| `/cabinet` | Личный кабинет |
| `/configurator` | Конфигуратор оборудования |
| `/calculator` | Калькулятор стоимости |
| `/admin` | Панель администратора |
| `/terms` | Условия использования |
| `/privacy` | Политика конфиденциальности |

---

## Функциональность

- Главная страница с описанием системы и уровнями надёжности TIA-942
- Регистрация и аутентификация через cookie-сессии
- Личный кабинет: управление конфигурациями, редактирование профиля, смена пароля
- Конфигуратор с фильтрацией по категориям, поиском и добавлением в корзину
- Калькулятор стоимости с интерактивными слайдерами
- Экспорт / импорт конфигураций в формат JSON
- Панель администратора: CRUD оборудования, просмотр всех конфигураций
- Управление пользователями: редактирование, удаление, смена роли
- **Защита смены роли:** повышение/понижение пользователя требует ввода специального пароля (`ROLE_CHANGE_PASSWORD`)
- Страницы «Условия использования» и «Политика конфиденциальности»
- Поддержка `.env` файла для хранения конфигурации
- Адаптивная вёрстка, анимации, SVG-иконки оборудования
