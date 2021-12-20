# アサイン募集中なissueを投稿する



## デプロイ

```
$ export GITHUB_APP_PRIVATE_KEY=$(cat private-key.pem)

$ gcloud functions deploy PostAssignableIssuesPubSub --set-env-vars GITHUB_APP_PRIVATE_KEY=$GITHUB_APP_PRIVATE_KEY,GITHUB_APP_ID=$GITHUB_APP_ID,GITHUB_OWNER=$GITHUB_OWNER,INSTALLATION_ID=$INSTALLATION_ID,VERIFY_ID_TOKEN=$VERIFY_ID_TOKEN,NOTE_TOKEN=$NOTE_TOKEN,NOTE_ID=$NOTE_ID,NOTE_HOST=$NOTE_HOST --runtime go116 --trigger-resource post_assignable_issues --trigger-event google.pubsub.topic.publish --region asia-northeast1
```
