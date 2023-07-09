#!/bin/sh

CURL='curl -s -w %{http_code}\n'

## Get from an empty DB
${CURL} http://${HOSTNAME}:${API_PORT}/api/v1/users
${CURL} http://${HOSTNAME}:${API_PORT}/api/v1/users/00000000-0000-0000-0000-000000000000
${CURL} http://${HOSTNAME}:${API_PORT}/api/v1/tweets
${CURL} http://${HOSTNAME}:${API_PORT}/api/v1/tweets/00000000-0000-0000-0000-000000000000

## Invalid UUIDs format
${CURL} http://${HOSTNAME}:${API_PORT}/api/v1/users/id1
${CURL} http://${HOSTNAME}:${API_PORT}/api/v1/tweets/123456

${CURL} -X POST -H "Content-Type: application/json" \
-d '{ "username": "foo", "name": "John Doe", "email": "jd@mail.com" }' \
http://${HOSTNAME}:${API_PORT}/api/v1/users

## Get all users
${CURL} http://${HOSTNAME}:${API_PORT}/api/v1/users

## Set user_id as an environment variable
user_id=$(${CURL} http://${HOSTNAME}:${API_PORT}/api/v1/users | jq -r -R 'fromjson? | .[0].id')

## Get a user by id
${CURL} http://${HOSTNAME}:${API_PORT}/api/v1/users/$user_id

## Create a tweet
${CURL} -X POST -H "Content-Type: application/json" \
-d '{"user_id":"'$user_id'", "content": "Hello World!" }' \
http://${HOSTNAME}:${API_PORT}/api/v1/tweets

# Get all tweets
${CURL} http://${HOSTNAME}:${API_PORT}/api/v1/tweets

## Set twitter_id as an environment variable
tweet_id=$(${CURL} http://${HOSTNAME}:${API_PORT}/api/v1/tweets | jq -r -R 'fromjson? | .[0].id')

## Get a tweet by id
${CURL} http://${HOSTNAME}:${API_PORT}/api/v1/tweets/$tweet_id

## Create a tweeter with an non-existing user_id
${CURL} -X POST -H "Content-Type: application/json" \
-d '{"user_id":"00000000-0000-0000-0000-000000000000", "content": "Hello World!" }' \
http://${HOSTNAME}:${API_PORT}/api/v1/tweets
