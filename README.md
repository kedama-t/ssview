# ssview
ssview is a small tool to glimpse data on Microsoft SQL Server.
(Because SSMS is too heavy on my pc...)

## Usage
Build:
```bash
go build ssview.go
```
This program is written with go 1.18.3.
And I'm verifying that can connect to SQL Server 2019(Windows) and Azure SQL Database.

In Windows:
```powershell
# Interactive CLI
./ssview.exe
# or just double-click ssview.exe

# Single action with a sql statement
./ssview.exe -s "select * from SalesLT.Customer"
```
Hit ```exit``` to exit interactive cli.

### Modes
ssview cli has two modes, "table" and "sql".
First, ssview runs with table mode.
If you want to change mode, hit ```table``` or ```sql```.
In both mode, listing is limted by ```limit``` value in db.json

#### Table Mode
```powershell
table? > {hit a table name you want to see}
```
In table mode, just hit a table name to list data in the table.

#### SQL Mode
```powershell
sql? > {write a sql statement you want to execute}
```
In SQL mode, you can execute a sql statement you want.

## db.json
db.json is config file to connect SQL Server.
```json
{
  "rdb": "sqlserver",
  "user": "hoge",
  "password": "fuga",
  "host": "localhost",
  "port": "1433",
  "database": "piyo",
  "limit": 50
}
```

Setting above provides a connection url below.
```
sqlserver://hoge:fuga@localhost:1433?database=piyo
```

If you set value "mysql" or "mariadb" on ```rdb```, you can connect to mysql or mariadb with ssview.
