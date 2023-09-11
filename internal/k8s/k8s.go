package k8s

import (
	"errors"
	"strings"
)

// Bio wrap container bio
type Bio struct {
	Name       string    `json:"name"`
	Containers []*ConBio `json:"containers"`
}

// ConBio container bio
type ConBio struct {
	Name  string `json:"name"`
	Image string `json:"img"`
	Pull  string `json:"pull"`
}

// ContainerPath ...
type ContainerPath struct {
	Ns    string `json:"namespace"`
	Kind  string `json:"kind" binding:"required"`
	Name  string `json:"name" binding:"required"`
	CName string `json:"container_name" binding:"required"`
	Img   string `json:"img" binding:"required"`
}

func (c ContainerPath) valid() bool {
	if c.CName == "" || c.Img == "" || c.Kind == "" || c.Name == "" || c.Ns == "" {
		return false
	}
	if c.Kind != "deploy" && c.Kind != "sts" {
		return false
	}
	return true
}

func (c ContainerPath) UpdateResource(labels ...string) error {
	if !c.valid() {
		return errors.New("invalid resource param")
	}
	switch c.Kind {
	case "deploy":
		return setDeployImg(&c, labels...)
	case "sts":
		return setStsImg(&c, labels...)
	}
	return errors.New("invalid resource kind")
}

func checkLabels(table map[string]string, labels []string) error {
	for _, l := range labels {
		kv := strings.Split(l, "=")
		if len(kv) != 2 {
			return errors.New("invalid label selector")
		}
		if kv[1] != table[kv[0]] {
			return errors.New("filter out by labels: " + l)
		}
	}
	return nil
}
