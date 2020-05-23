# do-api-proxy

Proxies DigitalOcean API calls restricting the allowed resources using a custom API token.

[DigitalOcean's API](https://developers.digitalocean.com/documentation/) offers a great number of automation possibilities. Unfortunately, as of today it is only possible to access the API via a unrestricted personal access token. This project tries to solve that problem.

*do-api-proxy* can be used as a proxy for DigitalOcean API clients, where you set up a custom API token which only allows restricted API access.

Currently, only access to the **volume** API is allowed: *All* volumes can be listed but only the configured volumes can be attached/detached or modified.

## Usage

### With Docker

Copy *.env.sample* to *.env* and adjust the values:

* `API_TOKEN`: Random token that clients will need to use in order to access the proxy
* `TARGET_API_TOKEN`: Personal access token from your DO account
* `VOLUMES`: Restrict volume access to given volumes (comma-separated list).

```bash
docker run -d -p 1338:1338 --env-file .env mazzolino/do-api-proxy
```

Alternatively, use the supplied [docker-compose.yml](docker-compose.yml).


### Without Docker

Download the latest version for your OS from the [releases page](https://github.com/djmaze/do-api-proxy/releases) and make it executable (`chmod u+x do-api-proxy`).

Run it like this:

```bash
env API_TOKEN=<API_TOKEN> TARGET_API_TOKEN=<TARGET_API_TOKEN> VOLUMES=<VOLUMES> ./do-api-proxy
```

### Proxy usage

When started, the proxy is available at port 1338. You can e.g. direct your local `doctl` command to it like this:


```bash
doctl -t <TARGET_API_TOKEN> -u http://localhost:1338 compute volume list
```

You can now make it available via https using a reverse proxy like [Traefik](https://containo.us/traefik/).

## Building

### With Docker

You can build the docker image using Docker Compose:

```bash
docker-compose build
```

### Without Docker

Install golang 1.13 and run the following:

```bash
go build
```

The binary should now be available at *./do-api-proxy*.
