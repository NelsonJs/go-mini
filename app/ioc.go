package app

import (
	"reflect"

	"github.com/berkaroad/ioc"
	"github.com/spf13/viper"
)

func getIocContainer() ioc.Container {
	v := viper.Get("runtime.container")
	if v == nil {
		c := ioc.NewContainer()
		viper.Set("runtime.container", c)
		return c
	}
	return v.(ioc.Container)
}

type Container struct {
	ioc ioc.Container
}

func GetContainer() Container {
	return Container{
		ioc: getIocContainer(),
	}
}

func (c Container) Invoke(f interface{}) ([]reflect.Value, error) {
	return c.ioc.Invoke(f)
}

func (c Container) Register(val interface{}, ifacePtr interface{}) {
	c.ioc.RegisterTo(val, ifacePtr, ioc.Singleton)
}
