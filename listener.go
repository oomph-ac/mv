package vers

import (
	"github.com/df-mc/dragonfly/server/session"
	"github.com/sandertv/gophertunnel/minecraft"
)

// listener is a custom minecraft.Listener for multi-version support.
type listener struct {
	*minecraft.Listener
}

// Accept accepts an incoming connection.
func (l listener) Accept() (session.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	return conn.(session.Conn), err
}

// Disconnect disconnects the connection with the given reason.
func (l listener) Disconnect(conn session.Conn, reason string) error {
	return l.Listener.Disconnect(conn.(*minecraft.Conn), reason)
}
