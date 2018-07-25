package main

import (
  "encoding/json"
  "github.com/gorilla/mux"
  "log"
  "net/http"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "time"
  "io/ioutil"
  "strconv"
  sq "github.com/Masterminds/squirrel"
)

var db *sql.DB

// The List Type
type List struct {
  Id int `json:"id,omitempty"`
  Name string `json:"name,omitempty"`
  Date_created time.Time `json:"date_created,omitempty"`
}

// The Task Type
type Task struct {
  Id int `json:"id,omitempty"`
  List_id int `json:"list_id,omitempty"`
  Description string `json:"description,omitempty"`
  Date_created time.Time `json:"date_created,omitempty"`
}

// The Query Response Type
type QueryResponse struct {
  Last_insert_id int64 `json:"last_insert_id,omitempty"`
  Row_count int64 `json:"rows_updated,omitempty"`
}

func logExitFatalError (e error) {
  if e != nil {
    log.Fatal(e)
  }
}

// Create List
func CreateList(w http.ResponseWriter, r *http.Request) {
  var list List
  body, err := ioutil.ReadAll(r.Body)
  logExitFatalError(err)
  err = json.Unmarshal(body, &list)
  logExitFatalError(err)
  stmt, err := db.Prepare("INSERT INTO list(name) VALUES(?)")
  logExitFatalError(err)
  res, err := stmt.Exec(list.Name)
  logExitFatalError(err)
  lastId, err := res.LastInsertId()
  logExitFatalError(err)
  rowCnt, err := res.RowsAffected()
  logExitFatalError(err)
  var queryResponse QueryResponse
  queryResponse.Last_insert_id = lastId
  queryResponse.Row_count = rowCnt
  json.NewEncoder(w).Encode(queryResponse)
}

// Create Task
func CreateTask(w http.ResponseWriter, r *http.Request) {
  var task Task
  body, err := ioutil.ReadAll(r.Body)
  logExitFatalError(err)
  err = json.Unmarshal(body, &task)
  logExitFatalError(err)
  params := mux.Vars(r)
  list_id, err := strconv.Atoi(params["list_id"])
  logExitFatalError(err)
  task.List_id = list_id

  stmt, err := db.Prepare("INSERT INTO task(list_id, description) VALUES(?, ?)")
  logExitFatalError(err)
  res, err := stmt.Exec(task.List_id, task.Description)
  logExitFatalError(err)
  lastId, err := res.LastInsertId()
  logExitFatalError(err)
  rowCnt, err := res.RowsAffected()
  logExitFatalError(err)
  var queryResponse QueryResponse
  queryResponse.Last_insert_id = lastId
  queryResponse.Row_count = rowCnt
  json.NewEncoder(w).Encode(queryResponse)
}

// Delete Task From List
func DeleteTask(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)
  list_id, err := strconv.Atoi(params["list_id"])
  logExitFatalError(err)
  task_id, err := strconv.Atoi(params["task_id"])
  logExitFatalError(err)

  stmt, err := db.Prepare("DELETE FROM task WHERE id = ? and list_id = ?")
  logExitFatalError(err)
  res, err := stmt.Exec(task_id, list_id)
  logExitFatalError(err)
  rowCnt, err := res.RowsAffected()
  logExitFatalError(err)
  var queryResponse QueryResponse
  queryResponse.Row_count = rowCnt
  json.NewEncoder(w).Encode(queryResponse)
}

// Update Task
func UpdateTask(w http.ResponseWriter, r *http.Request) {
  var task Task
  body, err := ioutil.ReadAll(r.Body)
  logExitFatalError(err)
  err = json.Unmarshal(body, &task)
  logExitFatalError(err)
  params := mux.Vars(r)
  task_id, err := strconv.Atoi(params["task_id"])
  logExitFatalError(err)

  statement := sq.Update("task").Where(sq.Eq{"id": task_id})
  if &task.List_id != nil && task.List_id != 0 {
      statement = statement.Set("list_id", task.List_id)
  }
  if &task.Description != nil && task.Description != "" {
      statement = statement.Set("description", task.Description)
  }

  res, err := statement.RunWith(db).Exec()
  logExitFatalError(err)

  rowCnt, err := res.RowsAffected()
  logExitFatalError(err)
  var queryResponse QueryResponse
  queryResponse.Row_count = rowCnt
  json.NewEncoder(w).Encode(queryResponse)
}

// Delete List
func DeleteList(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)
  list_id, err := strconv.Atoi(params["list_id"])
  logExitFatalError(err)

  stmt, err := db.Prepare("DELETE FROM list WHERE id = ?")
  logExitFatalError(err)
  res, err := stmt.Exec(list_id)
  logExitFatalError(err)
  rowCnt, err := res.RowsAffected()
  logExitFatalError(err)
  var queryResponse QueryResponse
  queryResponse.Row_count = rowCnt
  json.NewEncoder(w).Encode(queryResponse)
}

// connect to db
func dbConn() () {
    var err error
    dbDriver := "mysql"
    dbUser := "root"
    dbPass := "root"
    dbHost := "tcp(127.0.0.1:3306)"
    dbName := "restapi"
    db, err = sql.Open(dbDriver, dbUser+":"+dbPass+"@"+dbHost+"/"+dbName)
    logExitFatalError(err)
    log.Println("Connection Established")
}

// main function to boot up everything
func main() {
    dbConn()
    router := mux.NewRouter()
    router.HandleFunc("/list", CreateList).Methods("POST")
    router.HandleFunc("/list/{list_id}/task", CreateTask).Methods("POST")
    router.HandleFunc("/list/{list_id}/task/{task_id}", DeleteTask).Methods("DELETE")
    router.HandleFunc("/task/{task_id}", UpdateTask).Methods("PUT")
    router.HandleFunc("/list/{list_id}", DeleteList).Methods("DELETE")
    log.Fatal(http.ListenAndServe(":8000", router))
}
