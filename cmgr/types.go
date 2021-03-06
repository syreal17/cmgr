package cmgr

import (
	"context"
	"math/rand"

	"github.com/docker/docker/client"
	"github.com/jmoiron/sqlx"
)

const (
	DB_ENV           string = "CMGR_DB"
	DIR_ENV          string = "CMGR_DIR"
	ARTIFACT_DIR_ENV string = "CMGR_ARTIFACT_DIR"
	REGISTRY_ENV     string = "CMGR_REGISTRY"
	LOGGING_ENV      string = "CMGR_LOGGING"
	IFACE_ENV        string = "CMGR_INTERFACE"

	DYNAMIC_INSTANCES int = -1
	LOCKED            int = -2
)

type UnknownIdentifierError struct {
	Type string
	Name string
}

type Manager struct {
	cli                  *client.Client
	ctx                  context.Context
	log                  *logger
	chalDir              string
	artifactsDir         string
	db                   *sqlx.DB
	dbPath               string
	challengeDockerfiles map[string][]byte
	rand                 *rand.Rand
	challengeInterface   string
}

type PortInfo struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type HostInfo struct {
	Name   string `json:"name"`
	Target string `json:"target,omitempty"`
}

type ChallengeId string
type ChallengeMetadata struct {
	Id               ChallengeId         `json:"id"`
	Name             string              `json:"name,omitempty"`
	Namespace        string              `json:"namespace"`
	ChallengeType    string              `json:"challenge_type"`
	Description      string              `json:"description,omitempty"`
	Details          string              `json:"details,omitempty"`
	Hints            []string            `json:"hints,omitempty"`
	SourceChecksum   uint32              `json:"source_checksum"`
	MetadataChecksum uint32              `json:"metadata_checksum`
	Path             string              `json:"path"`
	Templatable      bool                `json:"templatable,omitempty"`
	PortMap          map[string]PortInfo `json:"port_map,omitempty"`
	Hosts            []HostInfo          `json:"hosts"`
	MaxUsers         int                 `json:"max_users,omitempty"`
	Category         string              `json:"category,omitempty"`
	Points           int                 `json:"points,omitempty"`
	Tags             []string            `json:"tags,omitempty"`
	Attributes       map[string]string   `json:"attributes,omitempty"`

	SolveScript bool             `json:"solve_script,omitempty"`
	Builds      []*BuildMetadata `json:"builds,omitempty"`
}
type ChallengeUpdates struct {
	Added      []*ChallengeMetadata `json:"added"`
	Refreshed  []*ChallengeMetadata `json:"refreshed"`
	Updated    []*ChallengeMetadata `json:"updated"`
	Removed    []*ChallengeMetadata `json:"removed"`
	Unmodified []*ChallengeMetadata `json:"unmodified"`
	Errors     []error              `json:"errors"`
}

type BuildId int64
type BuildMetadata struct {
	Id BuildId `json:"id"`

	Flag       string            `json:"flag"`
	LookupData map[string]string `json:"lookup_data,omitempty"`

	Seed         int                 `json:"seed"`
	Format       string              `json:"format"`
	Images       []Image             `json:"images"`
	HasArtifacts bool                `json:"has_artifacts"`
	LastSolved   int64               `json:"last_solved"`
	Challenge    ChallengeId         `json:"challenge_id"`
	Instances    []*InstanceMetadata `json:"instances,omitempty"`

	Schema        string `json:"schema"`
	InstanceCount int    `json:"instance_count"`
}

type ImageId int64
type Image struct {
	Id    ImageId  `json:"id"`
	Host  string   `json:"host"`
	Ports []string `json:"exposed_ports"`
	Build BuildId  `json:"build"`
}

type InstanceId int64
type InstanceMetadata struct {
	Id         InstanceId     `json:"id"`
	Ports      map[string]int `json:"ports,omitempty"`
	Containers []string       `json:"containers"`
	LastSolved int64          `json:"last_solved"`
	Build      BuildId        `json:"build_id"`
}

type Schema struct {
	Name       string                             `json:"name" yaml:"name"`
	FlagFormat string                             `json:"flag_format" yaml:"flag_format"`
	Challenges map[ChallengeId]BuildSpecification `json:"challenges" yaml:"challenges"`
}
type BuildSpecification struct {
	Seeds         []int `json:"seeds" yaml:"seeds"`
	InstanceCount int   `json:"instance_count" yaml:"instance_count"`
}
