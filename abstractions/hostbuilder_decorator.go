package abstractions

import "github.com/yoyofxteam/dependencyinjection"

// IHostBuilderDecorator Host Builder decorator or extension
type IHostBuilderDecorator interface {

	// OverrideConfigure is configure function by application builder.
	OverrideConfigure(configureFunc interface{}, builder IApplicationBuilder)
	// OverrideNewApplicationBuilder create application builder.
	OverrideNewApplicationBuilder(context *HostBuilderContext) IApplicationBuilder
	// OverrideNewHost Create IServiceHost.
	OverrideNewHost(server IServer, context *HostBuilderContext) IServiceHost
	// OverrideInnerConfigures inner configures for IOC
	OverrideIOCInnerConfigures(serviceCollection *dependencyinjection.ServiceCollection)
}
