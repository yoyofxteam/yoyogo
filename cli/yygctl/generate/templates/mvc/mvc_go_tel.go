package mvc

const DemoController_Tel = `
package {{.CurrentModelName}}

import (
	"github.com/yoyofx/yoyogo/web/context"
	"github.com/yoyofx/yoyogo/web/mvc"
)

type DemoController struct {
	mvc.ApiController // 必须继承
}

func NewDemoController() *DemoController {
	return &DemoController{}
}

//-------------------------------------------------------------------------------
type RegisterRequest struct {
	mvc.RequestBody
	UserName string` + " `param:\"UserName\"`\n" +
	"    Password string" + "`param:\"Password\"`\n" +
	`}

//GET URL  http://localhost:8080/app/v1/demo/register?UserName=max&Password=123
func (controller DemoController) Register(ctx *context.HttpContext, request *RegisterRequest) mvc.ApiResult {
	return mvc.ApiResult{Success: true, Message: "ok", Data: request}
}

//GET URL http://localhost:8080/app/v1/demo/getinfo
func (controller DemoController) GetInfo() mvc.ApiResult {
	return controller.OK("ok")
}

`

const Main_Tel = `
package {{.CurrentModelName}}

import (
	"github.com/yoyofx/yoyogo/abstractions"
	"github.com/yoyofxteam/dependencyinjection"
	"github.com/yoyofx/yoyogo/web"
	"github.com/yoyofx/yoyogo/web/actionresult/extension"
	"github.com/yoyofx/yoyogo/web/mvc"
    "{{.ModelName}}/controller"
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
				builder.AddController(controller.NewDemoController) // 注册mvc controller
			})
		}).
		ConfigureServices(func(serviceCollection *dependencyinjection.ServiceCollection) {
			// ioc
		})
}
`
const Mod_Tel = `

module {{.ModelName}}


go 1.16

require (
	github.com/yoyofxteam/dependencyinjection v1.0.0
	github.com/yoyofx/yoyogo {{.Version}}
)
`

const Config_Tel = `
yoyogo:
  application:
    name: yoyogo_demo_dev
    metadata: "develop"
    server:
      type: "fasthttp"
      address: ":8080"
      path: "app"
      max_request_size: 2096157
      session:
        name: "YOYOGO_SESSIONID"
        timeout: 3600
      tls:
        cert: ""
        key: ""
      mvc:
        template: "v1/{controller}/{action}"
        views:
          path: "./static/templates"
          includes: [ "","" ]
      static:
        patten: "/"
        webroot: "./static"
      jwt:
        header: "Authorization"
        secret: "12391JdeOW^%$#@"
        prefix: "Bearer"
        expires: 3
        enable: false
        skip_path: [
            "/info",
            "/v1/user/GetInfo",
            "/v1/user/GetSD"
        ]
      cors:
        allow_origins: ["*"]
        allow_methods: ["POST","GET","PUT", "PATCH"]
        allow_credentials: true
  cloud:
    apm:
      skyworking:
        address: localhost:11800
    discovery:
      type: "nacos"
      metadata:
        url: "127.0.0.1"
        port: 80
        namespace: "public"
        group_name: ""
    #    clusters: [""]
#      type: "consul"
#      metadata:
#        address: "localhost:8500"
#        health_check: "/actuator/health"
#        tags: [""]
#      type: "eureka"
#      metadata:
#        address: "http://localhost:5000/eureka"
  datasource:
      mysql:
        name: db1
        url: tcp(localhost:10042)/xxx?charset=utf8&parseTime=True
        username: root
        password: root
        debug: true
        pool:
          init_cap: 2
          max_cap: 5
          idle_timeout : 5
      redis:
        name: reids1
        url: 127.0.0.1:31379
        password:
        db: 0
        pool:
          init_cap: 2
          max_cap: 5
          idle_timeout: 5
`
