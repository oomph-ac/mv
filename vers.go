package vers

import (
	"github.com/df-mc/dragonfly/server"
	"github.com/sandertv/gophertunnel/minecraft"
)

// Vers is an instance for binding a multi-version Dragonfly server.
type Vers struct {
	addr string
}

// New creates a new Vers instance.
func New(localAddr string) *Vers {
	return &Vers{
		addr: localAddr,
	}
}

// Listen listens for incoming connections on the address.
func (v *Vers) Listen(conf *server.Config, name string, protocols []minecraft.Protocol, requirePacks bool) {
	conf.Listeners = nil
	conf.Listeners = append(conf.Listeners, func(_ server.Config) (server.Listener, error) {
		l, err := minecraft.ListenConfig{
			StatusProvider:       minecraft.NewStatusProvider(name, "???"),
			ResourcePacks:        conf.Resources,
			TexturePacksRequired: requirePacks,
			AcceptedProtocols:    protocols,
		}.Listen("raknet", v.addr)
		if err != nil {
			return nil, err
		}

		conf.Log.Infof("Server running on %v.", l.Addr())

		return listener{
			Listener: l,
		}, nil
	})
}
