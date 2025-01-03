## Сервис по сбору метрик (Сервер)

Проект состоит из двух частей: агент и сервер, в данном репозитории представлен Сервер. Сервер - принимает, обрабатывает и сохраняет метрики.

## Мотивация проекта

Расширить знания о разработке API-приложений, реализовать микросервисную архитектуру, написать приложение, поддерживающее разные способы хранения данных, реализовать протоколы gRPC и HTTP. Попрактиковаться в написании unit-тестов.

## Маршруты API

GET / — Возвращает HTML страницу со всеми метриками последней итерации.
POST /api/updates — Принимает batched (групповой) запрос метрик в формате JSON.
POST /api/update — Принимает метрику (DTO) в теле запроса в формате JSON.
POST /api/{type}/{name}/{value} — Принимает значение метрики через параметры URL.
GET /api/value/{type}/{name} — Возвращает значение метрики по её типу и названию из URL.

## Описание функционала приложения

Сервер умеет работать с протоколами gRPC и HTTP. Хранение метрик поддерживается разными способами: в памяти (slice), в файле, в базе данных (postgresql). Каждые m секунд Сервер принимает метрики, переданные Агентом и сохраняет их одним из способов. При получении метрик выполняется хэширование данных и если хэш совпадает с заголовком запроса Hash, то Сервер принимает запрос.

## Работа с приложением

Для запуска приложения необходимо создать файл .env с переменными, записанными в .env.example и поднять docker-compose.