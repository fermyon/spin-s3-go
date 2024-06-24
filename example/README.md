# Spin S3 Example

This example uses [LocalStack](https://github.com/localstack/localstack) rather
than requiring access to AWS. LocalStack is a cloud software development
framework to develop and test your AWS applications locally.

Start a localstack service using Docker

```
docker run \
  --rm -it \
  -p 4566:4566 \
  -p 4510-4559:4510-4559 \
  localstack/localstack
```

Build and start the Spin application

```
SPIN_VARIABLE_S3_ENDPOINT=http://s3.localhost.localstack.cloud:4566 spin build --up
```
Test the application

```
curl localhost:3000
```
