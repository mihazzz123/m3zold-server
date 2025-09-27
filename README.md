📦 m3zold-server
Backend-сервер на Go с чистой архитектурой для управления DIY-устройствами (ESP8266, NeoPixel и др.). Поддерживает регистрацию пользователей, CRUD-операции с устройствами, миграции, DI-контейнер и готов к JWT-аутентификации.

🚀 Быстрый старт
bash
# Клонируй репозиторий
git clone https://github.com/yourname/m3zold-server.git
cd m3zold-server

# Настрой .env
cp .env.example .env
# Укажи DB_URL, PORT и другие переменные

# Запусти миграции
go run migrations/migrate.go

# Запусти сервер
go run cmd/main.go
Или через Docker:

bash
docker-compose up --build
🧱 Архитектура
Проект построен по принципам Clean Architecture:

Код
internal/
├── domain/         # Сущности и интерфейсы
├── usecase/        # Бизнес-логика
├── infrastructure/ # Реализация репозиториев
├── delivery/       # HTTP-обработчики и роутер
├── container/      # DI-сборка
├── constants/      # Ошибки и константы
📁 Структура
bash
|-- cmd/main.go
|-- internal/
|   ├── domain/device/user
|   ├── usecase/device/user
|   ├── infrastructure/postgres
|   ├── delivery/http
|   ├── container
|   └── constants
|-- migrations/
|   ├── 001_create_users.sql
|   ├── 002_create_devices.sql
|   └── migrate.go
|-- .env
|-- Dockerfile
|-- docker-compose.yml
|-- README.md
🔐 Эндпоинты
Метод	URL	Описание
POST	/auth/register	Регистрация пользователя
POST	/devices	Добавить устройство
GET	/devices	Получить список
GET	/devices/:id	Получить по ID
PATCH	/devices/:id/status	Обновить статус
DELETE	/devices/:id	Удалить устройство
⚠️ Защищённые маршруты требуют заголовок X-User-ID (временно, позже — JWT).

🛠️ TODO
[ ] JWT-аутентификация и мидлвара

[ ] Swagger-документация

[ ] Тесты usecase'ов и репозиториев

[ ] WebSocket для реального времени

[ ] Поддержка ESP OTA / MQTT

[ ] Admin-панель и роли

📜 Лицензия
MIT — свободно используй, адаптируй и развивай.