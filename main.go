package mongo31

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// подключение к СУБД MongoDB
	mongoOpts := options.Client().ApplyURI("mongodb://0.0.0.0:27017/")
	client, err := mongo.Connect(context.Background(), mongoOpts)
	if err != nil {
		log.Fatal(err)
	}
	// не забываем закрывать ресурсы
	defer client.Disconnect(context.Background())
	// проверка связи с БД
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	type lang struct {
		Id   int
		Name string
	}

	langs := []lang{{6, "C++"}, {7, "Java"}}

	coll := client.Database("data").Collection("language")

}
