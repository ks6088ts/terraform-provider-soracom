package role

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceRole(t *testing.T) {
	r := ResourceRole()
	if r.Schema["operator_id"].Type != schema.TypeString {
		t.Error("Unexpected type")
	}
}
