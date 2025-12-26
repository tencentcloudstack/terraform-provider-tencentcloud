package bh

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bhv20230418 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bh/v20230418"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudBhDevice() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudBhDeviceCreate,
		Read:   resourceTencentCloudBhDeviceRead,
		Update: resourceTencentCloudBhDeviceUpdate,
		Delete: resourceTencentCloudBhDeviceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"device_set": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Asset parameter list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"os_name": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The operating system name can only be one of the following: Host (Linux, Windows), Database (MySQL, SQL Server, MariaDB, PostgreSQL, MongoDBReplicaSet, MongoDBSharded, Redis), or Container (TKE, EKS).",
						},
						"ip": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "IP address.",
						},
						"port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Management port.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Host name, can be empty.",
						},
						"department_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Department ID to which the asset belongs.",
						},
						"ip_port_set": {
							Type:        schema.TypeSet,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Asset multi-node: IP and port fields.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"enable_ssl": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Whether to enable SSL, 1: enable, 0: disable, only supports Redis assets.",
						},
						"ssl_cert": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "SSL certificate, required when EnableSSL is enabled.",
						},
						"ssl_cert_name": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "SSL certificate name, required when EnableSSL is enabled.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Asset instance ID.",
						},
						"ap_code": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Region to which the asset belongs.",
						},
						"ap_name": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Region name.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "VPC to which the asset belongs.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Subnet to which the asset belongs.",
						},
						"public_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Public IP.",
						},
					},
				},
			},

			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Cloud account ID to which the asset belongs.",
			},

			// computed
			"device_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the device.",
			},
		},
	}
}

func resourceTencentCloudBhDeviceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_device.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = bhv20230418.NewImportExternalDeviceRequest()
		response = bhv20230418.NewImportExternalDeviceResponse()
		deviceId string
	)

	if v, ok := d.GetOk("device_set"); ok {
		for _, item := range v.([]interface{}) {
			deviceSetMap := item.(map[string]interface{})
			externalDevice := bhv20230418.ExternalDevice{}
			if v, ok := deviceSetMap["os_name"].(string); ok && v != "" {
				externalDevice.OsName = helper.String(v)
			}

			if v, ok := deviceSetMap["ip"].(string); ok && v != "" {
				externalDevice.Ip = helper.String(v)
			}

			if v, ok := deviceSetMap["port"].(int); ok {
				externalDevice.Port = helper.IntUint64(v)
			}

			if v, ok := deviceSetMap["name"].(string); ok && v != "" {
				externalDevice.Name = helper.String(v)
			}

			if v, ok := deviceSetMap["department_id"].(string); ok && v != "" {
				externalDevice.DepartmentId = helper.String(v)
			}

			if v, ok := deviceSetMap["ip_port_set"]; ok {
				ipPortSetSet := v.(*schema.Set).List()
				for i := range ipPortSetSet {
					ipPortSet := ipPortSetSet[i].(string)
					externalDevice.IpPortSet = append(externalDevice.IpPortSet, helper.String(ipPortSet))
				}
			}

			if v, ok := deviceSetMap["enable_ssl"].(int); ok {
				externalDevice.EnableSSL = helper.IntInt64(v)
			}

			if v, ok := deviceSetMap["ssl_cert"].(string); ok && v != "" {
				externalDevice.SSLCert = helper.String(v)
			}

			if v, ok := deviceSetMap["ssl_cert_name"].(string); ok && v != "" {
				externalDevice.SSLCertName = helper.String(v)
			}

			if v, ok := deviceSetMap["instance_id"].(string); ok && v != "" {
				externalDevice.InstanceId = helper.String(v)
			}

			if v, ok := deviceSetMap["ap_code"].(string); ok && v != "" {
				externalDevice.ApCode = helper.String(v)
			}

			if v, ok := deviceSetMap["ap_name"].(string); ok && v != "" {
				externalDevice.ApName = helper.String(v)
			}

			if v, ok := deviceSetMap["vpc_id"].(string); ok && v != "" {
				externalDevice.VpcId = helper.String(v)
			}

			if v, ok := deviceSetMap["subnet_id"].(string); ok && v != "" {
				externalDevice.SubnetId = helper.String(v)
			}

			if v, ok := deviceSetMap["public_ip"].(string); ok && v != "" {
				externalDevice.PublicIp = helper.String(v)
			}

			request.DeviceSet = append(request.DeviceSet, &externalDevice)
		}
	}

	if v, ok := d.GetOkExists("account_id"); ok {
		request.AccountId = helper.IntUint64(v.(int))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhV20230418Client().ImportExternalDeviceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Import external device failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create bh device failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.DeviceIdSet == nil || len(response.Response.DeviceIdSet) == 0 {
		return fmt.Errorf("DeviceIdSet is nil.")
	}

	deviceId = helper.UInt64ToStr(*response.Response.DeviceIdSet[0])
	d.SetId(deviceId)
	return resourceTencentCloudBhDeviceRead(d, meta)
}

func resourceTencentCloudBhDeviceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_device.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service  = BhService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		deviceId = d.Id()
	)

	respData, err := service.DescribeBhDeviceById(ctx, deviceId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_bh_device` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	dMap := make(map[string]interface{}, 0)
	if respData.OsName != nil {
		dMap["os_name"] = respData.OsName
	}

	if respData.PrivateIp != nil {
		dMap["ip"] = respData.PrivateIp
	}

	if respData.Port != nil {
		dMap["port"] = respData.Port
	}

	if respData.Name != nil {
		dMap["name"] = respData.Name
	}

	if respData.Department != nil {
		if respData.Department.Id != nil {
			dMap["department_id"] = respData.Department.Id
		}
	}

	if respData.IpPortSet != nil {
		dMap["ip_port_set"] = respData.IpPortSet
	}

	if respData.EnableSSL != nil {
		dMap["enable_ssl"] = respData.EnableSSL
	}

	if respData.SSLCertName != nil {
		dMap["ssl_cert_name"] = respData.SSLCertName
	}

	if respData.InstanceId != nil {
		dMap["instance_id"] = respData.InstanceId
	}

	if respData.ApCode != nil {
		dMap["ap_code"] = respData.ApCode
	}

	if respData.ApName != nil {
		dMap["ap_name"] = respData.ApName
	}

	if respData.VpcId != nil {
		dMap["vpc_id"] = respData.VpcId
	}

	if respData.SubnetId != nil {
		dMap["subnet_id"] = respData.SubnetId
	}

	if respData.PublicIp != nil {
		dMap["public_ip"] = respData.PublicIp
	}

	_ = d.Set("device_set", []interface{}{dMap})

	if respData.CloudAccountId != nil {
		_ = d.Set("account_id", respData.CloudAccountId)
	}

	return nil
}

func resourceTencentCloudBhDeviceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_device.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		deviceId = d.Id()
	)

	needChange := false
	mutableArgs := []string{"device_set.0.port", "device_set.0.department_id"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := bhv20230418.NewModifyDeviceRequest()
		if v, ok := d.GetOkExists("device_set.0.port"); ok {
			request.Port = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("device_set.0.department_id"); ok {
			request.DepartmentId = helper.String(v.(string))
		}

		request.Id = helper.StrToUint64Point(deviceId)
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhV20230418Client().ModifyDeviceWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update bh device failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudBhDeviceRead(d, meta)
}

func resourceTencentCloudBhDeviceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_device.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = bhv20230418.NewDeleteDevicesRequest()
		deviceId = d.Id()
	)

	request.IdSet = append(request.IdSet, helper.StrToUint64Point(deviceId))
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhV20230418Client().DeleteDevicesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete bh device failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
