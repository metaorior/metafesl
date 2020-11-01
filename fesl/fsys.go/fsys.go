package fsys

import (
	"github.com/metaorior/metafesl/backend/network"
	"github.com/metaorior/metafesl/backend/network/codec"
	
)

// type ConnectSystem struct {
// 	ServerMode bool

// }

func (fsys *ConnectSystem) answer(client *network.Client, Pktnum uint32, fullPkt interface{} ) {
	client.WriteEncode(&codec.Answer{
		Type: codec.FeslSystem,
		PacketNumber: pnum,
		Payload: payload,
	})
}


type ansGetPing struct {

	TXN string `fesl:"TXN"`

	MinPings int `fesl:"minPingSitesToPing"`

	// PingSites defines a list of endpoints, which should be pinged,
	// accordiningly to minPingSitesToPing setting.
	PingSites string `fesl:"pingSites"` // doenst matter making it 127.0.0.1 //syn


}

// GetPingSites command called by the clinet

func GetPingSites(event network.EventClientCommand) {
	fsys.answer(
		event.Client,
		event.Command.PayloadID,
		ansGetPing{
			TXN: "GetPingSites",
			MingPing: 0,
			PingSites: "127.0.0.1"
		}
	)
}
