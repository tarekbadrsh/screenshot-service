# send single url
```
curl -d '{"url":"https://www.google.com/"}' -H "Content-Type: application/json" -X POST http://localhost:6060/json
```

# send multiple urls
```
curl -d '[{"url":"https://www.google.com/"},{"url":"https://www.twitter.com/"},{"url":"https://www.youtube.com/"},{"url":"https://www.github.com/"},{"url":"https://www.bitbucket.org/"},{"url":"https://www.chess.com/"},{"url":"https://www.microsoft.com/"},{"url":"https://www.samsung.com/"},{"url":"https://www.apple.com/"},{"url":"error"}]' -H "Content-Type: application/json" -X POST http://localhost:6060/json
```


# check data in the system 
```
curl  -H "Content-Type: application/json" -X GET http://localhost:7070/screenshots
```