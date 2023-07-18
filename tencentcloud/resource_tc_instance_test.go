package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("tencentcloud_instance", &resource.Sweeper{
		Name: "tencentcloud_instance",
		F:    testSweepCvmInstance,
	})
}

func testSweepCvmInstance(region string) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	sharedClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("getting tencentcloud client error: %s", err.Error())
	}
	client := sharedClient.(*TencentCloudClient)

	cvmService := CvmService{
		client: client.apiV3Conn,
	}

	instances, err := cvmService.DescribeInstanceByFilter(ctx, nil, nil)
	if err != nil {
		return fmt.Errorf("get instance list error: %s", err.Error())
	}

	for _, v := range instances {
		instanceId := *v.InstanceId
		instanceName := *v.InstanceName
		now := time.Now()
		createTime := stringTotime(*v.CreatedTime)
		interval := now.Sub(createTime).Minutes()

		if strings.HasPrefix(instanceName, keepResource) || strings.HasPrefix(instanceName, defaultResource) {
			continue
		}

		if needProtect == 1 && int64(interval) < 30 {
			continue
		}

		if err = cvmService.DeleteInstance(ctx, instanceId); err != nil {
			log.Printf("[ERROR] sweep instance %s error: %s", instanceId, err.Error())
		}
	}

	return nil
}

func TestAccTencentCloudInstanceResource_Basic(t *testing.T) {
	t.Parallel()

	id := "tencentcloud_instance.cvm_basic"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: id,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_COMMON) },
				Config:    testAccTencentCloudInstanceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttrSet(id, "private_ip"),
					resource.TestCheckResourceAttrSet(id, "vpc_id"),
					resource.TestCheckResourceAttrSet(id, "subnet_id"),
					resource.TestCheckResourceAttrSet(id, "project_id"),
				),
			},
			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_COMMON) },
				Config:    testAccTencentCloudInstanceModifyInstanceType,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttrSet(id, "instance_type"),
				),
			},
			{
				ResourceName:            id,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"disable_monitor_service", "disable_security_service", "hostname", "password", "force_delete"},
			},
		},
	})
}

func TestAccTencentCloudInstanceResource_WithDataDisk(t *testing.T) {
	t.Parallel()

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: id,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_COMMON) },
				Config:    testAccTencentCloudInstanceWithDataDisk,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr(id, "system_disk_size", "100"),
					resource.TestCheckResourceAttr(id, "system_disk_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr(id, "data_disks.0.data_disk_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr(id, "data_disks.0.data_disk_size", "100"),
					resource.TestCheckResourceAttr(id, "data_disks.0.data_disk_snapshot_id", ""),
					resource.TestCheckResourceAttr(id, "data_disks.1.data_disk_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr(id, "data_disks.1.data_disk_size", "100"),
				),
			},
			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_COMMON) },
				Config:    testAccTencentCloudInstanceWithDataDiskUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr(id, "system_disk_size", "100"),
					resource.TestCheckResourceAttr(id, "system_disk_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr(id, "data_disks.0.data_disk_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr(id, "data_disks.0.data_disk_size", "150"),
					resource.TestCheckResourceAttr(id, "data_disks.0.data_disk_snapshot_id", ""),
					resource.TestCheckResourceAttr(id, "data_disks.1.data_disk_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr(id, "data_disks.1.data_disk_size", "150"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceResource_WithNetwork(t *testing.T) {
	t.Parallel()

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: id,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_COMMON) },
				Config:    testAccTencentCloudInstanceWithNetworkFalse("false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckNoResourceAttr(id, "public_ip"),
				),
			},
			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_COMMON) },
				Config:    testAccTencentCloudInstanceWithNetwork("true", 5),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "internet_max_bandwidth_out", "5"),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttrSet(id, "public_ip"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceResource_WithPrivateIP(t *testing.T) {
	t.Parallel()
	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: id,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_COMMON) },
				Config:    testAccTencentCloudInstanceWithPrivateIP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceResource_WithKeyPairs(t *testing.T) {
	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: id,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_COMMON) },
				Config:    testAccTencentCloudInstanceWithKeyPair_withoutKeyPair,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
				),
			},
			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_COMMON) },
				Config: testAccTencentCloudInstanceWithKeyPair(
					"[tencentcloud_key_pair.key_pair_0.id, tencentcloud_key_pair.key_pair_1.id]",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr(id, "key_ids.#", "2"),
				),
			},
			{
				PreConfig: func() {
					testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_COMMON)
					time.Sleep(time.Second * 5)
				},
				Config: testAccTencentCloudInstanceWithKeyPair("[tencentcloud_key_pair.key_pair_2.id]"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr(id, "key_ids.#", "1"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceResource_WithPassword(t *testing.T) {
	t.Parallel()

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: id,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_COMMON) },
				Config:    testAccTencentCloudInstanceWithPassword("TF_test_123"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttrSet(id, "password"),
				),
			},
			{
				PreConfig: func() {
					testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_COMMON)
					time.Sleep(time.Second * 5)
				},
				Config: testAccTencentCloudInstanceWithPassword("TF_test_123456"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttrSet(id, "password"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceResource_WithImageLogin(t *testing.T) {

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: id,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_PREPAY) },
				Config:    testAccTencentCloudInstanceWithImageLogin,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr(id, "keep_image_login", "true"),
					resource.TestCheckResourceAttr(id, "disable_api_termination", "false"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceResource_WithName(t *testing.T) {
	t.Parallel()

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: id,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_COMMON) },
				Config:    testAccTencentCloudInstanceWithName(defaultInsName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr(id, "instance_name", defaultInsName),
				),
			},
			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_COMMON) },
				Config:    testAccTencentCloudInstanceWithName(defaultInsNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr(id, "instance_name", defaultInsNameUpdate),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceResource_WithHostname(t *testing.T) {
	t.Parallel()

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: id,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_COMMON) },
				Config:    testAccTencentCloudInstanceWithHostname,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr(id, "hostname", defaultInsName),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceResource_WithSecurityGroup(t *testing.T) {
	t.Parallel()

	instanceId := "tencentcloud_instance.foo"
	securitygroupId := "tencentcloud_security_group.foo"
	securitygroupRuleFooId := "tencentcloud_security_group_rule.foo"
	securitygroupRuleBarId := "tencentcloud_security_group_rule.bar"

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: instanceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_COMMON) },
				Config:    testAccTencentCloudInstanceWithSecurityGroup(`[tencentcloud_security_group.foo.id]`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(instanceId),
					testAccCheckTencentCloudInstanceExists(instanceId),
					resource.TestCheckResourceAttr(instanceId, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr(instanceId, "security_groups.#", "1"),
					resource.TestCheckResourceAttrSet(securitygroupId, "id"),
					resource.TestCheckResourceAttr(securitygroupRuleFooId, "type", "ingress"),
					resource.TestCheckResourceAttr(securitygroupRuleFooId, "port_range", "80,8080"),
					resource.TestCheckResourceAttr(securitygroupRuleBarId, "type", "ingress"),
					resource.TestCheckResourceAttr(securitygroupRuleBarId, "port_range", "3000"),
				),
			},
			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_COMMON) },
				Config: testAccTencentCloudInstanceWithSecurityGroup(`[
					tencentcloud_security_group.foo.id,
					tencentcloud_security_group.bar.id
				]`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(instanceId),
					testAccCheckTencentCloudInstanceExists(instanceId),
					resource.TestCheckResourceAttr(instanceId, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr(instanceId, "security_groups.#", "2"),
					resource.TestCheckResourceAttrSet(securitygroupId, "id"),
					resource.TestCheckResourceAttr(securitygroupRuleFooId, "type", "ingress"),
					resource.TestCheckResourceAttr(securitygroupRuleFooId, "port_range", "80,8080"),
					resource.TestCheckResourceAttr(securitygroupRuleBarId, "type", "ingress"),
					resource.TestCheckResourceAttr(securitygroupRuleBarId, "port_range", "3000"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceResource_WithOrderlySecurityGroup(t *testing.T) {
	t.Parallel()

	var sgId1, sgId2, sgId3 string
	instanceId := "tencentcloud_instance.cvm_with_orderly_sg"
	orderlySecurityGroupId1 := "tencentcloud_security_group.orderly_security_group1"
	orderlySecurityGroupId2 := "tencentcloud_security_group.orderly_security_group2"
	orderlySecurityGroupId3 := "tencentcloud_security_group.orderly_security_group3"

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: instanceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_COMMON) },
				Config: testAccTencentCloudInstanceOrderlySecurityGroups(`[
					tencentcloud_security_group.orderly_security_group1.id,
					tencentcloud_security_group.orderly_security_group2.id,
					tencentcloud_security_group.orderly_security_group3.id
				]`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudInstanceExists(instanceId),
					testAccCheckSecurityGroupExists(orderlySecurityGroupId1, &sgId1),
					testAccCheckSecurityGroupExists(orderlySecurityGroupId2, &sgId2),
					testAccCheckSecurityGroupExists(orderlySecurityGroupId3, &sgId3),

					resource.TestCheckResourceAttrPtr(instanceId, "orderly_security_groups.0", &sgId1),
					resource.TestCheckResourceAttrPtr(instanceId, "orderly_security_groups.1", &sgId2),
					resource.TestCheckResourceAttrPtr(instanceId, "orderly_security_groups.2", &sgId3),
				),
			},

			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_COMMON) },
				Config: testAccTencentCloudInstanceOrderlySecurityGroups(`[
					tencentcloud_security_group.orderly_security_group3.id,
					tencentcloud_security_group.orderly_security_group2.id,
					tencentcloud_security_group.orderly_security_group1.id
				]`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudInstanceExists(instanceId),
					testAccCheckSecurityGroupExists(orderlySecurityGroupId1, &sgId1),
					testAccCheckSecurityGroupExists(orderlySecurityGroupId2, &sgId2),
					testAccCheckSecurityGroupExists(orderlySecurityGroupId3, &sgId3),

					resource.TestCheckResourceAttrPtr(instanceId, "orderly_security_groups.0", &sgId3),
					resource.TestCheckResourceAttrPtr(instanceId, "orderly_security_groups.1", &sgId2),
					resource.TestCheckResourceAttrPtr(instanceId, "orderly_security_groups.2", &sgId1),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceResource_WithTags(t *testing.T) {
	t.Parallel()

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: id,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_COMMON) },
				Config: testAccTencentCloudInstanceWithTags(`{
					"hello" = "world"
					"happy" = "hour"
				}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr(id, "tags.hello", "world"),
					resource.TestCheckResourceAttr(id, "tags.happy", "hour"),
				),
			},
			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_COMMON) },
				Config: testAccTencentCloudInstanceWithTags(`{
					"hello" = "hello"
				}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr(id, "tags.hello", "hello"),
					resource.TestCheckNoResourceAttr(id, "tags.happy"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceResource_WithPlacementGroup(t *testing.T) {
	t.Parallel()

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_COMMON) },
				Config:    testAccTencentCloudInstanceWithPlacementGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttrSet(id, "placement_group_id"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceResource_WithSpotpaid(t *testing.T) {
	t.Parallel()

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: id,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_COMMON) },
				Config:    testAccTencentCloudInstanceWithSpotpaid,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceResource_DataDiskOrder(t *testing.T) {
	t.Parallel()

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: id,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_COMMON) },
				Config:    testAccTencentCloudInstanceWithDataDiskOrder,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "data_disks.0.data_disk_size", "100"),
					resource.TestCheckResourceAttr(id, "data_disks.1.data_disk_size", "50"),
					resource.TestCheckResourceAttr(id, "data_disks.2.data_disk_size", "70"),
				),
			},
		},
	})
}

func TestAccTencentCloudNeedFixInstancePostpaidToPrepaid(t *testing.T) {

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: id,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_PREPAY) },
				Config:    testAccTencentCloudInstancePostPaid,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
				),
			},
			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_PREPAY) },
				Config:    testAccTencentCloudInstanceBasicToPrepaid,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_charge_type", "PREPAID"),
					resource.TestCheckResourceAttr(id, "instance_charge_type_prepaid_period", "1"),
					resource.TestCheckResourceAttr(id, "instance_charge_type_prepaid_renew_flag", "NOTIFY_AND_MANUAL_RENEW"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceResource_PrepaidFallbackToPostpaid(t *testing.T) {

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: id,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_PREPAY) },
				Config:    testAccTencentCloudInstanceBasicToPrepaid,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_charge_type", "PREPAID"),
					resource.TestCheckResourceAttr(id, "instance_charge_type_prepaid_period", "1"),
					resource.TestCheckResourceAttr(id, "instance_charge_type_prepaid_renew_flag", "NOTIFY_AND_MANUAL_RENEW"),
				),
			},
			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_PREPAY) },
				Config:    testAccTencentCloudInstancePostPaid,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
				),
			},
		},
	})
}

func testAccCheckTencentCloudInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cvm instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cvm instance id is not set")
		}

		cvmService := CvmService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		instance, err := cvmService.DescribeInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				instance, err = cvmService.DescribeInstanceById(ctx, rs.Primary.ID)
				if err != nil {
					return retryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		if instance == nil {
			return fmt.Errorf("cvm instance id is not found")
		}
		return nil
	}
}

func testAccCheckInstanceDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cvmService := CvmService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_instance" {
			continue
		}

		instance, err := cvmService.DescribeInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				instance, err = cvmService.DescribeInstanceById(ctx, rs.Primary.ID)
				if err != nil {
					return retryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		if instance != nil && *instance.InstanceState != CVM_STATUS_SHUTDOWN && *instance.InstanceState != CVM_STATUS_TERMINATING {
			return fmt.Errorf("cvm instance still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

const testAccTencentCloudInstanceBasic = defaultInstanceVariable + `
resource "tencentcloud_instance" "cvm_basic" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  vpc_id            = var.cvm_vpc_id
  subnet_id         = var.cvm_subnet_id
  system_disk_type  = "CLOUD_PREMIUM"
  project_id        = 0
}
`

const testAccTencentCloudInstanceWithDataDiskOrder = defaultInstanceVariable + `
resource "tencentcloud_instance" "foo" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  vpc_id            = var.cvm_vpc_id
  subnet_id         = var.cvm_subnet_id
  system_disk_type  = "CLOUD_PREMIUM"
  project_id        = 0

  data_disks {
    data_disk_size         = 100
    data_disk_type         = "CLOUD_PREMIUM"
    delete_with_instance   = true
  }
  data_disks {
    data_disk_size         = 50
    data_disk_type         = "CLOUD_PREMIUM"
    delete_with_instance   = true
  }
  data_disks {
    data_disk_size         = 70
    data_disk_type         = "CLOUD_PREMIUM"
    delete_with_instance   = true
  }
}
`

const testAccTencentCloudInstancePostPaid = `
data "tencentcloud_instance_types" "default" {
  filter {
    name   = "instance-family"
    values = ["S1"]
  }

  cpu_core_count = 2
  memory_size    = 2
}

resource "tencentcloud_instance" "foo" {
  instance_name     = "` + defaultInsName + `"
  availability_zone = "` + defaultAZone + `"
  image_id          = "` + defaultTkeOSImageId + `"
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type  = "CLOUD_PREMIUM"
  force_delete = true
}
`

const testAccTencentCloudInstanceBasicToPrepaid = `
data "tencentcloud_instance_types" "default" {
  filter {
    name   = "instance-family"
    values = ["S1"]
  }

  cpu_core_count = 2
  memory_size    = 2
}

resource "tencentcloud_instance" "foo" {
  instance_name     = "` + defaultInsName + `"
  availability_zone = "` + defaultAZone + `"
  image_id          = "` + defaultTkeOSImageId + `"
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type  = "CLOUD_PREMIUM"
  instance_charge_type       = "PREPAID"
  instance_charge_type_prepaid_period = 1
  instance_charge_type_prepaid_renew_flag = "NOTIFY_AND_MANUAL_RENEW"
  force_delete = true
}
`

const testAccTencentCloudInstanceModifyInstanceType = defaultInstanceVariable + `
data "tencentcloud_instance_types" "new_type" {
	availability_zone = var.availability_cvm_zone
  
	cpu_core_count = 2
	memory_size    = 2
  }

resource "tencentcloud_instance" "cvm_basic" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.new_type.instance_types.0.instance_type
  vpc_id            = var.cvm_vpc_id
  subnet_id         = var.cvm_subnet_id
  system_disk_type  = "CLOUD_PREMIUM"
  project_id        = 0
}
`

const testAccTencentCloudInstanceWithDataDisk = defaultInstanceVariable + `
resource "tencentcloud_instance" "foo" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type

  system_disk_type = "CLOUD_PREMIUM"
  system_disk_size = 100

  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 100
    delete_with_instance  = true
	// encrypt = true
  } 
   
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 100
    # data_disk_snapshot_id = "snap-nvzu3dmh"
    delete_with_instance  = true
  }

  disable_security_service = true
  disable_monitor_service  = true
}
`

const testAccTencentCloudInstanceWithDataDiskUpdate = defaultInstanceVariable + `
resource "tencentcloud_instance" "foo" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type

  system_disk_type = "CLOUD_PREMIUM"
  system_disk_size = 100

  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  } 
   
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }

  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }

  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }
  
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }



  disable_security_service = true
  disable_monitor_service  = true
}
`

func testAccTencentCloudInstanceWithNetworkFalse(hasPublicIp string) string {
	return fmt.Sprintf(
		defaultInstanceVariable+`
resource "tencentcloud_instance" "foo" {
  instance_name              = var.instance_name
  availability_zone          = var.availability_cvm_zone
  image_id                   = data.tencentcloud_images.default.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  allocate_public_ip         = %s
  system_disk_type           = "CLOUD_PREMIUM"
}
`,
		hasPublicIp,
	)
}

func testAccTencentCloudInstanceWithNetwork(hasPublicIp string, maxBandWidthOut int64) string {
	return fmt.Sprintf(
		defaultInstanceVariable+`
resource "tencentcloud_instance" "foo" {
  instance_name              = var.instance_name
  availability_zone          = var.availability_cvm_zone
  image_id                   = data.tencentcloud_images.default.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  internet_max_bandwidth_out = %d
  allocate_public_ip         = %s
  system_disk_type           = "CLOUD_PREMIUM"
}
`,
		maxBandWidthOut, hasPublicIp,
	)
}

const testAccTencentCloudInstanceWithPrivateIP = defaultInstanceVariable + `
resource "tencentcloud_instance" "foo" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type  = "CLOUD_PREMIUM"
  vpc_id            = var.cvm_vpc_id
  subnet_id         = var.cvm_subnet_id
  private_ip        = "10.0.0.123"
}
`

const testAccTencentCloudInstanceWithKeyPair_withoutKeyPair = defaultInstanceVariable + `
resource "tencentcloud_instance" "foo" {
	instance_name     = var.instance_name
	availability_zone = var.availability_cvm_zone
	image_id          = data.tencentcloud_images.default.images.0.image_id
	instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
	system_disk_type  = "CLOUD_PREMIUM"
}
`

func testAccTencentCloudInstanceWithKeyPair(keyIds string) string {

	return fmt.Sprintf(
		defaultInstanceVariable+`
resource "tencentcloud_key_pair" "key_pair_0" {
  key_name = "key_pair_0"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDjd8fTnp7Dcuj4mLaQxf9Zs/ORgUL9fQxRCNKkPgP1paTy1I513maMX126i36Lxxl3+FUB52oVbo/FgwlIfX8hyCnv8MCxqnuSDozf1CD0/wRYHcTWAtgHQHBPCC2nJtod6cVC3kB18KeV4U7zsxmwFeBIxojMOOmcOBuh7+trRw=="
}

resource "tencentcloud_key_pair" "key_pair_1" {
  key_name = "key_pair_1"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCzwYE6KI8uULEvSNA2k1tlsLtMDe+x1Saw6yL3V1mk9NFws0K2BshYqsnP/BlYiGZv/Nld5xmGoA9LupOcUpyyGGSHZdBrMx1Dz9ajewe7kGowRWwwMAHTlzh9+iqeg/v6P5vW6EwK4hpGWgv06vGs3a8CzfbHu1YRbZAO/ysp3ymdL+vGvw/vzC0T+YwPMisn9wFD5FTlJ+Em6s9PzxqR/41t4YssmCwUV78ZoYL8CyB0emuB8wALvcXbdUVxMxpBEHd5U6ZP5+HPxU2WFbWqiFCuErLIZRuxFw8L/Ot+JOyNnadN1XU4crYDX5cML1i/ExXKVIDoBaLtgAJOpyeP"
}

resource "tencentcloud_key_pair" "key_pair_2" {
  key_name = "key_pair_2"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDJ1zyoM55pKxJptZBKceZSEypPN7BOunqBR1Qj3Tz5uImJ+dwfKzggu8PGcbHtuN8D2n1BH/GDkiGFaz/sIYUJWWZudcdut+ra32MqUvk953Sztf12rsFC1+lZ1CYEgon8Lt6ehxn+61tsS31yfUmpL1mq2vuca7J0NLdPMpxIYkGlifyAMISMmxi/m7gPYpbdZTmhQQS2aOhuLm+B4MwtTvT58jqNzIaFU0h5sqAvGQfzI5pcxwYvFTeQeXjJZfaYapDHN0MAg0b/vIWWNrDLv7dlv//OKBIaL0LIzIGQS8XXhF3HlyqfDuf3bjLBIKzYGSV/DRqlEsGBgzinJZXvJZug5oq1n2njDFsdXEvL6fYsP4WLvBLiQlceQ7oXi7m5nfrwFTaX+mpo7dUOR9AcyQ1AAgCcM67orB4E33ycaArGHtpjnCnWUjqQ+yCj4EXsD4yOL77wGsmhkbboVNnYAD9MJWsFP03hZE7p/RHY0C5NfLPT3mL45oZxBpC5mis="
}

resource "tencentcloud_instance" "foo" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  key_ids           = %s
  system_disk_type  = "CLOUD_PREMIUM"
}
`,
		keyIds,
	)
}

func testAccTencentCloudInstanceWithPassword(password string) string {
	return fmt.Sprintf(
		defaultInstanceVariable+`
resource "tencentcloud_instance" "foo" {
  instance_name              = var.instance_name
  availability_zone          = var.availability_cvm_zone
  image_id                   = data.tencentcloud_images.default.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  password                   = "%s"
  system_disk_type           = "CLOUD_PREMIUM"
}
`,
		password,
	)
}

const testAccTencentCloudInstanceWithImageLogin = defaultInstanceVariable + `
data "tencentcloud_images" "zoo" {
  image_type = ["PRIVATE_IMAGE"]
}
resource "tencentcloud_instance" "foo" {
  instance_name              = var.instance_name
  availability_zone          = var.availability_cvm_zone
  image_id                   = data.tencentcloud_images.zoo.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  keep_image_login 			 = true
  system_disk_type           = "CLOUD_PREMIUM"
  disable_api_termination    = false
}
`

func testAccTencentCloudInstanceWithName(instanceName string) string {
	return fmt.Sprintf(
		defaultInstanceVariable+`
resource "tencentcloud_instance" "foo" {
  instance_name     = "%s"
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type  = "CLOUD_PREMIUM"
}
`,
		instanceName,
	)
}

const testAccTencentCloudInstanceWithHostname = defaultInstanceVariable + `
resource "tencentcloud_instance" "foo" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  hostname          = var.instance_name
  system_disk_type  = "CLOUD_PREMIUM"
}
`

func testAccTencentCloudInstanceWithSecurityGroup(ids string) string {
	return fmt.Sprintf(
		defaultInstanceVariable+`
resource "tencentcloud_security_group" "foo" {
  name        = var.instance_name
  description = var.instance_name
}

resource "tencentcloud_security_group_rule" "foo" {
  security_group_id = tencentcloud_security_group.foo.id
  type              = "ingress"
  cidr_ip           = "0.0.0.0/0"
  ip_protocol       = "tcp"
  port_range        = "80,8080"
  policy            = "accept"
}

resource "tencentcloud_security_group" "bar" {
  name        = var.instance_name
  description = var.instance_name
}

resource "tencentcloud_security_group_rule" "bar" {
  security_group_id = tencentcloud_security_group.bar.id
  type              = "ingress"
  cidr_ip           = "0.0.0.0/0"
  ip_protocol       = "tcp"
  port_range        = "3000"
  policy            = "accept"
}

resource "tencentcloud_instance" "foo" {
  instance_name              = var.instance_name
  availability_zone          = var.availability_cvm_zone
  image_id                   = data.tencentcloud_images.default.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type           = "CLOUD_PREMIUM"
  security_groups            = %s
}
`,
		ids,
	)
}

func testAccTencentCloudInstanceWithTags(tags string) string {
	return fmt.Sprintf(
		defaultInstanceVariable+`
resource "tencentcloud_instance" "foo" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type  = "CLOUD_PREMIUM"
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }
  tags = %s
}
`,
		tags,
	)
}

const testAccTencentCloudInstanceWithPlacementGroup = defaultInstanceVariable + `
resource "tencentcloud_placement_group" "foo" {
  name = var.instance_name
  type = "HOST"
}

resource "tencentcloud_instance" "foo" {
  instance_name      = var.instance_name
  availability_zone  = var.availability_cvm_zone
  image_id           = data.tencentcloud_images.default.images.0.image_id
  instance_type      = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type   = "CLOUD_PREMIUM"
  placement_group_id = tencentcloud_placement_group.foo.id
}
`

const testAccTencentCloudInstanceWithSpotpaid = defaultInstanceVariable + `
resource "tencentcloud_instance" "foo" {
  instance_name        = var.instance_name
  availability_zone    = var.availability_cvm_zone
  image_id             = data.tencentcloud_images.default.images.0.image_id
  instance_type        = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type     = "CLOUD_PREMIUM"
  instance_charge_type = "SPOTPAID"
  spot_instance_type   = "ONE-TIME"
  spot_max_price       = "0.5"
}
`

func testAccTencentCloudInstanceOrderlySecurityGroups(sgs string) string {

	return fmt.Sprintf(defaultInstanceVariable+`
resource "tencentcloud_security_group" "orderly_security_group1" {
	name        = "test-cvm-orderly-sg1"
	description = "test-cvm-orderly-sg1"
}

resource "tencentcloud_security_group" "orderly_security_group2" {
	name        = "test-cvm-orderly-sg2"
	description = "test-cvm-orderly-sg2"
}

resource "tencentcloud_security_group" "orderly_security_group3" {
	name        = "test-cvm-orderly-sg3"
	description = "test-cvm-orderly-sg3"
}

resource "tencentcloud_instance" "cvm_with_orderly_sg" {
	instance_name              = "test-orderly-sg-cvm"
	availability_zone          = var.availability_cvm_zone
	image_id                   = data.tencentcloud_images.default.images.0.image_id
	instance_type              = data.tencentcloud_instance_types.default.instance_types.0.instance_type
	system_disk_type           = "CLOUD_PREMIUM"
	orderly_security_groups    = %s
}
`, sgs)
}
