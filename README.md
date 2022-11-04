# b2b2b



# Прототип системы автоматизированных взаимозачетов

В текущей реализации два микросервиса:
- repository общается с postgresql и принимает команды по gRPC
- http2repo общается с repository и принимает команды по http

и программа xml2repo, которая на вход принимает имена файлов с xml,
находит нужную информацию и отправляет в repository

Сервисы написаны в минимальном функционале, достаточном для того, чтобы проверить их взаимодействие.

**Для запуска:**

Если не создан том для postgres:
* sudo docker volume create --name=pg-data

ну или можно поправить docker-compose.yaml, убрав строчку "external: true" в параметрах тома

* sudo docker-compose build
* sudo docker-compose up -d

Получение токенов:
* curl -X POST -d 'grant_type=password&username=user01&password=12345' http://localhost:8888/token

Обновление auth токена
* curl -X POST -d 'grant_type=refresh_token&refresh_token={refresh_token при использовании curl в токене заменить все + на %2B}' http://localhost:8888/token


Потом можно постучаться на 8088 по http

* curl -X GET -H "Authorization: Bearer {auth_token}" http://localhost:8088/
* curl -X POST -H "Authorization: Bearer {auth_token}" -d '{"entity_inn":"1001", "entity_kpp":"1001", "short_name":"тест2"}' http://localhost:8088/entity
* curl -X GET -H "Authorization: Bearer {auth_token}" http://localhost:8088/obligation
* и т.д.

Или запустить xml2http, передав ей на вход имя xml-файла, который предусмотренно оставлен
* cd xml2repo
* go run cmd/app/main.go contract_1241101171822000013_74728698.xml






# прототип авторизации
поднял auth - сервис аутентификации
в http2repo токен проверяется и передаётся в репозиторий
сделал проверку токена в репозитории
сделал авторизацию на уровне репозитория для функции получения списка долгов:
показываются только те долги, в которых участвуют предприятия, с которыми связан пользователь


# Где я?
Оказывается попал на антипаттерн "наносервис".
Надо поработать над структурой кода и красотой
Надо сделать бота, работающего с системой.


# Ближайшие планы:
- 
- подумать над системой именования и структурой кода
- отработать ошибки синхронизации сервисов на старте
- написать внятную документацию
- подумать какой токен должен идти от xml2repo
- подумать про авторизацию телеграм-клиентов
- подумать про хранение пользовательских учётных записей
- секретные ключи получать в конфиге

# Неожиданности

1. Если в VSCode открыть проект, используя путь, в котором есть символические ссылки (сделанные командой ln -s), то debugger не останавливается на break point'ах.

2. можно сделать
type td struct {
}
func (t *td) testdb(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "я могу писать в базу\n")
}
r.Get("/", x.testdb)

хотя:
type HandlerFunc func(ResponseWriter, *Request)
func (mx *Mux) Get(pattern string, handlerFn http.HandlerFunc) 
