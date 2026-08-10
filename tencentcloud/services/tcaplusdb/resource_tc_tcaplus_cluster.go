package tcaplusdb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tcaplusdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcaplusdb/v20190823"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTcaplusCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcaplusClusterCreate,
		Read:   resourceTencentCloudTcaplusClusterRead,
		Update: resourceTencentCloudTcaplusClusterUpdate,
		Delete: resourceTencentCloudTcaplusClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"idl_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cluster data description language type, uniformly filled with `MIX`, enumeration value: `MIX`: supports both `PROTO` and `TDR` tables.",
			},

			"cluster_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cluster name, Chinese or English characters can be used, maximum length is 32 characters.",
			},

			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The private network instance ID bound to the cluster, such as: `vpc-f49l6u0z`.",
			},

			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The subnet instance ID bound to the cluster, such as: `subnet-pxir56ns`.",
			},

			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Cluster access password, must be `a-zA-Z0-9` characters, and must contain numbers, uppercase and lowercase letters.",
			},

			"old_password_expire_last": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3600,
				Description: "Expiration time of old password after password update, unit: second.",
			},

			"cluster_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Cluster type: `1` shared, `2` dedicated.",
			},

			"resource_tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "Cluster tag set. Note: this field cannot be modified after cluster creation via CreateCluster, but can be modified via ModifyClusterTags. Tags will be refreshed on Read via DescribeClusterTags.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tag value.",
						},
					},
				},
			},

			"server_list": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "Dedicated cluster occupied svr machines. Only valid when `cluster_type` is `2` (dedicated cluster). For creation, each element exposes `machine_type` and `machine_num`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"machine_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Machine type.",
						},
						"machine_num": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Machine quantity.",
						},
					},
				},
			},

			"proxy_list": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "Dedicated cluster occupied proxy machines. Only valid when `cluster_type` is `2` (dedicated cluster). For creation, each element exposes `machine_type` and `machine_num`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"machine_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Machine type.",
						},
						"machine_num": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Machine quantity.",
						},
					},
				},
			},

			// Computed values.
			"cluster_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster ID.",
			},
			"network_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Network type of the TcaplusDB cluster.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the TcaplusDB cluster.",
			},

			"password_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Password status of the TcaplusDB cluster. Valid values: `unmodifiable`, `modifiable`. `unmodifiable`. which means the password can not be changed in this moment; `modifiable`, which means the password can be changed in this moment.",
			},

			"api_access_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Access ID of the TcaplusDB cluster. For TcaplusDB SDK connect.",
			},

			"api_access_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Access IP of the TcaplusDB cluster. For TcaplusDB SDK connect.",
			},

			"api_access_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Access port of the TcaplusDB cluster. For TcaplusDB SDK connect.",
			},

			"old_password_expire_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Expiration time of the old password. If `password_status` is `unmodifiable`, it means the old password has not yet expired.",
			},
		},
	}
}

func resourceTencentCloudTcaplusClusterCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcaplus_cluster.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = tcaplusdb.NewCreateClusterRequest()
	)

	if v, ok := d.GetOk("idl_type"); ok {
		request.IdlType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cluster_name"); ok {
		request.ClusterName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("password"); ok {
		request.Password = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("cluster_type"); ok {
		request.ClusterType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("resource_tags"); ok {
		tags := v.(*schema.Set).List()
		for _, tag := range tags {
			tagMap := tag.(map[string]interface{})
			tagKey := tagMap["tag_key"].(string)
			tagValue := tagMap["tag_value"].(string)
			request.ResourceTags = append(request.ResourceTags, &tcaplusdb.TagInfoUnit{
				TagKey:   helper.String(tagKey),
				TagValue: helper.String(tagValue),
			})
		}
	}

	if v, ok := d.GetOk("server_list"); ok {
		servers := v.(*schema.Set).List()
		for _, server := range servers {
			serverMap := server.(map[string]interface{})
			machineType := serverMap["machine_type"].(string)
			machineNum := int64(serverMap["machine_num"].(int))
			request.ServerList = append(request.ServerList, &tcaplusdb.MachineInfo{
				MachineType: helper.String(machineType),
				MachineNum:  helper.Int64(machineNum),
			})
		}
	}

	if v, ok := d.GetOk("proxy_list"); ok {
		proxies := v.(*schema.Set).List()
		for _, proxy := range proxies {
			proxyMap := proxy.(map[string]interface{})
			machineType := proxyMap["machine_type"].(string)
			machineNum := int64(proxyMap["machine_num"].(int))
			request.ProxyList = append(request.ProxyList, &tcaplusdb.MachineInfo{
				MachineType: helper.String(machineType),
				MachineNum:  helper.Int64(machineNum),
			})
		}
	}

	var clusterId string
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTcaplusClient().CreateClusterWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create tcaplus cluster failed, Response is nil."))
		}

		if result.Response.ClusterId == nil || *result.Response.ClusterId == "" {
			return resource.NonRetryableError(fmt.Errorf("Create tcaplus cluster failed, ClusterId is nil or empty."))
		}

		clusterId = *result.Response.ClusterId
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create tcaplus cluster failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(clusterId)

	// wait
	waitReq := tcaplusdb.NewDescribeClustersRequest()
	waitReq.ClusterIds = []*string{&clusterId}
	reqErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTcaplusClient().DescribeClustersWithContext(ctx, waitReq)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Clusters == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe tcaplus cluster failed, Response is nil."))
		}

		if len(result.Response.Clusters) == 0 {
			return resource.NonRetryableError(fmt.Errorf("Describe tcaplus cluster failed, Clusters is empty."))
		}

		cluster := result.Response.Clusters[0]
		if cluster.ClusterStatus == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe tcaplus cluster failed, ClusterStatus is nil."))
		}

		if *cluster.ClusterStatus == 0 {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("cluster is still creating, current status: %d", *cluster.ClusterStatus))
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create tcaplus cluster failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudTcaplusClusterRead(d, meta)
}

func resourceTencentCloudTcaplusClusterRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcaplus_cluster.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service   = TcaplusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		clusterId = d.Id()
	)

	var clusterInfo tcaplusdb.ClusterInfo
	var has bool
	var reqErr error
	clusterInfo, has, reqErr = service.DescribeCluster(ctx, clusterId)
	if reqErr != nil {
		log.Printf("[CRITAL]%s read tcaplus cluster failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if !has {
		log.Printf("[WARN]%s resource `tencentcloud_tcaplus_cluster` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if clusterInfo.ClusterId != nil {
		_ = d.Set("cluster_id", clusterInfo.ClusterId)
	}

	if clusterInfo.IdlType != nil {
		_ = d.Set("idl_type", clusterInfo.IdlType)
	}

	if clusterInfo.ClusterName != nil {
		_ = d.Set("cluster_name", clusterInfo.ClusterName)
	}

	if clusterInfo.VpcId != nil {
		_ = d.Set("vpc_id", clusterInfo.VpcId)
	}

	if clusterInfo.SubnetId != nil {
		_ = d.Set("subnet_id", clusterInfo.SubnetId)
	}

	if clusterInfo.NetworkType != nil {
		_ = d.Set("network_type", clusterInfo.NetworkType)
	}

	if clusterInfo.CreatedTime != nil {
		_ = d.Set("create_time", clusterInfo.CreatedTime)
	}

	if clusterInfo.PasswordStatus != nil {
		_ = d.Set("password_status", clusterInfo.PasswordStatus)
	}

	if clusterInfo.ApiAccessId != nil {
		_ = d.Set("api_access_id", clusterInfo.ApiAccessId)
	}

	if clusterInfo.ApiAccessIp != nil {
		_ = d.Set("api_access_ip", clusterInfo.ApiAccessIp)
	}

	if clusterInfo.ApiAccessPort != nil {
		_ = d.Set("api_access_port", clusterInfo.ApiAccessPort)
	}

	if clusterInfo.OldPasswordExpireTime == nil || *clusterInfo.OldPasswordExpireTime == "" {
		_ = d.Set("old_password_expire_time", "-")
	} else {
		_ = d.Set("old_password_expire_time", clusterInfo.OldPasswordExpireTime)
	}

	if clusterInfo.ClusterType != nil {
		_ = d.Set("cluster_type", clusterInfo.ClusterType)
	}

	if clusterInfo.ServerList != nil {
		serverMachineNum := make(map[string]int)
		for _, server := range clusterInfo.ServerList {
			if server == nil || server.MachineType == nil {
				continue
			}
			serverMachineNum[*server.MachineType]++
		}
		serverList := make([]interface{}, 0, len(serverMachineNum))
		for machineType, count := range serverMachineNum {
			serverList = append(serverList, map[string]interface{}{
				"machine_type": machineType,
				"machine_num":  count,
			})
		}
		_ = d.Set("server_list", serverList)
	}

	if clusterInfo.ProxyList != nil {
		proxyMachineNum := make(map[string]int)
		for _, proxy := range clusterInfo.ProxyList {
			if proxy == nil || proxy.MachineType == nil {
				continue
			}
			proxyMachineNum[*proxy.MachineType]++
		}
		proxyList := make([]interface{}, 0, len(proxyMachineNum))
		for machineType, count := range proxyMachineNum {
			proxyList = append(proxyList, map[string]interface{}{
				"machine_type": machineType,
				"machine_num":  count,
			})
		}
		_ = d.Set("proxy_list", proxyList)
	}

	// Read cluster tags via service layer.
	tags, tagsReqErr := service.DescribeClusterTags(ctx, clusterId)
	if tagsReqErr != nil {
		log.Printf("[CRITAL]%s read tcaplus cluster tags failed, reason:%+v", logId, tagsReqErr)
		return tagsReqErr
	}

	resourceTags := make([]interface{}, 0, len(tags))
	for _, tag := range tags {
		if tag == nil {
			continue
		}

		tagMap := map[string]interface{}{}
		if tag.TagKey != nil {
			tagMap["tag_key"] = *tag.TagKey
		}

		if tag.TagValue != nil {
			tagMap["tag_value"] = *tag.TagValue
		}

		resourceTags = append(resourceTags, tagMap)
	}

	_ = d.Set("resource_tags", resourceTags)

	return nil
}

func resourceTencentCloudTcaplusClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcaplus_cluster.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		clusterId = d.Id()
	)

	if d.HasChange("cluster_type") || d.HasChange("server_list") || d.HasChange("proxy_list") {
		request := tcaplusdb.NewModifyClusterMachineRequest()
		request.ClusterId = helper.String(clusterId)

		if v, ok := d.GetOk("cluster_type"); ok {
			request.ClusterType = helper.Int64(int64(v.(int)))
		}

		if v, ok := d.GetOk("server_list"); ok {
			servers := v.(*schema.Set).List()
			for _, server := range servers {
				serverMap := server.(map[string]interface{})
				machineType := serverMap["machine_type"].(string)
				machineNum := int64(serverMap["machine_num"].(int))
				request.ServerList = append(request.ServerList, &tcaplusdb.MachineInfo{
					MachineType: helper.String(machineType),
					MachineNum:  helper.Int64(machineNum),
				})
			}
		}

		if v, ok := d.GetOk("proxy_list"); ok {
			proxies := v.(*schema.Set).List()
			for _, proxy := range proxies {
				proxyMap := proxy.(map[string]interface{})
				machineType := proxyMap["machine_type"].(string)
				machineNum := int64(proxyMap["machine_num"].(int))
				request.ProxyList = append(request.ProxyList, &tcaplusdb.MachineInfo{
					MachineType: helper.String(machineType),
					MachineNum:  helper.Int64(machineNum),
				})
			}
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTcaplusClient().ModifyClusterMachineWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify tcaplus cluster machine failed, Response is nil."))
			}
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update tcaplus cluster machine failed, reason:%+v", logId, reqErr)
			return reqErr
		}

		// wait
		waitReq := tcaplusdb.NewDescribeClustersRequest()
		waitReq.ClusterIds = []*string{&clusterId}
		reqErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTcaplusClient().DescribeClustersWithContext(ctx, waitReq)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil || result.Response.Clusters == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe tcaplus cluster failed, Response is nil."))
			}

			if len(result.Response.Clusters) == 0 {
				return resource.NonRetryableError(fmt.Errorf("Describe tcaplus cluster failed, Clusters is empty."))
			}

			cluster := result.Response.Clusters[0]
			if cluster.ClusterStatus == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe tcaplus cluster failed, ClusterStatus is nil."))
			}

			if *cluster.ClusterStatus == 0 {
				return nil
			}

			return resource.RetryableError(fmt.Errorf("cluster is still updating, current status: %d", *cluster.ClusterStatus))
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update tcaplus cluster failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	if d.HasChange("resource_tags") {
		oldTags, newTags := d.GetChange("resource_tags")
		oldTagsList := oldTags.(*schema.Set).List()
		newTagsList := newTags.(*schema.Set).List()

		oldTagsMap := make(map[string]string)
		for _, tag := range oldTagsList {
			tagMap := tag.(map[string]interface{})
			tagKey := tagMap["tag_key"].(string)
			tagValue := tagMap["tag_value"].(string)
			oldTagsMap[tagKey] = tagValue
		}

		newTagsMap := make(map[string]string)
		for _, tag := range newTagsList {
			tagMap := tag.(map[string]interface{})
			tagKey := tagMap["tag_key"].(string)
			tagValue := tagMap["tag_value"].(string)
			newTagsMap[tagKey] = tagValue
		}

		var replaceTags []*tcaplusdb.TagInfoUnit
		var deleteTags []*tcaplusdb.TagInfoUnit

		// Tags to add or update.
		for k, v := range newTagsMap {
			if oldV, ok := oldTagsMap[k]; !ok || oldV != v {
				replaceTags = append(replaceTags, &tcaplusdb.TagInfoUnit{
					TagKey:   helper.String(k),
					TagValue: helper.String(v),
				})
			}
		}

		// Tags to delete.
		for k := range oldTagsMap {
			if _, ok := newTagsMap[k]; !ok {
				deleteTags = append(deleteTags, &tcaplusdb.TagInfoUnit{
					TagKey: helper.String(k),
				})
			}
		}

		if len(replaceTags) > 0 || len(deleteTags) > 0 {
			tagsRequest := tcaplusdb.NewModifyClusterTagsRequest()
			tagsRequest.ClusterId = helper.String(clusterId)
			tagsRequest.ReplaceTags = replaceTags
			tagsRequest.DeleteTags = deleteTags

			tagsReqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				tagsResult, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTcaplusClient().ModifyClusterTagsWithContext(ctx, tagsRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, tagsRequest.GetAction(), tagsRequest.ToJsonString(), tagsResult.ToJsonString())
				}

				if tagsResult == nil || tagsResult.Response == nil {
					return resource.NonRetryableError(fmt.Errorf("Modify tcaplus cluster tags failed, Response is nil."))
				}
				return nil
			})

			if tagsReqErr != nil {
				log.Printf("[CRITAL]%s update tcaplus cluster tags failed, reason:%+v", logId, tagsReqErr)
				return tagsReqErr
			}
		}
	}

	if d.HasChange("cluster_name") {
		request := tcaplusdb.NewModifyClusterNameRequest()
		request.ClusterId = helper.String(clusterId)

		if v, ok := d.GetOk("cluster_name"); ok {
			request.ClusterName = helper.String(v.(string))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTcaplusClient().ModifyClusterNameWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify tcaplus cluster name failed, Response is nil."))
			}
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update tcaplus cluster name failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	if d.HasChange("password") {
		oldPwd, newPwd := d.GetChange("password")

		request := tcaplusdb.NewModifyClusterPasswordRequest()
		request.ClusterId = helper.String(clusterId)
		request.OldPassword = helper.String(oldPwd.(string))
		request.NewPassword = helper.String(newPwd.(string))
		request.Mode = helper.String("1")

		if v, ok := d.GetOkExists("old_password_expire_last"); ok {
			oldPasswordExpireLast := int64(v.(int))
			if oldPasswordExpireLast > 0 {
				expireTime := time.Now().Add(time.Second * time.Duration(oldPasswordExpireLast))
				loc, err := time.LoadLocation("Asia/Shanghai")
				if err != nil {
					return fmt.Errorf("Get Asia/Shanghai time group fail, %s", err.Error())
				}
				ex := expireTime.In(loc).Format("2006-01-02 15:04:05")
				request.OldPasswordExpireTime = helper.String(ex)
			}
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTcaplusClient().ModifyClusterPasswordWithContext(ctx, request)
			if e != nil {
				if sdkerr, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
					if sdkerr.Code == "FailedOperation.OldPasswordInUse" {
						return resource.NonRetryableError(fmt.Errorf("[TencentCloudSDKError] Code=FailedOperation.OldPasswordInUse, `password_status` is unmodifiable now, can modify after `old_password_expire_time`"))
					}
				}
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify tcaplus cluster password failed, Response is nil."))
			}
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update tcaplus cluster password failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudTcaplusClusterRead(d, meta)
}

func resourceTencentCloudTcaplusClusterDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcaplus_cluster.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service   = TcaplusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		clusterId = d.Id()
		request   = tcaplusdb.NewDeleteClusterRequest()
	)

	request.ClusterId = helper.String(clusterId)

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTcaplusClient().DeleteClusterWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete tcaplus cluster failed, Response is nil."))
		}

		if result.Response.TaskId == nil || *result.Response.TaskId == "" {
			return resource.NonRetryableError(fmt.Errorf("Delete tcaplus cluster failed, TaskId is nil or empty."))
		}
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete tcaplus cluster failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// Poll for deletion completion.
	pollErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		_, has, e := service.DescribeCluster(ctx, clusterId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if has {
			return resource.RetryableError(fmt.Errorf("delete cluster fail, cluster still exist from sdk DescribeClusters"))
		}

		return nil
	})

	if pollErr != nil {
		log.Printf("[CRITAL]%s wait tcaplus cluster delete failed, reason:%+v", logId, pollErr)
		return pollErr
	}

	return nil
}
