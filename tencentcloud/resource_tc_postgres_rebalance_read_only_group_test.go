package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPostgresRebalanceReadOnlyGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresRebalanceReadOnlyGroup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_postgres_rebalance_read_only_group.rebalance_read_only_group", "id")),
			},
			{
				ResourceName:      "tencentcloud_postgres_rebalance_read_only_group.rebalance_read_only_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPostgresRebalanceReadOnlyGroup = `

resource "tencentcloud_postgres_rebalance_read_only_group" "rebalance_read_only_group" {
  read_only_group_id = "pgrogrp-test"
}

`
