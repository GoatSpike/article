## 初めに
現在、クリーンアーキテクチャ設計を元にして開発している案件においてページネーション関数を実装する機会があったので自分用のメモとして記事を書きました。
拙い文章もあるかとは存じますが目を通してくれたら嬉しいです。

## クリーンアーキテクチャについて
![18ffbea5828e-20230915.png](https://qiita-image-store.s3.ap-northeast-1.amazonaws.com/0/3672923/8ada20fa-6733-4ee1-bd52-b37f2821a3e1.png)

[引用元：The Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

![25116acbdf7e-20230915.png](https://qiita-image-store.s3.ap-northeast-1.amazonaws.com/0/3672923/23833343-1d35-4d03-a6bd-e7d8c4f4b442.png)

[引用元: Clean Architecture　達人に学ぶソフトウェアの構造と設計](https://www.amazon.co.jp/dp/4048930656)

こちら二つの図はネットなどでクリーンアーキテクチャについて説明をされる時よく使われる図です。

| クラス              | 説明                                    |
|--------------------------|-----------------------------------------|
| Controller/Presenter           | APIのコントローラー                     |
| View                     | APIのレスポンス                         |
| Input Boundary/ Output Boundary           | ユースケースのインタフェース           |
| Use Case Interactor      | ユースケースの実装クラス               |
| Input Data               | ユースケースの関数の引数               |
| Output Data              | ユースケースの関数の戻り値             |
| Entities                 | ビジネスロジックの実装クラス           |
| Data Access Interface    | リポジトリのインタフェース             |
| Data Access              | リポジトリの実装クラス                 |

## 案件で実際に使用されているクリーンアーキテクチャ設計

案件ではすでにテンプレートが存在し、下記のような図のようにクラスを絞り実装しております。

![cleanArc.jpg](https://qiita-image-store.s3.ap-northeast-1.amazonaws.com/0/3672923/b8db269e-acae-42fa-8eab-50e713e58b65.jpeg)

| クラス                      | 説明                                                                         |
|----------------------------|------------------------------------------------------------------------------|
| Controller                 | APIのコントローラー（インターフェースアダプター層）                          |
| Use Case Interactor        | ユースケースの実装クラス（ユースケース層）                                   |
| Input                      | ユースケースの関数の引数（ユースケース層）                                   |
| Output                     | ユースケースの関数の戻り値のインターフェース（ユースケース層）               |
| Presenter                  | ユースケースの関数の戻り値の実装クラス（インターフェースアダプター層）      |
| Entities                   | ビジネスロジックの実装クラス（エンティティ層）                               |
| Repository Interface       | リポジトリのインターフェース（ユースケース層。本来であればインターフェースアダプター層）                |
| Repository                 | リポジトリの実装クラス（インターフェースアダプター層。フレームワークおよびドライバー層）                   |



```
project_root/
├── handler/
│   └── api/
│       └── main.go                     // エントリーポイント
│       └── api.go                     // APIのセットアップ(Controllerをまとめている)
├── pkg/
│   ├── interfaces/
│   │   ├── gateway/
│   │   │    ├── api/
│   │   │    │   └── example_controller.go    // Controller
│   │   │    └── repositories/
│   │   │         └── example_repository.go    // Repository (Implemention)
│   │   └── presenters/
│   │        └── example_presenter.go     // Presenter
│   ├── domain/
│   │   └── model/            
│   │        └── example_entity.go            // Entities
│   └── usecases/
│       ├── interactor/
│       │   └── example_interactor.go   // Use Case Interactor
│       ├── input/
│       │   └── example_input.go        // Input
│       ├── output/
│       │   └── example_output.go       // Output
│       └── repositories/
│            └── example_repository.go  // Repository Interface
└── configs/
    └── config.go                        // 設定関連
```

