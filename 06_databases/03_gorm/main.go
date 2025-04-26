package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"
)

type Item struct {
	Id          int `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Title       string
	Description string
	Updated     string `sql:"null"`
}

func (i *Item) TableName() string {
	return "items"
}

func (i *Item) BeforeSave() (err error) {
	fmt.Println("trigger on before save")
	return
}

// see also https://github.com/doug-martin/goqu
type Handler struct {
	DB   *gorm.DB
	Tmpl *template.Template
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	items := []*Item{}

	db := h.DB.Find(&items)
	err := db.Error
	panicOnErr(err)

	err = h.Tmpl.ExecuteTemplate(w, "index.html", struct {
		Items []*Item
	}{
		Items: items,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) AddForm(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, "create.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
	newItem := &Item{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
	}
	db := h.DB.Create(&newItem)
	err := db.Error
	panicOnErr(err)
	affected := db.RowsAffected

	fmt.Println("Insert - RowsAffected", affected, "LastInsertId: ", newItem.Id)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) Edit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	panicOnErr(err)

	post := &Item{
		Id: id,
	}

	db := h.DB.Find(post)
	err = db.Error
	if err == gorm.ErrRecordNotFound {
		fmt.Println("Record not found", id)
	} else {
		panicOnErr(err)
	}

	err = h.Tmpl.ExecuteTemplate(w, "edit.html", post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	panicOnErr(err)

	post := &Item{}
	h.DB.Find(post, id)

	post.Title = r.FormValue("title")
	post.Description = r.FormValue("description")
	post.Updated = "rvasily"

	db := h.DB.Save(post)
	err = db.Error
	panicOnErr(err)
	affected := db.RowsAffected

	fmt.Println("Update - RowsAffected", affected)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	panicOnErr(err)

	db := h.DB.Delete(&Item{Id: id})
	err = db.Error
	panicOnErr(err)
	affected := db.RowsAffected

	fmt.Println("Delete - RowsAffected", affected)

	w.Header().Set("Content-type", "application/json")
	resp := `{"affected": ` + strconv.Itoa(int(affected)) + `}`
	w.Write([]byte(resp))
}

func main() {

	// основные настройки подключения к базе
	dsn := "root:love@tcp(localhost:3306)/golang?"
	// указываем кодировку
	dsn += "&charset=utf8"
	// отказываемся от prepared statements
	// параметры подставляются сразу
	dsn += "&interpolateParams=true"

	db, err := gorm.Open("mysql", dsn)
	db.DB().Ping()
	panicOnErr(err)

	handlers := &Handler{
		DB:   db,
		Tmpl: template.Must(template.ParseGlob("templates/*")),
	}

	// В целях упрощения примера пропущена авторизация и csrf
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.List).Methods("GET")
	r.HandleFunc("/items", handlers.List).Methods("GET")
	r.HandleFunc("/items/new", handlers.AddForm).Methods("GET")
	r.HandleFunc("/items/new", handlers.Add).Methods("POST")
	r.HandleFunc("/items/{id}", handlers.Edit).Methods("GET")
	r.HandleFunc("/items/{id}", handlers.Update).Methods("POST")
	r.HandleFunc("/items/{id}", handlers.Delete).Methods("DELETE")

	fmt.Println("starting server at :8080")
	fmt.Println(http.ListenAndServe(":8080", r))
}

// Не используйте такой код в продакшене
// ошибка должна всегда явно обрабатываться
func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
