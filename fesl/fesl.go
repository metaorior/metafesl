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

// Fesl handles incoming and outgoing TCP data (called as FESL)
type Fesl struct {
	socket *network.socket //socket is the socket we are going to use
	cmds network.CmdRegistry //the commands
   //more
}

//new Create new FESL server

func New(bind string, server bool, db database.Adapter, mm *matchmaking.Pool) *Fesl {
	socket, err := network.NewSocketTLS(bind)
	if err != nil {
		logrus.Fatal(err)
		return nil
	}
	const (
		nOfCmds = 14 //make the commands fixed / a tecnique faster than case select

	)

	r := make(network.CmdRegistry, nOfCmds)
	{
		acct := &acct.Account{DB: db} //todo rename
		fsys := &fsys.ConnectSystem{ServerMode: server} //todo rename
		gsum := &gsum.GameSummary{} //todo rename
		pnow := &pnow.PlayNow{MM: mm} //todo rename
		rank := &rank.Ranking{DB: db} //todo rename
	}
	
	r.Register("newClient connected", func(ev network.EventClientCommand){
		//TLS Only
		fsys.MemCheck(ev.Client)

		ev.Client.HeartTicker = time.NewTicker(55 * time.Second)
		go func() {
			for ev.Client.Active {
				select {
				case <-ev.Client.HeartTicker.C:
					if !ev.Client.Active {
						return
					}
					fsys.MemCheck(ev.Client)
					logrus.Println("==Client Keep Alive==")
				}
			}
		}()
	} ) // close parenthesis lol


	// for r.Register {
	// 	r.Register = split("client.cmd")
	// }

	reg := r.Register
	reg("client.cmd.Hello", func(ev network.EventClientCommand) {
		if !server {
			gsum.Getgame_token(ev)
		}
		fsys.Hello(ev)
	})

	
	reg("client.cmd.Login", acct.Login)
	reg("client.cmd.AllHeroes", acct.AllHeroes)
	reg("client.cmd.GetAccount", acct.GetAccount)
	reg("client.cmd.LoginHero", acct.LoginHero)
	reg("client.cmd.AllStats", rank.AllStats
	reg("client.cmd.SingleStats", rank.SingleStats)
	reg("client.cmd.UserInfo", acct.UserInfo)
	reg("client.cmd.Sites", fsys.Sites)
	reg("client.cmd.UpdateStats", rank.UpdadeStats)
	reg("client.cmd.token", acct.token)
	reg("client.cmd.Pnow", func(ev network.EventClientCommand){
		pnow.Start(ev)
		pnow.Status(ev)
	})
	r.Register("client.cmd.MemCheck", func(ev network.EventClientCommand) {
	
	})
	return &Fesl{socket, r}

}

//listen and serve

func(fm *Fesl) ListenAndServe(ctx context.Context) {
	go fm.run(ctx)
}

// Run starts listening the socket for events and handles them upon receiving
// a message
func (fm *Fesl) Run(ctx context.Context) {
	for {
		select {
		case event := <-fm.socket.EventChan:
			fm.Handle(event)
		case <-ctx.Done():
			return
		}
	}
}


// Handle takes care of handling a single event
func (fm *Fesl) Handle(event network.SocketEvent) {
	ev, ok := event.Data.(network.EventClientCommand)
	if !ok {
		logrus.Println("Logic error: Cannot cast event to network.EventClientCommand")
		return
	}

	if !ev.Client.Active {
		logrus.WithField("cmd", ev.Command).Warn("Inactive client")
		return
	}

	fn, ok := fm.cmds.Find(event.Name)
	if !ok {
		logrus.
			WithFields(logrus.Fields{
				"event":   event.Name,
				"payload": ev.Command.Message,
				"query":   ev.Command.Query,
			}).
			Warn("fesl.UnhandledRequest")
		return
	}

	fn(ev)
}
