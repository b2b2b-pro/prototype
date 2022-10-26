# b2b2b

Прототип системы автоматизированных взаимозачетов

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

Потом можно постучаться на 8088 по http
* curl -X POST -d '{"entity_inn":"1000", "entity_kpp":"1000", "short_name":"тест"}' http://localhost:8088/entity
* curl http://localhost:8088/entity

Или запустить xml2http, передав ей на вход имя xml-файла, который предусмотренно оставлен
* cd xml2repo
* go run cmd/app/main.go contract_1241101171822000013_74728698.xml

Можно убедиться, что данные попали в базу:
* curl http://localhost:8088/entity



# Где я?
Оказывается попал на антипаттерн "наносервис".
Прикручу авторизацию, потом буду переделывать.



# Ближайшие планы:
- прикрутить авторизацию пока не поздно
- подумать над системой именования и структурой кода
- отработать ошибки синхронизации сервисов на старте
- написать внятную документацию

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
