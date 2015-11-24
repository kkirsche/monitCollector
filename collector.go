package main

import (
	"fmt"
	"github.com/clbanning/mxj"
	"github.com/paulrosania/go-charset/charset"
	_ "github.com/paulrosania/go-charset/data"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

type Map map[string]interface{}

func processEventQueue(event_queue_basepath string) {
	event_queue_path := fmt.Sprintf("%s/*", event_queue_basepath)
	str_arr, err := filepath.Glob(event_queue_path)
	if err != nil {
		log.Fatal(err)
	}
	for _, path := range str_arr {
		dat, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(dat[110:]))
	}
}

func collect(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if ok {
		_ = fmt.Sprintf("Username: %s | Password: %s", username, password)
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		mxj.XmlCharsetReader = charset.NewReader
		m, err := mxj.NewMapXml(body) // unmarshal
		if err != nil {
			log.Fatal(err)
		}
		// Single path
		path := m.PathForKeyShortest("avg01")
		val, err := m.ValueForPath(path)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(val)
		// Multi-path
		paths := m.PathsForKey("percent")
		for _, path := range paths {
			val, err := m.ValueForPath(path)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(val)
		}
		io.WriteString(w, "Success")
	}
}

func main() {
	http.HandleFunc("/collector", collect)
	http.ListenAndServe(":8085", nil)
}
