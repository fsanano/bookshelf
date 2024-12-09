# Bookshelf
The API for organizing the books

## Getting Started

### Prerequisites
- Go (version 1.16 or later)

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/fsanano/bookshelf.git
   cd bookshelf
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

### Running the Project

1. Start the API server:
   ```bash
   go run cmd/main.go api
   ```

2. The server will start on the default port `3006`. You can access it at `http://localhost:3006`.

## API Endpoints

### Documentation
https://documenter.getpostman.com/view/13739193/2s83zjri3P#1e841a8b-938a-4c7f-bc5c-5ffc1cd3087a

### Public Endpoints
- **POST** `/signup` - Register a new user

### Protected Endpoints
The following endpoints require authentication:

### User Management
- **GET** `/myself` - Get current user information

### Books Management
- **POST** `/books` - Add a new book
- **GET** `/books` - Get all books
- **GET** `/books/:title` - Search books by title
- **PUT** `/books/:id` - Update a book
- **DELETE** `/books/:id` - Delete a book

### Maintenance
- **GET** `/cleanup` - Run cleanup operations


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

