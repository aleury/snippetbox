# Snippetbox

## Getting Started

Copy `envrc.example` to `.envrc` and configure the environment variable `DATABASE_URL`:

```sh
export DATABASE_URL=postgres://username:password@localhost:5432/snippetbox
```

Run the following to start the application

```bash
go run ./cmd/web -dsn=$(DATABASE_URL)
```
