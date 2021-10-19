package PageHandlers

import (
	"embed"
	"net/http"
)

//go:embed static
var staticEmbedFS embed.FS

func Static(writer http.ResponseWriter, request *http.Request) {
	http.FileServer(http.FS(staticEmbedFS))
}
