package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
    "github.com/tc4mpbell/go-micro-auth"
    "fmt"
)

// Called first 
func SetupPostEndpoints() {
    http.HandleFunc("/view/", makeHandler(viewHandler))
    http.HandleFunc("/edit/", makeHandler(editHandler))
    http.HandleFunc("/save/", makeHandler(saveHandler))
    http.HandleFunc("/login", loginHandler)
    http.HandleFunc("/logout", logoutHandler)
    http.ListenAndServe(":8080", nil)
}

/* GLOBAL */
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9_-]+)$")
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Post) {
	err := templates.ExecuteTemplate(w, tmpl + ".html", p)
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

/* POSTS */

type Post struct {
	Title string
	Body []byte
}

func path(filename string) string {
	return "posts/" + filename + ".txt"
}

func (p* Post) save() error {
	filename := path(p.Title)
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPost(title string) (*Post, error) {
  filename := path(title)
  body, err := ioutil.ReadFile(filename)
  if err != nil {
  	return nil, err
  }
  return &Post{Title: title, Body: body}, nil
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        fmt.Println("Request!")
        if auth.Authenticated("taylor") {
            // extract page title from the request and call the passed handler 'fn'
            m := validPath.FindStringSubmatch(r.URL.Path)
            if m == nil {
                http.NotFound(w, r)
                return
            }

            fn(w, r, m[2])
        } else {
            fmt.Println("Not authenticated!")
            return 
        }
    }
}

/* WEB SERVER */
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPost(title)
    if err != nil {
        http.Redirect(w, r, "/edit/"+title, http.StatusFound)
        return
    }
    
    renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {    
    p, err := loadPost(title)
    if err != nil {
        p = &Post{Title: title}
    }
    renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
    body := r.FormValue("body")
    p := &Post{Title: title, Body: []byte(body)}
    err := p.save()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
}


func loginHandler(w http.ResponseWriter, r *http.Request) {
    username := r.FormValue("username")
    passwd := r.FormValue("password")

    fmt.Printf("Logging in %s, %s", username, passwd)


    auth.Login(username, passwd)
}


func logoutHandler(w http.ResponseWriter, r *http.Request) {
    username := r.FormValue("username")
    fmt.Println("Logging out %s", username)
    auth.Logout(username)
}
