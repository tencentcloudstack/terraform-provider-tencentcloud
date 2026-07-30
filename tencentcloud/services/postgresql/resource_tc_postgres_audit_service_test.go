package postgresql_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcpostgresql "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/postgresql"
)

type mockMetaAuditService struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaAuditService) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaAuditService{}

func newMockMetaAuditService() *mockMetaAuditService {
	return &mockMetaAuditService{client: &connectivity.TencentCloudClient{}}
}

func ptrStringAudit(s string) *string {
	return &s
}

func ptrUint64Audit(i uint64) *uint64 {
	return &i
}

func ptrFloat64Audit(f float64) *float64 {
	return &f
}

// go test ./tencentcloud/services/postgresql/ -run "TestPostgresAuditService" -v -count=1 -gcflags="all=-l"

// TestPostgresAuditService_Create tests the Create function
func TestPostgresAuditService_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	pgClient := &postgresql.Client{}
	patches.ApplyMethodReturn(newMockMetaAuditService().client, "UsePostgresqlClient", pgClient)

	patches.ApplyMethodFunc(pgClient, "OpenAuditService", func(request *postgresql.OpenAuditServiceRequest) (*postgresql.OpenAuditServiceResponse, error) {
		assert.NotNil(t, request.InstanceId)
		assert.Equal(t, "postgres-test123", *request.InstanceId)
		assert.NotNil(t, request.LogExpireDay)
		assert.Equal(t, uint64(30), *request.LogExpireDay)
		assert.NotNil(t, request.HotLogExpireDay)
		assert.Equal(t, uint64(7), *request.HotLogExpireDay)
		assert.NotNil(t, request.AuditType)
		assert.Equal(t, "simple", *request.AuditType)
		assert.NotNil(t, request.Product)
		assert.Equal(t, "postgres", *request.Product)

		resp := postgresql.NewOpenAuditServiceResponse()
		resp.Response = &postgresql.OpenAuditServiceResponseParams{
			RequestId: ptrStringAudit("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(pgClient, "DescribeAuditInstanceList", func(request *postgresql.DescribeAuditInstanceListRequest) (*postgresql.DescribeAuditInstanceListResponse, error) {
		resp := postgresql.NewDescribeAuditInstanceListResponse()
		resp.Response = &postgresql.DescribeAuditInstanceListResponseParams{
			TotalCount: ptrUint64Audit(1),
			Items: []*postgresql.AuditInstanceInfo{
				{
					InstanceId:       ptrStringAudit("postgres-test123"),
					AuditStatus:      ptrStringAudit("ON"),
					LogExpireDay:     ptrUint64Audit(30),
					HotLogExpireDay:  ptrUint64Audit(7),
					ColdLogExpireDay: ptrUint64Audit(23),
					HotLogSize:       ptrFloat64Audit(100.5),
					ColdLogSize:      ptrFloat64Audit(200.3),
					CreateTime:       ptrStringAudit("2024-01-01 00:00:00"),
				},
			},
			RequestId: ptrStringAudit("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaAuditService()
	res := svcpostgresql.ResourceTencentCloudPostgresAuditService()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":        "postgres-test123",
		"log_expire_day":     30,
		"hot_log_expire_day": 7,
		"audit_type":         "simple",
		"product":            "postgres",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "postgres-test123", d.Id())
}

// TestPostgresAuditService_Read tests the Read function
func TestPostgresAuditService_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	pgClient := &postgresql.Client{}
	patches.ApplyMethodReturn(newMockMetaAuditService().client, "UsePostgresqlClient", pgClient)

	patches.ApplyMethodFunc(pgClient, "DescribeAuditInstanceList", func(request *postgresql.DescribeAuditInstanceListRequest) (*postgresql.DescribeAuditInstanceListResponse, error) {
		assert.NotNil(t, request.Product)
		assert.Equal(t, "postgres", *request.Product)
		assert.NotNil(t, request.AuditSwitch)
		assert.Equal(t, uint64(1), *request.AuditSwitch)

		resp := postgresql.NewDescribeAuditInstanceListResponse()
		resp.Response = &postgresql.DescribeAuditInstanceListResponseParams{
			TotalCount: ptrUint64Audit(1),
			Items: []*postgresql.AuditInstanceInfo{
				{
					InstanceId:       ptrStringAudit("postgres-test123"),
					AuditStatus:      ptrStringAudit("ON"),
					LogExpireDay:     ptrUint64Audit(30),
					HotLogExpireDay:  ptrUint64Audit(7),
					ColdLogExpireDay: ptrUint64Audit(23),
					HotLogSize:       ptrFloat64Audit(100.5),
					ColdLogSize:      ptrFloat64Audit(200.3),
					CreateTime:       ptrStringAudit("2024-01-01 00:00:00"),
				},
			},
			RequestId: ptrStringAudit("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaAuditService()
	res := svcpostgresql.ResourceTencentCloudPostgresAuditService()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":        "postgres-test123",
		"log_expire_day":     30,
		"hot_log_expire_day": 7,
		"audit_type":         "simple",
		"product":            "postgres",
	})
	d.SetId("postgres-test123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "postgres-test123", d.Id())
}

// TestPostgresAuditService_Update tests the Update function
func TestPostgresAuditService_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	pgClient := &postgresql.Client{}
	patches.ApplyMethodReturn(newMockMetaAuditService().client, "UsePostgresqlClient", pgClient)

	patches.ApplyMethodFunc(pgClient, "ModifyAuditService", func(request *postgresql.ModifyAuditServiceRequest) (*postgresql.ModifyAuditServiceResponse, error) {
		assert.NotNil(t, request.InstanceId)
		assert.Equal(t, "postgres-test123", *request.InstanceId)
		assert.NotNil(t, request.LogExpireDay)
		assert.Equal(t, uint64(90), *request.LogExpireDay)
		assert.NotNil(t, request.HotLogExpireDay)
		assert.Equal(t, uint64(30), *request.HotLogExpireDay)
		assert.NotNil(t, request.AuditType)
		assert.Equal(t, "complex", *request.AuditType)

		resp := postgresql.NewModifyAuditServiceResponse()
		resp.Response = &postgresql.ModifyAuditServiceResponseParams{
			RequestId: ptrStringAudit("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(pgClient, "DescribeAuditInstanceList", func(request *postgresql.DescribeAuditInstanceListRequest) (*postgresql.DescribeAuditInstanceListResponse, error) {
		resp := postgresql.NewDescribeAuditInstanceListResponse()
		resp.Response = &postgresql.DescribeAuditInstanceListResponseParams{
			TotalCount: ptrUint64Audit(1),
			Items: []*postgresql.AuditInstanceInfo{
				{
					InstanceId:       ptrStringAudit("postgres-test123"),
					AuditStatus:      ptrStringAudit("ON"),
					LogExpireDay:     ptrUint64Audit(90),
					HotLogExpireDay:  ptrUint64Audit(30),
					ColdLogExpireDay: ptrUint64Audit(60),
					HotLogSize:       ptrFloat64Audit(150.0),
					ColdLogSize:      ptrFloat64Audit(300.0),
					CreateTime:       ptrStringAudit("2024-01-01 00:00:00"),
				},
			},
			RequestId: ptrStringAudit("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaAuditService()
	res := svcpostgresql.ResourceTencentCloudPostgresAuditService()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":        "postgres-test123",
		"log_expire_day":     90,
		"hot_log_expire_day": 30,
		"audit_type":         "complex",
		"product":            "postgres",
	})
	d.SetId("postgres-test123")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "postgres-test123", d.Id())
}

// TestPostgresAuditService_Delete tests the Delete function
func TestPostgresAuditService_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	pgClient := &postgresql.Client{}
	patches.ApplyMethodReturn(newMockMetaAuditService().client, "UsePostgresqlClient", pgClient)

	patches.ApplyMethodFunc(pgClient, "CloseAuditService", func(request *postgresql.CloseAuditServiceRequest) (*postgresql.CloseAuditServiceResponse, error) {
		assert.NotNil(t, request.InstanceId)
		assert.Equal(t, "postgres-test123", *request.InstanceId)
		assert.NotNil(t, request.Product)
		assert.Equal(t, "postgres", *request.Product)

		resp := postgresql.NewCloseAuditServiceResponse()
		resp.Response = &postgresql.CloseAuditServiceResponseParams{
			RequestId: ptrStringAudit("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaAuditService()
	res := svcpostgresql.ResourceTencentCloudPostgresAuditService()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":        "postgres-test123",
		"log_expire_day":     30,
		"hot_log_expire_day": 7,
		"audit_type":         "simple",
		"product":            "postgres",
	})
	d.SetId("postgres-test123")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}
