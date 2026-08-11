package postgresql_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	tag "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcpostgresql "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/postgresql"
)

type mockMetaReadonlyInstanceV2 struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaReadonlyInstanceV2) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaReadonlyInstanceV2{}

func newMockMetaReadonlyInstanceV2() *mockMetaReadonlyInstanceV2 {
	return &mockMetaReadonlyInstanceV2{client: &connectivity.TencentCloudClient{}}
}

func ptrStringROV2(s string) *string {
	return &s
}

func ptrUint64ROV2(i uint64) *uint64 {
	return &i
}

func ptrBoolROV2(b bool) *bool {
	return &b
}

// go test ./tencentcloud/services/postgresql/ -run "TestPostgresqlReadonlyInstanceV2" -v -count=1 -gcflags="all=-l"

func TestPostgresqlReadonlyInstanceV2_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	pgClient := &postgresql.Client{}
	patches.ApplyMethodReturn(newMockMetaReadonlyInstanceV2().client, "UsePostgresqlClient", pgClient)

	tagClient := &tag.Client{}
	patches.ApplyMethodReturn(newMockMetaReadonlyInstanceV2().client, "UseTagClient", tagClient)

	patches.ApplyMethodFunc(pgClient, "CreateReadOnlyDBInstance", func(request *postgresql.CreateReadOnlyDBInstanceRequest) (*postgresql.CreateReadOnlyDBInstanceResponse, error) {
		assert.NotNil(t, request.Zone)
		assert.Equal(t, "ap-guangzhou-3", *request.Zone)
		assert.NotNil(t, request.MasterDBInstanceId)
		assert.Equal(t, "pgs-xxxx", *request.MasterDBInstanceId)
		assert.NotNil(t, request.SpecCode)
		assert.Equal(t, "pgro.c2.large.ha", *request.SpecCode)
		assert.NotNil(t, request.Storage)
		assert.Equal(t, uint64(250), *request.Storage)
		assert.NotNil(t, request.InstanceCount)
		assert.Equal(t, uint64(1), *request.InstanceCount)
		assert.NotNil(t, request.Period)
		assert.Equal(t, uint64(1), *request.Period)

		resp := postgresql.NewCreateReadOnlyDBInstanceResponse()
		resp.Response = &postgresql.CreateReadOnlyDBInstanceResponseParams{
			DealNames:         []*string{ptrStringROV2("2023xxxx")},
			BillId:            ptrStringROV2("bill-xxxx"),
			DBInstanceIdSet:   []*string{ptrStringROV2("pgro-xxxx")},
			BillingParameters: ptrStringROV2("{}"),
			RequestId:         ptrStringROV2("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(pgClient, "DescribeDBInstanceAttribute", func(request *postgresql.DescribeDBInstanceAttributeRequest) (*postgresql.DescribeDBInstanceAttributeResponse, error) {
		assert.NotNil(t, request.DBInstanceId)
		assert.Equal(t, "pgro-xxxx", *request.DBInstanceId)
		resp := postgresql.NewDescribeDBInstanceAttributeResponse()
		resp.Response = &postgresql.DescribeDBInstanceAttributeResponseParams{
			DBInstance: &postgresql.DBInstance{
				DBInstanceId:       ptrStringROV2("pgro-xxxx"),
				DBInstanceName:     ptrStringROV2("tf_ro_instance_v2"),
				DBInstanceStatus:   ptrStringROV2("running"),
				Zone:               ptrStringROV2("ap-guangzhou-3"),
				VpcId:              ptrStringROV2("vpc-xxxx"),
				SubnetId:           ptrStringROV2("subnet-xxxx"),
				DBVersion:          ptrStringROV2("15.1"),
				DBInstanceStorage:  ptrUint64ROV2(250),
				ProjectId:          ptrUint64ROV2(0),
				SupportIpv6:        ptrUint64ROV2(0),
				DeletionProtection: ptrBoolROV2(false),
				MasterDBInstanceId: ptrStringROV2("pgs-xxxx"),
				AutoRenew:          ptrUint64ROV2(0),
				PayType:            ptrStringROV2("postpaid"),
			},
			RequestId: ptrStringROV2("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(pgClient, "DescribeReadOnlyGroups", func(request *postgresql.DescribeReadOnlyGroupsRequest) (*postgresql.DescribeReadOnlyGroupsResponse, error) {
		resp := postgresql.NewDescribeReadOnlyGroupsResponse()
		resp.Response = &postgresql.DescribeReadOnlyGroupsResponseParams{
			ReadOnlyGroupList: []*postgresql.ReadOnlyGroup{},
			RequestId:         ptrStringROV2("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(pgClient, "DescribeDBInstanceSecurityGroups", func(request *postgresql.DescribeDBInstanceSecurityGroupsRequest) (*postgresql.DescribeDBInstanceSecurityGroupsResponse, error) {
		resp := postgresql.NewDescribeDBInstanceSecurityGroupsResponse()
		resp.Response = &postgresql.DescribeDBInstanceSecurityGroupsResponseParams{
			SecurityGroupSet: []*postgresql.SecurityGroup{},
			RequestId:        ptrStringROV2("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(tagClient, "DescribeResourceTagsByResourceIds", func(request *tag.DescribeResourceTagsByResourceIdsRequest) (*tag.DescribeResourceTagsByResourceIdsResponse, error) {
		resp := tag.NewDescribeResourceTagsByResourceIdsResponse()
		resp.Response = &tag.DescribeResourceTagsByResourceIdsResponseParams{
			Tags:      []*tag.TagResource{},
			RequestId: ptrStringROV2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaReadonlyInstanceV2()
	res := svcpostgresql.ResourceTencentCloudPostgresqlReadonlyInstanceV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone":                  "ap-guangzhou-3",
		"master_db_instance_id": "pgs-xxxx",
		"spec_code":             "pgro.c2.large.ha",
		"storage":               250,
		"instance_count":        1,
		"period":                1,
		"instance_charge_type":  "POSTPAID_BY_HOUR",
		"vpc_id":                "vpc-xxxx",
		"subnet_id":             "subnet-xxxx",
		"name":                  "tf_ro_instance_v2",
		"auto_renew_flag":       0,
		"need_support_ipv6":     0,
		"project_id":            0,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "pgro-xxxx", d.Id())
	assert.Equal(t, "pgro-xxxx", d.Get("db_instance_id").(string))
	assert.Equal(t, "bill-xxxx", d.Get("bill_id").(string))
	assert.Equal(t, 1, d.Get("db_instance_id_set.#").(int))
	assert.Equal(t, "ap-guangzhou-3", d.Get("zone").(string))
}

func TestPostgresqlReadonlyInstanceV2_Create_EmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	pgClient := &postgresql.Client{}
	patches.ApplyMethodReturn(newMockMetaReadonlyInstanceV2().client, "UsePostgresqlClient", pgClient)

	patches.ApplyMethodFunc(pgClient, "CreateReadOnlyDBInstance", func(request *postgresql.CreateReadOnlyDBInstanceRequest) (*postgresql.CreateReadOnlyDBInstanceResponse, error) {
		resp := postgresql.NewCreateReadOnlyDBInstanceResponse()
		resp.Response = &postgresql.CreateReadOnlyDBInstanceResponseParams{
			DealNames:         []*string{ptrStringROV2("2023xxxx")},
			BillId:            ptrStringROV2("bill-xxxx"),
			DBInstanceIdSet:   []*string{},
			BillingParameters: ptrStringROV2("{}"),
			RequestId:         ptrStringROV2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaReadonlyInstanceV2()
	res := svcpostgresql.ResourceTencentCloudPostgresqlReadonlyInstanceV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone":                  "ap-guangzhou-3",
		"master_db_instance_id": "pgs-xxxx",
		"spec_code":             "pgro.c2.large.ha",
		"storage":               250,
		"instance_count":        1,
		"period":                1,
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Equal(t, "", d.Id())
}

func TestPostgresqlReadonlyInstanceV2_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	pgClient := &postgresql.Client{}
	patches.ApplyMethodReturn(newMockMetaReadonlyInstanceV2().client, "UsePostgresqlClient", pgClient)

	tagClient := &tag.Client{}
	patches.ApplyMethodReturn(newMockMetaReadonlyInstanceV2().client, "UseTagClient", tagClient)

	patches.ApplyMethodFunc(pgClient, "DescribeDBInstanceAttribute", func(request *postgresql.DescribeDBInstanceAttributeRequest) (*postgresql.DescribeDBInstanceAttributeResponse, error) {
		assert.NotNil(t, request.DBInstanceId)
		assert.Equal(t, "pgro-xxxx", *request.DBInstanceId)
		resp := postgresql.NewDescribeDBInstanceAttributeResponse()
		resp.Response = &postgresql.DescribeDBInstanceAttributeResponseParams{
			DBInstance: &postgresql.DBInstance{
				DBInstanceId:       ptrStringROV2("pgro-xxxx"),
				DBInstanceName:     ptrStringROV2("tf_ro_instance_v2"),
				DBInstanceStatus:   ptrStringROV2("running"),
				Zone:               ptrStringROV2("ap-guangzhou-3"),
				VpcId:              ptrStringROV2("vpc-xxxx"),
				SubnetId:           ptrStringROV2("subnet-xxxx"),
				DBVersion:          ptrStringROV2("15.1"),
				DBInstanceStorage:  ptrUint64ROV2(250),
				ProjectId:          ptrUint64ROV2(0),
				SupportIpv6:        ptrUint64ROV2(0),
				DeletionProtection: ptrBoolROV2(false),
				MasterDBInstanceId: ptrStringROV2("pgs-xxxx"),
				AutoRenew:          ptrUint64ROV2(0),
				PayType:            ptrStringROV2("postpaid"),
			},
			RequestId: ptrStringROV2("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(pgClient, "DescribeReadOnlyGroups", func(request *postgresql.DescribeReadOnlyGroupsRequest) (*postgresql.DescribeReadOnlyGroupsResponse, error) {
		resp := postgresql.NewDescribeReadOnlyGroupsResponse()
		resp.Response = &postgresql.DescribeReadOnlyGroupsResponseParams{
			ReadOnlyGroupList: []*postgresql.ReadOnlyGroup{},
			RequestId:         ptrStringROV2("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(pgClient, "DescribeDBInstanceSecurityGroups", func(request *postgresql.DescribeDBInstanceSecurityGroupsRequest) (*postgresql.DescribeDBInstanceSecurityGroupsResponse, error) {
		resp := postgresql.NewDescribeDBInstanceSecurityGroupsResponse()
		resp.Response = &postgresql.DescribeDBInstanceSecurityGroupsResponseParams{
			SecurityGroupSet: []*postgresql.SecurityGroup{},
			RequestId:        ptrStringROV2("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(tagClient, "DescribeResourceTagsByResourceIds", func(request *tag.DescribeResourceTagsByResourceIdsRequest) (*tag.DescribeResourceTagsByResourceIdsResponse, error) {
		resp := tag.NewDescribeResourceTagsByResourceIdsResponse()
		resp.Response = &tag.DescribeResourceTagsByResourceIdsResponseParams{
			Tags:      []*tag.TagResource{},
			RequestId: ptrStringROV2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaReadonlyInstanceV2()
	res := svcpostgresql.ResourceTencentCloudPostgresqlReadonlyInstanceV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("pgro-xxxx")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "pgro-xxxx", d.Id())
	assert.Equal(t, "ap-guangzhou-3", d.Get("zone").(string))
	assert.Equal(t, "pgs-xxxx", d.Get("master_db_instance_id").(string))
	assert.Equal(t, 250, d.Get("storage").(int))
	assert.Equal(t, "tf_ro_instance_v2", d.Get("name").(string))
	assert.Equal(t, "pgro-xxxx", d.Get("db_instance_id").(string))
	assert.Equal(t, false, d.Get("deletion_protection").(bool))
}

func TestPostgresqlReadonlyInstanceV2_Read_Empty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	pgClient := &postgresql.Client{}
	patches.ApplyMethodReturn(newMockMetaReadonlyInstanceV2().client, "UsePostgresqlClient", pgClient)

	patches.ApplyMethodFunc(pgClient, "DescribeDBInstanceAttribute", func(request *postgresql.DescribeDBInstanceAttributeRequest) (*postgresql.DescribeDBInstanceAttributeResponse, error) {
		resp := postgresql.NewDescribeDBInstanceAttributeResponse()
		resp.Response = &postgresql.DescribeDBInstanceAttributeResponseParams{
			DBInstance: nil,
			RequestId:  ptrStringROV2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaReadonlyInstanceV2()
	res := svcpostgresql.ResourceTencentCloudPostgresqlReadonlyInstanceV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("pgro-xxxx")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestPostgresqlReadonlyInstanceV2_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	pgClient := &postgresql.Client{}
	patches.ApplyMethodReturn(newMockMetaReadonlyInstanceV2().client, "UsePostgresqlClient", pgClient)

	isolated := false
	patches.ApplyMethodFunc(pgClient, "IsolateDBInstances", func(request *postgresql.IsolateDBInstancesRequest) (*postgresql.IsolateDBInstancesResponse, error) {
		assert.NotNil(t, request.DBInstanceIdSet)
		assert.Equal(t, 1, len(request.DBInstanceIdSet))
		assert.Equal(t, "pgro-xxxx", *request.DBInstanceIdSet[0])
		isolated = true
		resp := postgresql.NewIsolateDBInstancesResponse()
		resp.Response = &postgresql.IsolateDBInstancesResponseParams{
			RequestId: ptrStringROV2("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(pgClient, "DescribeDBInstanceAttribute", func(request *postgresql.DescribeDBInstanceAttributeRequest) (*postgresql.DescribeDBInstanceAttributeResponse, error) {
		resp := postgresql.NewDescribeDBInstanceAttributeResponse()
		resp.Response = &postgresql.DescribeDBInstanceAttributeResponseParams{
			DBInstance: &postgresql.DBInstance{
				DBInstanceId:     ptrStringROV2("pgro-xxxx"),
				DBInstanceStatus: ptrStringROV2("isolated"),
			},
			RequestId: ptrStringROV2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaReadonlyInstanceV2()
	res := svcpostgresql.ResourceTencentCloudPostgresqlReadonlyInstanceV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("pgro-xxxx")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.True(t, isolated)
}
