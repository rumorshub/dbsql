package dbsql

import (
	"context"

	"github.com/roadrunner-server/endure/v2/dep"
	"github.com/roadrunner-server/errors"
)

const PluginName = "db.sql"

type Plugin struct {
	opener *DBOpener
}

func (p *Plugin) Init(cfg Configurer) error {
	const op = errors.Op("db.sql_plugin_init")

	if !cfg.Has(PluginName) {
		return errors.E(op, errors.Disabled)
	}

	var channelsCfg ChannelsConfig
	if err := cfg.UnmarshalKey(PluginName, &channelsCfg); err != nil {
		return errors.E(op, err)
	}

	p.opener = NewOpener()
	for name, config := range channelsCfg {
		p.opener.AddChannel(name, config)
	}

	return nil
}

func (p *Plugin) Name() string {
	return PluginName
}

func (p *Plugin) Serve() chan error {
	return make(chan error, 1)
}

func (p *Plugin) Stop(context.Context) error {
	return p.opener.Close()
}

func (p *Plugin) Provides() []*dep.Out {
	return []*dep.Out{
		dep.Bind((*Opener)(nil), p.DBOpener),
	}
}

func (p *Plugin) DBOpener() *DBOpener {
	return p.opener
}
