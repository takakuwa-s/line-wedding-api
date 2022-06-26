#!/bin/sh

# Check number of arguments
if test $# -ne 2 ; then
  echo "number of arguments is not correct"
  exit 1
fi

if [[ "$1" != "wedding-api" ]] && [[ "$1" != "file-upload-api" ]] ; then
  echo "argument 1 is not correct: $1"
  exit 1
fi

if [[ "$2" != "local" ]] && [[ "$2" != "heroku" ]] && [[ "$2" != "docker-hub" ]] ; then
  echo "argument 2 is not correct: $2"
  exit 1
fi

dokcerFilePath="./Dockerfile"
if [[ "$1" == "file-upload-api" ]] ; then
  dokcerFilePath="./DockerfileWithFfmpeg"
fi

echo "----- show docker containers -----"
docker ps -a

echo "----- show docker images -----"
docker images

echo "----- remove stopped docker containers -----"
docker container prune -f

if test $2 = "local" ; then
  echo "----- dwploying to local -----"

  echo "----- remove docker image -----"
  docker rmi $1 -f

  echo "----- build the docker image -----"
  docker build . -t $1 -f $dokcerFilePath --build-arg app=$1

  echo "----- run the docker container -----"
  # docker run --rm -it $1 /bin/bash

  docker run --rm -it $1 -b localhost -p 8080:8080

elif test $2 = "docker-hub" ; then
  echo "----- dwploying to docker-hub -----"
  tag="takakuwa5docker/line-wedding-api"

  echo "----- remove docker image -----"
  docker rmi $tag -f

  echo "----- build the docker image -----"
  docker build . -t $tag -f $dokcerFilePath --build-arg app=$1

  echo "----- push the docker image to docker hub -----"
  docker push takakuwa5docker/line-wedding-api:latest
elif test $2 = "heroku" ; then
  echo "----- dwploying to heroku -----"
  ID="line-wedding-api"
  tag="registry.heroku.com/$ID/web"

  echo "----- remove docker image -----"
  docker rmi $tag -f

  echo "----- build the docker image -----"
  docker build . -t $tag -f $dokcerFilePath --build-arg app=$1

  echo "----- push the docker image to heroku Container Registry -----"
  docker push registry.heroku.com/$ID/web

  echo "----- release the heroku application -----"
  heroku container:release web -a $ID

  echo "----- show heroku log -----"
  heroku logs -a $ID --tail
else
  exit 1
fi