package mongodb_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	mongodb_sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/mongodb"
)

type mockMetaForAuditService struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForAuditService) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForAuditService{}

func newMockMetaForAuditService() *mockMetaForAuditService {
	return &mockMetaForAuditService{client: &connectivity.TencentCloudClient{}}
}

func ptrStringAudit(s string) *string { return &s }
func ptrBoolAudit(b bool) *bool       { return &b }
func ptrInt64Audit(v int64) *int64    { return &v }

func TestMongodbAuditService_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mongodbClient := &mongodb_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForAuditService().client, "UseMongodbClient", mongodbClient)

	patches.ApplyMethodFunc(mongodbClient, "OpenAuditService", func(request *mongodb_sdk.OpenAuditServiceRequest) (*mongodb_sdk.OpenAuditServiceResponse, error) {
		resp := mongodb_sdk.NewOpenAuditServiceResponse()
		resp.Response = &mongodb_sdk.OpenAuditServiceResponseParams{
			RequestId: ptrStringAudit("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(mongodbClient, "DescribeAuditConfig", func(request *mongodb_sdk.DescribeAuditConfigRequest) (*mongodb_sdk.DescribeAuditConfigResponse, error) {
		resp := mongodb_sdk.NewDescribeAuditConfigResponse()
		resp.Response = &mongodb_sdk.DescribeAuditConfigResponseParams{
			InstanceId:   ptrStringAudit("cmgo-test1234"),
			InstanceName: ptrStringAudit("test-instance"),
			AuditAll:     ptrBoolAudit(true),
			LogExpireDay: ptrInt64Audit(30),
			CreateTime:   ptrStringAudit("2024-01-01 00:00:00"),
			LogType:      ptrStringAudit("storage"),
			IsClosing:    ptrStringAudit("false"),
			IsOpening:    ptrStringAudit("false"),
			RequestId:    ptrStringAudit("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForAuditService()
	res := mongodb.ResourceTencentCloudMongodbAuditService()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":    "cmgo-test1234",
		"log_expire_day": 30,
		"audit_all":      true,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "cmgo-test1234", d.Id())
}

func TestMongodbAuditService_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mongodbClient := &mongodb_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForAuditService().client, "UseMongodbClient", mongodbClient)

	patches.ApplyMethodFunc(mongodbClient, "DescribeAuditConfig", func(request *mongodb_sdk.DescribeAuditConfigRequest) (*mongodb_sdk.DescribeAuditConfigResponse, error) {
		resp := mongodb_sdk.NewDescribeAuditConfigResponse()
		resp.Response = &mongodb_sdk.DescribeAuditConfigResponseParams{
			InstanceId:   ptrStringAudit("cmgo-test1234"),
			InstanceName: ptrStringAudit("test-instance"),
			AuditAll:     ptrBoolAudit(true),
			LogExpireDay: ptrInt64Audit(30),
			CreateTime:   ptrStringAudit("2024-01-01 00:00:00"),
			LogType:      ptrStringAudit("storage"),
			IsClosing:    ptrStringAudit("false"),
			IsOpening:    ptrStringAudit("false"),
			RequestId:    ptrStringAudit("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForAuditService()
	res := mongodb.ResourceTencentCloudMongodbAuditService()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":    "cmgo-test1234",
		"log_expire_day": 30,
		"audit_all":      true,
	})
	d.SetId("cmgo-test1234")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "cmgo-test1234", d.Get("instance_id"))
	assert.Equal(t, "test-instance", d.Get("instance_name"))
	assert.Equal(t, true, d.Get("audit_all"))
	assert.Equal(t, 30, d.Get("log_expire_day"))
	assert.Equal(t, "2024-01-01 00:00:00", d.Get("create_time"))
	assert.Equal(t, "storage", d.Get("log_type"))
	assert.Equal(t, "false", d.Get("is_closing"))
	assert.Equal(t, "false", d.Get("is_opening"))
}

func TestMongodbAuditService_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mongodbClient := &mongodb_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForAuditService().client, "UseMongodbClient", mongodbClient)

	patches.ApplyMethodFunc(mongodbClient, "ModifyAuditService", func(request *mongodb_sdk.ModifyAuditServiceRequest) (*mongodb_sdk.ModifyAuditServiceResponse, error) {
		resp := mongodb_sdk.NewModifyAuditServiceResponse()
		resp.Response = &mongodb_sdk.ModifyAuditServiceResponseParams{
			RequestId: ptrStringAudit("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(mongodbClient, "DescribeAuditConfig", func(request *mongodb_sdk.DescribeAuditConfigRequest) (*mongodb_sdk.DescribeAuditConfigResponse, error) {
		resp := mongodb_sdk.NewDescribeAuditConfigResponse()
		resp.Response = &mongodb_sdk.DescribeAuditConfigResponseParams{
			InstanceId:   ptrStringAudit("cmgo-test1234"),
			InstanceName: ptrStringAudit("test-instance"),
			AuditAll:     ptrBoolAudit(false),
			LogExpireDay: ptrInt64Audit(90),
			CreateTime:   ptrStringAudit("2024-01-01 00:00:00"),
			LogType:      ptrStringAudit("storage"),
			IsClosing:    ptrStringAudit("false"),
			IsOpening:    ptrStringAudit("false"),
			RequestId:    ptrStringAudit("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForAuditService()
	res := mongodb.ResourceTencentCloudMongodbAuditService()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":    "cmgo-test1234",
		"log_expire_day": 90,
		"audit_all":      false,
	})
	d.SetId("cmgo-test1234")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

func TestMongodbAuditService_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mongodbClient := &mongodb_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForAuditService().client, "UseMongodbClient", mongodbClient)

	patches.ApplyMethodFunc(mongodbClient, "CloseAuditService", func(request *mongodb_sdk.CloseAuditServiceRequest) (*mongodb_sdk.CloseAuditServiceResponse, error) {
		resp := mongodb_sdk.NewCloseAuditServiceResponse()
		resp.Response = &mongodb_sdk.CloseAuditServiceResponseParams{
			RequestId: ptrStringAudit("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(mongodbClient, "DescribeAuditConfig", func(request *mongodb_sdk.DescribeAuditConfigRequest) (*mongodb_sdk.DescribeAuditConfigResponse, error) {
		resp := mongodb_sdk.NewDescribeAuditConfigResponse()
		resp.Response = &mongodb_sdk.DescribeAuditConfigResponseParams{
			InstanceId:   ptrStringAudit("cmgo-test1234"),
			InstanceName: ptrStringAudit("test-instance"),
			AuditAll:     ptrBoolAudit(true),
			LogExpireDay: ptrInt64Audit(30),
			CreateTime:   ptrStringAudit("2024-01-01 00:00:00"),
			LogType:      ptrStringAudit("storage"),
			IsClosing:    ptrStringAudit("false"),
			IsOpening:    ptrStringAudit("false"),
			RequestId:    ptrStringAudit("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForAuditService()
	res := mongodb.ResourceTencentCloudMongodbAuditService()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":    "cmgo-test1234",
		"log_expire_day": 30,
		"audit_all":      true,
	})
	d.SetId("cmgo-test1234")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}
