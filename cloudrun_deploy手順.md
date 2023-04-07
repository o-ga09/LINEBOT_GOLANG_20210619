## Cloud Run deploy 手順


 シンプルなGo Web サーバをGoogle Cloud Run にデプロイする時の手順


- [ ] 0. 事前準備
    不要なDockerイメージの削除など
~~~bash
docker image ls
docker image rmi gcr.io/cloud-run-test 
docker image rmi taiti09/cloud-run-test
docker container rm taiti09/cloud-run-test
docker container rm 5c83bb584415
docker image rmi taiti09/cloud-run-test
~~~

- [x] 0.プロジェクトIDをセット
~~~
gcloud config set project　todo-app-20221107
~~~

- [x] 1. デプロイ用のDockerイメージをビルド

~~~bash
docker build --platform linux/amd64 -t asia.gcr.io/todo-app-20221107/linebot:latest --target deploy ./
~~~

- [ ] 2. GCPのCloud Registry にプッシュ

~~~bash
docker push asia.gcr.io/todo-app-20221107/linebot:latest
~~~

- [ ] 3. Cloud Run にデプロイ
~~~bash
gcloud run deploy --image asia.gcr.io/todo-app-20221107/linebot:latest --region asia-northeast1 --update-secrets=LINE_CHANNEL_SECRET=channel_secret:latest --update-secrets=LINE_ACCESS_TOKEN=access_token:latest
~~~

- [ ] 4. Cloud Runのサービスを削除

~~~bash
gcloud run services delete --platform managed --region us-west1 simple-go-web-server
~~~