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
	//go:embed sqlc/schema/schema.sql
	ddl string

	ctx    context.Context
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
	// 最初の一回目の計測は正確でないので捨てる
	gosi.RefreshCpu()
	for {
		time.Sleep(1 * time.Second)
		gosi.RefreshCpu()

		timestamp := time.Now().Format("2006-01-02 15:04:05")
		cpuStat := gosi.Cpu()
		param := queries.InsertCpuOriginalDataParams{
			Timestamp: timestamp,
			CpuUsage:  cpuStat.Total,
		}
		fmt.Println(cpuStat, param)

		err := qry.InsertCpuOriginalData(ctx, param)
		if err != nil {
			log.Printf("error!: %v\n", err)
			return err
		}
	}

	return nil
}

func insertDownSamplingCpuData() error {
	for {
		time.Sleep(5 * time.Second)

		rows, err := getDownsamplingData()
		if err != nil {
			log.Printf("error!: %v\n", err)
			return err
		}
		for _, v := range rows {
			//fmt.Println(i, v)
			param := queries.InsertCpuDownsampledDataParams{
				Timestamp:   v.Dstimestamp,
				AvgCpuUsage: v.AveCpuUsage,
				MaxCpuUsage: v.MaxCpuUsage,
			}

			err := qry.InsertCpuDownsampledData(ctx, param)
			if err != nil {
				log.Printf("error!: %v\n", err)
				return err
			}
		}
	}

	return nil
}

// getDownsamplingData CPU情報を読み込んで指定した間隔（秒）でダウンサンプリングして返す
func getDownsamplingData() ([]queries.SelectCpuDownsamplingDataRow, error) {
	var results []queries.SelectCpuDownsamplingDataRow
	var err error
	results, err = qry.SelectCpuDownsamplingData(ctx, "60")
	if err != nil {
		log.Printf("error!: %v\n", err)
		return results, err
	}

	return results, nil
}

func main() {
	setLogger()

	ctx = context.Background()
	if err := initializeDB(); err != nil {
		log.Fatal(err)
	}
	log.Println("initdb success")
	defer dbconn.Close()

	go insertCpuData()
	go insertDownSamplingCpuData()

	initializeServer()

}

func setLogger() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
