# Transflate

A comprehensive platform integrating OCR and translation services, supporting multi-language text recognition and translation.

Currently hosted on <https://translaterequest.com>

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Go compiler (recommended version >= 1.16)
- Node.js and npm

### Deployment Steps

1. Clone the repository:
   ```bash
   git clone https://github.com/oOSomnus/transflate.git
   cd transflate
   ```

2. Adjust nginx and viper configuration for your own environment.

3. Start all services:
   ```bash
   docker-compose up -d
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

