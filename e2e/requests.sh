#!/bin/sh

set -e

CURL="curl -s"

${CURL} -X POST -H "Content-Type: application/json" \
-d '{ "username": "foo", "name": "John Doe", "email": "jd@mail.com" }' \
http://api:${API_PORT}/api/v1/users | jq .

## Get all users
${CURL} http://api:${API_PORT}/api/v1/users  | jq .

## Set user_id as an environment variable
user_id=$(${CURL} http://api:${API_PORT}/api/v1/users | jq -r '.[0].id')

## Get a user by id
${CURL} http://api:${API_PORT}/api/v1/users/$user_id  | jq .

## Create a tweet
${CURL} -X POST -H "Content-Type: application/json" \
-d '{"user_id":"'$user_id'", "content": "Hello World!" }' \
http://api:${API_PORT}/api/v1/tweets  | jq .

# Get all tweets
${CURL} http://api:${API_PORT}/api/v1/tweets  | jq .

## Set twitter_id as an environment variable
tweet_id=$(${CURL} http://api:${API_PORT}/api/v1/tweets | jq -r '.[0].id')

## Get a tweet by id
${CURL} http://api:${API_PORT}/api/v1/tweets/$tweet_id  | jq .
