# Golang Backend Server

This repository contains a Golang backend server designed to work alongside a frontend, a database, and a microservice. It's primarily built using the Gin framework, and the backend handles request validation, passing it to a microservice, and interacts with the database as required.

## Features

- Validaion, processing requests from the frontend, sending replies.
- Interaction with a microservice.
- Read/Write operations with a database.
- JWT tokens for authorization.

## Tech Stack

- **Golang**: The primary programming language used.
- **Gin**: A web framework used for building the API.
- **jmoiron/sqlx**: For database interaction.
- **lib/pq**: A pure Go Postgres driver for Go's database/sql package.


## Architecture

1. General API Level: Has 3 specific APIs. This level is responsible for validating HTTP requests, setting the response's status code, and transferring it to services.
2. Models: Define the structure of incoming and outgoing data.
3. Services: Where the business logic begins. This includes interactions with the microservice, database operations, etc.

## Endpoints

### Authentication
- POST /auth/signup_number: Handle sign up with a phone number.
- POST /auth/refresh: Refresh authentication tokens.
- POST /auth/signin_number: Handle sign in with a phone number.
- POST /auth/code: Handle code from the phone.
- POST /auth/2fa: Authenticate with two-factor authentication.
- GET /auth/check_client: Check validity of tokens in cookies.

### Channels
- GET /channels/get_feed: Retrieve the feed channels.
- POST /channels/add_feed: Add channels to the feed.
- DELETE /channels/remove_feed: Remove сhannels to the feed.

### User
- GET /user/get_subs: Retrieve all user`s public subscribed channels.
- GET /user/form_feed: Form a feed based on added channels.
- GET /user/get_recommendation: Get recommendations for the user.
