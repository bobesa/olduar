package olduar

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	VERSION = "0.0.4"
)

type ServerConfig struct {
	Port string `json:"port"`
	Name string `json:"name"`

	DirSave string `json:"directorySave"`

	DirItems      string `json:"directoryItems"`
	DirLocations  string `json:"directoryLocations"`
	DirCharacters string `json:"directoryCharacters"`
	DirOther      string `json:"directoryOther"`
}

var MainServerMux *http.ServeMux
var MainServerInstance *http.Server
var MainServerConfig *ServerConfig = &ServerConfig{}

func UsernameCheck(username string) bool {
	return strings.IndexAny(username, " :") == -1
}

func Run(configFilename string) {
	fmt.Println("OLDUAR Server " + VERSION + "\n")

	//Loading config file
	data, err := ioutil.ReadFile(configFilename)
	if err != nil {
		fmt.Println("Unable to open \"" + configFilename + "\" config file")
		return
	}
	err = json.Unmarshal(data, &MainServerConfig)

	if err != nil {
		fmt.Println("Error in \"" + configFilename + "\" config file")
		return
	}

	//Randomize generation
	rand.Seed(time.Now().Unix())

	//Loading of files etc.
	if LoadLocations() && LoadItems() && LoadCharacters() {
		//Initialize actions
		LoadAttributes()

		//Prepare folders if not existing
		os.MkdirAll(MainServerConfig.DirSave+"/players", 0777)
		os.MkdirAll(MainServerConfig.DirSave+"/rooms", 0777)

		//Load all players
		LoadAllPlayers()

		//Prepared message .)
		fmt.Println("Everything is prepared, \"" + MainServerConfig.Name + "\" is running")
		fmt.Println()

		//Prepare server
		MainServerMux = http.NewServeMux()
		MainServerInstance = &http.Server{
			Addr:           ":" + MainServerConfig.Port,
			Handler:        MainServerMux,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		//REST API: Admin
		MainServerMux.HandleFunc("/admin/", func(w http.ResponseWriter, r *http.Request) {
			params := strings.Split(r.URL.Path[7:], "/")
			paramLen := len(params)
			if paramLen == 0 {
				http.NotFound(w, r)
				return
			}

			switch params[0] {
			case "give":
				if paramLen == 3 {
					player := ActivePlayersByUsername[params[1]]
					if player != nil {
						player.Give(params[2])
					}
				}
			}
		})

		//REST API: Player
		MainServerMux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
			params := strings.Split(r.URL.Path[5:], "/")
			paramLen := len(params)
			if paramLen == 0 {
				http.NotFound(w, r)
				return
			}

			//If no player found or player don`t have any room - return 404
			player, found := PlayerByAuthorization(r)
			if params[0] == "register" {
				if found {
					//Authorization is valid - register is false
					w.Write([]byte("false"))
				} else {
					authData, found := r.Header["Authorization"]
					if !found || len(authData) != 1 {
						http.NotFound(w, r)
						return
					}
					//Check BASE64 auth for data
					authInfo, err := base64.StdEncoding.DecodeString(strings.Replace(authData[0], "Basic ", "", 1))
					if err == nil {
						authData = strings.Split(string(authInfo), ":")
						if len(authData) == 2 && UsernameCheck(authData[0]) {
							player, found = ActivePlayersByUsername[strings.ToLower(authData[0])]
							if !found {
								//Register user
								player = &Player{Username: authData[0], Password: authData[1], Name: authData[0], Health: 50, MaxHealth: 50}
								player.Activate()
								w.Write([]byte("true"))
							} else {
								//Username used already
								w.Write([]byte("false"))
							}
						}
					}
					if player == nil {
						w.Write([]byte("false"))
					}
				}
				return
			}
			if player == nil {
				http.NotFound(w, r)
				return
			}

			//Process command
			w.Header().Set("Content-Type", "application/json")
			switch params[0] {
			case "save", "look", "do", "go", "inventory", "inspect", "pickup", "drop", "use", "equip", "stats", "attack", "defend", "ability":
				if player.Room == nil {
					w.Write([]byte("null"))
					return
				}
				resp := make(chan []byte)
				command := Command{Player: player, Command: strings.ToLower(params[0]), Response: resp}
				if paramLen > 1 {
					command.Parameter = strings.ToLower(params[1])
				}
				player.Room.queue <- &command
				w.Write(<-resp)

			case "players":
				playerList := make([]string, len(ActivePlayersByUsername))
				playerId := 0
				for username := range ActivePlayersByUsername {
					playerList[playerId] = username
					playerId++
				}
				data, err := json.Marshal(playerList)
				if err == nil {
					w.Write(data)
				} else {
					w.Write([]byte("[]"))
				}

			case "tell":
				if paramLen == 2 && params[1] != "" {
					target, found := ActivePlayersByUsername[strings.ToLower(params[1])]
					if found && target.Room != nil {
						defer r.Body.Close()
						messageData, err := ioutil.ReadAll(r.Body)
						message := string(messageData)
						if err != nil || message == "" {
							w.Write([]byte("false"))
						} else {
							target.Tell(player.Name + " (" + player.Username + "): " + message)
							w.Write([]byte("true"))
						}
					} else {
						w.Write([]byte("false"))
					}
				} else {
					w.Write([]byte("false"))
				}

			case "say":
				if player.Room != nil {
					defer r.Body.Close()
					messageData, err := ioutil.ReadAll(r.Body)
					message := string(messageData)
					if err != nil || message == "" {
						w.Write([]byte("false"))
					} else {
						player.Room.TellAll(player.Name + ": " + message)
						w.Write([]byte("true"))
					}
				} else {
					w.Write([]byte("false"))
				}

			case "party":
				if player.Room != nil {
					playerList := make([]string, len(player.Room.Players))
					for index, player := range player.Room.Players {
						playerList[index] = player.Username
					}
					data, err := json.Marshal(playerList)
					if err == nil {
						w.Write(data)
					} else {
						w.Write([]byte("[]"))
					}
				} else {
					w.Write([]byte("[]"))
				}

			case "rename":

				defer r.Body.Close()
				nameData, err := ioutil.ReadAll(r.Body)
				name := strings.Trim(string(nameData), "!?,.-= ")

				if err == nil && name != "" {
					w.Write([]byte("true"))
					player.Name = name
					if player.Room == nil {
						player.Save()
					}
				}

				if name == "" {
					w.Write([]byte("false"))
				}

			case "rooms":
				data, err := json.Marshal(GetRoomList())
				if err == nil {
					w.Write(data)
				} else {
					w.Write([]byte("[]"))
				}

			case "leave":
				if player.Room != nil {
					player.Room.Leave(player)
				}
				data, err := json.Marshal(GetRoomList())
				if err == nil {
					w.Write(data)
				} else {
					w.Write([]byte("[]"))
				}

			case "join":
				if paramLen > 1 && params[1] != "" {
					room, found := AllRooms[params[1]]
					if !found {
						room = CreateRoomWithName(params[1])
					}
					room.Join(player)
					resp := make(chan []byte)
					room.queue <- &Command{Player: player, Command: "look", Response: resp}
					w.Write(<-resp)

				} else {
					w.Write([]byte("null"))
				}
			}
		})

		//Start web playable client
		MainServerMux.Handle("/client/", http.StripPrefix("/client/", http.FileServer(http.Dir("../client/html/"))))

		//Start server
		MainServerInstance.ListenAndServe()
	}

}
