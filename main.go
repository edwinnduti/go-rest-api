/*
[*] Copyright © 2020
[*] Dev/Author ->  Edwin Nduti
[*] Description:
    The code is a REST API using mgo,mongo golang driver.
    Written in Golang.
*/

//START CODE
package main

//The necessary and needed libraries
import(
    "os"
    "fmt"
    "log"
    "time"
    "image"
    "bytes"
    "net/http"
    "image/jpeg"
    "encoding/json"
    "github.com/gorilla/mux"
    "github.com/urfave/negroni"
    go_qrcode "github.com/skip2/go-qrcode"

    mgo "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"

)

//User struct
type User struct{
    ID              bson.ObjectId      `bson:"_id" json:"id"`
    UserName        string             `json:"username"`
    Email           string              `json:"email"`
    Password        string             `json:"password"`
    DOB             string             `json:"dob"`
    Qrcode          []byte             `json:"qrcode"`
    CreatedAt       time.Time          `json:"createdat"`
    UpdatedAt       time.Time          `json:"lastupdatedat"`
}

// database and collection names are statically declared
const database, collection = "appservice", "user"

// DB connection
func connect() *mgo.Session {
    session, err := mgo.Dial("localhost")
    if err != nil {
        fmt.Println("session err:", err)
        os.Exit(1)
    }

    session.SetMode(mgo.Monotonic, true)

    return session
}

// save data
func Save(data interface{}) error {
    session := connect()

    defer session.Close()

    err := session.DB(database).C(collection).Insert(data)
    return err
}

//handle errors
func Check(e error){
    if e != nil{
        log.Fatalln(e)
    }
}

// Generating qrcode
func GenerateQrcode(dataString User) ([]byte,error) {

    out,err := json.Marshal(dataString)
    Check(err)

    values := fmt.Sprintf(string(out))

    var filename []byte
    filename , err = go_qrcode.Encode(values, go_qrcode.Medium, 256)
    Check(err)

    return filename,nil
}

// HTTP /POST /api
func PostHandler(w http.ResponseWriter,r *http.Request){
    var user User

    // Decode the incoming Data json
    err := json.NewDecoder(r.Body).Decode(&user)
    Check(err)

    user.ID = bson.NewObjectId()
    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()
    user.Qrcode,err = GenerateQrcode(user)
    Check(err)
    err = Save(user)
    Check(err)

    j, err := json.Marshal(user)
    Check(err)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    w.Write(j)
}

// HTTP /GET user record /api/{id}
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    obj_id := vars["id"]

    user := User{}
    session := connect()
    defer session.Close()

    s := session.DB(database).C(collection)
    err := s.Find(bson.M{"_id": bson.ObjectIdHex(obj_id)}).One(&user)
    Check(err)

    w.Header().Set("Content-Type", "application/json")
    j, err := json.Marshal(user)
    Check(err)

    w.WriteHeader(http.StatusOK)
    w.Write(j)
}


// HTTP /GET single user record /api
func GetAllUserHandler(w http.ResponseWriter, r *http.Request) {
    user := User{}
    users := []User{}
    session := connect()
    defer session.Close()

    s := session.DB(database).C(collection)
    iter := s.Find(nil).Iter()
    for iter.Next(&user) {
        users = append(users,user)
    }
    err := iter.Close()
    Check(err)

    w.Header().Set("Content-Type", "application/json")
    j, err := json.Marshal(users)
    Check(err)

    w.WriteHeader(http.StatusOK)
    w.Write(j)
}

// HTTP /PUT user record /api/{id}
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    obj_id := vars["id"]

    user := User{}
    session := connect()
    defer session.Close()

    var updatedUser User

    // Decode the incoming Data json
    err := json.NewDecoder(r.Body).Decode(&updatedUser)
    Check(err)

    s := session.DB(database).C(collection)

    err = s.Find(bson.M{"_id": bson.ObjectIdHex(obj_id)}).One(&user)
    Check(err)

    user.UserName = updatedUser.UserName
    user.UpdatedAt = time.Now()
    user.Qrcode,err = GenerateQrcode(user)
    Check(err)

    err = s.Update(bson.M{"_id": bson.ObjectIdHex(obj_id)},
        bson.M{"$set": bson.M{
        "username":user.UserName,
        "qrcode":user.Qrcode,
        "lastupdatedat":user.UpdatedAt,
        }})
    Check(err)

    w.Header().Set("Content-Type", "application/json")
    j, err := json.Marshal(user)
      Check(err)

    w.WriteHeader(http.StatusOK)
    w.Write(j)
}

// HTTP /DELETE user record /api/{id}
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    obj_id := vars["id"]

    var user User

    session := connect()
    defer session.Close()

    s := session.DB(database).C(collection)
    err := s.Remove(bson.M{"_id": bson.ObjectIdHex(obj_id)})
    Check(err)

    w.Header().Set("Content-Type", "application/json")
    j, err := json.Marshal(user)
    Check(err)

    w.WriteHeader(http.StatusOK)
    w.Write(j)
}


// HTTP /GET single qrcode /api/{id}/qrcode
func GetQrcodeHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    obj_id := vars["id"]

    user := User{}
    session := connect()
    defer session.Close()

    s := session.DB(database).C(collection)
    err := s.Find(bson.M{"_id": bson.ObjectIdHex(obj_id)}).One(&user)
    Check(err)

    // create image to be displayed as qrcode in the browser
    img, _, err := image.Decode(bytes.NewReader(user.Qrcode))
    Check(err)

    var opts jpeg.Options
    opts.Quality = 1

    w.Header().Set("Content-Type", "image/*")

    w.WriteHeader(http.StatusOK)

    err = jpeg.Encode(w, img, &opts)
    Check(err)
}

//Main Function
func main(){

    /*
    mgo.SetDebug(true)
    mgo.SetLogger(log.New(os.Stdout,"err",6))

    The above two lines are for debugging errors
    that occur straight from accessing the mongo db
     */

    //Register router
    r := mux.NewRouter().StrictSlash(false)

    // API routes,handlers and methods
    r.HandleFunc("/api", PostHandler).Methods("POST")
    r.HandleFunc("/api/{id}", GetUserHandler).Methods("GET")
    r.HandleFunc("/api/{id}", DeleteUserHandler).Methods("DELETE")
    r.HandleFunc("/api/{id}", UpdateUserHandler).Methods("PUT")
    r.HandleFunc("/api", GetAllUserHandler).Methods("GET")
    r.HandleFunc("/api/{id}/qrcode", GetQrcodeHandler).Methods("GET")


    //Get port
    Port := os.Getenv("PORT")
    if Port == "" {
        Port = "8094"
    }

    // establish logger
    n := negroni.Classic()
    n.UseHandler(r)
    server := &http.Server{
     Handler: n,
     Addr   : ":"+Port,
    }
    log.Printf("Listening on PORT: %s",Port)
    server.ListenAndServe()
}
