# Petly application

----

### To run the application, you need to have docker installed on your computer. Nothing else is needed, since all parts are automated, containerized and run by one common command.

#### To run the application please make an export of necessary environment variables with command bellow. You can change values if you want.

```shell
export POSTGRES_PASSWORD=some-password && \
export PETLY_DB_USER=petly_user && \
export PETLY_DB_USER_PASSWORD=super-secret-db-password && \
export PETLY_DB_NAME=petly && \
export APP_PORT=8080 && \
export JWT_TOKEN_SECRET=super-secret-secret && \
export JWT_TOKEN_EXPIRATION_TIME=24h && \
export USER_PASSWORD_SALT=super-secret-password-salt
```

#### Then execute the command `make app.start` and application will be started using docker-compose.
#### If you did not change the port, then the webserver will be available at `http://localhost:8080`
#### Swagger and endpoints description you can find by the address `http://localhost:8080/swagger/index.html`

## Application description:
Application has a very simple sign-in/sign-up functionality and authorization based on JWT token (see swagger how to use).
1. User can create as much as he wants dogs. 
2. User can like/dislike dogs of another users.
3. User can see matches with another dogs.

