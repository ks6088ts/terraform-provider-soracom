package eventhandler

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceHandlers(t *testing.T) {
	r := ResourceHandlers()
	if r.Schema["status"].Type != schema.TypeString {
		t.Error("Unexpected type")
	}
}
