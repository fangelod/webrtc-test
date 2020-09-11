package rtc

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	webrtc "github.com/pion/webrtc/v2"
	uuid "github.com/satori/go.uuid"

	"gopkg.in/yaml.v2"
)

const (
	iceServerYamlEnvString = "WEBRTC_TEST_ICE_SERVER_YAML"

	defaultIceServerYaml = "/opt/webrtc-test/configs/ice_servers.yaml"

	stunPrefix = "stun"
	turnPrefix = "turn"
)

var (
	iceServerYaml string
)

var (
	iceServers    *iceServerYamlType
	rtcIceServers []webrtc.ICEServer
)

type iceServerYamlType struct {
	Servers []IceServer `yaml:"servers"`
}

type IceServer struct {
	Credential     string   `json:"credential" yaml:"credential"`
	CredentialType string   `json:"credentialType" yaml:"credentialType"`
	URL            string   `json:"url" yaml:"url"`
	URLs           []string `json:"urls" yaml:"urls"`
	Username       string   `json:"username" yaml:"username"`
}

type CandidateFromClient struct {
	Candidate string `json:"candidate"`
	User      string `json:"user"`
}

func (i *IceServer) toIceServer() webrtc.ICEServer {
	newRIS := webrtc.ICEServer{}

	validURLs := make([]string, 0)
	if len(i.URLs) > 0 {
		for _, url := range i.URLs {
			if strings.HasPrefix(strings.TrimSpace(url), turnPrefix) {
				validURLs = append(validURLs, strings.TrimSpace(url))
			} else if strings.HasPrefix(strings.TrimSpace(url), stunPrefix) {
				validURLs = append(validURLs, strings.TrimSpace(url))
			}
		}
	}
	newRIS.URLs = validURLs

	if i.Credential != "" {
		newRIS.Credential = i.Credential
	}

	if strings.TrimSpace(i.CredentialType) != "" {
		switch strings.ToLower(strings.TrimSpace(i.CredentialType)) {
		case "password":
			newRIS.CredentialType = webrtc.ICECredentialTypePassword
		case "oauth":
			newRIS.CredentialType = webrtc.ICECredentialTypeOauth
		default:
			newRIS.CredentialType = webrtc.ICECredentialTypePassword
		}
	}

	if i.Username != "" {
		newRIS.Username = i.Username
	}

	return newRIS
}

func init() {
	readICEServerYaml()
}

func readICEServerYaml() {
	envYaml := os.Getenv(iceServerYamlEnvString)
	if envYaml == "" {
		if err := os.Setenv(iceServerYamlEnvString, defaultIceServerYaml); err != nil {
			panic(err)
		}
		iceServerYaml = defaultIceServerYaml
	} else {
		iceServerYaml = envYaml
	}

	if iceServers == nil {
		iceServers = &iceServerYamlType{}
	}

	if err := readYamlFile(iceServerYaml, iceServers); err != nil {
		panic(err)
	}

	if len(iceServers.Servers) > 0 {
		numURLs := 0
		for _, svr := range iceServers.Servers {
			for _, url := range svr.URLs {
				if strings.HasPrefix(url, "stun") || strings.HasPrefix(url, "turn") {
					numURLs++
				}
			}

			ris := svr.toIceServer()
			if len(ris.URLs) > 0 {
				rtcIceServers = append(rtcIceServers, ris)
			}
		}
	}
}

func readYamlFile(file string, i interface{}) error {
	return unmarshalFile(file, i, yaml.Unmarshal)
}

func unmarshalFile(file string, i interface{}, f func([]byte, interface{}) error) error {
	d, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	return f(d, i)
}

func AddICECandidateForUser(id uuid.UUID, user, candidate string) error {
	if _, exists := ongoingCalls[id]; !exists {
		return fmt.Errorf("unable to add ice candidate, call %v does not exist", id)
	}

	if _, exists := ongoingCalls[id].Participants[user]; !exists {
		return fmt.Errorf("unable to add ice candidate, connection for %s does not exist", user)
	}

	iceCandidate := webrtc.ICECandidateInit{Candidate: candidate}
        if err := ongoingCalls[id].Participants[user].PC.AddICECandidate(iceCandidate); err != nil {
		return err
	}

	return nil
}

func GetIceServers() (servers []IceServer) {
	for _, server := range rtcIceServers {
		cred := ""
		credType := ""

		if server.Credential != nil {
			cred = server.Credential.(string)
		}

		switch server.CredentialType.String() {
		case "password":
			credType = "password"
		case "token":
			credType = "token"
		}

		servers = append(servers, IceServer{
			Credential:     cred,
			CredentialType: credType,
			URLs:           server.URLs,
			Username:       server.Username,
		})
	}
	return
}
