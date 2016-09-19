package transport

// Handler contains a collection of handlers to wrap specific protocols like Slack or IRC,
// or even provide a dummy "standard out" interface for local testing.
type Handler interface {
	// ReadLine tokenizes a line consumed from the outside world.
	ReadLine() (string, error)

	// WriteMsg writes a message out to the outside world.
	WriteMsg(b []byte) error

	// RoomDescribe marshals the description of the given room in a format appropriate to the
	// underlying transport.
	// TODO RoomDescribe(r *room.Room) ([]byte, error)

	// MobEmoteMsg formats a message for a Mob performing a certain action that is appropriate
	// for the underlying transport.
	// TODO MobEmote(m *mob.Mob, string action) ([]byte, error)
}
