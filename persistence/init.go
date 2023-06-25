package persistence

import (
	"context"
	"database/sql"

	"github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateMongoConnection(ctx context.Context, connectionString string) (*mongo.Client, error){
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(connectionString).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func CreateSQLConnection(ctx context.Context, user, pass, addr, dbName string) (*sql.DB,error){
	cfg := mysql.Config{
		User: user,
		Passwd: pass,
		Net: "tcp",
		Addr: addr,
		DBName: "test",
	}
	// Get a database handle.
	var err error
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return nil, err
	}
	return db, nil
}