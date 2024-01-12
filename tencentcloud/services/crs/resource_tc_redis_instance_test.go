package crs_test

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/stretchr/testify/assert"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	svccrs "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/crs"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=redis_instance
	resource.AddTestSweepers("redis_instance", &resource.Sweeper{
		Name: "redis_instance",
		F: func(region string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			cli, _ := tcacctest.SharedClientForRegion(region)
			client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()

			service := svccrs.NewRedisService(client)

			instances, err := service.DescribeInstances(ctx, "ap-guangzhou-3", "", 0, 10)

			if err != nil {
				return err
			}

			for _, v := range instances {
				name := v.Name
				id := v.RedisId
				if tcacctest.IsResourcePersist(name, nil) {
					continue
				}
				// Collect infos before deleting action
				var chargeType string
				has, online, info, err := service.CheckRedisOnlineOk(ctx, id, tccommon.ReadRetryTimeout*20)
				if !has {
					continue
				}
				if online {
					chargeType = svccrs.REDIS_CHARGE_TYPE_NAME[*info.BillingMode]
				} else {
					log.Printf("Deleting ERROR: Creating redis task is processing.")
					continue
				}
				if err != nil {
					log.Printf("[CRITAL]%s redis querying before deleting task fail, reason:%s\n", logId, err.Error())
					continue
				}

				var wait = func(action string, taskInfo interface{}) (errRet error) {

					errRet = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
						var ok bool
						var err error
						switch v := taskInfo.(type) {
						case int64:
							ok, err = service.DescribeTaskInfo(ctx, id, v)
						case string:
							ok, _, err = service.DescribeInstanceDealDetail(ctx, v)
						}
						if err != nil {
							if _, ok := err.(*sdkErrors.TencentCloudSDKError); !ok {
								return resource.RetryableError(err)
							} else {
								return resource.NonRetryableError(err)
							}
						}
						if ok {
							return nil
						} else {
							return resource.RetryableError(fmt.Errorf("%s timeout.", action))
						}
					})

					if errRet != nil {
						log.Printf("[CRITAL]%s redis %s fail, reason:%s\n", logId, action, errRet.Error())
					}
					return errRet
				}

				if chargeType == svccrs.REDIS_CHARGE_TYPE_POSTPAID {
					taskId, err := service.DestroyPostpaidInstance(ctx, id)
					if err != nil {
						log.Printf("[CRITAL]%s redis %s fail, reason:%s\n", logId, "DestroyPostpaidInstance", err.Error())
						return err
					}
					if err = wait("DestroyPostpaidInstance", taskId); err != nil {
						return err
					}

				} else {
					if _, err := service.DestroyPrepaidInstance(ctx, id); err != nil {
						log.Printf("[CRITAL]%s redis %s fail, reason:%s\n", logId, "DestroyPrepaidInstance", err.Error())
						return err
					}

					// Deal info only support create and renew and resize, need to check destroy status by describing api.
					if errDestroyChecking := resource.Retry(20*tccommon.ReadRetryTimeout, func() *resource.RetryError {
						has, isolated, err := service.CheckRedisDestroyOk(ctx, id)
						if err != nil {
							log.Printf("[CRITAL]%s CheckRedisDestroyOk fail, reason:%s\n", logId, err.Error())
							return resource.NonRetryableError(err)
						}
						if !has || isolated {
							return nil
						}
						return resource.RetryableError(fmt.Errorf("instance is not ready to be destroyed"))
					}); errDestroyChecking != nil {
						log.Printf("[CRITAL]%s redis querying before deleting task fail, reason:%s\n", logId, errDestroyChecking.Error())
						return errDestroyChecking
					}
				}

				taskId, err := service.CleanUpInstance(ctx, id)
				if err != nil {
					log.Printf("[CRITAL]%s redis %s fail, reason:%s\n", logId, "CleanUpInstance", err.Error())
					return err
				}

				_ = wait("CleanUpInstance", taskId)
			}

			return nil
		},
	})
}

func TestAccTencentCloudRedisInstanceResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccTencentCloudRedisInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisInstanceBasic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTencentCloudRedisInstanceExists("tencentcloud_redis_instance.redis_instance_test"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_instance.redis_instance_test", "ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_instance.redis_instance_test", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "port", "6379"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "type_id", "2"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "redis_shard_num", "1"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "redis_replicas_num", "1"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "mem_size", "8192"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "name", "terraform_test"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "status", "online"),
				),
			},
			{
				Config: testAccRedisInstanceTags(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTencentCloudRedisInstanceExists("tencentcloud_redis_instance.redis_instance_test"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "tags.test", "test"),
				),
			},
			{
				Config: testAccRedisInstanceTagsUpdate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTencentCloudRedisInstanceExists("tencentcloud_redis_instance.redis_instance_test"),
					resource.TestCheckNoResourceAttr("tencentcloud_redis_instance.redis_instance_test", "tags.test"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "tags.abc", "abc"),
				),
			},
			{
				Config: testAccRedisInstanceUpdateName(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTencentCloudRedisInstanceExists("tencentcloud_redis_instance.redis_instance_test"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_instance.redis_instance_test", "ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_instance.redis_instance_test", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "port", "6379"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "type_id", "2"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "redis_shard_num", "1"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "redis_replicas_num", "1"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "mem_size", "8192"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "name", "terraform_test_update"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "status", "online"),
				),
			},
			{
				Config: testAccRedisInstanceUpdateMemSizeAndPasswordAndSg(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTencentCloudRedisInstanceExists("tencentcloud_redis_instance.redis_instance_test"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_instance.redis_instance_test", "ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_instance.redis_instance_test", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "port", "6379"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "type_id", "2"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "redis_shard_num", "1"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "redis_replicas_num", "1"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "type_id", "2"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "mem_size", "12288"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "name", "terraform_test_update"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "status", "online"),
				),
			},
			{
				Config: testAccRedisInstanceNetworkUpdate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTencentCloudRedisInstanceExists("tencentcloud_redis_instance.redis_instance_test"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "port", "6380"),
				),
			},
			{
				ResourceName:            "tencentcloud_redis_instance.redis_instance_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "type", "redis_shard_num", "redis_replicas_num", "force_delete", "operation_network"},
			},
		},
	})
}

func TestAccTencentCloudRedisInstanceParam(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccTencentCloudRedisInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisInstanceParamTemplate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTencentCloudRedisInstanceExists("tencentcloud_redis_instance.redis_instance_test"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_instance.redis_instance_test", "ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_instance.redis_instance_test", "create_time"),
				),
			},
			{
				Config: testAccRedisInstanceParamTemplateUpdate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTencentCloudRedisInstanceExists("tencentcloud_redis_instance.redis_instance_test"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_instance.redis_instance_test", "ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_instance.redis_instance_test", "create_time"),
				),
			},
		},
	})
}

func TestAccTencentCloudRedisInstance_Maz(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccTencentCloudRedisInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisInstanceMaz(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTencentCloudRedisInstanceExists("tencentcloud_redis_instance.redis_maz"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_maz", "mem_size", "2048"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_maz", "redis_replicas_num", "2"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_maz", "replica_zone_ids.#", "2"),
				),
			},
			{
				Config: testAccRedisInstanceMazUpdate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTencentCloudRedisInstanceExists("tencentcloud_redis_instance.redis_maz"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_maz", "mem_size", "4096"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_maz", "redis_replicas_num", "3"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_maz", "replica_zone_ids.#", "3"),
				),
			},
			{
				Config: testAccRedisInstanceMazUpdate2(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTencentCloudRedisInstanceExists("tencentcloud_redis_instance.redis_maz"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_maz", "mem_size", "2048"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_maz", "redis_replicas_num", "3"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_maz", "replica_zone_ids.#", "3"),
				),
			},
		},
	})
}

func TestAccTencentCloudRedisInstance_Cluster(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccTencentCloudRedisInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisInstanceCluster,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTencentCloudRedisInstanceExists("tencentcloud_redis_instance.redis_cluster"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_cluster", "redis_shard_num", "1"),
				),
			},
			{
				Config: testAccRedisInstanceClusterUpdateShard(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTencentCloudRedisInstanceExists("tencentcloud_redis_instance.redis_cluster"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_cluster", "redis_shard_num", "3"),
				),
			},
		},
	})
}

// go test -i; go test -test.run TestAccTencentCloudRedisInstance_ReplicasNum -v
func TestAccTencentCloudRedisInstance_ReplicasNum(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccTencentCloudRedisInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisReplicasNum,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTencentCloudRedisInstanceExists("tencentcloud_redis_instance.redis_instance_replicas"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_replicas", "redis_replicas_num", "1"),
				),
			},
			{
				Config: testAccRedisReplicasNumUp,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTencentCloudRedisInstanceExists("tencentcloud_redis_instance.redis_instance_replicas"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_replicas", "redis_replicas_num", "2"),
				),
			},
		},
	})
}

func TestAccTencentCloudRedisInstance_Prepaid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccTencentCloudRedisInstanceDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
				Config:    testAccRedisInstancePrepaidBasic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTencentCloudRedisInstanceExists("tencentcloud_redis_instance.redis_prepaid_instance_test"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_instance.redis_prepaid_instance_test", "ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_instance.redis_prepaid_instance_test", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_prepaid_instance_test", "port", "6379"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_prepaid_instance_test", "type_id", "2"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_prepaid_instance_test", "redis_shard_num", "1"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_prepaid_instance_test", "redis_replicas_num", "1"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_prepaid_instance_test", "mem_size", "8192"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_prepaid_instance_test", "name", "terraform_prepaid_test"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_prepaid_instance_test", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_prepaid_instance_test", "status", "online"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_prepaid_instance_test", "charge_type", "PREPAID"),
				),
			},
			{
				ResourceName:            "tencentcloud_redis_instance.redis_prepaid_instance_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "type", "redis_shard_num", "redis_replicas_num", "force_delete", "prepaid_period"},
			},
		},
	})
}

func TestAccTencentCloudRedisGetRemoveNodesByIds(t *testing.T) {
	mockNodes1 := []*redis.RedisNodeInfo{
		{
			NodeType: helper.IntInt64(0),
			ZoneId:   helper.IntUint64(100001),
			NodeId:   helper.IntInt64(101),
		},
		{
			NodeType: helper.IntInt64(1),
			ZoneId:   helper.IntUint64(100001),
			NodeId:   helper.IntInt64(102),
		},
		{
			NodeType: helper.IntInt64(1),
			ZoneId:   helper.IntUint64(100001),
			NodeId:   helper.IntInt64(103),
		},
		{
			NodeType: helper.IntInt64(1),
			ZoneId:   helper.IntUint64(100002),
			NodeId:   helper.IntInt64(104),
		},
		{
			NodeType: helper.IntInt64(1),
			ZoneId:   helper.IntUint64(100002),
			NodeId:   helper.IntInt64(105),
		},
		{
			NodeType: helper.IntInt64(1),
			ZoneId:   helper.IntUint64(100003),
			NodeId:   helper.IntInt64(106),
		},
	}

	origin := []int{
		100001,
		100001,
		100002,
		100002,
		100003,
	}
	mockAdds1, mockRemoves1 := tccommon.GetListDiffs(
		origin,
		[]int{
			100001,
			// -100001
			100002,
			// -100002
			100003,
			100004, // +
		},
	)
	assert.Contains(t, []int{100001, 100002}, mockRemoves1[0])
	assert.Contains(t, []int{100001, 100002}, mockRemoves1[1])
	assert.Equal(t, []int{100004}, mockAdds1)

	mockAdds2, mockRemoves2 := tccommon.GetListDiffs(origin, []int{100001, 100002})
	assert.Equal(t, len(mockRemoves2), 3)
	assert.Contains(t, mockRemoves2, 100001)
	assert.Contains(t, mockRemoves2, 100002)
	assert.Contains(t, mockRemoves2, 100003)
	assert.Equal(t, 0, len(mockAdds2))

	result1 := svccrs.TencentCloudRedisGetRemoveNodesByIds(mockRemoves1[:], mockNodes1)

	result1Len := len(result1)
	if result1Len != 2 {
		t.Errorf("result1 length expect %d, got %d", 2, result1Len)
		return
	}
	assert.Equal(t, int64(102), *result1[0].NodeId)
	assert.Equal(t, int64(104), *result1[1].NodeId)

	result2 := svccrs.TencentCloudRedisGetRemoveNodesByIds(mockRemoves2[:], mockNodes1)

	assert.Equal(t, int64(102), *result2[0].NodeId)
	assert.Equal(t, int64(104), *result2[1].NodeId)
	assert.Equal(t, int64(106), *result2[2].NodeId)
}

func testAccTencentCloudRedisInstanceExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svccrs.NewRedisService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		has, _, _, err := service.CheckRedisOnlineOk(ctx, rs.Primary.ID, time.Second)
		if has {
			return nil
		}
		if err != nil {
			return err
		}
		return fmt.Errorf("redis not exists.")
	}
}

func testAccTencentCloudRedisInstanceDestroy(s *terraform.State) error {

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svccrs.NewRedisService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_redis_instance" {
			continue
		}
		time.Sleep(5 * time.Second)
		has, isolated, err := service.CheckRedisDestroyOk(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if !has || isolated {
			return nil
		}
		return fmt.Errorf("redis not delete ok")
	}
	return nil
}

// FIXME use datasource instead
const testAccRedisDefaultTemplate = `
variable "redis_param_template" {
  default = "crs-cfg-1q38ngo0"
}

variable "redis_default_param_template" {
  default = "default-param-template-6"
}
`

func testAccRedisInstanceBasic() string {
	return tcacctest.DefaultVpcVariable + `
resource "tencentcloud_redis_instance" "redis_instance_test" {
  availability_zone  = "ap-guangzhou-3"
  type_id            = 2
  password           = "test12345789"
  mem_size           = 8192
  name               = "terraform_test"
  port               = 6379
  redis_shard_num    = 1
  redis_replicas_num = 1
  vpc_id 			 = var.vpc_id
  subnet_id			 = var.subnet_id
}`
}

func testAccRedisInstanceTags() string {
	return tcacctest.DefaultVpcVariable + `
resource "tencentcloud_redis_instance" "redis_instance_test" {
  availability_zone = "ap-guangzhou-3"
  type_id            = 2
  password           = "test12345789"
  mem_size           = 8192
  name               = "terraform_test"
  port               = 6379
  redis_shard_num    = 1
  redis_replicas_num = 1
  vpc_id 			 = var.vpc_id
  subnet_id			 = var.subnet_id

  tags = {
    test = "test"
  }
}`
}

func testAccRedisInstanceTagsUpdate() string {
	return tcacctest.DefaultVpcVariable + `
resource "tencentcloud_redis_instance" "redis_instance_test" {
  availability_zone  = "ap-guangzhou-3"
  type_id            = 2
  password           = "test12345789"
  mem_size           = 8192
  name               = "terraform_test"
  port               = 6379
  redis_shard_num    = 1
  redis_replicas_num = 1
  vpc_id 			 = var.vpc_id
  subnet_id			 = var.subnet_id

  tags = {
    abc = "abc"
  }
}`
}

func testAccRedisInstanceUpdateName() string {
	return tcacctest.DefaultVpcVariable + `
resource "tencentcloud_redis_instance" "redis_instance_test" {
  availability_zone  = "ap-guangzhou-3"
  type_id            = 2
  password           = "test12345789"
  mem_size           = 8192
  name               = "terraform_test_update"
  port               = 6379
  redis_shard_num    = 1
  redis_replicas_num = 1
  vpc_id 			 = var.vpc_id
  subnet_id			 = var.subnet_id

  tags = {
    abc = "abc"
  }
}`
}

func testAccRedisInstanceUpdateMemSizeAndPasswordAndSg() string {
	return tcacctest.DefaultVpcVariable + `
resource "tencentcloud_redis_instance" "redis_instance_test" {
  availability_zone = "ap-guangzhou-3"
  type_id            = 2
  password           = "AAA123456BBB"
  mem_size           = 12288
  name               = "terraform_test_update"
  port               = 6379
  redis_shard_num    = 1
  redis_replicas_num = 1
  vpc_id 			 = var.vpc_id
  subnet_id			 = var.subnet_id
  no_auth            = true
  security_groups    = [var.sg_id]

  tags = {
    "abc" = "abc"
  }
}`
}

func testAccRedisInstanceNetworkUpdate() string {
	return tcacctest.DefaultVpcVariable + `
resource "tencentcloud_redis_instance" "redis_instance_test" {
  availability_zone = "ap-guangzhou-3"
  type_id            = 2
  password           = "AAA123456BBB"
  mem_size           = 12288
  name               = "terraform_test_update"
  port               = 6380
  redis_shard_num    = 1
  redis_replicas_num = 1
  vpc_id 			 = var.vpc_id
  subnet_id			 = var.subnet_id
  no_auth            = true
  security_groups    = [var.sg_id]

  tags = {
    "abc" = "abc"
  }

  operation_network  = "changeVPort"
}`
}

var randMazInstanceName = fmt.Sprintf(`
variable "redis_maz_name" {
  default = "terraform_maz_%d"
}
`, rand.Intn(1000))

func testAccRedisInstanceMaz() string {
	return tcacctest.DefaultVpcVariable + randMazInstanceName + `
resource "tencentcloud_redis_instance" "redis_maz" {
  availability_zone = "ap-guangzhou-3"
  type_id            = 6 #7
  password           = "AAA123456BBB"
  mem_size           = 2048
  name               = var.redis_maz_name
  port               = 6379
  redis_shard_num    = 1
  redis_replicas_num = 2
  replica_zone_ids   = [100003, 100004]
  vpc_id 			 = var.vpc_id
  subnet_id			 = var.subnet_id
}`
}

func testAccRedisInstanceMazUpdate() string {
	return tcacctest.DefaultVpcVariable + randMazInstanceName + `
resource "tencentcloud_redis_instance" "redis_maz" {
  availability_zone = "ap-guangzhou-3"
  type_id            = 6 #7
  password           = "AAA123456BBB"
  mem_size           = 4096
  name               = var.redis_maz_name
  port               = 6379
  redis_shard_num    = 1
  redis_replicas_num = 3
  replica_zone_ids   = [100003, 100003, 100004]
  vpc_id 			 = var.vpc_id
  subnet_id			 = var.subnet_id
}`
}

func testAccRedisInstanceMazUpdate2() string {
	return tcacctest.DefaultVpcVariable + randMazInstanceName + `
resource "tencentcloud_redis_instance" "redis_maz" {
  availability_zone = "ap-guangzhou-3"
  type_id            = 6 #7
  password           = "AAA123456BBBC"
  mem_size           = 2048
  name               = var.redis_maz_name
  port               = 6379
  redis_shard_num    = 1
  redis_replicas_num = 3
  replica_zone_ids   = [100003, 100006, 100007]
  vpc_id 			 = var.vpc_id
  subnet_id 		 = var.subnet_id
}`
}

const testAccRedisInstanceCluster = tcacctest.DefaultVpcVariable + `
resource "tencentcloud_redis_instance" "redis_cluster" {
  availability_zone = "ap-guangzhou-3"
  type_id            = 7
  password           = "AAA123456BBB"
  mem_size           = 4096
  name               = "terraform_cluster"
  port               = 6379
  redis_shard_num    = 1
  redis_replicas_num = 2
  replica_zone_ids   = [100003, 100004]
  vpc_id 			 = var.vpc_id
  subnet_id			 = var.subnet_id
}`

func testAccRedisInstanceClusterUpdateShard() string {
	return tcacctest.DefaultVpcVariable + `
resource "tencentcloud_redis_instance" "redis_cluster" {
  availability_zone = "ap-guangzhou-3"
  type_id            = 7
  password           = "AAA123456BBB"
  mem_size           = 4096
  name               = "terraform_cluster"
  port               = 6379
  redis_shard_num    = 3
  redis_replicas_num = 2
  replica_zone_ids   = [100003, 100004]
  vpc_id 			 = var.vpc_id
  subnet_id			 = var.subnet_id
}`
}

func testAccRedisInstancePrepaidBasic() string {
	return tcacctest.DefaultVpcSubnets + `
resource "tencentcloud_redis_instance" "redis_prepaid_instance_test" {
  availability_zone  = "ap-guangzhou-3"
  type_id            = 2
  password           = "test12345789"
  mem_size           = 8192
  name               = "terraform_prepaid_test"
  port               = 6379
  redis_shard_num    = 1
  redis_replicas_num = 1
  charge_type        = "PREPAID"
  prepaid_period     = 2
  force_delete       = true
  vpc_id 			 = local.vpc_id
  subnet_id			 = local.subnet_id
}`
}

func testAccRedisInstanceParamTemplate() string {
	return tcacctest.DefaultVpcVariable + testAccRedisDefaultTemplate + `
resource "tencentcloud_redis_instance" "redis_instance_test" {
  availability_zone  = "ap-guangzhou-3"
  type_id            = 6
  password           = "test12345789"
  mem_size           = 8192
  name               = "terraform_test"
  port               = 6379
  redis_shard_num    = 1
  redis_replicas_num = 1
  params_template_id = var.redis_param_template
  vpc_id 			 = var.vpc_id
  subnet_id			 = var.subnet_id
}`
}

func testAccRedisInstanceParamTemplateUpdate() string {
	return tcacctest.DefaultVpcVariable + testAccRedisDefaultTemplate + `
resource "tencentcloud_redis_instance" "redis_instance_test" {
  availability_zone  = "ap-guangzhou-3"
  type_id            = 6
  password           = "test12345789"
  mem_size           = 8192
  name               = "terraform_test"
  port               = 6379
  redis_shard_num    = 1
  redis_replicas_num = 1
  params_template_id = var.redis_default_param_template
  vpc_id 			 = var.vpc_id
  subnet_id			 = var.subnet_id
}
`
}

const testAccRedisReplicasNum = tcacctest.DefaultCrsVar + `
resource "tencentcloud_redis_instance" "redis_instance_replicas" {
    auto_renew_flag    = 0
    availability_zone  = "ap-guangzhou-6"
    charge_type        = "POSTPAID"
	password           = "test12345789"
    mem_size           = 4096
    name               = "terraform_test_replicas"
    port               = 6379
    project_id         = 0
    redis_replicas_num = 1
    redis_shard_num    = 3
    replicas_read_only = false
    security_groups    = [
        "sg-edmur627",
    ]
    subnet_id          = var.subnet_id
    type_id            = 7
    vpc_id             = var.vpc_id
}`

const testAccRedisReplicasNumUp = tcacctest.DefaultCrsVar + `
resource "tencentcloud_redis_instance" "redis_instance_replicas" {
    auto_renew_flag    = 0
    availability_zone  = "ap-guangzhou-6"
    charge_type        = "POSTPAID"
	password           = "test12345789"
    mem_size           = 4096
    name               = "terraform_test_replicas"
    port               = 6379
    project_id         = 0
    redis_replicas_num = 2
    redis_shard_num    = 3
    replicas_read_only = false
    security_groups    = [
        "sg-edmur627",
    ]
    subnet_id          = var.subnet_id
    tags               = {}
    type_id            = 7
    vpc_id             = var.vpc_id
}`
