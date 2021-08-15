package main

import (
	"github.com/yoyofx/yoyogo/abstractions"
	"github.com/yoyofx/yoyogo/web"
	"github.com/yoyofx/yoyogo/web/actionresult/extension"
	"github.com/yoyofx/yoyogo/web/mvc"
	"github.com/yoyofxteam/dependencyinjection"
	"mvcdemo/controllers"
)

func main() {
	CreateMVCBuilder().Build().Run()
}

//* Create the builder of Web host
func CreateMVCBuilder() *abstractions.HostBuilder {
	configuration := abstractions.NewConfigurationBuilder().
		AddEnvironment().
		AddYamlFile("config").Build()

	return web.NewWebHostBuilder().
		UseConfiguration(configuration).
		Configure(func(app *web.ApplicationBuilder) {
			app.SetJsonSerializer(extension.CamelJson())
			app.UseMvc(func(builder *mvc.ControllerBuilder) {
				builder.AddViewsByConfig()                           //视图
				builder.AddController(controllers.NewDemoController) // 注册mvc controller
			})
		}).
		ConfigureServices(func(serviceCollection *dependencyinjection.ServiceCollection) {
			// ioc
		})
}
