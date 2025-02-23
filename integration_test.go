package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mcg-app-backend/service/models"
	"mime/multipart"
	"net/http"
	"net/url"
	"testing"
	"time"
	"unicode"
)

var authToken = "b"
var patientId = -1
var conditionId = -1

func TestApplication(t *testing.T) {
	go main()

	time.Sleep(time.Second)
	var results TestResults
	results = testUserCreate(results)
	results = testUserLogin(results)
	results = testPatientCreate(results)
	results = testPatientUpdate(results)
	results = testAddAttatchmentToPatient(results)
	results = testAddDiagnosedConditionToPatient(results)
	results = testSearchPatients(results)
	results = testDeleteCondition(results)
	results = testDeleteAttatchment(results)
	results = testDeletePatient(results)

	for _, res := range results {

		if res.Error != nil {
			t.Errorf("error in %v with error %v", res.Name, res.Error.Error())
		} else {
			fmt.Println(res.Name, "OK!")
		}
	}
}

type TestResult struct {
	Name  string
	Error error
}

type TestResults []TestResult

func (t *TestResults) Add(name string, err error) {
	*t = append(*t, TestResult{
		Name:  name,
		Error: err,
	})
}

func testDeletePatient(results TestResults) TestResults {
	path := fmt.Sprintf("/patients/%v", patientId)
	searchPath := "/patients"
	results.Add("delete patient", deleteAndEnsureStatus(path, 204, nil))
	var patients []models.Patient
	results.Add("test search patients by name after deletion", getAndEnsureStatus(searchPath, models.PatientSearch{
		Name: "Jane Smith",
	}, 200, &patients))
	results.Add("test that patients are not returned after deletion", func() error {
		if len(patients) != 0 {
			return fmt.Errorf("expected no patients, but got %v", len(patients))
		}
		return nil
	}())

	return results
}

func testDeleteAttatchment(results TestResults) TestResults {
	path := fmt.Sprintf("/attatchments/%v", conditionId)
	searchPath := "/patients"
	results.Add("delete diagnosed condition", deleteAndEnsureStatus(path, 204, nil))
	var patients []models.Patient
	results.Add("test search patients by diagnosed condition after deletion", getAndEnsureStatus(searchPath, models.PatientSearch{
		AttatchmentType: "MRI",
	}, 200, &patients))
	results.Add("test that returned patient has only 1 attatchment after delete", func() error {
		if len(patients) != 1 {
			return fmt.Errorf("expected 1 patient, but got %v", len(patients))
		}
		if len(patients[0].Attatchments) != 1 {
			return fmt.Errorf("expected 1 attatchment for patient, but got %v", len(patients[0].Attatchments))
		}

		return nil
	}())

	return results
}

func testDeleteCondition(results TestResults) TestResults {
	path := fmt.Sprintf("/diagnosedConditions/%v", conditionId)
	searchPath := "/patients"
	results.Add("delete diagnosed condition", deleteAndEnsureStatus(path, 204, nil))
	var patients []models.Patient
	results.Add("test search patients by diagnosed condition after deletion", getAndEnsureStatus(searchPath, models.PatientSearch{
		Name:                   "mismatch",
		DiagnosedConditionName: "some condition",
	}, 200, &patients))
	results.Add("test that patients are not returned for deleted condition", func() error {
		if len(patients) != 0 {
			return fmt.Errorf("expected no patients, but got %v", len(patients))
		}
		return nil
	}())

	return results
}

func testSearchPatients(results TestResults) TestResults {
	path := "/patients"
	realAuth := authToken
	authToken = "INVALID"
	results.Add("test search patients without auth token", getAndEnsureStatus(path, nil, 401, nil))
	authToken = realAuth

	var patients []models.Patient
	results.Add("test search patients with no conditions", getAndEnsureStatus(path, models.PatientSearch{}, 200, &patients))
	results.Add("test no patients are returned if no conditions are given", func() error {
		if len(patients) > 0 {
			return fmt.Errorf("expected no patients, but found %v", len(patients))
		}
		return nil
	}())

	results.Add("test search patients with no matching criteria", getAndEnsureStatus(path, models.PatientSearch{
		Name:                   "Invalid name",
		DiagnosedConditionName: "invalid name",
	}, 200, &patients))
	results.Add("test no patients are returned if no matching criteria", func() error {
		if len(patients) > 0 {
			return fmt.Errorf("expected no patients, but found %v", len(patients))
		}
		return nil
	}())

	results.Add("test search patients by name", getAndEnsureStatus(path, models.PatientSearch{
		Name: "Jane Smith",
	}, 200, &patients))

	results.Add("test that correct patients are returned by name", func() error {
		if len(patients) != 1 {
			return fmt.Errorf("expected 1 patient, but found %v", len(patients))
		}
		if len(patients[0].Attatchments) != 2 {
			return fmt.Errorf("expected patient to have 2 attatchments, but had %v", len(patients[0].Attatchments))
		}
		if len(patients[0].DiagnosedConditions) != 1 {
			return fmt.Errorf("expected patient to have 1 diagnosed condition, but had %v", len(patients[0].DiagnosedConditions))
		}
		return nil
	}())

	patients = nil

	results.Add("test search patients by diagnosed condition with mismatching name", getAndEnsureStatus(path, models.PatientSearch{
		Name:                   "mismatch",
		DiagnosedConditionName: "some condition",
	}, 200, &patients))

	results.Add("test that correct patients are returned by diagnosed condition", func() error {
		if len(patients) != 1 {
			return fmt.Errorf("expected 1 patient, but found %v", len(patients))
		}
		if len(patients[0].Attatchments) != 2 {
			return fmt.Errorf("expected patient to have 2 attatchments, but had %v", len(patients[0].Attatchments))
		}
		if len(patients[0].DiagnosedConditions) != 1 {
			return fmt.Errorf("expected patient to have 1 diagnosed condition, but had %v", len(patients[0].DiagnosedConditions))
		}
		return nil
	}())

	results.Add("test search patients by attatchment type ", getAndEnsureStatus(path, models.PatientSearch{
		AttatchmentType: "MRI",
	}, 200, &patients))

	results.Add("test that correct patients are returned by attatchment type", func() error {
		if len(patients) != 1 {
			return fmt.Errorf("expected 1 patient, but found %v", len(patients))
		}
		if len(patients[0].Attatchments) != 2 {
			return fmt.Errorf("expected patient to have 2 attatchments, but had %v", len(patients[0].Attatchments))
		}
		if len(patients[0].DiagnosedConditions) != 1 {
			return fmt.Errorf("expected patient to have 1 diagnosed condition, but had %v", len(patients[0].DiagnosedConditions))
		}
		return nil
	}())

	return results
}

func testAddDiagnosedConditionToPatient(results TestResults) TestResults {
	path := fmt.Sprintf("/patients/%v/diagnosedConditions", patientId)

	realAuth := authToken
	authToken = "INVALID"
	results.Add("test add condition without auth token", putAndEnsureStatus(path, nil, 401, nil))
	authToken = realAuth

	results.Add("test add condition with missing fields", postAndEnsureStatus(path, models.DiagnosedCondition{
		Name: "some condition",
	}, 400, nil))
	realPatientId := patientId
	patientId = -1
	results.Add("test add condition with invalid patient id", postAndEnsureStatus(path, models.DiagnosedCondition{
		Name: "some condition",
	}, 400, nil))
	patientId = realPatientId
	var diagnosedCondition models.DiagnosedCondition
	results.Add("test add condition with all fields", postAndEnsureStatus(path, models.DiagnosedCondition{
		Name:        "some condition",
		Code:        "ABCD",
		Description: "something",
		Date:        time.Now(),
	}, 200, &diagnosedCondition))
	conditionId = diagnosedCondition.Id
	return results

}

func testAddAttatchmentToPatient(results TestResults) TestResults {
	attatchment := models.Attatchment{
		Name:        "some-attatchment",
		Description: "some data",
		Type:        "MRI",
		Data:        []byte("some arbitrary data"),
	}
	var returnedAttatchment models.Attatchment
	realAuth := authToken
	authToken = "INVALID"
	results.Add("test add attatchment to patient without auth", postAttatchment(attatchment, 401, nil))

	authToken = realAuth

	results.Add("test add attatchment to patient", postAttatchment(attatchment, 200, nil))
	results.Add("test add second attatchment to patient", postAttatchment(attatchment, 200, &returnedAttatchment))
	results.Add("test ensure that added attatchment has id", func() error {
		if returnedAttatchment.Id < 1 {
			return fmt.Errorf("expected id > 0 but got %v", returnedAttatchment.Id)
		}
		return nil
	}())
	results.Add("test ensure that added attatchment has data", func() error {
		if len(returnedAttatchment.Data) == 0 {
			return fmt.Errorf("expected attatchment data, but got none")
		}
		return nil
	}())

	realPatientId := patientId
	patientId = -1
	results.Add("test add attatchment to patient with invalid patientId", postAttatchment(attatchment, 400, nil))
	patientId = realPatientId
	return results
}

func postAttatchment(attatchment models.Attatchment, status int, respBodyPntr any) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writeFormField(writer, "name", attatchment.Name)
	writeFormField(writer, "description", attatchment.Description)
	writeFormField(writer, "type", attatchment.Type)
	dataReader := bytes.NewReader(attatchment.Data)
	part, _ := writer.CreateFormFile("data", attatchment.Name)
	io.Copy(part, dataReader)
	writer.Close()
	url := fmt.Sprintf("http://localhost:8080/patients/%v/attatchments", patientId)
	r, _ := http.NewRequest(http.MethodPost, url, body)
	r.Header.Add("Content-Type", writer.FormDataContentType())
	r.Header.Add("Authorization", "Bearer "+authToken)
	return doAndEnsureStatus(r, status, respBodyPntr)
}

func writeFormField(writer *multipart.Writer, name string, value string) error {
	descWriter, err := writer.CreateFormField(name)
	if err != nil {
		return fmt.Errorf("error creating form field %v %w", name, err)
	}
	_, err = descWriter.Write([]byte(value))
	if err != nil {
		return fmt.Errorf("error writing form field %v %w", name, err)
	}
	return nil
}

func testPatientUpdate(results TestResults) TestResults {
	path := fmt.Sprintf("/patients/%v", patientId)
	realAuth := authToken
	authToken = "INVALID"
	results.Add("test update patient without auth token", putAndEnsureStatus(path, models.Patient{
		Name: "john smith",
	}, 401, nil))
	authToken = realAuth
	results.Add("test update patient with missing fields", putAndEnsureStatus(path, models.Patient{
		Name: "john smith",
	}, 400, nil))
	patient := models.PatientRequest{
		Name:               "Jane Smith",
		Address:            "185 main street",
		PhoneNumber:        "8044955579",
		ExternalIdentifier: "123",
		DateOfBirth:        time.Now(),
	}

	results.Add("test update patient with valid fields", putAndEnsureStatus(path, patient, 200, &patient))
	results.Add("test ensure updated patient data is returned", func() error {
		if patient.Address == "185 main street" {
			return nil
		}
		return fmt.Errorf("expected address to equal 185 main street but was %v", patient.Address)
	}())

	path = fmt.Sprintf("/patients/%v", -1)
	results.Add("test update patient with invalid Id", putAndEnsureStatus(path, models.Patient{
		Name: "john smith",
	}, 400, nil))

	return results
}

func testPatientCreate(results TestResults) TestResults {
	path := "/patients"
	realAuth := authToken
	authToken = "INVALID"
	results.Add("test create patient without auth token", postAndEnsureStatus(path, models.Patient{
		Name: "john smith",
	}, 401, nil))
	authToken = realAuth
	results.Add("test create patient with missing fields", postAndEnsureStatus(path, models.Patient{
		Name: "john smith",
	}, 400, nil))
	patient := models.Patient{
		Name:               "John Smith",
		Address:            "132 main street",
		PhoneNumber:        "8044955579",
		ExternalIdentifier: "123",
		DateOfBirth:        time.Now(),
	}
	results.Add("test create patient with all fields", postAndEnsureStatus(path, patient, 200, nil))
	results.Add("test create patient with duplicate external id", postAndEnsureStatus(path, patient, 409, nil))

	patient.ExternalIdentifier = "abc"
	patient.Name = "Jane Smith"
	results.Add("test create second patient", postAndEnsureStatus(path, patient, 200, &patient))
	fmt.Println("patientID", patient.Id)
	results.Add("test ensure id on new patient", func() error {
		if patient.Id > 0 {
			return nil
		}
		return fmt.Errorf("invalid Id %v on patient, expected > 0", patient.Id)
	}())
	patientId = patient.Id

	return results
}

func testUserLogin(results TestResults) TestResults {
	path := "/public/users/login"
	results.Add("test login invalid username", postAndEnsureStatus(path, models.UserRequest{
		Username: "aaaaaaaa",
	}, 400, nil))
	results.Add("test login invalid password", postAndEnsureStatus(path, models.UserRequest{
		Username: "abcdefg",
		Password: "aaaaaaa",
	}, 400, nil))
	var loginResponse models.LoginResponse
	results.Add("test login valid password", postAndEnsureStatus(path, models.UserRequest{
		Username: "abcdefg",
		Password: "abcdefg",
	}, 200, &loginResponse))
	authToken = loginResponse.Token
	return results
}

func testUserCreate(results TestResults) TestResults {
	path := "/public/users"
	results.Add("test post user with incomplete body", postAndEnsureStatus(path, models.UserRequest{
		Username: "abcdefg",
	}, 400, nil))

	results.Add("test post user with success body", postAndEnsureStatus(path, models.UserRequest{
		Username: "abcdefg",
		Password: "abcdefg",
	}, 204, nil))

	results.Add("test post user with duplicate username", postAndEnsureStatus(path, models.UserRequest{
		Username: "abcdefg",
		Password: "abcdefg",
	}, 409, nil))

	return results
}

func buildQueryStringRequestAndDo(method string, path string, body any, status int, respObjPtr any) error {
	jsonData, _ := json.Marshal(body)
	keyValMap := make(map[string]string)
	json.Unmarshal(jsonData, &keyValMap)

	values := url.Values{}
	for key, value := range keyValMap {
		values.Set(lowerFirst(key), value)
	}

	url := "http://localhost:8080" + path + "?" + values.Encode()

	req, _ := http.NewRequest(method, url, bytes.NewReader(jsonData))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+authToken)
	return doAndEnsureStatus(req, status, respObjPtr)
}

func lowerFirst(s string) string {
	if len(s) == 0 {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

func buildJsonRequestAndDo(method string, path string, body any, status int, respObjPtr any) error {
	jsonData, _ := json.Marshal(body)
	url := "http://localhost:8080" + path

	req, _ := http.NewRequest(method, url, bytes.NewReader(jsonData))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+authToken)
	return doAndEnsureStatus(req, status, respObjPtr)
}

func doAndEnsureStatus(req *http.Request, status int, respObjPtr any) error {

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error calling %v %w", req.URL.Path, err)
	}

	if resp.StatusCode != status {
		var body string
		if resp != nil {
			bts, _ := io.ReadAll(resp.Body)
			body = string(bts)
		}
		return fmt.Errorf("error calling %v.  expected status %v but got %v with body %v", req.URL.Path, status, resp.Status, body)
	}
	if respObjPtr != nil {
		bts, _ := io.ReadAll(resp.Body)
		return json.Unmarshal(bts, respObjPtr)
	}

	return nil
}

func deleteAndEnsureStatus(path string, status int, respObjPtr any) error {
	return buildQueryStringRequestAndDo(http.MethodDelete, path, nil, status, respObjPtr)
}

func getAndEnsureStatus(path string, body any, status int, respObjPtr any) error {
	return buildQueryStringRequestAndDo(http.MethodGet, path, body, status, respObjPtr)
}

func putAndEnsureStatus(path string, body any, status int, respObjPtr any) error {
	return buildJsonRequestAndDo(http.MethodPut, path, body, status, respObjPtr)
}

func postAndEnsureStatus(path string, body any, status int, respObjPtr any) error {
	return buildJsonRequestAndDo(http.MethodPost, path, body, status, respObjPtr)
}
