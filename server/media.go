package server

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-midway/midway"
	"github.com/go-serve/goserve/assets"
)

var tplVideo *template.Template
var srtDateReg *regexp.Regexp

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
	mime.AddExtensionType(".srt", "text/srt")

	// format of an SRT date line
	srtDateReg = regexp.MustCompile("(\\d{2}:\\d{2}:\\d{2}),(\\d{3}) --> (\\d{2}:\\d{2}:\\d{2}),(\\d{3})")
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
							"Path":     fileBasename + ".srt?mode=vtt",
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

// ServeSrt serves translates srt files to webvtt and write
// to browser on-the-go
func ServeSrt(root http.FileSystem) midway.Middleware {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if strings.ToLower(path.Ext(r.URL.Path)) == ".srt" {
				if r.URL.Query().Get("mode") == "vtt" {
					// TODO: check file extension, should be srt

					// Open SRT and wrap with SrtWebvttReader
					srt, err := root.Open(r.URL.Path)
					if err != nil {
						// defer to inner handler
						inner.ServeHTTP(w, r)
						return
					}
					defer srt.Close()

					// mask the SRT with masking reader
					r, err := NewSrtWebvttReader(srt)
					if err != nil {
						// TODO: handler error
						return
					}

					// use io.Copy to pipe out the masked reader
					w.Header().Set("Content-Type", "text/vtt; charset=utf-8")
					w.WriteHeader(http.StatusOK)
					io.Copy(w, r)
					return
				}
			}

			// use inner handler
			inner.ServeHTTP(w, r)
			return
		})
	}
}

// SrtWebvttReader masks inner reader stream of supposed
// SRT files into WEBVTT stream reader
type SrtWebvttReader struct {
	src  io.Reader
	buff []byte
	p    int
}

// Read implements io.Reader
func (r *SrtWebvttReader) Read(b []byte) (n int, err error) {

	// first read, prefix to WEBVTT
	if r.p == 0 {

		// add prefix text
		var prefN, buffN int
		prefN = 8

		// create buffer for first read
		buff := make([]byte, len(b)-8)
		buffN, err = r.src.Read(buff)
		if err != nil {
			return
		}

		// TODO; handle date conversion in buff
		buff = srtDateReg.ReplaceAll(buff, []byte("$1.$2 --> $3.$4"))

		// copy buffered first read to output byte slice
		copy(b, append([]byte("WEBVTT\n\n"), buff...))
		n = prefN + buffN
		r.p = n
		return
	}

	// resize buffer to fit reader
	if len(r.buff) != len(b) {
		r.buff = make([]byte, len(b))
	}

	// read to buffer
	n, err = r.src.Read(r.buff)
	r.p += n

	// handle date conversion in buff
	buff := srtDateReg.ReplaceAll(r.buff, []byte("$1.$2 --> $3.$4"))

	// copy to output
	copy(b, buff)
	return
}

// NewSrtWebvttReader converts
func NewSrtWebvttReader(inner io.Reader) (r io.Reader, err error) {
	if inner == nil {
		err = fmt.Errorf("inner reader is empty")
		return
	}
	r = &SrtWebvttReader{
		src:  inner,
		buff: make([]byte, 1024),
		p:    0,
	}
	return
}
