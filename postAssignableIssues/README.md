# アサイン募集中なissueを投稿する

<img width="480" alt="スクリーンショット 2021-12-29 13 53 32" src="https://user-images.githubusercontent.com/19209314/147628553-5d193c28-2418-4afa-8116-40a2c980893e.png">

以下のリポジトリのissueを元に作成
https://github.com/wheatandcat/tool-test/issues


## 仕組み

 - 以下のCloud Functionsを作成
   - 以下の条件でissueを抽出
     - **アサイン募集中** のラベルが設定されているissue
     - アサインがされていないissue
   - [Kibela](https://kibe.la/)  APIを使用して特定の記事を以下の内容で更新
     - アサイン募集中のissue一覧を表示
   - Cloud Scheduleで定期的に自動で記事を更新 

## 運用中のラベル

以下のラベル毎にカテゴリーを分けて表示

|  ラベル名  |  内容  |
| ---- | ---- |
|  RFP  |  Request For Proposalの略。問題を記載しているのでプロポーザルを記載して全体に共有する  |
| ドキュメンテーション  |  ドキュメントにまとめて共有する  |
| 設計前  |  機能の設計書を作成する  |
| 運用改善 |  運用改善系/効率化の対応 |

## デプロイ

```
$ export GITHUB_APP_PRIVATE_KEY=$(cat private-key.pem)

$ gcloud functions deploy PostAssignableIssuesPubSub --set-env-vars GITHUB_APP_PRIVATE_KEY=$GITHUB_APP_PRIVATE_KEY,GITHUB_APP_ID=$GITHUB_APP_ID,GITHUB_OWNER=$GITHUB_OWNER,INSTALLATION_ID=$INSTALLATION_ID,NOTE_TOKEN=$NOTE_TOKEN,NOTE_ID=$NOTE_ID,NOTE_HOST=$NOTE_HOST --runtime go116 --trigger-resource post_assignable_issues --trigger-event google.pubsub.topic.publish --region asia-northeast1
```
