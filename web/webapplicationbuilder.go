package web

import (
	"github.com/yoyofx/yoyogo/abstractions"
	"github.com/yoyofx/yoyogo/web/actionresult"
	"github.com/yoyofx/yoyogo/web/actionresult/extension"
	"github.com/yoyofx/yoyogo/web/context"
	"github.com/yoyofx/yoyogo/web/middlewares"
	"github.com/yoyofx/yoyogo/web/mvc"
	"github.com/yoyofx/yoyogo/web/router"
	"github.com/yoyofx/yoyogo/web/view"
	"net/http"
)

// HTTP methods

//application builder struct
type ApplicationBuilder struct {
	hostContext     *abstractions.HostBuilderContext // host build 's context
	routerBuilder   router.IRouterBuilder            // route builder of interface
	middleware      middleware
	handlers        []MiddlewareHandler
	routeConfigures []func(router.IRouterBuilder)          // endpoints router configure functions
	mvcConfigures   []func(builder *mvc.ControllerBuilder) // mvc router configure functions
}

// create classic application builder
func UseClassic() *ApplicationBuilder {
	return &ApplicationBuilder{}
}

//region Create the builder of Web host
func CreateHttpBuilder(routerConfig func(router router.IRouterBuilder)) *abstractions.HostBuilder {
	return NewWebHostBuilder().
		UseServer(DefaultHttpServer(DefaultAddress)).
		Configure(func(app *ApplicationBuilder) {
			app.UseStatic("/Static", "./Static")
			app.UseEndpoints(routerConfig)
		})
}

func CreateMvcBuilder(appFunc func(app *ApplicationBuilder)) *abstractions.HostBuilder {
	configuration := abstractions.NewConfigurationBuilder().
		AddEnvironment().
		AddYamlFile("config").Build()
	return NewWebHostBuilder().
		UseConfiguration(configuration).
		Configure(func(app *ApplicationBuilder) {
			app.UseStaticAssets()
			app.UseMvc(func(builder *mvc.ControllerBuilder) {
				builder.AddViewsByConfig()
			})
		}).Configure(appFunc)
}

func CreateBlankWebBuilder() *WebHostBuilder {
	return NewWebHostBuilder()
}

// create application builder when combo all handlers to middleware
func New(handlers ...MiddlewareHandler) *ApplicationBuilder {
	return &ApplicationBuilder{
		handlers: handlers,
	}
}

// create new web application builder
func NewWebApplicationBuilder() *ApplicationBuilder {
	routerBuilder := router.NewRouterBuilder()
	recovery := middlewares.NewRecovery()
	logger := middlewares.NewLogger()
	router := middlewares.NewRouter(routerBuilder)
	jwt := middlewares.NewJwt()
	self := New(logger, recovery, jwt, router)
	self.routerBuilder = routerBuilder
	actionresult.SetJsonSerializeEncoder(extension.DefaultJsonEncoder{})
	return self
}

// SetJsonSerializer set json serializer for response
func (self *ApplicationBuilder) SetJsonSerializer(encoder extension.Encoder) *ApplicationBuilder {
	actionresult.SetJsonSerializeEncoder(encoder)
	return self
}

// UseMvc after create builder , apply router and logger and recovery middleware
func (self *ApplicationBuilder) UseMvc(configure func(builder *mvc.ControllerBuilder)) *ApplicationBuilder {
	if !self.routerBuilder.IsMvc() {
		self.routerBuilder.UseMvc(true)
	}
	self.mvcConfigures = append(self.mvcConfigures, configure)
	return self
}

func (self *ApplicationBuilder) UseEndpoints(configure func(router.IRouterBuilder)) *ApplicationBuilder {
	self.routeConfigures = append(self.routeConfigures, configure)
	return self
}

func (this *ApplicationBuilder) buildEndPoints() {
	this.routerBuilder.SetConfiguration(this.hostContext.Configuration)
	for _, configure := range this.routeConfigures {
		configure(this.routerBuilder)
	}
}

func (this *ApplicationBuilder) buildMvc() {
	if this.routerBuilder.IsMvc() {
		controllerBuilder := this.routerBuilder.GetMvcBuilder()
		// add config for controller builder
		controllerBuilder.SetConfiguration(this.hostContext.Configuration)
		for _, configure := range this.mvcConfigures {
			configure(controllerBuilder)
		}
		// add view engine
		var viewEngine view.IViewEngine
		err := this.hostContext.HostServices.GetServiceByName(&viewEngine, "viewEngine")
		if err == nil {
			option := this.routerBuilder.GetMvcBuilder().GetRouterHandler().Options.ViewOption
			if option != nil {
				viewEngine.SetTemplatePath(option)
				controllerBuilder.SetViewEngine(viewEngine)
			}
		}

		// add controllers to application services
		controllerDescriptorList := controllerBuilder.GetControllerDescriptorList()
		for _, descriptor := range controllerDescriptorList {
			this.hostContext.
				ApplicationServicesDef.
				AddSingletonByNameAndImplements(descriptor.ControllerName, descriptor.ControllerType, new(mvc.IController))
		}
	}
}

func (this *ApplicationBuilder) buildMiddleware() {
	for _, handler := range this.handlers {
		if configurationMdw, ok := handler.(middlewares.IConfigurationMiddleware); ok {
			configurationMdw.SetConfiguration(this.hostContext.Configuration)
		}
	}
	this.middleware = build(this.handlers)
}

//  this time is not build host.context.HostServices , that add services define
func (this *ApplicationBuilder) innerConfigures() {
	this.hostContext.
		ApplicationServicesDef.
		AddSingletonByNameAndImplements("viewEngine", view.CreateViewEngine, new(view.IViewEngine))
	//-------------------------  view engine ----------------------------------
}

// build and combo all middleware to request delegate (ServeHTTP(w http.ResponseWriter, r *http.Request))
// return abstractions.IRequestDelegate type
func (this *ApplicationBuilder) Build() interface{} {
	if this.hostContext == nil {
		panic("hostContext is nil! please set.")
	}
	this.buildMiddleware()
	this.buildEndPoints()
	this.buildMvc()
	return this
}

func (this *ApplicationBuilder) SetHostBuildContext(context *abstractions.HostBuilderContext) {
	this.hostContext = context
	// has host.context.HostServices
	if this.hostContext.ApplicationServicesDef != nil {
		this.innerConfigures()
	}
}

// apply middleware in builder
func (app *ApplicationBuilder) UseMiddleware(handler MiddlewareHandler) {
	if handler == nil {
		panic("handler cannot be nil")
	}
	app.handlers = append(app.handlers, handler)
}

// apply static middleware in builder
func (app *ApplicationBuilder) UseStatic(patten string, path string) {
	app.UseMiddleware(middlewares.NewStatic(patten, path))
}

func (app *ApplicationBuilder) UseStaticAssets() {
	app.UseMiddleware(middlewares.NewStaticWithConfig(app.hostContext.Configuration))
}

// apply handler middleware in builder
func (app *ApplicationBuilder) UseHandler(handler http.Handler) {
	app.UseMiddleware(wrap(handler))
}

// apply handler func middleware in builder
func (app *ApplicationBuilder) UseHandlerFunc(handlerFunc func(rw http.ResponseWriter, r *http.Request)) {
	app.UseMiddleware(wrapFunc(handlerFunc))
}

// apply handler func middleware in builder
func (app *ApplicationBuilder) UseFunc(handlerFunc MiddlewareHandlerFunc) {
	app.UseMiddleware(handlerFunc)
}

/*
middlewares of Server MiddlewareHandler , request port.
*/
func (app *ApplicationBuilder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.middleware.Invoke(context.NewContext(w, r, app.hostContext.ApplicationServices))
}
