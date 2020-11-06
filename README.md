# Go GQL Server for Menagerie

## Usage
Don't bother

Development on this terminated several months before archiving on github,
it has been archived for educational/portfolio/storage reclamation purporses
and I don't really know if it was left in a functioning state, but if you 
define the following env vars when running the Dockerfile:
* PSQL_PORT
* PSQL_HOST
* PSQL_DBNAME
* PSQL_USER
* PSQL_PASSWORD

It might work if that would grant it access to a PGSQL server on the docker net
with the schema described in graph/init.sql - would work with any other SQL
server with a bit of modification to server.go

## Takeaways
* Don't use graphql (for the entire API)
* Don't use Golang (for graphql)

