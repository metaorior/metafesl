package main 
// Generate setters and getters for stats
// $ go generate ./cmd/goheroes
//go:generate go run ../stats-codegen/main.go -scan ../../backend/model --getters ../../backend/ranking/getters.go --setters ../../backend/ranking/setters.go --adders ../../backend/ranking/adders.go


// import is like #include or #import from python
import ( 
	"context" //context defines  deadlines, cancelsignals
	"flag" // for parsing flags like --debug 

	"github.com/google/gops/agent"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/subosito/gotenv"


	"github.com/metaorior/metafesl/backend/config"
	"github.com/metaorior/metafesl/backend/network"
	"github.com/metaorior/metafesl/backend/fesl"
	"github.com/metaorior/metafesl/backend/matchmaking"
	"github.com/metaorior/metafesl/backend/storage/database"
	"github.com/metaorior/metafesl/backend/storage/kvstore"
	"github.com/metaorior/metafesl/backend/theater"

)

func main() {
	setupConfig()
	setupzeroLog()

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()
	startFeslServer(ctx)
	
	zerolog.Print("Fesl server started.. ")
	<-ctx.Done()
}


//setup zeroLog so we can call it later
func setupzeroLog() {
	//logs will write with UNIX time
	//TODO browse its vendor package
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

}

func setupConfig() {
	var (
		configFile string
	)
	flag.StringVar(&configFile, "config", ".env", "Path to configuration file")
	flag.Parse()

	gotenv.Load(configFile)
	config.Initialize()
}


func startFeslServer(ctx context.Context) {
	db, err := database.New()
	if err != nil {
		logrus.Fatal(err)
	}

	network.InitClientData()
	kvs := kvstore.NewInMemory()
	mm := matchmaking.NewPool()

	fesl.New(config.feslCli(), false, db, mm).ListenAndServe(ctx)
	fesl.New(config.feslServ(), true, db, mm).ListenAndServe(ctx)

	
	theater.New(config.thtrCli(), db, kvs, mm).ListenAndServe(ctx)
	theater.New(config.thtrServ(), db, kvs, mm).ListenAndServe(ctx)
}
