[![Travis](https://img.shields.io/travis/beevelop/docker-taiga-back.svg?style=flat-square)](https://travis-ci.org/beevelop/docker-taiga-back)
[![Docker Pulls](https://img.shields.io/docker/pulls/beevelop/taiga-back.svg?style=flat-square)](https://links.beevelop.com/d-taiga-back)
[![ImageLayer](https://badge.imagelayers.io/beevelop/taiga-back:latest.svg)](https://imagelayers.io/?images=beevelop/taiga-back:latest)
[![Beevelop](https://links.beevelop.com/honey-badge)](https://beevelop.com)

# beevelop/taiga-back (adapted from htdvisser/taiga-back)
> [Taiga](https://taiga.io/) is a project management platform for startups and agile developers & designers who want a simple, beautiful tool that makes work truly enjoyable.

This Docker image can be used for running the Taiga backend. It works together with the [beevelop/taiga-front](https://registry.hub.docker.com/u/beevelop/taiga-front/) image.

## Running
A [postgres](https://registry.hub.docker.com/_/postgres/) container should be linked to the taiga-back container. The taiga-back container will use the ``POSTGRES_USER`` and ``POSTGRES_PASSWORD`` environment variables that are supplied to the postgres container.

```
docker run --name taiga_back_container --link postgres_container:postgres beevelop/taiga-back
```

## [Docker-Compose](https://gist.github.com/beevelop/1975674c22ce8948c895#file-docker-compose-yml)

## Environment

* ``SECRET_KEY`` defaults to ``"insecurekey"``, but you might want to change this.
* ``PUBLIC_REGISTER_ENABLED`` defaults to ``True``

* `FEEDBACK_ENABLED`: enables Feedback functionality (default: True)
* `FEEDBACK_EMAIL`: (default *support@taiga.io*)

* `STATS_ENABLED`: enable / disable stats (default: False)

* `FRONT_SITEMAP_ENABLED`: enable / disable sitemap (default: False)

* `GITHUB_API_CLIENT_ID`: enable GitHub signin/login by providing a client ID (see [register a new application](https://github.com/settings/applications/new))
* `GITHUB_API_CLIENT_SECRET`: the respective client secret

### URLs for static files and media files from taiga-back:
* ``MEDIA_URL`` defaults to ``"$API_SCHEME://$HOSTNAME/media/"``
* ``STATIC_URL`` defaults to ``"$API_SCHEME://$HOSTNAME/media/"``

### Domain configuration
* ``API_SCHEME`` defaults to ``"http"``. Use ``https`` if ``htdvisser/taiga-front-dist`` is used and SSL enabled.
* ``API_DOMAIN`` defaults to ``"$HOSTNAME"``
* ``FRONT_SCHEME`` defaults to ``"http"``. Use ``https`` if ``htdvisser/taiga-front-dist`` is used and SSL enabled.
* ``FRONT_DOMAIN`` defaults to ``"$HOSTNAME"``

### Database
If you want to use an existing PostgreSQL database, you can specify the database connection with the following environment variables
* ``DB_HOST``
* ``DB_USER``
* ``DB_NAME``
* ``DB_PASS``
* ``DB_PORT`` (default: 5432)

### Email configuration
* ``EMAIL_USE_TLS`` defaults to ``False``
* ``EMAIL_HOST`` defaults to ``"localhost"``
* ``EMAIL_PORT`` defaults to ``"25"``
* ``EMAIL_HOST_USER`` defaults to ``""``
* ``EMAIL_HOST_PASSWORD`` defaults to ``""``
* ``DEFAULT_FROM_EMAIL`` defaults to ``"no-reply@example.com"``

### Sample data
* ``SAMPLE_DATA``: set to any value (e.g. True) to initialize the Database with sample data

### Debugging
* ``DEBUG_ENTRYPOINT``: set to True to enable Debugging (`set -x`)
* ``DEBUG`` defaults to ``False``
* ``TEMPLATE_DEBUG`` defaults to ``False``
