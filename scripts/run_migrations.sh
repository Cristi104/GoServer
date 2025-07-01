psql -h 127.0.0.1 -p 5432 -U root -d GoServer -c "$(cat ./migrations/*)"
