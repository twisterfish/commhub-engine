package datastores

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "os"
    "fmt"
)

var RDS *sql.DB
var ERROR error

//////////////////////////////////////////////////////////////////////////////////////////
// Opens a connection to Cloud SQL MySQL instance
// The connection parameters are specified in the app.yaml
//////////////////////////////////////////////////////////////////////////////////////////
func OpenRDS()( *sql.DB, error ) {
	
	var (
            connectionName = os.Getenv("CLOUDSQL_CONNECTION_NAME")
            user           = os.Getenv("CLOUDSQL_USER")
            password       = os.Getenv("CLOUDSQL_PASSWORD") 
        )
	
    // Open up our database connection.
    RDS, ERROR = sql.Open("mysql", fmt.Sprintf("%s:%s@cloudsql(%s)/", user, password, connectionName))
    
    // If there is an error opening the connection, handle it
    if ERROR != nil {
    	defer RDS.Close()       
    }
    
    return RDS, ERROR

}