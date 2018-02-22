package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var db *sql.DB
var err error
var arr int

type Tag struct {
	ID int `json:"id"`
}

type Questions struct {
	Question struct {
		Name                 string `json:"name"`
		Section              string `json:"section"`
		Position             int    `json:"position"`
		Title                string `json:"title"`
		TitleSpanish         string `json:"titleSpanish"`
		SubmitedValue        string `json:"submited_value"`
		SpanishSubmitedValue string `json:"spanish_submited_value"`
		Des                  string `json:"des"`
		Ans                  string `json:"ans"`
		ViewType             int    `json:"view_type"`
		ParentID             int    `json:"parent_id"`
		IsRequired           int    `json:"isRequired"`
		IsSubmitField        int    `json:"is_submit_field"`
		IsActive             int    `json:"is_active"`
	} `json:"question"`
	Options []struct {
		Name                 string `json:"name"`
		Section              string `json:"section"`
		Position             int    `json:"position"`
		Title                string `json:"title"`
		TitleSpanish         string `json:"titleSpanish"`
		SubmitedValue        string `json:"submited_value"`
		SpanishSubmitedValue string `json:"spanish_submited_value"`
		Des                  string `json:"des"`
		Ans                  string `json:"ans"`
		ViewType             int    `json:"view_type"`
		IsRequired           int    `json:"isRequired"`
		IsSubmitField        int    `json:"is_submit_field"`
		IsActive             int    `json:"is_active"`
	} `json:"options"`
	Validation struct {
		Messgae        string `json:"messgae"`
		MessageSpanish string `json:"messageSpanish"`
		Regx           string `json:"regx"`
		Format         string `json:"format"`
	} `json:"validation"`
}

type userAddHandler struct {
	db *sql.DB
}

type userGetHandler struct {
	db *sql.DB
}

func (u userAddHandler) insertInDatabase(data Questions) error {

	_, err = u.db.Exec("INSERT INTO profile_questions(name, section, position,title,titleSpanish,submited_value,spanish_submited_value,des,ans,view_type,parent_id,isRequired,is_submit_field,is_active) VALUES(?, ?, ?,?,?,?,?,?,?,?,?,?,?,?)", data.Question.Name, data.Question.Section, data.Question.Position, data.Question.Title, data.Question.TitleSpanish, data.Question.SubmitedValue, data.Question.SpanishSubmitedValue, data.Question.Des, data.Question.Ans, data.Question.ViewType, data.Question.ParentID, data.Question.IsRequired, data.Question.IsSubmitField, data.Question.IsActive)

	if len(data.Options) > 0 {
		results, err := u.db.Query("SELECT LAST_INSERT_ID()")
		if err != nil {
			// panic(err.Error())
			fmt.Println("err")
		}
		for results.Next() {
			var tag Tag
			err = results.Scan(&tag.ID)
			if err != nil {
				// panic(err.Error())
				fmt.Println("er")
			}
			arr = tag.ID
		}
		for i := 0; i <= 1; i++ {
			_, err = u.db.Exec("INSERT INTO profile_questions(name, section, position,title,titleSpanish,submited_value,spanish_submited_value,des,ans,view_type,parent_id,isRequired,is_submit_field,is_active) VALUES(?, ?, ?,?,?,?,?,?,?,?,?,?,?,?)", data.Options[i].Name, data.Options[i].Section, data.Options[i].Position, data.Options[i].Title, data.Options[i].TitleSpanish, data.Options[i].SubmitedValue, data.Options[i].SpanishSubmitedValue, data.Options[i].Des, data.Options[i].Ans, data.Options[i].ViewType, arr, data.Options[i].IsRequired, data.Options[i].IsSubmitField, data.Options[i].IsActive)
		}
	} else if data.Validation.Messgae != "" {
		_, err = u.db.Exec("INSERT INTO input_types(messgae,messageSpanish,regx,format) VALUES(?,?,?,?)", data.Validation.Messgae, data.Validation.MessageSpanish, data.Validation.Regx, data.Validation.Format)
	}

	return err
}

func (u userAddHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("reading the Request error")
	}

	var k Questions

	if err = json.Unmarshal(body, &k); err != nil {
		fmt.Println("unmarshall error ")
	}

	err = u.insertInDatabase(k)
	w.Write([]byte(`{"code ":"success"}`))
}

func (v userGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	type Validation struct {
		Regx   string `json:"regx"`
		Format string `json:"format"`
	}

	w.Header().Set("Content-Type", "application/json")
	rows, err := v.db.Query(`SELECT regx, format FROM input_types`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	enc := json.NewEncoder(w)
	for rows.Next() {
		validation := Validation{}
		err = rows.Scan(&validation.Regx, &validation.Format)

		if err != nil {
			panic(err)
		}
		json.NewEncoder(os.Stdout).Encode(validation)

		enc.Encode(validation)

	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

}

func main() {
	var err error
	db, err = sql.Open("mysql", "root:nfn@tcp(127.0.0.1:3306)/api")
	if err != nil {
		log.Fatalf("failed to open db: %s", err)
	}

	defer db.Close()

	handler := userAddHandler{
		db: db,
	}
	handler2 := userGetHandler{
		db: db,
	}

	http.Handle("/add", handler)
	http.Handle("/get", handler2)
	http.ListenAndServe(":1269", nil)
}

