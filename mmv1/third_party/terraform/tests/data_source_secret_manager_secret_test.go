package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceSecretManagerSecret_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerSecretDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecretManagerSecret_basic(context),
				Check: resource.ComposeTestCheckFunc(
					CheckDataSourceStateMatchesResourceState("data.google_secret_manager_secret.foo", "google_secret_manager_secret.bar"),
				),
			},
		},
	})
}

func testAccDataSourceSecretManagerSecret_basic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_secret_manager_secret" "bar" {
  secret_id = "tf-test-secret-%{random_suffix}"
  
  labels = {
    label = "my-label"
  }

  replication {
    user_managed {
      replicas {
        location = "us-central1"
      }
      replicas {
        location = "us-east1"
      }
    }
  }
}

data "google_secret_manager_secret" "foo" {
    secret_id = google_secret_manager_secret.bar.secret_id
}
`, context)
}
