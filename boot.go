package core

import (
	"log"
	"reflect"

	"github.com/fiber-mvc/fiber-mvc/database"
	"github.com/fiber-mvc/fiber-mvc/routing"
	"github.com/gofiber/fiber/v2"
	"github.com/goioc/di"
)

type Config struct {
	Database    database.DBConfig
	Fiber       fiber.Config
	Services    map[string]interface{}
	Routers     []routing.Router
	Controllers map[string]interface{}
}

type _setup interface {
	Setup(app *fiber.App)
}

func register(name string, value interface{}) {
	_, err := di.RegisterBean(name, reflect.TypeOf(value))
	if err != nil {
		log.Fatal("Failed to register bean ", name, err)
	}
}

func setupIOC(config Config) {

	for key, s := range config.Services {
		log.Println("Registering service:", key)
		register(key, s)

	}

	for key, c := range config.Controllers {
		log.Println("Registering controller:", key)
		register(key, c)

	}

	err := di.InitializeContainer()
	if err != nil {
		log.Println("Failed to initialize container", err)
	}

}

func Boot(config Config) {
	if config.Database.Driver != "" {
		database.Connect(config.Database)
	}

	app := fiber.New(config.Fiber)

	setupIOC(config)

	for _, service := range config.Services {
		service.(_setup).Setup(app)
	}
	for _, controller := range config.Controllers {
		controller.(_setup).Setup(app)
	}
	for _, router := range config.Routers {
		router.Setup(app)
	}

	log.Fatalln(app.Listen(":3000"))
}
