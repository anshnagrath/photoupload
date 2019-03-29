package main

import (
	"crypto/sha1"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	uuid "github.com/satori/go.uuid"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("*.html"))
}
func main() {
	http.HandleFunc("/", userData)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("./public"))))
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
func userData(w http.ResponseWriter, req *http.Request) {
	c := getUserCookie(w, req)
	// fmt.Println(c.Value)
	if req.Method == http.MethodPost {
		mf, mh, err := req.FormFile("image")
		if err != nil {
			log.Panic(err)
		}
		defer mf.Close()
		ext := strings.Split(mh.Filename, ".")[1]
		h := sha1.New()
		io.Copy(h, mf)
		fname := fmt.Sprintf("%x", h.Sum(nil)) + "." + ext
		wd, err := os.Getwd()
		if err != nil {
			log.Panic(err)
		}
		path := filepath.Join(wd, "public", "images", fname)
		nf, err := os.Create(path)
		if err != nil {
			log.Panic(err)
		}
		defer nf.Close()
		mf.Seek(0, 0)
		io.Copy(nf, mf)
		c = appendValues(w, fname, c)
		// http.Redirect(w, req, "/public/", 302)
		//file, err := filepath.Glob(wd + "/public/images/*")
		// if err != nil {
		// 	log.Fatalln(err)
		// }
		// fmt.Printf("%T", file)
		// tpl.ExecuteTemplate(w, "photoblog.html", file)

	}
	// if req.Method == http.MethodGet {
	xs := strings.Split(c.Value, "|")
	xs = append(xs[:0], xs[1:]...)
	fmt.Println(xs)
	tpl.ExecuteTemplate(w, "photoblog.html", xs)
	// }
}
func appendValues(w http.ResponseWriter, fname string, c *http.Cookie) *http.Cookie {

	s := c.Value
	if !strings.Contains(s, fname) {
		s += "|" + fname
	}
	c.Value = s
	http.SetCookie(w, c)
	return c

}
func getUserCookie(w http.ResponseWriter, req *http.Request) *http.Cookie {
	c, err := req.Cookie("session")
	if err != nil {
		uid, _ := uuid.NewV4()
		c = &http.Cookie{
			Name:     "session",
			Value:    uid.String(),
			HttpOnly: true,
		}
		http.SetCookie(w, c)

	}
	return c
}

// func appendValues(c *http.Cookie)  {
// 	ck:=c.Value
// 	if !strings.Contains("")

// }
