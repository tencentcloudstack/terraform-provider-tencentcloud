/*
Provides a resource to create a postgres network_access

Example Usage

```hcl
resource "tencentcloud_postgres_network_access" "network_access" {
  read_only_group_id = "pgro-xxxx"
  vpc_id = "vpc-xxx"
  subnet_id = "subnet-xxx"
  is_assign_vip = false
  vip = ""
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

postgres network_access can be imported using the id, e.g.

```
terraform import tencentcloud_postgres_network_access.network_access network_access_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudPostgresNetworkAccess() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresNetworkAccessCreate,
		Read:   resourceTencentCloudPostgresNetworkAccessRead,
		Delete: resourceTencentCloudPostgresNetworkAccessDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"read_only_group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "RO group identifier.",
			},

			"vpc_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Unified VPC ID.",
			},

			"subnet_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Subnet ID.",
			},

			"is_assign_vip": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to manually assign the VIP. Valid values:true (manually assign), false (automatically assign).",
			},

			"vip": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Target VIP.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudPostgresNetworkAccessCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_network_access.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request         = postgres.NewCreateReadOnlyGroupNetworkAccessRequest()
		response        = postgres.NewCreateReadOnlyGroupNetworkAccessResponse()
		readOnlyGroupId string
	)
	if v, ok := d.GetOk("read_only_group_id"); ok {
		readOnlyGroupId = v.(string)
		request.ReadOnlyGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("is_assign_vip"); ok {
		request.IsAssignVip = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("vip"); ok {
		request.Vip = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresClient().CreateReadOnlyGroupNetworkAccess(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create postgres NetworkAccess failed, reason:%+v", logId, err)
		return err
	}

	readOnlyGroupId = *response.Response.ReadOnlyGroupId
	d.SetId(readOnlyGroupId)

	service := PostgresService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"opened"}, 180*readRetryTimeout, time.Second, service.PostgresNetworkAccessStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::postgres:%s:uin/:dbInstanceId/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudPostgresNetworkAccessRead(d, meta)
}

func resourceTencentCloudPostgresNetworkAccessRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_network_access.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := PostgresService{client: meta.(*TencentCloudClient).apiV3Conn}

	networkAccessId := d.Id()

	NetworkAccess, err := service.DescribePostgresNetworkAccessById(ctx, readOnlyGroupId)
	if err != nil {
		return err
	}

	if NetworkAccess == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `PostgresNetworkAccess` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if NetworkAccess.ReadOnlyGroupId != nil {
		_ = d.Set("read_only_group_id", NetworkAccess.ReadOnlyGroupId)
	}

	if NetworkAccess.VpcId != nil {
		_ = d.Set("vpc_id", NetworkAccess.VpcId)
	}

	if NetworkAccess.SubnetId != nil {
		_ = d.Set("subnet_id", NetworkAccess.SubnetId)
	}

	if NetworkAccess.IsAssignVip != nil {
		_ = d.Set("is_assign_vip", NetworkAccess.IsAssignVip)
	}

	if NetworkAccess.Vip != nil {
		_ = d.Set("vip", NetworkAccess.Vip)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "postgres", "dbInstanceId", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudPostgresNetworkAccessDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_network_access.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := PostgresService{client: meta.(*TencentCloudClient).apiV3Conn}
	networkAccessId := d.Id()

	if err := service.DeletePostgresNetworkAccessById(ctx, readOnlyGroupId); err != nil {
		return err
	}

	service := PostgresService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"closed"}, 180*readRetryTimeout, time.Second, service.PostgresNetworkAccessStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
