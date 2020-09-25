/*
Provides an elasticsearch instance resource.

Example Usage

```hcl
resource "tencentcloud_elasticsearch_instance" "foo" {
  instance_name     = "tf-test"
  availability_zone = "ap-guangzhou-3"
  version           = "7.5.1"
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  password          = "Test12345"
  license_type      = "oss"

  node_info_list {
    node_num  = 2
    node_type = "ES.S1.SMALL2"
	encrypt = false
  }

  tags = {
    test = "test"
  }
}
```

Import

Elasticsearch instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_elasticsearch_instance.foo es-17634f05
```
*/
package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	es "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/es/v20180416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func resourceTencentCloudElasticsearchInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudElasticsearchInstanceCreate,
		Read:   resourceTencentCloudElasticsearchInstanceRead,
		Update: resourceTencentCloudElasticsearchInstanceUpdate,
		Delete: resourceTencentCloudElasticsearchInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(1, 50),
				Description:  "Name of the instance, which can contain 1 to 50 English letters, Chinese characters, digits, dashes(-), or underscores(_).",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Availability zone.",
			},
			"version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Version of the instance. Valid values are `5.6.4`, `6.4.3`, `6.8.2` and `7.5.1`.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of a VPC network.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of a VPC subnetwork.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Password to an instance.",
			},
			"charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      ES_CHARGE_TYPE_POSTPAID_BY_HOUR,
				ValidateFunc: validateAllowedStringValue(ES_CHARGE_TYPE),
				Description:  "The charge type of instance. Valid values are `PREPAID` and `POSTPAID_BY_HOUR`.",
			},
			"charge_period": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     0,
				Description: "The tenancy of the prepaid instance, and uint is month. NOTE: it only works when charge_type is set to `PREPAID`.",
			},
			"renew_flag": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(ES_RENEW_FLAG),
				Description:  "When enabled, the instance will be renew automatically when it reach the end of the prepaid tenancy. Valid values are `RENEW_FLAG_AUTO` and `RENEW_FLAG_MANUAL`. NOTE: it only works when charge_type is set to `PREPAID`.",
			},
			"deploy_mode": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      ES_DEPLOY_MODE_SINGLE_REGION,
				ValidateFunc: validateAllowedIntValue(ES_DEPLOY_MODE),
				Description:  "Cluster deployment mode. Valid values are `0` and `1`, `0` is single-AZ deployment, and `1` is multi-AZ deployment. Default value is `0`.",
			},
			"multi_zone_infos": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Details of AZs in multi-AZ deployment mode (which is required when deploy_mode is `1`).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Availability zone.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The id of a VPC subnetwork.",
						},
					},
				},
			},
			"license_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      ES_LICENSE_TYPE_PLATINUM,
				ValidateFunc: validateAllowedStringValue(ES_LICENSE_TYPE),
				Description:  "License type. Valid values are `oss`, `basic` and `platinum`, and default value is `platinum`.",
			},
			"node_info_list": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				ForceNew:    true,
				Description: "Node information list, which is used to describe the specification information of various types of nodes in the cluster, such as node type, node quantity, node specification, disk type, and disk size.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_num": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Number of nodes.",
						},
						"node_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Node specification, and valid values refer to [document of tencentcloud](https://intl.cloud.tencent.com/document/product/845/18376).",
						},
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      ES_NODE_TYPE_HOT_DATA,
							ValidateFunc: validateAllowedStringValue(ES_NODE_TYPE),
							Description:  "Node type. Valid values are `hotData`, `warmData` and `dedicatedMaster`, and default value is 'hotData`.",
						},
						"disk_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CVM_DISK_TYPE_CLOUD_SSD,
							ValidateFunc: validateAllowedStringValue(ES_NODE_DISK_TYPE),
							Description:  "Node disk type. Valid values are `CLOUD_SSD` and `CLOUD_PREMIUM`, and default value is `CLOUD_SSD`.",
						},
						"disk_size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     100,
							Description: "Node disk size. Unit is GB, and default value is `100`.",
						},
						"encrypt": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Decides to encrypt this disk or not.",
						},
					},
				},
			},
			"basic_security_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      ES_BASIC_SECURITY_TYPE_OFF,
				ValidateFunc: validateAllowedIntValue(ES_BASIC_SECURITY_TYPE),
				Description:  "Whether to enable X-Pack security authentication in Basic Edition 6.8 and above. Valid values are `1` and `2`, `1` is disabled, `2` is enabled, and default value is `1`.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "A mapping of tags to assign to the instance. For tag limits, please refer to [Use Limits](https://intl.cloud.tencent.com/document/product/651/13354).",
			},

			// computed
			"elasticsearch_domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Elasticsearch domain name.",
			},
			"elasticsearch_vip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Elasticsearch VIP.",
			},
			"elasticsearch_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Elasticsearch port.",
			},
			"kibana_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Kibana access URL.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance creation time.",
			},
		},
	}
}

func resourceTencentCloudElasticsearchInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_elasticsearch_instance.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	elasticsearchService := ElasticsearchService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	request := es.NewCreateInstanceRequest()
	request.Zone = helper.String(d.Get("availability_zone").(string))
	request.EsVersion = helper.String(d.Get("version").(string))
	request.VpcId = helper.String(d.Get("vpc_id").(string))
	request.SubnetId = helper.String(d.Get("subnet_id").(string))
	request.Password = helper.String(d.Get("password").(string))
	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}
	if v, ok := d.GetOk("charge_type"); ok {
		chargeType := v.(string)
		request.ChargeType = &chargeType
		if chargeType == ES_CHARGE_TYPE_PREPAID {
			if v, ok := d.GetOk("charge_period"); ok {
				request.ChargePeriod = helper.IntUint64(v.(int))
			} else {
				return fmt.Errorf("elasticsearch charge period can not be empty when charge type is %s", chargeType)
			}
			if v, ok := d.GetOk("renew_flag"); ok {
				request.RenewFlag = helper.String(v.(string))
			}
		}
	}
	if v, ok := d.GetOk("deploy_mode"); ok {
		deployMode := v.(int)
		request.DeployMode = helper.IntUint64(deployMode)
		if deployMode == ES_DEPLOY_MODE_MULTI_REGION {
			if v, ok := d.GetOk("multi_zone_infos"); ok {
				infos := v.([]interface{})
				request.MultiZoneInfo = make([]*es.ZoneDetail, 0, len(infos))
				for _, item := range infos {
					value := item.(map[string]interface{})
					info := es.ZoneDetail{
						Zone:     helper.String(value["availability_zone"].(string)),
						SubnetId: helper.String(value["subnet_id"].(string)),
					}
					request.MultiZoneInfo = append(request.MultiZoneInfo, &info)
				}
			} else {
				return fmt.Errorf("elasticsearch multi_zone_infos can not be empty when deploy mode is %d", deployMode)
			}
		}
	}
	if v, ok := d.GetOk("license_type"); ok {
		request.LicenseType = helper.String(v.(string))
	}
	if v, ok := d.GetOk("basic_security_type"); ok {
		request.BasicSecurityType = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOk("node_info_list"); ok {
		infos := v.([]interface{})
		request.NodeInfoList = make([]*es.NodeInfo, 0, len(infos))
		for _, item := range infos {
			value := item.(map[string]interface{})
			info := es.NodeInfo{
				NodeNum:  helper.IntUint64(value["node_num"].(int)),
				NodeType: helper.String(value["node_type"].(string)),
			}
			if v := value["type"].(string); v != "" {
				info.Type = &v
			}
			if v := value["disk_type"].(string); v != "" {
				info.DiskType = &v
			}
			if v := value["disk_size"].(int); v > 0 {
				info.DiskSize = helper.IntUint64(v)
			}
			if v := value["encrypt"].(bool); v {
				info.DiskEncrypt = helper.BoolToInt64Pointer(v)
			}
			request.NodeInfoList = append(request.NodeInfoList, &info)
		}
	}

	instanceId := ""
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := meta.(*TencentCloudClient).apiV3Conn.UseEsClient().CreateInstance(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return retryError(err)
		}
		instanceId = *response.Response.InstanceId
		return nil
	})
	if err != nil {
		return err
	}
	d.SetId(instanceId)

	err = resource.Retry(5*readRetryTimeout, func() *resource.RetryError {
		instance, errRet := elasticsearchService.DescribeInstanceById(ctx, instanceId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if instance == nil || *instance.Status == ES_INSTANCE_STATUS_PROCESSING {
			return resource.RetryableError(errors.New("elasticsearch instance status is processing, retry..."))
		}
		return nil
	})
	if err != nil {
		return err
	}

	// tags
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		client := meta.(*TencentCloudClient).apiV3Conn
		tagService := TagService{client: client}
		region := client.Region
		resourceName := fmt.Sprintf("qcs::es:%s:uin/:instance/%s", region, instanceId)
		err := tagService.ModifyTags(ctx, resourceName, tags, nil)
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudElasticsearchInstanceRead(d, meta)
}

func resourceTencentCloudElasticsearchInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_elasticsearch_instance.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	instanceId := d.Id()
	elasticsearchService := ElasticsearchService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	var instance *es.InstanceInfo
	var errRet error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		instance, errRet = elasticsearchService.DescribeInstanceById(ctx, instanceId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if instance == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("instance_name", instance.InstanceName)
	_ = d.Set("availability_zone", instance.Zone)
	_ = d.Set("version", instance.EsVersion)
	_ = d.Set("vpc_id", instance.VpcUid)
	_ = d.Set("subnet_id", instance.SubnetUid)
	_ = d.Set("charge_type", instance.ChargeType)
	_ = d.Set("charge_period", instance.ChargePeriod)
	_ = d.Set("renew_flag", instance.RenewFlag)
	_ = d.Set("deploy_mode", instance.DeployMode)
	_ = d.Set("license_type", instance.LicenseType)
	_ = d.Set("basic_security_type", instance.SecurityType)
	_ = d.Set("elasticsearch_domain", instance.EsDomain)
	_ = d.Set("elasticsearch_vip", instance.EsVip)
	_ = d.Set("elasticsearch_port", instance.EsPort)
	_ = d.Set("kibana_url", instance.KibanaUrl)
	_ = d.Set("create_time", instance.CreateTime)

	multiZoneInfos := make([]map[string]interface{}, 0, len(instance.MultiZoneInfo))
	for _, item := range instance.MultiZoneInfo {
		info := make(map[string]interface{}, 2)
		info["availability_zone"] = item.Zone
		info["subnet_id"] = item.SubnetId
		multiZoneInfos = append(multiZoneInfos, info)
	}
	_ = d.Set("multi_zone_infos", multiZoneInfos)

	nodeInfoList := make([]map[string]interface{}, 0, len(instance.NodeInfoList))
	for _, item := range instance.NodeInfoList {
		info := make(map[string]interface{}, 5)
		info["node_num"] = item.NodeNum
		info["node_type"] = item.NodeType
		info["type"] = item.Type
		info["disk_type"] = item.DiskType
		info["disk_size"] = item.DiskSize
		info["encrypt"] = *item.DiskEncrypt > 0
		nodeInfoList = append(nodeInfoList, info)
	}
	_ = d.Set("node_info_list", nodeInfoList)

	if len(instance.TagList) > 0 {
		tags := make(map[string]string)
		for _, tag := range instance.TagList {
			tags[*tag.TagKey] = *tag.TagValue
		}
		_ = d.Set("tags", tags)
	}

	return nil
}

func resourceTencentCloudElasticsearchInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_elasticsearch.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	instanceId := d.Id()
	elasticsearchService := ElasticsearchService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	d.Partial(true)

	if d.HasChange("instance_name") {
		instanceName := d.Get("instance_name").(string)
		// Update operation support at most one item at the same time
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			errRet := elasticsearchService.UpdateInstance(ctx, instanceId, instanceName, "", 0)
			if errRet != nil {
				return retryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}
		d.SetPartial("instance_name")
	}
	if d.HasChange("password") {
		password := d.Get("password").(string)
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			errRet := elasticsearchService.UpdateInstance(ctx, instanceId, "", password, 0)
			if errRet != nil {
				return retryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}
		d.SetPartial("password")
	}

	if d.HasChange("version") {
		version := d.Get("version").(string)
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			errRet := elasticsearchService.UpdateInstanceVersion(ctx, instanceId, version)
			if errRet != nil {
				return retryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}
		d.SetPartial("version")

		err = resource.Retry(10*readRetryTimeout, func() *resource.RetryError {
			instance, errRet := elasticsearchService.DescribeInstanceById(ctx, instanceId)
			if errRet != nil {
				return retryError(errRet, InternalError)
			}
			if instance != nil && *instance.Status == ES_INSTANCE_STATUS_PROCESSING {
				return resource.RetryableError(errors.New("elasticsearch instance status is processing, retry..."))
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	if d.HasChange("license_type") {
		licenseType := d.Get("license_type").(string)
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			errRet := elasticsearchService.UpdateInstanceLicense(ctx, instanceId, licenseType)
			if errRet != nil {
				return retryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}
		d.SetPartial("license_type")

		err = resource.Retry(10*readRetryTimeout, func() *resource.RetryError {
			instance, errRet := elasticsearchService.DescribeInstanceById(ctx, instanceId)
			if errRet != nil {
				return retryError(errRet, InternalError)
			}
			if instance != nil && *instance.Status == ES_INSTANCE_STATUS_PROCESSING {
				return resource.RetryableError(errors.New("elasticsearch instance status is processing, retry..."))
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	if d.HasChange("basic_security_type") {
		basicSecurityType := d.Get("basic_security_type").(int)
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			errRet := elasticsearchService.UpdateInstance(ctx, instanceId, "", "", int64(basicSecurityType))
			if errRet != nil {
				return retryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}
		d.SetPartial("basic_security_type")
	}

	if d.HasChange("tags") {
		oldInterface, newInterface := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldInterface.(map[string]interface{}), newInterface.(map[string]interface{}))
		tagService := TagService{
			client: meta.(*TencentCloudClient).apiV3Conn,
		}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::es:%s:uin/:instance/%s", region, instanceId)
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
		d.SetPartial("tags")
	}

	d.Partial(false)

	return nil
}

func resourceTencentCloudElasticsearchInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_elasticsearch.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	instanceId := d.Id()
	elasticsearchService := ElasticsearchService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		errRet := elasticsearchService.DeleteInstance(ctx, instanceId)
		if errRet != nil {
			return retryError(errRet)
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
		instance, errRet := elasticsearchService.DescribeInstanceById(ctx, instanceId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if instance == nil {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("elasticsearch instance status is %d, retry...", *instance.Status))
	})
	if err != nil {
		return err
	}

	return nil
}
