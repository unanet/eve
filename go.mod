module gitlab.unanet.io/devops/eve

go 1.16

//replace gitlab.unanet.io/devops/go => ../../devops/go

require (
	github.com/aws/aws-sdk-go v1.37.25
	github.com/dghubble/sling v1.3.0
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-chi/cors v1.1.1
	github.com/go-chi/render v1.0.1
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-migrate/migrate/v4 v4.14.1
	github.com/google/go-cmp v0.5.2 // indirect
	github.com/jmoiron/sqlx v1.3.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/kr/pretty v0.2.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/robfig/cron/v3 v3.0.1
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.6.1
	gitlab.unanet.io/devops/go v1.4.0
	go.uber.org/zap v1.16.0
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110 // indirect
	golang.org/x/text v0.3.5 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
