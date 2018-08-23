package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID   bson.ObjectId `bson:"_id"`
	Name string        `bson:"userName"`
}

type Pod struct {
	ID    bson.ObjectId `bson:"_id"`
	User  bson.ObjectId `bson:"user"`
	Name  string        `bson:"name"`
	Users []User        `bson:"users" json:"users"`
}

func main() {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}

	defer session.Close()
	db := session.DB("pratice")

	user := User{
		Name: "YoYo",
	}

	pod := Pod{
		Name: "MyPod",
	}
	pod.ID = bson.NewObjectId()
	user.ID = bson.NewObjectId()
	pod.User = user.ID
	err = db.C("pods").Insert(pod)
	err = db.C("users").Insert(user)
	defer db.C("pods").RemoveId(pod.ID)
	defer db.C("users").RemoveId(user.ID)
	fmt.Println("Try to lookup")
	pipeline := []bson.M{
		//bson.M{"$lookup": bson.M{"from": "users", "localField": "user", "foreignField": "_id", "as": "users"}},
		//	bson.M{"$match": bson.M{"name": "MyPod"}},
		{"$lookup": bson.M{"from": "users", "localField": "user", "foreignField": "_id", "as": "users"}},
		{"$match": bson.M{"name": "MyPod"}},
	}

	var resp Pod
	err = db.C("pods").Pipe(pipeline).One(&resp)
	fmt.Printf("%v\n %+v", err, resp)
}
