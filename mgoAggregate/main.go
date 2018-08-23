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

func pipeline(c *mgo.Collection, user User) {
	start := time.Now()

	pipeline := []bson.M{
		{"$lookup": bson.M{"from": "pods", "localField": "_id", "foreignField": "createdBy", "as": "pods"}},
		{"$match": bson.M{"_id": user.ID}},
	}

	var resp User
	err := c.Pipe(pipeline).One(&resp)
	if err != nil {
		fmt.Println(err)
	}
	elapsed := time.Since(start)
	fmt.Printf("Pipeline took %s: %d\n", elapsed, len(resp.Pods))
}

func find(c *mgo.Collection, user User) {
	start := time.Now()

	c.Find(bson.M{"createdBy": user.ID}).All(&user.Pods)

	elapsed := time.Since(start)
	fmt.Printf("Find took %s: %d\n", elapsed, len(user.Pods))
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
	for i := 0; i < 30000; i++ {
		pod := Pod{
			ID:        bson.NewObjectId(),
			Name:      namesgenerator.GetRandomName(0),
			CreatedBy: user.ID,
		}
		db.C("pods").Insert(pod)
	}
	defer db.C("pods").DropCollection()

	pipeline(db.C("users"), user)
	find(db.C("pods"), user)
}
