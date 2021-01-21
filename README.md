# dev-portal

A StackOverflow / Reddit / Disqus / Talkyard clone.
For a short description of the project, look at [dev-portal-frontend](https://github.com/denysvitali/dev-portal-frontend).

## Requirements

- PostgreSQL
- Go

## Running

```bash
go run https://github.com/denysvitali/dev-portal/cmd/
```

## Remarks

- Still WIP
- For now, the credentials are hardcoded, a proper main.go will come at a later stage
- `docker-compose up -d` + `go run cmd/main.go` and you should have something on `http://127.0.0.1:8081/api/v1/topics/`
- Use the [frontend](https://github.com/denysvitali/dev-portal-frontend)
- PRs / Issues / Feedback are always welcome
