package sim

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceTags(t *testing.T) {
	r := ResourceTags()
	if r.Schema["tags"].Type != schema.TypeMap {
		t.Error("Unexpected type")
	}
}
