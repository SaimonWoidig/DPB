package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

const ScoreboardKey = "scoreboard"

const MinScore = 0
const MaxScore = 999

var MinScoreStr = strconv.Itoa(MinScore)
var MaxScoreStr = strconv.Itoa(MaxScore)

func main() {
	ctx := context.Background()

	// create client
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	// check connectivity
	if _, err := client.Ping(ctx).Result(); err != nil {
		panic(err.Error())
	}
	// deferred function that will clear redis and close connection
	defer func() {
		if _, err := client.FlushAll(ctx).Result(); err != nil {
			panic(err.Error())
		}
		if err := client.Close(); err != nil {
			panic(err.Error())
		}
	}()

	// add Alfred with score 888
	client.ZAdd(ctx, ScoreboardKey, redis.Z{
		Score:  888,
		Member: "Alfred",
	})

	// add other ten players
	players := []redis.Z{
		{Score: 123, Member: "Tom"},
		{Score: 111, Member: "Bob"},
		{Score: 222, Member: "Alice"},
		{Score: 333, Member: "Theresa"},
		{Score: 444, Member: "Jim"},
		{Score: 555, Member: "Tim"},
		{Score: 666, Member: "Martin"},
		{Score: 777, Member: "Joanna"},
		{Score: 889, Member: "Garfield"},
		{Score: 999, Member: "Maurice"},
	}
	client.ZAdd(ctx, ScoreboardKey, players...)

	// top 3 by score
	topThree, err := client.ZRevRangeByScore(ctx, ScoreboardKey, &redis.ZRangeBy{
		Min:    MinScoreStr,
		Max:    MaxScoreStr,
		Offset: 0,
		Count:  3,
	}).Result()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Top 3:")
	for pos, player := range topThree {
		fmt.Printf("%v. %v\n", pos+1, player)
	}

	// worst score
	worstPlayers, err := client.ZRangeByScoreWithScores(ctx, ScoreboardKey, &redis.ZRangeBy{
		Min:    MinScoreStr,
		Max:    MaxScoreStr,
		Offset: 0,
		Count:  1,
	}).Result()
	if err != nil {
		panic(err.Error())
	}
	worstPlayer := worstPlayers[0]
	fmt.Printf("The worst score is: %v\n", worstPlayer.Score)

	// number of players with score < 100 (noobies)
	numPlayersSub100, err := client.ZCount(ctx, ScoreboardKey, MinScoreStr, "100").Result()
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Number of noobies: %v\n", numPlayersSub100)

	// players with score > 850 (pros)
	playersOver850, err := client.ZRangeByScore(ctx, ScoreboardKey, &redis.ZRangeBy{
		Min: "850",
		Max: MaxScoreStr,
	}).Result()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("The pros:")
	for _, player := range playersOver850 {
		fmt.Printf("- %v\n", player)
	}

	// get Alfred's position
	alfredsPos, err := client.ZRevRank(ctx, ScoreboardKey, "Alfred").Result()
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Alfred's position now: %v\n", alfredsPos+1)

	// increment Alfred's score and check position
	if _, err = client.ZIncrBy(ctx, ScoreboardKey, 12, "Alfred").Result(); err != nil {
		panic(err.Error())
	}
	fmt.Println("Alfred played some games...")
	alfredsPos, err = client.ZRevRank(ctx, ScoreboardKey, "Alfred").Result()
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Alfred's position now: %v\n", alfredsPos+1)
}
