# line-wedding-api

## overview
The wedding application using LINE messaging API.
This project aims a clean architecture.

## contents (plan)

- environments
  - store the environmental variables
- entity
  - business data
- usecase
  - logic for the application
- resource
  - data files
- interface
  - gateway
    - message data repository
    - user data repository
    - file data repository
  - presenter
    - send line messaing
  - controller
  　- routing the line webhook event
- driver
  - router