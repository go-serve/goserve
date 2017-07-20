package server

import (
	"html/template"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-midway/midway"
	"github.com/go-serve/goserve/assets"
)

var tplVideo *template.Template

func init() {

	fs := assets.FileSystem()
	fh, err := fs.Open("/html/video.html")
	if err != nil {
		log.Print("Failed to load template")
		panic(err)
	}

	b, err := ioutil.ReadAll(fh)
	if err != nil {
		log.Print("Failed to read template file")
		panic(err)
	}

	tplVideo = template.New("video.html")
	tplVideo = tplVideo.Funcs(template.FuncMap{
		"joinPath": func(parts ...string) string {
			return path.Join(parts...)
		},
	})
	tplVideo, err = tplVideo.Parse(string(b))
	if err != nil {
		log.Print("Failed to parse video.html into template")
		panic(err)
	}

	// side-effect: add extension to mime types
	mime.AddExtensionType(".vtt", "text/vtt")
	mime.AddExtensionType(".srt", "text/vtt")
}

// ServeVideo displays HTML5 compatible video files with proper HTML player page
func ServeVideo(root http.FileSystem) midway.Middleware {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if r.URL.Query().Get("mode") == "videoplayer" {

				log.Printf("open videoplayer!")

				file, err := root.Open(r.URL.Path)
				if err != nil {
					// TODO: handle not found / other error
					return
				}
				log.Printf("open videoplayer checkpoint 2")

				if stat, err := file.Stat(); err != nil {
					// TODO: handle not found / other error
					return
				} else if !stat.Mode().IsRegular() {
					// TODO: not a file
					return
				}
				log.Printf("open videoplayer checkpoint 3")

				// TODO: detect if the file is an mp4 / ogg / ogv / vp8 / vp9
				// find if there is srt / webvtt file in the same folder
				subtitles := make([]map[string]string, 0, 1)
				fileBasename := strings.TrimSuffix(r.URL.Path, filepath.Ext(r.URL.Path))
				log.Printf("fileBasename: %s", fileBasename)

				if vtt, err := root.Open(fileBasename + ".vtt"); err == nil {
					if stat, err := vtt.Stat(); err == nil && !stat.Mode().IsDir() {
						subtitles = append(subtitles, map[string]string{
							"Path":     fileBasename + ".vtt",
							"Language": "en",
							"Label":    "English",
						})
					}
					vtt.Close()
				}
				if srt, err := root.Open(fileBasename + ".srt"); err == nil {
					if stat, err := srt.Stat(); err == nil && !stat.Mode().IsDir() {
						subtitles = append(subtitles, map[string]string{
							"Path":     fileBasename + ".srt",
							"Language": "en",
							"Label":    "English",
						})
					}
					srt.Close()
				}

				// display HTML5 video page with track definition for srt / webvtt
				w.Header().Add("Content-Type", "text/html; charset=utf-8")
				err = tplVideo.Execute(w, map[string]interface{}{
					"Name":        r.URL.Path,
					"Path":        r.URL.Path,
					"MetaType":    "video/mp4",
					"Stylesheets": stylesheets,
					"Scripts":     scripts,
					"Subtitles":   subtitles,
				})
				if err != nil {
					log.Printf("error executing template video.html: %s", err.Error())
				}
				return
			}

			// defers to inner handler
			inner.ServeHTTP(w, r)
		})
	}
}
