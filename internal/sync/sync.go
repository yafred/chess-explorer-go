package sync

import (
	"context"
	"log"
	"time"

	"github.com/spf13/viper"
	"github.com/yafred/chess-explorer/internal/chesscom"
	"github.com/yafred/chess-explorer/internal/lichess"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type user struct {
	Site     string `json:"site,omitempty"`
	Username string `json:"username,omitempty"`
}

// All ... Download recent games for all users in database
func All() {
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
		log.Fatal("Cannot connect to DB " + viper.GetString("mongo-url"))
	}

	// Gather names of users whose games we must not delete
	lastgamesCollection := client.Database(viper.GetString("mongo-db-name")).Collection("lastgames")
	findOptions := options.Find().SetProjection(bson.M{"site": 1, "username": 1})
	cursor, err := lastgamesCollection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	var users []user
	if err = cursor.All(ctx, &users); err != nil {
		log.Fatal(err)
	}

	// Call the right download command in a sequence
	for _, user := range users {
		log.Println("Synchronizing", user.Username, " (", user.Site, ")")
		switch user.Site {
		case "lichess.org":
			lichess.DownloadGames(user.Username, "")
			break
		case "chess.com":
			chesscom.DownloadGames(user.Username, "")
			break
		default:
			// Do nothing
		}
	}

}
