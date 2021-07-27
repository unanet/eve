package crud

import (
	"encoding/json"
	"github.com/unanet/eve/internal/data"
	"github.com/unanet/eve/pkg/eve"
	"testing"
)

const (
	defSpecA = "{\"spec\": {\"template\": {\"spec\": {\"nodeSelector\": {\"node-group\": \"shared\"}}}}}"
	defSpecB = "{\"spec\": {\"template\": {\"spec\": {\"containers\": [{\"livenessProbe\": {\"httpGet\": {\"path\": \"/analytics-api/Api.asmx\", \"port\": 8080}, \"periodSeconds\": 10, \"initialDelaySeconds\": 30}}]}}}}"
	defSpecC = "{\"spec\": {\"template\": {\"spec\": {\"containers\": [{\"readinessProbe\": {\"httpGet\": {\"path\": \"/analytics-api/Api.asmx\", \"port\": 8080 }, \"periodSeconds\": 10, \"initialDelaySeconds\": 45}}]}}}}"
	defSpecD = "{\"spec\": {\"template\": {\"spec\": {\"containers\": [{\"readinessProbe\": {\"httpGet\": {\"path\": \"/analytics-api/Api.asmx\", \"port\": 8080 }, \"periodSeconds\": 10, \"initialDelaySeconds\": 45}}]}}}}"
	defSpecE = "{\"spec\": {\"minReplicas\":2, \"maxReplicas\": 10}}"
)

func dummySpecData(spec string) map[string]interface{} {
	jsonMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(spec), &jsonMap)
	if err != nil {
		panic(err)
	}
	return jsonMap
}

func dummyDefSpec() []eve.DefinitionResult {
	var defSpecs = make([]eve.DefinitionResult, 0)

	defSpecs = append(defSpecs, eve.DefinitionResult{Class: "apps", Version: "v1", Kind: "Deployment", Order: "main", Data: dummySpecData(defSpecA)})
	defSpecs = append(defSpecs, eve.DefinitionResult{Class: "apps", Version: "v1", Kind: "Deployment", Order: "main", Data: dummySpecData(defSpecB)})
	defSpecs = append(defSpecs, eve.DefinitionResult{Class: "batch", Version: "v1", Kind: "Job", Order: "main", Data: dummySpecData(defSpecC)})
	defSpecs = append(defSpecs, eve.DefinitionResult{Class: "apps", Version: "v1", Kind: "Deployment", Order: "main", Data: dummySpecData(defSpecD)})
	defSpecs = append(defSpecs, eve.DefinitionResult{Class: "autoscaling", Version: "v1", Kind: "HorizontalPodAutoscaler", Order: "post", Data: dummySpecData(defSpecE)})

	return defSpecs
}

func TestManager_mergeDefinitionData(t *testing.T) {
	type fields struct {
		repo *data.Repo
	}
	type args struct {
		defSpecs []eve.DefinitionResult
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    eve.MetadataField
		wantErr bool
	}{
		{
			name:   "happy",
			fields: fields{repo: nil},
			args:   args{defSpecs: dummyDefSpec()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Manager{
				repo: tt.fields.repo,
			}
			_, err := m.mergeDefinitionData(tt.args.defSpecs)
			if (err != nil) != tt.wantErr {
				t.Errorf("mergeDefinitionData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
