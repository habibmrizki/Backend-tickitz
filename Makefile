include ./.env
DBURL=postgres://$(DBUSER):$(DBPASS)@$(DBHOST):$(DBPORT)/$(DBNAME)?sslmode=disable
MIGRATIONPATH=db/migrations

migrate-create :
	migrate create -ext sql -dir $(MIGRATIONPATH) -seq create_$(NAME)_table

migrate-up:
	migrate -database $(DBURL) -path $(MIGRATIONPATH) up

migrate-down:
	migrate -database $(DBURL) -path $(MIGRATIONPATH) down

migrate-force:
	migrate -database $(DBURL) -path $(MIGRATIONPATH) force ${v}

migrate-status:
	migrate -database $(DBURL) -path $(MIGRATIONPATH) version

# Seed database
seed:
	psql -h $(DBHOST) -U $(DBUSER) -d $(DBNAME) -f movie_seeds_full.sql
