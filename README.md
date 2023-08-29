# avito-trainee

# start docker
docker compose up --build

# curl commands

## Добавление пользователя (на вход user_id)

curl -X POST http://localhost:8080/users -H "Content-type: application/json" -d '{ "user_id": "1000"}'

## Добавление сегмента (на вход slug сегмента)

curl -X POST http://localhost:8080/segments -H "Content-type: application/json" -d '{ "name": "AVITO_VOICE_MESSAGES"}'

## Удаление сегмента (на вход slug сегмента)

curl -X DELETE http://localhost:8080/segments -H "Content-type: application/json" -d '{ "name": "AVITO_VOICE_MESSAGES"}'

## Добавление + удаление пользователя в сегмент 

curl -X PUT http://localhost:8080/users -H "Content-type: application/json" -d '{"user_id":"1000", "segment_to_join_names":[{"name":"AVITO_PERFORMANCE_VAS"},{"name":"AVITO_DISCOUNT_30"}],"segment_to_delete_names":["AVITO_VOICE_MESSAGES"]}'

## Получить активные сегменты пользователя

curl -X GET http://localhost:8080/users -H "Content-type: application/json" -d '{"user_id":"1000"}'

## Доп. задание 2

curl -X PUT http://localhost:8080/users -H "Content-type: application/json" -d '{"user_id":"1001", "segment_to_join_names":[{"name":"AVITO_DISCOUNT_30","deleted":"2023-08-30T00:00:00Z"},{"name":"AVITO_VOICE_MESSAGES"}]}'
