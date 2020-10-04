package acct

import (
	"github.com/sirupsen/logrus"
	"math/rand"
	"github.com/metaorior/metafesl/network"
)

type reqNuLogin struct {
	// TXN=NuLogin
	TXN string `fesl:"TXN"`
	// returnEncryptedInfo=0
	ReturnEncryptedInfo int `fesl:"returnEncryptedInfo"`
	// macAddr=$0a0027000000
	MacAddr string `fesl:"macAddr"`
}

type reqNuLoginServer struct {
	reqNuLogin

	AccountName     string `fesl:"nuid"`     // Value specified in +eaAccountName
	AccountPassword string `fesl:"password"` // Value specified in +eaAccountPassword
}

type reqNuLoginClient struct {
	reqNuLogin

	EncryptedInfo string `fesl:"encryptedInfo"` // Value specified in +sessionId
}

type ansNuLogin struct {
	Txn       string `fesl:"TXN"`
	ProfileID int    `fesl:"profileId"`
	UserID    int    `fesl:"userId"`
	NucleusID int    `fesl:"nuid"`
	LobbyKey  string `fesl:"lkey"`
}

type ansNuLoginErr struct {
	Txn     string                `fesl:"TXN"`
	Message string                `fesl:"localizedMessage"`
	Errors  []nuLoginContainerErr `fesl:"errorContainer"`
	Code    int                   `fesl:"errorCode"`
}

type nuLoginContainerErr struct {
	Value      string `fesl:"value"`
	FieldError string `fesl:"fieldError"`
	FieldName  string `fesl:"fieldName"`
}


var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}


// NuLogin handles acct.NuLogin cmd
func (acct *Account) NuLogin(event network.EventClientCommand) {
	lkey  := randString(10) //random string.. doesnt matter w/ security
	//+ only sessionId is important 4 security
	
	
	event.Client.PlayerData.LobbyKey = lkey
	err := network.Lobby.Add(lkey, event.Client.PlayerData)
	if err != nil {
		logrus.WithError(err).Warn("Cannot add ClientData in acct.NuLogin")
		return
	}

	switch event.Client.GetClientType() {
	case clientTypeServer:
		acct.serverNuLogin(event)
	default:
		acct.clientNuLogin(event)
	}
}

func (acct *Account) clientNuLogin(event network.EventClientCommand) {
	player, err := acct.DB.GetPlayerByToken(
		acct.DB.NewSession(),
		event.Command.Message["encryptedInfo"],
	)
	if err != nil {
		logrus.WithError(err).Warn("Client cannot sign in the acct.NuLogin")
		acct.clientNuLoginNotAuthorized(&event)
		return
	}

	event.Client.PlayerData.PlayerID = player.ID

	acct.answer(
		event.Client,
		event.Command.PayloadID,
		ansNuLogin{
			Txn:       acctNuLogin,
			UserID:    event.Client.PlayerData.PlayerID,
			ProfileID: event.Client.PlayerData.PlayerID,
			NucleusID: event.Client.PlayerData.PlayerID,
			LobbyKey:  event.Client.PlayerData.LobbyKey,
		},
	)
}

func (acct *Account) clientNuLoginNotAuthorized(event *network.EventClientCommand) {
	acct.answer(
		event.Client,
		event.Command.PayloadID,
		ansNuLoginErr{
			Txn:     acctNuLogin,
			Message: `"The user is not entitled to access this game"`,
			Code:    120,
		},
	)
}

// acctNuLoginServer - login cmd for servers
func (acct *Account) serverNuLogin(event network.EventClientCommand) {
	srv, err := acct.DB.GetServerLogin(
		acct.DB.NewSession(),
		event.Command.Message["nuid"],
	)
	if err != nil {
		logrus.WithError(err).Warn("Server cannot sign in the acct.NuLogin")
		acct.serverNuLoginNotAuthorized(&event)
		return
	}

	// TODO: Raw passwords are really, really insecure
	// 1. Validate credentials using some bcrypt or argon2
	// 2. DO NOT STORE raw passwords
	if srv.AccountPassword != event.Command.Message["password"] {
		acct.serverNuLoginNotAuthorized(&event)
		return
	}

	event.Client.PlayerData.ServerID = srv.ID
	event.Client.PlayerData.ServerSoldierName = srv.SoldierName
	event.Client.PlayerData.ServerUserName = srv.AccountUsername

	acct.answer(
		event.Client,
		event.Command.PayloadID,
		ansNuLogin{
			Txn:       acctNuLogin,
			ProfileID: srv.ID,
			UserID:    srv.ID,
			NucleusID: srv.ID,
			LobbyKey:  event.Client.PlayerData.LobbyKey,
		},
	)
}

func (acct *Account) serverNuLoginNotAuthorized(event *network.EventClientCommand) {
	acct.answer(
		event.Client,
		event.Command.PayloadID,
		ansNuLoginErr{
			Txn:     acctNuLogin,
			Message: `"The password the user specified is incorrect"`,
			Code:    122,
		},
	)
}
