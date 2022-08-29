package cockroachdb

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

const (
	//COCKROACH_USERNAME=root;COCKROACH_PASSWORD=;COCKROACH_HOST=127.0.0.1;COCKROACH_PORT=26257;COCKROACH_SCHEMA=cbdc
	cockroachUsername = "COCKROACH_USERNAME"
	cockroachPassword = "COCKROACH_PASSWORD"
	cockroachHost     = "COCKROACH_HOST"
	cockroachPort     = "COCKROACH_PORT"
	cockroachSchema   = "COCKROACH_SCHEMA"
)

var (
	DBClient *gorm.DB

	username = os.Getenv(cockroachUsername)
	password = os.Getenv(cockroachPassword)
	host     = os.Getenv(cockroachHost)
	port     = os.Getenv(cockroachPort)
	schema   = os.Getenv(cockroachSchema)
)

func InitCockroachDB() {
	dataSourceName := fmt.Sprintf(
		"postgresql://%s@%s:%s/%s?sslmode=disable&charset=utf8&parseTime=true",
		username,
		host,
		port,
		schema,
	)

	var err error
	DBClient, err = gorm.Open(postgres.Open(dataSourceName+"&application_name=$ docs_simplecrud_gorm"), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("failed to connect with database, %w", err))
	}

	log.Println("Database successfully configured")
}
