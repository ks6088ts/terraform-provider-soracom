package sim

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceConfigurationGroup(t *testing.T) {
	r := ResourceConfigurationGroup()
	if r.Schema["sim_id"].Type != schema.TypeString {
		t.Error("Unexpected type")
	}
}
