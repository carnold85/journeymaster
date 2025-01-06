# JourneyMaster

This repository contains the **JourneyMaster** application, written as part of a coding challenge for a major company. The objective was to implement a RESTful API based on a provided **Swagger/OpenAPI** specification. The program is written in Go and includes features like containerization, deployment scripts, and Kubernetes configuration.

---

## **Features**
- Implements RESTful endpoints as defined in the Swagger specification.
- Secured with environment-based configurations (e.g., credentials, TLS).
- Fully containerized using Docker.
- Deployable locally or in Kubernetes using provided scripts.

---

## **Local Development**

### **Prerequisites**
- [Go](https://golang.org/) (version 1.19 or later)
- [Docker](https://www.docker.com/) (for containerization)
- [Kubernetes CLI (kubectl)](https://kubernetes.io/docs/tasks/tools/) (if using Kubernetes)

### **Steps to Compile and Run Locally**

1. Clone the repository:
   ```bash
   git clone
   cd journeymaster
   ```

2. Initialize Go modules:
   ```bash
   go mod tidy
   ```

3. Run unit tests:
   ```bash
   go test ./...
   ```

4. Build the application:
   ```bash
   go build -o journeymaster
   ```

5. Run the application:
   ```bash
   ./journeymaster
   ```

---

## **Using the Provided Scripts**

### **Build the Application and Container**
Use the `build.sh` script to compile the application and build the Docker container:
```bash
./build.sh
```

### **Deploy Locally or to Kubernetes**
- **Locally**: Use the `deploy.sh` script to deploy the application on your local machine using Docker:
  ```bash
  ./deploy.sh
  ```

- **Kubernetes**: Apply the provided `journeymaster-deployment.yaml` file to deploy the application in a Kubernetes cluster:
  ```bash
  kubectl apply -f journeymaster-deployment.yaml
  ```

---

## **Environment Variables**
The application behavior is controlled by the following environment variables:

| Variable          | Description                                    | Default Value |
|--------------------|------------------------------------------------|---------------|
| `API_PORT`         | Port on which the API listens                 | `8080`        |
| `API_PREFIX`       | Prefix for all API endpoints                  | `/v1`         |
| `API_CREDENTIALS`  | Basic authentication credentials (`username:password`) | `test:test`   |
| `GIN_MODE`         | Gin framework mode (e.g., `release`, `debug`) | `release`     |
| `TLS_PEM`          | Path to the TLS certificate file              | `conf/cert.pem` |
| `TLS_KEY`          | Path to the TLS private key file              | `conf/key.pem` |
| `TLS_ENABLED`      | Whether to enable TLS (`true` or `false`)     | `false`       |