package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func RetrieveStuff(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		row := db.QueryRowContext(context.Background(), "SELECT test_some_number FROM test_pgx_conn WHERE test_name = 'hellosaiful'")
		if err := row.Err(); err != nil {
			fmt.Println("db.QueryRowContext", err)
			return
		}

		var someNumber int

		if err := row.Scan(&someNumber); err != nil {
			fmt.Println("row.Scan", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]int)
		resp["yourNumber"] = someNumber
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("JSON MARSHALL ERROR")
		}

		w.Write(jsonResp)
	})
}

func RetrieveStuffSingular(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := `SELECT test_name, test_some_number FROM test_pgx_conn WHERE test_id = $1`
		//pathVar := (chi.URLParam(r, "id"))
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()
		row := db.QueryRowContext(ctx, query, (chi.URLParam(r, "id")))
		if err := row.Err(); err != nil {
			fmt.Println("db.QueryRowContext", err)
			return
		}

		//fmt.Println(&row)

		var p PgxConnTest

		if err := row.Scan(&p.TestName, &p.TestSomeNumber); err != nil {
			fmt.Println("row.Scan", err)
			return
		}

		fmt.Println(p)

		w.Header().Set("Content-Type", "application/json")
		//resp := make(map[string]int)
		// resp["yourNumber"] = someNumber
		// jsonResp, err := json.Marshal(resp)
		jsonResp, err := json.Marshal(p)
		if err != nil {
			log.Fatalf("JSON MARSHALL ERROR")
		}

		w.Write(jsonResp)
		//w.Write([]byte("Successful single get"))
	})
}

type PgxConnTest struct {
	TestName       string
	TestSomeNumber int
}

func PostStuff(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var p PgxConnTest

		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		query := `INSERT INTO test_pgx_conn(test_name, test_some_number) VALUES ($1, $2) RETURNING test_id`
		row := db.QueryRowContext(context.Background(), query, p.TestName, p.TestSomeNumber)

		if err := row.Err(); err != nil {
			fmt.Println("db.QueryRowContext", err)
			fmt.Println("Failed in insertion")
			log.Fatal(err)
			return
		}

		var rowID int

		if err := row.Scan(&rowID); err != nil && err != sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println(rowID)

		w.Write([]byte("Successful insertion"))
		//w.Write([]byte("\nSuccessful test"))
	})
}

func RegisterUserStandard(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var p UserStandard

		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tx, err := db.BeginTx(context.Background(), nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		defer tx.Rollback()

		HashedPw, err := HashPassword(p.UserPassword)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		UserTableQuery := `INSERT INTO app_user(user_name, user_phone, user_email, user_password, user_role, activated) VALUES ($1, $2, $3, $4, 1, true) RETURNING user_id`
		row := tx.QueryRowContext(context.Background(), UserTableQuery, &p.UserName, &p.UserPhone, &p.UserEmail, HashedPw)

		if err := row.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var rowID int

		if err := row.Scan(&rowID); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//fmt.Println(rowID)

		HomelessTableQuery := `INSERT INTO homeless(primary_emergency_contact, secondary_emergency_contact, user_id) VALUES ($1, $2, $3)`
		_, err = tx.ExecContext(context.Background(), HomelessTableQuery, &p.UserPrimaryEmergencyContact, &p.UserSecondaryEmergencyContact, rowID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err = tx.Commit(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	})
}

func RegisterUserOrganization(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var p UserOrganization

		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tx, err := db.BeginTx(context.Background(), nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		defer tx.Rollback()

		HashedPw, err := HashPassword(p.UserPassword)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		UserTableQuery := `INSERT INTO app_user(user_name, user_phone, user_email, user_password, user_role, activated) VALUES ($1, $2, $3, $4, 2, false) RETURNING user_id`
		//remember to hash the password later!!!!
		//or better make a separate function to hash then call it in handler
		row := tx.QueryRowContext(context.Background(), UserTableQuery, &p.UserName, &p.UserPhone, &p.UserEmail, HashedPw)

		if err := row.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var rowID int

		if err := row.Scan(&rowID); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//fmt.Println(rowID)

		OrganizationTableQuery := `INSERT INTO organization(org_secondary_contact, user_id) VALUES ($1, $2)`
		_, err = tx.ExecContext(context.Background(), OrganizationTableQuery, &p.OrgSecondaryContact, rowID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err = tx.Commit(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	})
}

func PutStuff(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var p PgxConnTest

		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		query := `UPDATE test_pgx_conn SET test_some_number = $1 WHERE test_id = $2`
		_, err = db.ExecContext(context.Background(), query, p.TestSomeNumber, (chi.URLParam(r, "id")))
		if err != nil {
			log.Fatal(err)
			fmt.Println("Failed in update")
			return
		}

		w.Write([]byte("Successful update"))
	})
}
