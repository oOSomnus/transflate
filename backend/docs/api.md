# API Documentation

## User APIs

### Login

- **Endpoint:** `/login`
- **Method:** POST
- **Description:** Handles user authentication by validating credentials and generating a JWT token upon successful
  authentication.
- **Request Body:**
  ```json
  {
    "username": "string",
    "password": "string",
    "turnstile_token": "string"
  }
  ```
- **Response:**
    - **200 OK:**
      ```json
      {
        "message": "Login successful",
        "username": "string",
        "token": "string"
      }
      ```
    - **401 Unauthorized:** Invalid credentials or turnstile verification failure.

### Register

- **Endpoint:** `/register`
- **Method:** POST
- **Description:** Handles user registration by validating input and creating a user through the UserUsecase.
- **Request Body:**
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```
- **Response:**
    - **200 OK:**
      ```json
      {
        "message": "User created successfully"
      }
      ```
    - **500 Internal Server Error:** Failure to create the user.

### User Info

- **Endpoint:** `/user/info`
- **Method:** GET
- **Description:** Retrieves information about the authenticated user, including username and account balance.
- **Headers:**
    - `Authorization: Bearer <token>`
- **Response:**
    - **200 OK:**
      ```json
      {
        "username": "string",
        "balance": "integer"
      }
      ```
    - **401 Unauthorized:** Missing or invalid authentication.

## Task APIs

### Task Submit

- **Endpoint:** `/submit`
- **Method:** POST
- **Description:** Handles task submissions by processing an uploaded PDF file, performing OCR and translation, and
  returning a download link.
- **Headers:**
    - `Authorization: Bearer <token>`
- **Request Body:**
    - **Form-Data:**
        - `document`: PDF file
        - `lang`: Language code (optional, default: `eng`)
- **Response:**
    - **200 OK:**
      ```json
      {
        "data": "string (download link)"
      }
      ```
    - **400 Bad Request:** Invalid file format or processing failure.
    - **401 Unauthorized:** User not authorized.
    - **500 Internal Server Error:** Failed to generate download link.

### Task Status Check

- **Endpoint:** Not Implemented
- **Description:** Placeholder for a task status check endpoint.
