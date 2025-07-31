package config

import (
	"context"
	"order/model"

	// "database/sql"
	"fmt"
	"net"
	"os"

	"cloud.google.com/go/cloudsqlconn"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectWithGORM() (*gorm.DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPwd := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	instanceConnectionName := os.Getenv("INSTANCE_CONNECTION_NAME")
	usePrivate := os.Getenv("PRIVATE_IP")

	fmt.Printf("DB_USER: %s\nDB_PASS: %s\nDB_NAME: %s\nINSTANCE_CONNECTION_NAME: %s\nPRIVATE_IP: %s\n",
		dbUser, dbPwd, dbName, instanceConnectionName, usePrivate)

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
	dsn := fmt.Sprintf("%s:%s@cloudsqlconn(34.9.104.215:3306)/%s?parseTime=true",
		dbUser, dbPwd, dbName)

	// Initialize GORM with the DSN
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gorm.Open: %w", err)
	}
	db.AutoMigrate(&model.Order{})

	return db, nil
}
