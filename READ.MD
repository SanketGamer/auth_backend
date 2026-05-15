docker exec -it postgres_db psql -U postgres -d authdb


# cicd
- build the docker image
- push the docker image to docker hub
- ssh into our vm and start new image