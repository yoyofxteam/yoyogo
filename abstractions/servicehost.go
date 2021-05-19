package abstractions

import (
	"fmt"
	"github.com/yoyofx/yoyogo"
	"github.com/yoyofx/yoyogo/abstractions/platform/consolecolors"
	"github.com/yoyofx/yoyogo/abstractions/servicediscovery"
	"github.com/yoyofx/yoyogo/abstractions/xlog"
	"github.com/yoyofx/yoyogo/utils"
	"strconv"
)

type IServiceHost interface {
	Run()
	Shutdown()
	StopApplicationNotify()
	SetAppMode(mode string)
}

var sdRegEnable bool = true

// host base
type ServiceHost struct {
	HostContext *HostBuilderContext
	Server      IServer
	logger      xlog.ILogger
}

func NewServiceHost(server IServer, hostContext *HostBuilderContext) ServiceHost {
	var sdconfig *servicediscovery.Config
	_ = hostContext.HostServices.GetService(&sdconfig)
	if sdconfig != nil {
		sdRegEnable = sdconfig.RegisterWithSelf
	}

	return ServiceHost{Server: server, HostContext: hostContext, logger: xlog.GetXLogger("Application")}
}

func (host ServiceHost) Run() {
	hostEnv := host.HostContext.HostingEnvironment
	host.logger.SetCustomLogFormat(nil)
	RunningHostEnvironmentSetting(hostEnv)
	PrintLogo(host.logger, hostEnv)
	//exithooksignals.HookSignals(host)
	HostRunning(host.logger, host.HostContext)
	//application running
	_ = host.Server.Run(host.HostContext)
	//application ending
	HostEnding(host.logger, host.HostContext)
}

func (host ServiceHost) Shutdown() {
	host.Server.Shutdown()
}

func (host ServiceHost) StopApplicationNotify() {
	HostEnding(host.logger, host.HostContext)
	host.HostContext.ApplicationCycle.StopApplication()
}

func (host ServiceHost) SetAppMode(mode string) {
	host.HostContext.HostingEnvironment.Profile = mode
}

func HostRunning(log xlog.ILogger, context *HostBuilderContext) {
	go hostStarting(log, context)
}

func HostEnding(log xlog.ILogger, context *HostBuilderContext) {
	hostEnding(log, context)
}

func hostStarting(log xlog.ILogger, context *HostBuilderContext) {
	//Service Discovery
	if sdRegEnable {
		var sd servicediscovery.IServiceDiscovery
		_ = context.HostServices.GetService(&sd)
		if sd != nil {
			_ = sd.Register()
		}
	}
	//---------------------------------------------------
	//Host Services
	var services []IHostService
	_ = context.HostServices.GetService(&services)
	for _, service := range services {
		_ = service.Run()
	}
}

func hostEnding(log xlog.ILogger, context *HostBuilderContext) {
	//Service Discovery

	var sdcache servicediscovery.Cache
	err := context.HostServices.GetService(&sdcache)
	if err == nil {
		sdcache.Stop()
	}

	if sdRegEnable {
		var sd servicediscovery.IServiceDiscovery
		err = context.HostServices.GetService(&sd)
		if err == nil && sd != nil {
			_ = sd.Destroy()
		}
	}
	//---------------------------------------------------
	//Host Services
	var services []IHostService
	_ = context.HostServices.GetService(&services)
	for _, service := range services {
		_ = service.Stop()
	}
}

func PrintLogo(l xlog.ILogger, env *HostEnvironment) {
	//logo, _ := base64.StdEncoding.DecodeString(yoyogo.Logo)
	logo := yoyogo.Logo
	fmt.Println(consolecolors.Blue(string(logo)))
	fmt.Println(" ")
	fmt.Printf("%s   (version:  %s)", consolecolors.Green(":: YoyoGo ::"), consolecolors.Blue(env.Version))

	fmt.Print(consolecolors.Blue(`
light and fast , dependency injection based micro-service framework written in Go.
`))

	fmt.Println(" ")
	l.Info(consolecolors.Green("Welcome to YoyoGo, starting application ..."))
	l.Info("yoyogo framework version :  %s", consolecolors.Blue(env.Version))
	l.Info("server & protocol        :  %s", consolecolors.Green(env.Server))
	l.Info("machine host ip          :  %s", consolecolors.Blue(env.Host))
	l.Info("listening on port        :  %s", consolecolors.Blue(env.Port))
	l.Info("application running pid  :  %s", consolecolors.Blue(strconv.Itoa(env.PID)))
	l.Info("application name         :  %s", consolecolors.Blue(env.ApplicationName))
	l.Info("application exec path    :  %s", consolecolors.Yellow(utils.GetCurrentDirectory()))
	l.Info("application config path  :  %s", consolecolors.Yellow(env.MetaData["config.path"]))
	l.Info("application environment  :  %s", consolecolors.Yellow(consolecolors.Blue(env.Profile)))
	l.Info("running in %s mode , change (Dev,tests,Prod) mode by HostBuilder.SetEnvironment .", consolecolors.Red(env.Profile))
	l.Info(consolecolors.Green("Starting server..."))
	l.Info("server setting map       :  %v", env.MetaData)
	l.Info(consolecolors.Green("Server is Started."))
}
