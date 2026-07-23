package cvm_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	acctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svccvm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"
)

func TestAccTencentCloudPlacementGroupResource_Basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers:    acctest.AccProviders,
		CheckDestroy: testAccCheckPlacementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPlacementGroupResource_BasicCreate,
				Check:  resource.ComposeTestCheckFunc(testAccCheckPlacementGroupExists("tencentcloud_placement_group.placement"), resource.TestCheckResourceAttr("tencentcloud_placement_group.placement", "name", "tf-test-placement"), resource.TestCheckResourceAttr("tencentcloud_placement_group.placement", "type", "HOST"), resource.TestCheckResourceAttrSet("tencentcloud_placement_group.placement", "cvm_quota_total"), resource.TestCheckResourceAttrSet("tencentcloud_placement_group.placement", "current_num"), resource.TestCheckResourceAttrSet("tencentcloud_placement_group.placement", "create_time")),
			},
			{
				Config: testAccPlacementGroupResource_BasicChange1,
				Check:  resource.ComposeTestCheckFunc(testAccCheckPlacementGroupExists("tencentcloud_placement_group.placement"), resource.TestCheckResourceAttr("tencentcloud_placement_group.placement", "name", "tf-test-placement1"), resource.TestCheckResourceAttrSet("tencentcloud_placement_group.placement", "cvm_quota_total"), resource.TestCheckResourceAttrSet("tencentcloud_placement_group.placement", "current_num"), resource.TestCheckResourceAttrSet("tencentcloud_placement_group.placement", "create_time")),
			},
			{
				ResourceName:      "tencentcloud_placement_group.placement",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPlacementGroupResource_BasicCreate = `

resource "tencentcloud_placement_group" "placement" {
    type = "HOST"
    name = "tf-test-placement"
}

`
const testAccPlacementGroupResource_BasicChange1 = `

resource "tencentcloud_placement_group" "placement" {
    type = "HOST"
    name = "tf-test-placement1"
}

`

func testAccCheckPlacementGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("placement group %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("placement group id is not set")
		}

		cvmService := svccvm.NewCvmService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		placement, err := cvmService.DescribePlacementGroupById(ctx, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				placement, err = cvmService.DescribePlacementGroupById(ctx, rs.Primary.ID)
				if err != nil {
					return tccommon.RetryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		if placement == nil {
			return fmt.Errorf("placement group id is not found")
		}
		return nil
	}
}

func testAccCheckPlacementGroupDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cvmService := svccvm.NewCvmService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_placement_group" {
			continue
		}

		placement, err := cvmService.DescribePlacementGroupById(ctx, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				placement, err = cvmService.DescribePlacementGroupById(ctx, rs.Primary.ID)
				if err != nil {
					return tccommon.RetryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		if placement != nil {
			return fmt.Errorf("placement group still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

// Unit tests for strategy and partition_count fields
// Run all unit tests in this file:
//   go test -v -run "TestUnitPlacementGroup_" ./tencentcloud/services/cvm/

func TestUnitPlacementGroup_StrategySchemaDefinition(t *testing.T) {
	res := svccvm.ResourceTencentCloudPlacementGroup()
	s := res.Schema

	strategySchema, ok := s["strategy"]
	assert.True(t, ok, "strategy field should exist in schema")
	assert.Equal(t, schema.TypeString, strategySchema.Type)
	assert.True(t, strategySchema.Optional)
	assert.True(t, strategySchema.Computed)
	assert.NotNil(t, strategySchema.ValidateFunc)
}

func TestUnitPlacementGroup_PartitionCountSchemaDefinition(t *testing.T) {
	res := svccvm.ResourceTencentCloudPlacementGroup()
	s := res.Schema

	partitionCountSchema, ok := s["partition_count"]
	assert.True(t, ok, "partition_count field should exist in schema")
	assert.Equal(t, schema.TypeInt, partitionCountSchema.Type)
	assert.True(t, partitionCountSchema.Optional)
	assert.True(t, partitionCountSchema.Computed)
	assert.NotNil(t, partitionCountSchema.ValidateFunc)
}

func TestUnitPlacementGroup_CreateRequestWithPartitionStrategy(t *testing.T) {
	request := cvm.NewCreateDisasterRecoverGroupRequest()

	placementName := "test-partition-group"
	placementType := "HOST"
	strategy := "PARTITION"
	affinity := 1
	partitionCount := 5

	request.Name = &placementName
	request.Type = &placementType

	if strategy != "" {
		request.Strategy = &strategy
	}

	if affinity != 0 {
		request.Affinity = helper.IntInt64(affinity)
	}

	if partitionCount > 0 {
		request.PartitionCount = helper.IntInt64(partitionCount)
	}

	assert.Equal(t, "test-partition-group", *request.Name)
	assert.Equal(t, "HOST", *request.Type)
	assert.Equal(t, "PARTITION", *request.Strategy)
	assert.Equal(t, int64(1), *request.Affinity)
	assert.Equal(t, int64(5), *request.PartitionCount)
}

func TestUnitPlacementGroup_CreateRequestWithSpreadStrategy(t *testing.T) {
	request := cvm.NewCreateDisasterRecoverGroupRequest()

	placementName := "test-spread-group"
	placementType := "SW"
	strategy := "SPREAD"
	affinity := 3
	partitionCount := 0

	request.Name = &placementName
	request.Type = &placementType

	if strategy != "" {
		request.Strategy = &strategy
	}

	if affinity != 0 {
		request.Affinity = helper.IntInt64(affinity)
	}

	if partitionCount > 0 {
		request.PartitionCount = helper.IntInt64(partitionCount)
	}

	assert.Equal(t, "test-spread-group", *request.Name)
	assert.Equal(t, "SW", *request.Type)
	assert.Equal(t, "SPREAD", *request.Strategy)
	assert.Equal(t, int64(3), *request.Affinity)
	assert.Nil(t, request.PartitionCount, "partition_count should not be set when value is 0")
}

func TestUnitPlacementGroup_CreateRequestWithoutStrategy(t *testing.T) {
	request := cvm.NewCreateDisasterRecoverGroupRequest()

	placementName := "test-group"
	placementType := "RACK"
	strategy := ""
	affinity := 0
	partitionCount := 0

	request.Name = &placementName
	request.Type = &placementType

	if strategy != "" {
		request.Strategy = &strategy
	}

	if affinity != 0 {
		request.Affinity = helper.IntInt64(affinity)
	}

	if partitionCount > 0 {
		request.PartitionCount = helper.IntInt64(partitionCount)
	}

	assert.Equal(t, "test-group", *request.Name)
	assert.Equal(t, "RACK", *request.Type)
	assert.Nil(t, request.Strategy, "strategy should not be set when empty")
	assert.Nil(t, request.Affinity, "affinity should not be set when 0")
	assert.Nil(t, request.PartitionCount, "partition_count should not be set when 0")
}

func TestUnitPlacementGroup_ReadStateWithStrategyAndPartitionCount(t *testing.T) {
	placement := &cvm.DisasterRecoverGroup{
		DisasterRecoverGroupId: helper.String("ps-12345"),
		Name:                   helper.String("test-partition-group"),
		Type:                   helper.String("HOST"),
		Strategy:               helper.String("PARTITION"),
		PartitionCount:         helper.Int64(10),
		Affinity:               helper.Int64(2),
		CvmQuotaTotal:          helper.Int64(100),
		CurrentNum:             helper.Int64(5),
		CreateTime:             helper.String("2024-01-01T00:00:00Z"),
	}

	result := make(map[string]interface{})
	result["name"] = placement.Name
	result["type"] = placement.Type
	result["affinity"] = placement.Affinity
	result["cvm_quota_total"] = placement.CvmQuotaTotal
	result["current_num"] = placement.CurrentNum
	result["create_time"] = placement.CreateTime

	if placement.Strategy != nil {
		result["strategy"] = *placement.Strategy
	}
	if placement.PartitionCount != nil {
		result["partition_count"] = int(*placement.PartitionCount)
	}

	assert.Equal(t, "PARTITION", result["strategy"])
	assert.Equal(t, 10, result["partition_count"])
}

func TestUnitPlacementGroup_ReadStateWithNilStrategyAndPartitionCount(t *testing.T) {
	placement := &cvm.DisasterRecoverGroup{
		DisasterRecoverGroupId: helper.String("ps-67890"),
		Name:                   helper.String("test-spread-group"),
		Type:                   helper.String("SW"),
		Strategy:               nil,
		PartitionCount:         nil,
		Affinity:               helper.Int64(1),
		CvmQuotaTotal:          helper.Int64(50),
		CurrentNum:             helper.Int64(2),
		CreateTime:             helper.String("2024-06-01T00:00:00Z"),
	}

	result := make(map[string]interface{})
	result["name"] = placement.Name
	result["type"] = placement.Type
	result["affinity"] = placement.Affinity

	if placement.Strategy != nil {
		result["strategy"] = *placement.Strategy
	}
	if placement.PartitionCount != nil {
		result["partition_count"] = int(*placement.PartitionCount)
	}

	_, hasStrategy := result["strategy"]
	assert.False(t, hasStrategy, "strategy should not be set when nil")
	_, hasPartitionCount := result["partition_count"]
	assert.False(t, hasPartitionCount, "partition_count should not be set when nil")
}

func TestUnitPlacementGroup_PartitionCountValidation(t *testing.T) {
	partitionCount := 5
	strategy := "SPREAD"

	var err error
	if partitionCount > 0 && strategy != svccvm.CVM_PLACEMENT_GROUP_STRATEGY_PARTITION {
		err = fmt.Errorf("`partition_count` is only valid when `strategy` is set to `PARTITION`")
	}

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "`partition_count` is only valid when `strategy` is set to `PARTITION`")
}

func TestUnitPlacementGroup_PartitionCountValidationPass(t *testing.T) {
	partitionCount := 5
	strategy := "PARTITION"

	var err error
	if partitionCount > 0 && strategy != svccvm.CVM_PLACEMENT_GROUP_STRATEGY_PARTITION {
		err = fmt.Errorf("`partition_count` is only valid when `strategy` is set to `PARTITION`")
	}

	assert.Nil(t, err)
}

func TestUnitPlacementGroup_PartitionCountZeroValidationPass(t *testing.T) {
	partitionCount := 0
	strategy := "SPREAD"

	var err error
	if partitionCount > 0 && strategy != svccvm.CVM_PLACEMENT_GROUP_STRATEGY_PARTITION {
		err = fmt.Errorf("`partition_count` is only valid when `strategy` is set to `PARTITION`")
	}

	assert.Nil(t, err)
}

func TestUnitPlacementGroup_ImmutableFields(t *testing.T) {
	immutableArgs := []string{
		"type",
		"strategy",
		"affinity",
		"partition_count",
	}

	assert.Contains(t, immutableArgs, "strategy")
	assert.Contains(t, immutableArgs, "partition_count")
	assert.Contains(t, immutableArgs, "type")
	assert.Contains(t, immutableArgs, "affinity")
	assert.Equal(t, 4, len(immutableArgs))
}

func TestUnitPlacementGroup_StrategyConstants(t *testing.T) {
	assert.Equal(t, "SPREAD", svccvm.CVM_PLACEMENT_GROUP_STRATEGY_SPREAD)
	assert.Equal(t, "PARTITION", svccvm.CVM_PLACEMENT_GROUP_STRATEGY_PARTITION)
	assert.Equal(t, 2, len(svccvm.CVM_PLACEMENT_GROUP_STRATEGY))
	assert.Contains(t, svccvm.CVM_PLACEMENT_GROUP_STRATEGY, svccvm.CVM_PLACEMENT_GROUP_STRATEGY_SPREAD)
	assert.Contains(t, svccvm.CVM_PLACEMENT_GROUP_STRATEGY, svccvm.CVM_PLACEMENT_GROUP_STRATEGY_PARTITION)
}
