package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
)

//TODO: completion用のテキストを作成する関数を用意する
//TODO: どういう感じでfixtureのアウトプットをしたいか(適当にどのようなyamlが作成されるかいくつかパターンを用意する、次回のゴールを設定するイメージ)

//1. casoneのデータベースを用意する必要がある→ダンプで最新のデータをとってくる必要がある
//2. casoneのデータベースダンプファイルからテーブルの定義を取得する
//3. 上で取得したテーブル名、生成個数をOpenAPIに投げる
//4. numで定義した引数個数分のfixtureを生成して出力する（全カラム分）

//detail: casoneのスタッフのテーブルを作成する
// fixtureのデータは日本語文字列を含むことができる
// カラムの値を指定できるようにする

//課題: 全テーブルをschema.SystemChatMessageに投げるのは、無理そう？RateLimit?の関係
//まだよくわかっていないが、欲しいテーブルだけなげるのが良さそう

const msg = `
defined table users

CREATE TABLE users (
  uuid uuid NOT NULL,
  name varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  email varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT ''
)
`

func readFile() (lines string, err error) {
	b, err := ioutil.ReadFile("casone_lite_local_2023-09-08.sql")
	if err != nil {
		fmt.Println(os.Stderr, err)
		os.Exit(1)
	}
	lines = string(b)
	return lines, err
}

func writeFile(lines *schema.AIChatMessage, filePath string) (err error) {
	f, err := os.Open(filePath)
	data := []byte(lines.Content)
	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func connectLocalMysqlDB() (*sqlx.DB, error) {
	dsn := "root:@tcp(localhost:3306)/otelsql?parseTime=true"
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to connect to the MySQL database: %v", err)
	}
	defer db.Close()
	return db, nil
}

func main() {
	// フラグを用意する
	table := flag.String("table", "staffs", "string flag")
	num := flag.Int("n", 10, "int flag")
	//fill := flag.String("fill", "", "string flag")
	//language := flag.String("lang", "", "string flag")

	//file作成のフラグを用意する
	flag.Parse()

	//ファイル情報を読み込む
	sqlLines, err := readFile()

	llm, err := openai.NewChat()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	content := "Please generate " + string(*num) + " fixtures for the " + *table + " table in .yaml format without the model name."
	completion, err := llm.Call(ctx, []schema.ChatMessage{
		schema.SystemChatMessage{Content: sqlLines},
		schema.HumanChatMessage{Content: content},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(completion)
	writeFile(completion, "output.yaml")
}
