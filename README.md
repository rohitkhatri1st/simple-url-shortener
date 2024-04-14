# simple-url-shortener
This is a simple URL shortener application written in Go. It provides functionality to shorten long URLs into shorter, more manageable ones.

## Features

- Shorten long URLs into short, easy-to-share links.
- Retrieve the original long URL from a shortened link.
- Only one short URL will be generated for each long URL, ensuring that the same long URL always maps to the same short URL.
- Get statistics on top domains used for shortened URLs.

## Getting Started

### Prerequisites

- Go programming language installed on your system. You can download it from [here](https://golang.org/dl/).
- Git installed on your system to clone the repository.

### Installation

1. Clone the repository to your local machine:

```bash
git clone https://github.com/rohitkhatri1st/simple-url-shortener.git
```

2. Navigate to the project directory:
```bash
cd simple-url-shortener
```

3. Build the project:
```bash
go build .
```

4. Run the executable:
```bash
./simple-url-shortener
```
By default, the server will start running on port 8001.

#### Running the application with Docker
To run the application using Docker:

1. Make sure you have Docker installed on your system.
2. Clone the repository to your local machine (if you haven't already using above mentioned steps):
3. Navigate to the project directory using above mentioned steps:
4. Build the Docker image: 
```bash
docker build -t simple-url-shortener .
```
If above command doesn't work you may try:
```bash
docker build --pull --rm -f "dockerfile" -t simple-url-shortener:latest "."
```
5. Run the Docker container:
```bash
docker run -p 8001:8001 simple-url-shortener
```
By default, the server will start running on port 8001.


### Mocking
You can generate mocks using mockery package of golang.

Use the below command to generate mocks.
```bash
mockery --dir=server/storage --name=InMemoryDb --filename=in_memory_db_mock.go --output=server/mocks --outpkg=mocks
```

## Usage
Once the server is running, you can interact with it using HTTP requests. Here are the available endpoints:

- POST /shorten: Shorten a long URL. Send a JSON payload with the original URL to shorten.
- GET /original/{shortKey}: Retrieve the original long URL associated with a given short key.
- GET /top-domains: Get statistics on top domains used for shortened URLs.

Example usage:
```bash
curl -X POST -H "Content-Type: application/json" -d '{"originalUrl": "http://example.com"}' http://localhost:8001/shorten
```

## Testing
To run the tests for the project, navigate to the project directory and run:
```bash
go test ./...
```
This will run all the tests in the project.

## API Documentation
You may find the postman API documentation [here](https://documenter.getpostman.com/view/32520716/2sA3Bj7Yzc).
