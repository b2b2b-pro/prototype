## Прототип сервиса, который обеспечивает взаимодействие с хранилищем на базе postgresql

Для запуска:

sudo docker volume create --name=pg-data
sudo docker-compose build
sudo docker-compose up -d

TODO:
- авторизация

P.S. миграции захардкодил пока

P. P.S. не нравится, что большие куски кода сильно дублируются, надо поискать другое решение