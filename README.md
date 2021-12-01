# Docker Metadata Server
Starter project with REST API and SQL storage

## Local Usage
### Create
```
curl -X POST --data '@metadata.json' -H 'Content-Type=application/json' http://user:pass@localhost:8080/api/docker/metadata
```
### Get Metadata
```
curl http://user:pass@localhost:8080/api/docker/metadata/<record id>
```

### Get All Metadata
```
curl http://user:pass@localhost:8080/api/docker/metadata
```

### Approve Metadata
```
curl -X POST http://user:pass@localhost:8080/api/docker/metadata/<record id>/approve
```

### Deny Metadata
```
curl -X POST http://user:pass@localhost:8080/api/docker/metadata/<record id>/deny
```

## App Limitations
- I did not create any unit tests at this time
- Errors are returned as plan text instead of json
- Errors returned from the API may expose implementation details that shouldn't be exposed
- When retrieving all records, they are all loaded into memory
- No pagination, all records will be returned
- Create call also allows you to immediately approve the metadata, not sure if this is desirable

### Running Docker locally
Use a local `config.env` file, something like this:
```
DATABASEHOST=host.docker.internal
DATABASEPORT=5432
DATABASENAME=docker_metadata
DATABASEUSER=postgres
DATABASESSLMODE=disable
SERVERPORT=8080
BASICAUTHUSER=nobody
BASICAUTHPASS=s3cr3t
```

Run the docker container
```
docker run --name metadata --rm -p 8080:8080 --env-file config.env csquire/go-rest-example:1.0.0
```