package api

import (
	"fmt"
	"journeymaster/models"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	apiPortFallback = "8080"
)

// JourneyMasterAPI struct
type JourneyMasterAPI struct {
	apiPort string

	// (Gin) API Server
	server *http.Server

	// Basic Auth Map
	apiAuth map[string]string
}

// NewJourneyMasterAPI creates a new API instance
func NewJourneyMasterAPI() (*JourneyMasterAPI, error) {
	api := new(JourneyMasterAPI)
	return api, nil
}

// StartServer starts the HTTP server and listen
func (api *JourneyMasterAPI) StartServer() error {
	// set the listener port
	api.setListenerPort()
	if err := api.setCredentials(); err != nil {
		return err
	}

	// setup routing
	router := api.getRouter()
	api.setRoutes(router)

	// set Server variable
	api.server = &http.Server{
		Addr:    api.apiPort,
		Handler: router,
	}

	// TLS secured server
	if os.Getenv("TLS_ENABLED") == "true" {
		err := api.server.ListenAndServeTLS(os.Getenv("TLS_PEM"), os.Getenv("TLS_KEY"))
		if err != nil {
			log.Printf("JourneyMaster: server not available. Reason %v", err)
		}

	} else {
		err := api.server.ListenAndServe()
		if err != nil {
			log.Printf("JourneyMaster: server not available. Reason %v", err)
		}
	}

	return nil
}

// setAPIListenerPort Sets the JourneyMaster API port (fallbacks to 8080)
func (api *JourneyMasterAPI) setListenerPort() {
	value := ":" + os.Getenv("API_PORT")
	if len(value) == 0 {
		api.apiPort = ":" + apiPortFallback
		return
	}
	api.apiPort = value

	log.Printf("JourneyMaster: Setting up Port to %s", api.apiPort)
}

// setAPICredentials sets the credentials loaded in API_CREDENTIALS and prepare them for ginRouter
func (api *JourneyMasterAPI) setCredentials() error {
	api.apiAuth = make(map[string]string)

	authData := os.Getenv("API_CREDENTIALS")
	if authData == "" {
		log.Printf("Error: API_CREDENTIALS environment variable is not set")
		return fmt.Errorf("API_CREDENTIALS environment variable is not set")
	}

	lines := strings.Split(authData, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			log.Printf("Invalid format in credentials: %s\n", line)
			continue
		}
		api.apiAuth[parts[0]] = parts[1]
	}

	return nil
}

// getRouter returns a new gin router
func (api JourneyMasterAPI) getRouter() *gin.Engine {
	gin.SetMode(os.Getenv(`GIN_MODE`))
	router := gin.Default()

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	// deny all unknown routes
	router.NoRoute(func(c *gin.Context) {
		api.SendError(c, "VTA0001", "Method not implemented", http.StatusNotImplemented)
	})

	return router
}

// setRoutes sets all routes with its handler functions
func (api *JourneyMasterAPI) setRoutes(router *gin.Engine) {
	// set authorized group /v1
	authorized := router.Group(os.Getenv("API_PREFIX"), gin.BasicAuth(api.apiAuth))

	// analyze a vehicle trip
	authorized.POST("/trip", api.analyzeTripHandler)
}

// SendResponse returned by the api
func (api *JourneyMasterAPI) SendResponse(c *gin.Context, payload interface{}) {
	// build api json results with error and payload (not needed here)
	/*var response = models.APIResponse{Payload: payload}
	api.send(c, http.StatusOK, response)*/
	api.send(c, http.StatusOK, payload)
}

// SendError send an error
func (api *JourneyMasterAPI) SendError(c *gin.Context, errorCode string, errorMessage string, httpStatus int) {
	apiContext := c.GetString(`apiContext`)
	log.Printf("JourneyMaster: SendError(%s) Error code '%s' with message '%s' and http status code %d set", apiContext, errorCode, errorMessage, httpStatus)

	// build errors json result
	response := models.APIErrors{}
	response.Errors = append(response.Errors, models.APIErrorElement{ErrorCode: errorCode, ErrorMessage: errorMessage})

	api.send(c, httpStatus, response)
}

// send data to the client
func (api *JourneyMasterAPI) send(c *gin.Context, httpStatus int, response interface{}) {
	// use c.PureJSON for literally encoding
	c.JSON(httpStatus, response)
}
