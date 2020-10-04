package main 
// Generate setters and getters for stats
// $ go generate ./cmd/goheroes
//go:generate go run ../stats-codegen/main.go -scan ../../model --getters ../../ranking/getters.go --setters ../../ranking/setters.go --adders ../../ranking/adders.go


// import is like #include or #import from python
import ( 
	"context" //context defines  deadlines, cancelsignals
	"flag" // for parsing flags like --debug    
	
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/subosito/gotenv"


	"github.com/metaorior/metafesl/config"
	"github.com/metaorior/metafesl/backend/fesl"
	"github.com/metaorior/metafesl/backend/network"
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
	
    log.Print("Fesl server started.. ")
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
		log.Print(err)
	}

	network.InitClientData()
	kvs := kvstore.NewInMemory()
	mm := matchmaking.NewPool()

	fesl.New(config.FeslCli(), false, db, mm).ListenAndServe(ctx)
	fesl.New(config.FeslServ(), true, db, mm).ListenAndServe(ctx)

	
	theater.New(config.ThtrCli(), db, kvs, mm).ListenAndServe(ctx)
	theater.New(config.ThtrServ(), db, kvs, mm).ListenAndServe(ctx)
}
