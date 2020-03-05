package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/MariusVanDerWijden/mvp/simple"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

func main() {
	conn := connectToDB()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	for {
		rows, err := query(ctx, "revenue", conn)
		if err != nil {
			log.Fatal("unable to query table: ", err)
		}

		_, flt, err := scanRows(rows)
		if err != nil {
			log.Fatal("unable to scan row: ", err)
		}
		if err := simple.AnalyzeNum(1, flt); err != nil {
			log.Fatal("could not analyze output: ", err)
		}
		time.Sleep(5 * time.Second)
	}
}

func connectToDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	pool, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("unable to use data source name", err)
	}

	pool.SetConnMaxLifetime(0)
	pool.SetMaxIdleConns(3)
	pool.SetMaxOpenConns(3)

	return pool
}

func query(ctx context.Context, tablename string, conn *sql.DB) (*sql.Rows, error) {
	query := fmt.Sprintf("select * from %s", tablename)
	return conn.QueryContext(ctx, query)
}

func scanRows(rows *sql.Rows) ([]string, []float64, error) {
	var outStr []string
	var outFlt []float64
	for rows.Next() {
		var str sql.NullString
		var flt sql.NullFloat64
		err := rows.Scan(&str, &flt)
		if err != nil {
			return nil, nil, err
		}
		outStr = append(outStr, str.String)
		outFlt = append(outFlt, flt.Float64)
	}
	return outStr, outFlt, nil
}

/*
func analyzeOutput(output [][]interface{}) error {
	for i, t := range output {
		switch reflect.TypeOf(t) {
		case reflect.TypeOf([]sql.NullFloat64{}):
			log.Print("tick")
			data := make([]float64, len(t))
			for i, elem := range t {
				ele := elem.(sql.NullFloat64)
				data[i] = ele.Float64
			}
			if err := simple.AnalyzeNum(i, data); err != nil {
				return err
			}
		default:
			log.Print("tock")
			log.Println(reflect.TypeOf(t))
		}
	}
	return nil
}
*/
