package authentication

import (
	"log"
	"login/core/database"
    "login/core/models"
	"login/core/models/functions"
    "login/core/models/server"
	"login/core/sessions"
    "html/template"
	"net/http"
)

func Signin(w http.ResponseWriter, r *http.Request) {
    type Page struct {
		Name   string
		Title  string
		Script template.HTML
	}
    // Parse form data
    err := r.ParseForm()
    if err != nil {
		functions.Render(Page{
			Name:  models.Config.Name,
			Title: "Login",
			Script: template.HTML(functions.Toast(functions.Toastr{
				Icon:  "error",
				Title: "Error!",
				Text:  "All forms must be filled!",
			})),
		}, w, "auth/auth.html", server.TemplateCache)
        return
    }

    username := r.Form.Get("username")
    password := r.Form.Get("password")
    log.Printf("username: %s", username)
    // Authenticate user
    user, err := database.Container.AuthenticateUser(username, password)
    if err != nil {
        if err == database.ErrUserNotFound || err == database.ErrInvalidPassword {
            functions.Render(Page{
                Name:  models.Config.Name,
                Title: "Login",
                Script: template.HTML(functions.Toast(functions.Toastr{
                    Icon:  "error",
                    Title: "Error!",
                    Text:  "Invalid credentials!",
                })),
            }, w, "auth/auth.html", server.TemplateCache)
            return
        }
        return
    }

    cookie := sessions.CreateSession(user)
    http.SetCookie(w, cookie)
    http.Redirect(w, r, "/dash/", http.StatusSeeOther)
}
