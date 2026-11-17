package postgresql_test

import (
	"log"
	"os"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcpostgresql "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/postgresql"
)

type mockMetaDbInstanceSgDS struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaDbInstanceSgDS) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaDbInstanceSgDS{}

func newMockMetaDbInstanceSgDS() *mockMetaDbInstanceSgDS {
	return &mockMetaDbInstanceSgDS{client: &connectivity.TencentCloudClient{}}
}

func ptrStrSg(s string) *string {
	return &s
}

func ptrInt64Sg(i int64) *int64 {
	return &i
}

// go test ./tencentcloud/services/postgresql/ -run "TestPostgresqlDbInstanceSecurityGroupsDS" -v -count=1 -gcflags="all=-l"

func TestPostgresqlDbInstanceSecurityGroupsDS_ReadByDbInstanceId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	pgClient := &postgresql.Client{}
	patches.ApplyMethodReturn(newMockMetaDbInstanceSgDS().client, "UsePostgresqlClient", pgClient)

	patches.ApplyMethodFunc(pgClient, "DescribeDBInstanceSecurityGroups", func(request *postgresql.DescribeDBInstanceSecurityGroupsRequest) (*postgresql.DescribeDBInstanceSecurityGroupsResponse, error) {
		assert.NotNil(t, request.DBInstanceId)
		assert.Equal(t, "postgres-gzg9jb2n", *request.DBInstanceId)
		assert.True(t, request.ReadOnlyGroupId == nil)

		resp := postgresql.NewDescribeDBInstanceSecurityGroupsResponse()
		resp.Response = &postgresql.DescribeDBInstanceSecurityGroupsResponseParams{
			SecurityGroupSet: []*postgresql.SecurityGroup{
				{
					ProjectId:                ptrInt64Sg(0),
					CreateTime:               ptrStrSg("2024-01-01 00:00:00"),
					SecurityGroupId:          ptrStrSg("sg-1"),
					SecurityGroupName:        ptrStrSg("test-sg-1"),
					SecurityGroupDescription: ptrStrSg("test sg 1"),
					Inbound: []*postgresql.PolicyRule{
						{
							Action:      ptrStrSg("ACCEPT"),
							CidrIp:      ptrStrSg("10.0.0.0/16"),
							PortRange:   ptrStrSg("5432"),
							IpProtocol:  ptrStrSg("TCP"),
							Description: ptrStrSg("allow tcp 5432"),
						},
					},
					Outbound: []*postgresql.PolicyRule{
						{
							Action:      ptrStrSg("DROP"),
							CidrIp:      ptrStrSg("0.0.0.0/0"),
							PortRange:   ptrStrSg("all"),
							IpProtocol:  ptrStrSg("UDP"),
							Description: ptrStrSg("deny udp all"),
						},
					},
				},
				{
					ProjectId:         ptrInt64Sg(0),
					SecurityGroupId:   ptrStrSg("sg-2"),
					SecurityGroupName: ptrStrSg("test-sg-2"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaDbInstanceSgDS()
	res := svcpostgresql.DataSourceTencentCloudPostgresqlDbInstanceSecurityGroups()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"db_instance_id": "postgres-gzg9jb2n",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	securityGroupSet := d.Get("security_group_set").([]interface{})
	assert.Len(t, securityGroupSet, 2)

	sg0 := securityGroupSet[0].(map[string]interface{})
	assert.Equal(t, 0, sg0["project_id"].(int))
	assert.Equal(t, "2024-01-01 00:00:00", sg0["create_time"].(string))
	assert.Equal(t, "sg-1", sg0["security_group_id"].(string))
	assert.Equal(t, "test-sg-1", sg0["security_group_name"].(string))
	assert.Equal(t, "test sg 1", sg0["security_group_description"].(string))

	inbound0 := sg0["inbound"].([]interface{})
	assert.Len(t, inbound0, 1)
	inbound0Map := inbound0[0].(map[string]interface{})
	assert.Equal(t, "ACCEPT", inbound0Map["action"].(string))
	assert.Equal(t, "10.0.0.0/16", inbound0Map["cidr_ip"].(string))
	assert.Equal(t, "5432", inbound0Map["port_range"].(string))
	assert.Equal(t, "TCP", inbound0Map["ip_protocol"].(string))
	assert.Equal(t, "allow tcp 5432", inbound0Map["description"].(string))

	outbound0 := sg0["outbound"].([]interface{})
	assert.Len(t, outbound0, 1)
	outbound0Map := outbound0[0].(map[string]interface{})
	assert.Equal(t, "DROP", outbound0Map["action"].(string))
	assert.Equal(t, "0.0.0.0/0", outbound0Map["cidr_ip"].(string))
	assert.Equal(t, "all", outbound0Map["port_range"].(string))
	assert.Equal(t, "UDP", outbound0Map["ip_protocol"].(string))
	assert.Equal(t, "deny udp all", outbound0Map["description"].(string))

	sg1 := securityGroupSet[1].(map[string]interface{})
	assert.Equal(t, "sg-2", sg1["security_group_id"].(string))
	assert.Equal(t, "test-sg-2", sg1["security_group_name"].(string))
}

func TestPostgresqlDbInstanceSecurityGroupsDS_ReadByReadOnlyGroupId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	pgClient := &postgresql.Client{}
	patches.ApplyMethodReturn(newMockMetaDbInstanceSgDS().client, "UsePostgresqlClient", pgClient)

	patches.ApplyMethodFunc(pgClient, "DescribeDBInstanceSecurityGroups", func(request *postgresql.DescribeDBInstanceSecurityGroupsRequest) (*postgresql.DescribeDBInstanceSecurityGroupsResponse, error) {
		assert.True(t, request.DBInstanceId == nil)
		assert.NotNil(t, request.ReadOnlyGroupId)
		assert.Equal(t, "pgro-xxxxx", *request.ReadOnlyGroupId)

		resp := postgresql.NewDescribeDBInstanceSecurityGroupsResponse()
		resp.Response = &postgresql.DescribeDBInstanceSecurityGroupsResponseParams{
			SecurityGroupSet: []*postgresql.SecurityGroup{
				{
					ProjectId:         ptrInt64Sg(0),
					SecurityGroupId:   ptrStrSg("sg-ro-1"),
					SecurityGroupName: ptrStrSg("test-sg-ro-1"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaDbInstanceSgDS()
	res := svcpostgresql.DataSourceTencentCloudPostgresqlDbInstanceSecurityGroups()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"read_only_group_id": "pgro-xxxxx",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	securityGroupSet := d.Get("security_group_set").([]interface{})
	assert.Len(t, securityGroupSet, 1)

	sg0 := securityGroupSet[0].(map[string]interface{})
	assert.Equal(t, "sg-ro-1", sg0["security_group_id"].(string))
	assert.Equal(t, "test-sg-ro-1", sg0["security_group_name"].(string))
}

func TestPostgresqlDbInstanceSecurityGroupsDS_ReadWithNilFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	pgClient := &postgresql.Client{}
	patches.ApplyMethodReturn(newMockMetaDbInstanceSgDS().client, "UsePostgresqlClient", pgClient)

	patches.ApplyMethodFunc(pgClient, "DescribeDBInstanceSecurityGroups", func(request *postgresql.DescribeDBInstanceSecurityGroupsRequest) (*postgresql.DescribeDBInstanceSecurityGroupsResponse, error) {
		resp := postgresql.NewDescribeDBInstanceSecurityGroupsResponse()
		resp.Response = &postgresql.DescribeDBInstanceSecurityGroupsResponseParams{
			SecurityGroupSet: []*postgresql.SecurityGroup{
				{
					SecurityGroupId: ptrStrSg("sg-nil"),
					// Inbound and Outbound are nil
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaDbInstanceSgDS()
	res := svcpostgresql.DataSourceTencentCloudPostgresqlDbInstanceSecurityGroups()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"db_instance_id": "postgres-gzg9jb2n",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	securityGroupSet := d.Get("security_group_set").([]interface{})
	assert.Len(t, securityGroupSet, 1)

	sg0 := securityGroupSet[0].(map[string]interface{})
	assert.Equal(t, "sg-nil", sg0["security_group_id"].(string))

	// Inbound is nil, should have empty list
	inboundField := sg0["inbound"]
	if inboundField != nil {
		inboundList := inboundField.([]interface{})
		assert.Len(t, inboundList, 0)
	}

	// Outbound is nil, should have empty list
	outboundField := sg0["outbound"]
	if outboundField != nil {
		outboundList := outboundField.([]interface{})
		assert.Len(t, outboundList, 0)
	}
}

func TestPostgresqlDbInstanceSecurityGroupsDS_ReadWithEmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	pgClient := &postgresql.Client{}
	patches.ApplyMethodReturn(newMockMetaDbInstanceSgDS().client, "UsePostgresqlClient", pgClient)

	patches.ApplyMethodFunc(pgClient, "DescribeDBInstanceSecurityGroups", func(request *postgresql.DescribeDBInstanceSecurityGroupsRequest) (*postgresql.DescribeDBInstanceSecurityGroupsResponse, error) {
		resp := postgresql.NewDescribeDBInstanceSecurityGroupsResponse()
		resp.Response = &postgresql.DescribeDBInstanceSecurityGroupsResponseParams{
			SecurityGroupSet: []*postgresql.SecurityGroup{},
		}
		return resp, nil
	})

	meta := newMockMetaDbInstanceSgDS()
	res := svcpostgresql.DataSourceTencentCloudPostgresqlDbInstanceSecurityGroups()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"db_instance_id": "postgres-gzg9jb2n",
	})

	err := res.Read(d, meta)
	// When response is empty, the datasource should return an error (NonRetryableError)
	assert.Error(t, err)
}

func TestPostgresqlDbInstanceSecurityGroupsDS_ReadWithResultOutputFile(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	pgClient := &postgresql.Client{}
	patches.ApplyMethodReturn(newMockMetaDbInstanceSgDS().client, "UsePostgresqlClient", pgClient)

	patches.ApplyMethodFunc(pgClient, "DescribeDBInstanceSecurityGroups", func(request *postgresql.DescribeDBInstanceSecurityGroupsRequest) (*postgresql.DescribeDBInstanceSecurityGroupsResponse, error) {
		resp := postgresql.NewDescribeDBInstanceSecurityGroupsResponse()
		resp.Response = &postgresql.DescribeDBInstanceSecurityGroupsResponseParams{
			SecurityGroupSet: []*postgresql.SecurityGroup{
				{
					ProjectId:         ptrInt64Sg(0),
					SecurityGroupId:   ptrStrSg("sg-file"),
					SecurityGroupName: ptrStrSg("test-sg-file"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaDbInstanceSgDS()
	res := svcpostgresql.DataSourceTencentCloudPostgresqlDbInstanceSecurityGroups()
	outputFile := "./test_output_sg.json"
	defer func() {
		if err := os.Remove(outputFile); err != nil {
			log.Printf("failed to remove %s: %v", outputFile, err)
		}
	}()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"db_instance_id":     "postgres-gzg9jb2n",
		"result_output_file": outputFile,
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	_, e := os.Stat(outputFile)
	assert.NoError(t, e)
}

func TestPostgresqlDbInstanceSecurityGroupsDS_Schema(t *testing.T) {
	res := svcpostgresql.DataSourceTencentCloudPostgresqlDbInstanceSecurityGroups()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "db_instance_id")
	assert.Contains(t, res.Schema, "read_only_group_id")
	assert.Contains(t, res.Schema, "security_group_set")
	assert.Contains(t, res.Schema, "result_output_file")

	dbInstanceIdSchema := res.Schema["db_instance_id"]
	assert.Equal(t, schema.TypeString, dbInstanceIdSchema.Type)
	assert.True(t, dbInstanceIdSchema.Optional)

	readOnlyGroupIdSchema := res.Schema["read_only_group_id"]
	assert.Equal(t, schema.TypeString, readOnlyGroupIdSchema.Type)
	assert.True(t, readOnlyGroupIdSchema.Optional)

	securityGroupSetSchema := res.Schema["security_group_set"]
	assert.Equal(t, schema.TypeList, securityGroupSetSchema.Type)
	assert.True(t, securityGroupSetSchema.Computed)

	elemRes := securityGroupSetSchema.Elem.(*schema.Resource)
	assert.Contains(t, elemRes.Schema, "project_id")
	assert.Contains(t, elemRes.Schema, "create_time")
	assert.Contains(t, elemRes.Schema, "security_group_id")
	assert.Contains(t, elemRes.Schema, "security_group_name")
	assert.Contains(t, elemRes.Schema, "security_group_description")
	assert.Contains(t, elemRes.Schema, "inbound")
	assert.Contains(t, elemRes.Schema, "outbound")

	inboundRes := elemRes.Schema["inbound"].Elem.(*schema.Resource)
	assert.Contains(t, inboundRes.Schema, "action")
	assert.Contains(t, inboundRes.Schema, "cidr_ip")
	assert.Contains(t, inboundRes.Schema, "port_range")
	assert.Contains(t, inboundRes.Schema, "ip_protocol")
	assert.Contains(t, inboundRes.Schema, "description")

	outboundRes := elemRes.Schema["outbound"].Elem.(*schema.Resource)
	assert.Contains(t, outboundRes.Schema, "action")
	assert.Contains(t, outboundRes.Schema, "cidr_ip")
	assert.Contains(t, outboundRes.Schema, "port_range")
	assert.Contains(t, outboundRes.Schema, "ip_protocol")
	assert.Contains(t, outboundRes.Schema, "description")
}
