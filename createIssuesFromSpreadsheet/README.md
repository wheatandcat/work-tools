# スプレットシートからissueを作成する

## 準備

```
$ go mod download
```

## 設定

以下のコマンドで設定ファイルをコピー

```
$ mv config.template.toml config.toml 
```

config.tomlを書き換え

```
[GitHub]
token = ""
owner = ""
repository = ""
url = ""

```



| 名前 | 内容 |
----|---- 
|  token  |  GitHub APIのトークンを設定  |
|  owner  |  オーナー名を設定  |
|  repository  |  レポジトリー名を設定  |
|  url  |  ./gasがデプロイされたURL  |


## 実行

```
$ run main.go
```

