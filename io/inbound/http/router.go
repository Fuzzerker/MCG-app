package inboundhttp

import (
	"github.com/riandyrn/otelchi"
)

func (server HttpServer) setupRoutes() {

	server.webService.Use(otelchi.Middleware("mcg-application"))
	server.webService.Use(server.RequireValidToken)
	server.webService.Post("/public/users", server.handlePostUser())
	server.webService.Post("/public/users/login", server.handleLogin())
	server.webService.Post("/patients", server.handlePostPatient())
	server.webService.Put("/patients/{id}", server.handlePutPatient())
	server.webService.Post("/patients/{patientId}/attatchments", server.handlePostPatientAttatchment())
	server.webService.Post("/patients/{patientId}/diagnosedConditions", server.handlePostDiagnosedCondition())

	server.webService.Get("/patients", server.handleGetPatients())
	server.webService.Delete("/diagnosedConditions/{id}", server.handleDeleteDiagnosedCondition())

	server.webService.Delete("/patients/{id}", server.handleDeletePatient())
	server.webService.Delete("/attatchments/{id}", server.handleDeleteAttatchment())

}
