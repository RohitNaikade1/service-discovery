package server

import (
	"crypto/tls"
	"flag"
	"net/http"
	"service-discovery/env"
	"service-discovery/middlewares"
	"service-discovery/routes"

	"github.com/rs/cors"
)

var Logger = middlewares.Logger()

func HttpsServer() {
	port := env.PORT

	addr := flag.String("addr", ":"+port, "HTTPS network address")
	certFile := flag.String("certfile", "server.crt", "certificate file")
	keyFile := flag.String("keyfile", "server.key", "key file")
	flag.Parse()

	router := routes.NewRoutes()

	r := cors.AllowAll().Handler(router.Router)

	//	routes.NewRoutes().Run("localhost:" + port)
	srv := &http.Server{
		Addr:    *addr,
		Handler: r,
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS13,
			PreferServerCipherSuites: true,
		},
	}

	Logger.Info("Starting server on " + *addr)

	err := srv.ListenAndServeTLS(*certFile, *keyFile)
	Logger.Fatal(err.Error())

}
