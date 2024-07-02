module github.com/spin_s3_example

go 1.22.3

toolchain go1.22.4

require (
	github.com/fermyon/spin-go-sdk v0.0.0-20240220234050-48ddef7a2617
	github.com/fermyon/spin-aws-go v0.0.0
)

require github.com/julienschmidt/httprouter v1.3.0 // indirect

replace github.com/fermyon/spin-aws-go => ../.
