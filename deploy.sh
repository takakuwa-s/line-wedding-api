#!/bin/sh

# Check number of arguments
if test $# -ne 1 ; then
  exit 1
fi

echo "----- show docker containers -----"
docker ps -a

echo "----- show docker images -----"
docker images

echo "----- remove stopped docker containers -----"
docker container prune -f


if test $1 = "local" ; then
  echo "----- local env is selected -----"

  echo "----- remove docker image -----"
  docker rmi line-wedding-api -f

  echo "----- build the docker image -----"
  docker build . -t line-wedding-api

  # echo "----- run the docker container -----"
  # docker run --rm -it line-wedding-api -b localhost -p 8080:8080

elif test $1 = "main" ; then
  echo "----- main env is selected -----"

  echo "----- remove docker image -----"
  docker rmi registry.heroku.com/line-wedding-api/web -f

  echo "----- push the docker image to heroku Container Registry -----"
  heroku container:push web -a line-wedding-api

  echo "----- release the heroku application -----"
  heroku container:release web -a line-wedding-api

  echo "----- show heroku log -----"
  heroku logs -a line-wedding-api --tail
else
  exit 1
fi




