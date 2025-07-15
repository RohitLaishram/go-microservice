package config

import (
	"context"
	"payment/model"

	// "database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"cloud.google.com/go/cloudsqlconn"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectWithGORM() (*gorm.DB, error) {
	mustGetenv := func(k string) string {
		v := os.Getenv(k)
		if v == "" {
			log.Fatalf("Fatal Error: %s environment variable not set.", k)
		}
		return v
	}

	dbUser := mustGetenv("DB_USER")
	dbPwd := mustGetenv("DB_PASS")
	dbName := mustGetenv("DB_NAME")
	instanceConnectionName := mustGetenv("INSTANCE_CONNECTION_NAME")
	usePrivate := os.Getenv("PRIVATE_IP")

	// Create Cloud SQL dialer
	d, err := cloudsqlconn.NewDialer(context.Background(), cloudsqlconn.WithLazyRefresh())
	if err != nil {
		return nil, fmt.Errorf("cloudsqlconn.NewDialer: %w", err)
	}

	var opts []cloudsqlconn.DialOption
	if usePrivate != "" {
		opts = append(opts, cloudsqlconn.WithPrivateIP())
	}

	// Register Cloud SQL connection with Go SQL
	mysqlDriver.RegisterDialContext("cloudsqlconn", func(ctx context.Context, addr string) (net.Conn, error) {
		return d.Dial(ctx, instanceConnectionName, opts...)
	})

	// DSN format
	dsn := fmt.Sprintf("%s:%s@cloudsqlconn(localhost:3306)/%s?parseTime=true",
		dbUser, dbPwd, dbName)

	// Initialize GORM with the DSN
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gorm.Open: %w", err)
	}
	db.AutoMigrate(&model.Payment{})

	return db, nil
}
