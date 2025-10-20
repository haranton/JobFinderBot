Telegram Job Bot

Telegram Job Bot — это бот для поиска и подписки на вакансии. Пользователи могут выполнять поиск по ключевым словам, подписываться на новые вакансии и получать уведомления при их появлении.
Проект написан на Go с использованием PostgreSQL и Telegram Bot API.

Features

Регистрация пользователей

Поиск вакансий по ключевым словам

Подписка на вакансии

Удаление всех подписок

Автоматическая отправка новых вакансий

Хранение данных в PostgreSQL

Миграции через golang-migrate

Project Structure
tgbot/
├── cmd/
│   └── main.go                # Точка входа
├── internal/
│   ├── bot/                   # Работа с Telegram API
│   ├── config/                # Конфигурация окружения
│   ├── db/                    # Подключение к базе данных
│   ├── dto/                   # Data Transfer Objects
│   ├── fetcher/               # Работа с внешним API вакансий
│   ├── handler/               # Обработка сообщений Telegram
│   ├── models/                # Модели базы данных
│   ├── repo/                  # Репозиторий (доступ к данным)
│   ├── sender/                # Отправка уведомлений пользователям
│   └── service/               # Бизнес-логика
└── migrations/                # SQL миграции


Bot Commands
Command	Description
/start	Регистрация пользователя
/find <запрос>	Поиск вакансий
/subscribe <запрос>	Подписка на вакансии
/subscribes	Просмотр всех подписок
/deletesubscribes	Удаление всех подписок
/help	Список доступных команд
Example

Поиск вакансий:

/find golang remote


Подписка на вакансии:

/subscribe python developer


Просмотр активных подписок:

/subscribes

Technologies

Go 1.22+

PostgreSQL

golang-migrate

Docker / Docker Compose

Telegram Bot API

slog (structured logging)
