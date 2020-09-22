# Leader board

This project is build for demo purposes

## Getting started

1. Clone the project
2. Using Docker and docker-compose.
    ```bash
    docker-compose up -d --build
   ```
3. Then the web-service will be available under `http://localhost:8000`.

### Getting authorisation token

1. Open SQLite3 database in `db/leader_board.db`
2. Open `config` table
3. Search for `auth_token` config name and user the value as Bearer token

### Available endpoints
* GET, `/api/v1/leader-board` 
* POST, `/api/v1/player/score` 

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
