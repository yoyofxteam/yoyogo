package main

import (
	"github.com/yoyofx/yoyogo/Abstractions"
	"github.com/yoyofx/yoyogo/WebFramework"
	"github.com/yoyofx/yoyogo/WebFramework/Endpoints"
	"github.com/yoyofx/yoyogo/WebFramework/Router"
)

func main() {
	configuration := Abstractions.NewConfigurationBuilder().AddYamlFile("config").Build()
	// --profile=prod , to read , config.yml
	YoyoGo.NewWebHostBuilder().
		UseConfiguration(configuration).
		Configure(func(app *YoyoGo.WebApplicationBuilder) {
			app.UseEndpoints(func(router Router.IRouterBuilder) {
				router.POST("/alert", PostAlert)
				Endpoints.UsePrometheus(router)
			})
		}).Build().Run()
}
