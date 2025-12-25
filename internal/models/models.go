package models

import (
	"time"
)

type User struct {
    ID          int       `json:"id"`
	Username    string    `json:"username"`
	Password	string	  `json:"password"`
	Age         int       `json:"age"`
    Description string    `json:"description"`
    MostLikeGame       string    `json:"mostlikegame"`
    MostLikeGenre       string    `json:"mostlikegenre"`
	App			string 	  `json:"app"`
	Language			string 	  `json:"language"`
    CreatedAt   time.Time `json:"created_at"`
}

type UserUpdate struct {
    ID          int       `json:"id"`
	Age         int       `json:"age"`
    Description string    `json:"description"`
    MostLikeGame        string    `json:"mostlikegame"`
    MostLikeGenre       string    `json:"mostlikegenre"`
	App			string 	  `json:"app"`
	Language	string 	  `json:"language"`
}

type UserListShow struct {
	Username         string       `json:"username"`
	Age         int       `json:"age"`
    Description string    `json:"description"`
    MostLikeGame        string    `json:"mostlikegame"`
    MostLikeGenre       string    `json:"mostlikegenre"`
	Language	string 	  `json:"language"`
}

type FilterData struct {
	Age int `json:"age"`
	Game string `json:"game"`
	Genre string `json:"genre"`
	Language string `json:"language"`
	App string `json:"app"`
}

type UpdateUserData struct {
	UserUpdateOld UserUpdate
	UserUpdateNew UserUpdate
}


type Language struct {
    ID          int       `json:"id_language"`
	Lang		string 	  `json:"language"`
}

type Genres struct {
    ID          int       `json:"id_genre"`
	Genre		string 	  `json:"genre"`
}

type Games struct {
    ID          int       `json:"id_game"`
	Game		string 	  `json:"game"`
}

type Apps struct {
    ID          int       `json:"id_app"`
	App		string 	  `json:"app"`
}

