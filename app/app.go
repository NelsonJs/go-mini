package app

import "sync"

type App struct {
	cls map[string]ServiceRegistry
}

type ServiceRegistry struct {
	Ctor     interface{}
	IfacePtr interface{}
}

var (
	application *App
	appOnce     sync.Once
)

func Application() *App {
	appOnce.Do(func() {
		if application == nil {
			application = &App{
				cls: make(map[string]ServiceRegistry),
			}
		}
	})
	return application
}

func (app *App) Register(serviceName string, ctor interface{}, ifacePtr interface{}) {
	app.cls[serviceName] = ServiceRegistry{
		Ctor:     ctor,
		IfacePtr: ifacePtr,
	}
}

func (app *App) Services() map[string]ServiceRegistry {
	return app.cls
}
