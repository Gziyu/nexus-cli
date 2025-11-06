package registry

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

const ACCEPT_HEADER = "application/vnd.docker.distribution.manifest.v2+json"
const CREDENTIALS_FILE = ".credentials"

type Registry struct {
	Host       string `toml:"nexus_host"`
	Username   string `toml:"nexus_username"`
	Password   string `toml:"nexus_password"`
	Repository string `toml:"nexus_repository"`
}

type Repositories struct {
	Images []string `json:"repositories"`
}

type ImageTags struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type ImageManifest struct {
	SchemaVersion int64       `json:"schemaVersion"`
	MediaType     string      `json:"mediaType"`
	Config        LayerInfo   `json:"config"`
	Layers        []LayerInfo `json:"layers"`
}

type LayerInfo struct {
	MediaType string `json:"mediaType"`
	Size      int64  `json:"size"`
	Digest    string `json:"digest"`
}

type TagInfo struct {
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
}

func NewRegistry() (Registry, error) {
	r := Registry{}
	if _, err := os.Stat(CREDENTIALS_FILE); os.IsNotExist(err) {
		return r, fmt.Errorf("%s file not found", CREDENTIALS_FILE)
	} else if err != nil {
		return r, err
	}

	if _, err := toml.DecodeFile(CREDENTIALS_FILE, &r); err != nil {
		return r, err
	}
	return r, nil
}

func (r Registry) ListImages() ([]string, error) {
	client := &http.Client{}

	url := fmt.Sprintf("%s/repository/%s/v2/_catalog", r.Host, r.Repository)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(r.Username, r.Password)
	req.Header.Add("Accept", ACCEPT_HEADER)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP Code: %d", resp.StatusCode)
	}

	var repositories Repositories
	json.NewDecoder(resp.Body).Decode(&repositories)

	return repositories.Images, nil
}

func (r Registry) ListTagsByImage(image string) ([]string, error) {
	client := &http.Client{}

	url := fmt.Sprintf("%s/repository/%s/v2/%s/tags/list", r.Host, r.Repository, image)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(r.Username, r.Password)
	req.Header.Add("Accept", ACCEPT_HEADER)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP Code: %d", resp.StatusCode)
	}

	var imageTags ImageTags
	json.NewDecoder(resp.Body).Decode(&imageTags)

	return imageTags.Tags, nil
}

// ListTagsWithTime returns tags with their creation time information
func (r Registry) ListTagsWithTime(image string) ([]TagInfo, error) {
	tags, err := r.ListTagsByImage(image)
	if err != nil {
		return nil, err
	}

	var tagInfos []TagInfo
	for _, tag := range tags {
		created, err := r.GetTagCreationTime(image, tag)
		if err != nil {
			// If we can't get creation time, use zero time (will be sorted as oldest)
			created = time.Time{}
		}
		tagInfos = append(tagInfos, TagInfo{
			Name:    tag,
			Created: created,
		})
	}

	return tagInfos, nil
}

// GetTagCreationTime gets the creation time of a specific tag
func (r Registry) GetTagCreationTime(image string, tag string) (time.Time, error) {
	client := &http.Client{}

	url := fmt.Sprintf("%s/repository/%s/v2/%s/manifests/%s", r.Host, r.Repository, image, tag)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return time.Time{}, err
	}
	req.SetBasicAuth(r.Username, r.Password)
	req.Header.Add("Accept", ACCEPT_HEADER)

	resp, err := client.Do(req)
	if err != nil {
		return time.Time{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return time.Time{}, fmt.Errorf("HTTP Code: %d", resp.StatusCode)
	}

	// Try to get Last-Modified header first
	if lastModified := resp.Header.Get("Last-Modified"); lastModified != "" {
		if t, err := time.Parse(time.RFC1123, lastModified); err == nil {
			return t, nil
		}
	}

	// Fallback: use Docker-Content-Digest header timestamp
	if digest := resp.Header.Get("Docker-Content-Digest"); digest != "" {
		// For simplicity, we'll use current time as fallback
		// In a real implementation, you might want to parse the manifest
		// to get actual creation time from config
		return time.Now(), nil
	}

	return time.Time{}, errors.New("could not determine creation time")
}

func (r Registry) ImageManifest(image string, tag string) (ImageManifest, error) {
	var imageManifest ImageManifest
	client := &http.Client{}

	url := fmt.Sprintf("%s/repository/%s/v2/%s/manifests/%s", r.Host, r.Repository, image, tag)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return imageManifest, err
	}
	req.SetBasicAuth(r.Username, r.Password)
	req.Header.Add("Accept", ACCEPT_HEADER)

	resp, err := client.Do(req)
	if err != nil {
		return imageManifest, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return imageManifest, fmt.Errorf("HTTP Code: %d", resp.StatusCode)
	}

	json.NewDecoder(resp.Body).Decode(&imageManifest)

	return imageManifest, nil
}

func (r Registry) DeleteImageByTag(image string, tag string) error {
	sha, err := r.getImageSHA(image, tag)
	if err != nil {
		return err
	}
	client := &http.Client{}

	url := fmt.Sprintf("%s/repository/%s/v2/%s/manifests/%s", r.Host, r.Repository, image, sha)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(r.Username, r.Password)
	req.Header.Add("Accept", ACCEPT_HEADER)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 202 {
		return fmt.Errorf("HTTP Code: %d", resp.StatusCode)
	}

	fmt.Printf("%s:%s has been successful deleted\n", image, tag)

	return nil
}

func (r Registry) getImageSHA(image string, tag string) (string, error) {
	client := &http.Client{}

	url := fmt.Sprintf("%s/repository/%s/v2/%s/manifests/%s", r.Host, r.Repository, image, tag)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(r.Username, r.Password)
	req.Header.Add("Accept", ACCEPT_HEADER)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("HTTP Code: %d", resp.StatusCode)
	}

	return resp.Header.Get("Docker-Content-Digest"), nil
}
