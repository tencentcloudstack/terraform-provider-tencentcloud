package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMysqlRemoteBackupConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlRemoteBackupConfigCreate,
		Read:   resourceTencentCloudMysqlRemoteBackupConfigRead,
		Update: resourceTencentCloudMysqlRemoteBackupConfigUpdate,
		Delete: resourceTencentCloudMysqlRemoteBackupConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID, in the format: cdb-c1nl9rpv. Same instance ID as displayed in the ApsaraDB for Console page.",
			},

			"remote_backup_save": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Remote data backup switch, off - disable remote backup, on - enable remote backup.",
			},

			"remote_binlog_save": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Off-site log backup switch, off - off off-site backup, on-on off-site backup, only when the parameter RemoteBackupSave is on, the RemoteBinlogSave parameter can be set to on.",
			},

			"remote_region": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "User settings off-site backup region list.",
			},

			"expire_days": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Remote backup retention time, in days.",
			},
		},
	}
}

func resourceTencentCloudMysqlRemoteBackupConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_remote_backup_config.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudMysqlRemoteBackupConfigUpdate(d, meta)
}

func resourceTencentCloudMysqlRemoteBackupConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_remote_backup_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	remoteBackupConfig, err := service.DescribeMysqlRemoteBackupConfigById(ctx, instanceId)
	if err != nil {
		return err
	}

	if remoteBackupConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MysqlRemoteBackupConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if remoteBackupConfig.RemoteBackupSave != nil {
		_ = d.Set("remote_backup_save", remoteBackupConfig.RemoteBackupSave)
	}

	if remoteBackupConfig.RemoteBinlogSave != nil {
		_ = d.Set("remote_binlog_save", remoteBackupConfig.RemoteBinlogSave)
	}

	if remoteBackupConfig.RemoteRegion != nil {
		_ = d.Set("remote_region", remoteBackupConfig.RemoteRegion)
	}

	if remoteBackupConfig.ExpireDays != nil {
		_ = d.Set("expire_days", remoteBackupConfig.ExpireDays)
	}

	return nil
}

func resourceTencentCloudMysqlRemoteBackupConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_remote_backup_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := mysql.NewModifyRemoteBackupConfigRequest()

	instanceId := d.Id()

	request.InstanceId = &instanceId

	if v, ok := d.GetOk("remote_backup_save"); ok {
		request.RemoteBackupSave = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remote_binlog_save"); ok {
		request.RemoteBinlogSave = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remote_region"); ok {
		remoteRegionSet := v.(*schema.Set).List()
		for i := range remoteRegionSet {
			remoteRegion := remoteRegionSet[i].(string)
			request.RemoteRegion = append(request.RemoteRegion, &remoteRegion)
		}
	}

	if v, ok := d.GetOkExists("expire_days"); ok {
		request.ExpireDays = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMysqlClient().ModifyRemoteBackupConfig(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update mysql remoteBackupConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMysqlRemoteBackupConfigRead(d, meta)
}

func resourceTencentCloudMysqlRemoteBackupConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_remote_backup_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
