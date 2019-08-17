package database

import (
	"fmt"
	"github.com/globalsign/mgo"
	"log"
	"strings"
)

var session *mgo.Session
var Database *mgo.Database

func InitDB() {
	session, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		log.Fatal(err)
	}

	session.SetMode(mgo.Monotonic, true)

	Database = session.DB("giggle-test")
}

func CloseSession() {
	session.Close()
}

func GetCollection(name string) *mgo.Collection {
	if Database == nil {
		InitDB()
	}

	return Database.C(name)
}

func PrepareFilters(filters string) map[string]interface{} {
	var filterTypes []string
	var filtersListOfEachType = make(map[string]interface{})
	filterTypes = strings.Split(filters, ", ")
	if len(filterTypes) == 1 {
		filterTypes = strings.Split(filters, ",")
	}

	fmt.Println(filterTypes)

	for _, v := range filterTypes {
		if string(v[0]) == " " {
			strings.TrimSpace(v)
		}
		d := strings.Split(v, ":")
		if string(d[1][0]) == "-" {
			filtersListOfEachType[d[0]] = map[string]string{"$ne": d[1][1:]}
		} else {
			filtersListOfEachType[d[0]] = d[1]
		}
	}
	return filtersListOfEachType
}
