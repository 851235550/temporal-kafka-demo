## Preparation

1. Make sure GOROOT and GOPATH set correctly.
   Modify and source .bash_profile file.
   ```
   export GOROOT={Your GOROOT}
   export GOPATH={Your GOPATH}
   export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
   ```

2. Install golang dependencies.
   ```
   go get ./...
   ```
----

## Development & Test

### **1. First time**

1. Make sure `docker` is installed and started
2. Make sure `docker-compose` is installed

### **2. Start container**

`docker-compose build && docker-compose up`

And then you can access temporal web UI: http://localhost:8088/


### **3. Create Kafka topic**

1. Exec `docker ps` and copy kafka docker `CONTAINER ID`
2. Exec `docker exec -it {CONTAINER ID} /bin/bash` enter kafka container
3. Create a topic
```
   // In kafka container exec:
   kafka-topics.sh --create --topic temporal-topic --bootstrap-server kafka:9092 --partitions 1 --replication-factor 1
``` 
to create a topic named `temporal-topic`

4. Check topic was created. 
```
// In kafka container exec:
kafka-topics.sh --list --bootstrap-server localhost:9092
```
5. Check has message
```
// In kafka container exec:
kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic temporal-topic --from-beginning
```


### **4. Execute go binary**

To build
```
go build
```
The output binary file is **./suger**, it's for Linux platform.

To run
```
./suger
```
You can also specify workflow count with `--wf-cnt` option. 
```
./suger --wf-cnt=2
```
**--wf-cnt=2** indicates that 2 sub workflows will be initiated for both the production workflow and the consumption workflow

To check

You can see the producer and consumer workflow on the temporal web UI(http://localhost:8088/), showing the parent workflow first, then the child production flow every minute, and the child consumption flow every 2 minutes.


----

## Other Information

### Check consumer
```
kafka-consumer-groups.sh --bootstrap-server localhost:9092 --describe --group consumer-group-1
```
