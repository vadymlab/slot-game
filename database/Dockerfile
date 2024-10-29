FROM migrate/migrate:v4.18.1

ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=postgres
ENV POSTGRES_DB=postgres
ENV POSTGRES_HOST=localhost:5432

COPY ./database/migration /migration/

ENTRYPOINT ["sh", "-c", "/migrate -path /migration -database \"postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}/${POSTGRES_DB}?sslmode=disable\" up"]
