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
  - logic to get the data
- interface
  - gateway
    - image uploading interface
  - presenter
    - create line messaing api request
  - controller
  　- routing the message event
- driver
  - image upload api call