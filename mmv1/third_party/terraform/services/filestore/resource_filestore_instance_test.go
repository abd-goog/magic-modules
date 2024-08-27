package filestore_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/services/filestore"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func testResourceFilestoreInstanceStateDataV0() map[string]interface{} {
	return map[string]interface{}{
		"zone": "us-central1-a",
	}
}

func testResourceFilestoreInstanceStateDataV1() map[string]interface{} {
	v0 := testResourceFilestoreInstanceStateDataV0()
	return map[string]interface{}{
		"location": v0["zone"],
		"zone":     v0["zone"],
	}
}

func TestFilestoreInstanceStateUpgradeV0(t *testing.T) {
	expected := testResourceFilestoreInstanceStateDataV1()
	// linter complains about nil context even in a test setting
	actual, err := filestore.ResourceFilestoreInstanceUpgradeV0(context.Background(), testResourceFilestoreInstanceStateDataV0(), nil)
	if err != nil {
		t.Fatalf("error migrating state: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("\n\nexpected:\n\n%#v\n\ngot:\n\n%#v\n\n", expected, actual)
	}
}

func TestAccFilestoreInstance_update(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckFilestoreInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFilestoreInstance_update(name),
			},
			{
				ResourceName:            "google_filestore_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone", "location", "labels", "terraform_labels"},
			},
			{
				Config: testAccFilestoreInstance_update2(name),
			},
			{
				ResourceName:            "google_filestore_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone", "location"},
			},
		},
	})
}

func testAccFilestoreInstance_update(name string) string {
	return fmt.Sprintf(`
resource "google_filestore_instance" "instance" {
  name        = "tf-instance-%s"
  zone        = "us-central1-b"
  tier        = "BASIC_HDD"
  description = "An instance created during testing."

  file_shares {
    capacity_gb = 1024
    name        = "share"
  }

  networks {
    network = "default"
    modes   = ["MODE_IPV4"]
  }

  labels = {
    baz = "qux"
  }
}
`, name)
}

func testAccFilestoreInstance_update2(name string) string {
	return fmt.Sprintf(`
resource "google_filestore_instance" "instance" {
  name        = "tf-instance-%s"
  zone        = "us-central1-b"
  tier        = "BASIC_HDD"
  description = "A modified instance created during testing."

  file_shares {
    capacity_gb = 1536
    name        = "share"
  }

  networks {
    network = "default"
    modes   = ["MODE_IPV4"]
  }
}
`, name)
}

func TestAccFilestoreInstance_reservedIpRange_update(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckFilestoreInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFilestoreInstance_reservedIpRange_update(name),
			},
			{
				ResourceName:            "google_filestore_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone", "location", "networks.0.reserved_ip_range"},
			},
			{
				Config: testAccFilestoreInstance_reservedIpRange_update2(name),
			},
			{
				ResourceName:            "google_filestore_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone", "location", "networks.0.reserved_ip_range"},
			},
		},
	})
}

func testAccFilestoreInstance_reservedIpRange_update(name string) string {
	return fmt.Sprintf(`
resource "google_filestore_instance" "instance" {
  name = "tf-instance-%s"
  zone = "us-central1-b"
  tier = "BASIC_HDD"

  file_shares {
    capacity_gb = 1024
    name        = "share1"
  }

  networks {
    network           = "default"
    modes             = ["MODE_IPV4"]
    reserved_ip_range = "172.19.30.0/29"
  }
}
`, name)
}

func testAccFilestoreInstance_reservedIpRange_update2(name string) string {
	return fmt.Sprintf(`
resource "google_filestore_instance" "instance" {
  name = "tf-instance-%s"
  zone = "us-central1-b"
  tier = "BASIC_HDD"

  file_shares {
    capacity_gb = 1024
    name        = "share1"
  }

  networks {
    network           = "default"
    modes             = ["MODE_IPV4"]
    reserved_ip_range = "172.19.31.0/29"
  }
}
`, name)
}


func TestAccFilestoreInstance_tags(t *testing.T) {
	t.Skip()

	t.Parallel()

	instance := fmt.Sprintf("tf-test%s.org1.com", acctest.RandString(t, 5))
	context := map[string]interface{}{
		"instance": instance,
		"resource_name": "instance",
	}

	resourceName := acctest.Nprintf("google_filestore_instance.%{resource_name}", context)
	org := envvar.GetTestOrgFromEnv(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckFilestoreInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFileInstanceTags(context, map[string]string{org + "/env": "test"}),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "name", "terraform_labels", "zone"},
			},
			{
				Config: testAccFileInstanceTags_allowDestroy(context, map[string]string{org + "/env": "test"}),
			},
		},
	})
}

func testAccFileInstanceTags(context map[string]interface{}, tags map[string]string) string {

	r := acctest.Nprintf(`
	resource "google_filestore_instance" "%{resource_name}" {
	  name = "tf-instance-%s"
          zone = "us-central1-b"
          tier = "BASIC_HDD"

          file_shares {
            capacity_gb = 1024
            name        = "share1"
          }

          networks {
            network           = "default"
            modes             = ["MODE_IPV4"]
            reserved_ip_range = "172.19.31.0/24"
          }
	  tags = {`, context)

	l := ""
	for key, value := range tags {
		l += fmt.Sprintf("%q = %q\n", key, value)
	}

	l += fmt.Sprintf("}\n}")
	return r + l
}

func testAccFileInstanceTags_allowDestroy(context map[string]interface{}, tags map[string]string) string {

	r := acctest.Nprintf(`
	resource "google_filestore_instance" "%{resource_name}" {
	  name = "tf-instance-%s"
          zone = "us-central1-b"
          tier = "BASIC_HDD"

          file_shares {
            capacity_gb = 1024
            name        = "share1"
          }

          networks {
            network           = "default"
            modes             = ["MODE_IPV4"]
            reserved_ip_range = "172.19.31.0/24"
          }
	  tags = {`, context)

	l := ""
	for key, value := range tags {
		l += fmt.Sprintf("%q = %q\n", key, value)
	}

	l += fmt.Sprintf("}\n}")
	return r + l
}

