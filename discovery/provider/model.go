//go:generate go run ../../pkg/sdk/runable/steampipe_es_client_generator/main.go -pluginPath ../../steampipe-plugin-REPLACEME/REPLACEME -file $GOFILE -output ../../pkg/sdk/es/resources_clients.go -resourceTypesFile ../resource_types/resource-types.json

// Implement types for each resource

package provider

type Metadata struct{}

type ListAppsResponse struct {
	Apps []AppJSON `json:"apps"`
}

type AppJSON struct {
	ID           string `json:"id"`
	MachineCount int    `json:"machine_count,omitempty"`
	Name         string `json:"name"`
	Network      string    `json:"network,omitempty"`
}

type AppDescription struct {
	ID           string
	MachineCount int
	Name         string
	Network      string
}

type CheckJSON struct {
	Name      string `json:"name"`
	Output    string `json:"output"`
	Status    string `json:"status"`
	UpdatedAt string `json:"updated_at"`
}

type Check struct {
	Name      string
	Output    string
	Status    string
	UpdatedAt string
}

type ConfigCheckJSON struct {
	GracePeriod   any           `json:"grace_period"`
	Headers       []interface{} `json:"headers"`
	Interval      any           `json:"interval"`
	Kind          string        `json:"kind"`
	Method        string        `json:"method"`
	Path          string        `json:"path"`
	Port          int           `json:"port"`
	Protocol      string        `json:"protocol"`
	Timeout       any           `json:"timeout"`
	TLSServerName string        `json:"tls_server_name"`
	TLSSkipVerify bool          `json:"tls_skip_verify"`
	Type          string        `json:"type"`
}

type ConfigCheck struct {
	GracePeriod   any
	Headers       []interface{}
	Interval      any
	Kind          string
	Method        string
	Path          string
	Port          int
	Protocol      string
	Timeout       any
	TLSServerName string
	TLSSkipVerify bool
	Type          string
}

type ContainerJSON struct {
	Cmd        []interface{} `json:"cmd"`
	DependsOn  []interface{} `json:"depends_on"`
	Entrypoint []interface{} `json:"entrypoint"`
	Env        interface{}   `json:"env"`
	EnvFrom    []interface{} `json:"env_from"`
	Exec       []interface{} `json:"exec"`
	Files      []interface{} `json:"files"`
	Image      string        `json:"image"`
	Name       string        `json:"name"`
	Restart    interface{}   `json:"restart"`
	Secrets    []interface{} `json:"secrets"`
	Stop       interface{}   `json:"stop"`
	User       string        `json:"user"`
}

type Container struct {
	Cmd        []interface{}
	DependsOn  []interface{}
	Entrypoint []interface{}
	Env        interface{}
	EnvFrom    []interface{}
	Exec       []interface{}
	Files      []interface{}
	Image      string
	Name       string
	Restart    interface{}
	Secrets    []interface{}
	Stop       interface{}
	User       string
}

type DNSForwardRuleJSON struct {
	Addr     interface{} `json:"addr"`
	Basename interface{} `json:"basename"`
}

type DNSForwardRule struct {
	Addr     interface{}
	Basename interface{}
}

type DNSOptionJSON struct {
	Name  interface{} `json:"name"`
	Value interface{} `json:"value"`
}

type DNSOption struct {
	Name  interface{}
	Value interface{}
}

type DNSJSON struct {
	DNSForwardRules  []DNSForwardRuleJSON `json:"dns_forward_rules"`
	Hostname         string               `json:"hostname"`
	HostnameFQDN     string               `json:"hostname_fqdn"`
	Nameservers      []string             `json:"nameservers"`
	Options          []DNSOptionJSON      `json:"options"`
	Searches         []string             `json:"searches"`
	SkipRegistration bool                 `json:"skip_registration"`
}

type DNS struct {
	DNSForwardRules  []DNSForwardRule
	Hostname         string
	HostnameFQDN     string
	Nameservers      []string
	Options          []DNSOption
	Searches         []string
	SkipRegistration bool
}

type FileJSON struct {
	GuestPath  string `json:"guest_path"`
	Mode       int    `json:"mode"`
	RawValue   string `json:"raw_value"`
	SecretName string `json:"secret_name"`
}

type File struct {
	GuestPath  string `json:"guest_path"`
	Mode       int    `json:"mode"`
	RawValue   string `json:"raw_value"`
	SecretName string `json:"secret_name"`
}

type GuestJSON struct {
	CPUKind          string   `json:"cpu_kind"`
	CPUs             int      `json:"cpus"`
	GPUKind          string   `json:"gpu_kind"`
	GPUs             int      `json:"gpus"`
	HostDedicationID string   `json:"host_dedication_id"`
	KernelArgs       []string `json:"kernel_args"`
	MemoryMB         int      `json:"memory_mb"`
}

type Guest struct {
	CPUKind          string
	CPUs             int
	GPUKind          string
	GPUs             int
	HostDedicationID string
	KernelArgs       []string
	MemoryMB         int
}

type InitJSON struct {
	Cmd        []string `json:"cmd"`
	Entrypoint []string `json:"entrypoint"`
	Exec       []string `json:"exec"`
	KernelArgs []string `json:"kernel_args"`
	SwapSizeMB int      `json:"swap_size_mb"`
	TTY        bool     `json:"tty"`
}

type Init struct {
	Cmd        []string
	Entrypoint []string
	Exec       []string
	KernelArgs []string
	SwapSizeMB int
	TTY        bool
}

type MetricJSON struct {
	HTTPS bool   `json:"https"`
	Path  string `json:"path"`
	Port  int    `json:"port"`
}

type Metric struct {
	HTTPS bool
	Path  string
	Port  int
}

type MountJSON struct {
	AddSizeGB              int    `json:"add_size_gb"`
	Encrypted              bool   `json:"encrypted"`
	ExtendThresholdPercent int    `json:"extend_threshold_percent"`
	Name                   string `json:"name"`
	Path                   string `json:"path"`
	SizeGB                 int    `json:"size_gb"`
	SizeGBLimit            int    `json:"size_gb_limit"`
	Volume                 string `json:"volume"`
}

type Mount struct {
	AddSizeGB              int
	Encrypted              bool
	ExtendThresholdPercent int
	Name                   string
	Path                   string
	SizeGB                 int
	SizeGBLimit            int
	Volume                 string
}

type ProcessJSON struct {
	Cmd              []interface{} `json:"cmd"`
	Entrypoint       []interface{} `json:"entrypoint"`
	Env              interface{}   `json:"env"`
	EnvFrom          []interface{} `json:"env_from"`
	Exec             []interface{} `json:"exec"`
	IgnoreAppSecrets bool          `json:"ignore_app_secrets"`
	Secrets          []interface{} `json:"secrets"`
	User             string        `json:"user"`
}

type Process struct {
	Cmd              []interface{}
	Entrypoint       []interface{}
	Env              interface{}
	EnvFrom          []interface{}
	Exec             []interface{}
	IgnoreAppSecrets bool
	Secrets          []interface{}
	User             string
}

type RestartJSON struct {
	GPUBidPrice int    `json:"gpu_bid_price"`
	MaxRetries  int    `json:"max_retries"`
	Policy      string `json:"policy"`
}

type Restart struct {
	GPUBidPrice int
	MaxRetries  int
	Policy      string
}

type ServiceConcurrencyJSON struct {
	HardLimit interface{} `json:"hard_limit"`
	SoftLimit interface{} `json:"soft_limit"`
	Type      interface{} `json:"type"`
}

type ServiceConcurrency struct {
	HardLimit interface{}
	SoftLimit interface{}
	Type      interface{}
}

type ServiceJSON struct {
	Autostart                bool                   `json:"autostart"`
	Autostop                 bool                   `json:"autostop"`
	Checks                   []interface{}          `json:"checks"`
	Concurrency              ServiceConcurrencyJSON `json:"concurrency"`
	ForceInstanceDescription string                 `json:"force_instance_description"`
	ForceInstanceKey         string                 `json:"force_instance_key"`
	InternalPort             int                    `json:"internal_port"`
	MinMachinesRunning       int                    `json:"min_machines_running"`
	Ports                    []interface{}          `json:"ports"`
	Protocol                 string                 `json:"protocol"`
}

type Service struct {
	Autostart                bool
	Autostop                 bool
	Checks                   []interface{}
	Concurrency              ServiceConcurrency
	ForceInstanceDescription string
	ForceInstanceKey         string
	InternalPort             int
	MinMachinesRunning       int
	Ports                    []interface{}
	Protocol                 string
}

type StaticJSON struct {
	GuestPath     string `json:"guest_path"`
	IndexDocument string `json:"index_document"`
	TigrisBucket  string `json:"tigris_bucket"`
	URLPrefix     string `json:"url_prefix"`
}

type Static struct {
	GuestPath     string
	IndexDocument string
	TigrisBucket  string
	URLPrefix     string
}

type StopConfigJSON struct {
	Signal  string `json:"signal"`
	Timeout struct {
		TimeDuration int `json:"time.Duration"`
	} `json:"timeout"`
}

type StopConfig struct {
	Signal  string
	Timeout struct {
		TimeDuration int
	}
}

type ConfigJSON struct {
	AutoDestroy             bool              `json:"auto_destroy"`
	Checks                  ConfigCheckJSON   `json:"checks"`
	Containers              []ContainerJSON   `json:"containers"`
	DisableMachineAutostart bool              `json:"disable_machine_autostart"`
	DNS                     DNSJSON           `json:"dns"`
	Env                     map[string]string `json:"env"`
	Files                   []FileJSON        `json:"files"`
	Guest                   GuestJSON         `json:"guest"`
	Image                   string            `json:"image"`
	Init                    InitJSON          `json:"init"`
	Metadata                map[string]string `json:"metadata"`
	Metrics                 MetricJSON        `json:"metrics"`
	Mounts                  []MountJSON       `json:"mounts"`
	Processes               []ProcessJSON     `json:"processes"`
	Restart                 RestartJSON       `json:"restart"`
	Schedule                string            `json:"schedule"`
	Services                []ServiceJSON     `json:"services"`
	Size                    string            `json:"size"`
	Standbys                []string          `json:"standbys"`
	Statics                 []StaticJSON      `json:"statics"`
	StopConfig              StopConfigJSON    `json:"stop_config"`
}

type Config struct {
	AutoDestroy             bool
	Checks                  ConfigCheck
	Containers              []Container
	DisableMachineAutostart bool
	DNS                     DNS
	Env                     map[string]string
	Files                   []File
	Guest                   Guest
	Image                   string
	Init                    Init
	Metadata                map[string]string
	Metrics                 Metric
	Mounts                  []Mount
	Processes               []Process
	Restart                 Restart
	Schedule                string
	Services                []Service
	Size                    string
	Standbys                []string
	Statics                 []Static
	StopConfig              StopConfig
}

type EventJSON struct {
	ID        string      `json:"id"`
	Request   interface{} `json:"request"`
	Source    string      `json:"source"`
	Status    string      `json:"status"`
	Timestamp int         `json:"timestamp"`
	Type      string      `json:"type"`
}

type Event struct {
	ID        string
	Request   interface{}
	Source    string
	Status    string
	Timestamp int
	Type      string
}

type ImageRefJSON struct {
	Digest     string            `json:"digest"`
	Labels     map[string]string `json:"labels"`
	Registry   string            `json:"registry"`
	Repository string            `json:"repository"`
	Tag        string            `json:"tag"`
}

type ImageRef struct {
	Digest     string
	Labels     map[string]string
	Registry   string
	Repository string
	Tag        string
}

type MachineJSON struct {
	Checks           []CheckJSON  `json:"checks"`
	Config           ConfigJSON   `json:"config"`
	CreatedAt        string       `json:"created_at"`
	Events           []EventJSON  `json:"events"`
	HostStatus       string       `json:"host_status"`
	ID               string       `json:"id"`
	ImageRef         ImageRefJSON `json:"image_ref"`
	IncompleteConfig ConfigJSON   `json:"incomplete_config"`
	InstanceID       string       `json:"instance_id"`
	Name             string       `json:"name"`
	Nonce            string       `json:"nonce"`
	PrivateIP        string       `json:"private_ip"`
	Region           string       `json:"region"`
	State            string       `json:"state"`
	UpdatedAt        string       `json:"updated_at"`
}

type MachineDescription struct {
	Checks           []Check
	Config           Config
	CreatedAt        string
	Events           []Event
	HostStatus       string
	ID               string
	ImageRef         ImageRef
	IncompleteConfig Config
	InstanceID       string
	Name             string
	Nonce            string
	PrivateIP        string
	Region           string
	State            string
	UpdatedAt        string
}

type VolumeJSON struct {
	AttachedAllocID   string `json:"attached_alloc_id"`
	AttachedMachineID string `json:"attached_machine_id"`
	AutoBackupEnabled bool   `json:"auto_backup_enabled"`
	BlockSize         int    `json:"block_size"`
	Blocks            int    `json:"blocks"`
	BlocksAvail       int    `json:"blocks_avail"`
	BlocksFree        int    `json:"blocks_free"`
	CreatedAt         string `json:"created_at"`
	Encrypted         bool   `json:"encrypted"`
	FSType            string `json:"fstype"`
	HostStatus        string `json:"host_status"`
	ID                string `json:"id"`
	Name              string `json:"name"`
	Region            string `json:"region"`
	SizeGB            int    `json:"size_gb"`
	SnapshotRetention int    `json:"snapshot_retention"`
	State             string `json:"state"`
	Zone              string `json:"zone"`
}

type VolumeDescription struct {
	AttachedAllocID   string
	AttachedMachineID string
	AutoBackupEnabled bool
	BlockSize         int
	Blocks            int
	BlocksAvail       int
	BlocksFree        int
	CreatedAt         string
	Encrypted         bool
	FSType            string
	HostStatus        string
	ID                string
	Name              string
	Region            string
	SizeGB            int
	SnapshotRetention int
	State             string
	Zone              string
}

type SecretJSON struct {
	Label     string `json:"label"`
	PublicKey []int  `json:"publickey"`
	Type      string `json:"type"`
}

type SecretDescription struct {
	Label     string
	PublicKey []int
	Type      string
}
