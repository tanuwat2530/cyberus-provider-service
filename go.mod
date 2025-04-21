module cyberus/provider-service

go 1.22.2

replace CyberusGolangShareLibrary => ../cyberus-common-library

require (
	CyberusGolangShareLibrary v1.2.0
	github.com/google/uuid v1.6.0
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/redis/go-redis/v9 v9.7.3 // indirect
)
