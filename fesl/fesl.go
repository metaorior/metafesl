package fesl

import (
	"context"
	"time"
	"github.com/metaorior/metafesl/fesl"
	"github.com/metaorior/metafesl/fesl/login"
	"github.com/metaorior/metafesl/fesl/system"
	"github.com/metaorior/metafesl/fesl/playnow"
	"github.com/metaorior/metafesl/fesl/queryStats"
	"github.com/metaorior/metafesl/fesl/updateStats"
	"github.com/metaorior/metafesl/fesl/network"
	"github.com/metaorior/metafesl/fesl/storage/database"

	log "github.com/rs/zerolog"
)

type Fesl struct {
	socket *network.socket
	cmds network.CmdRegistry
}

// New create new Fesl
func NewFesl(bind string, serverMode bool, db database.Adapter, mm *playnow.Pool) *Fesl {
	socket, err := network.NewSocketTLS(bind) //bind is the variable (
		if err != nil {
			log.Println("error binding socket in NewFESl LN 27", err)
			return nil
		}
	)

	const (
		nOfCmds = 14 // cmds = commands
	)

	r := make(network.CmdRegistry, nOfCmds)