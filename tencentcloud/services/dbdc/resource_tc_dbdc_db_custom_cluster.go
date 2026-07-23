package dbdc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbdcv20201029 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDbdcDbCustomCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDbdcDbCustomClusterCreate,
		Read:   resourceTencentCloudDbdcDbCustomClusterRead,
		Update: resourceTencentCloudDbdcDbCustomClusterUpdate,
		Delete: resourceTencentCloudDbdcDbCustomClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cluster name. Up to 128 characters, only Chinese, English and underscore are allowed.",
			},

			"container_network": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Container network. All pods in this cluster are connected to this network.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "VPC ID of the container network.",
						},
						"subnet_ids": {
							Type:        schema.TypeList,
							Required:    true,
							ForceNew:    true,
							Description: "Subnet ID list of the container network.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"api_server_network": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Network information of the cluster API Server. Must be a network owned by this account, and can be the same as the container network.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "VPC ID of the API Server network.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Subnet ID of the API Server network.",
						},
					},
				},
			},

			"cluster_description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Cluster description.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Cluster tags.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			// computed
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Region that the cluster belongs to.",
			},

			"cluster_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DB Custom cluster status. Valid values: `Creating`, `Running`, `Destroying`.",
			},

			"cluster_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster version.",
			},

			"cluster_node_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of nodes in the cluster.",
			},

			"cluster_level": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster level.",
			},

			"created_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time.",
			},
		},
	}
}

func resourceTencentCloudDbdcDbCustomClusterCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dbdc_db_custom_cluster.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service   = DbdcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request   = dbdcv20201029.NewCreateDBCustomClusterRequest()
		response  = dbdcv20201029.NewCreateDBCustomClusterResponse()
		clusterId string
	)

	if v, ok := d.GetOk("cluster_name"); ok {
		request.ClusterName = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "container_network"); ok {
		containerNetwork := dbdcv20201029.ContainerNetwork{}
		if v, ok := dMap["vpc_id"]; ok {
			containerNetwork.VpcId = helper.String(v.(string))
		}

		if v, ok := dMap["subnet_ids"]; ok {
			subnetIdsList := v.([]interface{})
			for i := range subnetIdsList {
				subnetId := subnetIdsList[i].(string)
				containerNetwork.SubnetIds = append(containerNetwork.SubnetIds, &subnetId)
			}
		}

		request.ContainerNetwork = &containerNetwork
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "api_server_network"); ok {
		apiServerNetwork := dbdcv20201029.ApiServerNetwork{}
		if v, ok := dMap["vpc_id"]; ok {
			apiServerNetwork.VpcId = helper.String(v.(string))
		}

		if v, ok := dMap["subnet_id"]; ok {
			apiServerNetwork.SubnetId = helper.String(v.(string))
		}

		request.ApiServerNetwork = &apiServerNetwork
	}

	if v, ok := d.GetOk("cluster_description"); ok {
		request.ClusterDescription = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tags"); ok {
		for tagKey, tagValue := range v.(map[string]interface{}) {
			tag := dbdcv20201029.Tag{
				Key:   helper.String(tagKey),
				Value: helper.String(tagValue.(string)),
			}

			request.Tags = append(request.Tags, &tag)
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDbdcV20201029Client().CreateDBCustomClusterWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create dbdc db custom cluster failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create dbdc db custom cluster failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.ClusterId == nil {
		return fmt.Errorf("ClusterId is nil.")
	}

	clusterId = *response.Response.ClusterId
	d.SetId(clusterId)

	// Create is async, wait for the task to succeed.
	if response.Response.TaskId != nil {
		if err := waitDBCustomTaskSucceeded(ctx, &service, *response.Response.TaskId, d.Timeout(schema.TimeoutCreate)); err != nil {
			return err
		}
	}

	return resourceTencentCloudDbdcDbCustomClusterRead(d, meta)
}

func resourceTencentCloudDbdcDbCustomClusterRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dbdc_db_custom_cluster.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service   = DbdcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		clusterId = d.Id()
	)

	respData, err := service.DescribeDBCustomClusterById(ctx, clusterId)
	if err != nil {
		return err
	}

	if respData == nil || respData.ClusterId == nil {
		log.Printf("[WARN]%s resource `tencentcloud_dbdc_db_custom_cluster` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.ClusterName != nil {
		_ = d.Set("cluster_name", respData.ClusterName)
	}

	if respData.ContainerNetwork != nil {
		containerNetworkMap := map[string]interface{}{}
		if respData.ContainerNetwork.VpcId != nil {
			containerNetworkMap["vpc_id"] = respData.ContainerNetwork.VpcId
		}

		if respData.ContainerNetwork.SubnetIds != nil {
			containerNetworkMap["subnet_ids"] = helper.StringsInterfaces(respData.ContainerNetwork.SubnetIds)
		}

		_ = d.Set("container_network", []interface{}{containerNetworkMap})
	}

	if respData.ApiServerNetwork != nil {
		apiServerNetworkMap := map[string]interface{}{}
		if respData.ApiServerNetwork.VpcId != nil {
			apiServerNetworkMap["vpc_id"] = respData.ApiServerNetwork.VpcId
		}

		if respData.ApiServerNetwork.SubnetId != nil {
			apiServerNetworkMap["subnet_id"] = respData.ApiServerNetwork.SubnetId
		}

		_ = d.Set("api_server_network", []interface{}{apiServerNetworkMap})
	}

	if respData.ClusterDescription != nil {
		_ = d.Set("cluster_description", respData.ClusterDescription)
	}

	if respData.Tags != nil {
		tags := make(map[string]interface{}, len(respData.Tags))
		for _, tag := range respData.Tags {
			if tag == nil || tag.Key == nil {
				continue
			}

			if tag.Value != nil {
				tags[*tag.Key] = *tag.Value
			} else {
				tags[*tag.Key] = ""
			}
		}

		_ = d.Set("tags", tags)
	}

	if respData.Region != nil {
		_ = d.Set("region", respData.Region)
	}

	if respData.ClusterStatus != nil {
		_ = d.Set("cluster_status", respData.ClusterStatus)
	}

	if respData.ClusterVersion != nil {
		_ = d.Set("cluster_version", respData.ClusterVersion)
	}

	if respData.ClusterNodeNum != nil {
		_ = d.Set("cluster_node_num", respData.ClusterNodeNum)
	}

	if respData.ClusterLevel != nil {
		_ = d.Set("cluster_level", respData.ClusterLevel)
	}

	if respData.CreatedTime != nil {
		_ = d.Set("created_time", respData.CreatedTime)
	}

	return nil
}

func resourceTencentCloudDbdcDbCustomClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dbdc_db_custom_cluster.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		clusterId = d.Id()
	)

	// Only tags are mutable via the API, all other arguments are ForceNew.
	if d.HasChange("tags") {
		oldRaw, newRaw := d.GetChange("tags")
		oldTags := oldRaw.(map[string]interface{})
		newTags := newRaw.(map[string]interface{})

		request := dbdcv20201029.NewModifyDBCustomClusterTagsRequest()
		request.ClusterId = helper.String(clusterId)

		for tagKey, tagValue := range newTags {
			tag := dbdcv20201029.Tag{
				Key:   helper.String(tagKey),
				Value: helper.String(tagValue.(string)),
			}

			request.AddTags = append(request.AddTags, &tag)
		}

		for tagKey := range oldTags {
			if _, ok := newTags[tagKey]; !ok {
				request.DeleteTagKeys = append(request.DeleteTagKeys, helper.String(tagKey))
			}
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDbdcV20201029Client().ModifyDBCustomClusterTagsWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify dbdc db custom cluster tags failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update dbdc db custom cluster tags failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudDbdcDbCustomClusterRead(d, meta)
}

func resourceTencentCloudDbdcDbCustomClusterDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dbdc_db_custom_cluster.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service   = DbdcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request   = dbdcv20201029.NewDestroyDBCustomClusterRequest()
		response  = dbdcv20201029.NewDestroyDBCustomClusterResponse()
		clusterId = d.Id()
	)

	request.ClusterId = helper.String(clusterId)

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDbdcV20201029Client().DestroyDBCustomClusterWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete dbdc db custom cluster failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete dbdc db custom cluster failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// Delete is async, wait for the task to succeed.
	if response.Response.TaskId != nil {
		if err := waitDBCustomTaskSucceeded(ctx, &service, *response.Response.TaskId, d.Timeout(schema.TimeoutDelete)); err != nil {
			return err
		}
	}

	return nil
}

// waitDBCustomTaskSucceeded polls DescribeDBCustomTaskStatus until the task
// status becomes `Succeeded`. It returns an error if the task fails. The poll
// loop honors the given timeout.
func waitDBCustomTaskSucceeded(ctx context.Context, service *DbdcService, taskId uint64, timeout time.Duration) error {
	return resource.Retry(timeout, func() *resource.RetryError {
		status, e := service.DescribeDBCustomTaskStatusById(ctx, taskId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		switch status {
		case "Succeeded":
			return nil
		case "Failed":
			return resource.NonRetryableError(fmt.Errorf("dbdc db custom task [%d] execution failed, status is `Failed`.", taskId))
		default:
			return resource.RetryableError(fmt.Errorf("dbdc db custom task [%d] is still running, status is `%s`.", taskId, status))
		}
	})
}
