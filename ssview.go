package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"io/ioutil"
	"net/url"
	"os"
)

const (
	tablemode = iota
	sqlmode
)

func main() {
	conf, err := ReadConfig()
	var prompt = [2]string{"table? > ", "sql? > "}

	if err != nil {
		fmt.Println("Reading config file has failed.")
		fmt.Printf("err: %v\n", err)
		return
	}

	// connect to db
	u := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(conf.User, conf.Password),
		Host:     conf.Host + ":" + conf.Port,
		RawQuery: "database=" + conf.Database,
	}
	fmt.Println("Connection URL:" + u.String())
  
	db, err := sql.Open("sqlserver", u.String())
	if err != nil {
		fmt.Println("Cannot connect to server.")
		fmt.Println(u.String())
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
					query = fmt.Sprintf("select top %v * from %s", conf.Limit, tableName)
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
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
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