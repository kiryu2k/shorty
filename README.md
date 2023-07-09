# Shorty
Shorty is a simple URL shortener written in Go using PostgreSQL as storage. The project can be deployed locally and in Docker containers.
The service is implemented in such a way that only one alias can be assigned to one URL.
This means that before trying to shorten URL that already has an alias, a new alias will not be generated.
Instead, an existing one will be returned to the user.
The alias is generated as follows:
1) A hash value is generated using the SHA256 algorithm for the entered URL.
2) Derive a big integer number from the hash value bytes generated during the hasing.
3) Binary-to-text encoding is then applied. Base58 is used as standard encoding.

## Usage
Create a short URL:
```
curl -X POST -d '{"url":"https://github.com/kiryu-dev/shorty"}' localhost:8080/url
```
Redirect from a short URL:
```
curl -IX GET localhost:8080/url/<existing_alias>
```

## Tech Stack
1. Router: [chi](https://github.com/go-chi/chi).
2. DB and stuff: PostgreSQL, [migrate cli util](https://github.com/golang-migrate/migrate), [database/sql golang package](https://pkg.go.dev/database/sql) and [pq driver](https://github.com/lib/pq).
3. Test: [golang test package](https://pkg.go.dev/testing) and [testify lib](https://github.com/stretchr/testify) for unit-testing.
4. Configuration: [cleanenv](https://github.com/ilyakaznacheev/cleanenv) and [godotenv](https://github.com/joho/godotenv).
5. Other: [validator](https://github.com/go-playground/validator), [base58 algorithm](https://github.com/mr-tron/base58).
