package main

import (
	"github.com/yoyofx/yoyogo/abstractions"
	"github.com/yoyofx/yoyogo/web"
	"github.com/yoyofx/yoyogo/web/endpoints"
	"github.com/yoyofx/yoyogo/web/router"
)

func main() {
	configuration := abstractions.NewConfigurationBuilder().AddYamlFile("config").Build()
	// --profile=prod , to read , config.yml
	YoyoGo.NewWebHostBuilder().
		UseConfiguration(configuration).
		Configure(func(app *YoyoGo.WebApplicationBuilder) {
			app.UseEndpoints(func(router router.IRouterBuilder) {
				router.POST("/alert", PostAlert)
				endpoints.UsePrometheus(router)
			})
		}).Build().Run()
}
