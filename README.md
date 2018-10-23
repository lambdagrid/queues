# LambdaGrid Queues

## Setup Instructions

There are two dependencies of the current prototype: Amazon SQS and Postgres.

### AWS

Get the AWS CLI as documented [here](https://docs.aws.amazon.com/cli/latest/userguide/installing.html). Then, run `aws configure` on the command line. It will ask you for your access key ID and secret. You can find these in the [IAM console](https://console.aws.amazon.com/iam/home?region=us-west-2#/users) under the user's Security Credentials tab. You may want to create a user with SQS permissions. Right now, that's the only permissions it needs.

Once you have configured the `aws cli`, the AWS portion of the setup is done.

### Postgres

Get a Postgres database running locally. If you are on Mac OSX, one easy
way to do this is with [Postgres.app](https://postgresapp.com/).

Once you have the app, you need to create the `lambdagrid` user (usually just with `createuser lambdagrid`) and the
`lambdagrid`database (with `psql -U lambdagrid` and `CREATE DATABASE lambdagrid;`).

Then, you need to execute some SQL statements. Hop into `psql` with the lambdagrid user and lambdagrid DB with `psql -U lambdagrid lambdagrid` and execute the statements below. One easy thing to do would be to
automate the schema creation.

```sql
CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    account_name text NOT NULL UNIQUE,
    auth_key text NOT NULL UNIQUE,
    auth_secret text NOT NULL
);

CREATE UNIQUE INDEX accounts_pkey ON accounts(id int4_ops);
CREATE UNIQUE INDEX accounts_account_name_key ON accounts(account_name text_ops);
CREATE UNIQUE INDEX accounts_auth_key_key ON accounts(auth_key text_ops);

CREATE TABLE queues (
    name text PRIMARY KEY,
    owner_id integer NOT NULL REFERENCES accounts(id) ON DELETE CASCADE ON UPDATE CASCADE,
    queue_url text
);

CREATE UNIQUE INDEX queues_pkey ON queues(name text_ops);
```

### Golang

You need the Go tools installed to build the app. [Follow the instructions on this page](https://golang.org/doc/install#install). It basically consists of downloading the installer for your OS and then creating the directory `$HOME/go`.

### Cloning the repository

Clone the repository into `$HOME/go/src/github.com/lambdagrid/queues`.

### Downloading the dependencies

Run `cd $HOME/go/src/github.com/lambdagrid/queues; go get ./...`.

### Running the app

Run `go run main.go`. The web server will start listening on port 8080.

### Exercising the API

I like to use [Postman](https://www.getpostman.com/) to query the API. You can find the collection of requests I used to test it in the file `lambdagrid.postman_collection.json` in this repository. Use the "Register new account" request to get an API key and secret. Then edit the request group variables to set these as the `authKey` and `authSecret`. The rest of the request should just work then.

## Design



## Some easy things to do next

* Automate schema creation
* Re-enable tests by creating a DB handler interface/in-memory
* Add request logging and monitoring middleware
* Deduplicate lots of shared code in requests
