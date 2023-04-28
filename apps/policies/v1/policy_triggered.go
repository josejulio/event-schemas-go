// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    policyTriggered, err := UnmarshalPolicyTriggered(bytes)
//    bytes, err = policyTriggered.Marshal()

package policies

import "encoding/json"

func UnmarshalPolicyTriggered(data []byte) (PolicyTriggered, error) {
	var r PolicyTriggered
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *PolicyTriggered) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Event data for triggered policies.
type PolicyTriggered struct {
	// Triggered policies for a system         
	Policies                          []Policy `json:"policies"`
	System                            System   `json:"system"`
}

type Policy struct {
	Condition   string `json:"condition"`
	Description string `json:"description"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
}

// A RHEL system managed by console.redhat.com
type System struct {
	// Timestamp of when the system did a check in. Must adhere to RFC 3339.                
	CheckIn                                                                 string          `json:"check_in"`
	DisplayName                                                             string          `json:"display_name"`
	Tags                                                                    []RHELSystemTag `json:"tags"`
	HostURL                                                                 *string         `json:"host_url,omitempty"`
	Hostname                                                                *string         `json:"hostname,omitempty"`
	InventoryID                                                             string          `json:"inventory_id"`
	RHELVersion                                                             *string         `json:"rhel_version,omitempty"`
}

type RHELSystemTag struct {
	Key       string  `json:"key"`
	Namespace string  `json:"namespace"`
	Value     *string `json:"value,omitempty"`
}
