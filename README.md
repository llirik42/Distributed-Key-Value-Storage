# Распределённое key-value хранилище

Программа представляет собой узел распределённого key-value хранилища, использующего RAFT. Узлы общаются между собой через gRPC, описание протокола лежит [тут](proto/raft.proto).

## Запуск

### Генерация protobuf

Для начала нужно сгенерировать код для protobuf, который используется в gRPC:

```bash
mkdir generated
```

```bash
protoc --go_out=generated --go_opt=paths=source_relative --go-grpc_out=generated --go-grpc_opt=paths=source_relative proto/raft.proto
```

### Конфигурация

Каждый узел конфигурируется в своём .env-файле (**1.env**, **2.env**, **3.env** и т.д. — файлы для конфигурации первого, второго, третьего и других узлов). Конфигурационный файл должен содержать следующие поля:
* `EXECUTED_COMMANDS_KEY`
* `SELF_ADDRESS`
* `SELF_ID`
* `OTHER_NODES`
* `BROADCAST_TIME_MS`
* `MIN_ELECTION_TIMEOUT_MS`
* `MAX_ELECTION_TIMEOUT_MS`



## Взаимодействие с клиентами



## Запуск

```bash
mkdir generated
```

```bash
protoc --go_out=generated --go_opt=paths=source_relative --go-grpc_out=generated --go-grpc_opt=paths=source_relative proto/raft.proto
```
