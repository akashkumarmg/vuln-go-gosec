package main

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Path traversal vulnerability
		file := r.URL.Path
		data, err := ioutil.ReadFile(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(data)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		// SQL injection vulnerability
		username := r.FormValue("username")
		password := r.FormValue("password")
		db, err := sql.Open("sqlite3", "./example.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		query := fmt.Sprintf("SELECT * FROM users WHERE username='%s' AND password='%s'", username, password)
		rows, err := db.Query(query)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// Weak password storage
		hash := md5.Sum([]byte(password))
		fmt.Printf("MD5 hash: %x\n", hash)
	})

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		// Unrestricted file upload vulnerability
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		filename := header.Filename
		data, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = ioutil.WriteFile(filename, data, 0644)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/execute", func(w http.ResponseWriter, r *http.Request) {
		// Command injection vulnerability
		cmd := r.FormValue("cmd")
		output, err := exec.Command("sh", "-c", cmd).CombinedOutput()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(output)
	})

	http.HandleFunc("/traverse", func(w http.ResponseWriter, r *http.Request) {
		// Directory traversal vulnerability
		dir := r.FormValue("dir")
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, file := range files {
			fmt.Fprintln(w, file.Name())
		}
	})

	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		// Arbitrary file deletion vulnerability
		file := r.FormValue("file")
		err := os.RemoveAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		// Sensitive information disclosure vulnerability
		dir := "/etc"
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, file := range files {
			fmt.Fprintln(w, file.Name())
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
