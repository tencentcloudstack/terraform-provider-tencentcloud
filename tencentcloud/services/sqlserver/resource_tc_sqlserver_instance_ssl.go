package sqlserver

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSqlserverInstanceSsl() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverInstanceSslCreate,
		Read:   resourceTencentCloudSqlserverInstanceSslRead,
		Update: resourceTencentCloudSqlserverInstanceSslUpdate,
		Delete: resourceTencentCloudSqlserverInstanceSslDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
			"type": {
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: tccommon.ValidateAllowedStringValue(SSL_TYPE),
				Description:  "Operation type. enable: turn on SSL; disable: turn off SSL; renew: update the certificate validity period.",
			},
		},
	}
}

func resourceTencentCloudSqlserverInstanceSslCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_instance_ssl.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var instanceId string

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudSqlserverInstanceSslUpdate(d, meta)
}

func resourceTencentCloudSqlserverInstanceSslRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_instance_ssl.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
		sslType    string
	)

	if v, ok := d.GetOk("type"); ok {
		sslType = v.(string)
	}

	instanceSsl, err := service.DescribeSqlserverInstanceSslById(ctx, instanceId)
	if err != nil {
		return err
	}

	if instanceSsl == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverInstanceSsl` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if instanceSsl.InstanceId != nil {
		_ = d.Set("instance_id", instanceSsl.InstanceId)
	}

	if instanceSsl.SSLConfig.Encryption != nil {
		if sslType == SSL_TYPE_ENABLE || sslType == SSL_TYPE_DISABLE {
			_ = d.Set("type", instanceSsl.SSLConfig.Encryption)
		} else {
			_ = d.Set("type", SSL_TYPE_RENEW)
		}
	}

	return nil
}

func resourceTencentCloudSqlserverInstanceSslUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_instance_ssl.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		request     = sqlserver.NewModifyDBInstanceSSLRequest()
		attrRequest = sqlserver.NewDescribeDBInstancesAttributeRequest()
		instanceId  = d.Id()
		sslType     string
	)

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
		sslType = v.(string)
	}

	request.InstanceId = &instanceId
	request.WaitSwitch = helper.IntUint64(0)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().ModifyDBInstanceSSL(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver instanceSsl failed, reason:%+v", logId, err)
		return err
	}

	// wait
	attrRequest.InstanceId = &instanceId
	err = resource.Retry(tccommon.ReadRetryTimeout*10, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().DescribeDBInstancesAttribute(attrRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, attrRequest.GetAction(), attrRequest.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("Describe DBInstances Attribute not exists.")
			return resource.NonRetryableError(e)
		}

		if sslType == SSL_TYPE_ENABLE {
			if *result.Response.SSLConfig.Encryption == SSL_TYPE_ENABLE {
				return nil
			}
		} else if sslType == SSL_TYPE_DISABLE {
			if *result.Response.SSLConfig.Encryption == SSL_TYPE_DISABLE {
				return nil
			}
		} else {
			if *result.Response.SSLConfig.Encryption == SSL_TYPE_ENABLE {
				return nil
			}
		}

		return resource.RetryableError(fmt.Errorf("Modify DBInstance SSL is processing, status is: %s", *result.Response.SSLConfig.Encryption))
	})

	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver instanceSsl failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverInstanceSslRead(d, meta)
}

func resourceTencentCloudSqlserverInstanceSslDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_instance_ssl.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
