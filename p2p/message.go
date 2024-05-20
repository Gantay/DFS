package p2p

const (
	IncomingMessage = 0x1
	IncomingStream  = 0x2
)

// Message holds any arbitrary data that is sent over the
// each transport between two node ion the network.
type RPC struct {
	From    string
	Payload []byte
	Stream  bool
}
