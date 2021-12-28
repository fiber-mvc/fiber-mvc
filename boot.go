package core

import (
	"github.com/gofiber/template/html"
	"log"
	"reflect"

	"github.com/fiber-mvc/fiber-mvc/database"
	"github.com/fiber-mvc/fiber-mvc/routing"
	"github.com/gofiber/fiber/v2"
	"github.com/goioc/di"
)

type setupViewEngine func(e interface{})

type Config struct {
	Database        database.DBConfig
	Fiber           fiber.Config
	Services        map[string]interface{}
	Routers         []routing.Router
	Controllers     map[string]interface{}
	SetupViewEngine setupViewEngine
	Debug           bool
}

var AppConfig Config

type _setup interface {
	Setup(app *fiber.App)
}

func register(name string, value interface{}) {
	_, err := di.RegisterBean(name, reflect.TypeOf(value))
	if err != nil {
		log.Fatal("Failed to register bean ", name, err)
	}
}

func setupIOC() {

	for key, s := range AppConfig.Services {
		log.Println("Registering service:", key)
		register(key, s)

	}

	for key, c := range AppConfig.Controllers {
		log.Println("Registering controller:", key)
		register(key, c)

	}

	err := di.InitializeContainer()
	if err != nil {
		log.Println("Failed to initialize container", err)
	}

}

func Boot(c Config) {
	AppConfig = c

	if AppConfig.Database.Driver != "" {
		database.Connect(AppConfig.Database)
	}

	if AppConfig.Fiber.Views == nil {
		log.Println("Setting up default view engine")
		engine := html.New("./views", ".tpl")
		engine.Reload(AppConfig.Debug)
		if AppConfig.SetupViewEngine != nil {
			AppConfig.SetupViewEngine(engine)
		}
		AppConfig.Fiber.Views = engine
	}

	app := fiber.New(AppConfig.Fiber)

	setupIOC()

	for _, service := range AppConfig.Services {
		service.(_setup).Setup(app)
	}
	for _, controller := range AppConfig.Controllers {
		controller.(_setup).Setup(app)
	}
	for _, router := range AppConfig.Routers {
		router.Setup(app)
	}

	log.Fatalln(app.Listen(":3000"))
}
