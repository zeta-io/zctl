package builder

import "github.com/zeta-io/zctl/api/schema"

type Builder struct {
	dir string
	template string

	schema *schema.Schema
}

