package pgntodb

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Process ... process a single file or all the files of a folder
func Process(filepath string, lastGame *LastGame) bool {
	goOn := true

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

	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		log.Fatal("Cannot access " + filepath)
	}

	if info.IsDir() {
		fileinfos, err := ioutil.ReadDir(filepath)
		if err != nil {
			log.Fatal("Cannot list files in " + filepath)
		}
		for _, info := range fileinfos {
			if !info.IsDir() {
				log.Println(path.Join(filepath, info.Name()))
				goOn = processFile(path.Join(filepath, info.Name()), client, lastGame)
				if goOn == false {
					break
				}
			}
		}
	} else {
		goOn = processFile(filepath, client, lastGame)
	}

	return goOn
}

// ProcessFile ... does everything
func processFile(filepath string, client *mongo.Client, lastGame *LastGame) bool {

	// Open file
	file, err := os.Open(filepath)
	defer file.Close()

	if err != nil {
		log.Fatal("Cannot open file " + filepath)
	}

	// Do the work
	return pgnFileToDB(file, client, lastGame)
}

// FindLastGame ... find last game (allowing prevention of duplicates)
func FindLastGame(username string, site string) *LastGame {
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

	return findLastGame(username, site, client)
}
