package cynosdb

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCynosdbSsl() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbSslCreate,
		Read:   resourceTencentCloudCynosdbSslRead,
		Update: resourceTencentCloudCynosdbSslUpdate,
		Delete: resourceTencentCloudCynosdbSslDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster id.",
			},
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},
			"status": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Whether to enable SSL. `ON` means enabled, `OFF` means not enabled.",
			},
			"download_url": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Certificate download address.",
			},
		},
	}
}

func resourceTencentCloudCynosdbSslCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_ssl.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	clusterId := d.Get("cluster_id").(string)
	instanceId := d.Get("instance_id").(string)

	d.SetId(clusterId + tccommon.FILED_SP + instanceId)
	return resourceTencentCloudCynosdbSslUpdate(d, meta)
}

func resourceTencentCloudCynosdbSslRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_ssl.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	instanceId := idSplit[1]

	ssl, err := service.DescribeSSLStatus(ctx, clusterId, instanceId)
	if err != nil {
		return err
	}

	if ssl == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `tencentcloud_cynosdb_ssl` [%s] not found, please check if it has been deleted.",
			logId, instanceId,
		)
		return nil
	}

	_ = d.Set("cluster_id", clusterId)
	_ = d.Set("instance_id", instanceId)

	if ssl.IsOpenSSL != nil {
		if *ssl.IsOpenSSL == "yes" {
			_ = d.Set("status", "ON")
		} else {
			_ = d.Set("status", "OFF")
		}
	}
	if ssl.DownloadUrl != nil {
		_ = d.Set("download_url", ssl.DownloadUrl)
	}

	return nil
}

func resourceTencentCloudCynosdbSslUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_ssl.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	instanceId := idSplit[1]

	var taskId *int64
	if v, ok := d.GetOk("status"); ok {
		status := v.(string)
		if status == "ON" {
			request := cynosdb.NewOpenSSLRequest()
			request.ClusterId = helper.String(clusterId)
			request.InstanceId = helper.String(instanceId)

			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().OpenSSL(request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				taskId = result.Response.TaskId
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s update cynosdb ssl failed, reason:%+v", logId, err)
				return err
			}
		} else if status == "OFF" {
			request := cynosdb.NewCloseSSLRequest()
			request.ClusterId = helper.String(clusterId)
			request.InstanceId = helper.String(instanceId)

			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().CloseSSL(request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				taskId = result.Response.TaskId
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s update cynosdb ssl failed, reason:%+v", logId, err)
				return err
			}
		} else {
			return fmt.Errorf("[CRITAL]%s update cynosdb ssl failed, reason:your status must be ON or OFF!", logId)
		}
		service := CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		conf := tccommon.BuildStateChangeConf([]string{}, []string{"success"}, 10*tccommon.ReadRetryTimeout, time.Second, service.taskStateRefreshFunc(strconv.FormatInt(*taskId, 10), []string{}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}

	return resourceTencentCloudCynosdbSslRead(d, meta)
}

func resourceTencentCloudCynosdbSslDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_ssl.delete")()

	return nil
}
