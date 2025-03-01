# Go言語のゴルーチン（Goroutine）について

## 初めに

今回Go言語での実装においてゴルーチンとerrgroupを扱う場面がありましたので、勉強用として自分用にまとめてみました。

## (goroutine)ゴルーチンについて

Go言語は、`並行処理`（並行プログラミングとも言われる）を簡単に実装できるように設計されています。Goでは、並行処理を行うための軽量なスレッド「ゴルーチン（goroutine）」というものがあります。

通常プログラミングにおいてタスク1が完了したら、タスク2に取り掛かる逐次処理が基本となっています。並列処理は逐次処理とは違い、タスク1とタスク2を交互に進めることができます。この場合、並列処理はシングルコア・マルチスレッドで処理を行います。
- 並行処理：複数の処理を1つの主体が切り替えながらこなすこと。実行時間が早くなる（かもしれない）
  - 引用:[並行処理、並列処理、同期処理、非同期処理についてまとめ](https://qiita.com/kyabetsuda/items/384a57ff6b7250de40ad)
- シングルコア：1つのCPUに対して1つのコアが内蔵されているCPU。コアとはCPUにある機能の事で実際に処理を行う部品になります。（
  - 引用:[CPUの「コア数」「クロック周波数」「スレッド数」とは？](https://dosparaplus.com/library/details/000646.html)
- マルチスレッド：1つのコンピュータープログラムを実行する際に、アプリケーションのプロセス（タスク）を複数のスレッドに分けて並行処理する流れの事
  - 引用:[マルチスレッドとは](https://wa3.i-3-i.info/word12455.html)


調べてて勉強になった事して、並行処理と**並列処理**はぜんぜん違うという事です。
今回は並行処理には触れず、Go言語が提供している**並行処理 goroutine**について焦点を当てていきますのでご了承ください。

一応この二つの違いについてまとめられている[並行処理と並列処理](https://zenn.dev/hsaki/books/golang-concurrency/viewer/term)という記事も置いておきます。

## 特徴

1. **軽量性**:
   ゴルーチンは非常に軽量なスレッドです。  
   一般的なスレッドの場合、CPUコアに対してマッピングされOSによって管理されます。スレッドにの切り替えに伴うコンテキストスイッチによってオーバヘットが発生し得るとの事です。一方ゴルーチンはカーネルスレッドに対してマッピングされGoランタイムによって管理されます。ゴルーチンの切り替えはスレッド内部の処理に留まるためオーバヘットが極めて小さくなります。これが軽量スレッドと言われる理由です。
   - 引用：[Goroutine はなぜ軽量スレッドと称されるのか](https://www.ren510.dev/blog/goroutine-concurrency)という記事においてより詳しく説明されていました。

2. **メモリ消費**:
   ゴルーチンは新規作成時点では一般に2KBのスタック領域となります。また、自動的にスタックサイズが増減（スタックが動的に再割り当て）され、メモリが足りない場合はヒープを使用します。
   - 引用：[【Go言語入門】goroutineとは？ 実際に手を動かしながら goroutineの基礎を理解しよう](https://www.ariseanalytics.com/activities/report/20221005/)


## 基本的な使用方法

### 1 **ゴルーチンの開始**
   ゴルーチンを開始するには `go` キーワードを使います。

関数への利用
```go
go function()
```

無名関数への利用
```go
go func() {
...
}()
```

goroutine.goプログラムを実行してみます。  
メイン関数で2つのゴルーチンを起動していますが、メイン関数が終了すると全てのゴルーチンも終了するため、実際には何も表示されない可能性があります。そこで`time.Sleep(time.Second)`を追加しゴルーチンが作業を完了するのを待つ必要があります。

goroutine.go
```go
package main

import (
	"fmt"
	"time"
)

func say(s string) {
	fmt.Println(s)
}

func main() {
    go say("hello")
	go say("world")
	time.Sleep(time.Second)
}
```

```
$ go run goroutine.go
hello
world
$ go run goroutine.go
world
hello
```

何回か実行するとわかるのですが、goroutineはそれぞれ独立しているため実行の順序性が担保されていません。実行の順序性を制御する為`channel`を利用します

### 2 **チャネルを用いた通信**  
ゴルーチン間の通信には「チャネル（Channel）」を使用します。チャネルを使用することでゴルーチン同士の値をやり取りすることができたり、実行の順序を制御することができます

```go
ch := make(chan string)

ch <- data    // channelにdata変数の値を送る
data := <- ch // channelから値を取り出し、その値をdata変数に入れる
```

では実際にchannelを使ってデータの送受信、および実行順序を制御してみます。

```go
package main

import (
	"fmt"
	"time"
)

func Say(c chan string) {
	data := <-c // 2
	fmt.Println(data)

	data = "world"
	c <- data // 3
}

func main() {

	ch := make(chan string)

	go func(c chan string) {
		data := "hello"
		c <- data // 1

		data = <-c // 4
		fmt.Println(data)
	}(ch)

	go Say(ch)

	time.Sleep(time.Second)
}

```

実行の流れ
1. メイン関数でチャネルを作成される
2. 無名関数が起動され、`"hello"`をチャネルに送信する（1）
3. Sayゴルーチンが起動され、チャネルから`"hello"`を受信する（2）
4. Sayゴルーチンは、受信した`"hello"をコンソールに表示する。それから、チャネルに"world"`を送信する（3）
5. 無名関数がチャネルから`"world"`を受信し（4）、それをコンソールに表示する

全体としての出力  
hello（`Say`関数によって表示）  
world（無名関数によって表示）

**channelの方向**  
基本的にchannelは送信と受信を行うことができますが、どちらか一方の機能を持つchannelを作ることも可能です
```go
// 受信用channel
c1 := make(<-chan Type)

// 送信用channel
c2 := make(chan<- Type)
```

**close()**  
close()を使用してchannelを閉めることができます。channelを閉めたら、該当channelには二度と送信することはできません。しかし、channelに値が存在する限り受信は可能です
```go
close(myChannel)
```

下のコードを使用してchannelが閉じているかどうか確認することができます。閉じていたらcheckがfalseになり、開いていたらcheckがtrueになります
```go
data, check := <-myChannel
```

**for range**  
for rangeを使用してchannelが閉じる時まで値を受信することができます。channelが開いていたらrangeはchannelに値が入るまで待機します。channelが閉じられたらループは終了になります

```go
for data := range myChannel {
  ...
}
```


**select**  
switchと似ていますが、selectでcaseはchannelで送信または受信作業を意味します。selectはcaseのいずれかが実行されるまで待機します。selectにdefaultがあれば、caseが用意されていなくても待機せずにdefaultを実行します
```go
select {
  case <-ch1:
    // ch1に値が入った時に実行
  case <-ch2:
    // ch2に値が入った時に実行
  default:
    // 全てのchannelに値が入らなかった時に実行
}
```

## 同期化オブジェクト
同期化オブジェクトとは複数のゴルーチン間でリソースやデータの一貫性を確保し、同期を取るための機能です。また実行順序についても制御することができます

### Mutex
```go
package main

import (
	"fmt"
	"time"
)

func main() {

	var data = []int{}

	go func() {
		for i := 0; i < 10000; i++ {
			data = append(data, 1)
		}
	}()

	go func() {
		for i := 0; i < 10000; i++ {
			data = append(data, 1)
		}
	}()

	time.Sleep(2 * time.Second)

	fmt.Println(len(data)) // スライスの長さを出力
}
```

上記のコードを実行すると
```
go run goroutine.go
10002
go run goroutine.go
10000
 go run goroutine.go
10016
go run goroutine.go
9542
```

結果としては20000を期待する所、結果として様々な値が出力されています。これは二つのゴルーチンが競合し、同時に値にアクセスした為append()が正確に処理されなかったことが原因です。
MutexのLock()、Unlock()を使用する事で排他制御で実行することができます

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	var data = []int{}
	var mutex = new(sync.Mutex)

	go func() {
		for i := 0; i < 10000; i++ {
			mutex.Lock() // スライスを保護
			data = append(data, 1)
			mutex.Unlock() // スライスを保護解除
		}
	}()

	go func() {
		for i := 0; i < 10000; i++ {
			mutex.Lock() // スライスを保護
			data = append(data, 1)
			mutex.Unlock() // スライスを保護解除
		}
	}()

	time.Sleep(2 * time.Second) // 2秒待機

	fmt.Println(len(data)) // スライスの長さを出力
}
```

注意点としてLock()とUnlock()はペアを合わせる必要があります。ペアが合わない場合デットロックが発生します

### WaitGroup
​sync.WaitGroupの利用​:
- sync.WaitGroupは、複数のゴルーチンの完了を待つためのカウンタを提供します。
- wg.Add(1)でカウンタを増やし、新しいゴルーチンが開始されることを示します。
- 各ゴルーチンが`wg.Done()`を呼び出すと、カウンタが減少します。

```go
package main

import (
	"fmt"
	"sync"
)

func say(s string, wg *sync.WaitGroup) {
	defer wg.Done() // ゴルーチンが完了したことを通知
	fmt.Println(s)
}

func main() {
	var wg sync.WaitGroup

	wg.Add(1) // "hello" ゴルーチンを追加
	go say("hello", &wg)

	wg.Add(1) // "world" ゴルーチンを追加
	go say("world", &wg)

	wg.Wait() // すべてのゴルーチンが完了するのを待機
}
```


### errgroup  
今回紹介したsync.WaitGroupはエラーハンドリング機能がないため、それぞれのゴルーチン内で発生したエラーを個別に管理する必要があります

errgroupは`golang.org/x/sync/errgroup`パッケージに含まれている同期化オブジェクトです  
これは、複数のゴルーチンをグループとして管理し、それらのゴルーチンが完了するのを待つとともに、ゴルーチンの中で発生した最初のエラーを一括して処理するためのものです。`sync.WaitGroup`と同様の機能を拡張し、エラーハンドリング機能を提供することが可能です。  
今回案件で使用したのはerrgroupになります

インストール方法
```
$ go get golang.org/x/sync/errgroup
```

```go
package main

import (
	"fmt"
	"net/http"

	"golang.org/x/sync/errgroup"
)

// URLからデータを取得する例の関数
func fetchURL(url string) error {
	fmt.Println("Start fetching URL:", url)
	// HTTP GETリクエストを送信
	resp, err := http.Get(url)
	if err != nil {
		// エラーが発生した場合、そのエラーを返す
		return err
	}

	// 関数終了時にレスポンスボディを閉じる
	defer resp.Body.Close()
	return nil
}

func main() {
	// errgroup.Groupの作成
	var g errgroup.Group
	// チェックしたいURLのリスト（2つは存在しないURLに設定）
	urls := []string{
		"https://www.google.com",
		"https://www.invalid-url.com", // 無効なURL（エラーが発生）
		"https://www.github.com",
		"https://www.another-invalid-url.com", // 無効なURL（エラーが発生）
	}
	// 各URLに対してゴルーチンを起動
	for _, url := range urls {
		// urlの値をキャプチャするためにローカル変数を使用
		url := url
		// URLごとにゴルーチンを起動
		g.Go(func() error {
			return fetchURL(url)
		})
	}
	// 全てのゴルーチンが完了するのを待ち、エラーをチェック
	if err := g.Wait(); err != nil {
		// 一つ以上のゴルーチンでエラーが発生した場合、そのエラーを出力
		fmt.Println("Encountered error:", err)
	} else {
		// 全てのURLが正常にフェッチされた場合
		fmt.Println("All URLs fetched successfully")
	}
}

```

出力結果
```
Start fetching URL: https://www.another-invalid-url.com
Start fetching URL: https://www.google.com
Start fetching URL: https://www.github.com
Start fetching URL: https://www.invalid-url.com
Encountered error: Get "https://www.another-invalid-url.com": dial tcp: lookup www.another-invalid-url.com on 169.254.169.254:53: dial udp 169.254.169.254:53: connect: no route to host
```

[playGround](https://goplay.tools/snippet/nv5tCbjE6Pz)

使い方としては
- 通常のgoroutineの構文は go f() ですが、errgroup.Groupのインスタンスを生成し、g.GO(f()) でgoroutineを起動します。引数に関数を渡すのが特徴です
- fetchURLでエラーが発生した場合、エラーが返却されます。エラーがない場合`err==nil`です
- g.Wait()で全てのゴルーチンが終了するまで待ちます。エラーが複数あった場合は最初のエラーが返却することになります


またsetLimit()を使えばと特定の時間内に操作を何回実行できるかを制限したりすることもできます。
```go
package main

import (
	"fmt"
	"time"
)

// `setLimit` 関数は特定の時間内に操作回数を制限する
func setLimit(limit int) chan struct{} {
	ch := make(chan struct{}, limit)
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for {
			for i := 0; i < limit; i++ {
				ch <- struct{}{}
			}
			<-ticker.C
		}
	}()
	return ch
}

func main() {
	// レートリミットを設定（ここでは5秒間に2回まで）
	limitCh := setLimit(2)

	for i := 0; i < 10; i++ {
		// レートリミットのトークンを取得
		<-limitCh
		// 制限付きでプリント文を実行
		fmt.Println("Processing item", i+1)
	}

	fmt.Println("All items processed.")
}
```


引用
- [並行処理、並列処理、同期処理、非同期処理についてまとめ](https://qiita.com/kyabetsuda/items/384a57ff6b7250de40ad)）
- [CPUの「コア数」「クロック周波数」「スレッド数」とは？](https://dosparaplus.com/library/details/000646.html)）
- [マルチスレッドとは](https://wa3.i-3-i.info/word12455.html)
- [並行処理と並列処理](https://zenn.dev/hsaki/books/golang-concurrency/viewer/term)
- [Goroutine はなぜ軽量スレッドと称されるのか](https://www.ren510.dev/blog/goroutine-concurrency)
- [【Go言語入門】goroutineとは？ 実際に手を動かしながら goroutineの基礎を理解しよう](https://www.ariseanalytics.com/activities/report/20221005/)
