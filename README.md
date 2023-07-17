# gameserver-go

This repository is **golang**-implemented [KLabServerCamp/gameserver](https://github.com/KLabServerCamp/gameserver).

|          | Language | repo |
|:--       |:--       |:--   |
| Original | Python   | <https://github.com/KLabServerCamp/gameserver> |
| This     | Golang   | this repo |

## Development

- Board: <https://github.com/users/pollenjp/projects/5>

### local run

```sh
make up
```

access to open port (See `docker-compose.yml` for more details.)

### Debug DB

```sh
make db-exec
```

```sql
USE webapp;
SHOW tables;
```

### local build

```sh
GOOS=linux GOARCH=amd64 go build
```

### golangci-lint

Install <https://github.com/golangci/golangci-lint> ([repo](https://github.com/golangci/golangci-lint))
