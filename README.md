### [Schema](https://drawsql.app/teams/johns-team-27/diagrams/storeapi)

# Getting started

## Build the app
```
docker compose up
```
## Open docs
### [Docs](http://127.0.0.1:8080/api)

## For dev

### Install goose
```
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### Install air
```
go install github.com/cosmtrek/air@latest
```

### Install redocly
```
yarn global add @redocly/cli
```

## Make Targets

- **make dev**: Starts app in dev mode, fresh
- **make dev-db**: Starts databases for dev
- **make prod**: docker compose up --build
- **make migrate**: Migrations, goose
- **make build-openapi**: Generates html from openapi, redocly 
