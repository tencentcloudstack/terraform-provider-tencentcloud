package ckafka_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	localckafka "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ckafka"
)

func TestAccTencentCloudCkafkaAclResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCkafkaAclDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCkafkaAcl,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCkafkaAclExists("tencentcloud_ckafka_acl.foo"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_acl.foo", "resource_type", "TOPIC"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_acl.foo", "operation_type", "WRITE"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_acl.foo", "permission_type", "ALLOW"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_acl.foo", "host", "10.10.10.0"),
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_acl.foo", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_acl.foo", "resource_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_acl.foo", "principal"),
				),
			},
			{
				ResourceName:      "tencentcloud_ckafka_acl.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCkafkaAclExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		ckafkaService := localckafka.NewCkafkaService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("ckafka acl %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("ckafka acl id is not set")
		}

		_, has, err := ckafkaService.DescribeAclByAclId(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if !has {
			return fmt.Errorf("ckafka acl doesn't exist: %s", rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckCkafkaAclDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	ckafkaService := localckafka.NewCkafkaService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_ckafka_acl" {
			continue
		}

		_, has, err := ckafkaService.DescribeAclByAclId(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if !has {
			return nil
		}
		return fmt.Errorf("ckafka acl still exists: %s", rs.Primary.ID)
	}
	return nil
}

const testAccCkafkaAcl = tcacctest.DefaultKafkaVariable + `
resource "tencentcloud_ckafka_user" "foo" {
	instance_id  = var.instance_id
	account_name = "tf-test-acl-resource"
	password     = "test1234"
  }

resource "tencentcloud_ckafka_topic" "kafka_topic_acl" {
	instance_id                     = var.instance_id
	topic_name                      = "ckafka-topic-acl-test"
	replica_num                     = 2
	partition_num                   = 1
	note                            = "test topic"
	enable_white_list               = true
	ip_white_list                   = ["192.168.1.1"]
	clean_up_policy                 = "delete"
	sync_replica_min_num            = 1
	unclean_leader_election_enable  = false
	segment                         = 86400000
	retention                       = 60000
	max_message_bytes               = 8388608
}

resource "tencentcloud_ckafka_acl" foo {
  instance_id     = var.instance_id
  resource_type   = "TOPIC"
  resource_name   = tencentcloud_ckafka_topic.kafka_topic_acl.topic_name
  operation_type  = "WRITE"
  permission_type = "ALLOW"
  host            = "10.10.10.0"
  principal       = tencentcloud_ckafka_user.foo.account_name
}
`

// -----------------------------------------------------------------------------
// Mock-based tests for CreateAcl FailedOperation retry via resource.Retry
// with a max of 5 FailedOperation attempts before giving up.
//
// go test ./tencentcloud/services/ckafka/ -run "TestCkafkaAclCreateAclFailedOperation" -v -count=1 -gcflags="all=-l"
// -----------------------------------------------------------------------------

type mockMetaCkafkaAclFailedOperationRetry struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaCkafkaAclFailedOperationRetry) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaCkafkaAclFailedOperationRetry{}

func newMockMetaCkafkaAclFailedOperationRetry() *mockMetaCkafkaAclFailedOperationRetry {
	return &mockMetaCkafkaAclFailedOperationRetry{client: &connectivity.TencentCloudClient{}}
}

// TestCkafkaAclCreateAclFailedOperationRetryThenSuccess verifies that when the
// CreateAcl API returns FailedOperation for the first 3 calls and succeeds on the
// 4th call, CkafkaService.CreateAcl returns nil and the API is invoked 4 times
// (within the 5 FailedOperation retry limit).
func TestCkafkaAclCreateAclFailedOperationRetryThenSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ckafkaClient := &ckafka.Client{}
	patches.ApplyMethodReturn(newMockMetaCkafkaAclFailedOperationRetry().client, "UseCkafkaClient", ckafkaClient)

	var callCount int
	patches.ApplyMethodFunc(ckafkaClient, "CreateAcl", func(request *ckafka.CreateAclRequest) (*ckafka.CreateAclResponse, error) {
		callCount++
		if callCount <= 3 {
			return nil, &sdkErrors.TencentCloudSDKError{Code: "FailedOperation", Message: "operation failed"}
		}
		resp := ckafka.NewCreateAclResponse()
		resp.Response = &ckafka.CreateAclResponseParams{
			Result: &ckafka.JgwOperateResponse{
				ReturnCode: dpPtrString("0"),
			},
			RequestId: dpPtrString("fake-request-id"),
		}
		return resp, nil
	})

	service := localckafka.NewCkafkaService(newMockMetaCkafkaAclFailedOperationRetry().client)
	err := service.CreateAcl(context.TODO(), "ckafka-instance-retry", "TOPIC", "topic-acl-retry", "WRITE", "ALLOW", "*", "root")
	assert.NoError(t, err)
	assert.Equal(t, 4, callCount)
}

// TestCkafkaAclCreateAclFailedOperationAlwaysFails verifies that when the
// CreateAcl API always returns FailedOperation, CkafkaService.CreateAcl gives up
// after 5 FailedOperation errors (the max retry limit) and returns the
// FailedOperation error without further retries.
func TestCkafkaAclCreateAclFailedOperationAlwaysFails(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ckafkaClient := &ckafka.Client{}
	patches.ApplyMethodReturn(newMockMetaCkafkaAclFailedOperationRetry().client, "UseCkafkaClient", ckafkaClient)

	var callCount int
	patches.ApplyMethodFunc(ckafkaClient, "CreateAcl", func(request *ckafka.CreateAclRequest) (*ckafka.CreateAclResponse, error) {
		callCount++
		return nil, &sdkErrors.TencentCloudSDKError{Code: "FailedOperation", Message: "operation failed"}
	})

	service := localckafka.NewCkafkaService(newMockMetaCkafkaAclFailedOperationRetry().client)
	err := service.CreateAcl(context.TODO(), "ckafka-instance-retry", "TOPIC", "topic-acl-retry", "WRITE", "ALLOW", "*", "root")
	assert.Error(t, err)
	sdkErr, ok := err.(*sdkErrors.TencentCloudSDKError)
	assert.True(t, ok)
	assert.Equal(t, "FailedOperation", sdkErr.Code)
	assert.Equal(t, 5, callCount)
}

// TestCkafkaAclCreateAclNonFailedOperationError verifies that when the CreateAcl
// API returns a non-FailedOperation error, CkafkaService.CreateAcl returns the
// error immediately via NonRetryableError and the API is invoked only once.
func TestCkafkaAclCreateAclNonFailedOperationError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ckafkaClient := &ckafka.Client{}
	patches.ApplyMethodReturn(newMockMetaCkafkaAclFailedOperationRetry().client, "UseCkafkaClient", ckafkaClient)

	var callCount int
	patches.ApplyMethodFunc(ckafkaClient, "CreateAcl", func(request *ckafka.CreateAclRequest) (*ckafka.CreateAclResponse, error) {
		callCount++
		return nil, &sdkErrors.TencentCloudSDKError{Code: "InvalidParameter", Message: "invalid parameter"}
	})

	service := localckafka.NewCkafkaService(newMockMetaCkafkaAclFailedOperationRetry().client)
	err := service.CreateAcl(context.TODO(), "ckafka-instance-retry", "TOPIC", "topic-acl-retry", "WRITE", "ALLOW", "*", "root")
	assert.Error(t, err)
	sdkErr, ok := err.(*sdkErrors.TencentCloudSDKError)
	assert.True(t, ok)
	assert.Equal(t, "InvalidParameter", sdkErr.Code)
	assert.Equal(t, 1, callCount)
}
