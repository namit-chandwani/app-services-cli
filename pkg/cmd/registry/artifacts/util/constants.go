package util

import (
	"strings"
)

var DefaultArtifactGroup = "default"

var AllowedArtifactTypeEnumValues = []string{
	"AVRO",
	"PROTOBUF",
	"JSON",
	"OPENAPI",
	"ASYNCAPI",
	"GRAPHQL",
	"KCONNECT",
	"WSDL",
	"XSD",
	"XML",
}

// Get AllowedArtifactTypeEnumValues as string
func GetAllowedArtifactTypeEnumValuesAsString() string {
	return strings.Join(AllowedArtifactTypeEnumValues, ", ")
}
