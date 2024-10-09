# Тестовое задание

## Описание
Необходимо спроектировать и разработать сервис, который работает с реляционным хранилищем данных PostgreSQL.
Основной упор необходимо сделать на программную архитектуру: слои, маппинги структур данных между слоями приложения.

## Основные сущности

### [User](internal/service/entities/user.go)
* id
* firstname
* lastname
* fullname (firstname + lastname) ([as method](internal/service/entities/user.go:39))
* ~~age~~ birthdate
* ~~is_married~~ marital status
* password

###  [Product](internal/service/entities/product.go)
* id
* description
* tags
* quantity

## Что надо сделать

Реализовать следующий функционал:
1. [x] Регистрация пользователя (не младше 18 лет) ([проверка возраста при создании пользователя в юзкейсе userRegistration](./internal/service/usecases/user/registration/usecase.go:70));
2. [x] Пароль не меньше 8 символов ([проверка пароля при создании пользователя в юзкейсе userRegistration](./internal/service/usecases/user/registration/usecase.go:70));
3. [x] Пользователь может заказать продукт ([юзкейс addProductToOrder](internal/service/usecases/order/add_product_to_order/usecase.go));
4. [x] У пользователя может быть много заказов ([в сущности пользователя](internal/service/entities/user.go:29) и в [сущности заказа](internal/service/entities/order.go:27));
5. [x] Заказ может содержать множество продуктов (в [сущности заказа](internal/service/entities/order.go:28));
6. [x] Если, продуктов не осталось на складе – его нельзя заказать ([проверка на остатки на складах при добавлении в заказ](internal/service/usecases/order/add_product_to_order/usecase.go:94));
7. [x] Нужна историчность заказов (у заказа имеется [жизненный цикл](internal/service/entities/value_objects/order_status.go:16), который определяется переходами между статусами) и продуктов в заказе (например старая цена) (сущность orderProduct, которая меняется при изменении количества заказываемых продуктов в заказе и фиксирует цену товара при заказе).

## Тесты

* [x] Покрыть тестами несколько (2-3) функционально важных методов. (100% покрытие сформированных юзкейсов [1](internal/service/usecases/user/registration), [2](internal/service/usecases/order/add_product_to_order))

## Тезисно

Не все из перечисленного ниже обязательно реализовывать.
1. [ ] REST API
2. [x] Слоеная архитектура (без транспортного слоя, в котором дёргаются [юзкейсы](internal/service/usecases), которые используют [команды](internal/service/commands) и [запросы](internal/service/queries) (CQRS), которые взаимодействуют с [репозиториями](internal/service/repositories) и всё это имплементируется в [ioc-контейнере](internal/service/ioc/container.go))
3. [ ] Логирование (~~в контексте~~ через DI) - ~~middleware~~
4. [x] Трасировка, opentelemetry - ~~middleware~~ (как пример применения трассировки в [db модуле](internal/pkg/db/db.go:30) только без реализации)
5. [ ] Sentry - ловить паники в middleware
6. [x] На каждом слое своя структура данных
7. [x] Поток данных идет как в чистой или гексогональной архитектуре
