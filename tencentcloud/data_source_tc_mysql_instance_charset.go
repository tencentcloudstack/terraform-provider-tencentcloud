package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
)

func dataSourceTencentCloudMysqlInstanceCharset() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMysqlInstanceCharsetRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID, the format is: cdb-c1nl9rpv, which is the same as the instance ID displayed on the cloud database console page, and you can use the [query instance list] (https://cloud.tencent.com/document/api/236/15872) interface Gets the value of the field InstanceId in the output parameter.",
			},

			"charset": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The default character set of the instance, such as `latin1`, `utf8` etc.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudMysqlInstanceCharsetRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mysql_instance_charset.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	var instanceCharset *cdb.DescribeDBInstanceCharsetResponseParams
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMysqlInstanceCharsetByFilter(ctx, instanceId)
		if e != nil {
			return retryError(e)
		}
		instanceCharset = result
		return nil
	})
	if err != nil {
		return err
	}

	if instanceCharset != nil {
		_ = d.Set("charset", instanceCharset.Charset)
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), d); e != nil {
			return e
		}
	}
	return nil
}
