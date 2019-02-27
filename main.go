package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
)

func main() {
	handler := http.FileServer(http.Dir("static"))
	log.Println("Server started at port 3000...")

	http.Handle("/static/", http.StripPrefix("/static/", handler))
	http.HandleFunc("/", serveTemplates)

	errs := make(chan error, 2)
	server := &http.Server{Addr: ":3000"}

	go func() {
		errs <- server.ListenAndServe()
	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt)
		errs <- fmt.Errorf("%s", <-c)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-errs
	log.Println("\nGracefully shutting down...")
}

// Parses templates and serve on respective routes.
func serveTemplates(w http.ResponseWriter, r *http.Request) {
	lp := filepath.Join("static/templates", "layout.html")
	fp := filepath.Join("static/templates", "sample.html")
	tem, err := template.ParseFiles(lp, fp)
	if err != nil {
		log.Println(err)
	}
	err = tem.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		log.Println(err)
	}
}
