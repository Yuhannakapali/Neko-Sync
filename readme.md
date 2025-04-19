# Neko-Sync

Neko-Sync is a Go-based application that allows users to stream movies and music seamlessly. It is designed to provide a smooth and efficient media streaming experience.

## Features

- **Movie Streaming**: Stream your favorite movies in high quality.
- **Music Streaming**: Enjoy your favorite tracks with minimal buffering.
- **Cross-Platform Support**: Works on multiple platforms.
- **User-Friendly Interface**: Simple and intuitive design for easy navigation.
- **High Performance**: Built with Go for speed and reliability.

## Installation

1. Clone the repository:
  ```bash
  git clone https://github.com/yourusername/Neko-Sync.git
  ```
2. Navigate to the project directory:
  ```bash
  cd Neko-Sync
  ```
3. Build the project:
  ```bash
  go build
  ```
4. Run the application:
  ```bash
  ./Neko-Sync
  ```

## Usage

1. Open the application in your browser at `http://localhost:8080`.
2. Browse and select movies or music to stream.
3. Enjoy your media!

## Requirements

- Go 1.18 or later
- A modern web browser

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with ❤️ using Go.
- Inspired by the need for a lightweight media streaming solution.
- Thanks to the open-source community for their support.

## Project Structure

The project follows a modular structure to ensure maintainability and scalability:

```
Neko-Sync/
│
├── cmd/
│   └── nekosync/                  # Main entrypoint for the application (main.go)
│
├── internal/
│   ├── app/                       # Application layer (use cases, service orchestration)
│   │   ├── sync/                  # Example: sync-related use cases
│   │   └── user/                  # Example: user-related use cases
│   │
│   ├── domain/                    # Core domain layer (entities + interfaces)
│   │   ├── sync/                  # Sync entity and interfaces (repo, service)
│   │   └── user/                  # User entity and interfaces
│   │
│   ├── infra/                     # Infrastructure (db, email, third-party, etc.)
│   │   ├── db/                    # PostgreSQL setup + migrations
│   │   ├── persistence/           # Actual implementation of repo interfaces
│   │   └── logger/                # Logger setup ()
│   │
│   ├── interfaces/                # Delivery layer (HTTP handlers, controllers)
│   │   ├── http/                  # Echo HTTP handlers
│   │   │   ├── middleware/
│   │   │   ├── routes.go
│   │   │   └── handlers.go
│   │   └── grpc/                  # Optional gRPC support
│   │
│   └── config/                    # App configs, env loading
│
├── migrations/                    # SQL migration files
│
├── go.mod
├── go.sum
└── .env
```

This structure separates concerns and organizes the codebase into logical layers, making it easier to navigate and extend.