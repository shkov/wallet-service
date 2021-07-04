# wallet-service

This is RESTful API service that provides some abilities of making payments between accounts.

### Installing

Since some docker ["networking features"](https://docs.docker.com/docker-for-mac/networking/#there-is-no-docker0-bridge-on-macos) are not available on macOS, you need to put your ip address in the `POSTGRES_HOST` [variable](https://github.com/shkov/wallet-service/blob/main/deployments/docker-compose.yml#L24) in docker-compose.


Then run `make run` and Compose will start the entire app with all dependencies(postgresql):

### Usage examples:

1) `POST /api/v1/payments` applies the new payment to accounts. Note: if there is no "from" account in the system, it is considered that it has a balance of 1000.

```shell
curl --request POST \
  --url http://127.0.0.1:80/api/v1/payments \
  --header 'Content-Type: application/json' \
  --data '{
	"Amount": "1000",
	"From": 1,
	"To": 2
}'
```

2) `GET /api/v1/accounts/{id}` returns an account by the given id.

```shell
curl --request GET \
  --url http://127.0.0.1:80/api/v1/accounts/1
```

3) `GET /api/v1/payments/{accountId}` returns all payments by the account id.

```shell
curl --request GET \
  --url http://127.0.0.1:80/api/v1/payments/1
```
