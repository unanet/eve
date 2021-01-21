module gitlab.unanet.io/devops/eve

go 1.15

//replace gitlab.unanet.io/devops/go => ../../devops/go

//replace (
//	github.com/Sirupsen/logrus v1.7.0 => github.com/sirupsen/logrus v1.7.0
//	github.com/docker/docker v1.13.2-0.20170601211448-f5ec1e2936dc => github.com/docker/engine v17.12.0-ce-rc1.0.20180718150940-a3ef7e9a9bda+incompatible
//)

require (
	github.com/aws/aws-sdk-go v1.27.0
	github.com/dghubble/sling v1.3.0
	github.com/go-chi/chi v4.1.0+incompatible
	github.com/go-chi/cors v1.1.1
	github.com/go-chi/render v1.0.1
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0
	github.com/golang-migrate/migrate/v4 v4.10.0
	github.com/jmoiron/sqlx v1.2.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/lib/pq v1.3.0
	github.com/pkg/errors v0.9.1
	github.com/robfig/cron/v3 v3.0.0
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.4.0
	gitlab.unanet.io/devops/go v1.0.10
	go.uber.org/zap v1.16.0
)
