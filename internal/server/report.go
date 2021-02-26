package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/yafred/chess-explorer/internal/pgntodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type result struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type userResult struct {
	SiteName string `json:"sitename"`
	Name     string `json:"name"`
	Count    int    `json:"count"`
}

type report struct {
	TotalGames   int64 `json:"totalgames,omitempty"`
	Sites        []result
	Users        []userResult
	UsersAsWhite []result
	TimeControls []result
}

type reportResponse struct {
	Error string `json:"error"`
	Data  report `json:"data"`
}

var filter GameFilter

func reportHandler(w http.ResponseWriter, r *http.Request) {

	defer timeTrack(time.Now(), "reportHandler")

	// allow cross origin
	w.Header().Set("Access-Control-Allow-Origin", "*")

	filter.white = strings.TrimSpace(r.FormValue("white"))
	filter.black = strings.TrimSpace(r.FormValue("black"))
	filter.from = strings.TrimSpace(r.FormValue("from"))
	filter.to = strings.TrimSpace(r.FormValue("to"))

	response := reportResponse{}

	// Connect to DB
	client, err := mongo.NewClient(options.Client().ApplyURI(viper.GetString("mongo-url")))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

	games := client.Database(viper.GetString("mongo-db-name")).Collection("games")
	lastgames := client.Database(viper.GetString("mongo-db-name")).Collection("lastgames")

	// Total games
	totalGames, error := games.CountDocuments(ctx, bson.M{})
	if error != nil {
		log.Fatal(error)
	}
	report := report{}
	report.TotalGames = totalGames

	if filter.black == "" && filter.white == "" {
		reportGames(ctx, games, &report)
		reportSites(ctx, games, &report)
		reportUsers(ctx, games, lastgames, &report)
		reportUsersAsWhite(ctx, games, &report)
		reportTimeControls(ctx, filter, games, &report)
	} else {
		reportTimeControls(ctx, filter, games, &report)
	}

	// send the response
	response.Data = report
	json.NewEncoder(w).Encode(response)
}

// Games
func reportGames(ctx context.Context, games *mongo.Collection, report *report) {
	totalGames, error := games.CountDocuments(ctx, bson.M{})
	if error != nil {
		log.Fatal(error)
	}
	report.TotalGames = totalGames
}

// Sites
func reportSites(ctx context.Context, games *mongo.Collection, report *report) {
	filter := bson.M{"$match": bson.M{}}
	pipeline := make([]bson.M, 0)
	pipeline = append(pipeline, filter)

	groupStage := bson.M{
		"$group": bson.M{
			"_id":   bson.M{"site": "$site"},
			"count": bson.M{"$sum": 1},
		},
	}
	pipeline = append(pipeline, groupStage)

	sortStage := bson.M{
		"$sort": bson.M{
			"count": -1,
		},
	}
	pipeline = append(pipeline, sortStage)

	projectStage := bson.M{
		"$project": bson.M{
			"_id":   false,
			"name":  "$_id.site",
			"count": "$count",
		},
	}
	pipeline = append(pipeline, projectStage)

	aggregateCursor, err := games.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
	}

	defer aggregateCursor.Close(ctx)

	var siteResults []result
	if err = aggregateCursor.All(ctx, &siteResults); err != nil {
		log.Fatal(err)
	}

	report.Sites = siteResults
}

// Users
func reportUsers(ctx context.Context, games *mongo.Collection, lastgames *mongo.Collection, report *report) {
	cursor, err := lastgames.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	var results []pgntodb.LastGame
	if err = cursor.All(ctx, &results); err != nil {
		log.Fatal(err)
	}

	report.Users = make([]userResult, 0)
	for _, aUser := range results {
		filter := bson.M{"site": aUser.Site, "$or": []bson.M{{"white": aUser.Username}, {"black": aUser.Username}}}
		count, err := games.CountDocuments(ctx, filter)
		if err != nil {
			log.Fatal(err)
		}
		report.Users = append(report.Users, userResult{SiteName: aUser.Site, Name: aUser.Username, Count: int(count)})
	}
}

// Users as white
func reportUsersAsWhite(ctx context.Context, games *mongo.Collection, report *report) {
	filter := bson.M{"$match": bson.M{}}
	pipeline := make([]bson.M, 0)
	pipeline = append(pipeline, filter)

	groupStage := bson.M{
		"$group": bson.M{
			"_id":   bson.M{"white": "$white"},
			"count": bson.M{"$sum": 1},
		},
	}
	pipeline = append(pipeline, groupStage)

	reduceStage := bson.M{"$match": bson.M{
		"count": bson.M{"$gte": 10},
	}}
	pipeline = append(pipeline, reduceStage)

	sortStage := bson.M{
		"$sort": bson.M{
			"count": -1,
		},
	}
	pipeline = append(pipeline, sortStage)

	projectStage := bson.M{
		"$project": bson.M{
			"_id":   false,
			"name":  "$_id.white",
			"count": "$count",
		},
	}
	pipeline = append(pipeline, projectStage)

	aggregateCursor, err := games.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
	}

	defer aggregateCursor.Close(ctx)

	var usersAsWhiteResult []result
	if err = aggregateCursor.All(ctx, &usersAsWhiteResult); err != nil {
		log.Fatal(err)
	}

	report.UsersAsWhite = usersAsWhiteResult
}

// Time controls
func reportTimeControls(ctx context.Context, gameFilter GameFilter, games *mongo.Collection, report *report) {
	filter := bson.M{"$match": processGameFilter(gameFilter)}
	pipeline := make([]bson.M, 0)
	pipeline = append(pipeline, filter)

	groupStage := bson.M{
		"$group": bson.M{
			"_id":   bson.M{"timecontrol": "$timecontrol"},
			"count": bson.M{"$sum": 1},
		},
	}
	pipeline = append(pipeline, groupStage)

	sortStage := bson.M{
		"$sort": bson.M{
			"count": -1,
		},
	}
	pipeline = append(pipeline, sortStage)

	projectStage := bson.M{
		"$project": bson.M{
			"_id":   false,
			"name":  "$_id.timecontrol",
			"count": true,
		},
	}
	pipeline = append(pipeline, projectStage)

	aggregateCursor, err := games.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
	}

	defer aggregateCursor.Close(ctx)

	var timeControlResults []result
	if err = aggregateCursor.All(ctx, &timeControlResults); err != nil {
		log.Fatal(err)
	}

	report.TimeControls = timeControlResults
}
