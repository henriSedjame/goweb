package core

import (
	"goweb/src/core/db"
	"goweb/src/core/server"
	"goweb/src/core/web"
)

type Properties struct {
	Server *server.Properties `json:"server"`
	Db *db.Properties `json:"db"`
	Cors *web.CorsProperties `json:"cors"`
	ActiveProfiles string `json:"profiles"`
}


func DefaultAppProperties() *Properties {

	return &Properties{
		Server: &server.Properties{
			Port: 8080,
			Host: "localhost",
		},
		Cors: &web.CorsProperties{
			AllowedHeaders: "*",
			AllowedMethods: "*",
			AllowedOrigins: "*",
		},
	}
}
