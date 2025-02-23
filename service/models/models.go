package models

import (
	"mime/multipart"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserRequest struct {
	Username string `json:"username" required:"true" minLength:"6" description:"Username of the user"`
	Password string `json:"password" required:"true" minLength:"6" description:"Password the user will use to log in"`
}

type LoginResponse struct {
	Token string `json:"token" descripiton:"access token generated for the given credentials.  Should be sent as a bearer token on all future requests"`
}

type Empty struct {
}

type User struct {
	Username string
	Password string
}

type UserClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type CreateAttatchmentRequest struct {
	Name        string         `formData:"name" description:"name of this attatchment" required:"true" minLength:"5"`
	Description string         `formData:"descripiton" description:"description of this attatchment"`
	Type        string         `formData:"type" description:"type of attatchment" minLength:"1"`
	Data        multipart.File `formData:"data" description:"data associated with this attatchment"`
	PatientId   int            `path:"patientId"`
}

type UpdatePatientRequest struct {
	PatientRequest
	Id int `path:"id"`
}

type DeleteByIdRequest struct {
	Id int `path:"id"`
}

type PatientRequest struct {
	Name               string    `json:"name" required:"true" minLength:"3" description:"name of the patient"`
	Address            string    `json:"address"  description:"address of the patient"`
	PhoneNumber        string    `json:"phoneNumber" required:"true" minLength:"10" description:"phone number of the patient"`
	ExternalIdentifier string    `json:"externalIdentifier" required:"true" minLength:"3" description:"External identifier of the patient (for example - social security number)"`
	DateOfBirth        time.Time `json:"dateOfBirth" description:"date of birth of patient" required:"true"`
}

type Patient struct {
	Name                string               `json:"name" required:"true" minLength:"3" description:"name of the patient"`
	Address             string               `json:"address"  description:"address of the patient"`
	PhoneNumber         string               `json:"phoneNumber" required:"true" minLength:"10" description:"phone number of the patient"`
	ExternalIdentifier  string               `json:"externalIdentifier" required:"true" minLength:"3" description:"External identifier of the patient (for example - social security number)"`
	DateOfBirth         time.Time            `json:"dateOfBirth" description:"date of birth of patient" required:"true"`
	Id                  int                  `json:"id" description:"Internal id of the patient"`
	DiagnosedConditions []DiagnosedCondition `json:"diagnosedConditions" description:"conditions with which the patient has been diagnosed"`
	Attatchments        []Attatchment        `json:"attatchments" description:"attatchments for theph patient.  Could be any form of medical imaging or doctor's reports"`
}

type CreateDiagnosedConditionRequest struct {
	DiagnosedCondition
	PatientId int `path:"patientId"`
}

type DiagnosedCondition struct {
	Id          int       `json:"id" description:"internal id of the diagnosed condition"`
	PatientId   int       `json:"patientId" description:"internal id of the patient for whom this condition was diagnosed"`
	Name        string    `json:"name" description:"name of the condition" required:"true" minLength:"3"`
	Code        string    `json:"code" description:"medical code to identify the condition" required:"true" minLength:"3"`
	Description string    `json:"description" description:"description of the condition"`
	Date        time.Time `json:"date" description:"date on which this condition was diagnosed" required:"true"`
}

type Attatchment struct {
	Id          int    `json:"id" description:"id of the attatchment"`
	PatientId   int    `json:"patientId" description:"id of the patient to whom this attatchment belongs"`
	Name        string `json:"name" description:"name of this attatchment" required:"true" minLength:"5"`
	Description string `json:"descripiton" description:"description of this attatchment"`
	Type        string `json:"type" description:"type of attatchment" minLength:"4"`
	Data        []byte `json:"data" description:"data associated with this attatchment"`
}

type PatientSearch struct {
	Name                   string `query:"name" description:"name to search for"`
	Address                string `query:"address" description:"address to search for"`
	Phone                  string `query:"phone" description:"phone to search for"`
	ExternalIdentifier     string `query:"externalIdentifier" description:"externalIdentifier to search for"`
	DiagnosedConditionName string `query:"diagnosedConditionName" description:"name of the medical condition to search for"`
	DiagnosedConditionCode string `query:"diagnosedConditionCode" description:"code of the medical condition to search for"`
	AttatchmentName        string `query:"attatchmentName" description:"attatchment name to search for"`
	AttatchmentType        string `query:"attatchmentType" description:"attatchment type to search for"`
}
