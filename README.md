# Nexus

Nexus is a Data Distribution System designed for event streaming, dataset access, and value storage. It allows publishers to inform the Nexus server about event streams and datasets, while consumers can request specific data paths.

## Features

1. **Event Streaming (Kafka)**
   - Publishers can register event streams with the Nexus server, providing details such as the server address, data structure of events, and the path to publish.

2. **Dataset Access**
   - Users can run scripts to inform the Nexus server about static datasets, including how to access them and their data structure.

3. **Value Store/Publishing (Redis)**
   - Small values can be sent directly to the Nexus server, allowing consumers to retrieve data directly or act as a cache for event streaming.
  
4. **Go and Python Clients**
   - Go client is available here
   - Python client in [this repo](https://github.com/EliasLittle/NexusPython)

## Installation

- Install Golang: 
    - `sudo apt install golang-go`
- Install protoc: 
    - `sudo apt install -y protobuf-compiler`
- Install go-protoc: 
    - `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`
    - `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`
Run `make all`


## Usage

### Running the Server

To run the Nexus server, use the following command:

```shell
Usage: 'nexus-server <load_file_path> <save_file_path>' or 'nexus-server <save_file_path>'
Options:
  <load_file_path>   Path to the file to load.
  <save_file_path>   Path to the file to save.
  --help             Show this help message.
```

### Running the Client

To run the Nexus client, use the following command:

```shell
make run-client
```

### Publishing Data

You can publish values and datasets using the Nexus client. For example:

```shell
./nexus-client publish value /testing/a 123
./nexus-client publish dataset /testing/dataset/a ./tests/example_a.csv
```

## Example Use Cases

- Single location for all data sharing. Can organize by app, team, user, etc.
- Monitoring
    - Build system progress
    - Git branches
    - Docker container status
    - Kubernetes
- Inspect app internal state
- Share data between apps (pub-sub)
- Save queries
- Data versioning
- Personal scripting   
    - e.g. a script that searches a directory for "TODO:_" statements and then publishes these to /users/$USER/todo/path/of/project
- JIRA integration?


## Roadmap

### General
- [x] Save server state to disk and load from disk on startup
- [ ] Implement user authentication and authorization
- [x] Improve error handling and logging
- [ ] Watch and auto-publish files updates, directories (new files, subdirs), and DBs (new tables)
- [ ] Symbolic path links
    - i.e. /status/my_app/ -> /my_app/prod/bin/status
- [ ] Support functions
    - i.e. publish a bash command or binary to a path. When that path is accessed, the command is executed and the output is returned.
- [ ] Remote hosting
    - [ ] Local data storage (Redis)
    - [ ] Remote data storage (S3, GCS, etc.)
- [ ] Shard Nexus server
    - [ ] Path can point to a specific path on a different server
    - i.e. americas:my/main/path/alt -> europe:my/european/path
- [ ] Documentation for API endpoints
- [ ] Implement data filtering
- [ ] Implement data transformation
- [ ] Implement data compression
- [ ] Implement data encryption
- [ ] Implement data deduplication
- [ ] Seperate out registering and publishing

### Values
- [x] Publishing Values
    - [x] Int32
    - [x] Float64
    - [x] String
- [x] Accessing Values
    - [x] Int32
    - [x] Float64
    - [x] String

### Datasets
- [ ] Registering Datasets
    - [x] Registering individual files
    - [x] Registering directories
    - [ ] Registering datasets from a remote source
    - [x] Registering DB tables
    - [ ] Registering table queries
- [ ] Accessing Datasets
    - [x] Accessing individual files
    - [ ] Accessing directories
    - [ ] Accessing datasets from a remote source
    - [x] Accessing DB tables

### Event Streaming
- [ ] Register __ streams
    - [x] kafka
    - [ ] Websocket
    - [ ] file
    - [ ] DB
    - [ ] direct tcp
- [ ] Access __ streams
    - [x] kafka
    - [ ] Websocket
    - [ ] file
    - [ ] DB
    - [ ] direct tcp
- [ ] Add support for stream filtering

### Yukon
- [x] Minimal navigation and info display
- [x] Type out path
- [x] Command line path argument
- [x] Browse data
    - [x] Browse leaf node
        - [x] Browse values
        - [x] Browse datasets
            - [x] Browse individual files
            - [x] Browse directories
            - [x] Browse DB tables
        - [x] Browse streams
    - [ ] Display data of non-leaf node
- [ ] Display data
    - [ ] Display leaf node
        - [x] Display values
        - [ ] Display datasets
        - [x] Display streams
    - [ ] Display data of non-leaf node
- [ ] Add method to __ data
    - [ ] add
    - [ ] delete
    - [ ] update
    - [ ] query
- [ ] Add help message with keybindings
- [ ] Tab completion if only one child
- [ ] Create a web interface for easier interaction


## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or features you'd like to add.
