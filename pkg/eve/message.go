package eve

import (
	"context"
	"encoding/json"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/unanet/eve/pkg/s3"
	"github.com/unanet/go/pkg/errors"
)

type CloudUploader interface {
	Upload(ctx context.Context, key string, body []byte) (*s3.Location, error)
}

type CloudDownloader interface {
	Download(ctx context.Context, location *s3.Location) ([]byte, error)
}

type DeploymentCallbackMessage struct {
	DeploymentID uuid.UUID            `json:"deployment_id"`
	Status       DeploymentPlanStatus `json:"status"`
	Type         PlanType             `json:"type"`
	Messages     []string             `json:"messages"`
}

type CallbackMessage struct {
	Messages []string `json:"messages"`
}

func UnMarshalNSDeploymentFromS3LocationBody(ctx context.Context, cd CloudDownloader, b []byte) (*NSDeploymentPlan, error) {
	var location s3.Location
	err := json.Unmarshal(b, &location)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	planText, err := cd.Download(ctx, &location)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	var nsDeploymentPlan NSDeploymentPlan
	err = json.Unmarshal(planText, &nsDeploymentPlan)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return &nsDeploymentPlan, nil
}

func MarshalNSDeploymentPlanToS3LocationBody(ctx context.Context, cu CloudUploader, plan *NSDeploymentPlan) ([]byte, error) {
	nsDeploymentJson, err := json.Marshal(plan)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	location, err := cu.Upload(ctx, fmt.Sprintf("%s.json", plan.DeploymentID), nsDeploymentJson)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	locationJson, err := json.Marshal(location)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return locationJson, nil
}
