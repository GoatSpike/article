## 初めに
現在、クリーンアーキテクチャ設計を元にして開発している案件においてページネーション関数を実装する機会があったので自分用のメモとして記事を書きました。
拙い文章もあるかとは存じますがご了承ください。


### 1.ページネーション関数

interactor配下に下記のファイルを追加。

```go:pagination_interactor.go
package interactor

import (
	"math"

	"github.com/example/pkg/domain/model"
	"github.com/example/pkg/interfaces/gateway/database"
)

func Paginate(q database.Queryable, page, limit int) (database.Queryable, *model.PageInfo, error) {
	var (
		totalCount int64
		err        error
	)

	// ページ番号とリミットが 0 以下の場合はデフォルト値を設定
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	if totalCount, err = q.Count(); err != nil {
		return nil, nil, err
	}

	pageInfo := &model.PageInfo{
		TotalCount:  int(totalCount),
		CurrentPage: page,
	}

	// 最後のページ番号を計算
	pageInfo.LastPage = int(math.Ceil(float64(totalCount) / float64(limit)))

	// 次のページが存在するかどうかを判定
	pageInfo.HasNext = page < pageInfo.LastPage

	// オフセットを計算
	offset := (page - 1) * limit

	// クエリに Limit と Offset を設定
	q = q.Limit(limit).Offset(offset)

	return q, pageInfo, nil
}
```


### 2.実装ファイルの例
下記の例はconfigファイルやwireの実装、及び説明を省いております。
案件ではwireを使って依存関係を管理したり、DB接続の情報などはconfigファイルに記述すべきですがそれだと説明が非常に長くなってしまった為、流れが把握できる最小限のコード絞っております。



apiのセットアップファイル
```go:api.go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/example/pkg/interfaces/api"
)

type API struct {
	Engine *gin.Engine
	Port   string
}

func NewAPI(
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


func NewGenresController(
	base ApplicationController,
	genreSearchInteractor interactor.GenreSearchInteractor,
) GenresController {
	return GenresController{
		ApplicationController: base,
		genreSearchInteractor: genreSearchInteractor,
	}
}

```

Interactor層
本来はvalidateチェックなどをすべきですが、今回は設計とページネーション関数について説明したいので省きます。
```go:genre_interactor.go
package interactor

import (
	"context"

	"github.com/example/pkg/domain/model"
	"github.com/example/pkg/usecase/output"
	"github.com/example/pkg/usecase/repository"
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
	GenreCode      string  `json:"GenreCode" gorm:"not null; size:20"`
	GenreName      string  `json:"GenreName" gorm:"not null; size:50"`
	GenreRomanName *string `json:"GenreRomanName" gorm:"size:50"`
    CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
```

```go:domain/page_info.go
package model

type PageInfo struct {
	Count       int  `json:"count"`
	TotalCount  int  `json:"totalCount"`
	HasNext     bool `json:"hasNext"`
	CurrentPage int  `json:"currentPage"`
	LastPage    int  `json:"lastPage"`
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

type PageInfo struct {
	PageInfo PageInfoDTO `json:"pageInfo"`
}

type PageInfoDTO struct {
	Count       int  `json:"count"`
	TotalCount  int  `json:"totalCount"`
	HasNext     bool `json:"hasNext"`
	CurrentPage int  `json:"currentPage"`
	LastPage    int  `json:"lastPage"`
}
```

Presenter層

```go:genre_presenter.go
package presenter
import (
	"github.com/example/pkg/domain/model"
	"github.com/example/pkg/usecase/output"
)

type genreSearch struct{}

func (d genreSearch) Output(genres []*model.Genre, pageInfo *model.PageInfo) (*output.Genres, error) {
	var outputGenres []output.GenreDTO
	for _, genre := range genres {
		outputGenres = append(outputGenres, TransformGenreDTO(genre))
	}

	outputPageIngo := TransformPageInfoDTO(pageInfo)

	return &output.Genres{
		PageInfo: outputPageIngo,
		Genres:   outputGenres,
	}, nil
}

func (d genreSearch) Error(err error) (*output.Genres, error) { return nil, err }

func NewGenreSearchPresenter() output.GenreSearchPresenter {
	return genreSearch{}
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

### 実際の使用例
contorollerでpageとlimitをリクエストから受け取ります。
そしてそれをinteractorに渡してあげます。
```go:genres_controller.go
import (
	"net/http"
	"strconv"

	"github.com/example/pkg/usecase/input"
	"github.com/example/pkg/usecase/interactor"
)
type GenresController struct {
	ApplicationController
	genreSearchInteractor interactor.GenreSearchInteractor
}
func (c GenresController) Search(ctx *Context) {
	context := c.CreateContext(ctx)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	genres, err := c.genreSearchInteractor.Invoke(context, page, limit)
	if err != nil {
		c.Abort(ctx, err)

		return
	}
	ctx.JSON(http.StatusOK, genres)
}
```

```go:genre_interactor.go
func (i GenreSearchInteractor) Invoke(ctx context.Context, page int, limit int) (*output.Genres, error) {
	genres, pageInfo, err := i.genreRepository.Search(ctx, page, limit)
	if err != nil {
		return i.presenter.Error(err)
	}

	return i.presenter.Output(genres, pageInfo)
}
```

```go:genre_repository.go
// 省略
type GenreRepository struct {
	db database.Client
    interactor interactor.pagination_interactor
}

func (r *GenreRepository) Search(ctx context.Context, page int, limit int) ([]*model.Genre, *model.PageInfo, error) {
    var (
		err      error
		conn     database.Connection
		genres   []*model.Genre
		pageInfo *model.PageInfo
	)

    if conn, err = r.db.Reader(ctx); err != nil {
       return nil, nil, err
    }

    q, pageInfo, err = interactor.Paginate(q, page, limit)
    if err != nil {
       return nil, nil, err
    }
    
    q = q.Query().From("genres").OrderBy("genre_code")
	if err = q.Find(&genres); err != nil {
		return nil, err
	}

    pageInfo.Count = len(genres)

    return genres, pageInfo, nil
}

func NewGenreRepository(
    db database.Client,
    interactor interactor.pagination_interactor
) repo.GenreRepository {
	return &GenreRepository{
		db: db,
        interactor: interactor
	}
}
```

ページを1、リミットを5に設定することでDBに保存されてある15のデータのうち、0~4個目のデータを取得することができました。

request
```
curl --location 'http://localhost/genres?page=1&limit=5' 
```

response
```
{
    "pageInfo": {
        "pageInfo": {
            "count": 5,
            "currentPage": 1,
            "hasNext": true,
            "lastPage": 3,
            "totalCount": 15
        }
    },
    "genres": [
        {
            "genreCode": "1",
            "genreName": "ジャンル1",
            "genreRomanName": "genre 1",
            "id": 1
        },
        {
            "genreCode": "2",
            "genreName": "ジャンル2",
            "genreRomanName": "genre 2",
            "id": 2
        },
        {
            "genreCode": "3",
            "genreName": "ジャンル3",
            "genreRomanName": "",
            "id": 3
        },
        {
            "genreCode": "4",
            "genreName": "ジャンル4",
            "genreRomanName": "",
            "id": 4
        },
        {
            "genreCode": "5",
            "genreName": "ジャンル5,
            "genreRomanName": "",
            "id": 5
        }
    ]
}
