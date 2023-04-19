// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    common, err := UnmarshalCommon(bytes)
//    bytes, err = common.Marshal()

package core

import "encoding/json"

type Common interface{}

func UnmarshalCommon(data []byte) (Common, error) {
	var r Common
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Common) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
