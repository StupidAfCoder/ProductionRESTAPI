package sqlconnect

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"schoolREST/internal/models"
	"schoolREST/pkg/utils"
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
		return nil, utils.ErrorHandler(err, "Error Connecting to the database")
	}

	defer db.Close()

	query := "SELECT id , first_name , last_name , email , class , subject FROM teachers WHERE 1=1"
	var args []interface{}

	query, args = filterQuery(r, query, args)

	query = sortByParams(r, query)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Database Query Error")
	}
	defer rows.Close()

	// teachers := make([]models.Teacher, 0)
	for rows.Next() {
		var teacher models.Teacher
		err := rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error Retrieving data from the database")
		}
		teachers = append(teachers, teacher)
	}
	return teachers, nil
}

func GetOneTeacherDbHandler(id int) (models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "Error Connecting to the database")
	}

	defer db.Close()

	var teacher models.Teacher
	err = db.QueryRow("SELECT id , first_name , last_name , email , class , subject FROM teachers WHERE id = ?", id).Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
	if err == sql.ErrNoRows {
		return models.Teacher{}, utils.ErrorHandler(err, "No Data Found Associated")
	} else if err != nil {
		fmt.Println(err.Error())
		return models.Teacher{}, utils.ErrorHandler(err, "DataBase Query Error")
	}
	return teacher, nil
}

func AddTeachersDBHandler(newTeachers []models.Teacher) ([]models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error connecting to the database")
	}

	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO teachers (first_name , last_name , email , class , subject) VALUES (?,?,?,?,?)")
	if err != nil {
		return nil, utils.ErrorHandler(err, "Internal Server Error")
	}

	defer stmt.Close()

	addedTeachers := make([]models.Teacher, len(newTeachers))
	for i, teacher := range newTeachers {
		resp, err := stmt.Exec(teacher.FirstName, teacher.LastName, teacher.Email, teacher.Class, teacher.Subject)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Database Error")
		}
		lastInsertID, err := resp.LastInsertId()
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error retrieving the last ID")
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
		return models.Teacher{}, utils.ErrorHandler(err, "Error Connecting to the database")
	}
	defer db.Close()

	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT id , first_name , last_name , email , class , subject FROM teachers WHERE id = ?", id).Scan(&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Teacher{}, utils.ErrorHandler(err, "No Data Found Associated")
		}
		log.Println(err)
		return models.Teacher{}, utils.ErrorHandler(err, "Error Querying Information")
	}

	updatedTeacher.ID = existingTeacher.ID
	_, err = db.Exec("UPDATE teachers SET first_name = ? , last_name = ? , email = ? , class = ? , subject = ? WHERE id = ?", updatedTeacher.FirstName, updatedTeacher.LastName, updatedTeacher.Email, updatedTeacher.Class, updatedTeacher.Subject, updatedTeacher.ID)
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "Error Retreving Information")
	}
	return updatedTeacher, nil
}

func PatchTeachersDBHandler(updates []map[string]interface{}) error {
	db, err := ConnectDB()
	if err != nil {
		return utils.ErrorHandler(err, "Error Connecting to the database")
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return utils.ErrorHandler(err, "Internal Transaction Error")
	}

	for _, update := range updates {
		idStr, ok := update["id"].(string)
		if !ok {
			tx.Rollback()
			return utils.ErrorHandler(err, "Invalid Teacher ID")
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			tx.Rollback()
			return utils.ErrorHandler(err, "Internal Server Error")
		}

		var teacherFromDb models.Teacher
		err = db.QueryRow("SELECT id , first_name , last_name , email , class , subject FROM teachers WHERE id = ?", id).Scan(&teacherFromDb.ID, &teacherFromDb.FirstName, &teacherFromDb.LastName, &teacherFromDb.Email, &teacherFromDb.Class, &teacherFromDb.Subject)
		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				return utils.ErrorHandler(err, "Teacher Not Found")
			}
			return utils.ErrorHandler(err, "Error Retrieving Teacher")
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
			return utils.ErrorHandler(err, "Error Updating Teacher")
		}
	}
	err = tx.Commit()
	if err != nil {
		return utils.ErrorHandler(err, "Transaction Error")
	}
	return nil
}

func PatchOneTeacherDBHandler(id int, updates map[string]interface{}) (models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "Error connecting to the database")
	}
	defer db.Close()

	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT id , first_name , last_name , email , class , subject FROM teachers WHERE id = ?", id).Scan(&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Teacher{}, utils.ErrorHandler(err, "No Teacher Found")
		}
		return models.Teacher{}, utils.ErrorHandler(err, "Error retriveing information")
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
		return models.Teacher{}, utils.ErrorHandler(err, "Error Updating Teacher")
	}
	return existingTeacher, nil
}

func DeleteOneTeacherDB(id int) error {
	db, err := ConnectDB()
	if err != nil {
		return utils.ErrorHandler(err, "Error connecting to the database")
	}
	defer db.Close()

	result, err := db.Exec("DELETE from teachers WHERE id = ?", id)
	if err != nil {
		return utils.ErrorHandler(err, "Internal Server Error")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utils.ErrorHandler(err, "Error getting the affected rows")
	}
	if rowsAffected == 0 {
		return utils.ErrorHandler(err, "No Rows were affected")
	}
	return nil
}

func DeleteMultipleTeachesDB(ids []int) ([]int, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error connecting to the database")
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error starting the transaction")
	}

	stmt, err := tx.Prepare("DELETE from teachers WHERE id = ?")
	if err != nil {
		tx.Rollback()
		return nil, utils.ErrorHandler(err, "Error during transaction")
	}
	defer stmt.Close()

	deletedIDs := []int{}
	for _, id := range ids {
		result, err := stmt.Exec(id)
		if err != nil {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, "Error executing the statement")
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, "Error retreving the rows affected")
		}
		if rowsAffected > 0 {
			deletedIDs = append(deletedIDs, id)
		}
		if rowsAffected < 1 {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, fmt.Sprintf("ID %d does not exist", id))
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error during transaction")
	}

	if len(deletedIDs) < 1 {
		return nil, utils.ErrorHandler(err, "Id does not exist")
	}
	return deletedIDs, nil
}
