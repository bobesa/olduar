package olduar

import (
	"fmt"
	"math/rand"
	"time"
	"net/http"
	"encoding/json"
	"io/ioutil"
)

const (
	VERSION = "0.01a"
)

type ServerConfig struct {
	Port string				`json:"port"`
	Name string				`json:"name"`

	DirItems string			`json:"directory_items"`
	DirLocations string		`json:"directory_locations"`
	DirAttributes string	`json:"directory_attributes"`
}

var MainServerMux *http.ServeMux
var MainServerInstance *http.Server

func Run(configFilename string) {
	fmt.Println("OLDUAR Server "+VERSION+"\n")

	//Loading config file
	data, err := ioutil.ReadFile(configFilename);
	if(err != nil) {
		fmt.Println("Unable to open \""+configFilename+"\" config file")
		return
	}
	config := ServerConfig{}
	err = json.Unmarshal(data,&config)

	if(err != nil) {
		fmt.Println("Error in \""+configFilename+"\" config file")
		return
	}

	//Randomize generation
	rand.Seed(time.Now().Unix())

	//Loading of files etc.
	if(LoadLocations(config.DirLocations) && LoadItems(config.DirItems)) {
		InitializeActions()
		fmt.Println("Everything is prepared, \""+config.Name+"\" is running")
		fmt.Println()

		//Prepare server
		MainServerMux = http.NewServeMux()
		MainServerInstance = &http.Server{
			Addr:           ":8080",
			Handler:        MainServerMux,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		//TODO: Remove this (testing)
		//Prepare test player & game
		player := &Player{Name:"Test",Username:"test",HashPass:"test"}
		game := CreateGameStateFromName("test")
		game.Join(player)

		//Start server
		MainServerInstance.ListenAndServe()
	}

}
