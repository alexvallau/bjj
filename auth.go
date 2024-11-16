package main

import(
	"log"
	"golang.org/x/crypto/bcrypt"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"time"
)

type  Utilisateur struct{
	Id	int `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}


func (u *Utilisateur) Login()(bool, error){
	db, err := connectDB()
	if err != nil{
		log.Panic(err)
	}
	defer db.Close()
	var hashedPassword string
	query := "SELECT password FROM utilisateur WHERE username = ?"
	err = db.QueryRow(query, u.Username).Scan(&hashedPassword)
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(u.Password))
	if err != nil {
		fmt.Printf("Failed Connexion with username %s and password %s \n", u.Username, u.Password)
		return false, err
	}
	//Maintenant que le user est authentifié, on attribu son ID DB à son ID struct
	queryIdUser := "SELECT id FROM utilisateur WHERE username = ?"
	err = db.QueryRow(queryIdUser, u.Username).Scan(&u.Id)
	if err != nil {
		log.Fatal(err)
		return false, err
	}

	fmt.Printf("User %s correctly logged in", u.Username)
	return true, nil
}

func CreateUser(username, password string) {

	db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	hashedPassword, err := hashPassword(password)

	if err != nil {
		log.Fatal(err)
	}

	query := "INSERT INTO utilisateur (username, password) VALUES (?, ?)"
	_, err = db.Exec(query, username, hashedPassword)
	if err != nil {
		log.Fatal(err)
	}
	
}

func hashPassword(UnhashedPassword string) (hashedPassword []byte, err error) {
	return bcrypt.GenerateFromPassword([]byte(UnhashedPassword), bcrypt.DefaultCost)
}

func (u *Utilisateur) GenerateJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": u.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtKey, err := GetVarEnv("JWT_SECRET_KEY")
	if err != nil {
		log.Panic(err)
		return "", err
	}

	return token.SignedString(jwtKey)
}
func VerifyToken(tokenString string) (int,error) {

	secretKey, err := GetVarEnv("JWT_SECRET_KEY")
	if err != nil {
		log.Panic(err)
		return -1, err
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return -1, err
	}

	if !token.Valid {
		fmt.Println("MOtherfucker")
		return -1, err
	}
	claims := token.Claims.(jwt.MapClaims)
	fmt.Println("claims id usre", claims["user_id"])
	user_id := claims["user_id"]
	return int(user_id.(float64)), nil
}

func GetVarEnv(name string) ([]byte, error) {
	envFile, err := godotenv.Read(".env")
	if err != nil {
		return nil, err
	}
	return []byte(envFile[name]), nil
}

