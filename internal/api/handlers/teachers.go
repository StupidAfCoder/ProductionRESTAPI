package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"schoolREST/internal/models"
	"schoolREST/internal/repository/sqlconnect"
	"strconv"
	"strings"
)

func GetTeachersHandler(w http.ResponseWriter, r *http.Request) {

	var teachers []models.Teacher
	teachers, err := sqlconnect.GetTeachersDb(teachers, r)
	if err != nil {
		return
	}
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "Success",
		Count:  len(teachers),
		Data:   teachers,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetOneTeacherHandler(w http.ResponseWriter, r *http.Request) {

	idstr := r.PathValue("id")
	fmt.Println(idstr)

	//Handling the path parameter
	id, err := strconv.Atoi(idstr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	teacher, err := sqlconnect.GetOneTeacherDbHandler(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("Contc1ent-Type", "application/json")
	json.NewEncoder(w).Encode(teacher)
}

func AddTeachersHandler(w http.ResponseWriter, r *http.Request) {
	var newTeachers []models.Teacher
	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "Invalid Body Parameters", http.StatusBadRequest)
		return
	}

	addedTeachers, err := sqlconnect.AddTeachersDBHandler(newTeachers)
	if err != nil {
		log.Println(err)
		return
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

func UpdateTeachersHandler(w http.ResponseWriter, r *http.Request) {
	idstr := strings.TrimPrefix(r.URL.Path, "/teachers/")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Teacher ID", http.StatusBadRequest)
		return
	}

	var updatedTeacher models.Teacher
	err = json.NewDecoder(r.Body).Decode(&updatedTeacher)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Request Paylod", http.StatusBadRequest)
		return
	}

	updatedTeacherFromDB, err := sqlconnect.UpdateTeachersDBHandler(id, updatedTeacher)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTeacherFromDB)
}

func PatchTeachersHandler(w http.ResponseWriter, r *http.Request) {
	var updates []map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	err = sqlconnect.PatchTeachersDBHandler(updates)
	if err != nil {
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}

func PatchOneTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idstr := strings.TrimPrefix(r.URL.Path, "/teachers/")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Teacher ID", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Request Paylod", http.StatusBadRequest)
		return
	}

	updatedTeacher, err := sqlconnect.PatchOneTeacherDBHandler(id, updates)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTeacher)

}

func DeleteTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idstr := strings.TrimPrefix(r.URL.Path, "/teachers/")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Teacher ID", http.StatusBadRequest)
		return
	}

	err = sqlconnect.DeleteOneTeacherDB(id)
	if err != nil {
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func DeleteMultipleTeachersHandler(w http.ResponseWriter, r *http.Request) {
	var ids []int
	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error converting the data ", http.StatusBadRequest)
		return
	}

	deletedIDs, err := sqlconnect.DeleteMultipleTeachesDB(ids)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Status     string `json:"status"`
		DeletedIDs []int  `json:"deleted_ids"`
	}{
		Status:     "Teachers successfully deleted",
		DeletedIDs: deletedIDs,
	}
	json.NewEncoder(w).Encode(response)
}
