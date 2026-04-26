# Тестовое задание №5

Это тестовое задание HiTalent на позицию Junior Go Developer.

- [Описание задачи](docs/description.md)

## API

[![Run in Yaak](https://yaak.app/static/button.svg)](https://yaak.app/button/run?name=HiTalent+Go+Junior+Task&url=https%3A%2F%2Fgitverse.ru%2Fapi%2Frepos%2Fdigit4lsh4d0w%2Ftasks-go-junior-05-hitalent%2Fraw%2Fbranch%2Fmain%2Fyaak%252Fyaak.tasks-go-05-junior-hitalent.json)

| Действие                           | Метод  | Путь                    | Параметры |
| ---------------------------------- | ------ | ----------------------- | --------- |
| Создание чата                      | POST   | /chat                   |           |
| Добавление сообщения в чат         | POST   | /chat/{chat_id}/message |           |
| Просмотр чата с сообщениями        | GET    | /chat/{chat_id}         | limit     |
| Удаление чата вместе с сообщениями | DELETE | /chat/{chat_id}         |           |

### Создание чата

Пример запроса:

```bash
curl --json '{"title": "'$(openssl rand -hex 8)'"}' http://localhost:3000/chat
```

Пример ответа:

```json
{
  "id": 1290,
  "title": "eiShi0daideezae2",
  "created_at": "2026-04-24T22:24:55.585974872+10:00"
}
```

### Добавление сообщения в чат

Пример запроса:

```bash
curl --json '{"text": "'$(openssl rand -hex 8)'"}' http://localhost:3000/chat/1290/message
```

Пример ответа:

```json
{
  "id": 113,
  "chat_id": 1290,
  "text": "ec234599893410ae",
  "created_at": "2026-04-24T22:26:25.071143979+10:00"
}
```

### Просмотр чата с сообщениями

Пример запроса:

```bash
curl http://localhost:3000/chat/1290
```

Пример ответа:

```json
{
  "id": 1290,
  "title": "eiShi0daideezae2",
  "created_at": "2026-04-24T22:24:55.585974872+10:00",
  "messages": [
    {
      "id": 113,
      "chat_id": 1290,
      "text": "ec234599893410ae",
      "created_at": "2026-04-24T22:26:25.071143979+10:00"
    }
  ]
}
```

### Удаление чата вместе с сообщениями

Пример запроса:

```bash
curl -X DELETE http://localhost:3000/chat/1290
```

Пример ответа:

```json
{
  "success": "chat deleted successfully"
}
```

## Docker

Сборка:

```bash
docker buildx build --platform linux/amd64 -t tasks-go-junior-05-hitalent:latest .
```

Запуск (CLI):

```bash
docker run -it --rm \
    -p 3000:3000 \
    -v ./config.yaml:/app/config.yaml \
    -v ./sqlite3.db:/app/sqlite3.db \
    -v ./sqlite3.db-shm:/app/sqlite3.db-shm \
    -v ./sqlite3.db-wal:/app/sqlite3.db-wal \
    tasks-go-junior-05-hitalent:latest
```
