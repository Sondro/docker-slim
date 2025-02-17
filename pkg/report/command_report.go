package report

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/docker-slim/docker-slim/internal/app/master/docker/dockerfile"
	"github.com/docker-slim/docker-slim/pkg/util/errutil"
)

// Command state constants
const (
	CmdStateUnknown   = "unknown"
	CmdStateError     = "error"
	CmdStateStarted   = "started"
	CmdStateCompleted = "completed"
	CmdStateExited    = "exited"
	CmdStateDone      = "done"
)

// Command type constants
const (
	CmdTypeBuild   CmdType = "build"
	CmdTypeProfile CmdType = "profile"
	CmdTypeInfo    CmdType = "info"
)

// CmdType is the command name data type
type CmdType string

// Command is the common command report data
type Command struct {
	reportLocation string
	Type           CmdType `json:"type"`
	State          string  `json:"state"`
	Error          string  `json:"error,omitempty"`
}

// ImageMetadata provides basic image metadata
type ImageMetadata struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Size          int64    `json:"size"`
	SizeHuman     string   `json:"size_human"`
	CreateTime    string   `json:"create_time"`
	AllNames      []string `json:"all_names"`
	Author        string   `json:"Author,omitempty"`
	DockerVersion string   `json:"docker_version"`
	Architecture  string   `json:"architecture"`
	User          string   `json:"user,omitempty"`
	ExposedPorts  []string `json:"exposed_ports,omitempty"`
}

// SystemMetadata provides basic system metadata
type SystemMetadata struct {
	Type    string `json:"type"`
	Release string `json:"release"`
	OS      string `json:"os"`
}

// BuildCommand is the 'build' command report data
type BuildCommand struct {
	Command
	ImageReference         string                  `json:"image_reference"`
	System                 SystemMetadata          `json:"system"`
	SourceImage            ImageMetadata           `json:"source_image"`
	MinifiedImageSize      int64                   `json:"minified_image_size"`
	MinifiedImageSizeHuman string                  `json:"minified_image_size_human"`
	MinifiedImage          string                  `json:"minified_image"`
	MinifiedImageHasData   bool                    `json:"minified_image_has_data"`
	MinifiedBy             float64                 `json:"minified_by"`
	ArtifactLocation       string                  `json:"artifact_location"`
	ContainerReportName    string                  `json:"container_report_name"`
	SeccompProfileName     string                  `json:"seccomp_profile_name"`
	AppArmorProfileName    string                  `json:"apparmor_profile_name"`
	ImageStack             []*dockerfile.ImageInfo `json:"image_stack"`
}

// ProfileCommand is the 'profile' command report data
type ProfileCommand struct {
	Command
	OriginalImage          string  `json:"original_image"`
	OriginalImageSize      int64   `json:"original_image_size"`
	OriginalImageSizeHuman string  `json:"original_image_size_human"`
	MinifiedImageSize      int64   `json:"minified_image_size"`
	MinifiedImageSizeHuman string  `json:"minified_image_size_human"`
	MinifiedImage          string  `json:"minified_image"`
	MinifiedImageHasData   bool    `json:"minified_image_has_data"`
	MinifiedBy             float64 `json:"minified_by"`
	ArtifactLocation       string  `json:"artifact_location"`
	ContainerReportName    string  `json:"container_report_name"`
	SeccompProfileName     string  `json:"seccomp_profile_name"`
	AppArmorProfileName    string  `json:"apparmor_profile_name"`
}

// InfoCommand is the 'info' command report data
type InfoCommand struct {
	Command
	OriginalImage          string  `json:"original_image"`
	OriginalImageSize      int64   `json:"original_image_size"`
	OriginalImageSizeHuman string  `json:"original_image_size_human"`
	MinifiedImageSize      int64   `json:"minified_image_size"`
	MinifiedImageSizeHuman string  `json:"minified_image_size_human"`
	MinifiedImage          string  `json:"minified_image"`
	MinifiedImageHasData   bool    `json:"minified_image_has_data"`
	MinifiedBy             float64 `json:"minified_by"`
	ArtifactLocation       string  `json:"artifact_location"`
	ContainerReportName    string  `json:"container_report_name"`
	SeccompProfileName     string  `json:"seccomp_profile_name"`
	AppArmorProfileName    string  `json:"apparmor_profile_name"`
}

// NewBuildCommand creates a new 'build' command report
func NewBuildCommand(reportLocation string) *BuildCommand {
	return &BuildCommand{
		Command: Command{
			reportLocation: reportLocation,
			Type:           CmdTypeBuild,
			State:          CmdStateUnknown,
		},
	}
}

// NewProfileCommand creates a new 'profile' command report
func NewProfileCommand(reportLocation string) *ProfileCommand {
	return &ProfileCommand{
		Command: Command{
			reportLocation: reportLocation,
			Type:           CmdTypeProfile,
			State:          CmdStateUnknown,
		},
	}
}

// NewInfoCommand creates a new 'info' command report
func NewInfoCommand(reportLocation string) *InfoCommand {
	return &InfoCommand{
		Command: Command{
			reportLocation: reportLocation,
			Type:           CmdTypeInfo,
			State:          CmdStateUnknown,
		},
	}
}

func (p *Command) saveInfo(info interface{}) {
	if p.reportLocation != "" {
		dirName := filepath.Dir(p.reportLocation)
		baseName := filepath.Base(p.reportLocation)

		if baseName == "." {
			fmt.Printf("no build command report location: %v\n", p.reportLocation)
			return
		}

		if dirName != "." {
			_, err := os.Stat(dirName)
			if os.IsNotExist(err) {
				os.MkdirAll(dirName, 0777)
				_, err = os.Stat(dirName)
				errutil.FailOn(err)
			}
		}

		var reportData bytes.Buffer
		encoder := json.NewEncoder(&reportData)
		encoder.SetEscapeHTML(false)
		encoder.SetIndent("", "  ")
		err := encoder.Encode(info)
		errutil.FailOn(err)

		err = ioutil.WriteFile(p.reportLocation, reportData.Bytes(), 0644)
		errutil.FailOn(err)
	}
}

// Save saves the report data to the configured location
func (p *Command) Save() {
	p.saveInfo(p)
}

// Save saves the Build command report data to the configured location
func (p *BuildCommand) Save() {
	p.saveInfo(p)
}

// Save saves the Profile command report data to the configured location
func (p *ProfileCommand) Save() {
	p.saveInfo(p)
}

// Save saves the Info command report data to the configured location
func (p *InfoCommand) Save() {
	p.saveInfo(p)
}
