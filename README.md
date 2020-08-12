# Continuous
Project Continuous is a durable, load balancing distributed system that runs timers. It is implemented using consistent hashing, gRPC, SWIM protocol, and has a REST API 
Usage is simple.
Run Main.Go, enter what kind of connection (LAN,WAN,LOCAL), then enter a text file with known addresses to connect to.

API:
To call the API it will be the address with port 8080
i.e address is 1.2.3.4, 
1.2.3.4:8080 will be used

Create

POST- address:8080/api/v1/create


(startTime and Count are optional)


Body{
  "namespace": "Name",
  "interval": "hh:mm:ss",
  "startTime": "yyyy-mm-dd hh:mm:ss",
  "count": count
}


Success - 201
{
  UUID: uuid
}

Delete

DELETE - address:8080/api/v1/:namespace/:uuid


Success -202
{
  success: Timer Deleted
}

Get

GET - address:8080/api/v1/:namespace/:uuid

Success -202
{
  TimerInfo: TimerInfo
}
