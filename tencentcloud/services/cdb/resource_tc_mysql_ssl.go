package cdb

import (
	"context"
	"fmt"
	"log"
	"strings"

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
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"ro_group_id"},
				Description:  "Instance ID. Example value: cdb-c1nl9rpv.",
			},

			"ro_group_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"instance_id"},
				Description:  "RO group ID. Example value: cdbrg-k9a6gup3.",
			},

			"status": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"ON", "OFF"}),
				Description:  "Whether to enable SSL. `ON` means enabled, `OFF` means not enabled.",
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

	var (
		instanceId string
		roGroupId  string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		if !strings.HasPrefix(instanceId, "cdb-") {
			return fmt.Errorf("`instance_id` parameter value is invalid. Example value: cdb-c1nl9rpv.")
		}
	}

	if v, ok := d.GetOk("ro_group_id"); ok {
		roGroupId = v.(string)
		if !strings.HasPrefix(roGroupId, "cdbrg-") {
			return fmt.Errorf("`ro_group_id` parameter value is invalid. Example value: cdbrg-k9a6gup3.")
		}
	}

	if instanceId != "" {
		d.SetId(instanceId)
	} else if roGroupId != "" {
		d.SetId(roGroupId)
	} else {
		return fmt.Errorf("`instance_id` or `ro_group_id` must set one of.")
	}

	return resourceTencentCloudMysqlSslUpdate(d, meta)
}

func resourceTencentCloudMysqlSslRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_ssl.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		resID      = d.Id()
		instanceId string
		roGroupId  string
	)

	if strings.HasPrefix(resID, "cdb-") {
		instanceId = resID
		_ = d.Set("instance_id", instanceId)
	} else {
		roGroupId = resID
		_ = d.Set("ro_group_id", roGroupId)
	}

	ssl, err := service.DescribeMysqlSslById(ctx, instanceId, roGroupId)
	if err != nil {
		return err
	}

	if ssl == nil {
		log.Printf("[WARN]%s resource `tencentcloud_mysql_ssl` [%s] not found, please check if it has been deleted.", logId, instanceId)
		d.SetId("")
		return nil
	}

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

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		resID      = d.Id()
		instanceId string
		roGroupId  string
	)

	if strings.HasPrefix(resID, "cdb-") {
		instanceId = resID
	} else {
		roGroupId = resID
	}

	if v, ok := d.GetOk("status"); ok {
		status := v.(string)
		if status == "ON" {
			request := mysql.NewOpenSSLRequest()
			if instanceId != "" {
				request.InstanceId = helper.String(instanceId)
			}

			if roGroupId != "" {
				request.RoGroupId = helper.String(roGroupId)
			}

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
				log.Printf("[CRITAL]%s Open mysql ssl failed, reason:%+v", logId, err)
				return err
			}
		} else if status == "OFF" {
			request := mysql.NewCloseSSLRequest()
			if instanceId != "" {
				request.InstanceId = helper.String(instanceId)
			}

			if roGroupId != "" {
				request.RoGroupId = helper.String(roGroupId)
			}

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
				log.Printf("[CRITAL]%s Close mysql ssl failed, reason:%+v", logId, err)
				return err
			}
		} else {
			return fmt.Errorf("[CRITAL]%s update mysql ssl failed, reason:your status must be ON or OFF!", logId)
		}

		// wait
		err := resource.Retry(10*tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ssl, err := service.DescribeMysqlSslById(ctx, instanceId, roGroupId)
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

	return resourceTencentCloudMysqlSslRead(d, meta)
}

func resourceTencentCloudMysqlSslDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_ssl.delete")()

	return nil
}
