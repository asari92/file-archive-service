# file-archive-service

## Описание проекта

В данном проекте реализован REST API для работы с архивами и отправки файлов по электронной почте. API поддерживает загрузку, анализ, создание архивов и их отправку на указанные электронные адреса. Разработка ведется с соблюдением стандартов чистого кода и архитектуры, обеспечивая высокую производительность и минимальное потребление ресурсов.

## Технологии

- **Язык программирования:** Go

### Библиотеки:
- `github.com/mholt/archiver/v3`: Используется для работы с архивными файлами, включая их распаковку.
- `gopkg.in/mail.v2`: Используется для отправки электронной почты через SMTP с поддержкой множества функций, таких как вложения, HTML-тела и пользовательские заголовки.
- `net/http`: Для обработки HTTP-запросов и ответов.
- `mime/multipart`: Для обработки многокомпонентных форм данных, необходимых для загрузки файлов.
- `log/slog`: Для логирования операций и ошибок, что улучшает отладку и мониторинг приложения.


## Установка и запуск

Клонирование репозитория:
```
git clone https://github.com/asari92/file-archive-service.git && cd file-archive-service
```
Установка зависимостей:
```
make go-get
```
Запуск сервера:
```
make go-run
```

### Настройка среды

Перед запуском приложения необходимо настроить переменные окружения. Для этого:

1. Создайте файл `.env` в корневой директории проекта, используя пример из файла `.env.example`.
   
2. Заполните файл `.env` соответствующими значениями переменных:

   ```plaintext
   HOST=localhost
   PORT=10000
   MAX_UPLOAD_SIZE_INFO=10485760    # 10MB
   MAX_UPLOAD_SIZE_CREATE=33554432  # 32MB
   MAX_UPLOAD_SIZE_MAIL=10485760    # 10MB
   MAX_SEND_FILE_SIZE=26214400      # 25MB
   DIALER_TIMEOUT=60                # Таймаут для SMTP соединения
   MAIL_FROM=file-archive-service   # Отправитель почты
   SMTP_HOST=smtp.example.com       # SMTP сервер
   SMTP_PORT=587                    # SMTP порт
   SMTP_USER=user@example.com       # SMTP пользователь
   SMTP_PASSWORD=securepassword     # SMTP пароль

Эти настройки включают конфигурации для SMTP, что позволяет отправлять электронные письма через указанный почтовый сервер.

## API Роуты

### 1. Получение информации об архиве

- **Метод:** POST
- **URL:** `/api/archive/information`
- **Описание:** Принимает файл архива через `multipart/form-data` и возвращает структурированную информацию о содержимом архива.
- **Тело запроса:** Файл должен быть отправлен в формате `multipart/form-data` с ключом `file`.
- **Пример запроса:**

```
POST /api/archive/information HTTP/1.1 Content-Type: multipart/form-data; boundary=-{some-random-boundary}

-{some-random-boundary} Content-Disposition: form-data; name="file"; filename="my_archive.zip" Content-Type: application/zip

{Binary data of ZIP file} -{some-random-boundary}--
```

- **Пример успешного ответа (HTTP 200):**
```
{
    "filename": "my_archive.zip",
    "archive_size": 4102029,
    "total_size": 6836715,
    "total_files": 2,
    "files": [
        {
            "file_path": "photo.jpg",
            "size": 2516582,
            "mimetype": "image/jpeg"
        },
        {
            "file_path": "directory/document.docx",
            "size": 4320133,
            "mimetype": "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
        }
    ]
}
```
### 2. Формирование архива из файлов

- **Метод:** POST
- **URL:** `/api/archive/files`
- **Описание:** Принимает список файлов, объединяет их в .zip архив и возвращает архив клиенту.
- **Тело запроса:** Файлы должны быть отправлены в формате `multipart/form-data` с ключом `files[]`.
- **Пример запроса:**

```
POST /api/archive/files HTTP/1.1
Content-Type: multipart/form-data; boundary=-{some-random-boundary}

-{some-random-boundary}
Content-Disposition: form-data; name="files[]"; filename="document.docx"
Content-Type: application/vnd.openxmlformats-officedocument.wordprocessingml.document

{Binary data of file}
-{some-random-boundary}
Content-Disposition: form-data; name="files[]"; filename="avatar.png"
Content-Type: image/png

{Binary data of file}
-{some-random-boundary}--
```
- **Пример успешного ответа (HTTP 200):**
```
    HTTP/1.1 200 OK
    Content-Type: application/zip

    {Binary data of ZIP file}
```

### 3. Отправка файла на несколько почт

- **Метод:** POST
- **URL:** `/api/mail/file`
- **Описание:** Принимает файл и список почт, отправляет файл на все указанные адреса.
- **Тело запроса:** Файл и почты должны быть отправлены в формате `multipart/form-data` с ключами `file` и `emails`.
- **Пример запроса:**

```
POST /api/mail/file HTTP/1.1
Content-Type: multipart/form-data; boundary=-{some-random-boundary}

-{some-random-boundary}
Content-Disposition: form-data; name="file"; filename="document.docx"
Content-Type: application/vnd.openxmlformats-officedocument.wordprocessingml.document

{Binary data of file}
-{some-random-boundary}
Content-Disposition: form-data; name="emails"

elonmusk@x.com,jeffbezos@amazon.com,zuckerberg@meta.com
-{some-random-boundary}--
```
- **Пример успешного ответа (HTTP 200):**
```
HTTP/1.1 200 OK
```

## Тестирование


Для запуска тестов используйте следующую команду:

```bash
make test
```
Эта команда выполнит все тесты, находящиеся в тестовых файлах проекта.


## Лицензия

Данный проект распространяется под лицензией MIT. Это позволяет использовать, копировать, модифицировать, сливать, публиковать, распространять, выдавать подлицензии и/или продавать копии программного обеспечения, а также разрешать лицам, которым программное обеспечение предоставляется, делать это, при условии указания следующего уведомления об авторских правах и данного разрешения во всех копиях или существенных частях программного обеспечения.

Подробнее с текстом лицензии можно ознакомиться [здесь](LICENSE).