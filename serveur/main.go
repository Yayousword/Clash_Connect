package main

import (
	"log"
	"net/http"
	"github.com/rs/cors"
	"fmt"
	"encoding/json"
	"crypto/sha256"
	"github.com/dgrijalva/jwt-go"
	"database/sql"
	_ "github.com/lib/pq"
	"strconv"
    "encoding/hex"
	"io/ioutil"
	"errors"

)

const (
    host     = "localhost"
    port     = 5432
    user     = "clash_connect"
    password = "52fdc5a882ad0cc490297a43dce208cc36639f0c5224fc47bc849a978bd16d98"
    dbname   = "clash_connect"
)

const Log = "LOG : "

var mux = http.NewServeMux()

var secret string = "clash royale connect"
func generateToken(id string,tokenClashRoyal string) string{
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["tokenClashRoyal"] = tokenClashRoyal
	
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("Error creating token string:", err)
		return ""
	}
	return tokenString
}

type token_Clash struct {
	Identity string
	TokenClashRoyal string

}

func verifyToken(token string) (token_Clash, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("Error parsing token:", err)
		return token_Clash{Identity:"" , TokenClashRoyal:""}, err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		id := claims["id"].(string)
		tokenClashRoyal := claims["tokenClashRoyal"].(string)

		return token_Clash{Identity:id , TokenClashRoyal:tokenClashRoyal},nil
	} else {
		fmt.Println("Invalid token")
		return token_Clash{Identity:"" , TokenClashRoyal:""}, errors.New("Invalid token")
	}
}







func main() {


	/*
	* Connect to database
	*/
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	
	fmt.Println(Log + "Info BDD : " + psqlInfo)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	type Response struct {
		Success int   `json:"success"`
		Message string `json:"message,omitempty"`
	}


	/*
	* signIn
	* METHOD : POST , ROOT : "/auth/signIn"
	* need : {
		identity : string, 
		password : string, 
		tokenAPI : string 
	}
	* search the identity in authentication table
	* verify password
	* generate token
	* return {token : token,response : 200} || {response : 404, error : "identity not found"} 
	*/
	mux.HandleFunc("/auth/signIn",func(w http.ResponseWriter, r *http.Request) {
		
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			fmt.Println(Log + "Method not allowed")
			return
		}

		type Args struct {
			Identity string `json:"identity"`
			Password string `json:"password"`
			TokenAPI string `json:"tokenAPI"`
		}

		

		var args Args
		if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println(Log + "Parse args : " + err.Error())
			return
		}

		
		row := db.QueryRow("SELECT password FROM authentication WHERE id=$1", args.Identity)

		var password string
		if err := row.Scan(&password); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println(Log + err.Error())
			return
		}

		hash := sha256.Sum256([]byte(args.Password))

		if password != hex.EncodeToString(hash[:]) {
			http.Error(w, "Invalid id or password", http.StatusBadRequest)
			fmt.Println(Log + "Invalid id or password")
			return
		}

		

		client := &http.Client{}

		req, err := http.NewRequest("GET", "https://api.clashroyale.com/v1/players/"+ args.Identity, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(Log + "Req Clash API error");
			return
		}
	

		req.Header.Set("Authorization", "Bearer " + args.TokenAPI)
	
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(Log + "Response clash API error");
			return
		}

		defer resp.Body.Close()
		

		token := generateToken(args.Identity,args.TokenAPI)
		
		fmt.Println(Log + "token : " + token)

		response := Response{
			Success: 200,
			Message: "Successfully logged in",
		}
	
		jsonBytes, err := json.Marshal(map[string]interface{}{
			"token": token,
			"resp":  response,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(Log + err.Error())
			return
		}
	
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)
		

		fmt.Println(Log + "signIn Success")

	})
	
	/*
	* signUp
	* METHOD : POST , ROOT : "/auth/signUp"
	* need :  {
		identity : string,
		email : string, 
		password : string
	}
	* search the identity in authentication table 
	* set identity and password in authentication table
	* return {response : 200} || {response : 402, error : "identity already exists"} 
	*/
	mux.HandleFunc("/auth/signUp",func(w http.ResponseWriter, r *http.Request) {
				
		if r.Method != "POST" {

			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			fmt.Println("LOG : Method not allowed");
			return
		}

		type Args struct {
			Identity string `json:"identity"`
			Password string `json:"password"`
			Email string `json:"email"`
		}

	
		var args Args
		if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		hash := sha256.Sum256([]byte(args.Password))

		_,err := db.Exec("INSERT into authentication (id,email,password) VALUES ($1, $2, $3) ", args.Identity, args.Email, hex.EncodeToString(hash[:]))

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println("LOG : BDD Error "+ err.Error());
			return
		}
		
		response := Response{
			Success: 200,
			Message: "Successfully  signUp",
		}
	
		jsonBytes, err := json.Marshal(response)
		
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(Log + "Parse response error");
			return
		}
	
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)
		
		fmt.Println(Log + "signUp Success")

	})

	
	type publication struct {
		ID       string    `json:"id"`
		Time     string `json:"time"`
		Rate     int       `json:"rate"`
		Nb_Rate int     `json:"nombre_rate"`
	}



	// Define the structure of a card response from the Clash Royale API
	type Card struct {
		Name     string `json:"name"`
		Level    int    `json:"level"`
		MaxLevel int    `json:"maxLevel"`
		IconURL  string `json:"iconUrl"`
	}

	// Define the structure of a player response from the Clash Royale API
	type Player struct {
		Name         string `json:"name"`
		Tag          string `json:"tag"`
		StartingTrophies int `json:"startingTrophies"`
		Trophies     int    `json:"trophies"`
		Clan         struct {
			Tag  string `json:"tag"`
			Name string `json:"name"`
		} `json:"clan"`
		Cards []*Card `json:"cards"`
	}
	
	

	type BattleResponse struct {
		Type              string `json:"type"`
		BattleTime        string `json:"battleTime"`
		IsLadderTournament bool   `json:"isLadderTournament"`
		Arena             struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"arena"`
		GameMode struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"gameMode"`
		DeckSelection string `json:"deckSelection"`
		Team          []struct {
			Tag                 string   `json:"tag"`
			Name                string   `json:"name"`
			Crowns              int      `json:"crowns"`
			KingTowerHitPoints  int      `json:"kingTowerHitPoints"`
			PrincessTowersHitPoints []int `json:"princessTowersHitPoints"`
			Cards               []struct {
				Name       string `json:"name"`
				ID         int    `json:"id"`
				Level      int    `json:"level"`
				StarLevel  int    `json:"starLevel"`
				MaxLevel   int    `json:"maxLevel"`
				IconUrls   struct {
					Medium string `json:"medium"`
				} `json:"iconUrls"`
			} `json:"cards"`
			ElixirLeaked float64 `json:"elixirLeaked"`
		} `json:"team"`
		Opponent []struct {
			Tag                 string   `json:"tag"`
			Name                string   `json:"name"`
			TrophyChange        int      `json:"trophyChange"`
			Crowns              int      `json:"crowns"`
			KingTowerHitPoints  int      `json:"kingTowerHitPoints"`
			PrincessTowersHitPoints []int `json:"princessTowersHitPoints"`
			Clan struct {
				Tag     string `json:"tag"`
				Name    string `json:"name"`
				BadgeID int    `json:"badgeId"`
			} `json:"clan"`
			Cards []struct {
				Name     string `json:"name"`
				ID       int    `json:"id"`
				Level    int    `json:"level"`
				StarLevel  int  `json:"starLevel"`
				MaxLevel int    `json:"maxLevel"`
				IconUrls struct {
					Medium string `json:"medium"`
				} `json:"iconUrls"`
			} `json:"cards"`
			ElixirLeaked float64 `json:"elixirLeaked"`
		} `json:"opponent"`
		IsHostedMatch bool `json:"isHostedMatch"`
	}

	type Comment struct {
		Time  int `json:"time_comment"`
		Comment string `json:"comment"`
	}
		
	type BattleResponseClient struct {
			BattleInfo BattleResponse `json:"battleClash"`
			Comments []Comment `json:"comments"`
			Rate int `json:"rate"`
			Nb_Rate int `json:"nb_rate"`
	}


	/*
	* search
	* METHOD : GET , ROOT : "/search/combat"
	* need : { 
		token : string,
		rate : double
	}
	* search in publication table combats where rate is lower or equal to 'rate' and carts are used in the combat
	* return {combats : [combat], response : 200} 
	where combat = {
		{
					BattleInfo :battlesResponseClash ,
					Comments : comments,
					Rate : pub.Rate,
					Nb_Rate : pub.Nb_Rate,
				}
	where comment = {
		id: comment-sender,
		comment : string,
		time : timestamp
	}
	|| {response : 404, error : "response not found"} 
	|| {response : 403, error : "authorization denied"}

	* USE CLASH API
	*/
	mux.HandleFunc("/search/combat",func(w http.ResponseWriter, r *http.Request) {
		
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			fmt.Println(Log + "Method not allowed");
			return
		}


		type Args struct {
			Token string `json:"token"`
			Rate  int `json:"rate"`
		}
		rate, err1:= strconv.Atoi(r.URL.Query().Get("rate"))

		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusMethodNotAllowed)
			fmt.Println(Log + err1.Error());
			return
		}
		args := Args{Token : r.URL.Query().Get("token"), Rate : rate}


		if args.Rate < 0 {
			http.Error(w,  "Rate must be positive", http.StatusBadRequest)
			fmt.Println(Log + "Args value error:  " + "Rate must be positive");

			return
		}

		connectInfo, err := verifyToken(args.Token)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(Log + "Token error :" + err.Error());
			return
		}


		row,err := db.Query("SELECT * FROM publication WHERE rate >= $1", args.Rate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println(Log + "BDD error : " + err.Error());
			return
		}
		defer row.Close()

		
		
		var battlesResponseClient []BattleResponseClient
		for row.Next() {
			var pub publication
			if err := row.Scan(pub.ID, pub.Time,pub.Rate,pub.Nb_Rate ); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				fmt.Println(Log + "BDD error : " + err.Error());
				return
			}

			client := &http.Client{}

			req, err := http.NewRequest("GET", "https://api.clashroyale.com/v1/players/%23"+ pub.ID+ "/battlelog", nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				fmt.Println(Log + "Request API error : " + err.Error());
				return
			}
		
			req.Header.Set("Authorization", "Bearer " + connectInfo.TokenClashRoyal)
		
			resp, err := client.Do(req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				fmt.Println(Log + "Response API error : " + err.Error());
				return
			}
			defer resp.Body.Close()


			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				fmt.Println(Log + "Read error : " + err.Error());
				return
			}
			type BattlesResponseClashTmp []BattleResponse 

			var battlesResponseClashTmp BattlesResponseClashTmp
			if err := json.Unmarshal(body, &battlesResponseClashTmp); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
	

			var battlesResponseClash BattleResponse
			for _, battleResponse := range battlesResponseClashTmp {
				if battleResponse.BattleTime == pub.Time {
					battlesResponseClash = battleResponse
					break
				}
			}


			var comments []Comment
			row,err := db.Query("SELECT * FROM comment WHERE id_publication=$1 AND time_publication=$2", pub.ID, pub.Time)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				fmt.Println(Log + "BDD error : " + err.Error());
				return
			}
			defer row.Close()

			for row.Next() {
				var com Comment
				if err := row.Scan(com.Time, com.Comment); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					fmt.Println(Log + "BDD error : " + err.Error());
					return
				}
				comments = append(comments, com)

			}


			battlesResponseClient = append(battlesResponseClient, BattleResponseClient{
				BattleInfo :battlesResponseClash ,
				Comments : comments,
				Rate : pub.Rate,
				Nb_Rate : pub.Nb_Rate,
			})			
			
		}


		jsonBytes, err := json.Marshal(battlesResponseClient)
		
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(Log + "Json parse error : " + err.Error());
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)

	})


	/*
	* publication
	* METHOD : GET , ROOT : "/publication"
	* need : {token : string}
	* get random publications from  publication  table
	* return {publications : [publication], response : 200} 
	where publication = {
		id : id-combatant,
		time : timestamp,
		rate : double
		comments : [comment],
		combat_clash: combat_from_api 
	}
	where comment = {
		id: comment-sender,
		comment : string,
		time : timestamp
	}
	|| {response : 404, error : "response not found"} 
	|| {response : 403, error : "authorization denied"}

	* USE CLASH API
	*/
	mux.HandleFunc("/publication",func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}


		type Args struct {
			Token string `json:"token"`
		}

		args := Args{Token : r.URL.Query().Get("token")}
		connectInfo, err := verifyToken(args.Token)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}


		row,err := db.Query("SELECT * FROM publication")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer row.Close()

		
		
		var battlesResponseClient []BattleResponseClient
		for row.Next() {
			var pub publication
			if err := row.Scan(pub.ID, pub.Time,pub.Rate,pub.Nb_Rate ); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			client := &http.Client{}

			req, err := http.NewRequest("GET", "https://api.clashroyale.com/v1/players/%23"+ pub.ID+ "/battlelog", nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		
			req.Header.Set("Authorization", "Bearer " + connectInfo.TokenClashRoyal)
		
			resp, err := client.Do(req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer resp.Body.Close()


			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			type BattlesResponseClashTmp []BattleResponse 

			var battlesResponseClashTmp BattlesResponseClashTmp
			if err := json.Unmarshal(body, &battlesResponseClashTmp); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
	

			var battlesResponseClash BattleResponse
			for _, battleResponse := range battlesResponseClashTmp {
				if battleResponse.BattleTime == pub.Time {
					battlesResponseClash = battleResponse
					break
				}
			}


			var comments []Comment
			row,err := db.Query("SELECT * FROM comment WHERE id_publication=$1 AND time_publication=$2", pub.ID, pub.Time)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer row.Close()

			for row.Next() {
				var com Comment
				if err := row.Scan(com.Time, com.Comment); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				comments = append(comments, com)

			}


			battlesResponseClient = append(battlesResponseClient, BattleResponseClient{
				BattleInfo :battlesResponseClash ,
				Comments : comments,
				Rate : pub.Rate,
				Nb_Rate : pub.Nb_Rate,
			})			
			
		}


		jsonBytes, err := json.Marshal(battlesResponseClient)
		
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)
	})

	/*
	* publication
	* METHOD : POST , ROOT : "/publication/publish"
	* need : {combat-time : string,token : string}
	* add publication into publication table  
	* return {response : 200} 
	|| {response : 404, error : "primary key (id,timestamp) not found"} 
	|| {response : 403, error : "authorization denied"}

	* USE CLASH API
	*/
	mux.HandleFunc("/publication/publish",func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}


		type Args struct {
			Token string `json:"token"`
			Time string `json:"combat-time"`
		}

		var args Args
		if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		connectInfo, err := verifyToken(args.Token)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}


		_, err = db.Exec("INSERT INTO publication (id, time, rate, nombre_rate) VALUES ($1,$2,$3,$4)", connectInfo.Identity, args.Time,0,0)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		
		response := Response{
			Success: 200,
			Message: "Successfully  publish",
		}
	
		jsonBytes, err := json.Marshal(response)
		
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)
		
	})

	/*
	* publication
	* METHOD : POST , ROOT : "/publication/comment"
	* need : {id-fighter : string ,combat-time : timestamp, comment : string, token : string}
	* add comment into comment table  
	* return {response : 200} 
	|| {response : 404, error : "primary key (id,timestamp) not found"} 
	|| {response : 403, error : "authorization denied"}

	* USE CLASH API
	*/
	mux.HandleFunc("/publication/comment",func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}


		type Args struct {
			Token string `json:"token"`
			ID string `json:"id-fighter"`
			Time string `json:"combat-time"`
			Comment string `json:"comment"`
		}

		var args Args
		if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if args.Comment == "" {
			http.Error(w,  "empty comment", http.StatusBadRequest)
			return
		}

		connectInfo, err := verifyToken(args.Token)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		
		row := db.QueryRow("SELECT id FROM publication WHERE id=$1", args.ID )

		var ID string
		if err := row.Scan(&ID); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		

		_, err = db.Exec("INSERT INTO comment (id_sender ,id_publication ,time_publication,comment) VALUES ($1,$2,$3,%4)", connectInfo.Identity, args.ID,args.Time,args.Comment)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}


		response := Response{
			Success: 200,
			Message: "Successfully  Comment",
		}
	
		jsonBytes, err := json.Marshal(response)
		
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)
	})
	
	/*
	* publication
	* METHOD : POST , ROOT : "/publication/rate"
	* need : {id-fighter : string ,combat-time : timestamp, rate : integer, token : string}
	* add rate into publication table  
	* return {response : 200} 
	|| {response : 404, error : "primary key (id,timestamp) not found"} 
	|| {response : 403, error : "authorization denied"}

	* USE CLASH API
	*/
	mux.HandleFunc("/publication/rate",func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}


		type Args struct {
			Token string `json:"token"`
			ID string `json:"id-fighter"`
			Time string `json:"combat-time"`
			Rate int `json:"comment"`
		}

		var args Args
		if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if args.Rate < 0  && args.Rate > 5{
			http.Error(w,  "Rate must be positive and lower than 5", http.StatusBadRequest)
			return
		}

		_, err := verifyToken(args.Token)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		type responseTp struct{
			Rate int `json:"rate"`
			Nb_Rate int `json:"nombre_rate"`
		}

		row := db.QueryRow("SELECT rate, nombre_rate FROM publication WHERE id=$1", args.ID )

		var pub responseTp
		if err := row.Scan(pub.Rate, pub.Nb_Rate); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}


		
		pre, err := db.Prepare("UPDATE publication SET rate = $1, nombre_rate = $2 WHERE id=$3 AND time= $4")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer pre.Close()

		_, err = pre.Exec(args.Rate + pub.Rate, pub.Nb_Rate + 1,args.ID,args.Time)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}


		response := Response{
			Success: 200,
			Message: "Successfully  Comment",
		}
	
		jsonBytes, err := json.Marshal(response)
		
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)
	})

	/*
	* chat
	* METHOD : POST , ROOT : "/chat/send"
	* need : {message : string, token : string}
	* add message into chat table  
	* return {response : 200} 
	|| {response : 403, error : "authorization denied"}
	*/
	mux.HandleFunc("/chat/send",func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			fmt.Println(Log + "Chat/send " + "Method not allowed")
			return
		}

		type Args struct {
			Token string `json:"token"`
			Message string `json:"message"`
		}

		var args Args
		if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println(Log + "Chat/send "+ err.Error())

			return
		}

		
		if args.Message == "" {
			http.Error(w,  "empty message", http.StatusBadRequest)
			fmt.Println(Log + "Chat/send "+ "empty message")
			return
		}

		connectInfo, err := verifyToken(args.Token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(Log + "Chat/send "+ err.Error())
			
			return
		}

	

		_, err = db.Exec("INSERT INTO chat (id_sender ,message) VALUES ($1,$2)", connectInfo.Identity, args.Message)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println(Log + "Chat/send "+ err.Error())
			
			return
		}


		response := Response{
			Success: 200,
			Message: "Successfully  send Message",
		}
	
		jsonBytes, err := json.Marshal(response)
		
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(Log + "Chat/send "+ err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)


	})


	type Chat struct{
		Id_sender string `json:"id_sender"`
		Message string `json:"message"`
		Timestamp string `json:"timestamp"`
	}

	/*
	* chat
	* METHOD : GET , ROOT : "/chat"
	* need : {token : string}
	* get last 30 messages into chat table 
	* return {chats: [chat],response : 200} 
	|| {response : 403, error : "authorization denied"}

	** 
	*/
	mux.HandleFunc("/chat",func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			fmt.Println(Log + "Chat/ "+ "Method not allowed")
			return
		}


		type Args struct {
			Token string `json:"token"`
		}

		args := Args{Token : r.URL.Query().Get("token")}

		_, err := verifyToken(args.Token)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(Log + "Chat/ "+ err.Error())
			return
		}


		row,err := db.Query("SELECT * FROM chat")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println(Log + "SELECT Chat/ "+ err.Error())
			return
		}
		defer row.Close()

		var chats []Chat
		for row.Next() {
			var chat Chat
			if err := row.Scan(&(chat.Id_sender),&(chat.Timestamp),&(chat.Message)); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				fmt.Println(Log + "Chat/ "+ err.Error())
				return
			}
			chats = append(chats, chat)
		}

		result ,err:=json.Marshal(chats)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(Log + "Chat/ "+ err.Error())
			return
		}
	
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(result)

	})


	type PlayerInfo struct {
		Tag              string `json:"tag"`
		Name             string `json:"name"`
		ExpLevel         int    `json:"expLevel"`
		Trophies         int    `json:"trophies"`
		BestTrophies     int    `json:"bestTrophies"`
		Wins             int    `json:"wins"`
		Losses           int    `json:"losses"`
		BattleCount      int    `json:"battleCount"`
		ThreeCrownWins   int    `json:"threeCrownWins"`
		ChallengeCardsWon int   `json:"challengeCardsWon"`
		ChallengeMaxWins int    `json:"challengeMaxWins"`
		TournamentCardsWon int `json:"tournamentCardsWon"`
		TournamentBattleCount int `json:"tournamentBattleCount"`
		Role             string `json:"role"`
		Donations        int    `json:"donations"`
		DonationsReceived int   `json:"donationsReceived"`
		TotalDonations   int    `json:"totalDonations"`
		WarDayWins       int    `json:"warDayWins"`
		ClanCardsCollected int `json:"clanCardsCollected"`
		Clan              struct {
			Tag    string `json:"tag"`
			Name   string `json:"name"`
			Badge  struct {
				ID    int    `json:"id"`
				Name  string `json:"name"`
				Image string `json:"image"`
			} `json:"badge"`
		} `json:"clan"`
		Achievements []struct {
			Name     string `json:"name"`
			Stars    int    `json:"stars"`
			Value    int    `json:"value"`
			Target   int    `json:"target"`
			Info     string `json:"info"`
			CompletionInfo string `json:"completionInfo"`
		} `json:"achievements"`
		Cards []struct {
			Name   string `json:"name"`
			Level  int    `json:"level"`
			MaxLevel int  `json:"maxLevel"`
			Count  int    `json:"count"`
			IconUrls struct {
				Medium string `json:"medium"`
			} `json:"iconUrls"`
		} `json:"cards"`
		CurrentDeck []struct {
			Name   string `json:"name"`
			Level  int    `json:"level"`
			MaxLevel int  `json:"maxLevel"`
			Count  int    `json:"count"`
			IconUrls struct {
				Medium string `json:"medium"`
			} `json:"iconUrls"`
		} `json:"currentDeck"`
	}

	/*
	* User
	* METHOD : GET , ROOT : "/profile"
	* need : {token : string,userId : string}
	* add rate into publication table  
	* return all information about user and their recent battles 

	|| {response : 403, error : "authorization denied"}

	* USE CLASH API
	*/
	mux.HandleFunc("/profile",func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}


		type Args struct {
			Token string `json:"token"`
			UserID string `json:"userID"`
		}

		args := Args{Token : r.URL.Query().Get("token"), UserID : r.URL.Query().Get("userID")}


		connectInfo, err := verifyToken(args.Token)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req, err := http.NewRequest("GET", "https://api.clashroyale.com/v1/players/%23"+ args.UserID, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	
		req.Header.Set("Authorization", "Bearer " + connectInfo.TokenClashRoyal)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer resp.Body.Close()


		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var player PlayerInfo
		if err := json.Unmarshal(body, &player); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	

		req, err = http.NewRequest("GET", "https://api.clashroyale.com/v1/players/%23"+ args.UserID + "/battlelog", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	
		req.Header.Set("Authorization", "Bearer " + connectInfo.TokenClashRoyal)
	
		resp, err = client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer resp.Body.Close()


		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		type BattlesResponseClashTmp []BattleResponse 

		var battlesResponseClashTmp BattlesResponseClashTmp
		if err := json.Unmarshal(body, &battlesResponseClashTmp); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		
		type response struct {
			PlayerInfo PlayerInfo `json:"playerInfo"`
			BattlesResponseClash BattlesResponseClashTmp `json:"battlesResponse"`
		}

		respClient := response{
			PlayerInfo : player,
			BattlesResponseClash : battlesResponseClashTmp,
		}


		jsonBytes, err := json.Marshal(respClient)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)
	})


	/*
	* User
	* METHOD : GET , ROOT : "/profile"
	* need : {token : string,userId : string}
	* add rate into publication table  
	* return all information about user and their recent battles 

	|| {response : 403, error : "authorization denied"}

	* USE CLASH API
	*
	*
	* Used to retrieve only the battles of users that are published
	*/
	mux.HandleFunc("/profile/other",func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}


		type Args struct {
			Token string `json:"token"`
			Identity string `json:"identity"`
		}

		args := Args{Token : r.URL.Query().Get("token")}


		connectInfo, err := verifyToken(args.Token)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req, err := http.NewRequest("GET", "https://api.clashroyale.com/v1/players/%23"+ args.Identity, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	
		req.Header.Set("Authorization", "Bearer " + connectInfo.TokenClashRoyal)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer resp.Body.Close()


		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var player PlayerInfo
		if err := json.Unmarshal(body, &player); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	


		row,err := db.Query("SELECT * FROM publication WHERE id=$1", args.Identity)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer row.Close()

		
		
		var battlesResponseClient []BattleResponseClient
		for row.Next() {
			var pub publication
			if err := row.Scan(pub.ID, pub.Time,pub.Rate,pub.Nb_Rate ); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			client := &http.Client{}

			req, err := http.NewRequest("GET", "https://api.clashroyale.com/v1/players/%23"+ pub.ID+ "/battlelog", nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		
			req.Header.Set("Authorization", "Bearer " + connectInfo.TokenClashRoyal)
		
			resp, err := client.Do(req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer resp.Body.Close()


			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			type BattlesResponseClashTmp []BattleResponse 

			var battlesResponseClashTmp BattlesResponseClashTmp
			if err := json.Unmarshal(body, &battlesResponseClashTmp); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
	

			var battlesResponseClash BattleResponse
			for _, battleResponse := range battlesResponseClashTmp {
				if battleResponse.BattleTime == pub.Time {
					battlesResponseClash = battleResponse
					break
				}
			}


			var comments []Comment
			row,err := db.Query("SELECT * FROM comment WHERE id_publication=$1 AND time_publication=$2", pub.ID, pub.Time)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer row.Close()

			for row.Next() {
				var com Comment
				if err := row.Scan(com.Time, com.Comment); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				comments = append(comments, com)

			}


			battlesResponseClient = append(battlesResponseClient, BattleResponseClient{
				BattleInfo :battlesResponseClash ,
				Comments : comments,
				Rate : pub.Rate,
				Nb_Rate : pub.Nb_Rate,
			})			
			
		}

		
		type response struct {
			PlayerInfo PlayerInfo `json:"playerInfo"`
			BattlesResponseClash []BattleResponseClient `json:"battlesResponse"`
		}

		respClient := response{
			PlayerInfo : player,
			BattlesResponseClash : battlesResponseClient,
		}

		jsonBytes, err := json.Marshal(respClient)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)
	})

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Origin", "Content-Type", "Accept"},
	})

	handler := c.Handler(mux)
	

	// Démarrage du serveur
	log.Fatal(http.ListenAndServe(":8080", handler))
}



