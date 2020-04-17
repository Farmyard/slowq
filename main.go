package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/urfave/cli/v2"
)

type Option struct {
	Host     string
	Port     int
	Username string
	Password string
	Time     int64
}

func open(opt *Option) (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", opt.Username, opt.Password, opt.Host, opt.Port, "default", "utf8mb4"))
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	app := cli.NewApp()
	app.Name = "slowq"
	app.Usage = "Kill mysql slow query process"
	app.Version = "0.0.0"
	app.UsageText = "slowq [GLOBAL OPTIONS]"
	app.UseShortOptionHandling = true
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "host",
			Aliases: []string{"h"},
			Value:   "127.0.0.1",
			Usage:   "Host address of the database.",
		},
		&cli.IntFlag{
			Name:    "port",
			Aliases: []string{"P"},
			Value:   3306,
			Usage:   "Port number to use for connection. Honors $MYSQL_TCP_PORT.",
		},
		&cli.StringFlag{
			Name:    "username",
			Aliases: []string{"u"},
			Value:   "root",
			Usage:   "User name to connect to the database.",
		},
		&cli.StringFlag{
			Name:    "password",
			Aliases: []string{"p"},
			Usage:   "Password to connect to the database.",
		},
		&cli.Int64Flag{
			Name:    "time",
			Aliases: []string{"t"},
			Value:   15,
			Usage:   "slow timeout.",
		},
	}
	app.Action = func(c *cli.Context) error {
		opt := &Option{
			Host:     c.String("host"),
			Port:     c.Int("port"),
			Username: c.String("username"),
			Password: c.String("password"),
			Time:     c.Int64("time"),
		}
		if err := handle(opt); err != nil {
			fmt.Printf(err.Error())
		}
		return nil
	}
	cli.HelpFlag = &cli.BoolFlag{
		Name:  "help",
		Usage: "Show this message and exit",
	}
	cli.VersionFlag = &cli.BoolFlag{
		Name:  "version",
		Usage: "Output mycli's version.",
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func handle(opt *Option) error {
	f := bufio.NewReader(os.Stdin)
	if opt.Password == "" {
		fmt.Printf("Password:")
		input, _ := f.ReadString('\n')
		input = strings.TrimRight(input, "\n")
		opt.Password = input
	}
	if opt.Time < 10 {
		return errors.New("执行失败，参数 [-time] or [-t] 必须大于 10s")
	}
	// Open mysql
	db, err := open(opt)
	if err != nil {
		return err
	}
	fmt.Println("程序运行中...")
	for {
		show(db, opt)
		time.Sleep(10 * time.Millisecond)
	}
	return nil
}

type Processlist struct {
	Id      sql.NullInt64  `json:"id"`
	User    sql.NullString `json:"user"`
	Host    sql.NullString `json:"host"`
	Db      sql.NullString `json:"db"`
	Command sql.NullString `json:"command"`
	Time    sql.NullInt64  `json:"time"`
	State   sql.NullString `json:"state"`
	Info    sql.NullString `json:"info"`
}

func show(db *sql.DB, opt *Option) {
	rows, err := db.Query("SHOW PROCESSLIST")
	defer rows.Close()
	if err != nil {
		fmt.Println("show processlist error:", err)
	}
	for rows.Next() {
		var row Processlist
		if err := rows.Scan(&row.Id, &row.User, &row.Host, &row.Db, &row.Command, &row.Time, &row.State, &row.Info); err != nil {
			fmt.Println("rows scan error:", err)
		}
		if row.Command.String == "Query" && row.Time.Int64 >= opt.Time {
			id := strconv.FormatInt(row.Id.Int64, 10)
			if _, err := db.Exec("KILL " + id); err != nil {
				fmt.Println("kill process "+id+":", err)
			}
			fmt.Println("killed:'" + row.Info.String + "',id:" + id + ",time:" + strconv.FormatInt(row.Time.Int64, 10))
			msg, _ := json.Marshal(row)
			fmt.Println(string(msg))
		}
	}
}
