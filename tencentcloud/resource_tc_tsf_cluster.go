/*
Provides a resource to create a tsf cluster

Example Usage

```hcl
resource "tencentcloud_tsf_cluster" "cluster" {
  cluster_name = ""
  cluster_type = ""
  vpc_id = ""
  cluster_c_i_d_r = ""
  cluster_desc = ""
  tsf_region_id = ""
  tsf_zone_id = ""
  subnet_id = ""
  cluster_version = ""
  max_node_pod_num =
  max_cluster_service_num =
  program_id = ""
  kubernete_api_server = ""
  kubernete_native_type = ""
  kubernete_native_secret = ""
  program_id_list =
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

tsf cluster can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_cluster.cluster cluster_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudTsfCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfClusterCreate,
		Read:   resourceTencentCloudTsfClusterRead,
		Update: resourceTencentCloudTsfClusterUpdate,
		Delete: resourceTencentCloudTsfClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster name.",
			},

			"cluster_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster type.",
			},

			"vpc_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Vpc id.",
			},

			"cluster_c_i_d_r": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "CIDR assigned to cluster containers and service IP.",
			},

			"cluster_desc": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Cluster notes.",
			},

			"tsf_region_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The TSF region to which the cluster belongs.",
			},

			"tsf_zone_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The TSF availability zone to which the cluster belongs.",
			},

			"subnet_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Subnet id.",
			},

			"cluster_version": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Cluster version.",
			},

			"max_node_pod_num": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The maximum number of Pods on each Node in the cluster. The value ranges from 4 to 256. When the value is not a power of 2, the nearest power of 2 will be taken up.",
			},

			"max_cluster_service_num": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The maximum number of services in the cluster. The value ranges from 32 to 32768. If it is not a power of 2, the nearest power of 2 will be taken up.",
			},

			"program_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The dataset ID to be bound.",
			},

			"kubernete_api_server": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Api address.",
			},

			"kubernete_native_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "`K` : kubeconfig, `S` : service account.",
			},

			"kubernete_native_secret": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Native secret.",
			},

			"program_id_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Program Id List.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudTsfClusterCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_cluster.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = tsf.NewCreateClusterRequest()
		response  = tsf.NewCreateClusterResponse()
		clusterId string
	)
	if v, ok := d.GetOk("cluster_name"); ok {
		request.ClusterName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cluster_type"); ok {
		request.ClusterType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cluster_c_i_d_r"); ok {
		request.ClusterCIDR = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cluster_desc"); ok {
		request.ClusterDesc = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tsf_region_id"); ok {
		request.TsfRegionId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tsf_zone_id"); ok {
		request.TsfZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cluster_version"); ok {
		request.ClusterVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("max_node_pod_num"); ok {
		request.MaxNodePodNum = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("max_cluster_service_num"); ok {
		request.MaxClusterServiceNum = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("program_id"); ok {
		request.ProgramId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("kubernete_api_server"); ok {
		request.KuberneteApiServer = helper.String(v.(string))
	}

	if v, ok := d.GetOk("kubernete_native_type"); ok {
		request.KuberneteNativeType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("kubernete_native_secret"); ok {
		request.KuberneteNativeSecret = helper.String(v.(string))
	}

	if v, ok := d.GetOk("program_id_list"); ok {
		programIdListSet := v.(*schema.Set).List()
		for i := range programIdListSet {
			programIdList := programIdListSet[i].(string)
			request.ProgramIdList = append(request.ProgramIdList, &programIdList)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().CreateCluster(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tsf cluster failed, reason:%+v", logId, err)
		return err
	}

	clusterId = *response.Response.ClusterId
	d.SetId(clusterId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"Running"}, 8*readRetryTimeout, time.Second, service.TsfClusterStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::tsf:%s:uin/:cluster/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudTsfClusterRead(d, meta)
}

func resourceTencentCloudTsfClusterRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_cluster.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	clusterId := d.Id()

	cluster, err := service.DescribeTsfClusterById(ctx, clusterId)
	if err != nil {
		return err
	}

	if cluster == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfCluster` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if cluster.ClusterName != nil {
		_ = d.Set("cluster_name", cluster.ClusterName)
	}

	if cluster.ClusterType != nil {
		_ = d.Set("cluster_type", cluster.ClusterType)
	}

	if cluster.VpcId != nil {
		_ = d.Set("vpc_id", cluster.VpcId)
	}

	if cluster.ClusterCIDR != nil {
		_ = d.Set("cluster_c_i_d_r", cluster.ClusterCIDR)
	}

	if cluster.ClusterDesc != nil {
		_ = d.Set("cluster_desc", cluster.ClusterDesc)
	}

	if cluster.TsfRegionId != nil {
		_ = d.Set("tsf_region_id", cluster.TsfRegionId)
	}

	if cluster.TsfZoneId != nil {
		_ = d.Set("tsf_zone_id", cluster.TsfZoneId)
	}

	if cluster.SubnetId != nil {
		_ = d.Set("subnet_id", cluster.SubnetId)
	}

	if cluster.ClusterVersion != nil {
		_ = d.Set("cluster_version", cluster.ClusterVersion)
	}

	if cluster.MaxNodePodNum != nil {
		_ = d.Set("max_node_pod_num", cluster.MaxNodePodNum)
	}

	if cluster.MaxClusterServiceNum != nil {
		_ = d.Set("max_cluster_service_num", cluster.MaxClusterServiceNum)
	}

	if cluster.ProgramId != nil {
		_ = d.Set("program_id", cluster.ProgramId)
	}

	if cluster.KuberneteApiServer != nil {
		_ = d.Set("kubernete_api_server", cluster.KuberneteApiServer)
	}

	if cluster.KuberneteNativeType != nil {
		_ = d.Set("kubernete_native_type", cluster.KuberneteNativeType)
	}

	if cluster.KuberneteNativeSecret != nil {
		_ = d.Set("kubernete_native_secret", cluster.KuberneteNativeSecret)
	}

	if cluster.ProgramIdList != nil {
		_ = d.Set("program_id_list", cluster.ProgramIdList)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "tsf", "cluster", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudTsfClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_cluster.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tsf.NewModifyClusterRequest()

	clusterId := d.Id()

	request.ClusterId = &clusterId

	immutableArgs := []string{"cluster_name", "cluster_type", "vpc_id", "cluster_c_i_d_r", "cluster_desc", "tsf_region_id", "tsf_zone_id", "subnet_id", "cluster_version", "max_node_pod_num", "max_cluster_service_num", "program_id", "kubernete_api_server", "kubernete_native_type", "kubernete_native_secret", "program_id_list"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("cluster_name") {
		if v, ok := d.GetOk("cluster_name"); ok {
			request.ClusterName = helper.String(v.(string))
		}
	}

	if d.HasChange("cluster_desc") {
		if v, ok := d.GetOk("cluster_desc"); ok {
			request.ClusterDesc = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().ModifyCluster(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tsf cluster failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("tsf", "cluster", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudTsfClusterRead(d, meta)
}

func resourceTencentCloudTsfClusterDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_cluster.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
	clusterId := d.Id()

	if err := service.DeleteTsfClusterById(ctx, clusterId); err != nil {
		return err
	}

	return nil
}
