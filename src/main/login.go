package main

import (
    "fmt"
    "html/template"
	"net/http"
    "github.com/gorilla/sessions"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

func signin(r *http.Request) (theusername string, success bool){
	success = true;
	theusername = r.Form["tfname"][0]
	playerinfo, isexist := allPlayers[theusername]
	if isexist{
		if playerinfo.password == r.Form["tfpswd"][0]{
			fmt.Println(theusername, "login")
		} else{
			fmt.Println("wrong password")
			success = false
		}
	}else{
		var playerinfo Player
		playerinfo.username = r.Form["tfname"][0]
		playerinfo.password = r.Form["tfpswd"][0]
		allPlayers[theusername] = playerinfo
		fmt.Println(theusername, "create and login");
	}
	return
}

type logindata struct {
	Tips string
}

func resLogin(w http.ResponseWriter, r *http.Request){
	if(isLogedIn(r)){
		http.Redirect(w, r, "/welcome", http.StatusFound);
	}
    fmt.Println("Pg Login", r.Method) //get request method
    session, _ := store.Get(r, "cookie-name");
    if r.Method == "GET" {
        t, _ := template.ParseFiles("login.gtpl")
        t.Execute(w, nil)
    } else {
        // logic part of log in
		r.ParseForm();
        theusername, sisuccess := signin(r)
        if sisuccess == true{
			session.Values["authenticated"] = true
			session.Values["username"] = theusername;
			session.Values["password"] = r.Form["tfpswd"][0];
			session.Save(r, w)
        	http.Redirect(w, r, "/welcome", http.StatusFound)
        } else{
        	var thetip = logindata{Tips: "wrong password"};
	        t, _ := template.ParseFiles("login.gtpl")
	        t.Execute(w, thetip)
        }
    }
}