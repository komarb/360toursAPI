# 360 tours RESTful API
Created with Golang, Gorilla and MongoDB

## Installation using Docker
Install Docker and in project root directory write this code:
```
docker-compose up -d --build
```
If you’re using Docker natively on Linux, Docker Desktop for Mac, or Docker Desktop for Windows, then the server will be running on
```http://localhost:8080```

If you’re using Docker Machine on a Mac or Windows, use ```docker-machine ip MACHINE_VM``` to get the IP address of your Docker host. Then, open ```http://MACHINE_VM_IP:8080``` in a browser

## Requests
You can get Postman requests collection [here](https://www.getpostman.com/collections/81ff1e587a5d76dd0a5d)

