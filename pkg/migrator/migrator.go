package migrator

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const StringFieldLenghtTag = "strlen"
const maxStringFieldValue = 100

var pool *pgxpool.Pool

func Init() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	defer conn.Close(context.Background())
	if err != nil {
		conn.Close(context.Background())
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
}

func Migrate(types []any) {
	var err error
	pool, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	query := strings.Builder{}
	err = pool.QueryRow(context.Background(), fmt.Sprintf("DROP DATABASE IF EXISTS %s", os.Getenv("DATABASE_NAME"))).Scan()
	_, err = pool.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s", os.Getenv("DATABASE_NAME")))
	if err != nil {
		log.Fatal(err)
	}
	_, err = pool.Exec(context.Background(), fmt.Sprintf("DROP SCHEMA %s cascade\n", os.Getenv("SCHEMA_NAME")))
	if err != nil {
		log.Fatal(err)
	}
	_, err = pool.Exec(context.Background(), fmt.Sprintf("CREATE SCHEMA %s\n", os.Getenv("SCHEMA_NAME")))
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range types {
		sv := reflect.ValueOf(v)
		st := reflect.TypeOf(v)
		typeName := sv.Type().Name()
		query.WriteString(fmt.Sprintf("CREATE TABLE %s.%s (\n", os.Getenv("SCHEMA_NAME"), typeName))
		for i := 0; i < sv.NumField(); i++ {
			f := st.Field(i)
			fv := sv.Field(i)
			switch fv.Type().Kind() {
			case reflect.String:
				if tag, ok := f.Tag.Lookup(StringFieldLenghtTag); ok && tag != "" {
					strlen, err := strconv.ParseInt(tag, 10, 32)
					if err != nil {
						log.Fatalf("Migrate: incorrect field tag. %v", err)
					}
					query.WriteString(fmt.Sprintf("%s varchar(%d)", f.Name, strlen))
				} else {
					query.WriteString(fmt.Sprintf("%s varchar(%d)", f.Name, maxStringFieldValue))
				}
			case reflect.Int:
				query.WriteString(fmt.Sprintf("%s int", f.Name))
			}
			if i != sv.NumField()-1 {
				query.WriteByte(',')
			}
			query.WriteByte('\n')
		}
		query.WriteString(");\n")
	}
	log.Printf("query: %s\n", query.String())
	_, err = pool.Exec(context.Background(), query.String())
	if err != nil {
		log.Fatal(err)
	}
}

func Seed(seeds []any) {
	query := strings.Builder{}

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	defer conn.Close(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range seeds {
		query.WriteString(fmt.Sprintf("INSERT INTO %s.%s ", os.Getenv("SCHEMA_NAME"), reflect.TypeOf(v).Name()))
		query.WriteString(fmt.Sprintf("VALUES (\n"))

		sv := reflect.ValueOf(v)
		for i := 0; i < sv.NumField(); i++ {
			field := sv.Field(i)
			query.WriteString(fmt.Sprintf("'%v'", reflect.Indirect(field)))
			if i != sv.NumField()-1 {
				query.WriteByte(',')
			}
		}
		query.WriteString(fmt.Sprintf(");\n"))
		log.Printf("%s", query.String())

		err = conn.QueryRow(context.Background(), query.String()).Scan()
		query.Reset()
	}
	log.Printf("%s", query.String())
}
