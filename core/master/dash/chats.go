package dash

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
    "time"
	"login/core/database"
	authentication "login/core/master/api/auth"
    "login/core/models/validation"
	"login/core/models"
	"login/core/models/functions"
	"login/core/models/server"
	"login/core/sessions"
)

func init() {
    Route.NewSub(server.NewRoute("/chats", func(w http.ResponseWriter, r *http.Request) {
        ok, user := sessions.IsLoggedIn(w, r)
        if !ok {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        avatarPath, _ := database.Container.GetAvatar(user.ID)
        personIDStr := r.URL.Query().Get("userID")
        if personIDStr == "" {
            http.Error(w, "userID not provided", http.StatusBadRequest)
            return
        }

        personID, err := strconv.Atoi(personIDStr)
        if err != nil {
            http.Error(w, "Invalid userID", http.StatusBadRequest)
            return
        }

        chats, _, err := database.Container.GetChatAndMessages(user.ID, personID)
        if err != nil {
            http.Error(w, "Error fetching chats", http.StatusInternalServerError)
            return
        }
        person, _ := database.Container.GetUserByID(personID)
        if chats == nil {
            Sanitize.LogMessage("error", "Chat is nil", nil)
            http.Error(w, "Chat not found", http.StatusNotFound)
            return
        }

        chatsID := strconv.Itoa(chats.ID)

        chatIDCookie := http.Cookie{
            Name:     "chatID",
            Value:    chatsID,
            Path:     "/",
            HttpOnly: true,
            Secure:   r.TLS != nil,
        }
        http.SetCookie(w, &chatIDCookie)
        switch strings.ToLower(r.Method) {
        case "get":
            type Page struct {
                Name      string
                Avatar    string
                Title     string
                *sessions.Session
                Time, PersonTime      string
                Script    template.JS
                Person   *database.User
                Chat     *database.Chat
            }
            formattedTime := time.Now().Format("03:04PM") + " " + person.TMZ
            formattedTime2 := time.Now().Format("03:04PM") + " " + user.TMZ
            functions.Render(Page{
                Name:     models.Config.Name,
                Title:   "Auth",
                Session:  user,
				Avatar:	 "/_assets/media/avatar/"+avatarPath,
                Person:   person,
                PersonTime: formattedTime,
                Time:     formattedTime2,
                Chat:     chats,
            }, w, "dash/chats.html", server.TemplateCache)
        case "post":
            authentication.Signin(w, r)
        }
    }))
}
