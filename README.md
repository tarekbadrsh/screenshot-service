# Shutter - Screenshot as a Service
Shutter is a service that accepts a list of URLs and returns a screenshot of each web page. The service is designed to handle up to 1 million screenshots per day, but it is important to consider the `robots.txt` file of the URLs being processed. For more information on `robots.txt`, please refer to [Wikipedia](https://en.wikipedia.org/wiki/Robots_exclusion_standard).

## How to Use Shutter
Shutter has been built using a dockerized strategy, with all services running in Docker containers managed by `docker-compose`.

### Components
- `receiver`: An API that acts as the entry point for Shutter, responsible for receiving JSON data for processing.
- `Kafka`: A pipeline data management service to handle data between services. Apache Kafka is a popular stream-processing platform. For more information, please refer to [Wikipedia](https://en.wikipedia.org/wiki/Apache_Kafka).
- `screen-shot-service`: The core service of Shutter, responsible for processing URLs and producing screenshots of each URL. The screenshots are stored in the file system.
- `screen-shot-api`: An API that acts as the storage interface for Shutter.
- `screen-shot-db`: A database storage for Shutter that only contains records/meta-data about the screenshots, but not the actual images.
- `deployment`: Shutter includes `docker files`, `docker-compose`, `Makefile`, and `.env` for easy local deployment.

### Build
- To build the images, run `make build` or `docker-compose build`.
- To build from source, each service is built using Go. To build, run `go build .` in the directories for `receiver`, `screen-shot-api`, and `screen-shot-service`. Please use Go version 1.12 or later.

### Run
- To run all services in Docker containers, run `make up` or `docker-compose up`.
- The command `make run` will call `docker-compose down`, `docker-compose build`, and `docker-compose up` in sequence, which is useful in development mode.

#### Notes
- The command `make up` creates a Docker bridge network called `host_machine` for communication between services. It will not work locally with `docker-compose` without this network.
- `Docker-compose` includes all external services required (Kafka, Zookeeper, Scrapy-splash, and PostgreSQL). If running the application manually, make sure to run all external services and update the `.env` or `config.json` with the appropriate configuration.
- Each service has its own configuration file, but there is a main configuration file `.env` that works with `docker-compose` for easier updates.
- Requests can be sent using `test.http` or with `curl`.
#### Single URL Request
```
curl -d '{"url":"https://www.google.com/"}' -H "Content-Type: application/json" -X POST http://localhost:6060/json
```
#### Multiple URL Request
```
/"},{"url":"https://www.youtube.com/"},{"url":"https://www.github.com/"},{"url":"https://www.bitbucket.org/"},{"url":"https://www.chess.com/"},{"url":"https://www.microsoft.com/"},{"url":"https://www.samsung.com/"},{"url":"https://www.apple.com/"},{"url":"error"}]' -H "Content-Type: application/json" -X POST http://localhost:6060/json
```
#### Generated Images Location
The generated images can be found in the **data** directory.
#### check data in the system 
```
curl -H "Content-Type: application/json" -X GET http://localhost:7070/screenshots
```
- you can follow `Makefile` it's include a lot of useful commands e.g.
  - `make logs` to follow the logs from Shutter, and you can get logs from specific service by run `make logs {{service-name}}` e.g. `make logs screen-shot-service`
  - `make psql` to login into `screen-shot-db` database and fetch and filter the data.
  - `make bash` to login into any container while is running and use the container bash e.g. `make bash screen-shot-service`, please note: it will not work with services dockerfile finally build with **scratch** because scratch does not have bash command.

### Down
- run `make down` which will call `docker-compose down` and kill all docker containers.
- or run `docker rm -f $(docker ps -aq)` which will kill docker containers in your system (not only Shutters containers)

# How did I build Shutter

This part is a random account of what happened since I received the task until I reached the current situation. You can skip it and directly look at the code.

## Steps
1. I started by thinking about building the core service (screen-shot-service), which is the screenshot generator.
2. After searching for a screenshot service built in Golang, I found [gowitness](https://en.wikipedia.org/wiki/Apache_Kafka), an awesome library that executes the command `google-chrome` to generate screenshots of any URL.
example to use:

```
docker run --rm -it -v $(pwd)/screenshots:/screenshots leonjza/gowitness:latest single --url=https://www.google.com
```
I decided to try it and after extracting the `chrome` package from `gowitness`, it worked fine, so I continued building the other packages: `main`, `config`, `logger`, and `storage`.
3. The screen-shot-service had to handle up to 1,000,000 screenshots per day, meaning around 1000 images per minute. To accommodate this, I designed the storage strategy to save the images in the file system, with a new directory created each minute. The images are tracked by when they were generated and the URL hash.
4. I used the `md5 algorithm` to ensure that the URL length wouldn't affect the length of the generated image names.
5. Next, I worked on the `receiver` service and continued building the Kafka service, adding Kafka clients to both the screen-shot-service and the receiver.
6. At first, everything worked fine, but with a lot of requests, the `chrome` package in the screen-shot-service couldn't handle the requests. The issue was that Google Chrome couldn't generate more than one screenshot at a time.
7. I then decided to revisit my old tool, `scrapy-splash`, to solve this issue. You can find more information [here](https://github.com/scrapy-plugins/scrapy-splash). I used `scrapy-splash` 2 years ago in my project (Elwizara). `scrapy-splash` is built in Python, but you can call the service through the network. They use `Lua` as a scripting language to render the result, which worked well for me except for a minor issue I faced ([more information here](https://github.com/scrapy-plugins/scrapy-splash/issues/186)).
8. The design pattern of the screen-shot-service made it easy to accept multiple generators. `Scrapy-splash` worked much better than Google Chrome, especially in concurrency mode.
9. After this, I started implementing `screen-shot-api` and `screen-shot-db` to provide a data system interface for Shutter.
10. `Screen-shot-api` was built using a `3-Tier Architecture` (more information [here](https://en.wikipedia.org/wiki/Multitier_architecture)).
11. I used my code generator, [modelgen](https://github.com/tarekbadrshalaan/modelgen), to generate `screen-shot-api`. `Modelgen` is compatible with (mysql, postgres, mssql, sqlite, oracle) and provides a full API implementation that works with the database.

## Summary

This is a summary of how I implemented Shutter. Unfortunately, I didn't have much time to add more features or tests, so feel free to look at the code and provide feedback. I look forward to hearing from you.

## Notes
- Please note that this service is a POC and should not be put into a production environment without stress testing.
- Many parts of this service were built previously, making it easier to put together and not build from scratch.
- The code is open source, so feel free to use it.
