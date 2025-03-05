## Реализация
* Для того, чтобы облегчить возможную миграцию при будущих обратно-несовместимых изменениях, API сервиса доступно по двум префиксам — `/api` и `/api/v1`. Предполагаем, что в случае если API изменится, то его новая версия будет доступна по `/api/v2`, а по адресу `/api/v1` некоторое время будет поддерживаться deprecated версия, совместимая с сервисами, которые не успели обновиться. По `/api` всегда поддерживаем последнюю версию.
* `.env` для удобства проверки закоммичен в репозиторий. В реальной жизни так разумеется делать не надо.

## Требования
* Golang 1.24

## Конфигурация
Образец переменных окружения/ образец содержимого .env файла находится в .env.example.
Поддерживаются следующие настройки:
* `EMLIB_DB_HOST` — хост Postgres (по умолчанию `localhost`)
* `EMLIB_DB_PORT` — порт подключения к Postgres (по умолчанию `5432`)
* `EMLIB_DB_USER` — пользователь базы данных (по умолчанию `postgres`)
* `EMLIB_DB_PASSWORD` — пароль подключения к базе данных (по умолчанию пустой)
* `EMLIB_DB_NAME` — название базы на сервере Postgres.
* `EMLIB_SERVER_PORT` — порт, на котором будет доступно API микросервиса (по умолчанию `8080`)
* `EMLIB_SERVER_READ_TIMEOUT`
* `EMLIB_SERVER_WRITE_TIMEOUT`
* `EMLIB_SERVER_MODE` — режим работы Gin сервера (по умолчанию `debug`)
* `EMLIB_LOG_LEVEL` — уровень логирования в логике приложения (по умолчанию `debug`)
