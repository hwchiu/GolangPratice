package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/moby/moby/pkg/namesgenerator"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID   bson.ObjectId `bson:"_id"`
	Name string        `bson:"name"`
	Pods []Pod         `bson:"pods"`
}

type Pod struct {
	ID        bson.ObjectId `bson:"_id"`
	CreatedBy bson.ObjectId `bson:"createdBy"`
	Name      string        `bson:"name"`
}

//For namesgenerator
func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	db := session.DB("pratice")

	//Create User
	user := User{
		Name: namesgenerator.GetRandomName(0),
		ID:   bson.NewObjectId(),
	}

	db.C("users").Insert(user)
	defer db.C("users").RemoveId(user.ID)
	//Create Pods
	for i := 0; i <= 5; i++ {
		pod := Pod{
			ID:        bson.NewObjectId(),
			Name:      namesgenerator.GetRandomName(0),
			CreatedBy: user.ID,
		}
		db.C("pods").Insert(pod)
		defer db.C("pods").Remove(bson.M{"name": pod.Name})
	}

	pipeline := []bson.M{
		{"$lookup": bson.M{"from": "pods", "localField": "_id", "foreignField": "createdBy", "as": "pods"}},
		{"$match": bson.M{"_id": user.ID}},
	}

	var resp User
	db.C("users").Pipe(pipeline).One(&resp)
	for _, v := range resp.Pods {
		fmt.Printf("%+v\n", v)
	}
}
