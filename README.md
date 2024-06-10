# S3 client for Spin applications

Mostly working...

## TODO

- AWS request signing is broken for getting objects but not putting or listing
- Signing needs some refactoring
- Add relevant parameters for making requests. (list filters, etc)
- Ensure correct parsing for response objects
- Remove httputil tracing due to tinygo breaking net/http
- Support paging results
- Move internals to `internal` packages to allow unit tests to work
