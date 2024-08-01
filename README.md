# VKYC-Backend

## Follow the below steps to run this application:

- Install golang 1.21.6+, refer: https://go.dev
- This application uses PostgreSQL as database, it should be installed and running in the system
- Add a `.env` file in the project's root directly and copy the contents of `.env.example` file into `.env`
- Update the values of `.env` for the mentioned keys
- Application is containerized using docker, to run the application run:
  ```
  make build-run
  ```
- To cleanup docker container run:
  ```
  make clean
  ```
To build the app, run:
  ```
  CGO_ENABLED=0 go build -o vkyc-backend
  ```