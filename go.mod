module github.com/marekaf/gcr-lifecycle-policy

go 1.14

require (
	github.com/go-openapi/strfmt v0.19.3 // indirect
	github.com/google/go-cmp v0.4.0 // indirect
	github.com/imdario/mergo v0.3.10 // indirect
	github.com/jedib0t/go-pretty v4.3.0+incompatible
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.5.1 // indirect
	go.mongodb.org/mongo-driver v1.1.2 // indirect
	golang.org/x/net v0.0.0-20200226121028-0de0cce0169b // indirect
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	golang.org/x/sys v0.0.0-20200122134326-e047566fdf82 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	k8s.io/api v0.0.0-20190918155943-95b840bb6a1f
	k8s.io/apimachinery v0.0.0-20190913080033-27d36303b655
	k8s.io/client-go v0.0.0-20190918160344-1fbdaa4c8d90
	k8s.io/utils v0.0.0-20200731180307-f00132d28269 // indirect
)

replace github.com/marekaf/gcr-lifecycle-policy/pkg/worker v0.0.0 => ../pkg/worker
