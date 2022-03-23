# line-wedding-api

## overview
The wedding application using LINE messaging API.
This project aims a clean architecture.

## contents (plan)

- environments
  - store the environmental variables
- entity
  - store the business logic
    - wedding data such as Course menu, seat, and message template
    - image upload logic
    - check in logic
    - sign up logic
- usecase
  - logic to handle the data
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