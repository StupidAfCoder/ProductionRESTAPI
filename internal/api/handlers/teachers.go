package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"schoolREST/internal/models"
	"strconv"
	"strings"
	"sync"
)

var (
	teachers = make(map[int]models.Teacher)
	mutex    = &sync.Mutex{}
	nextID   = 1
)

func init() {
	//Creating some dummy data for the API
	teachers[nextID] = models.Teacher{
		ID:        nextID,
		FirstName: "Krish",
		LastName:  "Jain",
		Class:     "10A",
		Subject:   "Math",
	}
	nextID++
	teachers[nextID] = models.Teacher{
		ID:        nextID,
		FirstName: "Sonali",
		LastName:  "Sharma",
		Class:     "5A",
		Subject:   "English",
	}
	nextID++
	teachers[nextID] = models.Teacher{
		ID:        nextID,
		FirstName: "Rakshash",
		LastName:  "Satan",
		Class:     "12B",
		Subject:   "Physics",
	}
	nextID++
}

func TeachersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	switch r.Method {
	case http.MethodGet:
		getTeachersHandler(w, r)
	case http.MethodPost:
		addTeachersHandler(w, r)
	case http.MethodDelete:
		w.Write([]byte("This is the DELETE Method for teachers routes"))
	case http.MethodPatch:
		w.Write([]byte("This is the PATCH Method for teachers routes"))
	}

}

func getTeachersHandler(w http.ResponseWriter, r *http.Request) {

	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	idstr := strings.TrimSuffix(path, "/")
	fmt.Println(idstr)

	if idstr == "" {
		firstName := r.URL.Query().Get("first_name")
		lastName := r.URL.Query().Get("last_name")

		teachersList := make([]models.Teacher, 0, len(teachers))
		for _, teacher := range teachers {
			if (firstName == "" || teacher.FirstName == firstName) && (lastName == "" || teacher.LastName == lastName) {
				teachersList = append(teachersList, teacher)
			}
		}

		response := struct {
			Status string           `json:"status"`
			Count  int              `json:"count"`
			Data   []models.Teacher `json:"data"`
		}{
			Status: "Success",
			Count:  len(teachersList),
			Data:   teachersList,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	}
	//Handling the path parameter
	id, err := strconv.Atoi(idstr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	teacher, exists := teachers[id]
	if !exists {
		http.Error(w, "Teacher with that ID was not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(teacher)
}

func addTeachersHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	var newTeachers []models.Teacher
	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "Invalid Body Parameters", http.StatusBadRequest)
		return
	}

	addedTeachers := make([]models.Teacher, len(newTeachers))
	for i, teacher := range newTeachers {
		teacher.ID = nextID
		teachers[nextID] = teacher
		addedTeachers[i] = teacher
		nextID++
	}

	w.Header().Set("Content-Type", "applicaton/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "Success",
		Count:  len(addedTeachers),
		Data:   addedTeachers,
	}
	json.NewEncoder(w).Encode(response)
}
