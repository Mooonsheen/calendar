package application

import (
	"context"
	"errors"
	"html/template"
	"net/http"
	"net/url"
	"path/filepath"
	"server/internal/repository"
	"strconv"
	"time"

	//"crypto/sha256"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/julienschmidt/httprouter"
)

//var flag bool

type app struct {
	ctx   context.Context
	repo  *repository.Repository
	cache map[string]repository.User
}

func (a app) Routes(r *httprouter.Router) {
	r.ServeFiles("/public/*filepath", http.Dir("public"))
	r.GET("/", a.authorized(a.Home))
	r.GET("/login", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		a.LoginPage(w, "")
	})
	r.GET("/sign_in", a.Login)
	r.POST("/login", a.Registration)

}

func (a app) Home(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	lp := filepath.Join("public", "html", "calendar.html")
	common := filepath.Join("public", "html", "common_calendar.html")
	tmpl, err := template.ParseFiles(lp, common)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = tmpl.ExecuteTemplate(w, "calendar", 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (a app) LoginPage(rw http.ResponseWriter, message string) {
	lp := filepath.Join("public", "html", "login.html")
	common := filepath.Join("public", "html", "common_login.html")
	tmpl, err := template.ParseFiles(lp, common)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	type answer struct {
		Message string
	}
	data := answer{message}
	err = tmpl.ExecuteTemplate(rw, "login", data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}

func (a app) Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	email := r.FormValue("email")
	passwd := r.FormValue("psw")

	if email == "" || passwd == "" {
		a.LoginPage(w, "Необходимо указать логин и пароль")
		return
	}

	_, err := a.repo.Login(a.ctx, email, passwd)
	if err != nil {
		a.LoginPage(w, "Вы ввели неверный логин или пароль")
		return
	}

	time64 := time.Now().Unix()
	timeInt := strconv.Itoa(int(time64))
	token := email + passwd + timeInt

	// hashToken := md5.Sum([]byte(token))
	// hashedToken := hex.EncodeToString(hashToken[:])

	//a.cache[hashedToken] = user
	livingTime := 60 * time.Minute
	expiration := time.Now().Add(livingTime)
	//кука будет жить 1 час
	cookie := http.Cookie{Name: "token", Value: url.QueryEscape(token), Expires: expiration}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a app) Registration(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	email := r.FormValue("email")
	passwd := r.FormValue("psw")
	reppas := r.FormValue("psw-repeat")

	if email == "" || passwd == "" {
		a.LoginPage(w, "Необходимо указать логин и пароль!")
		return
	}

	if passwd != reppas {
		a.LoginPage(w, "Пароли не совпадают")
		return
	}
	_, err := a.repo.AddNewUser(a.ctx, email, passwd)
	if err != nil {
		a.LoginPage(w, "Аккаунт уже существует")
		return
	}

	a.LoginPage(w, "Аккаунт создан, можете войти в него")
	//http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a app) authorized(next httprouter.Handle) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		token, err := readCookie("token", r)
		if err != nil {
			http.Redirect(rw, r, "/login", http.StatusSeeOther)
			return
		}
		if _, ok := a.cache[token]; !ok {
			next(rw, r, ps)
			return
		} else {
			http.Redirect(rw, r, "/login", http.StatusSeeOther)
			return
		}
	}
}

func readCookie(name string, r *http.Request) (value string, err error) {
	if name == "" {
		return value, errors.New("you are trying to read empty cookie")
	}
	cookie, err := r.Cookie(name)
	if err != nil {
		return value, err
	}
	str := cookie.Value
	value, _ = url.QueryUnescape(str)
	return value, err
}

func NewApp(ctx context.Context, dbpool *pgxpool.Pool) *app {
	return &app{ctx, repository.NewRepository(dbpool), make(map[string]repository.User)}
}
