package inmemory

import (
	"context"
	"mcg-app-backend/service/customerrors"
	"mcg-app-backend/service/models"
	"sync"
)

type InMemoryRepo struct {
	mutex               sync.RWMutex
	patients            map[int]models.Patient
	attatchments        map[int]models.Attatchment
	diagnosedConditions map[int]models.DiagnosedCondition
	users               map[string]models.User
	nextPatientId       int
	nextAttatchmentId   int
	nextConditionId     int
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		patients:            make(map[int]models.Patient),
		attatchments:        make(map[int]models.Attatchment),
		diagnosedConditions: make(map[int]models.DiagnosedCondition),
		users:               make(map[string]models.User),
		nextPatientId:       1,
		nextAttatchmentId:   1,
		nextConditionId:     1,
	}
}

func (r *InMemoryRepo) InsertPatient(ctx context.Context, patient models.Patient) (int, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	id := r.nextPatientId
	r.nextPatientId++

	patient.Id = id
	r.patients[id] = patient

	return id, nil
}

func (r *InMemoryRepo) UpdatePatient(ctx context.Context, patient models.Patient) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.patients[patient.Id]; !exists {
		return customerrors.NewInvalidInputError("patient not found")
	}

	r.patients[patient.Id] = patient
	return nil
}

func (r *InMemoryRepo) DeletePatient(ctx context.Context, id int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.patients[id]; !exists {
		return customerrors.NewInvalidInputError("patient not found")
	}

	delete(r.patients, id)
	return nil
}

func (r *InMemoryRepo) GetCountOfPatientId(ctx context.Context, id int) (int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if _, exists := r.patients[id]; exists {
		return 1, nil
	}
	return 0, nil
}

func (r *InMemoryRepo) GetCountOfExternalIdentifier(ctx context.Context, externalIdentifier string) (int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	count := 0
	for _, patient := range r.patients {
		if patient.ExternalIdentifier == externalIdentifier {
			count++
		}
	}
	return count, nil
}

func (r *InMemoryRepo) InsertAttatchment(ctx context.Context, attatchment models.Attatchment) (int, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	id := r.nextAttatchmentId
	r.nextAttatchmentId++

	attatchment.Id = id
	r.attatchments[id] = attatchment

	return id, nil
}

func (r *InMemoryRepo) InsertDiagnosedCondition(ctx context.Context, condition models.DiagnosedCondition) (int, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	id := r.nextConditionId
	r.nextConditionId++

	condition.Id = id
	r.diagnosedConditions[id] = condition

	return id, nil
}

func (r *InMemoryRepo) DeleteAttatchmentsByPatientId(ctx context.Context, patientId int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for id, attatchment := range r.attatchments {
		if attatchment.PatientId == patientId {
			delete(r.attatchments, id)
		}
	}

	return nil
}

func (r *InMemoryRepo) DeleteDiagnosedConditionsByPatientId(ctx context.Context, patientId int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for id, condition := range r.diagnosedConditions {
		if condition.PatientId == patientId {
			delete(r.diagnosedConditions, id)
		}
	}

	return nil
}

func (r *InMemoryRepo) DeleteDiagnosedCondition(ctx context.Context, conditionId int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.diagnosedConditions[conditionId]; !exists {
		return customerrors.NewInvalidInputError("diagnosed condition not found")
	}

	delete(r.diagnosedConditions, conditionId)
	return nil
}

func (r *InMemoryRepo) DeleteAttatchment(ctx context.Context, attatchmentId int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.attatchments[attatchmentId]; !exists {
		return customerrors.NewInvalidInputError("attatchment not found")
	}

	delete(r.attatchments, attatchmentId)
	return nil
}

func (r *InMemoryRepo) SearchPatients(ctx context.Context, search models.PatientSearch) ([]models.Patient, error) {
	var patients []models.Patient

	matchedPatientIds := make(map[int]bool)

	conditionsByPatientId := make(map[int][]models.DiagnosedCondition)
	attatchmentsByPatientId := make(map[int][]models.Attatchment)

	for _, condition := range r.diagnosedConditions {
		conditionsByPatientId[condition.PatientId] = append(conditionsByPatientId[condition.PatientId], condition)
		if (search.DiagnosedConditionCode != "" && condition.Code == search.DiagnosedConditionCode) ||
			(search.DiagnosedConditionName != "" && condition.Name == search.DiagnosedConditionName) {
			matchedPatientIds[condition.PatientId] = true
		}
	}

	for _, attatchment := range r.attatchments {
		attatchmentsByPatientId[attatchment.PatientId] = append(attatchmentsByPatientId[attatchment.PatientId], attatchment)
		if (search.AttatchmentName != "" && attatchment.Name == search.AttatchmentName) ||
			(search.AttatchmentType != "" && attatchment.Type == search.AttatchmentType) {
			matchedPatientIds[attatchment.PatientId] = true
		}
	}

	for _, patient := range r.patients {
		if (search.Name != "" && patient.Name == search.Name) ||
			(search.Phone != "" && patient.PhoneNumber == search.Phone) ||
			(search.Address != "" && patient.Address == search.Address) ||
			(search.ExternalIdentifier != "" && patient.ExternalIdentifier == search.ExternalIdentifier) {
			matchedPatientIds[patient.Id] = true
		}

		if matchedPatientIds[patient.Id] {
			patient.Attatchments = attatchmentsByPatientId[patient.Id]
			patient.DiagnosedConditions = conditionsByPatientId[patient.Id]
			patients = append(patients, patient)
		}
	}
	return patients, nil
}

func (r *InMemoryRepo) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	return r.users[username], nil
}

func (r *InMemoryRepo) InsertUser(ctx context.Context, user models.User) error {
	r.users[user.Username] = user
	return nil
}
