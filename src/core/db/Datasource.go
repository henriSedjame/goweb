package db

import (
	"github.com/go-pg/pg/v10"
	"go.mongodb.org/mongo-driver/mongo"
	"goweb/src/core"
)

type Datasource struct {
	Properties *Properties
}

func (source *Datasource) CanStart() bool {
	return source.Properties != nil
}

func (source *Datasource) Stop(app *core.App) error {
	switch source.Properties.Type {
	case POSTGRES:
		if value := app.Context().Value(core.PsqlCtxKey); value != nil {
			if err := value.(pg.DB).Close(); err != nil {
				return err
			}
		}
	case MONGO:
		if value := app.Context().Value(core.MongoCtxKey); value != nil {
			client := value.(mongo.Client)
			if err := client.Disconnect(app.Context()); err != nil {
				return err
			}
		}
	}

	return nil
}

func (source *Datasource) Start(app *core.App) error {
	switch source.Properties.Type {
	case POSTGRES:
		if err := PsqlExecutors(app)(source.Properties); err != nil {
			return err
		}
	case MONGO:
		if err := MongoExecutors(app)(source.Properties); err != nil {
			return err
		}
	}

	return nil
}

func NewPsql(opts ...Option) *Datasource {
	props := Properties{
		Type: POSTGRES,
		Username: "",
		Password: "",
		Port: 5432,
		Host: "localhost",
	}

	for _, opt := range opts{
		opt(&props)
	}

	return &Datasource{
		Properties: &props,
	}
}

func NewMongo(opts ...Option) *Datasource {
	props := Properties{
		Type: MONGO,
		Username: "",
		Password: "",
		Port: 27017,
		Host: "localhost",
	}

	for _, opt := range opts{
		opt(&props)
	}

	return &Datasource{
		Properties: &props,
	}
}
