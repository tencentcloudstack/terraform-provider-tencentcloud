package cynosdb

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
)

func ResourceTencentCloudCynosdbIsolateInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbIsolateInstanceCreate,
		Read:   resourceTencentCloudCynosdbIsolateInstanceRead,
		Update: resourceTencentCloudCynosdbIsolateInstanceUpdate,
		Delete: resourceTencentCloudCynosdbIsolateInstanceDelete,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},

			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"operate": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "isolate, activate.",
			},
		},
	}
}

func resourceTencentCloudCynosdbIsolateInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_isolate_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		clusterId  string
		instanceId string
	)
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
	}

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(clusterId + tccommon.FILED_SP + instanceId)

	return resourceTencentCloudCynosdbIsolateInstanceUpdate(d, meta)
}

func resourceTencentCloudCynosdbIsolateInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_isolate_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCynosdbIsolateInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_isolate_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	instanceId := idSplit[1]

	var operate string
	if v, ok := d.GetOk("operate"); ok {
		operate = v.(string)
	}

	var flowId int64
	if operate == "isolate" {
		request := cynosdb.NewIsolateInstanceRequest()
		request.ClusterId = &clusterId
		request.InstanceIdList = []*string{&instanceId}
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().IsolateInstance(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			flowId = *result.Response.FlowId
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s isolate cynosdb instance failed, reason:%+v", logId, err)
			return err
		}
	} else if operate == "activate" {
		request := cynosdb.NewActivateInstanceRequest()
		request.ClusterId = &clusterId
		request.InstanceIdList = []*string{&instanceId}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().ActivateInstance(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			flowId = *result.Response.FlowId
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s activate cynosdb instance failed, reason:%+v", logId, err)
			return err
		}
	} else {
		return fmt.Errorf("[CRITAL]%s Operation type error", logId)
	}

	service := CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err := resource.Retry(6*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ok, err := service.DescribeFlow(ctx, flowId)
		if err != nil {
			if _, ok := err.(*sdkErrors.TencentCloudSDKError); !ok {
				return resource.RetryableError(err)
			} else {
				return resource.NonRetryableError(err)
			}
		}
		if ok {
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("isolate or activate cynosdb instance is processing"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s isolate or activate cynosdb instance fail, reason:%s\n", logId, err.Error())
		return err
	}

	return resourceTencentCloudCynosdbIsolateInstanceRead(d, meta)
}

func resourceTencentCloudCynosdbIsolateInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_isolate_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
