package role

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceRoleAttachment(t *testing.T) {
	r := ResourceRoleAttachment()
	if r.Schema["operator_id"].Type != schema.TypeString {
		t.Error("Unexpected type")
	}
}
