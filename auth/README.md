# Сервис аутентификации

Аутентификация осуществляется по http

curl -X POST -d 'grant_type=password&username=user01&password=12345' http://localhost:8888/token

curl -X POST -d 'grant_type=refresh_token&refresh_token={refresh_token при использовании curl в токене заменить все + на %2B}' http://localhost:8888/token

