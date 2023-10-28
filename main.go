package main

import (
	"Josh/database"
	"crypto/sha256"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type PageHome struct {
	Posts           []database.Post
	IsConnecter     bool
	ConnectUserInfo string
	ConnectUserImg  string
	Error           int
	Try             bool
}

// Error = 0 (Aucun problème)
// Error = 1 (Session deja ouverte)
// Error = 2 (Problème de connection)
// Error = 3 (Problème d'Inscription)
// Error = 4 (Session expiré)

type PagePost struct {
	OnePost         database.Post
	IdPost          string
	Responses       []database.Reponse
	ConnectUserName string
	ConnectUserImg  string
}

type PageNewPost struct {
	ConnectUserImg string
}

func ResetDB() {
	database.Database()

	yannimg := renderImg()
	elisaimg := renderImg()
	kevinimg := renderImg()
	lilianeimg := renderImg()
	joshuaimg := renderImg()

	database.DatabaseAndUsers([]string{"yann@ynov.com", "Yann", HashPassword2("yann"), yannimg})
	database.DatabaseAndUsers([]string{"elisa@ynov.com", "Elisa", HashPassword2("elisa"), elisaimg})
	database.DatabaseAndUsers([]string{"kevin@ynov.com", "Kévin", HashPassword2("kevin"), kevinimg})
	database.DatabaseAndUsers([]string{"liliane@ynov.com", "Liliane", HashPassword2("liliane"), lilianeimg})
	database.DatabaseAndUsers([]string{"joshua@ynov.com", "Joshua", HashPassword2("joshua"), joshuaimg})

	database.DatabaseAndPost([]string{"Yann", "film", "First Post", "Moi j'adore ET", strconv.Itoa(55), strconv.Itoa(3), "27 Mai 2023", yannimg})
	database.DatabaseAndPost([]string{"Yann", "serie", "Second Post", "Moi j'adore GOT", strconv.Itoa(2), strconv.Itoa(33), "27 Mai 2023", yannimg})
	database.DatabaseAndPost([]string{"Elisa", "film", "First Post", "Moi j'adore ET", strconv.Itoa(55), strconv.Itoa(3), "27 Mai 2023", elisaimg})
	database.DatabaseAndPost([]string{"Elisa", "serie", "Second Post", "Moi j'adore Got", strconv.Itoa(2), strconv.Itoa(33), "27 Mai 2023", elisaimg})
	database.DatabaseAndPost([]string{"Kevin", "film", "First Post", "Moi j'adore ET", strconv.Itoa(55), strconv.Itoa(3), "27 Mai 2023", kevinimg})
	database.DatabaseAndPost([]string{"Kevin", "serie", "Second Post", "Moi j'adore Got", strconv.Itoa(2), strconv.Itoa(33), "27 Mai 2023", kevinimg})
	database.DatabaseAndPost([]string{"Liliane", "film", "First Post", "Moi j'adore ET", strconv.Itoa(55), strconv.Itoa(3), "27 Mai 2023", lilianeimg})
	database.DatabaseAndPost([]string{"Liliane", "serie", "Second Post", "Moi j'adore Got", strconv.Itoa(2), strconv.Itoa(33), "27 Mai 2023", lilianeimg})
	database.DatabaseAndPost([]string{"Joshua", "film", "First Post", "Moi j'adore ET", strconv.Itoa(55), strconv.Itoa(3), "27 Mai 2023", joshuaimg})
	database.DatabaseAndPost([]string{"Joshua", "serie", "Second Post", "Moi j'adore Got", strconv.Itoa(2), strconv.Itoa(33), "27 Mai 2023", joshuaimg})

	database.DatabaseAndReponse([]string{strconv.Itoa(1), "Yann", "Moi aussi !!!!!", transformDate(), yannimg})
	database.DatabaseAndReponse([]string{strconv.Itoa(1), "Elisa", "Moi aussi !!!!!", transformDate(), elisaimg})
	database.DatabaseAndReponse([]string{strconv.Itoa(1), "Kevin", "Moi aussi !!!!!", transformDate(), kevinimg})
	database.DatabaseAndReponse([]string{strconv.Itoa(1), "Liliane", "Moi aussi !!!!!", transformDate(), lilianeimg})
	database.DatabaseAndReponse([]string{strconv.Itoa(1), "Joshua", "Moi aussi !!!!!", transformDate(), joshuaimg})
}

func initStruct() (PageHome, PagePost, PageNewPost) {
	var home PageHome
	home.Posts = database.GetAllPost()
	home.IsConnecter = false
	home.Error = 0
	home.ConnectUserInfo = ""

	var post PagePost
	post.ConnectUserName = ""
	post.ConnectUserImg = ""

	var newpost PageNewPost

	return home, post, newpost
}

var tmplHome = template.Must(template.ParseFiles("./html/home.html"))
var tmplPost = template.Must(template.ParseFiles("./html/post.html"))
var tmplNewPost = template.Must(template.ParseFiles("./html/newpost.html"))
var HomeStruct, PostStruct, NewPostStruct = initStruct()

// var HomeStruct PageHome
// var PostStruct PagePost
// var NewPostStruct PageNewPost

func main() {

	// ResetDB()
	// database.Database()

	fmt.Printf("\n")
	fmt.Println("http://localhost:8080/")
	fmt.Printf("\n")

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ho", homeTransiHandler) // fonction de transition
	http.HandleFunc("/post", postHandler)
	http.HandleFunc("/newpost", newPostHandler)
	http.HandleFunc("/newpostTransi", newPostTransiHandler) // fonction de transition
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)

	http.Handle("/assets/css/", http.StripPrefix("/assets/css/", http.FileServer(http.Dir("assets/css"))))
	http.Handle("/assets/images/", http.StripPrefix("/assets/images/", http.FileServer(http.Dir("assets/images"))))
	http.Handle("/assets/js/", http.StripPrefix("/assets/js/", http.FileServer(http.Dir("assets/js"))))

	http.ListenAndServe(":8080", nil)
}

var filtres = "home"

func filter(HomeStruct *PageHome) {
	if filtres == "home" {
		HomeStruct.Posts = database.GetAllPost()
	} else if filtres == "films" {
		HomeStruct.Posts = database.GetTagFilm()
	} else if filtres == "series" {
		HomeStruct.Posts = database.GetTagSerie()
	} else if filtres == "Plikes" {
		HomeStruct.Posts = database.SelectByDescending("nbrLikes")
	} else if filtres == "Pdislikes" {
		HomeStruct.Posts = database.SelectByDescending("NbrDislikes")
	} else if filtres == "Mlikes" {
		HomeStruct.Posts = database.SelectByAscending("nbrLikes")
	} else if filtres == "Mdislikes" {
		HomeStruct.Posts = database.SelectByAscending("NbrDislikes")
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	err := tmplHome.Execute(w, HomeStruct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func homeTransiHandler(w http.ResponseWriter, r *http.Request) {
	buLikesDislikes := r.FormValue("bulike/dislike")
	if buLikesDislikes != "" {
		idBu := strings.Split(buLikesDislikes, ",")
		if idBu[1] == "like" {
			nbr := database.RecupNbr("nbrLikes", idBu[0])
			nbr += 1
			database.UpdateNbr("nbrLikes", nbr, idBu[0])
			//fmt.Println(idBu[0], "Like")
		} else if idBu[1] == "dislike" {
			nbr := database.RecupNbr("nbrDislikes", idBu[0])
			nbr += 1
			database.UpdateNbr("nbrDislikes", nbr, idBu[0])
			// fmt.Println(idBu[0], "Dislike")
		}
	}

	headerLinks := r.FormValue("link")
	if headerLinks != "" {
		filtres = headerLinks
	}

	BuMenuDeroulant := r.FormValue("BuMenuDeroulant")
	if BuMenuDeroulant != "" {
		filtres = BuMenuDeroulant
	}

	trucs := r.FormValue("creerPost")
	if trucs != "" {
		fmt.Println(trucs)
	}

	filter(&HomeStruct)
	http.Redirect(w, r, "/", http.StatusFound)

}

func postHandler(w http.ResponseWriter, r *http.Request) {
	buLikesDislikes := r.FormValue("bulike/dislike")
	if buLikesDislikes != "" {
		idBu := strings.Split(buLikesDislikes, ",")
		if idBu[1] == "like" {
			nbr := database.RecupNbr("nbrLikes", idBu[0])
			nbr += 1
			database.UpdateNbr("nbrLikes", nbr, idBu[0])
			// fmt.Println(idBu[0], "Like")
		} else if idBu[1] == "dislike" {
			nbr := database.RecupNbr("nbrDislikes", idBu[0])
			nbr += 1
			database.UpdateNbr("nbrDislikes", nbr, idBu[0])
			// fmt.Println(idBu[0], "Dislike")
		}
	}

	idPost := r.FormValue("buPost")
	if idPost == "retour" {
		PostStruct.IdPost = ""
		http.Redirect(w, r, "/ho", http.StatusFound)
	} else if idPost == "Envoyer" {
		rep := r.FormValue("response")
		database.DatabaseAndReponse([]string{PostStruct.IdPost, PostStruct.ConnectUserName, rep, transformDate(), PostStruct.ConnectUserImg})
	}
	if idPost != "Envoyer" && idPost != "" {
		PostStruct.IdPost = idPost
	}
	PostStruct.OnePost = database.GetOnePost(PostStruct.IdPost)
	PostStruct.Responses = database.GetResponses(PostStruct.IdPost)
	err := tmplPost.Execute(w, PostStruct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func newPostHandler(w http.ResponseWriter, r *http.Request) {
	err := tmplNewPost.Execute(w, NewPostStruct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func newPostTransiHandler(w http.ResponseWriter, r *http.Request) {
	buton := r.FormValue("buCreerPost")
	if buton == "creer" {
		titre := r.FormValue("topic-name")
		tag := r.FormValue("category")
		description := r.FormValue("content")
		database.DatabaseAndPost([]string{HomeStruct.ConnectUserInfo, tag, titre, description, strconv.Itoa(0), strconv.Itoa(0), transformDate(), "117902422"})
	}
	filtres = "home"
	filter(&HomeStruct)
	http.Redirect(w, r, "/", http.StatusFound)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Connection -----------------------------------------------------------------
	connectionEmail := r.FormValue("ConnectionEmail")
	connectionMdp := r.FormValue("ConnectionMdp")
	if connectionEmail != "" {
		user := database.GetUser(connectionEmail)
		if database.GetSession(user.Email) {
			if user.Password == HashPassword2(connectionMdp) && user.Email == connectionEmail {
				// Generate a new session ID
				sessionID := uuid.New().String()
				// Set the session ID as a cookie with an expiration date
				expiration := time.Now().Add(1 * time.Hour) // Session expires after 1 minute
				cookie := &http.Cookie{
					Name:     "Session",
					Value:    sessionID,
					Expires:  expiration,
					HttpOnly: true,
				}
				http.SetCookie(w, cookie)
				database.DatabaseAndSession([]string{connectionEmail, sessionID})
				HomeStruct.ConnectUserInfo = user.Username
				HomeStruct.ConnectUserImg = user.Img
				HomeStruct.IsConnecter = true
				HomeStruct.Error = 0

				PostStruct.ConnectUserName = user.Username
				PostStruct.ConnectUserImg = user.Img

				NewPostStruct.ConnectUserImg = user.Img
			} else {
				HomeStruct.ConnectUserInfo = ""
				HomeStruct.ConnectUserImg = ""
				HomeStruct.IsConnecter = false
				HomeStruct.Error = 2

				PostStruct.ConnectUserName = ""
				PostStruct.ConnectUserImg = ""

				NewPostStruct.ConnectUserImg = ""
			}
		} else {
			HomeStruct.ConnectUserInfo = ""
			HomeStruct.ConnectUserImg = ""
			HomeStruct.IsConnecter = false
			HomeStruct.Error = 1

			PostStruct.ConnectUserName = ""
			PostStruct.ConnectUserImg = ""

			NewPostStruct.ConnectUserImg = ""
		}

	}

	// Inscription ----------------------------------------------------------------
	inscriptionName := r.FormValue("InscriptionName")
	inscriptionEmail := r.FormValue("InscriptionEmail")
	inscriptionMdp := r.FormValue("InscriptionMdp")

	if inscriptionName != "" {
		if !database.GetEmail(inscriptionEmail) {
			img := renderImg()
			database.DatabaseAndUsers([]string{inscriptionEmail, inscriptionName, HashPassword2(inscriptionMdp), img})
			// Generate a new session ID
			sessionID := uuid.New().String()
			// Set the session ID as a cookie with an expiration date
			expiration := time.Now().Add(1 * time.Hour) // Session expires after 1 minute
			cookie := &http.Cookie{
				Name:     "Session",
				Value:    sessionID,
				Expires:  expiration,
				HttpOnly: true,
			}
			http.SetCookie(w, cookie)
			database.DatabaseAndSession([]string{inscriptionEmail, sessionID})
			HomeStruct.ConnectUserInfo = inscriptionName
			HomeStruct.ConnectUserImg = img
			HomeStruct.IsConnecter = true
			HomeStruct.Error = 0
			PostStruct.ConnectUserName = inscriptionName
			PostStruct.ConnectUserImg = img
			NewPostStruct.ConnectUserImg = img
		} else {
			HomeStruct.Error = 3
		}
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Remove the session from the server-side sessions map
	sessionCookie, err := r.Cookie("Session")
	if err == nil {
		database.DeleteSession(sessionCookie.Value)
	}
	HomeStruct.ConnectUserInfo = ""
	HomeStruct.IsConnecter = false
	PostStruct.ConnectUserName = ""
	PostStruct.ConnectUserImg = ""
	NewPostStruct.ConnectUserImg = ""

	// Clear the session cookie
	cookie := &http.Cookie{
		Name:    "session",
		Value:   "",
		Expires: time.Now().Add(-time.Hour),
	}
	http.SetCookie(w, cookie)

	// Redirect to the login page after logging out
	http.Redirect(w, r, "/", http.StatusFound)
}

// Utilitaires *************************************************************************************************************************

func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func transformDate() string {
	Mois := []string{"", "Janvier ", "Février ", "Mars ", "Avril ", "Mai ", "Juin ", "Juillet ", "Août ", "Septembre ", "Octobre ", "Novembre ", "Décembre "}
	currentTime := time.Now()
	dateString := currentTime.Format("2006-01-02 ")
	data := strings.Split(dateString, "-")
	dateFinale := ""
	for i := 0; i < len(data); i++ {
		if i == 1 {
			if data[i][0] == 0 {
				dateFinale = Mois[int(data[i][0])] + dateFinale
			} else {
				indMois, _ := strconv.Atoi(data[i])
				dateFinale = Mois[indMois] + dateFinale
			}
		} else if i == 2 {
			if data[i][0] == 0 {
				dateFinale = string(data[i][1]) + dateFinale
			} else {
				dateFinale = data[i] + dateFinale
			}
		} else {
			dateFinale = data[i] + dateFinale
		}
		dateFinale = dateFinale + " "
	}
	return dateFinale
}

func renderImg() string {
	img := "1179024"
	nb := rand.Intn(15)
	boolean := true
	for boolean {
		if (nb != 0) && (nb != 1) {
			boolean = false
		} else {
			nb = rand.Intn(15)
		}
	}
	if nb < 10 {
		img += "0" + strconv.Itoa(nb)
	} else {
		img += strconv.Itoa(nb)
	}
	return img
}

func HashPassword2(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	hashed := h.Sum(nil)
	return string(hashed)
}
