module github.com/c12s/star

go 1.21

toolchain go1.21.0

require (
	github.com/c12s/kuiper v1.0.0
	github.com/c12s/magnetar v1.0.0
	github.com/c12s/oort v1.0.0
	github.com/nats-io/nats.go v1.28.0
	github.com/shirou/gopsutil v3.21.11+incompatible
)

require (
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/nats-io/nkeys v0.4.4 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	github.com/tklauser/go-sysconf v0.3.11 // indirect
	github.com/tklauser/numcpus v0.6.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.3 // indirect
	golang.org/x/crypto v0.11.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/text v0.11.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230525234030-28d5490b6b19 // indirect
	google.golang.org/grpc v1.57.0
	google.golang.org/protobuf v1.31.0
)

replace github.com/c12s/magnetar => ../magnetar

replace github.com/c12s/kuiper => ../kuiper

replace github.com/c12s/oort => ../oort
