package tcaplusdb_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	tcaplusdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcaplusdb/v20190823"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svctcaplusdb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tcaplusdb"
)

var testTcaplusClusterResourceName = "tencentcloud_tcaplus_cluster"
var testTcaplusClusterResourceKey = testTcaplusClusterResourceName + ".test_cluster"

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_tcaplus_cluster
	resource.AddTestSweepers("tencentcloud_tcaplus_cluster", &resource.Sweeper{
		Name: "tencentcloud_tcaplus_cluster",
		F: func(r string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			cli, _ := tcacctest.SharedClientForRegion(r)
			client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()
			service := svctcaplusdb.NewTcaplusService(client)

			clusters, err := service.DescribeClusters(ctx, "", "")
			if err != nil {
				return err
			}

			for i := range clusters {
				c := clusters[i]
				id := *c.ClusterId
				name := *c.ClusterName
				created, err := time.Parse("2006-01-02 15:04:05", *c.CreatedTime)
				if err != nil {
					created = time.Time{}
				}
				if tcacctest.IsResourcePersist(name, &created) {
					continue
				}
				_, err = service.DeleteCluster(ctx, id)
				if err != nil {
					continue
				}
			}

			return nil
		},
	})
}

func TestAccTencentCloudTcaplusClusterResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTcaplusClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTcaplusCluster,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTcaplusClusterExists(testTcaplusClusterResourceKey),
					resource.TestCheckResourceAttrSet(testTcaplusClusterResourceKey, "network_type"),
					resource.TestCheckResourceAttrSet(testTcaplusClusterResourceKey, "create_time"),
					resource.TestCheckResourceAttrSet(testTcaplusClusterResourceKey, "password_status"),
					resource.TestCheckResourceAttrSet(testTcaplusClusterResourceKey, "api_access_id"),
					resource.TestCheckResourceAttrSet(testTcaplusClusterResourceKey, "api_access_ip"),
					resource.TestCheckResourceAttrSet(testTcaplusClusterResourceKey, "api_access_port"),
					resource.TestCheckResourceAttr(testTcaplusClusterResourceKey, "idl_type", "PROTO"),
					resource.TestCheckResourceAttr(testTcaplusClusterResourceKey, "cluster_name", "tf_te1_guagua"),
					resource.TestCheckResourceAttr(testTcaplusClusterResourceKey, "old_password_expire_last", "3600"),
				),
			},
			{
				ResourceName:            testTcaplusClusterResourceKey,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"old_password_expire_last", "password"},
			},

			{
				Config: testAccTcaplusClusterUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTcaplusClusterExists(testTcaplusClusterResourceKey),
					resource.TestCheckResourceAttrSet(testTcaplusClusterResourceKey, "network_type"),
					resource.TestCheckResourceAttrSet(testTcaplusClusterResourceKey, "create_time"),
					resource.TestCheckResourceAttrSet(testTcaplusClusterResourceKey, "password_status"),
					resource.TestCheckResourceAttrSet(testTcaplusClusterResourceKey, "api_access_id"),
					resource.TestCheckResourceAttrSet(testTcaplusClusterResourceKey, "api_access_ip"),
					resource.TestCheckResourceAttrSet(testTcaplusClusterResourceKey, "api_access_port"),

					resource.TestCheckResourceAttr(testTcaplusClusterResourceKey, "idl_type", "PROTO"),
					resource.TestCheckResourceAttr(testTcaplusClusterResourceKey, "cluster_name", "tf_te1_guagua_2"),
					resource.TestCheckResourceAttr(testTcaplusClusterResourceKey, "password", "aQQ2345677888"),
					resource.TestCheckResourceAttr(testTcaplusClusterResourceKey, "old_password_expire_last", "300"),
				),
			},
		},
	})
}

func testAccCheckTcaplusClusterDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testTcaplusClusterResourceName {
			continue
		}
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		service := svctcaplusdb.NewTcaplusService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		_, has, err := service.DescribeCluster(ctx, rs.Primary.ID)
		if err != nil {
			_, has, err = service.DescribeCluster(ctx, rs.Primary.ID)
		}
		if err != nil {
			return err
		}
		if !has {
			return nil
		} else {
			return fmt.Errorf("delete tcaplus cluster %s fail, still on server", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTcaplusClusterExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		service := svctcaplusdb.NewTcaplusService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		_, has, err := service.DescribeCluster(ctx, rs.Primary.ID)
		if err != nil {
			_, has, err = service.DescribeCluster(ctx, rs.Primary.ID)
		}
		if err != nil {
			return err
		}
		if has {
			return nil
		} else {
			return fmt.Errorf("tcaplus cluster %s not found on server", rs.Primary.ID)

		}
	}
}

const testAccTcaplusCluster string = tcacctest.DefaultVpcSubnets + `
resource "tencentcloud_tcaplus_cluster" "test_cluster" {
  idl_type                 = "PROTO"
  cluster_name             = "tf_te1_guagua"
  vpc_id                   = local.vpc_id
  subnet_id                = local.subnet_id
  password                 = "1qaA2k1wgvfa3ZZZ"
  old_password_expire_last = 3600
}
`
const testAccTcaplusClusterUpdate string = tcacctest.DefaultVpcSubnets + `
resource "tencentcloud_tcaplus_cluster" "test_cluster" {
  idl_type                 = "PROTO"
  cluster_name             = "tf_te1_guagua_2"
  vpc_id                   = local.vpc_id
  subnet_id                = local.subnet_id
  password                 = "aQQ2345677888"
  old_password_expire_last = 300
}
`

// ---- gomonkey-based unit tests for tcaplus_cluster new parameters ----

type mockMetaTcaplusCluster struct {
	client *connectivity.TencentCloudClient
}

var _ tccommon.ProviderMeta = &mockMetaTcaplusCluster{}

func (m *mockMetaTcaplusCluster) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

func newMockMetaTcaplusCluster() *mockMetaTcaplusCluster {
	return &mockMetaTcaplusCluster{client: &connectivity.TencentCloudClient{}}
}

func ptrStringTcaplusCluster(s string) *string {
	return &s
}

func ptrInt64TcaplusCluster(i int64) *int64 {
	return &i
}

func resourceTagsTcaplusClusterRaw() []interface{} {
	return []interface{}{
		map[string]interface{}{
			"tag_key":   "env",
			"tag_value": "prod",
		},
		map[string]interface{}{
			"tag_key":   "owner",
			"tag_value": "terraform",
		},
	}
}

func serverListTcaplusClusterRaw() []interface{} {
	return []interface{}{
		map[string]interface{}{
			"machine_type": "S5.LARGE8",
			"machine_num":  2,
		},
	}
}

func proxyListTcaplusClusterRaw() []interface{} {
	return []interface{}{
		map[string]interface{}{
			"machine_type": "S5.LARGE8",
			"machine_num":  1,
		},
	}
}

// TestTcaplusClusterDedicatedCreate verifies cluster_type, resource_tags, server_list, proxy_list are passed to CreateCluster.
func TestTcaplusClusterDedicatedCreate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tcaplusClient := &tcaplusdb.Client{}
	patches.ApplyMethodReturn(newMockMetaTcaplusCluster().client, "UseTcaplusClient", tcaplusClient)

	var capturedRequest *tcaplusdb.CreateClusterRequest
	patches.ApplyMethodFunc(tcaplusClient, "CreateCluster", func(request *tcaplusdb.CreateClusterRequest) (*tcaplusdb.CreateClusterResponse, error) {
		capturedRequest = request
		resp := tcaplusdb.NewCreateClusterResponse()
		resp.Response = &tcaplusdb.CreateClusterResponseParams{
			ClusterId: ptrStringTcaplusCluster("1910000000"),
			RequestId: ptrStringTcaplusCluster("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(tcaplusClient, "DescribeClusters", func(request *tcaplusdb.DescribeClustersRequest) (*tcaplusdb.DescribeClustersResponse, error) {
		resp := tcaplusdb.NewDescribeClustersResponse()
		resp.Response = &tcaplusdb.DescribeClustersResponseParams{
			TotalCount: ptrInt64TcaplusCluster(1),
			Clusters: []*tcaplusdb.ClusterInfo{
				{
					ClusterId:      ptrStringTcaplusCluster("1910000000"),
					ClusterName:    ptrStringTcaplusCluster("tf_dedicated_cluster"),
					IdlType:        ptrStringTcaplusCluster("PROTO"),
					VpcId:          ptrStringTcaplusCluster("vpc-xxx"),
					SubnetId:       ptrStringTcaplusCluster("subnet-xxx"),
					NetworkType:    ptrStringTcaplusCluster("ccnz"),
					CreatedTime:    ptrStringTcaplusCluster("2025-01-01 00:00:00"),
					PasswordStatus: ptrStringTcaplusCluster("modifiable"),
					ApiAccessId:    ptrStringTcaplusCluster("access-id"),
					ApiAccessIp:    ptrStringTcaplusCluster("1.2.3.4"),
					ApiAccessPort:  ptrInt64TcaplusCluster(9999),
					ClusterType:    ptrInt64TcaplusCluster(2),
					ServerList: []*tcaplusdb.ServerDetailInfo{
						{
							ServerUid:   ptrStringTcaplusCluster("server-uid-1"),
							MachineType: ptrStringTcaplusCluster("S5.LARGE8"),
							MemoryRate:  ptrInt64TcaplusCluster(50),
							DiskRate:    ptrInt64TcaplusCluster(60),
							ReadNum:     ptrInt64TcaplusCluster(100),
							WriteNum:    ptrInt64TcaplusCluster(200),
							Version:     ptrStringTcaplusCluster("3.5.0"),
						},
					},
					ProxyList: []*tcaplusdb.ProxyDetailInfo{
						{
							ProxyUid:            ptrStringTcaplusCluster("proxy-uid-1"),
							MachineType:         ptrStringTcaplusCluster("S5.LARGE8"),
							ProcessSpeed:        ptrInt64TcaplusCluster(1000),
							AverageProcessDelay: ptrInt64TcaplusCluster(5),
							SlowProcessSpeed:    ptrInt64TcaplusCluster(10),
							Version:             ptrStringTcaplusCluster("3.5.0"),
						},
					},
				},
			},
			RequestId: ptrStringTcaplusCluster("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaTcaplusCluster()
	res := svctcaplusdb.ResourceTencentCloudTcaplusCluster()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"idl_type":      "PROTO",
		"cluster_name":  "tf_dedicated_cluster",
		"vpc_id":        "vpc-xxx",
		"subnet_id":     "subnet-xxx",
		"password":      "1qaA2k1wgvfa3ZZZ",
		"cluster_type":  2,
		"resource_tags": resourceTagsTcaplusClusterRaw(),
		"server_list":   serverListTcaplusClusterRaw(),
		"proxy_list":    proxyListTcaplusClusterRaw(),
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "1910000000", d.Id())

	assert.NotNil(t, capturedRequest.ClusterType)
	assert.Equal(t, int64(2), *capturedRequest.ClusterType)

	assert.NotNil(t, capturedRequest.ResourceTags)
	assert.Equal(t, 2, len(capturedRequest.ResourceTags))
	assert.Equal(t, "env", *capturedRequest.ResourceTags[0].TagKey)
	assert.Equal(t, "prod", *capturedRequest.ResourceTags[0].TagValue)
	assert.Equal(t, "owner", *capturedRequest.ResourceTags[1].TagKey)
	assert.Equal(t, "terraform", *capturedRequest.ResourceTags[1].TagValue)

	assert.NotNil(t, capturedRequest.ServerList)
	assert.Equal(t, 1, len(capturedRequest.ServerList))
	assert.Equal(t, "S5.LARGE8", *capturedRequest.ServerList[0].MachineType)
	assert.Equal(t, int64(2), *capturedRequest.ServerList[0].MachineNum)

	assert.NotNil(t, capturedRequest.ProxyList)
	assert.Equal(t, 1, len(capturedRequest.ProxyList))
	assert.Equal(t, "S5.LARGE8", *capturedRequest.ProxyList[0].MachineType)
	assert.Equal(t, int64(1), *capturedRequest.ProxyList[0].MachineNum)

	clusterType := d.Get("cluster_type").(int)
	assert.Equal(t, 2, clusterType)

	serverList := d.Get("server_list").([]interface{})
	assert.Equal(t, 1, len(serverList))
	serverMap := serverList[0].(map[string]interface{})
	assert.Equal(t, "server-uid-1", serverMap["server_uid"].(string))
	assert.Equal(t, "S5.LARGE8", serverMap["machine_type"].(string))
	assert.Equal(t, 50, serverMap["memory_rate"].(int))
	assert.Equal(t, 60, serverMap["disk_rate"].(int))
	assert.Equal(t, 100, serverMap["read_num"].(int))
	assert.Equal(t, 200, serverMap["write_num"].(int))
	assert.Equal(t, "3.5.0", serverMap["version"].(string))

	proxyList := d.Get("proxy_list").([]interface{})
	assert.Equal(t, 1, len(proxyList))
	proxyMap := proxyList[0].(map[string]interface{})
	assert.Equal(t, "proxy-uid-1", proxyMap["proxy_uid"].(string))
	assert.Equal(t, "S5.LARGE8", proxyMap["machine_type"].(string))
	assert.Equal(t, 1000, proxyMap["process_speed"].(int))
	assert.Equal(t, 5, proxyMap["average_process_delay"].(int))
	assert.Equal(t, 10, proxyMap["slow_process_speed"].(int))
	assert.Equal(t, "3.5.0", proxyMap["version"].(string))
}

// TestTcaplusClusterSharedCreateDefault verifies new params are not sent when not specified.
func TestTcaplusClusterSharedCreateDefault(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tcaplusClient := &tcaplusdb.Client{}
	patches.ApplyMethodReturn(newMockMetaTcaplusCluster().client, "UseTcaplusClient", tcaplusClient)

	var capturedRequest *tcaplusdb.CreateClusterRequest
	patches.ApplyMethodFunc(tcaplusClient, "CreateCluster", func(request *tcaplusdb.CreateClusterRequest) (*tcaplusdb.CreateClusterResponse, error) {
		capturedRequest = request
		resp := tcaplusdb.NewCreateClusterResponse()
		resp.Response = &tcaplusdb.CreateClusterResponseParams{
			ClusterId: ptrStringTcaplusCluster("1910000001"),
			RequestId: ptrStringTcaplusCluster("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(tcaplusClient, "DescribeClusters", func(request *tcaplusdb.DescribeClustersRequest) (*tcaplusdb.DescribeClustersResponse, error) {
		resp := tcaplusdb.NewDescribeClustersResponse()
		resp.Response = &tcaplusdb.DescribeClustersResponseParams{
			TotalCount: ptrInt64TcaplusCluster(1),
			Clusters: []*tcaplusdb.ClusterInfo{
				{
					ClusterId:      ptrStringTcaplusCluster("1910000001"),
					ClusterName:    ptrStringTcaplusCluster("tf_shared_cluster"),
					IdlType:        ptrStringTcaplusCluster("PROTO"),
					VpcId:          ptrStringTcaplusCluster("vpc-xxx"),
					SubnetId:       ptrStringTcaplusCluster("subnet-xxx"),
					NetworkType:    ptrStringTcaplusCluster("ccnz"),
					CreatedTime:    ptrStringTcaplusCluster("2025-01-01 00:00:00"),
					PasswordStatus: ptrStringTcaplusCluster("modifiable"),
					ApiAccessId:    ptrStringTcaplusCluster("access-id"),
					ApiAccessIp:    ptrStringTcaplusCluster("1.2.3.4"),
					ApiAccessPort:  ptrInt64TcaplusCluster(9999),
					ClusterType:    ptrInt64TcaplusCluster(1),
				},
			},
			RequestId: ptrStringTcaplusCluster("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaTcaplusCluster()
	res := svctcaplusdb.ResourceTencentCloudTcaplusCluster()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"idl_type":     "PROTO",
		"cluster_name": "tf_shared_cluster",
		"vpc_id":       "vpc-xxx",
		"subnet_id":    "subnet-xxx",
		"password":     "1qaA2k1wgvfa3ZZZ",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)

	assert.Nil(t, capturedRequest.ClusterType)
	assert.Nil(t, capturedRequest.ResourceTags)
	assert.Nil(t, capturedRequest.ServerList)
	assert.Nil(t, capturedRequest.ProxyList)

	clusterType := d.Get("cluster_type").(int)
	assert.Equal(t, 1, clusterType)
}

// TestTcaplusClusterReadServerListNil verifies Read does not panic when ServerList fields are nil.
func TestTcaplusClusterReadServerListNil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tcaplusClient := &tcaplusdb.Client{}
	patches.ApplyMethodReturn(newMockMetaTcaplusCluster().client, "UseTcaplusClient", tcaplusClient)

	patches.ApplyMethodFunc(tcaplusClient, "DescribeClusters", func(request *tcaplusdb.DescribeClustersRequest) (*tcaplusdb.DescribeClustersResponse, error) {
		resp := tcaplusdb.NewDescribeClustersResponse()
		resp.Response = &tcaplusdb.DescribeClustersResponseParams{
			TotalCount: ptrInt64TcaplusCluster(1),
			Clusters: []*tcaplusdb.ClusterInfo{
				{
					ClusterId:      ptrStringTcaplusCluster("1910000002"),
					ClusterName:    ptrStringTcaplusCluster("tf_nil_cluster"),
					IdlType:        ptrStringTcaplusCluster("PROTO"),
					VpcId:          ptrStringTcaplusCluster("vpc-xxx"),
					SubnetId:       ptrStringTcaplusCluster("subnet-xxx"),
					NetworkType:    ptrStringTcaplusCluster("ccnz"),
					CreatedTime:    ptrStringTcaplusCluster("2025-01-01 00:00:00"),
					PasswordStatus: ptrStringTcaplusCluster("modifiable"),
					ApiAccessId:    ptrStringTcaplusCluster("access-id"),
					ApiAccessIp:    ptrStringTcaplusCluster("1.2.3.4"),
					ApiAccessPort:  ptrInt64TcaplusCluster(9999),
					ClusterType:    ptrInt64TcaplusCluster(2),
					ServerList: []*tcaplusdb.ServerDetailInfo{
						{
							ServerUid:   ptrStringTcaplusCluster("server-uid-nil"),
							MachineType: ptrStringTcaplusCluster("S5.LARGE8"),
						},
					},
					ProxyList: []*tcaplusdb.ProxyDetailInfo{
						{
							ProxyUid:    ptrStringTcaplusCluster("proxy-uid-nil"),
							MachineType: ptrStringTcaplusCluster("S5.LARGE8"),
						},
					},
				},
			},
			RequestId: ptrStringTcaplusCluster("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaTcaplusCluster()
	res := svctcaplusdb.ResourceTencentCloudTcaplusCluster()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"idl_type":     "PROTO",
		"cluster_name": "tf_nil_cluster",
		"vpc_id":       "vpc-xxx",
		"subnet_id":    "subnet-xxx",
		"password":     "1qaA2k1wgvfa3ZZZ",
		"cluster_type": 2,
		"server_list":  serverListTcaplusClusterRaw(),
		"proxy_list":   proxyListTcaplusClusterRaw(),
	})
	d.SetId("1910000002")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	clusterType := d.Get("cluster_type").(int)
	assert.Equal(t, 2, clusterType)

	serverList := d.Get("server_list").([]interface{})
	assert.Equal(t, 1, len(serverList))
	serverMap := serverList[0].(map[string]interface{})
	assert.Equal(t, "server-uid-nil", serverMap["server_uid"].(string))
	assert.Equal(t, "S5.LARGE8", serverMap["machine_type"].(string))
	assert.Equal(t, 0, serverMap["memory_rate"].(int))
	assert.Equal(t, 0, serverMap["disk_rate"].(int))
	assert.Equal(t, 0, serverMap["read_num"].(int))
	assert.Equal(t, 0, serverMap["write_num"].(int))
	assert.Equal(t, "", serverMap["version"].(string))

	proxyList := d.Get("proxy_list").([]interface{})
	assert.Equal(t, 1, len(proxyList))
	proxyMap := proxyList[0].(map[string]interface{})
	assert.Equal(t, "proxy-uid-nil", proxyMap["proxy_uid"].(string))
	assert.Equal(t, "S5.LARGE8", proxyMap["machine_type"].(string))
	assert.Equal(t, 0, proxyMap["process_speed"].(int))
	assert.Equal(t, 0, proxyMap["average_process_delay"].(int))
	assert.Equal(t, 0, proxyMap["slow_process_speed"].(int))
	assert.Equal(t, "", proxyMap["version"].(string))
}

// TestTcaplusClusterReadProxyListNil verifies Read handles nil ProxyList fields gracefully.
func TestTcaplusClusterReadProxyListNil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tcaplusClient := &tcaplusdb.Client{}
	patches.ApplyMethodReturn(newMockMetaTcaplusCluster().client, "UseTcaplusClient", tcaplusClient)

	patches.ApplyMethodFunc(tcaplusClient, "DescribeClusters", func(request *tcaplusdb.DescribeClustersRequest) (*tcaplusdb.DescribeClustersResponse, error) {
		resp := tcaplusdb.NewDescribeClustersResponse()
		resp.Response = &tcaplusdb.DescribeClustersResponseParams{
			TotalCount: ptrInt64TcaplusCluster(1),
			Clusters: []*tcaplusdb.ClusterInfo{
				{
					ClusterId:   ptrStringTcaplusCluster("1910000003"),
					ClusterName: ptrStringTcaplusCluster("tf_proxy_nil"),
					IdlType:     ptrStringTcaplusCluster("PROTO"),
					VpcId:       ptrStringTcaplusCluster("vpc-xxx"),
					SubnetId:    ptrStringTcaplusCluster("subnet-xxx"),
					ClusterType: ptrInt64TcaplusCluster(2),
					ProxyList: []*tcaplusdb.ProxyDetailInfo{
						{
							ProxyUid: ptrStringTcaplusCluster("proxy-uid-only"),
						},
					},
				},
			},
			RequestId: ptrStringTcaplusCluster("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaTcaplusCluster()
	res := svctcaplusdb.ResourceTencentCloudTcaplusCluster()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"idl_type":     "PROTO",
		"cluster_name": "tf_proxy_nil",
		"vpc_id":       "vpc-xxx",
		"subnet_id":    "subnet-xxx",
		"password":     "1qaA2k1wgvfa3ZZZ",
		"cluster_type": 2,
	})
	d.SetId("1910000003")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	proxyList := d.Get("proxy_list").([]interface{})
	assert.Equal(t, 1, len(proxyList))
	proxyMap := proxyList[0].(map[string]interface{})
	assert.Equal(t, "proxy-uid-only", proxyMap["proxy_uid"].(string))
	assert.Equal(t, "", proxyMap["machine_type"].(string))
	assert.Equal(t, 0, proxyMap["process_speed"].(int))
	assert.Equal(t, 0, proxyMap["average_process_delay"].(int))
	assert.Equal(t, 0, proxyMap["slow_process_speed"].(int))
	assert.Equal(t, "", proxyMap["version"].(string))
}

// TestTcaplusClusterImmutableReject verifies Update rejects changes to immutable args.
func TestTcaplusClusterImmutableReject(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tcaplusClient := &tcaplusdb.Client{}
	patches.ApplyMethodReturn(newMockMetaTcaplusCluster().client, "UseTcaplusClient", tcaplusClient)

	patches.ApplyMethodFunc(tcaplusClient, "DescribeClusters", func(request *tcaplusdb.DescribeClustersRequest) (*tcaplusdb.DescribeClustersResponse, error) {
		resp := tcaplusdb.NewDescribeClustersResponse()
		resp.Response = &tcaplusdb.DescribeClustersResponseParams{
			TotalCount: ptrInt64TcaplusCluster(1),
			Clusters: []*tcaplusdb.ClusterInfo{
				{
					ClusterId:   ptrStringTcaplusCluster("1910000004"),
					ClusterName: ptrStringTcaplusCluster("tf_immutable"),
					IdlType:     ptrStringTcaplusCluster("PROTO"),
					VpcId:       ptrStringTcaplusCluster("vpc-xxx"),
					SubnetId:    ptrStringTcaplusCluster("subnet-xxx"),
					ClusterType: ptrInt64TcaplusCluster(2),
				},
			},
			RequestId: ptrStringTcaplusCluster("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaTcaplusCluster()
	res := svctcaplusdb.ResourceTencentCloudTcaplusCluster()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"idl_type":     "PROTO",
		"cluster_name": "tf_immutable",
		"vpc_id":       "vpc-xxx",
		"subnet_id":    "subnet-xxx",
		"password":     "1qaA2k1wgvfa3ZZZ",
		"cluster_type": 1,
	})
	d.SetId("1910000004")

	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "cluster_type"
	})

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cluster_type")
	assert.Contains(t, err.Error(), "cannot be changed")
}

// TestTcaplusClusterImmutableRejectServerList verifies Update rejects changes to server_list.
func TestTcaplusClusterImmutableRejectServerList(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tcaplusClient := &tcaplusdb.Client{}
	patches.ApplyMethodReturn(newMockMetaTcaplusCluster().client, "UseTcaplusClient", tcaplusClient)

	patches.ApplyMethodFunc(tcaplusClient, "DescribeClusters", func(request *tcaplusdb.DescribeClustersRequest) (*tcaplusdb.DescribeClustersResponse, error) {
		resp := tcaplusdb.NewDescribeClustersResponse()
		resp.Response = &tcaplusdb.DescribeClustersResponseParams{
			TotalCount: ptrInt64TcaplusCluster(1),
			Clusters: []*tcaplusdb.ClusterInfo{
				{
					ClusterId:   ptrStringTcaplusCluster("1910000005"),
					ClusterName: ptrStringTcaplusCluster("tf_immutable_sl"),
					IdlType:     ptrStringTcaplusCluster("PROTO"),
					VpcId:       ptrStringTcaplusCluster("vpc-xxx"),
					SubnetId:    ptrStringTcaplusCluster("subnet-xxx"),
				},
			},
			RequestId: ptrStringTcaplusCluster("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaTcaplusCluster()
	res := svctcaplusdb.ResourceTencentCloudTcaplusCluster()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"idl_type":     "PROTO",
		"cluster_name": "tf_immutable_sl",
		"vpc_id":       "vpc-xxx",
		"subnet_id":    "subnet-xxx",
		"password":     "1qaA2k1wgvfa3ZZZ",
	})
	d.SetId("1910000005")

	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "server_list"
	})

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "server_list")
	assert.Contains(t, err.Error(), "cannot be changed")
}

// TestTcaplusClusterSchema validates the new schema fields.
func TestTcaplusClusterSchema(t *testing.T) {
	res := svctcaplusdb.ResourceTencentCloudTcaplusCluster()

	clusterTypeField, ok := res.Schema["cluster_type"]
	assert.True(t, ok, "cluster_type should exist in schema")
	assert.Equal(t, schema.TypeInt, clusterTypeField.Type)
	assert.True(t, clusterTypeField.Optional)
	assert.True(t, clusterTypeField.Computed)

	resourceTagsField, ok := res.Schema["resource_tags"]
	assert.True(t, ok, "resource_tags should exist in schema")
	assert.Equal(t, schema.TypeList, resourceTagsField.Type)
	assert.True(t, resourceTagsField.Optional)

	serverListField, ok := res.Schema["server_list"]
	assert.True(t, ok, "server_list should exist in schema")
	assert.Equal(t, schema.TypeList, serverListField.Type)
	assert.True(t, serverListField.Optional)
	assert.True(t, serverListField.Computed)

	serverSchema := serverListField.Elem.(*schema.Resource).Schema
	_, ok = serverSchema["machine_type"]
	assert.True(t, ok, "machine_type should exist in server_list schema")
	_, ok = serverSchema["machine_num"]
	assert.True(t, ok, "machine_num should exist in server_list schema")
	_, ok = serverSchema["server_uid"]
	assert.True(t, ok, "server_uid should exist in server_list schema")

	proxyListField, ok := res.Schema["proxy_list"]
	assert.True(t, ok, "proxy_list should exist in schema")
	assert.Equal(t, schema.TypeList, proxyListField.Type)
	assert.True(t, proxyListField.Optional)
	assert.True(t, proxyListField.Computed)

	proxySchema := proxyListField.Elem.(*schema.Resource).Schema
	_, ok = proxySchema["machine_type"]
	assert.True(t, ok, "machine_type should exist in proxy_list schema")
	_, ok = proxySchema["machine_num"]
	assert.True(t, ok, "machine_num should exist in proxy_list schema")
	_, ok = proxySchema["proxy_uid"]
	assert.True(t, ok, "proxy_uid should exist in proxy_list schema")
}
