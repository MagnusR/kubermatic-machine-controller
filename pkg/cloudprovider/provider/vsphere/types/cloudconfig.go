/*
Copyright 2019 The Machine Controller Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package types

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig/v3"

	"github.com/magnusr/kubermatic-machine-controller/pkg/ini"
)

const (
	cloudConfigTpl = `[Global]
user              = {{ .Global.User | iniEscape }}
password          = {{ .Global.Password | iniEscape }}
port              = {{ .Global.VCenterPort | iniEscape }}
insecure-flag     = {{ .Global.InsecureFlag }}
working-dir       = {{ .Global.WorkingDir | iniEscape }}
datacenter        = {{ .Global.Datacenter | iniEscape }}
datastore         = {{ .Global.DefaultDatastore | iniEscape }}
server            = {{ .Global.VCenterIP | iniEscape }}
{{- if .Global.IPFamily }}
ip-family         = {{ .Global.IPFamily | iniEscape }}
{{- end }}

[Disk]
scsicontrollertype = {{ .Disk.SCSIControllerType | iniEscape }}

[Workspace]
server            = {{ .Workspace.VCenterIP | iniEscape }}
datacenter        = {{ .Workspace.Datacenter | iniEscape }}
folder            = {{ .Workspace.Folder | iniEscape }}
default-datastore = {{ .Workspace.DefaultDatastore | iniEscape }}
resourcepool-path = {{ .Workspace.ResourcePoolPath | iniEscape }}

{{ range $name, $vc := .VirtualCenter }}
[VirtualCenter {{ $name | iniEscape }}]
user = {{ $vc.User | iniEscape }}
password = {{ $vc.Password | iniEscape }}
port = {{ $vc.VCenterPort }}
datacenters = {{ $vc.Datacenters | iniEscape }}
{{- if $vc.IPFamily }}
ip-family = {{ $vc.IPFamily | iniEscape }}
{{- end }}
{{ end }}
`
)

type WorkspaceOpts struct {
	VCenterIP        string `gcfg:"server"`
	Datacenter       string `gcfg:"datacenter"`
	Folder           string `gcfg:"folder"`
	DefaultDatastore string `gcfg:"default-datastore"`
	ResourcePoolPath string `gcfg:"resourcepool-path"`
}

type DiskOpts struct {
	SCSIControllerType string `dcfg:"scsicontrollertype"`
}

type GlobalOpts struct {
	User             string `gcfg:"user"`
	Password         string `gcfg:"password"`
	InsecureFlag     bool   `gcfg:"insecure-flag"`
	VCenterPort      string `gcfg:"port"`
	WorkingDir       string `gcfg:"working-dir"`
	Datacenter       string `gcfg:"datacenter"`
	DefaultDatastore string `gcfg:"datastore"`
	VCenterIP        string `gcfg:"server"`
	ClusterID        string `gcfg:"cluster-id"`
	IPFamily         string `gcfg:"ip-family"` // NOTE: supported only in case of out-of-tree CCM
}

type VirtualCenterConfig struct {
	User        string `gcfg:"user"`
	Password    string `gcfg:"password"`
	VCenterPort string `gcfg:"port"`
	Datacenters string `gcfg:"datacenters"`
	IPFamily    string `gcfg:"ip-family"` // NOTE: supported only in case of out-of-tree CCM
}

// CloudConfig is used to read and store information from the cloud configuration file.
type CloudConfig struct {
	Global    GlobalOpts
	Disk      DiskOpts
	Workspace WorkspaceOpts

	VirtualCenter map[string]*VirtualCenterConfig
}

// String converts CloudConfig into its formatted string representation.
func (c *CloudConfig) String() (string, error) {
	funcMap := sprig.TxtFuncMap()
	funcMap["iniEscape"] = ini.Escape

	tpl, err := template.New("cloud-config").Funcs(funcMap).Parse(cloudConfigTpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse the cloud config template: %w", err)
	}

	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, c); err != nil {
		return "", fmt.Errorf("failed to execute cloud config template: %w", err)
	}

	return buf.String(), nil
}

// CloudConfigToString converts CloudConfig into its formatted string representation.
// Deprecated: use struct receiver function String() instead.
func CloudConfigToString(c *CloudConfig) (string, error) {
	funcMap := sprig.TxtFuncMap()
	funcMap["iniEscape"] = ini.Escape

	tpl, err := template.New("cloud-config").Funcs(funcMap).Parse(cloudConfigTpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse the cloud config template: %w", err)
	}

	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, c); err != nil {
		return "", fmt.Errorf("failed to execute cloud config template: %w", err)
	}

	return buf.String(), nil
}
