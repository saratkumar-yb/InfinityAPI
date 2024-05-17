# InfinityAPI Documentation

## Getting Started

### Step 1: Set Up the Project

First, initialize a new Go project and install the necessary dependencies:

```sh
go mod init infinityapi
go mod tidy
```

### Step 2: Configuration File

Create a file named `config.ini` with the following content:

```ini
[db]
host = localhost
port = 5432
user = saratkumar-yb
password = yourpassword
dbname = yourdbname
sslmode = disable

[server]
http_listener = 0.0.0.0
http_port = 8000
```

#### Configuration Details

#### Database Configuration (`[db]` section)

- `host`: The hostname or IP address of the PostgreSQL server.
- `port`: The port number on which the PostgreSQL server is listening.
- `user`: The username to connect to the PostgreSQL database.
- `password`: The password to connect to the PostgreSQL database.
- `dbname`: The name of the PostgreSQL database to connect to.
- `sslmode`: The SSL mode to use when connecting to the PostgreSQL database (e.g., `disable`, `require`, `verify-ca`, `verify-full`).

#### Server Configuration (`[server]` section)

- `http_listener`: The IP address to listen on (e.g., `0.0.0.0` to listen on all interfaces).
- `http_port`: The port number on which the server will listen for HTTP requests.

### Step 3: Initialize/Migrate the Database

To migrate the database, run the following command:

```sh
go run main.go -migrate
```

This command will read the `schema.sql` file and create the necessary tables in the PostgreSQL database.

### Step 4: Start the Server

To start the API server, run the following command:

```sh
go run main.go -startserver
```

The server will start listening for HTTP requests on the IP address and port specified in the `config.ini` file.

## Sample Data

### Insert Data into `yba` Table

```sh
curl -X POST -H "Content-Type: application/json" -d '{
    "version": "2.20.0.0",
    "type": "yba",
    "architecture": "x86_64",
    "platform": "linux",
    "commit": "xxxxxxxxxxxxx",
    "branch": "2.20"
}' http://localhost:8000/yba
```

### Insert Data into `ybdb` Table

```sh
curl -X POST -H "Content-Type: application/json" -d '{
    "version": "2.20.0.0",
    "type": "ybdb",
    "architecture": "x86_64",
    "platform": "linux",
    "download_url": "http://example.com/download",
    "commit": "xxxxxxxxxxxxx",
    "branch": "2.20"
}' http://localhost:8000/ybdb
```

```sh
curl -X POST -H "Content-Type: application/json" -d '{
    "version": "2.18.0.0",
    "type": "ybdb",
    "architecture": "x86_64",
    "platform": "linux",
    "download_url": "http://example.com/download",
    "commit": "xxxxxxxxxxxxx",
    "branch": "2.18"
}' http://localhost:8000/ybdb
```


### Create Compatibility Relationships

```sh
curl -X POST -H "Content-Type: application/json" -d '{
    "yba_versions": ["2.20.0.0"],
    "ybdb_versions": ["2.20.0.0"]
}' http://localhost:8000/compatibility
```

```sh
curl -X POST -H "Content-Type: application/json" -d '{
    "yba_versions": ["2.20.0.0"],
    "ybdb_versions": ["2.18.0.0"]
}' http://localhost:8000/compatibility
```


### Get Compatible `ybdb` Versions for a Given `yba` Version

```sh
curl -X POST -H "Content-Type: application/json" -d '{
    "yba_version": "2.20.0.0"
}' http://localhost:8000/compatibility_list
```

## Running Tests

To run the tests, use the following command:

```sh
go test -v
```

This command will execute all the tests in the `main_test.go` file and output the results.

### Test Functions

- `TestInsertYbaHandler`: Tests inserting data into the `yba` table.
- `TestInsertYbdbHandler`: Tests inserting data into the `ybdb` table.
- `TestInsertCompatibilityHandler`: Tests creating compatibility relationships.
- `TestGetCompatibleYbdbHandler`: Tests fetching compatible `ybdb` versions for a given `yba` version.

## Building binaries

### Build the Binary for `linux/amd64`

To build the binary for the `linux/amd64` architecture, use the following command:

```sh
GOOS=linux GOARCH=amd64 go build -o infinityapi-amd64 main.go
```

This command sets the `GOOS` environment variable to `linux` and the `GOARCH` environment variable to `amd64`, and then builds the binary, naming the output file `infinityapi-amd64`.

### Build the Binary for `linux/arm64`

To build the binary for the `linux/arm64` architecture, use the following command:

```sh
GOOS=linux GOARCH=arm64 go build -o infinityapi-arm64 main.go
```

This command sets the `GOOS` environment variable to `linux` and the `GOARCH` environment variable to `arm64`, and then builds the binary, naming the output file `infinityapi-arm64`.

## Conclusion

This documentation provides all the necessary steps to configure, start, and test the InfinityAPI server. Follow the instructions carefully to ensure everything is set up correctly.