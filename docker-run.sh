#!/bin/bash

## Supprimez le conteneur existant s'il existe
if [ "$(docker ps -a -q -f name=color-generator)" ]; then
  docker rm -f color-generator
fi

# Exécutez le conteneur Docker avec ou sans le volume en fonction de la condition précédente
docker run -v ./app/data:/app/data -p "8080:8080" --name color-generator color-generator
