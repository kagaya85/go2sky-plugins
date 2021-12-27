# Go2Sky with database/sql/driver

## Installation

```bash
go get -u github.com/SkyAPM/go2sky-plugins/sql/driver
```

## Usage

```go
import (
    "database/sql"

	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/reporter"
	sqlPlugin "github.com/SkyAPM/go2sky-plugins/sql/driver"
	"github.com/go-sql-driver/mysql"
)

// init reporter
re, err := reporter.NewLogReporter()
defer re.Close()

// init tracer
tracer, err := go2sky.NewTracer("service-name", go2sky.WithReporter(re))

// register go2sky sql wrapper
sql.Register("skywalking-sql", sqlPlugin.NewTracerDriver(&mysql.MySQLDriver{}, tracer, sqlPlugin.WithSqlDBType(sqlPlugin.MYSQL), sqlPlugin.WithQueryReport()))
db, err := sql.Open("skywalking-sql", "user:password@tcp(127.0.0.1:3306)/dbname")

// use db handle as usual.
```