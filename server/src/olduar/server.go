package olduar

import (
	"fmt"
	"math/rand"
	"time"
	"net/http"
)

const (
	VERSION = "0.01a"
)

var MainServerMux *http.ServeMux
var MainServerInstance *http.Server

func Run(databaseDirectory string) {
	fmt.Println("OLDUAR Server "+VERSION+"\n")

	rand.Seed(time.Now().Unix()) //Randomize generation

	if(LoadLocations(databaseDirectory+"/locations") && LoadItems(databaseDirectory+"/items")) {
		fmt.Println("Everything is prepared!\n")

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
