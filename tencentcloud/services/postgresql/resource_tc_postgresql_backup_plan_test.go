package postgresql_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcpostgresql "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/postgresql"
)

type mockMetaBackupPlan struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaBackupPlan) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaBackupPlan{}

func newMockMetaBackupPlan() *mockMetaBackupPlan {
	return &mockMetaBackupPlan{client: &connectivity.TencentCloudClient{}}
}

func ptrStringBackupPlan(s string) *string {
	return &s
}

func ptrUint64BackupPlan(i uint64) *uint64 {
	return &i
}

// go test ./tencentcloud/services/postgresql/ -run "TestPostgresqlBackupPlan" -v -count=1 -gcflags="all=-l"

func TestPostgresqlBackupPlan_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	pgClient := &postgresql.Client{}
	patches.ApplyMethodReturn(newMockMetaBackupPlan().client, "UsePostgresqlClient", pgClient)

	patches.ApplyMethodFunc(pgClient, "CreateBackupPlanWithContext", func(ctx context.Context, request *postgresql.CreateBackupPlanRequest) (*postgresql.CreateBackupPlanResponse, error) {
		assert.NotNil(t, request.DBInstanceId)
		assert.Equal(t, "postgres-test123", *request.DBInstanceId)
		assert.NotNil(t, request.PlanName)
		assert.Equal(t, "tf-example-plan", *request.PlanName)
		assert.NotNil(t, request.BackupPeriodType)
		assert.Equal(t, "month", *request.BackupPeriodType)
		assert.NotNil(t, request.BackupPeriod)
		assert.Equal(t, 3, len(request.BackupPeriod))
		assert.NotNil(t, request.MinBackupStartTime)
		assert.Equal(t, "01:00:00", *request.MinBackupStartTime)
		assert.NotNil(t, request.MaxBackupStartTime)
		assert.Equal(t, "02:00:00", *request.MaxBackupStartTime)
		assert.NotNil(t, request.BaseBackupRetentionPeriod)
		assert.Equal(t, uint64(30), *request.BaseBackupRetentionPeriod)

		resp := postgresql.NewCreateBackupPlanResponse()
		resp.Response = &postgresql.CreateBackupPlanResponseParams{
			PlanId:    ptrStringBackupPlan("plan-xxxx"),
			RequestId: ptrStringBackupPlan("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(pgClient, "DescribeBackupPlans", func(request *postgresql.DescribeBackupPlansRequest) (*postgresql.DescribeBackupPlansResponse, error) {
		assert.NotNil(t, request.DBInstanceId)
		assert.Equal(t, "postgres-test123", *request.DBInstanceId)

		backupPeriodJson, _ := json.Marshal([]string{"1", "2", "15"})
		resp := postgresql.NewDescribeBackupPlansResponse()
		resp.Response = &postgresql.DescribeBackupPlansResponseParams{
			Plans: []*postgresql.BackupPlan{
				{
					PlanId:                    ptrStringBackupPlan("plan-xxxx"),
					PlanName:                  ptrStringBackupPlan("tf-example-plan"),
					BackupPeriodType:          ptrStringBackupPlan("month"),
					BackupPeriod:              ptrStringBackupPlan(string(backupPeriodJson)),
					MinBackupStartTime:        ptrStringBackupPlan("01:00:00"),
					MaxBackupStartTime:        ptrStringBackupPlan("02:00:00"),
					BaseBackupRetentionPeriod: ptrUint64BackupPlan(30),
					LogBackupRetentionPeriod:  ptrUint64BackupPlan(7),
					BackupMethod:              ptrStringBackupPlan("physical"),
				},
			},
			RequestId: ptrStringBackupPlan("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaBackupPlan()
	res := svcpostgresql.ResourceTencentCloudPostgresqlBackupPlan()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"db_instance_id":               "postgres-test123",
		"plan_name":                    "tf-example-plan",
		"backup_period_type":           "month",
		"backup_period":                []interface{}{"1", "2", "15"},
		"min_backup_start_time":        "01:00:00",
		"max_backup_start_time":        "02:00:00",
		"base_backup_retention_period": 30,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "postgres-test123#plan-xxxx", d.Id())
	assert.Equal(t, "plan-xxxx", d.Get("plan_id").(string))
	assert.Equal(t, "tf-example-plan", d.Get("plan_name").(string))
	assert.Equal(t, "month", d.Get("backup_period_type").(string))
}

func TestPostgresqlBackupPlan_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	pgClient := &postgresql.Client{}
	patches.ApplyMethodReturn(newMockMetaBackupPlan().client, "UsePostgresqlClient", pgClient)

	patches.ApplyMethodFunc(pgClient, "DescribeBackupPlans", func(request *postgresql.DescribeBackupPlansRequest) (*postgresql.DescribeBackupPlansResponse, error) {
		assert.NotNil(t, request.DBInstanceId)
		assert.Equal(t, "postgres-test123", *request.DBInstanceId)

		backupPeriodJson, _ := json.Marshal([]string{"1", "2", "15"})
		resp := postgresql.NewDescribeBackupPlansResponse()
		resp.Response = &postgresql.DescribeBackupPlansResponseParams{
			Plans: []*postgresql.BackupPlan{
				{
					PlanId:                    ptrStringBackupPlan("plan-xxxx"),
					PlanName:                  ptrStringBackupPlan("tf-example-plan"),
					BackupPeriodType:          ptrStringBackupPlan("month"),
					BackupPeriod:              ptrStringBackupPlan(string(backupPeriodJson)),
					MinBackupStartTime:        ptrStringBackupPlan("01:00:00"),
					MaxBackupStartTime:        ptrStringBackupPlan("02:00:00"),
					BaseBackupRetentionPeriod: ptrUint64BackupPlan(30),
					LogBackupRetentionPeriod:  ptrUint64BackupPlan(7),
					BackupMethod:              ptrStringBackupPlan("physical"),
				},
			},
			RequestId: ptrStringBackupPlan("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaBackupPlan()
	res := svcpostgresql.ResourceTencentCloudPostgresqlBackupPlan()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"db_instance_id":     "postgres-test123",
		"plan_name":          "tf-example-plan",
		"backup_period_type": "month",
		"backup_period":      []interface{}{"1", "2", "15"},
	})
	d.SetId("postgres-test123#plan-xxxx")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "postgres-test123#plan-xxxx", d.Id())
	assert.Equal(t, "plan-xxxx", d.Get("plan_id").(string))
	assert.Equal(t, "tf-example-plan", d.Get("plan_name").(string))
	assert.Equal(t, "month", d.Get("backup_period_type").(string))
	assert.Equal(t, "01:00:00", d.Get("min_backup_start_time").(string))
	assert.Equal(t, "02:00:00", d.Get("max_backup_start_time").(string))
	assert.Equal(t, 30, d.Get("base_backup_retention_period").(int))
	assert.Equal(t, 7, d.Get("log_backup_retention_period").(int))
	assert.Equal(t, "physical", d.Get("backup_method").(string))
}

func TestPostgresqlBackupPlan_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	pgClient := &postgresql.Client{}
	patches.ApplyMethodReturn(newMockMetaBackupPlan().client, "UsePostgresqlClient", pgClient)

	patches.ApplyMethodFunc(pgClient, "DescribeBackupPlans", func(request *postgresql.DescribeBackupPlansRequest) (*postgresql.DescribeBackupPlansResponse, error) {
		resp := postgresql.NewDescribeBackupPlansResponse()
		resp.Response = &postgresql.DescribeBackupPlansResponseParams{
			Plans:     []*postgresql.BackupPlan{},
			RequestId: ptrStringBackupPlan("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaBackupPlan()
	res := svcpostgresql.ResourceTencentCloudPostgresqlBackupPlan()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"db_instance_id":     "postgres-test123",
		"plan_name":          "tf-example-plan",
		"backup_period_type": "month",
		"backup_period":      []interface{}{"1", "2", "15"},
	})
	d.SetId("postgres-test123#plan-xxxx")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestPostgresqlBackupPlan_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	pgClient := &postgresql.Client{}
	patches.ApplyMethodReturn(newMockMetaBackupPlan().client, "UsePostgresqlClient", pgClient)

	patches.ApplyMethodFunc(pgClient, "ModifyBackupPlanWithContext", func(ctx context.Context, request *postgresql.ModifyBackupPlanRequest) (*postgresql.ModifyBackupPlanResponse, error) {
		assert.NotNil(t, request.DBInstanceId)
		assert.Equal(t, "postgres-test123", *request.DBInstanceId)
		assert.NotNil(t, request.PlanId)
		assert.Equal(t, "plan-xxxx", *request.PlanId)
		assert.NotNil(t, request.PlanName)
		assert.Equal(t, "tf-example-plan-updated", *request.PlanName)
		assert.NotNil(t, request.BackupMethod)
		assert.Equal(t, "logical", *request.BackupMethod)
		assert.NotNil(t, request.LogBackupRetentionPeriod)
		assert.Equal(t, uint64(14), *request.LogBackupRetentionPeriod)

		resp := postgresql.NewModifyBackupPlanResponse()
		resp.Response = &postgresql.ModifyBackupPlanResponseParams{
			RequestId: ptrStringBackupPlan("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(pgClient, "DescribeBackupPlans", func(request *postgresql.DescribeBackupPlansRequest) (*postgresql.DescribeBackupPlansResponse, error) {
		backupPeriodJson, _ := json.Marshal([]string{"1", "2", "15"})
		resp := postgresql.NewDescribeBackupPlansResponse()
		resp.Response = &postgresql.DescribeBackupPlansResponseParams{
			Plans: []*postgresql.BackupPlan{
				{
					PlanId:                    ptrStringBackupPlan("plan-xxxx"),
					PlanName:                  ptrStringBackupPlan("tf-example-plan-updated"),
					BackupPeriodType:          ptrStringBackupPlan("month"),
					BackupPeriod:              ptrStringBackupPlan(string(backupPeriodJson)),
					MinBackupStartTime:        ptrStringBackupPlan("01:00:00"),
					MaxBackupStartTime:        ptrStringBackupPlan("02:00:00"),
					BaseBackupRetentionPeriod: ptrUint64BackupPlan(30),
					LogBackupRetentionPeriod:  ptrUint64BackupPlan(14),
					BackupMethod:              ptrStringBackupPlan("logical"),
				},
			},
			RequestId: ptrStringBackupPlan("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaBackupPlan()
	res := svcpostgresql.ResourceTencentCloudPostgresqlBackupPlan()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"db_instance_id":               "postgres-test123",
		"plan_name":                    "tf-example-plan-updated",
		"backup_period_type":           "month",
		"backup_period":                []interface{}{"1", "2", "15"},
		"min_backup_start_time":        "01:00:00",
		"max_backup_start_time":        "02:00:00",
		"base_backup_retention_period": 30,
		"log_backup_retention_period":  14,
		"backup_method":                "logical",
	})
	d.SetId("postgres-test123#plan-xxxx")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "postgres-test123#plan-xxxx", d.Id())
	assert.Equal(t, "tf-example-plan-updated", d.Get("plan_name").(string))
	assert.Equal(t, "logical", d.Get("backup_method").(string))
	assert.Equal(t, 14, d.Get("log_backup_retention_period").(int))
}

func TestPostgresqlBackupPlan_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	pgClient := &postgresql.Client{}
	patches.ApplyMethodReturn(newMockMetaBackupPlan().client, "UsePostgresqlClient", pgClient)

	patches.ApplyMethodFunc(pgClient, "DeleteBackupPlanWithContext", func(ctx context.Context, request *postgresql.DeleteBackupPlanRequest) (*postgresql.DeleteBackupPlanResponse, error) {
		assert.NotNil(t, request.DBInstanceId)
		assert.Equal(t, "postgres-test123", *request.DBInstanceId)
		assert.NotNil(t, request.PlanId)
		assert.Equal(t, "plan-xxxx", *request.PlanId)

		resp := postgresql.NewDeleteBackupPlanResponse()
		resp.Response = &postgresql.DeleteBackupPlanResponseParams{
			RequestId: ptrStringBackupPlan("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaBackupPlan()
	res := svcpostgresql.ResourceTencentCloudPostgresqlBackupPlan()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"db_instance_id":     "postgres-test123",
		"plan_name":          "tf-example-plan",
		"backup_period_type": "month",
		"backup_period":      []interface{}{"1", "2", "15"},
	})
	d.SetId("postgres-test123#plan-xxxx")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}
