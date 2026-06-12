package cynosdb

import (
	"context"
	"fmt"
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCynosdbClusterReadOnlyInstanceGroupAccesOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbClusterReadOnlyInstanceGroupAccesOperationCreate,
		Read:   resourceTencentCloudCynosdbClusterReadOnlyInstanceGroupAccesOperationRead,
		Delete: resourceTencentCloudCynosdbClusterReadOnlyInstanceGroupAccesOperationDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(300 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},
			"port": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Port.",
			},
			"security_group_ids": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Security group IDs.",
			},
		},
	}
}

func resourceTencentCloudCynosdbClusterReadOnlyInstanceGroupAccesOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_cluster_read_only_instance_group_acces_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		request   = cynosdb.NewOpenClusterReadOnlyInstanceGroupAccessRequest()
		clusterId string
		flowId    int64
	)

	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(clusterId)
	}

	if v, ok := d.GetOk("port"); ok {
		request.Port = helper.String(v.(string))
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroupIds := v.([]interface{})
		for _, item := range securityGroupIds {
			request.SecurityGroupIds = append(request.SecurityGroupIds, helper.String(item.(string)))
		}
	}

	var response *cynosdb.OpenClusterReadOnlyInstanceGroupAccessResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().OpenClusterReadOnlyInstanceGroupAccessWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cynosdb_cluster_read_only_instance_group_acces_operation failed, reason:%+v", logId, err)
		return err
	}

	log.Printf("[DEBUG]%s logId=%s, d.Id()=%s before checking response", logId, logId, d.Id())

	if response == nil || response.Response == nil || response.Response.FlowId == nil {
		return fmt.Errorf("create cynosdb_cluster_read_only_instance_group_acces_operation failed, Response or FlowId is nil")
	}

	d.SetId(clusterId)

	// wait
	flowId = *response.Response.FlowId
	service := CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		ok, e := service.DescribeFlow(ctx, flowId)
		if e != nil {
			if _, ok := e.(*sdkErrors.TencentCloudSDKError); !ok {
				return resource.RetryableError(e)
			} else {
				return resource.NonRetryableError(e)
			}
		}

		if ok {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("cynosdb_cluster_read_only_instance_group_acces_operation flow %d is still processing", flowId))
	})

	if err != nil {
		log.Printf("[CRITAL]%s cynosdb_cluster_read_only_instance_group_acces_operation flow polling failed, reason:%s\n", logId, err.Error())
		return err
	}

	return resourceTencentCloudCynosdbClusterReadOnlyInstanceGroupAccesOperationRead(d, meta)
}

func resourceTencentCloudCynosdbClusterReadOnlyInstanceGroupAccesOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_cluster_read_only_instance_group_acces_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCynosdbClusterReadOnlyInstanceGroupAccesOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_cluster_read_only_instance_group_acces_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
