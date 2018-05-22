# Test Application

## Installation

### Dependencies

Go version 1.9.2

#### Docker

[https://docs.docker.com/docker-for-mac/install/](https://docs.docker.com/docker-for-mac/install/)

#### Database Migration Tool

mattes/migrate

```bash
go get -u -d github.com/mattes/migrate/cli github.com/lib/pq
go build -tags 'postgres' -o $GOPATH/bin/migrate github.com/mattes/migrate/cli
```

#### Postgres Client

```bash
brew install postgres
```

#### Go dep management

Masterminds/glide

```bash
curl https://glide.sh/get | sh
```

## Setup

If you have multiple projects you may want to set up your Go environment

$HOME/project/go/.env

```bash
#!/usr/bin/env bash

# environment
export PATH=$(pwd)/bin:${PATH}
export GOPATH=$(pwd)
```

At the root of your `$GOPATH`

```bash
mkdir -p src/example.com
cd src/example.com
clone https://github.com/alienspaces/test
cd test
```

## Environment

Runtime environment variables are managed in a `.env` file.

```bash
cp .env.example .env
source .env
```

## Services

Scripts to start and stop dependent services such as postgres database.

### Start services

```bash
./dev-bin/start-services
```

### Stop services

```bash
./dev-bin/stop-services
```

## Development

### Install Go package dependencies

```bash
glide install
```

### Build / Run

```bash
./dev-bin/build
test-api
```

### Test

```bash
./dev-bin/test
```
