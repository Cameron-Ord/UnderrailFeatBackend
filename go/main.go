package main

import (
	"encoding/json"
	"fmt"
	"io"
	"main/calculation"
	"main/db"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var url string

func main() {
	if db.Production_Mode {
		url = "https://featgenerator.cameron-ord.online"
	} else if !db.Production_Mode {
		url = "http://localhost:5173"
	}
	// setting the router to handle paths
	r := mux.NewRouter()
	//sends the /api/endpoint to the handler
	r.HandleFunc("/api/{endpoint}", handler).Methods("POST", "GET", "DELETE", "UPDATE")
	//setting cors, currently set with wildcard for testing
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{url}),
		handlers.AllowedMethods([]string{"POST", "GET", "DELETE", "UPDATE"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)
	http.Handle("/", corsHandler(r))
	//listening on the set port
	port := "6969"
	fmt.Printf("Server is running on http://localhost:%s\n", port)
	http.ListenAndServe(":"+port, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	endpoint := vars["endpoint"]
	switch endpoint {

	case "get-profile-info":
		w.Header().Set("Access-Control-Allow-Origin", url)
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		session_token := r.URL.Query().Get("session_token")
		client_id := r.URL.Query().Get("client_id")
		if session_token == "" || client_id == "" {
			http.Error(w, "Missing or invalid params", http.StatusBadRequest)
			return
		}
		uintVal, err := strconv.ParseUint(client_id, 10, 64)
		if err != nil {
			http.Error(w, "Error converting string values", http.StatusInternalServerError)
			return
		}
		client_id_uint := uint(uintVal)
		user_session_data := db.User_Session_Data{
			Client_Session_Token: session_token,
			Client_ID_Value:      client_id_uint,
		}
		jsonified_data, err := db.GetProfileInfo(user_session_data)
		if err != nil {
			http.Error(w, "Error getting profile info", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonified_data)

	case "get-all-builds":
		w.Header().Set("Access-Control-Allow-Origin", url)
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		_, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request", http.StatusBadRequest)
			return
		}
		jsonified_data, err := db.ServeBuilds()
		if err != nil {
			http.Error(w, "Error during DB transaction", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonified_data)

	case "get-user-builds":
		w.Header().Set("Access-Control-Allow-Origin", url)
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		session_token := r.URL.Query().Get("session_token")
		client_id := r.URL.Query().Get("client_id")

		if session_token == "" || client_id == "" {
			http.Error(w, "Missing or invalid params", http.StatusBadRequest)
			return
		}
		//base 10, 64 bit
		uintVal, err := strconv.ParseUint(client_id, 10, 64)
		if err != nil {
			http.Error(w, "Error converting string values", http.StatusBadRequest)
			return
		}

		client_id_uint := uint(uintVal)
		user_session_data := db.User_Session_Data{
			Client_Session_Token: session_token,
			Client_ID_Value:      client_id_uint,
		}

		jsonified_data, err := db.GetUserBuilds(user_session_data)
		if err != nil {
			http.Error(w, "Error retrieving database data", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonified_data)

	case "savebuild":

		w.Header().Set("Access-Control-Allow-Origin", url)
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request", http.StatusBadRequest)
			return
		}
		var saveData db.SaveData
		if err := json.Unmarshal(body, &saveData); err != nil {
			http.Error(w, "Error decoding JSON data", http.StatusInternalServerError)
			return
		}
		err = db.SaveBuild(saveData)
		if err != nil {
			http.Error(w, "Error during DB transaction", http.StatusBadRequest)
			return
		}
		result := "Build saved successfully"
		resultjson, err := json.Marshal(result)
		if err != nil {
			http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resultjson)

	case "signup":

		w.Header().Set("Access-Control-Allow-Origin", url)
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request", http.StatusBadRequest)
			return
		}

		var signupQuery db.SignupData
		if err := json.Unmarshal(body, &signupQuery); err != nil {
			http.Error(w, "Error decoding JSON data", http.StatusInternalServerError)
			return
		}

		err = db.ConnectForSignup(signupQuery)
		if err != nil {
			http.Error(w, "Error during signup", http.StatusBadRequest)
			return
		}

		result := "Signup successful"
		signupResult, err := json.Marshal(result)

		if err != nil {
			http.Error(w, "Error during json marshalling", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(signupResult)
	case "login":

		w.Header().Set("Access-Control-Allow-Origin", url)
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request", http.StatusBadRequest)
			return
		}

		var loginQuery db.LoginData
		if err := json.Unmarshal(body, &loginQuery); err != nil {
			http.Error(w, "Error decoding JSON data", http.StatusInternalServerError)
			return
		}

		session_data, err := db.ConnectForLogin(loginQuery)
		if err != nil {
			http.Error(w, "Error during login", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(session_data)

	case "calculate":
		//setting cors headers
		w.Header().Set("Access-Control-Allow-Origin", url)
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		//reading all data sent from frontend
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request", http.StatusBadRequest)
			return
		}
		//unmarshalling the json data to the structs in calculation.go
		var data calculation.RequestData
		if err := json.Unmarshal(body, &data); err != nil {
			http.Error(w, "Error decoding JSON data", http.StatusBadRequest)
			return
		}
		//initializing the skill and stat checkers and assigning the returned result
		returnedFeats, err := calculation.Prepare_Data(data)
		//if the function returns an error, returns and writes a error
		if err != nil {
			http.Error(w, "Error during calculation", http.StatusInternalServerError)
			return
		}
		//if there is no error, statusOKs and writes makes a response using the returnedFeats json
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(returnedFeats)

	default:
		http.Error(w, "Invalid endpoint", http.StatusBadGateway)
		w.Header().Set("Content-Type", "application/json")
	}
}
