package Abstractions

import (
	"encoding/base64"
	"fmt"
	"github.com/yoyofx/yoyogo"
	"github.com/yoyofx/yoyogo/Abstractions/Platform/ConsoleColors"
	"github.com/yoyofx/yoyogo/Abstractions/xlog"
	"github.com/yoyofx/yoyogo/Utils"
	"github.com/yoyofx/yoyogo/WebFramework/Context"
	"strconv"
)

type IServiceHost interface {
	Run()
	Shutdown()
	StopApplicationNotify()
	SetAppMode(mode string)
}

func PrintLogo(l xlog.ILogger, env *Context.HostEnvironment) {
	logo, _ := base64.StdEncoding.DecodeString(YoyoGo.Logo)

	fmt.Println(ConsoleColors.Blue(string(logo)))
	fmt.Printf("%s                   (%s)", ConsoleColors.Green(":: YoyoGo ::"), ConsoleColors.Blue(env.Version))
	fmt.Println(" ")
	fmt.Println(" ")
	l.Debug(ConsoleColors.Green("Welcome to YoyoGo, starting application ..."))
	l.Debug("yoyogo framework version :  %s", ConsoleColors.Blue(env.Version))
	l.Debug("machine host ip          :  %s", ConsoleColors.Blue(env.Host))
	l.Debug("listening on port        :  %s", ConsoleColors.Blue(env.Port))
	l.Debug("application running pid  :  %s", ConsoleColors.Blue(strconv.Itoa(env.PID)))
	l.Debug("application name         :  %s", ConsoleColors.Blue(env.ApplicationName))
	l.Debug("application environment  :  %s", ConsoleColors.Blue(env.Profile))
	l.Debug("application exec path    :  %s", ConsoleColors.Yellow(Utils.GetCurrentDirectory()))
	l.Debug("running in %s mode , change (Dev,Test,Prod) mode by HostBuilder.SetEnvironment .", ConsoleColors.Blue(env.Profile))
	l.Debug(ConsoleColors.Green("Starting HTTP server..."))

}
