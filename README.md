[![Build Status](https://drone.monkiato.com/api/badges/monkiato/guardian/status.svg?ref=refs/heads/master)](https://drone.monkiato.com/monkiato/guardian)
[![codecov](https://codecov.io/gh/monkiato/guardian/branch/master/graph/badge.svg?token=JA4L3YMUVP)](https://codecov.io/gh/monkiato/guardian)



# Guardian

Authentication server in Go created mainly for forward authentication

## Features

 - Use external Postgres DB for user data managing
 - Forward headers
 
## Endpoints

**POST /auth/signin**   send user data for registration

**POST /auth/login**    login with user and password to generate the expected cookie

**POST /auth/approve**  approve existing user registration using the approval token provided during the signin

**GET  /auth/validate** validate current cookie and return status 200 if everything is ok
 
 # Environment Varaiables

 **SECRET_KEY** (mandatory) code used for token encoding and decoding
 
 **COOKIE_NAME** (optional) key value used to store cookie in browsers, it's recommended to declare this value to prevent conflicts

 **TOKEN_EXPIRATION_HOURS** (default '24' hours) specify token duration in hours, only integer numbers are allowed
 
 **REDIRECT_URL** (default 'http://localhost/') url used to redirect the origin request when the authentication validation has failed (generally it's a redirect to the login page)
 
 **DOMAIN_NAME** (default 'localhost') required to make cookies available across any services under the same domain

## Build Docker Image

No extra parameters are required for the docker image, so just run:

`docker build . -t guardian-auth`

## Run Docker Container

TBD

