DB_USER = ${PETLY_DB_USER}
DB_PASSWORD = ${PETLY_DB_USER_PASSWORD}
DB_NAME = ${PETLY_DB_NAME}

mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
mkfile_dir := $(dir $(mkfile_path))

app.start:
	docker-compose up -d

app.stop:
	docker-compose down -v

migrates.up:
	docker run \
	-v $(mkfile_dir)docker/migrations:/migrations \
	--network host \
	migrate/migrate \
	-path /migrations \
	-database postgresql://$(DB_USER):$(DB_PASSWORD)@localhost:5432/$(DB_NAME)?sslmode=disable \
	up $(num_migrations)

migrates.down:
	docker run \
	-v $(mkfile_dir)docker/migrations:/migrations \
	--network host \
	migrate/migrate \
	-path /migrations \
	-database postgresql://$(DB_USER):$(DB_PASSWORD)@localhost:5432/$(DB_NAME)?sslmode=disable \
	down $(num_migrations)

migrates.create:
	docker run \
	-v $(mkfile_dir)docker/migrations:/migrations \
	migrate/migrate \
	create \
	-dir /migrations \
	-ext sql \
	$(migration_name)