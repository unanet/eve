package crud

import (
	"encoding/json"
	"gitlab.unanet.io/devops/eve/internal/data"
	"gitlab.unanet.io/devops/eve/pkg/eve"
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
	//AdefSpec := make(map[string]map[string]interface{})
	err := json.Unmarshal([]byte(spec), &jsonMap)
	if err != nil {
		panic(err)
	}
	//AdefSpec["appsv1.Deployment"] = jsonMapA
	return jsonMap
}

func addDummyDefSpec(t string, d map[string]interface{}) map[string]map[string]interface{} {
	defSpec := make(map[string]map[string]interface{})
	defSpec[t] = d
	return defSpec
}

func dummyDefSpec() []eve.DefinitionSpec {
	var defSpecs = make([]eve.DefinitionSpec, 0)

	defSpecs = append(defSpecs, addDummyDefSpec("appsv1.Deployment", dummySpecData(defSpecA)))
	defSpecs = append(defSpecs, addDummyDefSpec("appsv1.Deployment", dummySpecData(defSpecB)))
	defSpecs = append(defSpecs, addDummyDefSpec("batchv1.Job", dummySpecData(defSpecC)))
	defSpecs = append(defSpecs, addDummyDefSpec("appsv1.Deployment", dummySpecData(defSpecD)))
	defSpecs = append(defSpecs, addDummyDefSpec("v2beta2.HorizontalPodAutoscaler", dummySpecData(defSpecE)))

	return defSpecs
}

func TestManager_mergeDefinitionData(t *testing.T) {
	type fields struct {
		repo *data.Repo
	}
	type args struct {
		defSpecs []eve.DefinitionSpec
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
