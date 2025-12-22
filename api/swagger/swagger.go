package swagger

import _ "embed"

//go:embed web.swagger.json
var registrySpec []byte

func Registry() []byte {
	return registrySpec
}
