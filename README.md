# my_go_post_data_crawler

go_crawler_setting.yml
```yml
defaultRoomName: "#timeline"

dsn: username:pass@/tablename?parseTime=true
postPath: http://localhost:8080/hubot/send_message

mongodbUrl: localhost

rootDir: /Users/ota42y/

twitter:
  databaseName: twitter
  collectionName: tweets
  screenNames: [LoveLive_staff]

chatLog:
  logFolder: log/chat/
  saveFolder: archive/chat/
```

(chat log for https://github.com/adragomir/hubot-logger)

twitter.yml
```yml
consumerKey:
consumerSecret:
accessToken:
accessTokenSecret:
```
