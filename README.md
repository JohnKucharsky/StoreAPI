### [Schema](https://drawsql.app/teams/johns-team-27/diagrams/storeapi)

# Getting started

## Build the app
```
docker compose up
```

## For dev

### Install goose
```
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### Install fresh
```
go install github.com/zzwx/fresh@latest
```

### Install redocly
```
yarn global add @redocly/cli
```

## Make Targets

- **make dev**: Starts app in dev mode, fresh
- **make dev-db**: Starts databases for dev
- **make prod**: docker compose up
- **make migrate**: Migrations, goose
- **make migrate-down**: Migration down
- **make api**: Generates html from openapi, redocly 
