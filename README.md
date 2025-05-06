## go-file-sharing

A simple file sharing server written in Golang, without any additional WEB frameworks or routers

Stack:

- Golang

- MongoDB (store file information)

- Redis (cache file information to reduce main database load)

Features:

- upload files (no file size cap)

- get file info via file ID

- download file via file ID

- files are stored on the disk

### Deploy

Clone repository and install dependencies

```shell script
cd ./go-file-sharing
gvm use go1.24
go mod download
```

### Environment variables

The `.env` file is required, see [.env.example](./.env.example) for details

### Launch

```shell script
go run ./
```

Alternatively can be launched with [AIR](https://github.com/air-verse/air)

Server will be available at http://localhost:9000

### License

[MIT](./LICENSE.md)
