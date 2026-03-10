# configrefresh

```bash
./build.sh            ;# build

./run-local-rabbit.sh ;# run rabbit

export AMQP_URL=amqp://guest:guest@localhost:5672/

configrefresh         ;# run the tool

## open http://localhost:15672/ -- login with guest:guest"
```

# Docker images

https://hub.docker.com/r/udhos/configrefresh

```bash
docker run --rm udhos/configrefresh:0.0.1
```