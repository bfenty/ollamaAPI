# Ollama Proxy

This is a Go-based proxy server that forwards requests to an Ollama service and includes Prometheus metrics for monitoring.

## Setup

1. **Install Go**: Ensure you have Go installed on your system. You can download it from [https://golang.org/dl/](https://golang.org/dl/).

2. **Clone the Repository**:
   ```sh
   git clone https://github.com/bfent/ollamaAPI.git
   cd ollamaAPI
   ```

3. **Install Dependencies**:
   ```sh
   go mod download
   ```

4. **Set Environment Variables**: Create a `.env` file in the root directory with the following content:
   ```
   OLLAMA_PROXY_KEY=your_api_key_here
   OLLAMA_URL=https://ollama.example.com/
   ```

5. **Run the Application**:
   ```sh
   go run main.go
   ```

6. **Build and Run Docker Container**:
   - Build the Docker image:
     ```sh
     docker build -t ollama-proxy .
     ```
   - Run the Docker container:
     ```sh
     docker run -p 8080:8080 ollama-proxy
     ```
# Ollama Proxy

This is a Go-based proxy server that forwards requests to an Ollama service and includes Prometheus metrics for monitoring.

## Setup

1. **Install Go**: Ensure you have Go installed on your system. You can download it from [https://golang.org/dl/](https://golang.org/dl/).

2. **Clone the Repository**:
   ```sh
   git clone https://github.com/bfent/ollamaAPI.git
   cd ollamaAPI
   ```

3. **Install Dependencies**:
   ```sh
   go mod download
   ```

4. **Set Environment Variables**: Create a `.env` file in the root directory with the following content:
   ```
   OLLAMA_PROXY_KEY=your_api_key_here
   OLLAMA_URL=https://ollama.example.com/
   ```

5. **Run the Application**:
   ```sh
   go run main.go
   ```

The proxy server will start on port `8080`. You can access Prometheus metrics at `/metrics`.

## Notes

- Ensure that your API key and Ollama URL are correctly set in the `.env` file.
- The proxy server logs all requests to the console. For production use, consider setting up a proper logging system.
