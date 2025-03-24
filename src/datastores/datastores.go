package datastores

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var RDS *sql.DB
var ERROR error

// ////////////////////////////////////////////////////////////////////////////////////////
// Opens a connection to Aurora MySQL instance
// The connection parameters are specified in the app.yaml
// ////////////////////////////////////////////////////////////////////////////////////////
func OpenRDS() (*sql.DB, error) {

	var (
		connectionName = os.Getenv("AURORASQL_CONNECTION_NAME")
		user           = os.Getenv("AURORASQL_USER")
		password       = os.Getenv("AURORASQL_PASSWORD")
	)

	/*	var (
			connectionName = "connection here"
			user           = "user here"
			password       = "password here"
		)
	*/
	// Open up our database connection.
	RDS, ERROR = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/", user, password, connectionName))

	// If there is an error opening the connection, handle it
	if ERROR != nil {
		defer RDS.Close()
	}

	return RDS, ERROR

}
