module github.com/spin_s3_example

go 1.22.3

require (
	github.com/fermyon/spin-s3-go v0.0.0
	github.com/fermyon/spin/sdk/go/v2 v2.2.0
)

require github.com/julienschmidt/httprouter v1.3.0 // indirect

replace github.com/fermyon/spin-s3-go => ../.
