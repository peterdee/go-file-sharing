## go-file-sharing

A simple file sharing server written in Golang, without any additional WEB frameworks or routers

Stack:

- Golang

- MongoDB (store file information)

- Redis (cache file information to reduce main database load)

Features:

- upload files (file size cap is set via environment variable)

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

Setting `MAX_FILE_SIZE_BYTES` variable to `0` disables file size cap 

### Launch

```shell script
go run ./
```

Alternatively can be launched with [AIR](https://github.com/air-verse/air)

Server will be available at http://localhost:9000

### License

[MIT](./LICENSE.md)
