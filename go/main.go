package main

import (
	"encoding/json"
	"fmt"
	"io"
	"main/calculation"
	"main/db"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// setting the router to handle paths
	r := mux.NewRouter()
	//sends the /api/endpoint to the handler
	r.HandleFunc("/api/{endpoint}", handler).Methods("POST", "GET", "DELETE", "UPDATE")
	//setting cors, currently set with wildcard for testing
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
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

	case "get-all-builds":
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		_, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request", http.StatusBadRequest)
			return
		}
		err = db.ServeBuilds()
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error during DB transaction", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

	case "get-user-builds":

	case "savebuild":

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request", http.StatusBadRequest)
			return
		}
		var saveData db.SaveData
		if err := json.Unmarshal(body, &saveData); err != nil {
			http.Error(w, "Error decoding JSON data", http.StatusBadRequest)
			return
		}
		err = db.SaveBuild(saveData)
		if err != nil {
			fmt.Println("Database error: ", err)
			http.Error(w, "Error during DB transaction", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		result := "Build saved successfully"
		resultjson, err := json.Marshal(result)
		if err != nil {
			return
		}
		w.Write(resultjson)

	case "signup":

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request", http.StatusBadRequest)
			return
		}
		var signupQuery db.SignupData
		if err := json.Unmarshal(body, &signupQuery); err != nil {
			http.Error(w, "Error decoding JSON data", http.StatusBadRequest)
			return
		}
		err = db.ConnectForSignup(signupQuery)
		if err != nil {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		result := "Signup successful"
		signupResult, err := json.Marshal(result)
		if err != nil {
			return
		}

		w.Write(signupResult)

	case "login":

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request", http.StatusBadRequest)
			return
		}

		var loginQuery db.LoginData
		if err := json.Unmarshal(body, &loginQuery); err != nil {
			http.Error(w, "Error decoding JSON data", http.StatusBadRequest)
			return
		}
		fmt.Println(loginQuery.Username)
		fmt.Println(loginQuery.Password)
		session_data, err := db.ConnectForLogin(loginQuery)
		if err != nil {
			http.Error(w, "Error during login", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(session_data)

	case "calculate":
		//setting cors headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
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
		returnedFeats, err := calculation.PrepareData(data)
		//if the function returns an error, returns and writes a error
		if err != nil {
			http.Error(w, "Error during calculation", http.StatusInternalServerError)
			return
		}
		//if there is no error, statusOKs and writes makes a response using the returnedFeats json
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Println(" ")
		fmt.Println("----------------------")
		fmt.Println("Finished... making response..")
		fmt.Println("----------------------")
		fmt.Println(" ")

		w.Write(returnedFeats)
		return
	default:
		http.Error(w, "Invalid endpoint", http.StatusNotFound)
	}
}
