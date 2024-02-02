package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocql/gocql"
	"log"
	"time"
)

type LocationData struct {
	FImei          int64
	FDate          int64
	FTime          int64
	FAddr          string
	FDirection     int64
	FDistance      int64
	FEndTime       int64
	FLat           int64
	FLon           int64
	FPType         int
	FSpeed         int64
	FStartTime     int64
	FTotalDistance int64
	FType          int
	FWgs           string
}

type LocationDataShow struct {
	FImei          int64
	FDate          int64
	FTime          time.Time
	FAddr          string
	FDirection     int64
	FDistance      int64
	FEndTime       time.Time
	FLat           int64
	FLon           int64
	FPType         int
	FSpeed         int64
	FStartTime     time.Time
	FTotalDistance int64
	FType          int
	FWgs           string
}

var g_imeis []uint64

func migrate() {
	sourceCluster := gocql.NewCluster("47.110.46.6") // Replace with source Cassandra node IP
	sourceCluster.Keyspace = "slxk"
	sourceCluster.Consistency = gocql.LocalOne

	sourceSession, err := sourceCluster.CreateSession()
	if err != nil {
		panic(err)
	}
	defer sourceSession.Close()

	targetCluster := gocql.NewCluster("172.19.39.156") // Replace with target Cassandra node IP
	targetCluster.Keyspace = "slxk"
	targetCluster.Consistency = gocql.LocalOne

	targetSession, err := targetCluster.CreateSession()
	if err != nil {
		panic(err)
	}
	defer targetSession.Close()

	for _, imei := range g_imeis {
		fmt.Printf("migrate %v\n", imei)
		query := fmt.Sprintf("SELECT * FROM tkv_location WHERE fimei = %d and fdate>=20231212", imei)

		iter := sourceSession.Query(query).Iter()

		var result LocationData
		for iter.Scan(
			&result.FImei, &result.FDate, &result.FTime, &result.FAddr,
			&result.FDirection, &result.FDistance, &result.FEndTime, &result.FLat,
			&result.FLon, &result.FPType, &result.FSpeed, &result.FStartTime,
			&result.FTotalDistance, &result.FType, &result.FWgs) {
			// Insert into the target Cassandra cluster
			targetQuery := fmt.Sprintf("INSERT INTO tkv_location (fimei, fdate, ftime, faddr, fdirection, fdistance, fendtime, flat, flon, fptype, fspeed, fstarttime, ftotaldistance, ftype, fwgs) VALUES (%d, %d, %d, '%s', %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, '%s');",
				result.FImei, result.FDate, result.FTime, result.FAddr, result.FDirection,
				result.FDistance, result.FEndTime, result.FLat, result.FLon, result.FPType,
				result.FSpeed, result.FStartTime, result.FTotalDistance, result.FType, result.FWgs)

			if err := targetSession.Query(targetQuery).Exec(); err != nil {
				fmt.Println("Error inserting into target cluster:", err)
			}
		}

		if err := iter.Close(); err != nil {
			fmt.Println("Error closing iterator:", err)
		}
	}
}

func get_imeis() {
	// Database connection parameters
	dbUser := "admin"
	dbPassword := "shht"
	dbHost := "114.215.190.173"
	dbPort := "8000"
	dbName := "slxk"

	// Create a MySQL database connection
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Ping the database to check if the connection is successful
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// Query the database
	rows, err := db.Query("SELECT distinct(fimei) from t_garagebinddevice_0082")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate over the result set
	for rows.Next() {
		var imei uint64
		if err := rows.Scan(&imei); err != nil {
			log.Fatal(err)
		}

		// Process the retrieved data
		g_imeis = append(g_imeis, imei)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func query() {
	cluster := gocql.NewCluster("47.110.46.6") // Replace with source Cassandra node IP
	cluster.Keyspace = "slxk"
	cluster.Consistency = gocql.LocalOne

	session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	imei := 58231017557
	date := 20240118

	query := fmt.Sprintf("SELECT * FROM tkv_location WHERE fimei = %d and fdate = %v", imei, date)

	iter := session.Query(query).Iter()

	var result LocationData
	for iter.Scan(
		&result.FImei, &result.FDate, &result.FTime, &result.FAddr,
		&result.FDirection, &result.FDistance, &result.FEndTime, &result.FLat,
		&result.FLon, &result.FPType, &result.FSpeed, &result.FStartTime,
		&result.FTotalDistance, &result.FType, &result.FWgs) {

		show := LocationDataShow{
			FImei:          result.FImei,
			FDate:          result.FDate,
			FTime:          time.Unix(result.FTime, 0),
			FAddr:          result.FAddr,
			FDirection:     result.FDirection,
			FDistance:      result.FDistance,
			FEndTime:       time.Unix(result.FEndTime, 0),
			FLat:           result.FLat,
			FLon:           result.FLon,
			FPType:         result.FPType,
			FSpeed:         result.FSpeed,
			FStartTime:     time.Unix(result.FStartTime, 0),
			FTotalDistance: result.FTotalDistance,
			FType:          result.FType,
			FWgs:           result.FWgs,
		}
		//fmt.Println(result)
		fmt.Println(show.FImei, show.FTime, show.FPType, show.FType, show.FSpeed, show.FWgs)
	}
}

func main() {
	query()
	//get_imeis()
	//migrate()
	//test()
}

func test() {
	t := time.Unix(1705559009, 0)
	fmt.Println(t)
}
