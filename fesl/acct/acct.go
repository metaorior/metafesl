package acct

import (
	"github.com/metaorior/metafesl/network"
	"github.com/metaorior/metafesl/network/codec"
	"github.com/metaorior/metafesl/storage/database"
)

const (
	acctNuGetAccount      = "NuGetAccount"
	acctNuGetPersonas     = "NuGetPersonas"
	acctNuLogin           = "NuLogin"
	acctNuLoginPersona    = "NuLoginPersona"
	acctNuLookupUserInfo  = "NuLookupUserInfo"
)

type Account struct {
	DB database.Adapter
}

func (acct *Account) answer(client *network.Client, pnum uint32, payload interface{}) {
	client.WriteEncode(&codec.Answer{
		Type:         codec.FeslAccount,
		PacketNumber: pnum,
		Payload:      payload,
	})
}


type answerToken struct {
	Txn            string `fesl:"TXN"`
	Token 		   string `fesl:"telemetryToken"`
	Enabled        string `fesl:"enabled"`
	Filters        string `fesl:"filters"`
	Disabled       bool   `fesl:"disabled"`
}

// GetTelemetryToken handles acct.GetTelemetryToken cmd
func (acct *Account) answerToken(event network.EventClientCommand) {
	acct.answer(
		event.Client,
		event.Command.PayloadID,
		ansGetTelemetryToken{
			Txn:            "GetTelemetryToken",
			Token: "1234",
			Enabled:        "US",
		},
	)
}


type reqGetAccount struct {
	// TXN=GetAccount
	TXN string `fesl:"TXN"`
}

type ansGetAccount struct {
	TXN string `fesl:"TXN"`

	DobDay         int    `fesl:"DOBDay"`
	DobMonth       int    `fesl:"DOBMonth"`
	DobYear        int    `fesl:"DOBYear"`
	Country        string `fesl:"country"`
	Language       string `fesl:"language"`
	GlobalOptIn    bool   `fesl:"globalOptin"`
	ThidPartyOptIn bool   `fesl:"thirdPartyOptin"`
}

type ansClientGetAccount struct {
	ansGetAccount

	NucleusID int    `fesl:"nuid"`
	UserID    int    `fesl:"userId"`
	HeroName  string `fesl:"heroName"`
}

// GetAccount handles acct.GetAccount cmd
func (acct *Account) GetAccount(event network.EventClientCommand) {
	switch event.Client.GetClientType() {
	case clientTypeServer:
		acct.serverGetAccount(event)
	default:
		acct.clientGetAccount(event)
	}
}

func (acct *Account) clientGetAccount(event network.EventClientCommand) {
	acct.answer(
		event.Client,
		event.Command.PayloadID,
		ansClientGetAccount{
			ansGetAccount: ansGetAccount{
				TXN:           "NuGetAccount",
				Country:        "US",
				Language:       "en_US",
				DobDay:         1,
				DobMonth:       1,
				DobYear:        1990,
				GlobalOptIn:    false,
				ThidPartyOptIn: false,
			},

			NucleusID: event.Client.PlayerData.PlayerID,
			UserID:    event.Client.PlayerData.PlayerID,
			HeroName:  event.Client.PlayerData.HeroName, // ?
		},
	)
}

func (acct *Account) serverGetAccount(event network.EventClientCommand) {
	acct.answer(
		event.Client,
		event.Command.PayloadID,
		ansGetAccount{
			TXN:            "NuGetAccount",
			Country:        "US",
			Language:       "en_US",
			DobDay:         1,
			DobMonth:       1,
			DobYear:        1990,
			GlobalOptIn:    false,
			ThidPartyOptIn: false,
		},
	)
}
