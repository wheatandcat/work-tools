# テスト内容のあるPRのissueを作成する

## 運用中のラベル

|  ラベル名  |  内容  |
| ---- | ---- |
|  テストを実施  |  テスト対象とするラベル  |

## デプロイ

```
$ export GITHUB_APP_PRIVATE_KEY=$(cat private-key.pem)

$ gcloud functions deploy CreateTestingIssuePub --set-env-vars GITHUB_APP_PRIVATE_KEY=$GITHUB_APP_PRIVATE_KEY,GITHUB_APP_ID=$GITHUB_APP_ID,GITHUB_OWNER=$GITHUB_OWNER,INSTALLATION_ID=$INSTALLATION_ID,SLACK_TOKEN=$SLACK_TOKEN,SLACK_CHANNEL=$SLACK_CHANNEL --runtime go116 --trigger-resource post_assignable_issues --trigger-event google.pubsub.topic.publish --region asia-northeast1
```
