# Please set up environment first
export WORST_DB_USERNAME=postgres
export WORST_DB_PASSWORD=postgres
export WORST_DB_NAME=postgres

# https://hub.docker.com/_/postgres
# Then spin up a postgres instance with
sudo docker run -it --name test-pg -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres

# And manually test connection with
psql -h localhost -U postgres postgres

# or
#docker start test-pg

#Then simply do 'go test -v'

# if needed to clean up fresh
# docker kill test-pg
# docker rm test-pg