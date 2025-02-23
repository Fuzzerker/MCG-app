package inboundhttp

import (
	"net/http"

	"github.com/swaggest/openapi-go/openapi31"
	"github.com/swaggest/rest/web"
	swgui "github.com/swaggest/swgui/v5emb"
	"go.uber.org/zap"
)

type HttpServer struct {
	authService               AuthService
	userService               UserService
	patientService            PatientService
	attatchmentService        AttatchmentService
	diagnosedConditionService DiagnosedConditionsService
	logger                    *zap.Logger
	webService                *web.Service
}

func NewServer(authService AuthService, userService UserService, patientService PatientService, attatchmentService AttatchmentService, diagnosedConditionService DiagnosedConditionsService, logger *zap.Logger) *HttpServer {
	return &HttpServer{
		authService:               authService,
		userService:               userService,
		patientService:            patientService,
		attatchmentService:        attatchmentService,
		diagnosedConditionService: diagnosedConditionService,
		logger:                    logger,
	}
}

func (server *HttpServer) Start() {
	server.webService = web.NewService(openapi31.NewReflector())

	server.webService.OpenAPISchema().SetTitle("MCG Patient API")
	server.webService.OpenAPISchema().SetDescription("This service provides an API to manage patient data.")
	server.webService.OpenAPISchema().SetVersion("v1.0.0")
	address := "localhost:8080"
	server.setupRoutes()
	server.logger.Info("listening at", zap.String("adress", address))
	server.webService.Docs("/public/docs", swgui.New)
	if err := http.ListenAndServe(address, server.webService); err != nil {
		zap.Error(err)
	}

}
