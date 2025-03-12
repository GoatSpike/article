

前提
myboxは2s

詳細は元々APIを叩くと4s前後かかっている
他の方のPCでは25~30s
原因は不明なのですが、モデルの修正やpreloadの見直しを行う

modelはDBと1:1
outputはjsonの型

やったこと
## step1
modelの再構築

具体的にはDBと照らし合わせ、リレーションがある構造体をモデルに持たせるよう修正。
json tagの廃止
gorm tag の追加。`foreignKey`は使用するが、referenceは使用しない


Mybox
```
curl -w "http_code: %{http_code}\n\
time_total: %{time_total}\n\
time_namelookup: %{time_namelookup}\n\
time_connect: %{time_connect}\n\
time_appconnect: %{time_appconnect}\n\
time_starttransfer: %{time_starttransfer}\n\
" --location 'http://localhost:8100/api/v1/ripples/1587' \
--header 'x-sato-accesstoken: - eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiNDA5ZDQ4YTctZWQ1OS00ZTBmLTg3MDAtNTNiMDFlNjE3YTgxIiwic2lkIjoiZTVlNzc3ZTYtMmEzZC00NTM2LWI4NTMtMDhjNjZjMDQyZjgyIiwiYWRtaW4iOmZhbHNlLCJleHAiOjE3NDE3NzE3OTJ9.9UPvWaf5LmQBHE3zFswyaoKDqVG-b9IIuJex0YetAoU'
```

# HTTP Response Data
| 回数 | http_code | time_total | time_namelookup | time_connect | time_appconnect | time_starttransfer |
|:----:|:---------:|:----------:|:---------------:|:------------:|:---------------:|:------------------:|
|  1   |    200    |  2.115472  |     0.000017    |    0.000524  |     0.000000    |      2.111988      |
|  2   |    200    |  2.085498  |     0.000014    |    0.000503  |     0.000000    |      2.083643      |
|  3   |    200    |  2.950818  |     0.000018    |    0.000747  |     0.000000    |      2.947696      |
|  4   |    200    |  1.799089  |     0.000009    |    0.000382  |     0.000000    |      1.796315      |
|  5   |    200    |  2.272939  |     0.000010    |    0.000477  |     0.000000    |      2.270973      |
|  平均 |    200    |  2.2443632 |     0.0000136   |    0.0005266 |     0.000000    |      2.242923      |

## step2
outputのmodelの再構築
適切なモデルに修正この時はあまり早くならない

## step3
MyBoxのpresenterの初期化の方法を例のように変えてみる
```
var ripples []output.Ripple
↓
ripples := make([]output.Ripple, 0, len(rippleReportPair))
```

変更したところは5箇所
# HTTP Response Data
| 回数  | http_code | time_total  | time_namelookup |  time_connect | time_appconnect | time_starttransfer |
|:----:|:---------:|:-----------:|:----------------:|:--------------:|:---------------:|:------------------:|
|  1  |  200  |  1.027011  |  0.000018  |  0.000551  |  0.000000  |  1.023025  |
|  2  |  200  |  2.133748  |  0.000020  |  0.000621  |  0.000000  |  2.122284  |
|  3  |  200  |  1.226195  |  0.000023  |  0.000638  |  0.000000  |  1.218340  |
|  4  |  200  |  1.300406  |  0.000013  |  0.000369  |  0.000000  |  1.291534  |
|  5  |  200  |  1.414999  |  0.000015  |  0.000559  |  0.000000  |  1.405155  |
| 平均 |  200  |  1.4204718  |  0.0000178  |  0.0005476  |  0.000000  |  1.4128676  |


## step1


詳細は約4sから下記のように推移
```
curl -w "http_code: %{http_code}\n\
time_total: %{time_total}\n\
time_namelookup: %{time_namelookup}\n\
time_connect: %{time_connect}\n\
time_appconnect: %{time_appconnect}\n\
time_starttransfer: %{time_starttransfer}\n\
" --location 'http://localhost:8100/api/v1/ripples/1587' \
--header 'x-sato-accesstoken: - eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiNDA5ZDQ4YTctZWQ1OS00ZTBmLTg3MDAtNTNiMDFlNjE3YTgxIiwic2lkIjoiZTVlNzc3ZTYtMmEzZC00NTM2LWI4NTMtMDhjNjZjMDQyZjgyIiwiYWRtaW4iOmZhbHNlLCJleHAiOjE3NDE3NzE3OTJ9.9UPvWaf5LmQBHE3zFswyaoKDqVG-b9IIuJex0YetAoU'
```

# HTTP Response Data

| 回数 | http_code | time_total | time_namelookup | time_connect | time_appconnect | time_starttransfer |
|:----:|:---------:|:----------:|:---------------:|:------------:|:---------------:|:------------------:|
|  1   |    200    |  3.036293  |     0.000016    |   0.000475   |     0.000000    |       3.034043     |
|  2   |    200    |  2.854623  |     0.000020    |   0.000567   |     0.000000    |       2.851780     |
|  3   |    200    |  4.443639  |     0.000012    |   0.000263   |     0.000000    |       4.439579     |
|  4   |    200    |  3.851207  |     0.000018    |   0.000483   |     0.000000    |       3.849224     |
|  5   |    200    |  2.737645  |     0.000015    |   0.000455   |     0.000000    |       2.734677     |
| 平均 |    200    | 3.3846814  |    0.0000162    |  0.0004486   |     0.000000    |      3.3814606     |

ちなみに25sくらいかかるpcだと23sくらい
```
http_code: 200
time_total: 22.552183
time_namelookup: 0.000016
time_connect: 0.002270
time_appconnect: 0.000000
time_starttransfer: 22.545252
```

## step2
outputのmodelの再構築
適切なモデルに修正この時はあまり早くならない

## step3
MyBoxのpresenterの初期化の方法を例のように変えてみる
```
var ripples []output.Ripple
↓
ripples := make([]output.Ripple, 0, len(rippleReportPair))
```

変更したところは5箇所
# HTTP Response Data
| 回数  | http_code | time_total  | time_namelookup |  time_connect | time_appconnect | time_starttransfer |
|:----:|:---------:|:-----------:|:---------------:|:-------------:|:---------------:|:------------------:|
|  1  |  200  |  2.230859  |  0.000014  |  0.000218  |  0.000000  |  2.228830  |
|  2  |  200  |  1.333433  |  0.000019  |  0.000375  |  0.000000  |  1.329613  |
|  3  |  200  |  1.333433  |  0.000019  |  0.000375  |  0.000000  |  1.329613  |
|  4  |  200  |  1.234983  |  0.000016  |  0.000487  |  0.000000  |  1.232957  |
|  5  |  200  |  1.395880  |  0.000020  |  0.000675  |  0.000000  |  1.393568  |
| 平均 |  200  |  1.5057176  |  0.0000176  |  0.000426  |  0.000000  |  1.5021162  |

