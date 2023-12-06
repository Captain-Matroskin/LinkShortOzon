# Тестовое задание Ozon: проект для сокращения ссылок 

---
Задача: реализовать сервис, предоставляющий API по созданию сокращённых ссылок.

Работу выполнил - Жданов Никита.

## Примеры запросов

Отправка POST запроса на сокращение ссылки:
```
curl --header "Content-Type: application/json" \
--request POST \
--data '{"link": "www.testSite.ru"}' \
0.0.0.0:5001/api/v1/linkShort/
```


Отправка GET запроса на получение полной ссылки по полученной сокращенной:
```
curl --header "Content-Type: application/json" \
  --request GET \
  --data '{"link": "$shortLink$"}' \
  0.0.0.0:5001/api/v1/linkShort/
```
где *$shortLink$* - это готовая сокращенная ссылка (например "ozon.click.ru/_FeLIUZ33Y").

## Запуска проекта

Запуск проекта происходит через *Docker* контейнеры посредством *docker-compose.yml*
```
docker-compose -f docker-compose.yml up -d --build
```