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
)

type Map map[string]interface{}

func collect(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if ok {
		output := fmt.Sprintf("Username: %s | Password: %s", username, password)
		fmt.Println(output)
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
		path := m.PathForKeyShortest("version")
		fmt.Println(m.ValueForPath(path))
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
	http.ListenAndServe(":2812", nil)
}
