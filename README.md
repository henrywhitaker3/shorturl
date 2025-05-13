# Shorturl

A simple shorturl service backed by postgres and redis

## Architecture

![Architecture Diagram](./diagrams/assets/architecture.png)

### App Servers

The app server/api creates and serves shorturls.

When a request comes in to create a url, a  job is queued in redis to create the
underlying shorturl.

When a shorturl is visited, a local in-memory LRU cache is first checked for
the url, if it doesn't exist, it is stored in the cache then the client is
redirected to the long url - this reduces reads forwarded to the database.

### Generator

A background process runs that generates aliases (the shorturl id). This way,
creating a new url only has to reserve an alias in the aliases table instead
of generating a new one itself and checking if it is available.

The number of free aliases available can be modified in the config, but defaults
to 100,000:

```yaml
generator:
    buffer_szie: 250000
```

### Click Tracking

If click tracking is turned on:

```yaml
tracking:
    enabled: true
```

Then each visit to a shorturl will queue a job to store the click. This will store:

- The shorturl
- The IP address
- The time the url was visited

Visit retention can be configured via config file:

```yaml
tracking:
    retention:
        # If disabled, history will be kept indefinitely
        enabled: true
        period: 48h
```
