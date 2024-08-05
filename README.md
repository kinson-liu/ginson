# Ginson

As the name, it's a simple REST framework based on [Gin](https://github.com/gin-gonic/gin)

## Features
1. Simple structure, move routers into each api file.
2. Restful apis, support authenticator, validator, paginator, sorter, filter.
3. Reload configs when the config file changes.
## Dir Structure
```shell
ginson
├── api               >> layer API
│   ├── request       >> request structs
│   └── response      >> response struct
├── config            >> service configs
├── core              >> core functions
│   └── const         >> consts
├── docs              >> swagger docs
├── logs              >> service logs
├── middleware        >> middlewares
│   └── auth
├── model             >> Layer Model
└── service           >> Layer Service
```

## Develop Plan
- [x] Config
- [x] Log
- [x] JWT Token
- [x] Redis
- [x] Mongo
- [x] Swagger
- [ ] Casbin
- [ ] Code Generate
