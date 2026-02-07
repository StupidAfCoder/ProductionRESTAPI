package sqlconnect

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"schoolREST/internal/models"
	"strconv"
	"strings"
)

func isValidSortOrder(order string) bool {
	return order == "asc" || order == "desc"
}

func isValidSortField(field string) bool {
	validFields := map[string]bool{
		"first_name": true,
		"last_name":  true,
		"email":      true,
		"class":      true,
		"subject":    true,
	}
	return validFields[field]
}

func sortByParams(r *http.Request, query string) string {
	sortParams := r.URL.Query()["sortby"]
	if len(sortParams) > 0 {
		query += " ORDER BY"
		for i, param := range sortParams {
			parts := strings.Split(param, ":")
			if len(parts) != 2 {
				continue
			}
			field, order := parts[0], parts[1]
			if !isValidSortField(field) || !isValidSortOrder(order) {
				continue
			}
			if i > 0 {
				query += ","
			}
			query += " " + field + " " + order
		}
	}
	return query
}

func filterQuery(r *http.Request, query string, args []interface{}) (string, []interface{}) {
	params := map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
		"email":      "email",
		"class":      "class",
		"subject":    "subject",
	}
	for param, dbField := range params {
		value := r.URL.Query().Get(param)
		if value != "" {
			query += " AND " + dbField + " = ?"
			args = append(args, value)
		}
	}
	return query, args
}

func GetTeachersDb(teachers []models.Teacher, r *http.Request) ([]models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		//http.Error(w, "Error Connecting to the database", http.StatusInternalServerError)
		return nil, err
	}

	defer db.Close()

	query := "SELECT id , first_name , last_name , email , class , subject FROM teachers WHERE 1=1"
	var args []interface{}

	query, args = filterQuery(r, query, args)

	query = sortByParams(r, query)

	rows, err := db.Query(query, args...)
	if err != nil {
		fmt.Println(err.Error())
		// http.Error(w, "Database Query Error", http.StatusInternalServerError)
		return nil, err
	}
	defer rows.Close()

	// teachers := make([]models.Teacher, 0)
	for rows.Next() {
		var teacher models.Teacher
		err := rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
		if err != nil {
			// http.Error(w, "Error scanning result database", http.StatusInternalServerError)
			return nil, err
		}
		teachers = append(teachers, teacher)
	}
	return teachers, nil
}

func GetOneTeacherDbHandler(id int) (models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		// http.Error(w, "Error Connecting to the database", http.StatusInternalServerError)
		return models.Teacher{}, err
	}

	defer db.Close()

	var teacher models.Teacher
	err = db.QueryRow("SELECT id , first_name , last_name , email , class , subject FROM teachers WHERE id = ?", id).Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
	if err == sql.ErrNoRows {
		// http.Error(w, "No Teacher associated with that ID ", http.StatusNotFound)
		return models.Teacher{}, err
	} else if err != nil {
		fmt.Println(err.Error())
		// http.Error(w, "Database query error", http.StatusInternalServerError)
		return models.Teacher{}, err
	}
	return teacher, nil
}

func AddTeachersDBHandler(newTeachers []models.Teacher) ([]models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		// http.Error(w, "Error Connecting to the database", http.StatusInternalServerError)
		return nil, err
	}

	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO teachers (first_name , last_name , email , class , subject) VALUES (?,?,?,?,?)")
	if err != nil {
		fmt.Println(err.Error())
		// http.Error(w, "Error preparing the statement ", http.StatusInternalServerError)
		return nil, err
	}

	defer stmt.Close()

	addedTeachers := make([]models.Teacher, len(newTeachers))
	for i, teacher := range newTeachers {
		resp, err := stmt.Exec(teacher.FirstName, teacher.LastName, teacher.Email, teacher.Class, teacher.Subject)
		if err != nil {
			// http.Error(w, "Error executing the Sql query", http.StatusInternalServerError)
			return nil, err
		}
		lastInsertID, err := resp.LastInsertId()
		if err != nil {
			// http.Error(w, "Error retrieving the last ID of the inserted  element", http.StatusInternalServerError)
			return nil, err
		}
		teacher.ID = int(lastInsertID)
		addedTeachers[i] = teacher
	}
	return addedTeachers, nil
}

func UpdateTeachersDBHandler(id int, updatedTeacher models.Teacher) (models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		log.Println(err)
		// http.Error(w, "Error connecting database", http.StatusInternalServerError)
		return models.Teacher{}, err
	}
	defer db.Close()

	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT id , first_name , last_name , email , class , subject FROM teachers WHERE id = ?", id).Scan(&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)

	if err != nil {
		if err == sql.ErrNoRows {
			// http.Error(w, "Teacher not found", http.StatusInternalServerError)
			return models.Teacher{}, err
		}
		log.Println(err)
		// http.Error(w, "Error Querying Information", http.StatusInternalServerError)
		return models.Teacher{}, err
	}

	updatedTeacher.ID = existingTeacher.ID
	_, err = db.Exec("UPDATE teachers SET first_name = ? , last_name = ? , email = ? , class = ? , subject = ? WHERE id = ?", updatedTeacher.FirstName, updatedTeacher.LastName, updatedTeacher.Email, updatedTeacher.Class, updatedTeacher.Subject, updatedTeacher.ID)
	if err != nil {
		// http.Error(w, "Error updating teacher", http.StatusInternalServerError)
		return models.Teacher{}, err
	}
	return updatedTeacher, nil
}

func PatchTeachersDBHandler(updates []map[string]interface{}) error {
	db, err := ConnectDB()
	if err != nil {
		log.Println(err)
		// http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		// http.Error(w, "Error during transaction starting", http.StatusInternalServerError)
		return err
	}

	for _, update := range updates {
		idStr, ok := update["id"].(string)
		if !ok {
			tx.Rollback()
			// http.Error(w, "Invalid Teacher ID", http.StatusBadRequest)
			return err
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			tx.Rollback()
			// http.Error(w, "Error converting String to int for ID", http.StatusBadRequest)
			return err
		}

		var teacherFromDb models.Teacher
		err = db.QueryRow("SELECT id , first_name , last_name , email , class , subject FROM teachers WHERE id = ?", id).Scan(&teacherFromDb.ID, &teacherFromDb.FirstName, &teacherFromDb.LastName, &teacherFromDb.Email, &teacherFromDb.Class, &teacherFromDb.Subject)
		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				// http.Error(w, "Teacher not found", http.StatusNotFound)
				return err
			}
			// http.Error(w, "Error retrieving teacher", http.StatusInternalServerError)
			return err
		}
		teacherVal := reflect.ValueOf(&teacherFromDb).Elem()
		teacherType := teacherVal.Type()

		for key, value := range update {
			if key == "id" {
				continue //skip this iteration
			}
			for i := 0; i < teacherVal.NumField(); i++ {
				field := teacherType.Field(i)
				if field.Tag.Get("json") == key+",omitempty" {
					fieldVal := teacherVal.Field(i)
					if fieldVal.CanSet() {
						val := reflect.ValueOf(value)
						if val.Type().ConvertibleTo(fieldVal.Type()) {
							fieldVal.Set(val.Convert(fieldVal.Type()))
						} else {
							tx.Rollback()
							log.Printf("Unsuccessful in converting the %v to %v", val.Type(), fieldVal.Type())
							return err
						}
					}
					break
				}
			}
		}

		_, err = tx.Exec("UPDATE teachers SET first_name = ? , last_name = ? , email = ? , class = ? , subject = ? WHERE id = ?", teacherFromDb.FirstName, teacherFromDb.LastName, teacherFromDb.Email, teacherFromDb.Class, teacherFromDb.Subject, teacherFromDb.ID)
		if err != nil {
			// http.Error(w, "Error updating teacher", http.StatusInternalServerError)
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		// http.Error(w, "Error commiting the transaction", http.StatusInternalServerError)
		return err
	}
	return nil
}

func PatchOneTeacherDBHandler(id int, updates map[string]interface{}) (models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		log.Println(err)
		// http.Error(w, "Error connecting database", http.StatusInternalServerError)
		return models.Teacher{}, err
	}
	defer db.Close()

	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT id , first_name , last_name , email , class , subject FROM teachers WHERE id = ?", id).Scan(&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)

	if err != nil {
		if err == sql.ErrNoRows {
			// http.Error(w, "Teacher not found", http.StatusInternalServerError)
			return models.Teacher{}, err
		}
		log.Println(err)
		// http.Error(w, "Error Querying Information", http.StatusInternalServerError)
		return models.Teacher{}, err
	}

	teacherVal := reflect.ValueOf(&existingTeacher).Elem()
	teacherType := teacherVal.Type()

	for k, v := range updates {
		for i := 0; i < teacherVal.NumField(); i++ {
			field := teacherType.Field(i)
			if field.Tag.Get("json") == k+",omitempty" {
				if teacherVal.Field(i).CanSet() {
					fieldVal := teacherVal.Field(i)
					fieldVal.Set(reflect.ValueOf(v).Convert(fieldVal.Type()))
				}
			}
		}
	}

	fmt.Println(existingTeacher)

	_, err = db.Exec("UPDATE teachers SET first_name = ? , last_name = ? , email = ? , class = ? , subject = ? WHERE id = ?", existingTeacher.FirstName, existingTeacher.LastName, existingTeacher.Email, existingTeacher.Class, existingTeacher.Subject, existingTeacher.ID)
	if err != nil {
		// http.Error(w, "Error updating teacher", http.StatusInternalServerError)
		return models.Teacher{}, err
	}
	return existingTeacher, nil
}

func DeleteOneTeacherDB(id int) error {
	db, err := ConnectDB()
	if err != nil {
		log.Println(err)
		// http.Error(w, "Error connecting database", http.StatusInternalServerError)
		return err
	}
	defer db.Close()

	result, err := db.Exec("DELETE from teachers WHERE id = ?", id)
	if err != nil {
		// http.Error(w, "Error executing Sql query", http.StatusInternalServerError)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// http.Error(w, "Error getting rows affected", http.StatusInternalServerError)
		return err
	}
	if rowsAffected == 0 {
		// http.Error(w, "Rows affected are 0", http.StatusNotFound)
		return err
	}
	return nil
}

func DeleteMultipleTeachesDB(ids []int) ([]int, error) {
	db, err := ConnectDB()
	if err != nil {
		log.Println(err)
		// http.Error(w, "Error connecting database", http.StatusInternalServerError)
		return nil, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		// http.Error(w, "Error starting transaction ", http.StatusInternalServerError)
		return nil, err
	}

	stmt, err := tx.Prepare("DELETE from teachers WHERE id = ?")
	if err != nil {
		log.Println(err)
		tx.Rollback()
		// http.Error(w, "Error preparing the statement ", http.StatusInternalServerError)
		return nil, err
	}
	defer stmt.Close()

	deletedIDs := []int{}
	for _, id := range ids {
		result, err := stmt.Exec(id)
		if err != nil {
			log.Println(err)
			tx.Rollback()
			// http.Error(w, "Error executing the statemenent", http.StatusInternalServerError)
			return nil, err
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Println(err)
			tx.Rollback()
			// http.Error(w, "Error retrieving the rows affected", http.StatusInternalServerError)
			return nil, err
		}
		if rowsAffected > 0 {
			deletedIDs = append(deletedIDs, id)
		}
		if rowsAffected < 1 {
			tx.Rollback()
			// http.Error(w, fmt.Sprintf("ID %d does not exist", id), http.StatusBadRequest)
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		// http.Error(w, "Error commiting transaction", http.StatusInternalServerError)
		return nil, err
	}

	if len(deletedIDs) < 1 {
		// http.Error(w, "IDs do not exist", http.StatusBadRequest)
		return nil, err
	}
	return deletedIDs, nil
}
