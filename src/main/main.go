package main

import (
	"net/http"
	"log"
	"flag"
	"fmt"
    "github.com/gorilla/sessions"
)

type Player struct {
   username string
   password string
}

type PlayerAccount struct {
	PlayerId int;
	Accounts map[string]int
}

func (pa * PlayerAccount) GetAccountString()string{
	var res string;
	for k,v := range pa.Accounts{
		res += k+":"+fmt.Sprintf("%d",v)+"; ";
	}
	return res;
}

type Room struct {
	roomid string
	maxPlayerIndex int
	StartResources map[string]int
}

var allPlayers map[string]Player;
var allRooms map[string]*Room;

var roomToPlayer map[string]map[string]*PlayerAccount;
var playerToRoom map[string]string;


var addr = flag.String("addr", ":8080", "http service address")


func isLogedIn(r *http.Request) bool{
    session, _ := store.Get(r, "cookie-name");
    recusername, ok := session.Values["username"].(string);
    if(!ok){
    	return false;
    }
    recpassword, _ := session.Values["password"].(string);
    plyInfo, plyExist := allPlayers[recusername];
    if(!plyExist){
    	return false;
    }
    if(plyInfo.password != recpassword){
    	return false;
    }
    return true;
}

func getSessionUsername(r *http.Request) string{
    session, _ := store.Get(r, "cookie-name");
    recusername, ok := session.Values["username"].(string);
    if(!ok){
    	return "__unknown_username__";
    }
    return recusername;
}

func resLogout(w http.ResponseWriter, r *http.Request){
    fmt.Println("logout") //get request method
    session, _ := store.Get(r, "cookie-name");
    session.Options = &sessions.Options{
	    Path:     "/",
	    MaxAge:   -1,
	    HttpOnly: true,
	}
    session.Save(r, w);
    http.Redirect(w, r, "/login", http.StatusFound);
}

func main(){
	go h.run();
	allPlayers = make(map[string]Player);
	allRooms = make(map[string]*Room);
	roomToPlayer = make(map[string]map[string]*PlayerAccount);
	playerToRoom = make(map[string]string);
	store.MaxAge(3*3600);
	http.HandleFunc("/", resLogin);
	http.HandleFunc("/login", resLogin);
	http.HandleFunc("/welcome", resWelcome);
	http.HandleFunc("/create", resCreateRoom);
	http.HandleFunc("/enter", resEnterRoom);
	http.HandleFunc("/exit", resExitRoom);
	http.HandleFunc("/logout", resLogout);
	http.HandleFunc("/gameroom", resGameRoom);
	http.HandleFunc("/ws", resWS);
	//http.HandleFunc("/ws", serveWs)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}