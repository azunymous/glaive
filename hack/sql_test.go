package hack

import (
	"fmt"
	"os"
	"testing"
)

func TestCreateSQLDSN(t *testing.T) {
	var (
		dbUser                 = "glv"                               // e.g. 'my-db-user'
		dbPwd                  = "<PASSWORD_GOES_HERE>"              // e.g. 'my-db-password'
		instanceConnectionName = "igiari-glv:us-central1:igiari-glv" // e.g. 'project:region:instance' (gcloud sql instances describe [INSTANCE_NAME])
		dbName                 = "asagi"                             // e.g. 'my-database'
	)

	socketDir, isSet := os.LookupEnv("DB_SOCKET_DIR")
	if !isSet {
		socketDir = "/cloudsql"
	}

	dbURI := fmt.Sprintf("%s:%s@unix(/%s/%s)/%s?parseTime=true", dbUser, dbPwd, socketDir, instanceConnectionName, dbName)
	println(dbURI)
}
