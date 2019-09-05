package server

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
)

func setLoadedToursCookie(w http.ResponseWriter, value []string) {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		log.Panic(err)
	}

	c := http.Cookie{}
	c.Name = "loadedTours"
	c.Value	= base64.StdEncoding.EncodeToString(jsonValue)
	c.Path = "/"

	http.SetCookie(w, &c)
}
func getLoadedToursCookie(r *http.Request) ([]string, error) {
	var loadedToursID []string
	c, err := r.Cookie("loadedTours")
	if err != nil {
		return nil, err
	}
	jsonValue, err :=  base64.StdEncoding.DecodeString(c.Value)
	err = json.Unmarshal(jsonValue, &loadedToursID)
	return loadedToursID, err
}
