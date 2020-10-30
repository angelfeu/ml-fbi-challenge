package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type satelliteType struct {
	Name     string   `json:"name"`
	Position position `json:"position"`
}

type allSatellitesType []satelliteType

var allSatellites allSatellitesType = []satelliteType{
	{
		Name: "Kenobi",
		Position: position{
			X: -500,
			Y: -200,
		},
	},
	{
		Name: "Skywalker",
		Position: position{
			X: 100,
			Y: -100,
		},
	},
	{
		Name: "Sato",
		Position: position{
			X: 500,
			Y: 100,
		},
	},
}

type satellitesReqType struct {
	Satellites []requestType `json:"satellites"`
}

var satellites satellitesReqType

type requestType struct {
	Name     string   `json:"name"`
	Distance float32  `json:"distance"`
	Message  []string `json:"message"`
}

type allRequestType []requestType

var requests allRequestType

type position struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type respUbicationAndMessage struct {
	Position position `json:"position"`
	Message  string   `json:"message"`
}

func GetLocation(distances ...float32) (x, y float32) {
	return 100, -75
}

func GetMessage(messages ...[]string) (message string) {
	//var message string
	//var positionInArray int = GetMax(len(messages[0]), len(messages[0]), len(messages[0]))
	var array1 []string = messages[0]
	var array2 []string = messages[1]
	var array3 []string = messages[2]

	var position2 int
	var position3 int

	// si el len de algún array es mas corto se completa
	// de acuerdo a la conciliación de posicion por una palabra encontrada
	for i := 0; i < len(array1); i++ {
		position2 = containsInPosition(array1[i], array2)
		if i < position2 {
			for j := i; j <= position2; j++ {
				array1 = append(array1, "")
				copy(array1[1:], array1[0:])
				array1[0] = ""
			}
			break
		} else if position2 < i {
			for j := position2; j <= i; j++ {
				array2 = append(array2, "")
				copy(array2[1:], array2[0:])
				array2[0] = ""
			}
			break
		}
	}

	for i := 0; i < len(array1); i++ {
		position3 = containsInPosition(array1[i], array3)
		if i < position3 {
			for j := i; j <= position3; j++ {
				array1 = append(array1, "")
				copy(array1[1:], array1[0:])
				array1[0] = ""
			}
			break
		} else if position3 < i {
			for j := position3; j <= i; j++ {
				array3 = append(array3, "")
				copy(array3[1:], array3[0:])
				array3[0] = ""
			}
			break
		}
	}

	// arma la frase con la primer palabra que encuentra en cada posicion
	for i := 0; i < GetMax(len(array1), len(array2), len(array3)); i++ {
		if message != "" {
			message = message + " "
		}
		message = message + GetValueNonBlank(array1[i], array2[i], array3[i])
	}

	return message
}

func GetSatellites(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/jspn")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(allSatellites)
}

func PostSatellites(w http.ResponseWriter, r *http.Request) {
	var newSatellites satellitesReqType
	err := json.NewDecoder(r.Body).Decode(&newSatellites)
	if err != nil {
		http.Error(w, "Error al leer los datos", http.StatusBadRequest)
		return
	}
	satellites = newSatellites

	w.Header().Set("Content-Type", "application/jspn")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newSatellites)
}

func GetUbicationAndMessage(w http.ResponseWriter, r *http.Request) {
	var satellitesMessages []requestType = satellites.Satellites
	if len(satellitesMessages) < 3 {
		http.Error(w, "Los datos son escasos", http.StatusBadRequest)
		return
	}

	var respUbicationAndMessage respUbicationAndMessage
	respUbicationAndMessage.Message = GetMessage(satellitesMessages[0].Message, satellitesMessages[1].Message, satellitesMessages[2].Message)
	respUbicationAndMessage.Position.X, respUbicationAndMessage.Position.Y = GetLocation(satellitesMessages[0].Distance, satellitesMessages[1].Distance, satellitesMessages[2].Distance)
	json.NewEncoder(w).Encode(respUbicationAndMessage)
}

func PostUbicationAndMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var satelliteName string = vars["satellite_name"]
	var newRequest requestType
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error al leer los datos", http.StatusBadRequest)
		return
	}
	json.Unmarshal(reqBody, &newRequest)
	newRequest.Name = satelliteName

	var satellitesMessages []requestType = satellites.Satellites
	var addMessage bool = true

	for i := 0; i < len(satellitesMessages); i++ {
		if newRequest.Name == satellitesMessages[i].Name {
			satellites.Satellites[i].Message = satellitesMessages[i].Message
			addMessage = false
		}
	}
	if addMessage {
		satellites.Satellites = append(satellites.Satellites, newRequest)
	}

	w.Header().Set("Content-Type", "application/jspn")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newRequest)
}

func containsInPosition(str string, arr []string) int {
	var position int = -1

	for i := 0; i < len(arr); i++ {
		position = position + 1
		if arr[i] == str {
			break
		}
	}

	return position
}

func GetValueNonBlank(values ...string) string {
	var valueStr string = ""

	for i := 0; i < len(values); i++ {
		if values[i] != "" {
			valueStr = values[i]
			break
		}
	}

	return valueStr
}

func GetMax(values ...int) int {
	var valueMax int = 0

	for i := 0; i < len(values); i++ {
		if values[i] > valueMax {
			valueMax = values[i]
		}
	}

	return valueMax
}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Bienvenido a la API ML-FBI-Challenge")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	// endpoints
	router.HandleFunc("/", indexRoute).Methods("GET")

	router.HandleFunc("/get_satellites", GetSatellites).Methods("GET")

	router.HandleFunc("/topsecret", PostSatellites).Methods("POST")

	router.HandleFunc("/topsecret_split", GetUbicationAndMessage).Methods("GET")
	router.HandleFunc("/topsecret_split/{satellite_name}", PostUbicationAndMessage).Methods("POST")

	// start the server
	log.Fatal(http.ListenAndServe(":3000", router))
}
