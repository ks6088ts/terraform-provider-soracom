package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestProvider(t *testing.T) {
	provider := Provider()
	if provider.Schema["profile"].Type != schema.TypeString {
		t.Error("unexpected type")
	}
}
