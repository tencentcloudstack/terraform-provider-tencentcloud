package cynosdb

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdbv20190107 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCynosdbClusterTransparentEncrypt() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbClusterTransparentEncryptCreate,
		Read:   resourceTencentCloudCynosdbClusterTransparentEncryptRead,
		Update: resourceTencentCloudCynosdbClusterTransparentEncryptUpdate,
		Delete: resourceTencentCloudCynosdbClusterTransparentEncryptDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cluster ID.",
			},

			"key_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Key type (cloud, custom).",
			},

			"key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Key Id.",
			},

			"key_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Key region.",
			},

			"is_open_global_encryption": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable global encryption.",
			},
		},
	}
}

func resourceTencentCloudCynosdbClusterTransparentEncryptCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_cluster_transparent_encrypt.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		clusterId string
	)

	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
	}

	d.SetId(clusterId)

	return resourceTencentCloudCynosdbClusterTransparentEncryptUpdate(d, meta)
}

func resourceTencentCloudCynosdbClusterTransparentEncryptRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_cluster_transparent_encrypt.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service   = CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		clusterId = d.Id()
	)

	respData, err := service.DescribeCynosdbClusterTransparentEncryptById(ctx, clusterId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_cynosdb_cluster_transparent_encrypt` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("cluster_id", clusterId)

	if respData.KeyId != nil {
		_ = d.Set("key_id", respData.KeyId)
	}

	if respData.KeyRegion != nil {
		_ = d.Set("key_region", respData.KeyRegion)
	}

	if respData.KeyType != nil {
		_ = d.Set("key_type", respData.KeyType)
	}

	if respData.IsOpenGlobalEncryption != nil {
		_ = d.Set("is_open_global_encryption", respData.IsOpenGlobalEncryption)
	}

	return nil
}

func resourceTencentCloudCynosdbClusterTransparentEncryptUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_cluster_transparent_encrypt.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		clusterId = d.Id()
	)

	needChange := false
	mutableArgs := []string{"key_type", "key_id", "key_region"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := cynosdbv20190107.NewOpenClusterTransparentEncryptRequest()
		response := cynosdbv20190107.NewOpenClusterTransparentEncryptResponse()
		if v, ok := d.GetOk("key_type"); ok {
			request.KeyType = helper.String(v.(string))
		}

		if v, ok := d.GetOk("key_id"); ok {
			request.KeyId = helper.String(v.(string))
		}

		if v, ok := d.GetOk("key_region"); ok {
			request.KeyRegion = helper.String(v.(string))
		}

		request.ClusterId = &clusterId
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().OpenClusterTransparentEncryptWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil || result.Response.TaskId == nil {
				return resource.NonRetryableError(fmt.Errorf("Update cynosdb cluster transparent encrypt failed, Response is nil."))
			}

			response = result
			return nil
		})
		if reqErr != nil {
			log.Printf("[CRITAL]%s update cynosdb cluster transparent encrypt failed, reason:%+v", logId, reqErr)
			return reqErr
		}

		taskId := *response.Response.TaskId
		service := CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		conf := tccommon.BuildStateChangeConf([]string{}, []string{"success"}, 10*tccommon.ReadRetryTimeout, time.Second, service.taskStateRefreshFunc(strconv.FormatInt(taskId, 10), []string{}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}

	if d.HasChange("is_open_global_encryption") {
		request := cynosdbv20190107.NewModifyClusterGlobalEncryptionRequest()
		response := cynosdbv20190107.NewModifyClusterGlobalEncryptionResponse()
		if v, ok := d.GetOkExists("is_open_global_encryption"); ok {
			request.IsOpenGlobalEncryption = helper.Bool(v.(bool))
		}

		request.ClusterId = &clusterId
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().ModifyClusterGlobalEncryptionWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil || result.Response.TaskId == nil {
				return resource.NonRetryableError(fmt.Errorf("Update cynosdb cluster is_open_global_encryption failed, Response is nil."))
			}

			response = result
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update cynosdb cluster is_open_global_encryption failed, reason:%+v", logId, reqErr)
			return reqErr
		}

		taskId := *response.Response.TaskId
		service := CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		conf := tccommon.BuildStateChangeConf([]string{}, []string{"success"}, 10*tccommon.ReadRetryTimeout, time.Second, service.taskStateRefreshFunc(strconv.FormatInt(taskId, 10), []string{}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}

	return resourceTencentCloudCynosdbClusterTransparentEncryptRead(d, meta)
}

func resourceTencentCloudCynosdbClusterTransparentEncryptDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_cluster_transparent_encrypt.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
