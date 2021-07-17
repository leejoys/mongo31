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
	workDB         = "gonews"
	workCollection = "posts"
)

type Post struct {
	ID      int    // номер записи
	Title   string // заголовок публикации
	Content string // содержание публикации
	PubTime int64  // время публикации
	Link    string // ссылка на источник
}

func main() {
	pwd := os.Getenv("Cloud0pass")
	connstr := fmt.Sprintf("mongodb+srv://sup:%s@cloud0.wspoq.mongodb.net/gonews?retryWrites=true&w=majority", pwd)

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

	posts := []Post{{1, "Вышел Microsoft Linux",
		"Как сообщают непроверенные источники, новая ОС будет бесплатной.",
		time.Now().Unix(), "https://github.com/microsoft/CBL-Mariner"},
		{2, "Инженеры Google не желают возвращаться в офисы",
			"Инженеры Google не желают возвращаться в офисы, заявляя, что они не менее продуктивны на удалёнке.",
			time.Now().Unix(), "https://habr.com/ru/news/t/568128/"}}

	err = insertPack(client, posts)

	if err != nil {
		log.Fatal(err)
	}

	posts, err = languages(client)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(posts)

	for _, p := range posts {
		err := deletePost(client, p)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func insertPack(c *mongo.Client, docs []Post) error {
	coll := c.Database(workDB).Collection(workCollection)

	for _, doc := range docs {
		_, err := coll.InsertOne(context.Background(), doc)
		if err != nil {
			return err
		}

	}
	return nil
}

func languages(c *mongo.Client) ([]Post, error) {
	coll := c.Database(workDB).Collection(workCollection)
	ctx := context.Background()
	filter := bson.D{}
	cur, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var ps []Post
	for cur.Next(ctx) {
		var l Post
		err = cur.Decode(&l)
		if err != nil {
			return nil, err
		}
		ps = append(ps, l)
	}
	return ps, cur.Err()
}

//DeletePost - удаляет пост по id
func deletePost(c *mongo.Client, p Post) error {
	coll := c.Database(workDB).Collection(workCollection)
	filter := bson.D{{Key: "id", Value: p.ID}}
	_, err := coll.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}
