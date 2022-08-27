## To show all the network already exists
- `docker network ls`
## To inspect any network
- `docker network inspect <network-name>`
## To create a new network
- `docker network create <new-network-name>`
- `docker network create bank-network`
## To connect any existing container with any network 
- `$ docker network connect <network-name> <container-name>`
-`$ docker network connect bank-network bank-postgres`
## To inspect any container's network
`$ docker container inspect <container>`
- it's totally okay if we connect single container in multiple node
```
docker run --name bank-application --network bank-network  -e GIN_MODE=release -e DB_SOURCE="postgresql://postgres:secret1234@bank-postgres:5432/simple_bank?sslmode=disable" -p 8080:8080 hremon331046/bank-app:latest
```