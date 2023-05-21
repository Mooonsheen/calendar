package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"server/internal/application"
	"server/internal/repository"

	"github.com/julienschmidt/httprouter"
)

func main() {
	ctx := context.Background()
	dbpool, err := repository.InitDBConn(ctx)
	if err != nil {
		log.Fatal("failed to init DB connection", err)
	}
	defer dbpool.Close()
	a := application.NewApp(ctx, dbpool)
	r := httprouter.New()
	a.Routes(r)
	srv := &http.Server{Addr: "0.0.0.0:8080", Handler: r}
	fmt.Println("It is alive! Try http://localhost:8080")
	srv.ListenAndServe()
}
