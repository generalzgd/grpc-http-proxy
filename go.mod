module grpcproxy

go 1.12

require (
	github.com/astaxie/beego v1.12.0
	github.com/generalzgd/comm-libs v0.0.0-20200109081424-f5b2eb752230
	github.com/generalzgd/grpc-svr-frame v0.0.0-20190929085824-d6a34f7d1057
	github.com/golang/protobuf v1.3.2
	github.com/grpc-ecosystem/grpc-gateway v1.12.1
	github.com/processout/grpc-go-pool v1.2.1
	google.golang.org/grpc v1.24.0
	hutte.zhanqi.tv/go/grpc-frame v0.0.0-20200108063655-fd79acc586bc
	hutte.zhanqi.tv/go/link v0.0.0-20190806080448-2db9d9ddb37d
	hutte.zhanqi.tv/go/svr-frame v0.0.0-20200108075438-9746e890da79
)

replace (
	cloud.google.com/go => github.com/googleapis/google-cloud-go v0.37.4
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190513172903-22d7a77e9e5f
	golang.org/x/exp => github.com/golang/exp v0.0.0-20190718202018-cfdd5522f6f6
	golang.org/x/image => github.com/golang/image v0.0.0-20190703141733-d6a02ce849c9
	golang.org/x/lint => github.com/golang/lint v0.0.0-20190409202823-959b441ac422
	golang.org/x/mobile => github.com/golang/mobile v0.0.0-20190719004257-d2bd2a29d028
	golang.org/x/mod => github.com/golang/mod v0.1.0
	golang.org/x/net => github.com/golang/net v0.0.0-20190827160401-ba9fcec4b297
	golang.org/x/oauth2 => github.com/golang/oauth2 v0.0.0-20190604053449-0f29369cfe45
	golang.org/x/sync => github.com/golang/sync v0.0.0-20190423024810-112230192c58
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190712062909-fae7ac547cb7
	golang.org/x/text => github.com/golang/text v0.3.2
	golang.org/x/time => github.com/golang/time v0.0.0-20190308202827-9d24e82272b4
	golang.org/x/tools => github.com/golang/tools v0.0.0-20191217033636-bbbf87ae2631
	golang.org/x/xerrors => github.com/golang/xerrors v0.0.0-20191204190536-9bdfabe68543
	google.golang.org/api v0.3.1 => github.com/googleapis/google-api-go-client v0.3.1
	google.golang.org/appengine => github.com/golang/appengine v1.6.1
	google.golang.org/genproto => github.com/googleapis/go-genproto v0.0.0-20190516172635-bb713bdc0e52
	google.golang.org/grpc => github.com/grpc/grpc-go v1.24.0
)
