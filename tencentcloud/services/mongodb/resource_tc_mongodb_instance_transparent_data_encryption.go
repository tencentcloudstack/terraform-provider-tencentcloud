package mongodb

import (
	"context"
	"fmt"
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMongodbInstanceTransparentDataEncryption() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMongodbInstanceTransparentDataEncryptionCreate,
		Read:   resourceTencentCloudMongodbInstanceTransparentDataEncryptionRead,
		Update: resourceTencentCloudMongodbInstanceTransparentDataEncryptionUpdate,
		Delete: resourceTencentCloudMongodbInstanceTransparentDataEncryptionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID, for example: cmgo-p8vn ***. Currently supported general versions include: 4.4 and 5.0, but the cloud disk version is not currently supported.",
			},

			"kms_region": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The region where the Key Management Service (KMS) serves, such as ap-shanghai.",
			},

			"key_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Key ID. If this parameter is not set and the specific key ID is not specified, Tencent Cloud will automatically generate the key and this key will be beyond the control of Terraform.",
			},
			"transparent_data_encryption_status": {
				Computed: true,
				Type:     schema.TypeString,
				Description: "Represents whether transparent encryption is turned on. Valid values:\n" +
					"- close: Not opened;\n" +
					"- open: It has been opened.",
			},
			"key_info_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of bound keys.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Master Key ID.",
						},
						"key_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Master key name.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance and key binding time.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Key status.",
						},
						"key_usage": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Purpose of the key.",
						},
						"key_origin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Key origin.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudMongodbInstanceTransparentDataEncryptionCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_transparent_data_encryption.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = mongodb.NewEnableTransparentDataEncryptionRequest()
		response = mongodb.NewEnableTransparentDataEncryptionResponse()
	)
	instanceId := d.Get("instance_id").(string)
	request.InstanceId = helper.String(instanceId)
	request.KmsRegion = helper.String(d.Get("kms_region").(string))

	if v, ok := d.GetOk("key_id"); ok {
		request.KeyId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMongodbClient().EnableTransparentDataEncryption(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s enable transparent data encryption failed, reason:%+v", logId, err)
		return nil
	}

	flowId := *response.Response.FlowId
	flowIdString := strconv.FormatInt(flowId, 10)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	timeout := d.Timeout(schema.TimeoutCreate)
	if response != nil && response.Response != nil {
		if err = service.DescribeAsyncRequestInfo(ctx, flowIdString, timeout); err != nil {
			return err
		}
	}

	d.SetId(instanceId)

	return resourceTencentCloudMongodbInstanceTransparentDataEncryptionRead(d, meta)
}

func resourceTencentCloudMongodbInstanceTransparentDataEncryptionRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_transparent_data_encryption.read")()
	defer tccommon.InconsistentCheck(d, meta)()
	request := mongodb.NewDescribeTransparentDataEncryptionStatusRequest()
	request.InstanceId = helper.String(d.Id())
	ratelimit.Check(request.GetAction())
	response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMongodbClient().DescribeTransparentDataEncryptionStatus(request)
	if err != nil {
		return err
	}
	_ = d.Set("instance_id", d.Id())
	_ = d.Set("transparent_data_encryption_status", response.Response.TransparentDataEncryptionStatus)
	var kmsRegion string
	if response.Response != nil && len(response.Response.KeyInfoList) > 0 {
		kmsInfoList := make([]map[string]interface{}, 0)
		for _, kmsInfoDetail := range response.Response.KeyInfoList {
			kmsInfoDetailMap := make(map[string]interface{})
			kmsInfoDetailMap["key_id"] = kmsInfoDetail.KeyId
			kmsInfoDetailMap["key_name"] = kmsInfoDetail.KeyName
			kmsInfoDetailMap["create_time"] = kmsInfoDetail.CreateTime
			kmsInfoDetailMap["status"] = kmsInfoDetail.Status
			kmsInfoDetailMap["key_usage"] = kmsInfoDetail.KeyUsage
			kmsInfoDetailMap["key_origin"] = kmsInfoDetail.KeyOrigin
			kmsInfoList = append(kmsInfoList, kmsInfoDetailMap)
			if kmsInfoDetail.KmsRegion != nil {
				kmsRegion = *kmsInfoDetail.KmsRegion
			}
		}
		_ = d.Set("key_info_list", kmsInfoList)
		_ = d.Set("kms_region", kmsRegion)
	}

	return nil
}

func resourceTencentCloudMongodbInstanceTransparentDataEncryptionUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_transparent_data_encryption.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	immutableFields := []string{"instance_id", "kms_region", "key_id"}
	for _, f := range immutableFields {
		if d.HasChange(f) {
			return fmt.Errorf("cannot update argument `%s`", f)
		}
	}
	return resourceTencentCloudMongodbInstanceTransparentDataEncryptionRead(d, meta)
}

func resourceTencentCloudMongodbInstanceTransparentDataEncryptionDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_transparent_data_encryption.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
