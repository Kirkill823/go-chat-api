# go-chat-api

Для запуска докер проекта небходимо выполнить docker-compose up --build

Адрес сервера будет http://localhost:8080

Методы API (Копия из тз)
1. POST /chats/ — создать чат
Body: title: str
Response: созданный чат
2. POST /chats/{id}/messages/ — отправить сообщение в чат
Body: text: str
Response: созданное сообщение
3. GET /chats/{id} — получить чат и последние N сообщений
Query: limit (по умолчанию 20, максимум 100)
Response:
● чат
● messages: [] (сообщения отсортированы по created_at)
4. DELETE /chats/{id} — удалить чат вместе со всеми сообщениями
Response: 204 No Content (или json-статус)
