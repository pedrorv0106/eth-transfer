# eth-transfer
Ethereum Transfer

## Get all dependencies

    go get ./...

## DB Migrations
#### Install migration tool goose

#    go get -v goose

### Run migration script

#	goose -dir migrations mysql "$DBUSER:$DBPASSWORD@/xentrade_development" up

## copy .env.mainnet to .env
  cp .env.mainnet .env

## change environment variables in .env  

## Run Application
    go run blockchain/main.go
    go run fee_transaction/main.go
