package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"net/url"
	"os"
)

const (
	tablemode = iota
	sqlmode
)

func connectToDB(conf *Config) (*sql.DB, error) {
	switch {
	case conf.Rdb == "sqlserver":
		// connect to sql server
		u := &url.URL{
			Scheme:   "sqlserver",
			User:     url.UserPassword(conf.User, conf.Password),
			Host:     conf.Host + ":" + conf.Port,
<<<<<<< HEAD
			Path:     conf.Instance,
=======
>>>>>>> 6a598ee5b1f0ab76b8d116665fd621f1524dcb77
			RawQuery: "database=" + conf.Database,
		}
		fmt.Println("Connection URL:" + u.String())
		db, err := sql.Open("sqlserver", u.String())
		return db, err
	case conf.Rdb == "mysql" || conf.Rdb == "mariadb":
		// connect to mysql or mariadb
		connectionstring := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", conf.User, conf.Password, conf.Host, conf.Port, conf.Database)
		fmt.Println("Connection String:" + connectionstring)
		db, err := sql.Open("mysql", connectionstring)
		return db, err
	default:
		return nil, errors.New("undefined rdb.")
	}
}

func main() {
	var prompt = [2]string{"table? > ", "sql? > "}
	conf, err := ReadConfig()
	if err != nil {
		fmt.Println("Reading config file has failed.")
		fmt.Printf("err: %v\n", err)
		return
	}

	db, err := connectToDB(conf)
	if err != nil {
		fmt.Println("Cannot connect to server.")
		fmt.Printf("err: %v\n", err)
		return
	}
	defer db.Close()

	var argStatement = flag.String("s", "", "sql statement")
	flag.Parse()
	if *argStatement != "" {
		// single command
		fmt.Printf("query: %v\n", *argStatement)
		rows, err := db.Query(*argStatement)
		if err != nil {
			fmt.Println("Error has Occured.")
			fmt.Printf("err: %v\n", err)
		} else {
			defer rows.Close()
			displayRows(rows, conf.Limit)
		}
		return
	} else {
		// interactive mode
		scanner := bufio.NewScanner(os.Stdin)
		mode := tablemode

		for {
			fmt.Print(prompt[mode])
			scanner.Scan()

			switch scanner.Text() {
			case "exit":
				return
			case "sql":
				mode = sqlmode
			case "table":
				mode = tablemode
			default:

				var query string
				switch mode {
				case tablemode:
					var tableName string = scanner.Text()
					switch {
					case conf.Rdb == "sqlserver":
						query = fmt.Sprintf("select top %v * from %s", conf.Limit, tableName)
					case conf.Rdb == "mysql" || conf.Rdb == "mariadb":
						query = fmt.Sprintf("select * from %s limit %v", tableName, conf.Limit)
					}
				case sqlmode:
					query = scanner.Text()
				}

				fmt.Printf("query: %v\n", query)
				rows, err := db.Query(query)
				if err != nil {
					fmt.Println("Error has Occured.")
					fmt.Printf("err: %v\n", err)
				} else {
					defer rows.Close()
					displayRows(rows, conf.Limit)
				}
			}
		}
	}
}

type Config struct {
	Rdb      string `json:"rdb"`
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Instance string `json:"instance"`
	Port     string `json:"port"`
	Database string `json:"database"`
	Limit    int    `json:"limit"`
}

// reading config file
func ReadConfig() (*Config, error) {
	const file = "./db.json"

	result := new(Config)

	conf, err := ioutil.ReadFile(file)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(conf), result)
	if err != nil {
		return result, err
	}

	return result, nil
}

//display query result on console
func displayRows(rows *sql.Rows, limit int) {
	cols, _ := rows.Columns()

	container := make([]sql.NullString, len(cols))
	ptrs := make([]interface{}, len(cols))

	for i := range ptrs {
		ptrs[i] = &container[i]
	}
	count := 0

	for rows.Next() {
		count++
		if count > limit {
			fmt.Println("Abort listing.")
			break
		}

		err := rows.Scan(ptrs...)
		if err != nil {
			fmt.Println("Reading row has failed.")
			fmt.Printf("err: %v\n", err)
		}
		for j := range container {
			if j != 0 {
				fmt.Printf(",")
			}
			if container[j].Valid {
				fmt.Printf(container[j].String)
			} else {
				fmt.Printf("NULL")
			}
		}
		fmt.Println()
	}
}
