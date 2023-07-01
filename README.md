# twitter-clone


## Generate Go code from the Open-API spec


## Docker

### build the docker image

```bash
docker build --platform linux/amd64 .

```

## Endpoints

```sh
# create a user
curl -v -X POST -H "Content-Type: application/json" \
-d '{ "username": "foo", "name": "John Doe", "email": "jd@mail.com" }' \
http://localhost:8888/users

# get all users
curl -v -X GET http://localhost:8888/users

# Create a tweet
curl -v -X POST -H "Content-Type: application/json" \
-d '{"user_id":"283d3731-2f90-4112-bcf9-c8d117f6b250", "content": "Hello World!" }' \
http://localhost:8888/tweets
```
