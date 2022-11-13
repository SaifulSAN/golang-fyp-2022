package main

import (
	db "decision-support-system-go/adapter"
	"decision-support-system-go/app/router"
	"fmt"
	"net/http"
)

func main() {
	dbPtr := db.NewDbConnection()
	appRouter := router.NewServer(dbPtr)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: appRouter,
	}

	srv.ListenAndServe()

	defer func() {
		dbPtr.Close()
		fmt.Println("DB Closed")
	}()
}
