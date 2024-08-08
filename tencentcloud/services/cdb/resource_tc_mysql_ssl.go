package cdb

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMysqlSsl() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlSslCreate,
		Read:   resourceTencentCloudMysqlSslRead,
		Update: resourceTencentCloudMysqlSslUpdate,
		Delete: resourceTencentCloudMysqlSslDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID. Example value: cdb-c1nl9rpv.",
			},

			"status": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Whether to enable SSL. `ON` means enabled, `OFF` means not enabled.",
			},

			"url": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The certificate download link. Example value: http://testdownload.url.",
			},
		},
	}
}

func resourceTencentCloudMysqlSslCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_ssl.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	d.SetId(d.Get("instance_id").(string))

	return resourceTencentCloudMysqlSslUpdate(d, meta)
}

func resourceTencentCloudMysqlSslRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_ssl.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	instanceId := d.Id()

	ssl, err := service.DescribeMysqlSslById(ctx, instanceId)
	if err != nil {
		return err
	}

	if ssl == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `tencentcloud_mysql_ssl` [%s] not found, please check if it has been deleted.",
			logId, instanceId,
		)
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if ssl.Status != nil {
		_ = d.Set("status", ssl.Status)
	}

	if ssl.Url != nil {
		_ = d.Set("url", ssl.Url)
	}

	return nil
}

func resourceTencentCloudMysqlSslUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_ssl.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	instanceId := d.Id()

	status := ""
	if v, ok := d.GetOk("status"); ok {
		status = v.(string)
		if status == "ON" {
			request := mysql.NewOpenSSLRequest()
			request.InstanceId = helper.String(instanceId)

			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().OpenSSL(request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s update mysql ssl failed, reason:%+v", logId, err)
				return err
			}
		} else if status == "OFF" {
			request := mysql.NewCloseSSLRequest()
			request.InstanceId = helper.String(instanceId)

			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().CloseSSL(request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s update mysql ssl failed, reason:%+v", logId, err)
				return err
			}
		} else {
			return fmt.Errorf("[CRITAL]%s update mysql ssl failed, reason:your status must be ON or OFF!", logId)
		}

		if status != "" {
			service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
			err := resource.Retry(7*tccommon.ReadRetryTimeout, func() *resource.RetryError {
				ssl, err := service.DescribeMysqlSslById(ctx, instanceId)
				if err != nil {
					return resource.NonRetryableError(err)
				}
				if ssl == nil {
					err = fmt.Errorf("mysqlid %s instance ssl not exists", instanceId)
					return resource.NonRetryableError(err)
				}
				if *ssl.Status != status {
					return resource.RetryableError(fmt.Errorf("mysql ssl status is (%v)", *ssl.Status))
				}
				if *ssl.Status == status {
					return nil
				}
				err = fmt.Errorf("mysql ssl status is %v,we won't wait for it finish", *ssl.Status)
				return resource.NonRetryableError(err)
			})

			if err != nil {
				log.Printf("[CRITAL]%s mysql switchForUpgrade fail, reason:%s\n ", logId, err.Error())
				return err
			}
		}
	}

	return resourceTencentCloudMysqlSslRead(d, meta)
}

func resourceTencentCloudMysqlSslDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_ssl.delete")()

	return nil
}
