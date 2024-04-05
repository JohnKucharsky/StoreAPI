### [Schema](https://drawsql.app/teams/johns-team-27/diagrams/storeapi)

# Getting started

## Build the app
```
docker compose up
```

## Make Targets

- **make dev**: Starts app in dev mode, need to install fresh
- **make dev-db**: Starts databases for dev
- **make prod**: docker compose up
- **make migrate**: Migrations, uses goose
- **make migrate-down**: Migration down
- **make api**: Generates html from openapi, uses redocly 
