package delete

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type user struct {
	Site     string `json:"site,omitempty"`
	Username string `json:"username,omitempty"`
}
type game struct {
	ID string `json:"_id,omitempty"`
}

// Games ... Delete games for user {username} or lichess.org:{username} or chess.com:{username}
func Games(username string) {
	// process argument
	site := ""

	username = strings.TrimSpace(username)
	if strings.Index(username, ":") != -1 {
		splitUserName := strings.Split(username, ":")
		site = splitUserName[0]
		username = splitUserName[1]
	}

	// Connect to DB
	client, err := mongo.NewClient(options.Client().ApplyURI(viper.GetString("mongo-url")))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Ping MongoDB
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal("Cannot connect to DB")
	}

	// Gather names of users whose games we must not delete
	lastgamesCollection := client.Database("chess-explorer").Collection("lastgames")
	findOptions := options.Find().SetProjection(bson.M{"site": 1, "username": 1})
	cursor, err := lastgamesCollection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	var users []user
	if err = cursor.All(ctx, &users); err != nil {
		log.Fatal(err)
	}

	// Delete games
	andClause := make([]bson.M, 0)

	if site != "" {
		andClause = append(andClause, bson.M{"site": site})
	}

	deleteBson := make([]bson.M, 0)
	deleteBson = append(deleteBson, bson.M{"white": username})
	deleteBson = append(deleteBson, bson.M{"black": username})
	andClause = append(andClause, bson.M{"$or": deleteBson})

	notIn := make([]string, 0)
	for _, user := range users {
		if strings.ToLower(user.Username) != strings.ToLower(username) {
			notIn = append(notIn, user.Username)
		}
	}

	if len(notIn) > 0 {
		andClause = append(andClause, bson.M{"white": bson.M{"$nin": notIn}})
		andClause = append(andClause, bson.M{"black": bson.M{"$nin": notIn}})
	}

	gameFilter := bson.M{}
	switch len(andClause) {
	case 0:
		log.Fatal("Unexpected")
		break
	case 1:
		gameFilter = andClause[0]
		break
	default:
		gameFilter = bson.M{"$and": andClause}
		break
	}

	gamesCollection := client.Database("chess-explorer").Collection("games")

	collation := options.Collation{Locale: "en", Strength: 2}
	deleteOptions := options.DeleteOptions{Collation: &collation} // case insensitive search

	_, err = gamesCollection.DeleteMany(ctx, gameFilter, &deleteOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Delete user
	deleteUsersFilter := bson.M{"username": username}
	if site != "" {
		deleteUsersFilter = bson.M{"username": username, "site": site}
	}
	_, err = lastgamesCollection.DeleteMany(ctx, deleteUsersFilter, &deleteOptions)
	if err != nil {
		log.Fatal(err)
	}

}
