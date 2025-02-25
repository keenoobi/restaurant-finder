# Поиск ресторанов в Москве

Этот проект представляет собой REST API для поиска ресторанов в Москве с использованием Elasticsearch. Сервис предоставляет возможность:
- Получения списка ресторанов с пагинацией.
- Поиска ближайших ресторанов по координатам.
- Аутентификации через JWT.

## Структура проекта

    .
    ├── bin
    ├── cmd
    ├── data.csv
    ├── docker-compose.yml
    ├── go.mod
    ├── go.sum
    ├── internal
    ├── templates
    ├── tests
    └── web-app

    7 directories


## Установка и запуск

### Требования
- Go 1.16 или выше
- Docker (для запуска Elasticsearch)

### Использование Makefile

Проект включает `Makefile` для упрощения сборки и запуска. Доступные команды:

| Команда           | Описание                                                                 |
|-------------------|-------------------------------------------------------------------------|
| `make build`      | Сборка проекта.                                                         |
| `make elastic-up` | Запуск Elasticsearch через Docker Compose.                              |
| `make elastic-down` | Остановка Elasticsearch.                                              |
| `make load-data`  | Загрузка данных в Elasticsearch.                                        |
| `make run`        | Запуск веб-приложения.                                                  |
| `make start`      | Запуск всех сервисов (Elasticsearch + веб-приложение).                  |
| `make stop`       | Остановка всех сервисов.                                                |
| `make clean`      | Очистка проекта (удаление скомпилированных файлов).                     |
| `make test`       | Запуск тестов.                                                          |
| `make help`       | Показать список доступных команд.                                       |

### Пошаговая инструкция

#### 1. Клонируйте репозиторий:
```bash
git clone https://github.com/ваш-username/restaurant-finder.git
cd restaurant-finder
```
#### 2. Запустите Elasticsearch:
```make
make elastic-up
```
#### 3. Загрузите данные в Elasticsearch:
```make
make load-data
```
#### 4. Запустите веб-приложение:
```make
make load-data
```
*Или запустите все сервисы одной командой:*
```make
make start
```
#### 5. Запустите веб-приложение:
```make
make stop
```
#### 6. Очистите проект (если нужно)::
```make
make clean
```
## Использование API
### Получение списка ресторанов
- GET `/api/places?page=<номер_страницы>`

Пример ответа:
```json
{
  "name": "Places",
  "total": 13649,
  "places": [
    {
      "id": 1,
      "name": "Sushi Wok",
      "address": "gorod Moskva, prospekt Andropova, dom 30",
      "phone": "(499) 754-44-44",
      "location": {
        "lat": 55.879001531303366,
        "lon": 37.71456500043604
      }
    }
  ],
  "prev_page": 1,
  "next_page": 3,
  "last_page": 1364
}
```
### Получение списка ресторанов
- GET `/api/recommend?lat=<широта>&lon=<долгота>`
Пример ответа:
```json
{
  "name": "Recommendation",
  "places": [
    {
      "id": 30,
      "name": "Ryba i mjaso na ugljah",
      "address": "gorod Moskva, prospekt Andropova, dom 35A",
      "phone": "(499) 612-82-69",
      "location": {
        "lat": 55.67396575768212,
        "lon": 37.66626689310591
      }
    }
  ]
}
```
### Получение JWT токена
- GET `/api/get_token`
Пример ответа:
```json
{
  "name": "Recommendation",
  "places": [
    {
      "id": 30,
      "name": "Ryba i mjaso na ugljah",
      "address": "gorod Moskva, prospekt Andropova, dom 35A",
      "phone": "(499) 612-82-69",
      "location": {
        "lat": 55.67396575768212,
        "lon": 37.66626689310591
      }
    }
  ]
}
```
### Аутентификация
Для доступа к эндпоинту `/api/recommend` необходимо передать JWT-токен в заголовке `Authorization`:
```
Authorization: Bearer <ваш_токен>
```
### Тестирование
Для тестирования можно использовать `curl` или Postman.

Пример запроса:
```bash
curl -X GET "http://localhost:8888/api/places?page=1"
```