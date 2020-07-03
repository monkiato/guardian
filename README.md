<img src="docs/logo.svg" width="30%"/>

# Guardian
[![Build Status](https://drone.monkiato.com/api/badges/monkiato/guardian/status.svg?ref=refs/heads/master)](https://drone.monkiato.com/monkiato/guardian)
[![codecov](https://codecov.io/gh/monkiato/guardian/branch/master/graph/badge.svg)](https://codecov.io/gh/monkiato/guardian)
[![Go Report Card](https://goreportcard.com/badge/github.com/monkiato/guardian)](https://goreportcard.com/report/github.com/monkiato/guardian)


Authentication server in Go created mainly for forward authentication

## Features

 - Session ID data stored in cookies
 - Using secure configuration for cookies to prevent CSRF or any other attacks
 - Use external Postgres DB for user data managing
 - Forward headers
 
## Endpoints

| Type | Path | Description |
| ---- | ---- | ----------- |
| POST | /auth/signin | send user data for registration |
| POST | /auth/login | login with user and password to generate the expected cookie |
| POST | /auth/logout | logout existing logged-in user (invalidates cookie) |
| POST | /auth/approve | approve existing user registration using the approval token provided during the signin |
| POST | /auth/disable | disable existing user, approval token is required to invalidate the user |
| GET | /auth/validate | validate current cookie and return status 200 if everything is ok |
| GET | /auth/me | get logged-in user data, otherwise 403 is returned |
 
 # Environment Varaiables

 **SECRET_KEY** (mandatory) code used for token encoding and decoding
 
 **COOKIE_NAME** (optional) key value used to store cookie in browsers, it's recommended to declare this value to prevent conflicts

 **TOKEN_EXPIRATION_HOURS** (default '24' hours) specify token duration in hours, only integer numbers are allowed
 
 **REDIRECT_URL** (default 'http://localhost/') url used to redirect the origin request when the authentication validation has failed (generally it's a redirect to the login page)
 
 **DOMAIN_NAME** (default 'localhost') required to make cookies available across any services under the same domain

 # Secure Configuration for Cookies
 
  - Cookies attached to a specific domain or subdomain only
  - HTTPS support only through 'Secure' property
  - Mitigate XSS attack using 'HttpOnly' property
 
## Build Docker Image

No extra parameters are required for the docker image, so just run:

`docker build . -t guardian-auth`

## Run Docker Container

`docker-compose.yml` is available as a sample to deploy a stack
 with Guardian and a Postgres DB
 
 `docker-compose.traefik.yml` is also available with Traefik used as a reverse proxy
 
 Run `docker-compose -f docker-compose.yml up -d`

# 'Guardian' Logo

<div>Icons made by <a href="https://www.flaticon.com/authors/freepik" title="Freepik">Freepik</a> from <a href="https://www.flaticon.com/" title="Flaticon">www.flaticon.com</a></div>
