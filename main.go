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

//課題: 全テーブルをschema.SystemChatMessageに投げるのは、無理そう？RateLimit?の関係
//前提として、全テーブル、全カラムについて定数の範囲やテーブル間のデータ構造を考慮するのは厳しいので
//エラーがでているカラムに関してはよしなに使用者の方で修正してもらう

//定数問題：
//存在しない定数が入っていたらNGなので、それらの値は1-3の間で出力する

//外部キー
//外部キーでの制約があるものに関しては、NULLの値を入れて出力する

//他カラム依存：
//全部考慮できないので、ロジック組むのは作業として重そうなので考慮しない方針でいく

func connectLocalMysqlDB() (*sqlx.DB, error) {
	dsn := "root:@tcp(localhost:3306)/casone_lite_local?parseTime=true"
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
		var a, b string
		if err := rows.Scan(&a, &b); err != nil {
			panic(err)
		}
		return b, nil
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
	// フラグを用意する
	table := flag.String("table", "staffs", "string flag")
	num := flag.Int("n", 10, "int flag")
	fileOutput := flag.String("file", "sample", "string flag")
	flag.Parse()

	llm, err := openai.NewChat(openai.WithModel("gpt-3.5-turbo-16k"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("db-debug")
	db, err := connectLocalMysqlDB()
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
		`Please generate " + %d + " fixtures for the " + %s + "table in .yaml format with id and filled fields without the model name.
		and branch_id, tenant_id is between 1 and 3,
		and time columns is in timedate type,
		and value is based on japanese,
		and if there is a foreign key constraint, set the value to null,
		and id columns are between 1 value.
		`, *num, *table)

	fmt.Println("call-chatgpt-debug")
	completion, err := llm.Call(ctx, []schema.ChatMessage{
		schema.SystemChatMessage{Content: schemaInfo},
		schema.HumanChatMessage{Content: content},
	})
	if err != nil {
		log.Fatal(err)
	}

	//ファイル出力
	if *fileOutput != "" {
		err := writeOnYamlFile(*fileOutput+".yaml", completion.GetContent())
		if err != nil {
			fmt.Println(err)
		}
	}
}
