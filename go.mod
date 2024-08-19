module github.com/oomph-ac/mv

go 1.22

toolchain go1.22.2

require (
	github.com/df-mc/dragonfly v0.9.17
	github.com/df-mc/worldupgrader v1.0.16
	github.com/go-gl/mathgl v1.1.0
	github.com/google/uuid v1.6.0
	github.com/pelletier/go-toml v1.9.5
	github.com/sandertv/gophertunnel v1.40.1
	github.com/sirupsen/logrus v1.9.3
	golang.org/x/exp v0.0.0-20240416160154-fe59bbe5cc7f
)

require (
	github.com/brentp/intintmap v0.0.0-20190211203843-30dc0ade9af9 // indirect
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/df-mc/atomic v1.10.0 // indirect
	github.com/df-mc/goleveldb v1.1.9 // indirect
	github.com/go-jose/go-jose/v3 v3.0.3 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/muhammadmuzzammil1998/jsonc v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/sandertv/go-raknet v1.14.1 // indirect
	github.com/segmentio/fasthash v1.0.3 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/crypto v0.24.0 // indirect
	golang.org/x/image v0.17.0 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/oauth2 v0.21.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
)

// replace github.com/sandertv/gophertunnel v1.37.0 => github.com/oomph-ac/gophertunnel v0.0.0-20240711192558-cd94bd1ba4b7

replace github.com/sandertv/go-raknet v1.13.0 => github.com/tedacmc/tedac-raknet v0.0.4
