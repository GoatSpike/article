package main

import (
	"fmt"
	"net/http"
	"time"
	// "golang.org/x/sync/errgroup"
)

const rateLimit = 2 // 1秒間に2リクエストまで許可

// URLからデータを取得する例の関数
func fetchURL(client *http.Client, url string) error {
	// URLをログ出力
	fmt.Println("Fetching URL:", url)
	// HTTP GETリクエストを送信
	resp, err := client.Get(url)
	if err != nil {
		// エラーが発生した場合、そのエラーを返す
		return err
	}
	// 関数終了時にレスポンスボディを閉じる
	defer resp.Body.Close()
	return nil
}

// `setLimit` 関数は特定の時間内にリクエストの数を制限する
func setLimit(limit int) chan struct{} {
	ch := make(chan struct{}, limit)
	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			for i := 0; i < limit; i++ {
				ch <- struct{}{}
			}
			<-ticker.C
		}
	}()
	return ch
}

// func main() {
// 	// errgroup.Groupの作成
// 	var g errgroup.Group
// 	// HTTPクライアントを作成
// 	client := &http.Client{}
// 	// チェックしたいURLのリスト（2つは存在しないURLに設定）
// 	urls := []string{
// 		"https://www.google.com",
// 		"https://www.invalid-url.com",         // 無効なURL（エラーが発生）
// 		"https://www.another-invalid-url.com", // 無効なURL（エラーが発生）
// 	}
// 	// レートリミット制限用のチャンネルを取得
// 	limitCh := setLimit(rateLimit)
// 	// 各URLに対してゴルーチンを起動
// 	for _, url := range urls {
// 		// urlの値をキャプチャするためにローカル変数を使用
// 		url := url
// 		// URLごとにゴルーチンを起動
// 		g.Go(func() error {
// 			<-limitCh // レートリミット制限を待つ
// 			return fetchURL(client, url)
// 		})
// 	}
// 	// 全てのゴルーチンが完了するのを待ち、エラーをチェック
// 	if err := g.Wait(); err != nil {
// 		// 一つ以上のゴルーチンでエラーが発生した場合、そのエラーを出力
// 		fmt.Println("Encountered error:", err)
// 	} else {
// 		// 全てのURLが正常にフェッチされた場合
// 		fmt.Println("All URLs fetched successfully")
// 	}
// }
