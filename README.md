# Leader board

This project is build for demo purposes

## Getting started

### Docker setup
1. Clone the project
2. Use Docker and docker-compose.
    ```bash
    # In project root 
    docker-compose up -d --build
   ```
3. Then the web-service will be available under `http://localhost:8000`.

### Native setup
1. Clone the project
2. Install dependencies
    ```bash
    # In project root
    go get .
    ```
3. Prepare SQLite3 database
    ```bash
    # In project root
   ./db/migrate.sh
   ```
4. Build project or run straight away
    ```bash
    go build -o leader-board .
   # Then run a binary file
   ./leader-board
   
   # OR
   
   go run main.go
   ```

### Getting authorisation token

Default token - 7191ba6933d2ea4d775dd31f6ea351abf794b9fe
#### Configure authorisation token
1. Open SQLite3 database in `db/leader_board.db`
2. Open `config` table
3. Search for `auth_token` config name and use the value as Bearer token

### Available endpoints
* GET, `/api/v1/leader-board`, available query parameters: page<int>, period<string>: monthly|all-time(default)
* POST, `/api/v1/player/score` 

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
