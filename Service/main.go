package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Application struct {
	ConfigDbConfig ConfigDbConfig
	ConfigDb       ConfigDb
	Handlers       Handlers
}

func main() {
	app := BuildApplication()
	router := BuildServer(&app)

	log.Println("Server starting on :8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func BuildApplication() Application {
	configDbConfig := ConfigDbConfig{
		User:     "redis",
		Password: "redis",
		Database: "configs",
	}
	configDb := ConfigDb{
		Config:    configDbConfig,
		Configs:   make(map[string]Config),
		Overrides: make(map[string]ConfigOverrides),
	}
	handlers := Handlers{
		ConfigDb: configDb,
	}

	return Application{
		ConfigDbConfig: configDbConfig,
		ConfigDb:       configDb,
		Handlers:       handlers,
	}
}

func BuildServer(app *Application) http.Handler {
	handlers := app.Handlers
	router := mux.NewRouter()

	// Configs
	router.Methods("GET").
		Path("/configs").
		HandlerFunc(CatchErrors(handlers.ListConfigs))
	router.Methods("POST").
		Path("/configs").
		HandlerFunc(CatchErrors(handlers.PostConfig))
	router.Methods("GET").
		Path("/configs/{service}/{name}").
		HandlerFunc(CatchErrors(handlers.GetConfig))
	router.Methods("DELETE").
		Path("/configs/{service}/{name}").
		HandlerFunc(CatchErrors(handlers.DeleteConfig))

	// Overrides
	router.Methods("GET").
		Path("/configs/{service}/{name}/overrides").
		HandlerFunc(CatchErrors(handlers.ListOverrides))
	router.Methods("POST").
		Path("/configs/{service}/{name}/overrides").
		HandlerFunc(CatchErrors(handlers.PostOverride))
	router.Methods("GET").
		Path("/configs/{service}/{name}/overrides/{entityType}/{entityId}").
		HandlerFunc(CatchErrors(handlers.GetOverride))
	router.Methods("DELETE").
		Path("/configs/{service}/{name}/overrides/{entityType}/{entityId}").
		HandlerFunc(CatchErrors(handlers.DeleteOverride))

	router.Methods("POST").
		Path("/configs/{service}/{name}/value").
		HandlerFunc(CatchErrors(handlers.GetConfigValue))

	var finalHandler http.Handler = router
	finalHandler = loggingMiddleware(finalHandler)
	return finalHandler
}

func CatchErrors(handler func(*http.Request) (*HttpResponse, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := handler(r)
		if panicErr := recover(); panicErr != nil {
			fmt.Fprintf(os.Stderr, "Panic recovered: %+v\n", panicErr)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error handling request: %+v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(res.Status)
		_, err = w.Write(res.Data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing response: %+v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}
