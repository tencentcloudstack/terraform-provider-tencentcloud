package tke

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tkev20180525 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudKubernetesLogConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudKubernetesLogConfigCreate,
		Read:   resourceTencentCloudKubernetesLogConfigRead,
		Delete: resourceTencentCloudKubernetesLogConfigDelete,
		Schema: map[string]*schema.Schema{
			"log_config": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "JSON expression of log collection configuration. For more details, please refer to the guide: https://www.tencentcloud.com/zh/document/product/457/64846.",
			},

			"log_config_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Log config name.",
			},

			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cluster ID.",
			},

			"logset_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "CLS log set ID.",
			},

			"cluster_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "tke",
				Description: "The current cluster type supports tke and eks, default is tke.",
			},
		},
	}
}

func resourceTencentCloudKubernetesLogConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_log_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request       = tkev20180525.NewCreateCLSLogConfigRequest()
		clusterId     string
		logConfigName string
		clusterType   string
	)

	if v, ok := d.GetOk("log_config_name"); ok {
		logConfigName = v.(string)
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterId = helper.String(v.(string))
		clusterId = v.(string)
	}

	if v, ok := d.GetOk("log_config"); ok {
		request.LogConfig = helper.String(v.(string))
	}

	if v, ok := d.GetOk("logset_id"); ok {
		request.LogsetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cluster_type"); ok {
		request.ClusterType = helper.String(v.(string))
		clusterType = v.(string)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().CreateCLSLogConfigWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create kubernetes log config failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{clusterId, logConfigName, clusterType}, tccommon.FILED_SP))

	// wait
	waitReq := tkev20180525.NewDescribeLogConfigsRequest()
	waitReq.ClusterId = &clusterId
	waitReq.ClusterType = &clusterType
	waitReq.LogConfigNames = &logConfigName
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().DescribeLogConfigsWithContext(ctx, waitReq)
		if e != nil {
			if v, ok := e.(*errors.TencentCloudSDKError); ok {
				if v.GetCode() == "FailedOperation.KubernetesGetOperationError" {
					return resource.RetryableError(fmt.Errorf("Waiting for kubernetes log config to be created ready."))
				}
			}

			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, waitReq.GetAction(), waitReq.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe log configs failed. Response is nil."))
		}

		if result.Response.LogConfigs != nil {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Waiting for kubernetes log config to be created ready."))
	})

	if err != nil {
		log.Printf("[CRITAL]%s create kubernetes log config failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudKubernetesLogConfigRead(d, meta)
}

func resourceTencentCloudKubernetesLogConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_log_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	clusterId := idSplit[0]
	logConfigName := idSplit[1]
	clusterType := idSplit[2]

	_ = d.Set("cluster_id", clusterId)
	_ = d.Set("log_config_name", logConfigName)
	_ = d.Set("cluster_type", clusterType)

	var respData *tkev20180525.DescribeLogConfigsResponseParams
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKubernetesLogConfigById(ctx, clusterId, logConfigName, clusterType)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if err := resourceTencentCloudKubernetesLogConfigReadRequestOnSuccess0(ctx, result); err != nil {
			return err
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s read kubernetes log config failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_kubernetes_log_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	return nil
}

func resourceTencentCloudKubernetesLogConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_log_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = tkev20180525.NewDeleteLogConfigsRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	clusterId := idSplit[0]
	logConfigName := idSplit[1]
	clusterType := idSplit[2]

	request.ClusterId = helper.String(clusterId)
	request.LogConfigNames = helper.String(logConfigName)
	request.ClusterType = helper.String(clusterType)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().DeleteLogConfigsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		if err := resourceTencentCloudKubernetesLogConfigDeletePostRequest0(ctx, request, result); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete kubernetes log config failed, reason:%+v", logId, err)
		return err
	}

	// wait
	waitReq := tkev20180525.NewDescribeLogConfigsRequest()
	waitReq.ClusterId = &clusterId
	waitReq.ClusterType = &clusterType
	waitReq.LogConfigNames = &logConfigName
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().DescribeLogConfigsWithContext(ctx, waitReq)
		if e != nil {
			if v, ok := e.(*errors.TencentCloudSDKError); ok {
				if v.GetCode() == "FailedOperation.KubernetesGetOperationError" {
					return nil
				}
			}

			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, waitReq.GetAction(), waitReq.ToJsonString(), result.ToJsonString())
		}

		return resource.RetryableError(fmt.Errorf("Waiting for kubernetes log config to be deleted ready."))
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete kubernetes log config failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
