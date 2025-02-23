package main

import (
	inboundhttp "mcg-app-backend/io/inbound/http"
	inmemory "mcg-app-backend/io/outbound/in-memory"
	"mcg-app-backend/service/attatchments"
	"mcg-app-backend/service/auth"
	diagnosedconditions "mcg-app-backend/service/diagnosedConditions"
	"mcg-app-backend/service/patients"
	"mcg-app-backend/service/tracing"
	"mcg-app-backend/service/users"
	"time"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	repo := inmemory.NewInMemoryRepo()
	tracer := tracing.NewService(logger)
	patientSrv := patients.NewPatientService(repo, tracer)
	attatchmentSrv := attatchments.NewAttachmentService(repo, patientSrv, tracer)
	diagnosedConditionSrv := diagnosedconditions.NewDiagnosedConditionService(repo, patientSrv, tracer)
	userService := users.NewService(repo, tracer)
	//in a real application, these would be fed by environment variables
	expirationTime := time.Minute * 10
	issuer := "localhost"
	tokenSecret := "asdf"
	authService := auth.NewService(userService, tracer, expirationTime, issuer, tokenSecret)
	inboundhttp.NewServer(authService, userService, patientSrv, attatchmentSrv, diagnosedConditionSrv, logger).Start()
}
