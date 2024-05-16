package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/SaimonWoidig/DPB/cv11/pkg/models"
	"github.com/gocql/gocql"
	"github.com/joho/godotenv"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
)

const (
	EnvCassandraHostKey     = "CASSANDRA_HOST"
	EnvCassandraUsernameKey = "CASSANDRA_USERNAME"
	EnvCassandraPasswordKey = "CASSANDRA_PASSWORD"
	EnvCassandraKeyspace    = "CASSANDRA_KEYSPACE"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err.Error())
	}
	cassandraHost := os.Getenv(EnvCassandraHostKey)
	cassandraUsername := os.Getenv(EnvCassandraUsernameKey)
	cassandraPassword := os.Getenv(EnvCassandraPasswordKey)
	cassandraKeyspace := os.Getenv(EnvCassandraKeyspace)
	if cassandraHost == "" || cassandraUsername == "" || cassandraPassword == "" || cassandraKeyspace == "" {
		panic("Cassandra host, username, password and keyspace must be set")
	}
	fmt.Printf("Connecting to Cassandra at %q\n", cassandraHost)
	cluster := gocql.NewCluster(cassandraHost)
	cluster.Keyspace = cassandraKeyspace
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: cassandraUsername, Password: cassandraPassword}
	session, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		panic(err.Error())
	}
	defer session.Close()
	ksMeta, err := session.KeyspaceMetadata(cassandraKeyspace)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Keyspace %q found as %q in Cassandra.\n", cassandraKeyspace, ksMeta.Name)

	// declare reused variables
	var message models.MessagesStruct
	var messages []models.MessagesStruct
	var roomId int64
	var speakerId int64
	var query *gocqlx.Queryx

	// fetch a message from the table
	query = session.Query("SELECT * FROM messages LIMIT 1", []string{})
	fmt.Printf("Executing query: %q\n", query.Statement())
	if err := query.GetRelease(&message); err != nil {
		fmt.Println("Failed to fetch message")
		panic(err.Error())
	}
	fmt.Printf("Message: \n%+v\n", message)

	// create secondary index on speaker_id to be able to use it in WHERE
	query = session.Query("CREATE INDEX IF NOT EXISTS ON messages(speaker_id)", []string{})
	fmt.Printf("Executing query: %q\n", query.Statement())
	if err := query.ExecRelease(); err != nil {
		fmt.Println("Failed to create index on speaker_id")
		panic(err.Error())
	}

	// fetch 5 most recent messages in room 1 of speaker 2
	roomId = 1
	speakerId = 2
	query = session.Query(models.Messages.SelectBuilder(models.Messages.Metadata().Columns...).
		Where(qb.Eq("speaker_id")).
		Limit(5).
		ToCql(),
	)
	fmt.Printf("Executing query: %q\n with bind values: roomId: %d, speakerId: %d\n", query.Statement(), roomId, speakerId)
	if err := query.Bind(roomId, speakerId).SelectRelease(&messages); err != nil {
		fmt.Println("Failed to fetch messages in room 1 of speaker 2")
		panic(err.Error())
	}
	fmt.Printf("Messages in room 1 of speaker 2:\n%+v\n", PrettyStringSlice(messages))

	// count number of messages in room 1 of speaker 2
	roomId = 1
	speakerId = 2
	query = session.Query(models.Messages.SelectBuilder().
		Where(qb.Eq("speaker_id")).
		CountAll().
		ToCql(),
	)
	var messageCount struct {
		Count int64
	}
	fmt.Printf("Executing query: %q\n with bind values: roomId: %d, speakerId: %d\n", query.Statement(), roomId, speakerId)
	if err := query.Bind(roomId, speakerId).GetRelease(&messageCount); err != nil {
		fmt.Println("Failed to count messages in room 1 of speaker 2")
		panic(err.Error())
	}
	fmt.Printf("Number of messages in room 1 of speaker 2:\n %+v\n\n", messageCount.Count)

	// count number of messages in each room
	// we could use a materialized view but here we use a COUNT with a GROUP BY
	query = session.Query("SELECT room_id, COUNT(*) FROM messages GROUP BY room_id", []string{})
	type messageCountByRoom struct {
		RoomId int64
		Count  int64
	}
	var messageCountByRoomSlice []messageCountByRoom
	fmt.Printf("Executing query: %q\n", query.Statement())
	if err := query.SelectRelease(&messageCountByRoomSlice); err != nil {
		fmt.Println("Failed to count messages in each room")
		panic(err.Error())
	}
	fmt.Printf("Number of messages in each room:\n%+v\n", PrettyStringSlice(messageCountByRoomSlice))

	// find all room ids
	query = session.Query("SELECT DISTINCT room_id FROM messages", []string{})
	fmt.Printf("Executing query: %q\n", query.Statement())
	var roomIds []int64
	if err := query.SelectRelease(&roomIds); err != nil {
		fmt.Println("Failed to find all room ids")
		panic(err.Error())
	}
	fmt.Printf("All room ids:\n%+v\n", PrettyStringSlice(roomIds))

	// BONUS

	// create a materialized view containing only room_id, time and message
	query = session.Query("CREATE MATERIALIZED VIEW IF NOT EXISTS messages_by_room AS SELECT room_id,time,message FROM messages WHERE room_id IS NOT NULL AND time IS NOT NULL PRIMARY KEY (room_id,time)", []string{})
	fmt.Printf("Executing query: %q\n", query.Statement())
	if err := query.ExecRelease(); err != nil {
		fmt.Println("Failed to create a materialized view")
		panic(err.Error())
	}
	// and fetch a message from the materialized view
	var messageMaterializedView models.MessagesStruct
	if err := session.Query("SELECT * FROM messages_by_room LIMIT 1", []string{}).GetRelease(&messageMaterializedView); err != nil {
		fmt.Println("Failed to fetch message from materialized view")
		panic(err.Error())
	}
	// since we are reusing the message struct, the speaker_id is present but set to 0 since it is not part of the materialized view
	fmt.Printf("Message from materialized view: \n%+v\n\n", messageMaterializedView)

	// create a UDF that returns whether a message contains profanity (hate for example)
	profaneWord := "hate"
	query = session.Query(fmt.Sprintf("CREATE OR REPLACE FUNCTION isProfane(message text) RETURNS NULL ON NULL INPUT RETURNS boolean LANGUAGE java AS $$ return message.toLowerCase().contains(%q); $$", profaneWord), []string{})
	fmt.Printf("Executing query: %q\n", query.Statement())
	if err := query.ExecRelease(); err != nil {
		fmt.Println("Failed to create a UDF")
		panic(err.Error())
	}
	// and use it - manually found a message that matches
	query = session.Query("SELECT room_id, time, speaker_id, message, isProfane(message) AS is_profane FROM messages WHERE room_id = 1 AND time = '2021-04-27 09:19:41.456000+0000' LIMIT 1", []string{})
	fmt.Printf("Executing query: %q\n", query.Statement())
	var messageIsProfane struct {
		Time      time.Time
		Message   string
		SpeakerId int64
		RoomId    int64
		IsProfane bool
	}
	if err := query.GetRelease(&messageIsProfane); err != nil {
		fmt.Println("Failed to fetch message with a use of the UDF")
		panic(err.Error())
	}
	fmt.Printf("Profane message: \n%+v\n\n", messageIsProfane)

	// find the oldest and latest message
	var minAndMaxTime struct {
		MinTime time.Time
		MaxTime time.Time
	}
	query = session.Query("SELECT min(time) AS min_time, max(time) AS max_time FROM messages", []string{})
	fmt.Printf("Executing query: %q\n", query.Statement())
	if err := query.GetRelease(&minAndMaxTime); err != nil {
		fmt.Println("Failed to find the oldest and latest message")
		panic(err.Error())
	}
	fmt.Printf("Oldest and latest message times: \n%+v\n\n", minAndMaxTime)

	// find the shortest and longest message
	// we need a UDF for getting a text length
	query = session.Query("CREATE OR REPLACE FUNCTION length(message text) RETURNS NULL ON NULL INPUT RETURNS int LANGUAGE java AS $$ return message.length(); $$", []string{})
	fmt.Printf("Executing query: %q\n", query.Statement())
	if err := query.ExecRelease(); err != nil {
		fmt.Println("Failed to create a UDF")
		panic(err.Error())
	}
	// now use it
	var minAndMaxLength struct {
		MinLength int
		MaxLength int
	}
	query = session.Query("SELECT min(length(message)) AS min_length, max(length(message)) AS max_length FROM messages", []string{})
	fmt.Printf("Executing query: %q\n", query.Statement())
	if err := query.GetRelease(&minAndMaxLength); err != nil {
		fmt.Println("Failed to find the shortest and longest message")
		panic(err.Error())
	}
	fmt.Printf("Shortest and longest message lengths: \n%+v\n\n", minAndMaxLength)

	// find an average message length for each speaker
	// first since speaker_id is not in the primary key we need to create a materialized view where it is
	query = session.Query("CREATE MATERIALIZED VIEW IF NOT EXISTS message_speaker AS SELECT room_id, time, speaker_id, message FROM messages WHERE speaker_id IS NOT NULL AND message IS NOT NULL AND room_id IS NOT NULL AND time IS NOT NULL PRIMARY KEY (speaker_id, room_id, time)", []string{})
	fmt.Printf("Executing query: %q\n", query.Statement())
	if err := query.ExecRelease(); err != nil {
		fmt.Println("Failed to create a materialized view")
		panic(err.Error())
	}

	// then we can group by speaker_id
	query = session.Query("SELECT speaker_id, avg(length(message)) AS avg_length FROM message_speaker GROUP BY speaker_id", []string{})
	fmt.Printf("Executing query: %q\n", query.Statement())
	var avgLengthBySpeaker []struct {
		SpeakerId int64
		AvgLength int64
	}
	if err := query.SelectRelease(&avgLengthBySpeaker); err != nil {
		fmt.Println("Failed to find an average message length for each speaker")
		panic(err.Error())
	}
	fmt.Printf("Average message length for each speaker: \n%+v\n", PrettyStringSlice(avgLengthBySpeaker))

}

func PrettyStringSlice[T any](slice []T) string {
	if len(slice) == 0 {
		return "- empty\n"
	}
	var sb strings.Builder
	for _, v := range slice {
		sb.WriteString(fmt.Sprintf("- %+v\n", v))
	}
	return sb.String()
}
