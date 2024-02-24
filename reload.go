package livereload

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/fsnotify/fsnotify"
)

//go:embed livereload.js
var liveReloadScript []byte

// HandleLiveReload adds the livereload script and handler to the given mux
// to make use of it you need to add a script with the src="/livereload" your html file
// see the LiveReloadScriptHTML func
func HandleLiveReload(mux *http.ServeMux, dirsToWatch ...string) {
	mux.HandleFunc("GET /livereload.js", serveLiveReloadScript)
	mux.HandleFunc("GET /livereload", liveReloadHandler(dirsToWatch...))
}

func LiveReloadScriptHTML() template.HTML {
	return `<script src="/livereload.js" type="text/javascript"></script>`
}

func serveLiveReloadScript(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	w.WriteHeader(http.StatusOK)
	w.Write(liveReloadScript)
}

func liveReloadHandler(dirsToWatch ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			panic(err)
		}
		defer watcher.Close()

		for _, dir := range dirsToWatch {
			err = watcher.Add(dir)
			if err != nil {
				panic(err)
			}
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.WriteHeader(http.StatusOK)

		type fileUpdate struct {
			FileName string
		}

		for event := range watcher.Events {
			update, err := json.Marshal(fileUpdate{FileName: event.Name})
			if err != nil {
				panic(err)
			}
			data := append([]byte("data:"), update...)
			data = append(data, []byte("\n\n")...)

			_, err = w.Write(data)
			if err != nil {
				fmt.Println(err, "closing connection")
				return
			}
			w.(http.Flusher).Flush()
		}
	}
}
