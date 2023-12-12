package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
)

func connectDB() (*sqlx.DB, error) {
	dsn := "root:@tcp(localhost:3306)/dammykun?parseTime=true"
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to connect to the MySQL database: %v", err)
	}
	return db, nil
}

func showSchemaInfo(db *sqlx.DB, table *string) (string, error) {
	rows, err := db.Query(fmt.Sprintf("show create table %s", *table))
	if err != nil {
		log.Fatalf("failed to show columns from the MySQL database: %v", err)
		return "", err
	}
	defer rows.Close()

	for rows.Next() {
		var tableName, schemaInfo string
		if err := rows.Scan(&tableName, &schemaInfo); err != nil {
			panic(err)
		}
		return schemaInfo, nil
	}
	return "", errors.New("not found")
}

func writeOnYamlFile(fileName string, data string) error {
	f, err := os.Create(fileName)
	_, err = f.Write([]byte(data))
	if err != nil {
		return err
	}
	return nil
}

func main() {
	num := flag.Int("n", 10, "int flag")
	table := flag.String("table", "users", "string flag")
	fileOutput := flag.String("file", "sample", "string flag")
	flag.Parse()

	llm, err := openai.NewChat(openai.WithModel("gpt-4-1106-preview"))
	if err != nil {
		log.Fatal(err)
	}

	db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	schemaInfo, err := showSchemaInfo(db, table)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	content := fmt.Sprintf(
		`Please generate " + %d + " fixtures for the " + %s + "table in .yaml format with id and filled fields without the model name,
		if there is a foreign key constraint, set the value to null,
		time columns is in timedate type,
		value is based on japanese,
		id columns are between 1 value,
		and without message.
		`, *num, *table)

	completion, err := llm.Call(ctx, []schema.ChatMessage{
		schema.SystemChatMessage{Content: schemaInfo},
		schema.HumanChatMessage{Content: content},
	})
	if err != nil {
		log.Fatal(err)
	}

	if *fileOutput != "" {
		err := writeOnYamlFile(*fileOutput+".yaml", completion.GetContent())
		if err != nil {
			fmt.Println(err)
		}
	}
}
