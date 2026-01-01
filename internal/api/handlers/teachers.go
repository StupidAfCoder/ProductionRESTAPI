package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"schoolREST/internal/models"
	"schoolREST/internal/repository/sqlconnect"
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
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		http.Error(w, "Error Connecting to the database", http.StatusInternalServerError)
	}

	defer db.Close()

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

	var teacher models.Teacher
	err = db.QueryRow("SELECT id , first_name , last_name , email , class , subject FROM teachers WHERE id = ?", id).Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
	if err == sql.ErrNoRows {
		http.Error(w, "No Teacher associated with that ID ", http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacher)
}

func addTeachersHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		http.Error(w, "Error Connecting to the database", http.StatusInternalServerError)
	}

	defer db.Close()

	var newTeachers []models.Teacher
	err = json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "Invalid Body Parameters", http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("INSERT INTO teachers (first_name , last_name , email , class , subject) VALUES (?,?,?,?,?)")
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Error preparing the statement ", http.StatusInternalServerError)
		return
	}

	defer stmt.Close()

	addedTeachers := make([]models.Teacher, len(newTeachers))
	for i, teacher := range newTeachers {
		resp, err := stmt.Exec(teacher.FirstName, teacher.LastName, teacher.Email, teacher.Class, teacher.Subject)
		if err != nil {
			http.Error(w, "Error executing the Sql query", http.StatusInternalServerError)
			return
		}
		lastInsertID, err := resp.LastInsertId()
		if err != nil {
			http.Error(w, "Error retrieving the last ID of the inserted  element", http.StatusInternalServerError)
			return
		}
		teacher.ID = int(lastInsertID)
		addedTeachers[i] = teacher
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
