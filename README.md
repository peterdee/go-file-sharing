## go-file-sharing

A simple file sharing server written in Golang, without any additional WEB frameworks or routers

Stack:

- Golang

- MongoDB (main storage)

- Redis (cache)

Public features:

- upload files (file size cap is set via environment variable)

- get file info

- download file

- files are stored on the disk

- files are automatically marked as deleted if their view and download stats were not updated for the past 14 days

Manager features:

- view file list (paginated)

- view single file

- delete file (mark as deleted)

- get own user profile

- change password

Root user features:

- create user (root user or manager)

- view user list (paginated)

- view single user

- update user data

- delete user (mark as deleted)

- all of the manager features are also available

### Deploy

Clone repository and install dependencies

```shell script
cd ./go-file-sharing
gvm use go1.24
go mod download
```

### Environment variables

The `.env` file is required, see [.env.example](./.env.example) for details

Setting `MAX_FILE_SIZE_BYTES` variable to `0` disables file size cap 

### Launch

```shell script
go run ./
```

Alternatively can be launched with [AIR](https://github.com/air-verse/air)

Server will be available at http://localhost:9000

### Launch with Docker

Build server application separately:

```shell script
docker build -t server:latest .
```

Run server application separately:

```shell script
docker run server:latest
```

Launch all of the services:

```shell script
dokcer compose up
```

Server will be available at http://localhost:9000

### License

[MIT](./LICENSE.md)
