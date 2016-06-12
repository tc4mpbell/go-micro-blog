package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
    "errors"
)

/* GLOBAL */
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

func getTitle(w http.ResponseWriter, r *http.Request) (string, error)  {
    m := validPath.FindStringSubmatch(r.URL.Path)
    if m == nil {
        http.NotFound(w, r)
        return  "", errors.New("Invalid Page Title")
    }
    return m[2], nil
}

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
        // extract page title from the request and call the passed handler 'fn'
        m := validPath.FindStringSubmatch(r.URL.Path)
        if m == nil {
            http.NotFound(w, r)
            return
        }
        fn(w, r, m[2])
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

func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.ListenAndServe(":8080", nil)
}


