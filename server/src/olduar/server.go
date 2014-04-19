package olduar

import (
	"fmt"
	"math/rand"
	"time"
	"net/http"
	"strings"
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
		LoadAllPlayers("./save/players/")

		fmt.Println("Everything is prepared, \""+config.Name+"\" is running")
		fmt.Println()

		//Prepare server
		MainServerMux = http.NewServeMux()
		MainServerInstance = &http.Server{
			Addr:           ":"+config.Port,
			Handler:        MainServerMux,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		//REST Api
		apiPath := "/api/"
		apiPathLen := len(apiPath)
		MainServerMux.HandleFunc(apiPath, func(w http.ResponseWriter, r *http.Request){
				//If no player found or player don`t have any room - return 404
				player, found := PlayerByAuthorization(r)
				if(!found) {
					http.NotFound(w,r)
					return
				}

				//Process command
				params := strings.Split(r.URL.Path[apiPathLen:],"/")
				paramLen := len(params)
				if(paramLen>0) {
					w.Header().Set("Content-Type", "application/json")
					switch(params[0]){
					case "room":
						if(player.Room == nil || paramLen < 2) {
							http.NotFound(w,r)
							return
						}
						resp := make(chan []byte)
						command := Command{Player:player, Command: strings.ToLower(params[1]), Response: resp}
						if(paramLen>2) {
							command.Parameter = strings.ToLower(params[2])
						}
						player.Room.queue <- &command
						data := <- resp
						w.Write(data)
					case "rooms":
						data, err := json.Marshal(GetRoomList())
						if(err == nil) {
							w.Write(data)
						} else {
							w.Write([]byte("[]"))
						}
					case "leave":
						if(player.Room != nil) {
							player.Room.Leave(player)
						}
						data, err := json.Marshal(GetRoomList())
						if(err == nil) {
							w.Write(data)
						} else {
							w.Write([]byte("[]"))
						}

					case "join":
						if(paramLen>1) {
							room, found := AllRooms[params[1]]
							if(!found) {
								room = CreateRoomWithName(params[1])
							}
							room.Join(player)
							resp := make(chan []byte)
							room.queue <- &Command{Player:player, Command: "look", Response: resp}
							data := <- resp
							w.Write(data)

						} else {
							w.Write([]byte("null"))
						}
					}
				} else {
					http.NotFound(w,r)
				}
			})

		//Start web playable client
		MainServerMux.Handle("/client/", http.StripPrefix("/client/", http.FileServer(http.Dir("../client/html/"))))

		//Start server
		MainServerInstance.ListenAndServe()
	}

}
