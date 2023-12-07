package tencentcloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	audit "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cloudaudit/v20190319"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func resourceTencentCloudAudit() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This resource has been deprecated in Terraform TencentCloud provider version 1.78.16. Please use 'tencentcloud_audit_track' instead.",
		Create:             resourceTencentCloudAuditCreate,
		Read:               resourceTencentCloudAuditRead,
		Update:             resourceTencentCloudAuditUpdate,
		Delete:             resourceTencentCloudAuditDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of audit. Valid length ranges from 3 to 128. Only alpha character or numbers or '_' supported.",
			},
			"cos_bucket": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the cos bucket to save audit log. Caution: the validation of existing cos bucket will not be checked by terraform.",
			},
			"cos_region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Region of the cos bucket.",
			},
			"enable_kms_encry": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicate whether the log is encrypt with KMS algorithm or not.",
			},
			"key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Existing CMK unique key. This field can be get by data source `tencentcloud_audit_key_alias`. Caution: the region of the KMS must be as same as the `cos_region`.",
			},
			"log_file_prefix": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The log file name prefix. The length ranges from 3 to 40. If not set, the account ID will be the log file prefix.",
			},
			"read_write_attribute": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Event attribute filter. Valid values: `1`, `2`, `3`. `1` for readonly, `2` for write-only, `3` for all.",
			},
			"audit_switch": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicate whether to turn on audit logging or not.",
			},
		},
	}
}

func resourceTencentCloudAuditCreate(d *schema.ResourceData, meta interface{}) (errRet error) {
	defer logElapsed("resource.tencentcloud_audit.create")()
	request := audit.NewCreateAuditRequest()

	name := d.Get("name").(string)
	cosBucketName := d.Get("cos_bucket").(string)
	cosRegion := d.Get("cos_region").(string)

	isEnableKmsEncry := d.Get("enable_kms_encry").(bool)
	keyId := ""
	if v, ok := d.GetOk("key_id"); ok {
		keyId = v.(string)
	}
	request.IsEnableKmsEncry = helper.BoolToInt64Ptr(isEnableKmsEncry)
	if isEnableKmsEncry {
		if keyId == "" {
			return fmt.Errorf("`key_id` must be set with valid value when `enable_kms_encry` is true")
		}
		request.KmsRegion = &cosRegion
		request.KeyId = &keyId
	} else {
		if keyId != "" {
			return fmt.Errorf("`key_id` can not be set when `enable_kms_encry` is false")
		}
	}

	readWriteAttribute := d.Get("read_write_attribute").(int)
	logFilePrefix := d.Get("log_file_prefix").(string)

	request.AuditName = &name
	request.IsCreateNewBucket = helper.BoolToInt64Ptr(false)
	request.CosBucketName = &cosBucketName
	request.CosRegion = &cosRegion
	//CMQ is not supported with terraform
	request.IsEnableCmqNotify = helper.IntInt64(0)
	request.ReadWriteAttribute = helper.IntInt64(readWriteAttribute)
	request.LogFilePrefix = &logFilePrefix

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := meta.(*TencentCloudClient).apiV3Conn.UseAuditClient().CreateAudit(request)
		if err != nil {
			return retryError(err)
		}
		if response != nil && response.Response != nil && int(*response.Response.IsSuccess) > 0 {
			d.SetId(name)
			return nil
		} else {
			return resource.NonRetryableError(fmt.Errorf("create audit %s failed", name))
		}
	})

	if err != nil {
		return nil
	}

	auditSwitch := d.Get("audit_switch").(bool)

	err = modifyAuditSwitch(name, auditSwitch, meta)
	if err != nil {
		errRet = err
		return
	}

	return resourceTencentCloudAuditRead(d, meta)
}

func resourceTencentCloudAuditRead(d *schema.ResourceData, meta interface{}) (errRet error) {
	defer logElapsed("resource.tencentcloud_audit.read")()
	defer inconsistentCheck(d, meta)()
	logId := getLogId(contextNil)
	request := audit.NewDescribeAuditRequest()

	auditId := d.Id()

	request.AuditName = &auditId

	ratelimit.Check(request.GetAction())
	var response *audit.DescribeAuditResponse
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseAuditClient().DescribeAudit(request)
		if e != nil {
			log.Printf("[CRITAL]%s %s fail, reason:%s\n", logId, request.GetAction(), e.Error())
			return retryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil || response.Response.AuditName == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("name", response.Response.AuditName)
	_ = d.Set("read_write_attribute", response.Response.ReadWriteAttribute)
	_ = d.Set("log_file_prefix", response.Response.LogFilePrefix)
	_ = d.Set("enable_kms_encry", *response.Response.IsEnableKmsEncry > 0)
	_ = d.Set("cos_region", response.Response.CosRegion)
	_ = d.Set("cos_bucket", response.Response.CosBucketName)
	if *response.Response.IsEnableKmsEncry > 0 {
		_ = d.Set("key_id", response.Response.KeyId)
	}
	_ = d.Set("audit_switch", *response.Response.AuditStatus > 0)

	return nil
}

func resourceTencentCloudAuditUpdate(d *schema.ResourceData, meta interface{}) (errRet error) {
	defer logElapsed("resource.tencentcloud_audit.update")()

	logId := getLogId(contextNil)
	request := audit.NewUpdateAuditRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	attributeSet := []string{"cos_region", "cos_bucket", "enable_kms_encry", "log_file_prefix", "read_write_attribute", "key_id"}
	attribuChange := false
	for _, attr := range attributeSet {
		if d.HasChange(attr) {
			attribuChange = true
			break
		}
	}
	d.Partial(true)
	if attribuChange {
		name := d.Get("name").(string)
		cosBucketName := d.Get("cos_bucket").(string)
		cosRegion := d.Get("cos_region").(string)

		isEnableKmsEncry := d.Get("enable_kms_encry").(bool)
		keyId := ""
		if v, ok := d.GetOk("key_id"); ok {
			keyId = v.(string)
		}
		request.IsEnableKmsEncry = helper.BoolToInt64Ptr(isEnableKmsEncry)
		if isEnableKmsEncry {
			if keyId == "" {
				return fmt.Errorf("`key_id` must be set with valid value when `enable_kms_encry` is true")
			}
			request.KmsRegion = &cosRegion
			request.KeyId = &keyId
		} else {
			if keyId != "" {
				return fmt.Errorf("`key_id` can not be set when `enable_kms_encry` is false")
			}
		}
		readWriteAttribute := d.Get("read_write_attribute").(int)
		logFilePrefix := d.Get("log_file_prefix").(string)

		request.AuditName = &name
		request.IsCreateNewBucket = helper.BoolToInt64Ptr(false)
		request.CosBucketName = &cosBucketName
		request.CosRegion = &cosRegion
		//CMQ is not supported with terraform
		request.IsEnableCmqNotify = helper.IntInt64(0)
		request.ReadWriteAttribute = helper.IntInt64(readWriteAttribute)
		request.LogFilePrefix = &logFilePrefix

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			response, err := meta.(*TencentCloudClient).apiV3Conn.UseAuditClient().UpdateAudit(request)
			if err != nil {
				return retryError(err)
			}
			if response != nil && response.Response != nil && int(*response.Response.IsSuccess) > 0 {
				return nil
			} else {
				return resource.NonRetryableError(fmt.Errorf("update audit %s failed", name))
			}
		})
		if err != nil {
			log.Printf("[CRITAL]%s delete audit %s failed, reason:%s\n", logId, name, err.Error())
			return err
		}
	}
	if d.HasChange("audit_switch") {
		auditSwitch := d.Get("audit_switch").(bool)

		err := modifyAuditSwitch(d.Id(), auditSwitch, meta)
		if err != nil {
			errRet = err
			return
		}
	}
	d.Partial(false)

	return resourceTencentCloudAuditRead(d, meta)
}

func resourceTencentCloudAuditDelete(d *schema.ResourceData, meta interface{}) (errRet error) {
	defer logElapsed("resource.tencentcloud_audit.delete")()

	logId := getLogId(contextNil)
	request := audit.NewDeleteAuditRequest()

	auditId := d.Id()

	request.AuditName = &auditId
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		_, e := meta.(*TencentCloudClient).apiV3Conn.UseAuditClient().DeleteAudit(request)
		if e != nil {
			log.Printf("[CRITAL]%s reason[%s]\n", logId, e.Error())
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete audit %s failed, reason:%s\n", logId, auditId, err.Error())
		return err
	}
	return nil
}

func modifyAuditSwitch(auditname string, auditSwitch bool, meta interface{}) (errRet error) {
	if auditSwitch {
		request := audit.NewStartLoggingRequest()
		request.AuditName = &auditname
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			response, err := meta.(*TencentCloudClient).apiV3Conn.UseAuditClient().StartLogging(request)
			if err != nil {
				return retryError(err)
			}
			if response != nil && response.Response != nil && int(*response.Response.IsSuccess) > 0 {
				return nil
			} else {
				return resource.NonRetryableError(fmt.Errorf("Start logging failed"))
			}
		})
		if err != nil {
			return err
		}
	} else {
		request := audit.NewStopLoggingRequest()
		request.AuditName = &auditname
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			response, err := meta.(*TencentCloudClient).apiV3Conn.UseAuditClient().StopLogging(request)
			if err != nil {
				return retryError(err)
			}
			if response != nil && response.Response != nil && int(*response.Response.IsSuccess) > 0 {
				return nil
			} else {
				return resource.NonRetryableError(fmt.Errorf("Stop logging failed"))
			}
		})
		if err != nil {
			return err
		}
	}
	return nil
}
