# API Documentation

## Overview
This documentation provides an overview of the API endpoints available for the Chirpy application. Each endpoint details the request methods, required parameters, and possible response formats.

## API Endpoints

### 1. Get All Posts
- **Method**: GET  
- **Endpoint**: `/api/posts`  
- **Request Parameters**:  
  - `page`: (optional) integer for pagination, default is 1  
  - `limit`: (optional) integer to set the number of posts per page, default is 10  
- **Response Format**:  
  ```json
  {
    "posts": [
      {
        "id": "string",
        "title": "string",
        "content": "string",
        "author": "string",
        "created_at": "string"
      }
    ],
    "total": "integer"
  }
  ```

### 2. Create a New Post
- **Method**: POST  
- **Endpoint**: `/api/posts`  
- **Request Body**:  
  ```json
  {
    "title": "string",
    "content": "string",
    "author": "string"
  }
  ```
- **Response Format**:  
  ```json
  {
    "id": "string",
    "title": "string",
    "content": "string",
    "author": "string",
    "created_at": "string"
  }
  ```

### 3. Update a Post
- **Method**: PUT  
- **Endpoint**: `/api/posts/{postId}`  
- **Request Body**:  
  ```json
  {
    "title": "string",
    "content": "string"
  }
  ```
- **Response Format**:  
  ```json
  {
    "id": "string",
    "title": "string",
    "content": "string",
    "author": "string",
    "created_at": "string"
  }
  ```

### 4. Delete a Post
- **Method**: DELETE  
- **Endpoint**: `/api/posts/{postId}`  
- **Response Format**:  
  ```json
  {
    "message": "string"
  }
  ```

## Error Handling
- **Response Format**:  
  ```json
  {
    "error": {
      "code": "integer",
      "message": "string"
    }
  }
  ```