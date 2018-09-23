package main

import (
	"net/http"
    "html/template"
	"fmt"
	"math/rand"
)

func createNewRoomID()(newRoomID string, success bool){
	success = true;
	if(len(allRooms) > 1000){
		success = false;
		return;
	}
	for i := 0; i < 9999; i++{
		var newroomnum int = rand.Intn(9999)
		newRoomID = fmt.Sprintf("%04d",newroomnum);
		_, isexist := allRooms[newRoomID]
		if !isexist {
			fmt.Println("create room:", newRoomID);
			return;
		}
	}
	success = false;
	return;
}

func exitRoom(playername string, roomid string){
	currentroomid, ok := playerToRoom[playername];
	if ok {
		if(currentroomid == roomid){
			delete(playerToRoom, playername);
			return;
		}
	}
	_, roomexist := roomToPlayer[roomid];
	if(roomexist){
		delete(roomToPlayer[roomid], playername);
	}
	if(len(roomToPlayer[roomid]) == 0){
		delete(roomToPlayer, roomid);
		delete(allRooms, roomid);
	}
}

func enterRoom(playername string, roomid string) bool{
	_, roomExist := roomToPlayer[roomid];
	if(!roomExist){
		return false;
	}
	currentroomid, ok := playerToRoom[playername];
	if ok {
		if(currentroomid == roomid){
			return true;
		}
	}
	playerToRoom[playername] = roomid;
	maxPI := allRooms[roomid].maxPlayerIndex;
	maxPI++;
	allRooms[roomid].maxPlayerIndex = maxPI;
	roomToPlayer[roomid][playername] = &PlayerAccount{
		PlayerId: maxPI,
		Accounts: make(map[string]int),
	};
	for resName, resValue := range allRooms[roomid].StartResources{
		roomToPlayer[roomid][playername].Accounts[resName] = resValue;
	}
	return true;
}

func resCreateRoom(w http.ResponseWriter, r *http.Request) {
	if(!isLogedIn(r)){
		http.Redirect(w, r, "/login", http.StatusFound);
	}
    session, _ := store.Get(r, "cookie-name");
    recusername, ok := session.Values["username"].(string);
    if(!ok){
    	return;
    }
    newRoomId, crtRoomSuc := createNewRoomID();
    if(!crtRoomSuc){
    	return;
    }
    
    allRooms[newRoomId] = &Room{newRoomId,0,make(map[string]int)};
    allRooms[newRoomId].StartResources["money"] = 47;
    roomToPlayer[newRoomId] = make(map[string]*PlayerAccount);
    enterRoom(recusername,newRoomId);
    http.Redirect(w, r, "/gameroom", http.StatusFound);
}

func resEnterRoom(w http.ResponseWriter, r *http.Request) {
	if(!isLogedIn(r)){
		http.Redirect(w, r, "/login", http.StatusFound);
	}
	//fmt.Fprintf(w, r.Form["lan_mode"][0])
    session, _ := store.Get(r, "cookie-name");
    recusername, ok := session.Values["username"].(string);
    if(!ok){
    	return;
    }
    r.ParseForm();
	targetRoom := r.Form["txtRoomNum"][0]
    if(!enterRoom(recusername, targetRoom)){
    	return;
    }
    http.Redirect(w, r, "/gameroom", http.StatusFound);
}

func resExitRoom(w http.ResponseWriter, r *http.Request) {
	if(!isLogedIn(r)){
		http.Redirect(w, r, "/login", http.StatusFound);
	}
	//fmt.Fprintf(w, r.Form["lan_mode"][0])
    session, _ := store.Get(r, "cookie-name");
    _, ok := session.Values["username"];
    if(!ok){
    	return;
    }
    return;
}

func generateRoomList()string{
	ret := "";
	for k, v := range roomToPlayer{
		var userlist string;
		userlist = "";
		for k1, _ := range v{
			userlist += k1+",";
		}
		ret += "<tr><td>"+k+"</td><td>"+userlist+"</td></tr>"
	}
	return ret;
}

type welcomeData struct{
	Roomlist map[string]map[string]*PlayerAccount;
}

func resWelcome(w http.ResponseWriter, r *http.Request) {
	if(!isLogedIn(r)){
		http.Redirect(w, r, "/login", http.StatusFound);
	}
	//fmt.Fprintf(w, r.Form["lan_mode"][0])

    fmt.Println("welcome method:", r.Method) //get request method
    //session, _ := store.Get(r, "cookie-name")
    if r.Method == "GET" {
    	wd := welcomeData{};
    	wd.Roomlist = roomToPlayer;
        t := template.Must(template.ParseFiles("welcome.gtpl"))
        t.Execute(w, wd);
    } else {
    	
    }
}