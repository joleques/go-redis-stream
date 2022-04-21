# Redis Stream com Golang

- POC para testar envio e consumo no Redis-Stream

## Vaiaveis de Ambiente: 
- Quando rodar local

``` 
export REDIS_HOST=localhost
export STREAM=events
export GROUP=groupLocal
```

# Run Producer

```
go run producer/main.go
```

# Run Consumer

```
go run consumer/main.go
```

# Docker

- Sobe 4 consumers divido em 2 grupos, para testar é necessario rodar somente o producer

``` 
docker-compose up -d
```