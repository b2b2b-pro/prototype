## http2repo - сервис - шлюз между http и gRPC

Находится в стадии раннего прототипа, функционал реализован не полностью, только в объёме
для проверки технологий.

Для получения токенов нужно постучаться в сервис аутентификации

curl -X GET -H "Authorization: Bearer {auth_token}" http://localhost:8088/

curl -X POST -H "Authorization: Bearer {auth_token}" -d '{"entity_inn":"1001", "entity_kpp":"1001", "short_name":"тест2"}' http://localhost:8088/entity

curl -X GET -H "Authorization: Bearer {auth_token}" http://localhost:8088/obligation


TODO:
- авторизация
- потом остальное

P.S. не нравится, что большие куски кода сильно дублируются, надо поискать другое решение

# мусом

curl http://localhost:8088/entity
curl -X POST -d '{"entity_inn":"1000", "entity_kpp":"1000", "short_name":"тест"}' http://localhost:8088/entity
