package tencentcloud

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudKubernetesEncryptionProtectionResource_basic(t *testing.T) {
	t.Parallel()
	rName := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				SkipFunc: func() (bool, error) {
					if strings.Contains(os.Getenv(PROVIDER_DOMAIN), "test") {
						fmt.Printf("[International station]skip TestAccTencentCloudKubernetesClusterEndpointResource, because the test station did not support this feature yet!\n")
						return true, nil
					}
					return false, errors.New("need test")
				},
				Config: fmt.Sprintf(testAccTkeEncryptionProtection, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_encryption_protection.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_encryption_protection.example", "cluster_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_encryption_protection.example", "kms_configuration.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_encryption_protection.example", "kms_configuration.0.key_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_encryption_protection.example", "kms_configuration.0.kms_region"),
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_encryption_protection.example", "status"),
				),
			},
		},
	})
}

const testAccTkeEncryptionProtection = `

variable "example_region" {
  default = "ap-guangzhou"
}

variable "example_cluster_cidr" {
  default = "10.32.0.0/16"
}

variable "availability_zone" {
  default = "ap-guangzhou-3"
}

variable "env_az" {
  type = string
}

variable "env_region" {
  type = string
}

data "tencentcloud_vpc_subnets" "vpc" {
  is_default        = true
  availability_zone = var.env_az != "" ? var.env_az : var.availability_zone
}

resource "tencentcloud_kubernetes_cluster" "example" {
  vpc_id                  = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  cluster_cidr            = var.example_cluster_cidr
  cluster_max_pod_num     = 32
  cluster_name            = "tf_example_cluster"
  cluster_desc            = "a tf example cluster for the kms test"
  cluster_max_service_num = 32
  cluster_deploy_type     = "MANAGED_CLUSTER"
}

resource "tencentcloud_kms_key" "example" {
  alias       = "tf-example-%s"
  description = "example of kms key instance"
  key_usage   = "ENCRYPT_DECRYPT"
  is_enabled  = true
}

resource "tencentcloud_kubernetes_encryption_protection" "example" {
  cluster_id = tencentcloud_kubernetes_cluster.example.id
  kms_configuration {
    key_id     = tencentcloud_kms_key.example.id
    kms_region = var.env_region != "" ? var.env_region : var.example_region
  }
}

`
