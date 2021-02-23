package registry

import (
	"github.com/gotomicro/ego-component/ekubernetes"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
)

type Option func(c *Container)

type Container struct {
	config *Config
	name   string
	logger *elog.Component
	client *ekubernetes.Component
}

func DefaultContainer() *Container {
	return &Container{
		config: DefaultConfig(),
		logger: elog.EgoLogger.With(elog.FieldComponent(ekubernetes.PackageName)),
	}
}

func Load(key string) *Container {
	c := DefaultContainer()
	if err := econf.UnmarshalKey(key, &c.config); err != nil {
		c.logger.Panic("parse Config error", elog.FieldErr(err), elog.FieldKey(key))
		return c
	}
	c.logger = c.logger.With(elog.FieldComponentName(key))
	c.name = key
	return c
}

func WithClientKubernetes(kubernetes *ekubernetes.Component) Option {
	return func(c *Container) {
		c.client = kubernetes
	}
}

// Build ...
func (c *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(c)
	}
	if c.client == nil {
		if c.config.OnFailHandle == "panic" {
			c.logger.Panic("client kubernetes nil", elog.FieldKey("use WithKubernetes method"))
		} else {
			c.logger.Error("client kubernetes nil", elog.FieldKey("use WithKubernetes method"))
		}
	}
	return newComponent(c.name, c.config, c.logger, c.client)
}
