package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"sheazuzu/common/src/cli"
	verrors "sheazuzu/common/src/errors"
	"sheazuzu/common/src/logging"
	"sheazuzu/common/src/metrics"
	"sheazuzu/common/src/swagger"
	"sheazuzu/common/src/tracing"
	"sheazuzu/sheazuzu/src/configuration"
	"sheazuzu/sheazuzu/src/controller"
	"sheazuzu/sheazuzu/src/database"
	"sheazuzu/sheazuzu/src/generated/sheazuzu"
	"sheazuzu/sheazuzu/src/repository"
	"sheazuzu/sheazuzu/src/service"
	"sync"
	"syscall"
)

var (
	ProjectVersion = "XXXX.XX.X"
)

func main() {

	config := configuration.New()

	verrors.ServiceId = verrors.Sheazuzu

	app := cli.Command{
		Usage:    "Sheazuzu",
		Version:  ProjectVersion,
		Flags:    config.SetupFlags("sheazuzu"),
		Validate: config.Validate,
		Run:      Run(config),
	}

	app.Execute(os.Args[1:]...)
}

func Run(cfg *configuration.Configuration) func(cmd *cli.Command, args ...string) {
	return func(cmd *cli.Command, args ...string) {

		logger := logging.GetLogger(cfg.Logging.Level, cfg.Logging.Format).
			With("version", cmd.Version)

		handleSigterm(logger)

		db := database.InitDB(cfg.Database.GetDatabaseConn())
		db.LogMode(true) //gorm log model
		defer db.Close()

		sheazuzuRepo := repository.ProvideSheazuzuRepository(db, logger)
		sheazuzuSerivce := service.ProvideSheazuzuService(sheazuzuRepo, logger)
		sheazuzuApi := controller.ProvideSheazuzuAPI(sheazuzuSerivce, logger)

		serverWithMiddleware := sheazuzu.NewServerWithMiddleware(sheazuzuApi)
		serverWithMiddleware.GetMatchDataByIdUsingGETMiddlewares = getMiddleWareChain("machineByIdUsingGET", logger)
		serverWithMiddleware.AllMatchDataUsingGETMiddlewares = getMiddleWareChain("allMachinesUsingGET", logger)
		serverWithMiddleware.UploadMatchDataUsingPOSTMiddlewares = getMiddleWareChain("uploadMatchDataUsingPOST", logger)

		contextPath := cfg.Server.GetContextPath()

		router := chi.NewRouter()

		router.Route("/"+contextPath, func(r chi.Router) {
			swaggerDoc, _ := sheazuzu.GetSwagger()
			swagger.RegisterSwaggerHandlers(r, swaggerDoc, contextPath)
			sheazuzu.HandlerFromMux(serverWithMiddleware, r)
		})

		logger.Infof("Starting HTTP service at %v", cfg.Server.Port)
		if contextPath != "" {
			logger.Infof("Context path is set to %s", contextPath)
		}

		// start the internal web server
		go func() {
			err := http.ListenAndServe(fmt.Sprintf(":%v", cfg.Server.Port), router) // Goroutine will block here
			if err != nil {
				logger.Panicf("Could not start webserver: %s", err)
			}
		}()
		// wait for the sigterm signal
		waitForExit()

	}
}

func getMiddleWareChain(endpoint string, logger *zap.SugaredLogger) chi.Middlewares {
	return chi.Chain(
		jsonContentTypeHeaderHandler,
		metrics.GetMetricsRecordingHandlerForEndpoint(endpoint),
		tracing.TraceHandler(logger, endpoint),
	)
}

var jsonContentTypeHeaderHandler = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("attaching header")
		w.Header().Add("Content-Type", "application/json")
		// for local frontend accessing the url with ip address
		w.Header().Add("Access-Control-Allow-Origin", "*")

		next.ServeHTTP(w, r)
	})
}

func waitForExit() {
	wg := sync.WaitGroup{} // Use a WaitGroup to block main() exit
	wg.Add(1)
	wg.Wait()
}

func handleSigterm(logger *zap.SugaredLogger) {

	c := make(chan os.Signal, 1) // Create a channel accepting os.Signal

	// Bind a given os.Signal to the channel we just created
	signal.Notify(c, os.Interrupt)    // Register os.Interrupt
	signal.Notify(c, syscall.SIGTERM) // Register syscall.SIGTERM

	go func() { // Start an anonymous func running in a goroutine
		<-c // that will block until a message is received on
		logger.Infof("Received shutdown signal")
		os.Exit(0) // de-registration and exit program.
	}()
}
