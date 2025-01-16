# Nexus

Nexus is a Data Distribution System designed for event streaming, dataset access, and value storage. It allows publishers to inform the Nexus server about event streams and datasets, while consumers can request specific data paths.

## Features

1. **Event Streaming (Kafka)**
   - Publishers can register event streams with the Nexus server, providing details such as the server address, data structure of events, and the path to publish.

2. **Dataset Access**
   - Users can run scripts to inform the Nexus server about static datasets, including how to access them and their data structure.

3. **Value Store/Publishing (Redis)**
   - Small values can be sent directly to the Nexus server, allowing consumers to retrieve data directly or act as a cache for event streaming.

## Usage

### Running the Server

To run the Nexus server, use the following command:

```shell
make run-server
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

## Roadmap

### General
- [ ] Save server state to disk and load from disk on startup
- [ ] Implement user authentication and authorization
- [ ] Improve error handling and logging
- [ ] Watch and auto-publish files updates, directories (new files, subdirs), and DBs (new tables)
- [ ] Remote hosting
    - [ ] Local data storage (Redis)
    - [ ] Remote data storage (S3, GCS, etc.)
- [ ] Documentation for API endpoints
- [ ] Implement data filtering
- [ ] Implement data transformation
- [ ] Implement data compression
- [ ] Implement data encryption
- [ ] Implement data deduplication

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
    - [ ] Registering directories
    - [ ] Registering datasets from a remote source
    - [ ] Registering DB tables
        - [ ] Registering table queries
- [ ] Accessing Datasets
    - [x] Accessing individual files
    - [ ] Accessing directories
    - [ ] Accessing datasets from a remote source
    - [ ] Accessing DB tables

### Event Streaming
- [ ] Register streams
    - [ ] Register kafka streams
    - [ ] Websocket streams
    - [ ] Register file streams
    - [ ] Register DB streams
    - [ ] Register direct tcp streams
- [ ] Access streams
    - [ ] Access kafka streams
    - [ ] Access file streams
    - [ ] Access DB streams
    - [ ] Access direct tcp streams
- [ ] Add support for stream filtering

### Yukon
- [x] Minimal navigation and info display
- [x] Type out path
- [ ] Display data
    - [ ] Display leaf node
        - [x] Display values
        - [ ] Display datasets
        - [ ] Display streams
    - [ ] Display data of non-leaf node
- [ ] Add help message with keybindings
- [ ] Tab completion if only one child
- [ ] Create a web interface for easier interaction


## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or features you'd like to add.