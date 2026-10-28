package dbdc_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	dbdcv20201029 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dbdc"
)

// go test ./tencentcloud/services/dbdc/ -run "TestDbdcDbCustomNodeSecurityGroupsDS" -v -count=1 -gcflags="all=-l"

func TestDbdcDbCustomNodeSecurityGroupsDS_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdcDS().client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomNodeSecurityGroups", func(request *dbdcv20201029.DescribeDBCustomNodeSecurityGroupsRequest) (*dbdcv20201029.DescribeDBCustomNodeSecurityGroupsResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomNodeSecurityGroupsResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomNodeSecurityGroupsResponseParams{
			Groups: []*dbdcv20201029.SecurityGroup{
				{
					SecurityGroupId:     ptrStr("sg-abc123"),
					SecurityGroupName:   ptrStr("test-sg-1"),
					SecurityGroupRemark: ptrStr("test security group 1"),
					ProjectId:           ptrInt64(0),
					CreateTime:          ptrStr("2024-01-01 00:00:00"),
					Inbound: []*dbdcv20201029.PolicyRule{
						{
							Action:     ptrStr("ACCEPT"),
							CidrIp:     ptrStr("0.0.0.0/0"),
							PortRange:  ptrStr("80"),
							IpProtocol: ptrStr("TCP"),
							Id:         ptrStr("sg-abc123"),
							Desc:       ptrStr("allow http"),
						},
					},
					Outbound: []*dbdcv20201029.PolicyRule{
						{
							Action:     ptrStr("ACCEPT"),
							CidrIp:     ptrStr("0.0.0.0/0"),
							PortRange:  ptrStr("ALL"),
							IpProtocol: ptrStr("ALL"),
							Id:         ptrStr("sg-abc123"),
							Desc:       ptrStr("allow all outbound"),
						},
					},
				},
				{
					SecurityGroupId:   ptrStr("sg-def456"),
					SecurityGroupName: ptrStr("test-sg-2"),
					ProjectId:         ptrInt64(1),
					CreateTime:        ptrStr("2024-06-01 00:00:00"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaDbdcDS()
	res := dbdc.DataSourceTencentCloudDbdcDbCustomNodeSecurityGroups()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"node_id": "dbcn-abc12345",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	groups := d.Get("groups").([]interface{})
	assert.Len(t, groups, 2)

	group0 := groups[0].(map[string]interface{})
	assert.Equal(t, "sg-abc123", group0["security_group_id"].(string))
	assert.Equal(t, "test-sg-1", group0["security_group_name"].(string))
	assert.Equal(t, "test security group 1", group0["security_group_remark"].(string))
	assert.Equal(t, int64(0), group0["project_id"].(int64))
	assert.Equal(t, "2024-01-01 00:00:00", group0["create_time"].(string))

	inbound := group0["inbound"].([]interface{})
	assert.Len(t, inbound, 1)
	inbound0 := inbound[0].(map[string]interface{})
	assert.Equal(t, "ACCEPT", inbound0["action"].(string))
	assert.Equal(t, "0.0.0.0/0", inbound0["cidr_ip"].(string))
	assert.Equal(t, "80", inbound0["port_range"].(string))
	assert.Equal(t, "TCP", inbound0["ip_protocol"].(string))
	assert.Equal(t, "allow http", inbound0["desc"].(string))

	outbound := group0["outbound"].([]interface{})
	assert.Len(t, outbound, 1)
	outbound0 := outbound[0].(map[string]interface{})
	assert.Equal(t, "ACCEPT", outbound0["action"].(string))
	assert.Equal(t, "ALL", outbound0["port_range"].(string))
	assert.Equal(t, "ALL", outbound0["ip_protocol"].(string))

	group1 := groups[1].(map[string]interface{})
	assert.Equal(t, "sg-def456", group1["security_group_id"].(string))
	assert.Equal(t, "test-sg-2", group1["security_group_name"].(string))
	assert.Equal(t, int64(1), group1["project_id"].(int64))
}

func TestDbdcDbCustomNodeSecurityGroupsDS_ReadEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdcDS().client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomNodeSecurityGroups", func(request *dbdcv20201029.DescribeDBCustomNodeSecurityGroupsRequest) (*dbdcv20201029.DescribeDBCustomNodeSecurityGroupsResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomNodeSecurityGroupsResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomNodeSecurityGroupsResponseParams{
			Groups: []*dbdcv20201029.SecurityGroup{},
		}
		return resp, nil
	})

	meta := newMockMetaDbdcDS()
	res := dbdc.DataSourceTencentCloudDbdcDbCustomNodeSecurityGroups()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"node_id": "dbcn-empty",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	groups := d.Get("groups").([]interface{})
	assert.Len(t, groups, 0)
}

func TestDbdcDbCustomNodeSecurityGroupsDS_ReadNilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdcDS().client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomNodeSecurityGroups", func(request *dbdcv20201029.DescribeDBCustomNodeSecurityGroupsRequest) (*dbdcv20201029.DescribeDBCustomNodeSecurityGroupsResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomNodeSecurityGroupsResponse()
		resp.Response = nil
		return resp, nil
	})

	meta := newMockMetaDbdcDS()
	res := dbdc.DataSourceTencentCloudDbdcDbCustomNodeSecurityGroups()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"node_id": "dbcn-nil-response",
	})

	err := res.Read(d, meta)
	assert.Error(t, err)
}

func TestDbdcDbCustomNodeSecurityGroupsDS_Schema(t *testing.T) {
	res := dbdc.DataSourceTencentCloudDbdcDbCustomNodeSecurityGroups()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "node_id")
	assert.Contains(t, res.Schema, "groups")
	assert.Contains(t, res.Schema, "result_output_file")

	nodeIdSchema := res.Schema["node_id"]
	assert.Equal(t, schema.TypeString, nodeIdSchema.Type)
	assert.True(t, nodeIdSchema.Required)

	groupsSchema := res.Schema["groups"]
	assert.Equal(t, schema.TypeList, groupsSchema.Type)
	assert.True(t, groupsSchema.Computed)

	elemRes := groupsSchema.Elem.(*schema.Resource)
	assert.Contains(t, elemRes.Schema, "security_group_id")
	assert.Contains(t, elemRes.Schema, "security_group_name")
	assert.Contains(t, elemRes.Schema, "security_group_remark")
	assert.Contains(t, elemRes.Schema, "project_id")
	assert.Contains(t, elemRes.Schema, "create_time")
	assert.Contains(t, elemRes.Schema, "inbound")
	assert.Contains(t, elemRes.Schema, "outbound")

	inboundSchema := elemRes.Schema["inbound"]
	assert.Equal(t, schema.TypeList, inboundSchema.Type)
	assert.True(t, inboundSchema.Computed)

	inboundElem := inboundSchema.Elem.(*schema.Resource)
	assert.Contains(t, inboundElem.Schema, "action")
	assert.Contains(t, inboundElem.Schema, "cidr_ip")
	assert.Contains(t, inboundElem.Schema, "port_range")
	assert.Contains(t, inboundElem.Schema, "ip_protocol")
	assert.Contains(t, inboundElem.Schema, "service_module")
	assert.Contains(t, inboundElem.Schema, "address_module")
	assert.Contains(t, inboundElem.Schema, "id")
	assert.Contains(t, inboundElem.Schema, "desc")
}
