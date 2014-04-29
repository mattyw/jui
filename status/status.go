package status

import (
	"fmt"
	"launchpad.net/goyaml"
	"os/exec"
	"strings"
)

type SubordinateCharm struct {
	AgentState string `agent-state`
}

type CharmData struct {
	AgentState    string `agent-state`
	PublicAddress string `public-address`
	Subordinates  map[string]SubordinateCharm
}

type ServiceUnits map[string]CharmData

type Service struct {
	Charm       string
	Exposed     bool
	Relations   map[string][]string
	Subordinate bool
	Units       ServiceUnits
}

func (s *Service) CharmVersion() string {
	charmName := strings.Split(s.Charm, "-")
	return charmName[len(charmName)-1:][0]
}

type JujuStatus struct {
	Services map[string]Service
}

type DeployerService struct {
	Branch string
}

type JujuDeployer struct {
	Services  map[string]DeployerService
	Inherits  string
	Relations [][]string
}

type Deployer map[string]JujuDeployer

//Call "juju status"
func CallJujuStatus() ([]byte, error) {
	out, err := exec.Command("juju", "status").Output()
	if err != nil {
		return []byte{0}, fmt.Errorf("Error calling juju status with %s", err)
	}

	return out, nil
}

func StatusFromApi(status []byte) (JujuStatus, error) {
	jujuState := JujuStatus{}
	err := goyaml.Unmarshal(status, &jujuState)
	if err != nil {
		return jujuState, err
	}
	return jujuState, nil
}

func StatusFromDeployer(name string, status []byte) (JujuDeployer, error) {
	jujuState := Deployer{}
	err := goyaml.Unmarshal(status, &jujuState)
	if err != nil {
		return jujuState[name], err
	}
	// Expand inherits
	fmt.Println(jujuState[name])
	inheritedName := jujuState[name].Inherits
	for sName, service := range jujuState[inheritedName].Services {
		jujuState[name].Services[sName] = service
	}
	return jujuState[name], nil
}

func GetStatus() (JujuStatus, error) {
	output, err := CallJujuStatus()
	if err != nil {
		return JujuStatus{}, err
	}
	status, err := StatusFromApi(output)
	if err != nil {
		return JujuStatus{}, err
	}
	return status, nil
}
