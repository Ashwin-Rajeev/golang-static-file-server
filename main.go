package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

func main() {
	handler := http.FileServer(http.Dir("static"))

	http.Handle("/static/", http.StripPrefix("/static/", handler))
	http.HandleFunc("/", serveTemplates)

	errs := make(chan error, 2)
	interrupt := make(chan os.Signal)
	server := &http.Server{Addr: ":3000"}

	go func() {
		log.Println("Server started at port 3000...")
		errs <- server.ListenAndServe()
	}()
	
	go func() {
		signal.Notify(interrupt, os.Interrupt)
		errs <- fmt.Errorf("%s", <-interrupt)
	}()

	select {
	case <-errs:
		close(errs)
		close(interrupt)
		log.Println("\nGracefully shutting down...")
		time.Sleep(2 * time.Second)
	}
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
