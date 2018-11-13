package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type Session struct {
	session *mgo.Session
}

type Conf struct {
	MongodbConnString string
	DbName            string
	Pwd               string
}

var conf *Conf

var mgoSession *mgo.Session

// GetMongoSession connect new session if session not exist, send cloned version of session
func GetMongoSession() *mgo.Session {

	if mgoSession == nil {
		log.Println("Creating new session")
		connString := conf.MongodbConnString
		dialInfo, err := mgo.ParseURL(connString)
		if conf.Pwd != "" {
			dialInfo.Password = conf.Pwd // override password form conn string
		}
		if err != nil {
			log.Println(err)
		}
		/*
			dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
				log.Println("Loading TLS config")
				tlsConfig := &tls.Config{}
				conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
				if err != nil {
					log.Println(err)
				} else {
					log.Println("Got connection to MongoDB")
				}
				return conn, err
			}
		*/
		log.Println("About to DialWithInfo")
		session, err := mgo.DialWithInfo(dialInfo)
		if err != nil {
			panic(err)
		} else {
			log.Println("Got session to MongoDB")
		}

		mgoSession = session
	}

	log.Println("Cloning mgoSession")
	return mgoSession.Clone()
}

// GetSessionDatabase ...
func GetSessionDatabase(dbName string) (*mgo.Session, *mgo.Database) {
	log.Println("Calling GetSessionDatabase")
	session := GetMongoSession()
	// session.SetMode(mgo.Monotonic, true)
	db := session.DB(dbName)
	return session, db
}

func main() {

	conf = &Conf{}
	conf.MongodbConnString = "mongodb://localhost:27017/test"
	conf.DbName = "test"
	conf.Pwd = ""

	adId := 1
	for adId < 100 {
		result := getAdId(123)
		log.Println("Ad for id=" + string(adId) + "=" + result.name)
		adId += adId
	}
	fmt.Println(adId)

}

type Ad struct {
	name string
}

func getAdId(id int) Ad {
	log.Println("Calling getAdId")
	session, db := GetSessionDatabase(conf.DbName)
	defer session.Close()
	var result Ad

	err := db.C("ad").Find(bson.M{"ad_id": id}).Select(bson.M{"_id": 0}).One(&result)
	if err != nil {
		log.Println("Error doing find")
		log.Println(err)
	}
	return result
}
