# Распределённые алгоритмы

## Задание

Реализовать распределённое ключ-значение хранилище, использующее алгоритм RAFT с распределённой блокировкой. Использовать [статью](In%20Search%20of%20an%20Understandable%20Consensus%20Algorithm.pdf).


## Запуск

```bash
mkdir generated
```

```bash
protoc --go_out=generated --go_opt=paths=source_relative --go-grpc_out=generated --go-grpc_opt=paths=source_relative proto/raft.proto
```
