package db

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"goweb/src/core"
	"goweb/src/data"
)

type Executors func(dbProps *Properties) error

func PsqlExecutors(app *core.App) Executors {
	return func(dbProps *Properties) error {

		// Get entities from app context
		models := app.Context().Value(core.ModelsCtxKey).([]data.Table)

		// If database name is filled
		if dbProps.Dbname != "" {

			// Create pqsl url
			optStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
				dbProps.Username, dbProps.Password, dbProps.Host, dbProps.Port, dbProps.Dbname)

			// Try to parse url
			if opt, err := pg.ParseURL(optStr); err != nil {
				return err
			} else {

				// Connect the database
				db := pg.Connect(opt)

				// Store database into the app context
				// by creating a new context based on app context
				ctx := context.WithValue(app.Context(), core.PsqlCtxKey, db)

				// Save the new context
				app.WithContext(ctx)

				// Try to query psql version
				// in order to verify if database is working
				var version string
				if _, err := db.QueryOneContext(ctx, pg.Scan(&version), "SELECT version()"); err != nil {
					return err
				} else {

					// Log message to notify database connection succeeded
					app.Logger.Println(" ===> Connection to Database succeeded.")
					app.Logger.Printf(" ===> Database version : %s", version)

					// Try to create a table for each entities
					for _, mod := range models {
						if err := db.Model(mod).CreateTable(&orm.CreateTableOptions{
							Temp:        false,
							IfNotExists: true,
						}); err != nil {
							return err
						}
					}
				}
			}
		} else {
			return core.AppError{
				Message: "Database is not set. Please consider to set property { \"db\" : {\"database\" : ****}}",
			}
		}
		return nil
	}
}

func MongoExecutors(app *core.App)  Executors{
	return func(dbProps *Properties) error {

		// Get entities from app context
		models := app.Context().Value(core.ModelsCtxKey).([]data.Table)

		if dbProps.Dbname != "" {

			// Create mongo connect uri
			uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s",
				dbProps.Username, dbProps.Password, dbProps.Host, dbProps.Port, dbProps.Dbname)

			// Create client options
			clientOptions := options.Client().ApplyURI(uri)

			// Connect the database
			if client, err := mongo.Connect(app.Context(), clientOptions); err != nil {
				return err
			} else {

				// Store database into the app context
				// by creating a new context based on app context
				ctx := context.WithValue(app.Context(), core.MongoCtxKey, client)

				// Save the new context
				app.WithContext(ctx)

				// Check the connection
				if err := client.Ping(app.Context(), nil); err != nil {
					return err
				} else {
					// Log message to notify database connection succeeded
					app.Logger.Println(" ===> Connection to Database succeeded.")

					for _, mod := range models {
						if err := client.Database(dbProps.Dbname).CreateCollection(app.Context(), mod.Name(), &options.CreateCollectionOptions{

						}); err != nil {
							return err
						}
					}
				}
			}

		} else {
			return core.AppError{
				Message: "Database is not set. Please consider to set property { \"db\" : {\"database\" : ****}}",
			}
		}

		return nil
	}
}



