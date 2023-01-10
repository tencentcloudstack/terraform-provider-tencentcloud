/*
Provides a resource to create a cynosdb security_group

Example Usage

```hcl
resource "tencentcloud_cynosdb_security_group" "test" {
  cluster_id = "cynosdbmysql-bws8h88b"
  security_group_ids = ["sg-baxfiao5"]
  instance_group_type = "RO"
}
```

Import

cynosdb security_group can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_security_group.security_group ${cluster_id}#${instance_group_type}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCynosdbSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbSecurityGroupCreate,
		Update: resourceTencentCloudCynosdbSecurityGroupUpdate,
		Read:   resourceTencentCloudCynosdbSecurityGroupRead,
		Delete: resourceTencentCloudCynosdbSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Cluster id.",
			},
			"instance_group_type": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance group type. Available values: \n-`HA` - HA group; \n-`RO` - Read-only group;\n-`ALL` - HA and RO group.",
			},
			"security_group_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of security group IDs to be modified, an array of one or more security group IDs.",
			},
		},
	}
}

func resourceTencentCloudCynosdbSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_security_group.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := cynosdb.NewModifyDBInstanceSecurityGroupsRequest()
	var clusterId string
	var instanceGroupType string

	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.InstanceId = helper.String(clusterId)
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroupIdsSet := v.(*schema.Set).List()
		for i := range securityGroupIdsSet {
			securityGroupId := securityGroupIdsSet[i].(string)
			request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroupId)
		}
	}

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	_, clusterIfo, has, err := service.DescribeClusterById(ctx, clusterId)
	if err != nil {
		return err
	}
	if has {
		request.Zone = clusterIfo.Zone
	}

	if v, ok := d.GetOk("instance_group_type"); ok {
		instanceGroupType = v.(string)
	}
	grpsResponse, err := service.DescribeClusterInstanceGrps(ctx, clusterId)
	if err != nil {
		return err
	}
	instanceGrpInfoList := grpsResponse.Response.InstanceGrpInfoList
	for _, instanceGrpInfo := range instanceGrpInfoList {
		posType := instanceGrpInfo.Type
		log.Printf("*posType: %v, %v", *posType, *posType != strings.ToLower(instanceGroupType))
		if *posType != strings.ToLower(instanceGroupType) && instanceGroupType != "ALL" {
			continue
		}
		request.InstanceId = instanceGrpInfo.InstanceGrpId
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().ModifyDBInstanceSecurityGroups(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s create cynosdb securityGroup failed, reason:%+v", logId, err)
			return err
		}
	}

	d.SetId(clusterId + FILED_SP + instanceGroupType)
	return resourceTencentCloudCynosdbSecurityGroupRead(d, meta)
}

func resourceTencentCloudCynosdbSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_security_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	instanceGroupType := idSplit[1]

	grpsResponse, err := service.DescribeClusterInstanceGrps(ctx, clusterId)
	if err != nil {
		return err
	}
	instanceGrpInfoList := grpsResponse.Response.InstanceGrpInfoList
	if len(instanceGrpInfoList) == 0 {
		return fmt.Errorf("Not fount instanceGrpInfoList")
	}

	var securityGroups []*cynosdb.SecurityGroup
	securityGroupIds := make([]string, 0)

	for _, instanceGrpInfo := range instanceGrpInfoList {
		if *instanceGrpInfo.Type != strings.ToLower(instanceGroupType) {
			continue
		}
		securityGroups, err = service.DescribeCynosdbSecurityGroups(ctx, *instanceGrpInfo.InstanceGrpId)
		if err != nil {
			return err
		}
	}

	for _, securityGroup := range securityGroups {
		securityGroupIds = append(securityGroupIds, *securityGroup.SecurityGroupId)
	}
	_ = d.Set("cluster_id", clusterId)
	_ = d.Set("instance_group_type", instanceGroupType)
	_ = d.Set("security_group_ids", securityGroupIds)

	return nil
}

func resourceTencentCloudCynosdbSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_security_group.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	instanceGroupType := idSplit[1]

	request := cynosdb.NewModifyDBInstanceSecurityGroupsRequest()
	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	grpsResponse, err := service.DescribeClusterInstanceGrps(ctx, clusterId)
	if err != nil {
		return err
	}

	d.Partial(true)
	if d.HasChange("security_group_ids") {
		securityGroupIdsSet := d.Get("security_group_ids").(*schema.Set).List()
		if len(securityGroupIdsSet) == 0 {
			return fmt.Errorf("`security_group_ids` can not empty")
		}
		for i := range securityGroupIdsSet {
			securityGroupId := securityGroupIdsSet[i].(string)
			request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroupId)
		}
	}
	_, clusterIfo, has, err := service.DescribeClusterById(ctx, clusterId)
	if err != nil {
		return err
	}
	if has {
		request.Zone = clusterIfo.Zone
	}
	instanceGrpInfoList := grpsResponse.Response.InstanceGrpInfoList
	for _, instanceGrpInfo := range instanceGrpInfoList {
		posType := instanceGrpInfo.Type
		log.Printf("*posType: %v, %v", *posType, *posType != strings.ToLower(instanceGroupType))
		if *posType != strings.ToLower(instanceGroupType) && instanceGroupType != "ALL" {
			continue
		}
		request.InstanceId = instanceGrpInfo.InstanceGrpId
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().ModifyDBInstanceSecurityGroups(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s create cynosdb securityGroup failed, reason:%+v", logId, err)
			return err
		}
	}
	d.SetPartial("security_group_ids")
	d.Partial(false)
	return resourceTencentCloudCynosdbSecurityGroupRead(d, meta)
}

func resourceTencentCloudCynosdbSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_security_group.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	instanceGroupType := idSplit[1]

	request := cynosdb.NewDisassociateSecurityGroupsRequest()
	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	_, clusterInfo, has, err := service.DescribeClusterById(ctx, clusterId)
	if err != nil {
		return err
	}
	if has {
		request.Zone = clusterInfo.Zone
	}
	grpsResponse, err := service.DescribeClusterInstanceGrps(ctx, clusterId)
	if err != nil {
		return err
	}

	instanceGrpInfoList := grpsResponse.Response.InstanceGrpInfoList
	for _, instanceGrpInfo := range instanceGrpInfoList {
		posType := instanceGrpInfo.Type
		if *posType != strings.ToLower(instanceGroupType) && instanceGroupType != "ALL" {
			continue
		}
		securityGroupIds := make([]*string, 0)
		securityGroups, err := service.DescribeCynosdbSecurityGroups(ctx, *instanceGrpInfo.InstanceGrpId)
		if err != nil {
			return err
		}

		for _, securityGroup := range securityGroups {
			securityGroupIds = append(securityGroupIds, securityGroup.SecurityGroupId)
		}

		request.InstanceIds = []*string{instanceGrpInfo.InstanceGrpId}
		request.SecurityGroupIds = securityGroupIds
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().DisassociateSecurityGroups(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s create cynosdb securityGroup failed, reason:%+v", logId, err)
			return err
		}
	}

	return nil
}
