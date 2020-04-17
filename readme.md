## 编译

```
$ GOOS=linux GOARCH=amd64 go build -o slowq main.go
```

## Usage

```
$ slowq --help
NAME:
   slowq - A golang convenient converter supports Database to Struct

USAGE:
   slowq [GLOBAL OPTIONS] [DATABASE]

VERSION:
   0.0.0

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --host value, -h value      Host address of the database. (default: "127.0.0.1")
   --port value, -P value      Port number to use for connection. Honors $MYSQL_TCP_PORT. (default: 3306)
   --username value, -u value  User name to connect to the database. (default: "root")
   --password value, -p value  Password to connect to the database.
   --time value, -t value      slow timeout.
   --help                      Show this message and exit (default: false)
   --version                   Output mycli's version. (default: false)
```

## 执行命令

```
$ ./slowq -h 127.0.0.1 -P 3306 -u root -p root -t 10
```