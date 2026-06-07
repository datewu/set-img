package k8s

import (
	"errors"
	"strings"

	corev1 "k8s.io/api/core/v1"
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

// EnvVar represents a simple environment variable name and value
type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ContainerPath ...
type ContainerPath struct {
	Ns        string   `json:"namespace"`
	Kind      string   `json:"kind" binding:"required"`
	Name      string   `json:"name" binding:"required"`
	CName     string   `json:"container_name" binding:"required"`
	Img       string   `json:"img,omitempty"`
	Env       []EnvVar `json:"env,omitempty"`
	UpdateEnv bool     `json:"update_env,omitempty"`
}

func (c ContainerPath) valid() bool {
	if c.CName == "" || c.Kind == "" || c.Name == "" || c.Ns == "" {
		return false
	}
	if c.Img == "" && !c.UpdateEnv {
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

// ParseEnvStr parses a multi-line string into a slice of EnvVar
func ParseEnvStr(s string) []EnvVar {
	lines := strings.Split(s, "\n")
	var res []EnvVar
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" || strings.HasPrefix(l, "#") {
			continue
		}
		parts := strings.SplitN(l, "=", 2)
		if len(parts) < 1 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		if key == "" {
			continue
		}
		var val string
		if len(parts) == 2 {
			val = strings.TrimSpace(parts[1])
		}
		res = append(res, EnvVar{
			Name:  key,
			Value: val,
		})
	}
	return res
}

// ParseEnvVarValue parses a string representation of env var value, restoring valueFrom fields if needed
func ParseEnvVarValue(name, valStr string, original *corev1.EnvVar) corev1.EnvVar {
	valStr = strings.TrimSpace(valStr)
	if strings.HasPrefix(valStr, "valueFrom(") && strings.HasSuffix(valStr, ")") {
		content := valStr[len("valueFrom(") : len(valStr)-1]
		parts := strings.Split(content, ":")
		if len(parts) > 0 {
			switch parts[0] {
			case "configMapKeyRef":
				if len(parts) >= 3 {
					return corev1.EnvVar{
						Name: name,
						ValueFrom: &corev1.EnvVarSource{
							ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{Name: parts[1]},
								Key:                  parts[2],
							},
						},
					}
				}
			case "secretKeyRef":
				if len(parts) >= 3 {
					return corev1.EnvVar{
						Name: name,
						ValueFrom: &corev1.EnvVarSource{
							SecretKeyRef: &corev1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{Name: parts[1]},
								Key:                  parts[2],
							},
						},
					}
				}
			case "fieldRef":
				if len(parts) >= 2 {
					return corev1.EnvVar{
						Name: name,
						ValueFrom: &corev1.EnvVarSource{
							FieldRef: &corev1.ObjectFieldSelector{
								FieldPath: parts[1],
							},
						},
					}
				}
			case "resourceFieldRef":
				if len(parts) >= 3 {
					return corev1.EnvVar{
						Name: name,
						ValueFrom: &corev1.EnvVarSource{
							ResourceFieldRef: &corev1.ResourceFieldSelector{
								ContainerName: parts[1],
								Resource:      parts[2],
							},
						},
					}
				}
			}
		}
		if original != nil && original.Name == name && original.ValueFrom != nil {
			return *original
		}
	}
	return corev1.EnvVar{
		Name:  name,
		Value: valStr,
	}
}
