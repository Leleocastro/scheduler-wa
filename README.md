# INTRODUCTION

## To run the project, follow the steps below:

- Clone the repository
- Go inside the kong folder run the commands:
  - `docker-compose up -d`
  - `KONG_DATABASE=postgres docker compose --profile database up -d`
- Go to the root and run the commands:
  - `docker-compose up -d`
- Now go inside the kong container with the command:
  - `docker exec -it api-gateway bash`
- and run the command:
  - `kong config db_import /opt/kong/config.yaml`
