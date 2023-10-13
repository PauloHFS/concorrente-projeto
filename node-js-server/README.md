# Node JS Server

## How to run

> [!IMPORTANT]
> This project uses nodejs, docker and docker-compose to run
> You should have this tools installed in your machine

1. Install dependencies

    ```bash
    yarn
    ```

2. Setup de Infrastructure

    ```bash
    docker run --name concorrente-projeto-db -p 5432:5432 -e POSTGRES_PASSWORD=postgres -d postgres
    ```

3. Run the application

    ```bash
    yarn dev
    ```

<!-- > [!NOTE]
> You should turn off the application before turn off the infrastructure at end
> do this running `docker-compose down` -->

## How to test

```bash
curl --request POST \
  --url http://localhost:8080/resources \
  --header 'Content-Type: application/json' \
  --data '{
 "name": "some name",
 "description": "some description",
 "values": [1,2,3,4]
}'
```
