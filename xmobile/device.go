package xmobile

import (
	"encoding/json"
	"errors"
	"strings"

	"go.olapie.com/x/xconv"
)

type DeviceInfo struct {
	Name       string `json:"name,omitempty"`
	Model      string `json:"model,omitempty"`
	ModelType  string `json:"model_type,omitempty"`
	Language   string `json:"language,omitempty"`
	SysName    string `json:"sys_name,omitempty"`
	SysVersion string `json:"sys_version,omitempty"`
	Carrier    string `json:"carrier,omitempty"`
}

func (i *DeviceInfo) Validate() error {
	i.Name = strings.TrimSpace(i.Name)
	i.Model = strings.TrimSpace(i.Model)
	i.ModelType = strings.TrimSpace(i.Name)
	i.Language = strings.TrimSpace(i.Language)
	i.SysName = strings.TrimSpace(i.SysName)
	i.SysVersion = strings.TrimSpace(i.SysVersion)
	i.Carrier = strings.TrimSpace(i.Carrier)
	if i.Name == "" {
		return errors.New("missing Name")
	}
	if i.ModelType == "" {
		return errors.New("missing ModelType")
	}
	if i.Language == "" {
		return errors.New("missing Language")
	}
	return nil
}

func NewDeviceInfo() *DeviceInfo {
	return new(DeviceInfo)
}

func (d *DeviceInfo) Attributes() map[string]string {
	m := make(map[string]string)
	err := json.Unmarshal(xconv.MustToJSONBytes(d), &m)
	if err != nil {
		panic(err)
	}
	return m
}

type AppInfo struct {
	AppID    string `json:"id,omitempty"`
	BundleID string `json:"bundleId,omitempty"`
	Name     string `json:"name,omitempty"`
	Version  string `json:"version,omitempty"`
}

func NewAppInfo() *AppInfo {
	return new(AppInfo)
}

func (i *AppInfo) Attributes() map[string]string {
	m := make(map[string]string)
	err := json.Unmarshal(xconv.MustToJSONBytes(i), &m)
	if err != nil {
		panic(err)
	}
	return m
}

func (i *AppInfo) Validate() error {
	i.AppID = strings.TrimSpace(i.AppID)
	i.BundleID = strings.TrimSpace(i.BundleID)
	i.Name = strings.TrimSpace(i.Name)
	i.Version = strings.TrimSpace(i.Version)
	if i.AppID == "" {
		return errors.New("missing AppID")
	}
	if i.BundleID == "" {
		return errors.New("missing BundleID")
	}
	if i.Name == "" {
		return errors.New("missing Name")
	}
	if i.Version == "" {
		return errors.New("missing Version")
	}
	return nil
}
