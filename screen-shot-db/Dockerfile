FROM postgres

## copy sql file into the container 
COPY init.sql /docker-entrypoint-initdb.d/init.sql

## NOTE: SQL files will be executed in alphabetical order
CMD ["postgres"]