#!/bin/sh

set -e

CURL="curl -s --connect-timeout 3 --max-time 10 -o /tmp/response_body.txt -w %{http_code}"
HOSTNAME="${HOSTNAME:-api}"
API="http://${HOSTNAME}:${API_PORT}/api/v1"

request() {
	name="$1"
	shift
	code=$(${CURL} "$@")
	echo "${name}:${code}"
}

check_jq_true() {
	name="$1"
	expr="$2"
	if jq -e "$expr" /tmp/response_body.txt >/dev/null 2>&1; then
		echo "${name}:true"
	else
		echo "${name}:false"
	fi
}

check_body_non_empty() {
	name="$1"
	if [ -s /tmp/response_body.txt ]; then
		echo "${name}:true"
	else
		echo "${name}:false"
	fi
}

check_error_shape() {
	name="$1"
	expected_code="$2"
	expected_message="$3"
	if jq -e ".code == ${expected_code} and .message == \"${expected_message}\"" /tmp/response_body.txt >/dev/null 2>&1; then
		echo "${name}:true"
	else
		echo "${name}:false"
	fi
}

request health http://${HOSTNAME}:${API_PORT}/health

request users_empty ${API}/users
request tweets_empty ${API}/tweets

request users_invalid_uuid ${API}/users/not-a-uuid
check_body_non_empty users_invalid_uuid_payload
request tweets_invalid_uuid ${API}/tweets/not-a-uuid
check_body_non_empty tweets_invalid_uuid_payload

request users_not_found ${API}/users/00000000-0000-0000-0000-000000000000
check_error_shape users_not_found_payload 404 "User not found"
request tweets_not_found ${API}/tweets/00000000-0000-0000-0000-000000000000
check_error_shape tweets_not_found_payload 404 "Tweet not found"

request users_bad_payload -X POST -H "Content-Type: application/json" \
	-d '{ "username": "foo" }' \
	${API}/users
check_body_non_empty users_bad_payload_payload

request users_create -X POST -H "Content-Type: application/json" \
	-d '{ "username": "foo", "name": "John Doe", "email": "jd@mail.com" }' \
	${API}/users

request users_duplicate_username -X POST -H "Content-Type: application/json" \
	-d '{ "username": "foo", "name": "John Doe 2", "email": "jd2@mail.com" }' \
	${API}/users
check_error_shape users_duplicate_username_payload 500 "Error creating user"

request users_duplicate_email -X POST -H "Content-Type: application/json" \
	-d '{ "username": "foo2", "name": "John Doe 3", "email": "jd@mail.com" }' \
	${API}/users
check_error_shape users_duplicate_email_payload 500 "Error creating user"

request users_list ${API}/users
check_jq_true users_list_payload 'length == 1 and .[0].username == "foo" and .[0].email == "jd@mail.com" and .[0].name == "John Doe" and (. [0].id | type == "string" and length > 0)'
user_id=$(jq -r '.[0].id // empty' /tmp/response_body.txt)

if [ -z "$user_id" ]; then
	echo "error:user_id_missing"
	exit 1
fi

request users_get ${API}/users/${user_id}
check_jq_true users_get_payload '.username == "foo" and .email == "jd@mail.com" and .name == "John Doe" and (.id | type == "string" and length > 0)'

request tweets_bad_payload -X POST -H "Content-Type: application/json" \
	-d '{"user_id":"'$user_id'"}' \
	${API}/tweets
check_body_non_empty tweets_bad_payload_payload

request tweets_invalid_user -X POST -H "Content-Type: application/json" \
	-d '{"user_id":"00000000-0000-0000-0000-000000000000", "content": "Hello World!" }' \
	${API}/tweets
check_error_shape tweets_invalid_user_payload 400 "Invalid user ID"

tweet_280=$(printf '%*s' 280 '' | tr ' ' 'a')
tweet_281=$(printf '%*s' 281 '' | tr ' ' 'b')

request tweets_len_280 -X POST -H "Content-Type: application/json" \
	-d '{"user_id":"'$user_id'", "content": "'$tweet_280'" }' \
	${API}/tweets

request tweets_len_281 -X POST -H "Content-Type: application/json" \
	-d '{"user_id":"'$user_id'", "content": "'$tweet_281'" }' \
	${API}/tweets
check_error_shape tweets_len_281_payload 500 "Error creating tweet"

request tweets_create -X POST -H "Content-Type: application/json" \
	-d '{"user_id":"'$user_id'", "content": "Hello World!" }' \
	${API}/tweets

request tweets_list ${API}/tweets
check_jq_true tweets_list_payload 'length >= 2 and (map(.content) | index("Hello World!")) != null and (map(.content) | index("'$tweet_280'")) != null and all(.[]; (.id | type == "string" and length > 0) and (.user_id | type == "string" and length > 0))'
tweet_id=$(jq -r '.[] | select(.content == "Hello World!") | .id' /tmp/response_body.txt | head -n 1)

if [ -z "$tweet_id" ]; then
	echo "error:tweet_id_missing"
	exit 1
fi

request tweets_get ${API}/tweets/${tweet_id}
check_jq_true tweets_get_payload '.content == "Hello World!" and (.id | type == "string" and length > 0) and (.user_id | type == "string" and length > 0)'
