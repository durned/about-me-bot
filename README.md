# 2.5 Subscription Bot 

### Prerequisites 
- latest `GoLang` version
- docker desktop
- `makefile`
- `.env` filled out

## Control 
`/about-me-bot/cmd/.env`: two operation modes possible, them being server and worker

## Make commands 
- `make local`: runs db in a docker container and go code locally 
- `make deployed`: runs both db and go code in two docker containers
- `make stop`: stops the container operation
- `make clean`: does what stop does but removes saved container data, too

