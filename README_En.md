# YoyoGo [中文介绍](https://github.com/yoyofx/yoyogo/blob/master/README-ZHCN.md "中文介绍")
YoyoGo is a simple, light and fast Web framework written in Go. 

![Release](https://img.shields.io/github/v/tag/yoyofx/yoyogo.svg?color=24B898&label=release&logo=github&sort=semver)
![Go](https://github.com/yoyofx/yoyogo/workflows/Go/badge.svg)
![GoVersion](https://img.shields.io/github/go-mod/go-version/maxzhang1985/yoyogo)
![DockerPull](https://img.shields.io/docker/pulls/maxzhang1985/yoyogo)
[![Report](https://goreportcard.com/badge/github.com/yoyofx/yoyogo)](https://goreportcard.com/report/github.com/maxzhang1985/yoyogo)
[![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg?color=24B898&logo=go&logoColor=ffffff)](https://godoc.org/github.com/yoyofx/yoyogo)
![Contributors](https://img.shields.io/github/contributors/yoyofx/yoyogo.svg)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)

# Features
- Pretty and fast router 
- Middleware Support
- Friendly to REST API
- No regexp or reflect
- Inspired by many excellent Go Web framework

[![](Resources/dingdingQR.jpg)](https://sourcerer.io/yoyofx)

# Installation
`go get github.com/yoyofx/yoyogo`

# Simple Example
```golang
package main
import ...

func main() {
    webHost := YoyoGo.CreateDefaultBuilder(func(router Router.IRouterBuilder) {
        router.GET("/info",func (ctx *Context.HttpContext) {
            ctx.JSON(200, Context.H{"info": "ok"})
        })
    }).Build().Run()       //default port :8080
}
```
![](Resources/yoyorun.jpg)


# ToDo
## Standard
* [X] Print Logo (YoyoGo)
* [X] Unified program entry (YoyoGo)
* [X] Simple router binded handle func
* [X] HttpContext (request,response)
* [X] Static File EndPoint（Static File Server）
* [X] JSON Serialization Struct （Context.H）
* [X] Get Request File And Save
* [X] Unite Get Post Data Away (form-data , x-www-form-urlencoded)
* [X] Binding Model for Http Request ( From, JSON ) 
### Response Render
* [X] Render Interface
* [X] JSON Render
* [X] JSONP Render
* [X] Indented Json Render
* [X] Secure Json Render
* [X] Ascii Json Render
* [X] Pure Json Render
* [X] Binary Data Render
* [X] TEXT
* [X] Protobuf
* [X] MessagePack
* [X] XML
* [X] YAML
* [X] File
* [X] Image
* [X] Template
* [X] Auto formater Render

## Middleware
* [X] Logger
* [X] StaticFile
* [X] Router
* [X] Router Middleware
* [X] CORS
* [X] JWT
* [X] Binding


## Router
* [x] GET、POST、HEAD、PUT、DELETE Method Support
* [x] Router Tree
* [x] Router Expression
* [x] RouteData (var)
* [x] Router Support Struct Refect Func Binded.
* [x] Router Support Prefix and Group Such as "/api/v1/endpoint"
* [X] Controller Router And Router Tempalte (Default)
* [X] Router Filter

## Dependency injection
* [X] Dependency injection Framework
* [X] Dependency injection Integration
* [X] Framework's factory and type in Dependency injection Integration

## Features
* [ ] configtion
* [ ] WebSocket
* [X] JWT 
* [ ] swagger
* [ ] GRpc
* [X] Prometheus 


# Advanced Example
```golang
package main
import ...

func main() {
	webHost := CreateCustomWebHostBuilder().Build()
	webHost.Run()
}

func CreateCustomBuilder() *Abstractions.HostBuilder {
    configuration := Abstractions.NewConfigurationBuilder().AddYamlFile("config").Build()
   	return YoyoGo.NewWebHostBuilder().
   		UseConfiguration(configuration).
   		Configure(func(app *YoyoGo.WebApplicationBuilder) {
   			app.UseMiddleware(Middleware.NewCORS())
   			app.UseMiddleware(Middleware.NewRequestTracker())
   			app.UseStaticAssets()
   			app.UseEndpoints(registerEndpointRouterConfig)
   			app.UseMvc(func(builder *Mvc.ControllerBuilder) {
   				//builder.AddViews(&View.Option{Path: "./Static/templates"})
   				builder.AddViewsByConfig()
   				builder.AddController(contollers.NewUserController)
   				builder.AddFilter("/v1/user/info", &contollers.TestActionFilter{})
   			})
   		}).
   		ConfigureServices(func(serviceCollection *DependencyInjection.ServiceCollection) {
   			serviceCollection.AddTransientByImplements(models.NewUserAction, new(models.IUserAction))
   			// Eureka.UseServiceDiscovery(serviceCollection)
   			//Consul.UseServiceDiscovery(serviceCollection)
   			Nacos.UseServiceDiscovery(serviceCollection)
   		}).
   		OnApplicationLifeEvent(getApplicationLifeEvent)
}

//region endpoint router config function
func registerEndpoints(router Router.IRouterBuilder) {
    Endpoints.UseHealth(router)
	Endpoints.UseViz(router)
	Endpoints.UsePrometheus(router)
	Endpoints.UsePprof(router)
	Endpoints.UseJwt(router)	

    router.GET("/error", func(ctx *Context.HttpContext) {
		panic("http get error")
	})

	router.POST("/info/:id", PostInfo)

	router.Group("/v1/api", func(router *Router.RouterGroup) {
		router.GET("/info", GetInfo)
	})

	router.GET("/info", GetInfo)
	router.GET("/ioc", GetInfoByIOC)
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
func GetInfo(ctx *Context.HttpContext) {
	ctx.JSON(200,  Context.H{"info": "ok"})
}

func GetInfoByIOC(ctx *Context.HttpContext) {
	var userAction models.IUserAction
	_ = ctx.RequiredServices.GetService(&userAction)
	ctx.JSON(200,  Context.H{"info": "ok " + userAction.Login("zhang")})
}

//HttpPost request: /info/:id ?q1=abc&username=123
func PostInfo(ctx *Context.HttpContext) {
	qs_q1 := ctx.Query("q1")
	pd_name := ctx.Param("username")

	userInfo := &UserInfo{}
	_ = ctx.Bind(userInfo)

	strResult := fmt.Sprintf("Name:%s , Q1:%s , bind: %s", pd_name, qs_q1, userInfo)

	ctx.JSON(200,  Context.H{"info": "hello world", "result": strResult})
}

func fireApplicationLifeEvent(life *YoyoGo.ApplicationLife) {
	printDataEvent := func(event YoyoGo.ApplicationEvent) {
		fmt.Printf("[yoyogo] Topic: %s; Event: %v\n", event.Topic, event.Data)
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

// Mvc 
type UserController struct {
	*Controller.ApiController
	userAction models.IUserAction    // IOC
}

// ctor for ioc
func NewUserController(userAction models.IUserAction) *UserController {
	return &UserController{userAction: userAction}
}

// reuqest param binder
type RegiserRequest struct {
	Controller.RequestParam
	UserName string `param:"username"`
	Password string `param:"password"`
}

// auto bind action param by ioc
func (this *UserController) Register(ctx *Context.HttpContext, request *RegiserRequest) ActionResult.IActionResult {
	result := Controller.ApiResult{Success: true, Message: "ok", Data: request}
	return ActionResult.Json{Data: result}
}

// use userAction interface by ioc  
func (this *UserController) GetInfo() Controller.ApiResult {
	return this.OK(this.userAction.Login("zhang"))
}

```
