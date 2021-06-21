package core

import (
	"context"
	"encoding/json"
	"fmt"
	"goweb/src/core/db"
	"goweb/src/web"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

type App struct {
	server      *http.Server
	properties  *Properties
	datasource  StartAndStoppable
	Logger      *log.Logger
	classpath   string
	ctx         AppCtx
	controllers web.RestController
}

type AppOptions func(app *App)

func Rest(controllers web.RestController) AppOptions {
	return func(app *App) {
		app.controllers = controllers
	}
}

func DB(datasource StartAndStoppable) AppOptions {
	return func(app *App) {
		app.datasource = datasource
	}
}

func New(logger *log.Logger, options ...AppOptions) *App {
	app := &App{Logger: logger}
	for _, opt := range options {
		opt(app)
	}
	return app
}

func (app *App) WithOption(option AppOptions) *App {
	option(app)
	return app
}

// Run the application
func (app *App) Run() {

	// set the application context
	app.WithContext(context.Background())

	// set aplication classpath
	if err := app.setClasspath(); err != nil {
		app.Logger.Fatal(err)
	}

	// load application properties
	if err := app.loadProperties(); err != nil {
		app.Logger.Fatal(err)
	}

	// set datasource
	if app.datasource == nil {
		app.setDatasource()
	}

	if app.datasource.CanStart() {
		// start application databases
		if err := app.datasource.Start(app); err != nil {
			app.Logger.Fatal(err)
		}
	}
	
	// configure application server
	//app.configureServer()

	// launch application server
	go func() {
		app.Logger.Fatal(app.server.ListenAndServe())
	}()

	app.shutDown()
}

// shutDown the application
func (app *App) shutDown() {

	signalChannel := make(chan os.Signal)

	// Envoyer un signal au channel lors :
	// * d'une interruption de la machine
	// * d'un arrêt de la machine
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, os.Kill)

	// Récupérer le signal envoyé
	_ = <- signalChannel

	app.Logger.Println(" ===> Arrêt du serveur")

	deadline, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	// Stop database
	if err := app.datasource.Stop(app); err != nil {
		app.Logger.Fatal(err)
	}

	// ShutDown server
	if err := app.server.Shutdown(deadline); err != nil {
		app.Logger.Fatal(err)
	}

}

// setClasspath define app classpath
func (app *App) setClasspath() error {
	if len(app.classpath) == 0 {
		if rootDir, err := os.Getwd(); err != nil {
			return err
		} else {
			app.classpath = fmt.Sprintf("%s/src/resources", rootDir)

			app.Logger.Printf(" ===> Application classpath configured to : %s", app.classpath)
		}
	}

	return nil
}

// loadProperties load the application properties
func (app *App) loadProperties() error {

	props := DefaultAppProperties()

	if err := app.setProfileProps("", props); err != nil {
		return err
	} else {

		profiles := app.properties.ActiveProfiles

		if profiles != "" {

			app.Logger.Printf(" ===> Application active profiles : %s", profiles)

			for _, profile := range strings.Split(profiles, ",") {
				if err := app.setProfileProps(profile, app.properties); err != nil {
					return err
				}
			}
		}
	}

	app.Logger.Println(" ===> Application Properties retrieval succeeded.")

	return nil
}

// setProfileProps
func (app *App) setProfileProps(profile string, props *Properties) error {

	suffix := "application"

	if profile != "" {
		suffix = fmt.Sprintf("%s-%s", suffix, profile)
	}

	propertiesLocation := fmt.Sprintf("%s/%s.json", app.classpath, suffix)

	if propertiesFile, err := os.OpenFile(propertiesLocation, os.O_RDONLY, 0777); err != nil {
		return err
	} else {
		defer  propertiesFile.Close()

		if bytes, err := ioutil.ReadAll(propertiesFile); err != nil {
			app.Logger.Println("An error occured while retrieving application properties")
			return err
		} else {
			if err := json.Unmarshal(bytes, &props); err != nil {
				app.Logger.Println("An error occured while retrieving application properties")
				return err
			}
		}
	}

	app.properties = props

	return nil
}

func (app *App) setDatasource() {
	props := app.properties.Db
	app.datasource = &db.Datasource{
		Properties: props,
	}
}

// Context get the app context
func (app App) Context() AppCtx{
	return app.ctx
}

// WithContext : set new application context
func (app *App) WithContext(ctx AppCtx)  {
	app.ctx = ctx
}
