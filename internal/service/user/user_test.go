package user

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceUser(t *testing.T) {
	r := ResourceUser()
	if r.Schema["operator_id"].Type != schema.TypeString {
		t.Error("Unexpected type")
	}
}
