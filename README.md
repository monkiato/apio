[![Build Status](https://drone.monkiato.com/api/badges/monkiato/apio/status.svg?ref=refs/heads/master)](https://drone.monkiato.com/monkiato/apio)


# Apio

A dynamic REST API server, using a manifest file to specify available collections

## Features

 - Multiple collections
 - Specify collection schema (field names and types)
 - Autogenerated generic REST API endpoints (GET, PUT, POST, DELETE)
 - Scheme validations on PUT or POST operations
 - List all available endpoints
 - MongoDB as main database
 
 
## Manifest Declaration

the manifest is a json formatted file with a list of collections to be populated in the API.

The default path is `/app/manifest.json`, or can be changed through the environment variable `MANIFEST_PATH`

expected fields per collection:

 - name: collection name used for the url
 - fields: list of field names and types to be used

e.g.

```go
[
  {
    "name": "books",
    "fields": {
      "title": "string",
      "author": "string",
      "year": "float",
      ...
    }
  },
  ... // more collections
]
```

Available endpoints will be:

```go
GET     http://myurl.com/api/books/{id}
PUT     http://myurl.com/api/books/
POST    http://myurl.com/api/books/{id}
DELETE  http://myurl.com/api/books/{id}
```

A sample file can be found in *manifest.sample.json*


## Available Field Types

 - string
 - float (any numeric field)
 - bool
 
 
## Build Docker Image

Create Docker Image:

`docker build . -t apio-server`


## Run Docker Container

An example is available in docker-compose.yml 

A MongoDB is required for the default storage mode

Custom environment arguments:

    MONGODB_HOST: "{host}:{port}"   //default localhost:27017
    MONGODB_NAME: {db_name}         //default 'apio'
    MANIFEST_PATH: {custom}         //default /app/manifest.json
    DEBUG_MODE: 1                   //default 0, enable verbose logs

A volume mapping is required in order to provide the manifest file:

    volumes:
      - "{your-local-path}:/app/manifest.json"
