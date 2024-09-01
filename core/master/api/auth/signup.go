package authentication

import (
    "login/core/database"
    "login/core/models/functions"
    "login/core/sessions"
    "login/core/models/server"
    "net/http"
    "login/core/models"
    "html/template"
    "github.com/google/uuid"
    "time"
)

func init() {
    Route.NewSub(server.NewRoute("/signup", func(w http.ResponseWriter, r *http.Request) {
        type Page struct {
            Name   string
            Title  string
            Script template.HTML
        }
        if r.Method == http.MethodPost {
            // Parse form data
            err := r.ParseForm()
            if err != nil {
                functions.Render(Page{
                    Name:  models.Config.Name,
                    Title: "Signup",
                    Script: template.HTML(functions.Toast(functions.Toastr{
                        Icon:  "error",
                        Title: "Error!",
                        Text:  "Form must be filled!",
                    })),
                }, w, "auth/auth.html", server.TemplateCache)
                return
            }

            username := r.Form.Get("username")
            password := r.Form.Get("password")
            password2 := r.Form.Get("password2")
            email := r.Form.Get("email")

            //Validation
            // Validate username length
            if len(username) < 4 {
                functions.Render(Page{
                    Name:  models.Config.Name,
                    Title: "Signup",
                    Script: template.HTML(functions.Toast(functions.Toastr{
                        Icon:  "error",
                        Title: "Error!",
                        Text:  "Username must be longer then 4 characters!",
                    })),
                }, w, "auth/auth.html", server.TemplateCache)
                return
            }

            // Validate password length
            if len(password) < 6 {
                functions.Render(Page{
                    Name:  models.Config.Name,
                    Title: "Signup",
                    Script: template.HTML(functions.Toast(functions.Toastr{
                        Icon:  "error",
                        Title: "Error!",
                        Text:  "Password must be longer then 6 characters!",
                    })),
                }, w, "auth/auth.html", server.TemplateCache)
                return
            }
            // Validate password match
            if password != password2 {
                functions.Render(Page{
                    Name:  models.Config.Name,
                    Title: "Signup",
                    Script: template.HTML(functions.Toast(functions.Toastr{
                        Icon:  "error",
                        Title: "Error!",
                        Text:  "Passwords do not match!",
                    })),
                }, w, "auth/auth.html", server.TemplateCache)
                return
            }

            // Check if username already exists
            existingUsers, err := database.Container.GetUsers()
            if err != nil {
                functions.Render(Page{
                    Name:  models.Config.Name,
                    Title: "Signup",
                    Script: template.HTML(functions.Toast(functions.Toastr{
                        Icon:  "error",
                        Title: "Error!",
                        Text:  "User already exists!",
                    })),
                }, w, "auth/auth.html", server.TemplateCache)
                return
            }

            for _, user := range existingUsers {
                if user.Username == username {
                    functions.Render(Page{
                        Name:  models.Config.Name,
                        Title: "Signup",
                        Script: template.HTML(functions.Toast(functions.Toastr{
                            Icon:  "error",
                            Title: "Error!",
                            Text:  "Username Taken!",
                        })),
                    }, w, "auth/auth.html", server.TemplateCache)
                    return
                }
            }

            //end of validation
            // Create new user
            newUser := &database.User{
                Username: username,
                Password: password,
                Email:    email,
            }

            // Call function to create new user in the database
            err = database.Container.NewUser(newUser)
            if err != nil {
                functions.Render(Page{
                    Name:  models.Config.Name,
                    Title: "Signup",
                    Script: template.HTML(functions.Toast(functions.Toastr{
                        Icon:  "error",
                        Title: "Error!",
                        Text:  "Failed to create user, Likely database issuse!",
                    })),
                }, w, "auth/auth.html", server.TemplateCache)
                return
            }
             // Set session token cookie
            expiresAt := time.Now().Add(30 * time.Minute)
            sessionToken := uuid.NewString()
            sessions.Sessions[sessionToken] = sessions.Session{
                User:   newUser,
                Expiry: expiresAt,
            }
            cookie := &http.Cookie{
                Name:    "session-token",
                Value:   sessionToken,
                Expires: expiresAt,
                Path:    "/",
            }
            http.SetCookie(w, cookie)
            http.Redirect(w, r, "/dash/", http.StatusSeeOther)
        } else {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    }))
}
