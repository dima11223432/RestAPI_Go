package apiserver

import (
	"RestApi/internal/app/store"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

//apiServer...

type apiServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
}

func NewApiServer(config *Config) *apiServer {
	return &apiServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

func (a *apiServer) Start() error {
	if err := a.configureLogger(); err != nil {
		return err
	}

	a.configureRouter()
	if err := a.configureStore(); err != nil {
		return err
	}

	a.logger.Info("Starting api server")
	return http.ListenAndServe(a.config.BinAddr, a.router)
}

func (a *apiServer) configureLogger() error {
	level, err := logrus.ParseLevel(a.config.LogLevel)
	if err != nil {
		return err
	}
	a.logger.SetLevel(level)
	return nil
}

func (a *apiServer) configureRouter() {

	a.router.HandleFunc("/", a.HandleHello())

}

func (a *apiServer) configureStore() error {
	st := store.NewStore(a.config.Store)
	if err := st.Open(); err != nil {
		return err
	}
	return nil
}

func (a *apiServer) HandleHello() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello")
	}

}
