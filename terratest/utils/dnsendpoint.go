package utils

import "fmt"

type DNSEndpoint struct {
	Metadata Metadata `json:"metadata"`
	Spec Spec `json:"spec"`
}

type Metadata struct {
	Name string `json:"name"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

type Spec struct {
	Endpoints []Endpoint `json:"endpoints,omitempty"`
}

type Endpoint struct{
	Name string `json:"dnsName"`
	RecordType string `json:"recordType"`
	RecordTTL int `json:"recordTTL"`
	Targets []string `json:"targets,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
}

func (dep DNSEndpoint) GetEndpointByName(name string) (Endpoint, error) {
	for _, ep := range dep.Spec.Endpoints {
		if ep.Name == name {
			return ep, nil
		}
	}
	return Endpoint{}, fmt.Errorf("not found")
}