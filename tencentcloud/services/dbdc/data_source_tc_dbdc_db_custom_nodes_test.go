package dbdc_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	dbdcv20201029 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dbdc"
)

// go test ./tencentcloud/services/dbdc/ -run "TestDbdcDbCustomNodesDS" -v -count=1 -gcflags="all=-l"

func TestDbdcDbCustomNodesDS_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdcDS().client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomNodes", func(request *dbdcv20201029.DescribeDBCustomNodesRequest) (*dbdcv20201029.DescribeDBCustomNodesResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomNodesResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomNodesResponseParams{
			TotalCount: ptrInt64(2),
			NodeSet: []*dbdcv20201029.DBCustomNode{
				{
					NodeId:      ptrStr("dbcn-abc12345"),
					NodeName:    ptrStr("test-node-1"),
					SSHEndpoint: ptrStr("10.0.0.1:22"),
					LanIP:       ptrStr("10.0.0.1"),
					ClusterId:   ptrStr("dbcc-nmtmsew8"),
					Zone:        ptrStr("ap-guangzhou-3"),
					NodeType:    ptrStr("DB.AT5.32XLARGE512"),
					CPU:         ptrInt64(128),
					Memory:      ptrInt64(512),
					SystemDisk: &dbdcv20201029.SystemDisk{
						DiskType: ptrStr("CLOUD_HSSD"),
						DiskSize: ptrInt64(50),
					},
					DataDisks: []*dbdcv20201029.DataDisk{
						{
							DiskType: ptrStr("LOCAL_NVME"),
							DiskSize: ptrInt64(7180),
							DiskName: ptrStr("data-disk-1"),
						},
					},
					OsName:       ptrStr("CentOS 7.6"),
					ImageId:      ptrStr("img-abc123"),
					VpcId:        ptrStr("vpc-abc123"),
					SubnetId:     ptrStr("subnet-abc123"),
					Status:       ptrStr("Running"),
					ChargeType:   ptrStr("PREPAID"),
					ExpireTime:   ptrStr("2025-01-01 00:00:00"),
					CreatedTime:  ptrStr("2024-01-01 00:00:00"),
					IsolatedTime: ptrStr(""),
					Tags: []*dbdcv20201029.Tag{
						{Key: ptrStr("env"), Value: ptrStr("production")},
					},
					AutoRenew: ptrInt64(1),
					SwitchId:  ptrStr("switch-abc123"),
					RackId:    ptrStr("rack-abc123"),
					HostIp:    ptrStr("192.168.1.1"),
				},
				{
					NodeId:     ptrStr("dbcn-def67890"),
					NodeName:   ptrStr("test-node-2"),
					ClusterId:  ptrStr("dbcc-nmtmsew8"),
					Status:     ptrStr("Creating"),
					ChargeType: ptrStr("PREPAID"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaDbdcDS()
	res := dbdc.DataSourceTencentCloudDbdcDbCustomNodes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	nodeSet := d.Get("node_set").([]interface{})
	assert.Len(t, nodeSet, 2)

	node0 := nodeSet[0].(map[string]interface{})
	assert.Equal(t, "dbcn-abc12345", node0["node_id"].(string))
	assert.Equal(t, "test-node-1", node0["node_name"].(string))
	assert.Equal(t, "10.0.0.1:22", node0["ssh_endpoint"].(string))
	assert.Equal(t, "10.0.0.1", node0["lan_ip"].(string))
	assert.Equal(t, "dbcc-nmtmsew8", node0["cluster_id"].(string))
	assert.Equal(t, "ap-guangzhou-3", node0["zone"].(string))
	assert.Equal(t, "DB.AT5.32XLARGE512", node0["node_type"].(string))
	assert.Equal(t, 128, node0["cpu"].(int))
	assert.Equal(t, 512, node0["memory"].(int))

	systemDisk0 := node0["system_disk"].([]interface{})
	assert.Len(t, systemDisk0, 1)
	systemDiskMap0 := systemDisk0[0].(map[string]interface{})
	assert.Equal(t, "CLOUD_HSSD", systemDiskMap0["disk_type"].(string))
	assert.Equal(t, 50, systemDiskMap0["disk_size"].(int))

	dataDisks0 := node0["data_disks"].([]interface{})
	assert.Len(t, dataDisks0, 1)
	dataDiskMap0 := dataDisks0[0].(map[string]interface{})
	assert.Equal(t, "LOCAL_NVME", dataDiskMap0["disk_type"].(string))
	assert.Equal(t, 7180, dataDiskMap0["disk_size"].(int))
	assert.Equal(t, "data-disk-1", dataDiskMap0["disk_name"].(string))

	assert.Equal(t, "CentOS 7.6", node0["os_name"].(string))
	assert.Equal(t, "img-abc123", node0["image_id"].(string))
	assert.Equal(t, "vpc-abc123", node0["vpc_id"].(string))
	assert.Equal(t, "subnet-abc123", node0["subnet_id"].(string))
	assert.Equal(t, "Running", node0["status"].(string))
	assert.Equal(t, "PREPAID", node0["charge_type"].(string))
	assert.Equal(t, "2025-01-01 00:00:00", node0["expire_time"].(string))
	assert.Equal(t, "2024-01-01 00:00:00", node0["created_time"].(string))
	assert.Equal(t, 1, node0["auto_renew"].(int))
	assert.Equal(t, "switch-abc123", node0["switch_id"].(string))
	assert.Equal(t, "rack-abc123", node0["rack_id"].(string))
	assert.Equal(t, "192.168.1.1", node0["host_ip"].(string))

	tags0 := node0["tags"].([]interface{})
	assert.Len(t, tags0, 1)
	tagMap0 := tags0[0].(map[string]interface{})
	assert.Equal(t, "env", tagMap0["key"].(string))
	assert.Equal(t, "production", tagMap0["value"].(string))

	node1 := nodeSet[1].(map[string]interface{})
	assert.Equal(t, "dbcn-def67890", node1["node_id"].(string))
	assert.Equal(t, "test-node-2", node1["node_name"].(string))
	assert.Equal(t, "Creating", node1["status"].(string))
}

func TestDbdcDbCustomNodesDS_ReadWithNilFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdcDS().client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomNodes", func(request *dbdcv20201029.DescribeDBCustomNodesRequest) (*dbdcv20201029.DescribeDBCustomNodesResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomNodesResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomNodesResponseParams{
			TotalCount: ptrInt64(1),
			NodeSet: []*dbdcv20201029.DBCustomNode{
				{
					NodeId:   ptrStr("dbcn-nil-fields"),
					NodeName: ptrStr("node-no-disks-tags"),
					Status:   ptrStr("Running"),
					// SystemDisk, DataDisks, Tags are nil
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaDbdcDS()
	res := dbdc.DataSourceTencentCloudDbdcDbCustomNodes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	nodeSet := d.Get("node_set").([]interface{})
	assert.Len(t, nodeSet, 1)

	node0 := nodeSet[0].(map[string]interface{})
	assert.Equal(t, "dbcn-nil-fields", node0["node_id"].(string))
	// SystemDisk is nil, should not have system_disk in output
	systemDiskField := node0["system_disk"]
	if systemDiskField != nil {
		systemDiskList := systemDiskField.([]interface{})
		assert.Len(t, systemDiskList, 0)
	}
	// DataDisks is nil, should not have data_disks in output
	dataDisksField := node0["data_disks"]
	if dataDisksField != nil {
		dataDisksList := dataDisksField.([]interface{})
		assert.Len(t, dataDisksList, 0)
	}
	// Tags is nil, should not have tags in output
	tagsField := node0["tags"]
	if tagsField != nil {
		tagsList := tagsField.([]interface{})
		assert.Len(t, tagsList, 0)
	}
}

func TestDbdcDbCustomNodesDS_ReadWithNodeIds(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdcDS().client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomNodes", func(request *dbdcv20201029.DescribeDBCustomNodesRequest) (*dbdcv20201029.DescribeDBCustomNodesResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomNodesResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomNodesResponseParams{
			TotalCount: ptrInt64(1),
			NodeSet: []*dbdcv20201029.DBCustomNode{
				{
					NodeId:   ptrStr("dbcn-abc12345"),
					NodeName: ptrStr("test-node-1"),
					Status:   ptrStr("Running"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaDbdcDS()
	res := dbdc.DataSourceTencentCloudDbdcDbCustomNodes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"node_ids": []interface{}{"dbcn-abc12345"},
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	nodeSet := d.Get("node_set").([]interface{})
	assert.Len(t, nodeSet, 1)
}

func TestDbdcDbCustomNodesDS_ReadWithFilters(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdcDS().client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomNodes", func(request *dbdcv20201029.DescribeDBCustomNodesRequest) (*dbdcv20201029.DescribeDBCustomNodesResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomNodesResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomNodesResponseParams{
			TotalCount: ptrInt64(1),
			NodeSet: []*dbdcv20201029.DBCustomNode{
				{
					NodeId:    ptrStr("dbcn-filter-test"),
					ClusterId: ptrStr("dbcc-nmtmsew8"),
					Status:    ptrStr("Running"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaDbdcDS()
	res := dbdc.DataSourceTencentCloudDbdcDbCustomNodes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"filters": []interface{}{
			map[string]interface{}{
				"name":   "cluster-id",
				"values": []interface{}{"dbcc-nmtmsew8"},
			},
		},
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	nodeSet := d.Get("node_set").([]interface{})
	assert.Len(t, nodeSet, 1)
}

func TestDbdcDbCustomNodesDS_ReadWithEmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdcDS().client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomNodes", func(request *dbdcv20201029.DescribeDBCustomNodesRequest) (*dbdcv20201029.DescribeDBCustomNodesResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomNodesResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomNodesResponseParams{
			TotalCount: ptrInt64(0),
			NodeSet:    []*dbdcv20201029.DBCustomNode{},
		}
		return resp, nil
	})

	meta := newMockMetaDbdcDS()
	res := dbdc.DataSourceTencentCloudDbdcDbCustomNodes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	// When response is empty, the datasource should return an error (NonRetryableError)
	assert.Error(t, err)
}

func TestDbdcDbCustomNodesDS_Schema(t *testing.T) {
	res := dbdc.DataSourceTencentCloudDbdcDbCustomNodes()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "node_ids")
	assert.Contains(t, res.Schema, "filters")
	assert.Contains(t, res.Schema, "tags")
	assert.Contains(t, res.Schema, "node_set")
	assert.Contains(t, res.Schema, "result_output_file")

	nodeIdsSchema := res.Schema["node_ids"]
	assert.Equal(t, schema.TypeList, nodeIdsSchema.Type)
	assert.True(t, nodeIdsSchema.Optional)

	nodeSetSchema := res.Schema["node_set"]
	assert.Equal(t, schema.TypeList, nodeSetSchema.Type)
	assert.True(t, nodeSetSchema.Computed)

	elemRes := nodeSetSchema.Elem.(*schema.Resource)
	assert.Contains(t, elemRes.Schema, "node_id")
	assert.Contains(t, elemRes.Schema, "node_name")
	assert.Contains(t, elemRes.Schema, "ssh_endpoint")
	assert.Contains(t, elemRes.Schema, "lan_ip")
	assert.Contains(t, elemRes.Schema, "cluster_id")
	assert.Contains(t, elemRes.Schema, "zone")
	assert.Contains(t, elemRes.Schema, "node_type")
	assert.Contains(t, elemRes.Schema, "cpu")
	assert.Contains(t, elemRes.Schema, "memory")
	assert.Contains(t, elemRes.Schema, "system_disk")
	assert.Contains(t, elemRes.Schema, "data_disks")
	assert.Contains(t, elemRes.Schema, "os_name")
	assert.Contains(t, elemRes.Schema, "image_id")
	assert.Contains(t, elemRes.Schema, "vpc_id")
	assert.Contains(t, elemRes.Schema, "subnet_id")
	assert.Contains(t, elemRes.Schema, "status")
	assert.Contains(t, elemRes.Schema, "charge_type")
	assert.Contains(t, elemRes.Schema, "expire_time")
	assert.Contains(t, elemRes.Schema, "created_time")
	assert.Contains(t, elemRes.Schema, "isolated_time")
	assert.Contains(t, elemRes.Schema, "tags")
	assert.Contains(t, elemRes.Schema, "auto_renew")
	assert.Contains(t, elemRes.Schema, "switch_id")
	assert.Contains(t, elemRes.Schema, "rack_id")
	assert.Contains(t, elemRes.Schema, "host_ip")
}
