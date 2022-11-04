## Прототип сервиса, который обеспечивает взаимодействие с хранилищем на базе postgresql

Ответвился, делаю прототип авторизации.
Пока прибил гвоздями секрет и NewSHA256RC4TokenSecurityProvider

Автоиизация сделана для получения списка долгов,
работает так: токен проверяется на валидность, 
из токена получается login и его id, для него ищутся предприятия, к которым этот login привязан,
ищутся долги этих предприятий и возвращаются
Вот заклинание:
select 	obligation.id, obligation.debtor_id, obligation.creditor_id, obligation.cost, obligation.payment_date, obligation.origin_id 
	from obligation
        join (select account_entity.account_id, account_entity.entity_id from account_entity 
			where account_entity.account_id = $1) as A 
		    ON obligation.debtor_id = A.entity_id OR obligation.creditor_id = A.entity_id;

$1 = id пользователя

Если login == 'admin', то отдаются все долги

# Для запуска:

sudo docker volume create --name=pg-data
sudo docker-compose build
sudo docker-compose up -d

TODO:
- авторизация

P.S. миграции захардкодил пока

P. P.S. не нравится, что большие куски кода сильно дублируются, надо поискать другое решение