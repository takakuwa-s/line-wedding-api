#!/bin/sh

# Check number of arguments
if test $# -ne 2 ; then
  echo "number of arguments is not correct"
  exit 1
fi

if [[ "$1" != "wedding-api" ]] && [[ "$1" != "background-process-api" ]] ; then
  echo "argument 1 is not correct: $1"
  exit 1
fi

if [[ "$2" != "local" ]] && [[ "$2" != "dev" ]] && [[ "$2" != "prod" ]] ; then
  echo "argument 2 is not correct: $2"
  exit 1
fi

dokcerFilePath="./Dockerfile"
if [[ "$1" == "background-process-api" ]] ; then
  dokcerFilePath="./DockerfileWithFfmpeg"
fi

echo "----- show docker containers -----"
docker ps -a

echo "----- show docker images -----"
docker images

echo "----- remove stopped docker containers -----"
docker container prune -f

if test $2 = "local" ; then
  echo "----- deploying to local -----"

  echo "----- remove docker image -----"
  docker rmi $1 -f

  echo "----- build the docker image -----"
  docker build . -t $1 -f $dokcerFilePath --build-arg app=$1 --build-arg env="dev"

  echo "----- run the docker container -----"
  # docker run --rm -it $1 /bin/bash

  docker run --rm -it $1 -b localhost -p 8080:8080

elif test $2 = "dev" ; then
  echo "----- deploying to docker-hub -----"
  repo="takakuwa5docker/line-wedding-api"
  tag=$1-dev

  echo "----- remove docker image -----"
  docker rmi $repo:$tag -f

  echo "----- build the docker image -----"
  docker build . -t $repo:$tag -f $dokcerFilePath --build-arg app=$1 --build-arg env=dev

  echo "----- push the docker image to docker hub -----"
  docker push $repo:$tag

  echo "----- restart azure app service -----"
  appName=line-wedding-api
  if test $1 = "background-process-api" ; then
    appName=wedding-background-process-api
  fi
  appName=$appName-dev
  az webapp restart --resource-group wedding-$2 --name $appName
else
  echo "----- deploying to docker-hub -----"
  repo="takakuwa5docker/line-wedding-api"
  tag=$1-prod-shimakame

  echo "----- remove docker image -----"
  docker rmi $repo:$tag -f

  echo "----- build the docker image -----"
  docker build . -t $repo:$tag -f $dokcerFilePath --build-arg app=$1 --build-arg env=prod

  echo "----- push the docker image to docker hub -----"
  docker push $repo:$tag

  echo "----- restart azure app service -----"
  appName=line-wedding-api-prod-shima
  if test $1 = "background-process-api" ; then
    appName=wedding-background-process-api-prod-shima
  fi
  az webapp restart --resource-group wedding-prod-shima --name $appName
fi