# roku-discovery-agent
Collects metadata from active Roku TV's on the local area network and logs to standard output, with options to publish to a [Google Cloud Pub/Sub Topic](https://cloud.google.com/pubsub)

## Installation
Build binary executable:
```shell script
go build -o agent cmd/agent/main.go
```
## Running
Run the executable:
```shell script
./agent
```

To publish metrics to a Google Cloud Pub/Sub Topic:
```shell script
./agent -p your-project-id -t your-topic-id
```

### Demo
![](./demo.gif)