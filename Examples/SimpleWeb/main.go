package main

import (
	"SimpleWeb/contollers"
	"SimpleWeb/models"
	"fmt"
	"github.com/yoyofx/yoyogo/abstractions"
	"github.com/yoyofx/yoyogo/abstractions/xlog"
	"github.com/yoyofx/yoyogo/dependencyinjection"
	_ "github.com/yoyofx/yoyogo/pkg/datasources/mysql"
	_ "github.com/yoyofx/yoyogo/pkg/datasources/redis"
	"github.com/yoyofx/yoyogo/pkg/servicediscovery/nacos"
	web "github.com/yoyofx/yoyogo/web"
	"github.com/yoyofx/yoyogo/web/actionresult/extension"
	"github.com/yoyofx/yoyogo/web/context"
	"github.com/yoyofx/yoyogo/web/endpoints"
	"github.com/yoyofx/yoyogo/web/middlewares"
	"github.com/yoyofx/yoyogo/web/mvc"
	"github.com/yoyofx/yoyogo/web/router"
)

func SimpleDemo() {
	web.CreateHttpBuilder(func(router router.IRouterBuilder) {
		endpoints.UsePrometheus(router)
		registerEndpointRouterConfig(router)

		router.GET("/info", func(ctx *context.HttpContext) {
			ctx.JSON(200, context.H{"info": "ok"})
		})
	}).Build().Run()
}

func main() {
	//SimpleDemo()

	webHost := CreateCustomBuilder().Build()
	webHost.Run()
}

//* Create the builder of Web host
func CreateCustomBuilder() *abstractions.HostBuilder {

	configuration := abstractions.NewConfigurationBuilder().
		AddEnvironment().
		AddYamlFile("config").Build()

	return web.NewWebHostBuilder().
		UseConfiguration(configuration).
		Configure(func(app *web.ApplicationBuilder) {
			app.UseMiddleware(middlewares.NewCORS())
			//web.UseMiddleware(middlewares.NewRequestTracker())
			app.UseStaticAssets()
			app.UseEndpoints(registerEndpointRouterConfig)
			app.SetJsonSerializer(extension.CamelJson())
			app.UseMvc(func(builder *mvc.ControllerBuilder) {
				//builder.AddViews(&view.Option{Path: "./Static/templates"})
				builder.AddViewsByConfig()
				builder.AddController(contollers.NewUserController)
				builder.AddFilter("/v1/user/info", &contollers.TestActionFilter{})
			})
		}).
		ConfigureServices(func(serviceCollection *dependencyinjection.ServiceCollection) {
			serviceCollection.AddTransientByImplements(models.NewUserAction, new(models.IUserAction))

			//eureka.UseServiceDiscovery(serviceCollection)
			//consul.UseServiceDiscovery(serviceCollection)
			nacos.UseServiceDiscovery(serviceCollection)
		}).
		OnApplicationLifeEvent(getApplicationLifeEvent)
}

//*/

//region router config function
func registerEndpointRouterConfig(routerBuilder router.IRouterBuilder) {
	endpoints.UseHealth(routerBuilder)
	endpoints.UseViz(routerBuilder)
	endpoints.UsePrometheus(routerBuilder)
	endpoints.UsePprof(routerBuilder)
	endpoints.UseReadiness(routerBuilder)
	endpoints.UseLiveness(routerBuilder)
	endpoints.UseJwt(routerBuilder)

	routerBuilder.GET("/error", func(ctx *context.HttpContext) {
		panic("http get error")
	})

	routerBuilder.POST("/info/:id", PostInfo)

	routerBuilder.Group("/v1/api", func(routergroup *router.RouterGroup) {
		routergroup.GET("/info", GetInfo)
	})

	routerBuilder.GET("/", GetInfo)

	routerBuilder.GET("/info", GetInfo)
	routerBuilder.GET("/ioc", GetInfoByIOC)
}

//endregion

//region Http Request Methods

type UserInfo struct {
	UserName string `param:"username"`
	Number   string `param:"q1"`
	Id       string `param:"id"`
}

//HttpGet request: /info  or /v1/api/info
//bind UserInfo for id,q1,username
func GetInfo(ctx *context.HttpContext) {
	ctx.JSON(200, context.H{"info": "ok"})
}

func GetInfoByIOC(ctx *context.HttpContext) {
	var userAction models.IUserAction
	_ = ctx.RequiredServices.GetService(&userAction)
	ctx.JSON(200, context.H{"info": "ok " + userAction.Login("zhang")})
}

//HttpPost request: /info/:id ?q1=abc&username=123
func PostInfo(ctx *context.HttpContext) {
	qs_q1 := ctx.Input.Query("q1")
	pd_name := ctx.Input.Param("username")
	id := ctx.Input.Param("id")
	userInfo := &UserInfo{}
	_ = ctx.Bind(userInfo)

	strResult := fmt.Sprintf("Name:%s , Q1:%s , bind: %s , routeData id:%s", pd_name, qs_q1, userInfo, id)

	ctx.JSON(200, context.H{"info": "hello world", "result": strResult})
}

func getApplicationLifeEvent(life *abstractions.ApplicationLife) {
	printDataEvent := func(event abstractions.ApplicationEvent) {
		xlog.GetXLogger("Application Life Event:").Debug("Topic: %s; Event: %v", event.Topic, event.Data)
		//fmt.Printf("[yoyogo] Topic: %s; Event: %v\n", event.Topic, event.Data)
	}

	for {
		select {
		case ev := <-life.ApplicationStarted:
			go printDataEvent(ev)
		case ev := <-life.ApplicationStopped:
			go printDataEvent(ev)
			break
		}
	}
}

//endregion
