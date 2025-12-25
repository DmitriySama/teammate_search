package swagger

import _ "embed"

//go:embed web.swagger.json
var tsSpec []byte

func Teammate_search() []byte {
	return tsSpec
}
