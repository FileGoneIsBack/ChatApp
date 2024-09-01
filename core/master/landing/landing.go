package landing

import (
    "login/core/models"
    "login/core/models/functions"
    "login/core/models/server"
    "net/http"
    "sync"
)

var (
    countMu sync.Mutex
    count   int
)

func init() {
    Route.NewSub(server.NewRoute("/", func(w http.ResponseWriter, r *http.Request) {
        incrementCount()
        type Page struct {
            Name, Title string
            Count       int
        }
        functions.Render(Page{
            Name:  models.Config.Name,
            Title: "Welcome!",
            Count: getCount() + 1000,
        }, w, "index.html", server.TemplateCache)
    }))
}

func incrementCount() {
    countMu.Lock()
    defer countMu.Unlock()

    count++
}

func getCount() int {
    countMu.Lock()
    defer countMu.Unlock()
    return count
}
