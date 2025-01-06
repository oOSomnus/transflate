# Transfalte

A comprehensive platform integrating OCR and translation services, supporting multi-language text recognition and translation.

## Directory Structure

- **backend**: Backend services, including implementations of OCR and translation functionalities.
    - `api`: Defines and generates gRPC service interfaces.
    - `cmd`: Entry points for the services.
    - `engine`: Includes OCR-related data and the Tesseract OCR engine.
    - `internal`: Core business logic and service implementations.
    - `pkg`: Shared utilities and middleware.
    - `bin`: Scripts for Docker builds and proto file generation.

- **frontend**: Frontend service, featuring a React-based interface and Nginx configuration.
## Quick Start

### Prerequisites

- Docker and Docker Compose
- Go compiler (recommended version >= 1.16)
- Node.js and npm

### Deployment Steps

1. Clone the repository:
   ```bash
   git clone https://github.com/oOSomnus/transflate.git
   cd transflateg
   ```

2. Adjust nginx and environmental variables configuration for your own environment.

3. Start all services:
   ```bash
   docker-compose up -d
   ```

### Local Development

#### Backend Development

1. Generate gRPC code:
   ```bash
   cd backend
   make proto
   ```

2. Build the backend services:
   ```bash
   cd backend
   make build
   ```

#### Frontend Development

1. Install dependencies:
   ```bash
   cd frontend
   npm install
   ```

2. Start the development server:
   ```bash
   npm start
   ```

## Features

- **OCR Service**: Provides text recognition capabilities supporting multiple languages (e.g., English, French, Russian).
- **Translation Service**: Integrates a GPT translation module for high-quality text translations.
- **Task Management**: Supports the creation, updating, and querying of user tasks.

## Dependencies

This project uses the following open source software:
- **[Tesseract OCR] https://github.com/tesseract-ocr/tesseract.git**: Used for text recognition, licensed under the Apache License 2.0.

## License

This project is licensed under the [MIT License](./LICENSE).

