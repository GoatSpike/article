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

### 実装ファイル
apiのセットアップ
```go:api.go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/example/configs"
	"github.com/example/pkg/interfaces/api"
)

type API struct {
	Engine *gin.Engine
	Port   string
	Config *configs.APIConfig
}

func NewAPI(
	config *configs.APIConfig,
	genresAPI api.GenresController,
) *API {
	router := &API{
		Engine: gin.New(),
		Port:   config.Port,
		Config: config,
	}
	genres := router.Engine.Group("/genres")
	{
		genres.GET("/", func(c *gin.Context) { genresAPI.Search(c) })
		genres.POST("/", func(c *gin.Context) { genresAPI.Create(c) })
	}
	return router
}

func (r *API) Run() error {
	return r.Engine.Run(r.Port)
}
```

Controller層
```go:genres_controller.go
package api

import (
	"net/http"

	"github.com/example/pkg/usecase/input"
	"github.com/example/pkg/usecase/interactor"
)

type GenresController struct {
	ApplicationController
	genreSearchInteractor interactor.GenreSearchInteractor
	genreCreateInteractor interactor.GenreCreateInteractor

}

func (c GenresController) Search(ctx *Context) {
	context := c.CreateContext(ctx)

	genres, err := c.genreSearchInteractor.Invoke(context)
	if err != nil {
		c.Abort(ctx, err)

		return
	}

	ctx.JSON(http.StatusOK, genres)
}

func (c GenresController) Create(ctx *Context) {
	context := c.CreateContext(ctx)

	var body input.GenreCreateInput
	if err := ctx.ShouldBindJSON(&body); err != nil {
		c.AbortBadRequest(ctx, err)

		return
	}

	genres, err := c.genreCreateInteractor.Invoke(context, body)
	if err != nil {
		c.Abort(ctx, err)

		return
	}

	ctx.JSON(http.StatusCreated, genres)
}

func NewGenresController(
	base ApplicationController,
	genreSearchInteractor interactor.GenreSearchInteractor,
	genreCreateInteractor interactor.GenreCreateInteractor,
) GenresController {
	return GenresController{
		ApplicationController: base,
		genreSearchInteractor: genreSearchInteractor,
		genreCreateInteractor: genreCreateInteractor,
	}
}

```

Interactor層
本来はvalidateチェックなどをすべきですが、今回は設計とページネーション関数について説明したいので省きます。
```go:genre_interactor.go
package interactor

import (
	"context"

	pkgErr "github.com/example/pkg/errors"

	"github.com/example/pkg/domain/model"
	"github.com/example/pkg/usecase/input"
	"github.com/example/pkg/usecase/output"
	"github.com/example/pkg/usecase/repository"
	"github.com/example/pkg/usecase/transaction"
	"github.com/example/pkg/usecase/validator"
)

type GenreSearchInteractor struct {
	genreRepository repository.GenreRepository
	presenter       output.GenreSearchPresenter
}

func (i GenreSearchInteractor) Invoke(ctx context.Context) (*output.Genres, error) {
	genres, err := i.genreRepository.Search(ctx)
	if err != nil {
		return i.presenter.Error(err)
	}

	return i.presenter.Output(genres)
}

func NewGenreSearchInteractor(
	genreRepository repository.GenreRepository,
	presenter output.GenreSearchPresenter,
) GenreSearchInteractor {
	return GenreSearchInteractor{
		genreRepository: genreRepository,
		presenter:       presenter,
	}
}

type GenreCreateInteractor struct {
	transaction     transaction.Transaction
	genreRepository repository.GenreRepository
	presenter       output.GenreCreatePresenter
}

func (i GenreCreateInteractor) Invoke(ctx context.Context, body input.GenreCreateInput) (*output.Genres, error) {
	var (
		genre  model.Genre
		genres []*model.Genre
		exist  bool
		e      error
	)

	if err := i.transaction.Tx(ctx, func(ctx context.Context) error {
		genre = model.Genre{
			GenreCode:      body.GenreCode,
			GenreName:      body.GenreName,
			GenreRomanName: body.GenreRomanName,
		}

		if _, e = i.genreRepository.Create(ctx, &genre); e != nil {
			return e
		}

		if genres, e = i.genreRepository.Search(ctx); e != nil {
			return e
		}

		return nil
	}); err != nil {
		return i.presenter.Error(err)
	}

	return i.presenter.Output(genres)
}

func NewGenreCreateInteractor(
	transaction transaction.Transaction,
	genreRepository repository.GenreRepository,
	presenter output.GenreCreatePresenter,
) GenreCreateInteractor {
	return GenreCreateInteractor{
		transaction:     transaction,
		genreRepository: genreRepository,
		presenter:       presenter,
	}
}
```

usecase/repositoryのインターフェース
```go:genre_repository.go(interfaces)
package repository

import (
	"context"

	"github.com/example/pkg/domain/model"
)

type GenreRepository interface {
	Search(ctx context.Context) ([]*model.Genre, error)
	Create(ctx context.Context, genre *model.Genre) (*model.Genre, error)
}
```

interfaces/gatewayにあるrepositoryの実装ファイル
ORMはGORMを使用しております。
詳しく知りたい方は[公式サイト](https://gorm.io/ja_JP/docs/index.html)をご覧ください。
```go:genre_repository.go(impl)
package repository

import (
	"context"

	"github.com/example/pkg/domain/model"
	"github.com/example/pkg/interfaces/gateway/database"
	repo "github.com/example/pkg/usecase/repository"
)

type GenreRepository struct {
	db database.Client
}

func (r *GenreRepository) Search(ctx context.Context) ([]*model.Genre, error) {
	var (
		err    error
		conn   database.Connection
		genres []*model.Genre
	)

	if conn, err = r.db.Reader(ctx); err != nil {
		return nil, err
	}

	q := conn.Query().From("genres").OrderBy("genre_code")

	if err = q.Find(&genres); err != nil {
		return nil, err
	}

	return genres, nil
}

func (r *GenreRepository) Create(ctx context.Context, genre *model.Genre) (*model.Genre, error) {
	var (
		err  error
		conn database.Connection
	)

	if conn, err = r.db.Writer(ctx); err != nil {
		return nil, err
	}

	if err = conn.Save(genre); err != nil {
		return nil, err
	}

	return genre, nil
}


func NewGenreRepository(db database.Client) repo.GenreRepository {
	return &GenreRepository{
		db: db,
	}
}
```


Entities層
```go:genre.go
package model

type Genre struct {
	ID        uint64    `json:"id" gorm:"primary_key;"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	GenreCode      string  `json:"GenreCode" gorm:"not null; size:20"`
	GenreName      string  `json:"GenreName" gorm:"not null; size:50"`
	GenreRomanName *string `json:"GenreRomanName" gorm:"size:50"`
}

```

Input層
```go:genre.go
package input

type GenreCreateInput struct {
	GenreCode      string  `json:"genreCode" validate:"required,max=20"`
	GenreName      string  `json:"genreName" validate:"required,max=50"`
	GenreRomanName *string `json:"genreRomanName" validate:"omitempty,max=50"`
}
```

Output層
```go:genre.go
package output

import "github.com/satoholdings-ripple-api/pkg/domain/model"

type Genres struct {
	Genres []GenreDTO `json:"genres"`
}

type Genre struct {
	Genre GenreDTO `json:"genre"`
}

type GenreDTO struct {
	ID             uint64  `json:"id"`
	GenreCode      string  `json:"genreCode"`
	GenreName      string  `json:"genreName"`
	GenreRomanName *string `json:"genreRomanName"`
}

type GenreSearchPresenter interface {
	Output(genres []*model.Genre) (*Genres, error)
	Error(err error) (*Genres, error)
}

type GenreCreatePresenter interface {
	Output(genres []*model.Genre) (*Genres, error)
	Error(err error) (*Genres, error)
}
```

Presenter層
```go:genre.go
package presenter

import (
	"github.com/satoholdings-ripple-api/pkg/domain/model"
	"github.com/satoholdings-ripple-api/pkg/usecase/output"
)

type genreSearch struct{}

func (d genreSearch) Output(genres []*model.Genre) (*output.Genres, error) {
	var outputGenres []output.GenreDTO
	for _, genre := range genres {
		outputGenres = append(outputGenres, TransformGenreDTO(genre))
	}

	return &output.Genres{
		Genres: outputGenres,
	}, nil
}

func (d genreSearch) Error(err error) (*output.Genres, error) { return nil, err }

func NewGenreSearchPresenter() output.GenreSearchPresenter {
	return genreSearch{}
}

type genreCreate struct{}

func (d genreCreate) Output(genres []*model.Genre) (*output.Genres, error) {
	var outputGenres []output.GenreDTO
	for _, genre := range genres {
		outputGenres = append(outputGenres, TransformGenreDTO(genre))
	}

	return &output.Genres{
		Genres: outputGenres,
	}, nil
}

func (d genreCreate) Error(err error) (*output.Genres, error) { return nil, err }

func NewGenreCreatePresenter() output.GenreCreatePresenter {
	return genreCreate{}
}

func TransformGenreDTO(genre *model.Genre) output.GenreDTO {
	return output.GenreDTO{
		ID:             genre.ID,
		GenreCode:      genre.GenreCode,
		GenreName:      genre.GenreName,
		GenreRomanName: genre.GenreRomanName,
	}
}

```