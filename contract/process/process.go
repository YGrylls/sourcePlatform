package process

import (
	"encoding/json"
)
import perror "github.com/pkg/errors"

//-----------State of process
type State int

const (
	InProcess State = iota
	Completed
)

var stateNames = []string{"InProcess", "Completed"}

func (s State) String() string {
	if s < InProcess || s > Completed {
		return "Unknown"
	}
	return stateNames[s]
}

const ClassName = "org.sourceplatform.process"

//-----------

//-----------Process definition
type Process struct {
	ProcessLocalId   string `json:"processLocalId"`
	OwnerOrg         string `json:"ownerOrg"`
	OptionName       string `json:"optionName"`
	StartTime        int64  `json:"startTime"`
	CompleteTime     int64  `json:"completeTime"`
	StartPosition    string `json:"startPosition"`
	CompletePosition string `json:"completePosition"`
	PreKey           string `json:"preKey"`
	State            State  `json:"state"`
	Class            string `json:"class"`
}

func (p *Process) Serialize() ([]byte, error) {
	return json.Marshal(p)
}

func Deserialize(bytes []byte, p *Process) error {
	err := json.Unmarshal(bytes, p)
	if err != nil {
		return perror.Errorf("Error deserializing process json bytes. %s", err)
	}
	return nil
}
