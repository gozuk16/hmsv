package main

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gozuk16/gosi"
	"github.com/gozuk16/hmsv/queries"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

var (
	ctx context.Context
	//go:embed sqlc/schema/schema.sql
	ddl string
	//dbconn *sql.DB
	dbconn *sql.DB
	qry    *queries.Queries
)

func initializeServer() error {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", hello)

	e.Logger.Fatal(e.Start(":8888"))

	return nil
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "hello world!")
}

// initializeDB sqlite3の初期化とコネクション作成
func initializeDB() error {
	db, err := sql.Open("sqlite3", "hm.db")
	if err != nil {
		return err
	}
	dbconn = db

	if _, err := dbconn.ExecContext(ctx, ddl); err != nil {
		return err
	}

	qry = queries.New(dbconn)

	return nil
}

// insertCpuData CPU情報の登録
func insertCpuData() error {
	/*
		tx, err := dbconn.BeginTx(ctx, nil)
		//tx, err := dbconn.Begin()
		if err != nil {
			log.Println(err)
			return err
		}
		defer func() {
			tx.Rollback()
			if recover() != nil {
				// panicの場合は別にログを出す
				log.Printf("error!: %v\n", recover())
			}
		}()

		//qtx := queries.New(dbconn)
		qtx := queries.New(dbconn).WithTx(tx)
	*/
	for {
		cpuStat := gosi.Cpu()
		fmt.Println(cpuStat)
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		param := queries.InsertCpuOriginalDataParams{
			Timestamp: timestamp,
			CpuUsage:  float64(cpuStat.Total),
		}

		err := qry.InsertCpuOriginalData(ctx, param)
		if err != nil {
			log.Printf("error!: %v\n", err)
			return err
		}
		time.Sleep(5 * time.Second)
	}

	//	_, _ = queries.InsertDownsampledData(ctx)

	return nil
}

func insertDownSamplingCpuData() error {
	for {
		time.Sleep(5 * time.Second)

		getDownsamplingData()
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		param := queries.InsertCpuDownsampledDataParams{
			Timestamp:   timestamp,
			AvgCpuUsage: 0,
			MaxCpuUsage: 0,
		}

		err := qry.InsertCpuDownsampledData(ctx, param)
		if err != nil {
			log.Printf("error!: %v\n", err)
			return err
		}
	}

	return nil
}

// getDownsamplingData CPU情報を読み込んで指定した間隔（秒）でダウンサンプリングして返す
func getDownsamplingData() error {
	results, _ := qry.SelectCpuDownsamplingData(ctx, "60")
	for i, v := range results {
		fmt.Println(i, v)
	}

	return nil
}

func main() {
	setLogger()

	ctx = context.Background()
	if err := initializeDB(); err != nil {
		log.Fatal(err)
	}
	log.Println("initdb success")
	defer dbconn.Close()

	/*
		if err := insertCpuData(); err != nil {
			log.Fatal(err)
		}

		if err := getDownsamplingData(); err != nil {
			log.Fatal(err)
		}
	*/
	go insertCpuData()

	initializeServer()

}

func setLogger() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
