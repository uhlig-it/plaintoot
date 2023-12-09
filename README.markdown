# plaintoot

Provides a plaintext representation of a Mastodon post ("toot").

# Print

`plaintoot print` prints a plain-text representation of a single post. Provide a post's URL as an argument.

# Run

Grab the [latest release](https://github.com/uhlig-it/plaintoot/releases/latest), unpack it and run `plaintoot` for more information.

There is also a Docker image; run it with:

```command
$ docker run --interactive --tty --rm suhlig/plaintoot
```

# Serve

`plaintoot serve` serves a plain-text representation of a single post via HTTP. It will start an HTTP server listening on `$PORT` (defaults to `8080`) and provide the same functionality as `print` (see above).

# Develop

Install [watchexec](https://github.com/watchexec/watchexec#install) and start the server like this in order to reload the server process on any file changes:

```command
$ STARTUP_DELAY=10s MAX_UPTIME=20s watchexec --restart go run . serve
```

The docker image can be built and pushed with

```command
$ docker buildx build \
  --push \
  --platform linux/arm/v7,linux/arm64/v8,linux/amd64 \
  --tag suhlig/plaintoot
```

# Bonus

This application is intended for use in my course "[Web Services](https://ws.uhlig.it/)" at [DHBW](https://www.ravensburg.dhbw.de/studienangebot/bachelor-studiengaenge/informatik). For this purpose, it has some endpoints that showcase Kubernetes' [liveness and readiness probes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/):

* `/liveness` will generally return `200`, unless the environment variable `MAX_UPTIME` was set to a value [`time.ParseDuration`](https://pkg.go.dev/time#ParseDuration) accepts and the given time since server start has elapsed. Other paths will still work, but `/liveness` will return `500` thereafter.
* If `STARTUP_DELAY` is set, the application will sleep for the given duration (see above for the expected format), and only after that the HTTP port will be opened.
* `/readiness` will return `200` if the authentication with Twitter is successful, otherwise `500` will be returned.

# TODO

* Add /metrics exposing
  - the number of requests since startup
  - uptime
