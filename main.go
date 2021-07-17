package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	workDB         = "data"
	workCollection = "languages"
)

type lang struct {
	Id   int
	Name string
}

func main() {
	pwd := os.Getenv("Cloud0pass")
	connstr := fmt.Sprintf("mongodb+srv://sup:%s@cloud0.wspoq.mongodb.net/posts?retryWrites=true&w=majority", pwd)

	// подключение к СУБД MongoDB в облаке
	clientOptions := options.Client().
		ApplyURI(connstr)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// // подключение к СУБД MongoDB
	// mongoOpts := options.Client().ApplyURI("mongodb://0.0.0.0:27017/")
	// client, err := mongo.Connect(context.Background(), mongoOpts)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// не забываем закрывать ресурсы
	defer client.Disconnect(context.Background())
	// проверка связи с БД
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	langs := []lang{{6, "C++"}, {7, "Java"}}

	err = insertPack(client, langs)

	if err != nil {
		log.Fatal(err)
	}

	langs, err = languages(client)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(langs)
}

func insertPack(c *mongo.Client, docs []lang) error {
	coll := c.Database(workDB).Collection(workCollection)

	for _, doc := range docs {
		_, err := coll.InsertOne(context.Background(), doc)
		if err != nil {
			return err
		}

	}
	return nil
}

func languages(c *mongo.Client) ([]lang, error) {
	coll := c.Database(workDB).Collection(workCollection)
	ctx := context.Background()
	filter := bson.D{}
	cur, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var ls []lang
	for cur.Next(ctx) {
		var l lang
		err = cur.Decode(&l)
		if err != nil {
			return nil, err
		}
		ls = append(ls, l)
	}
	return ls, cur.Err()
}
