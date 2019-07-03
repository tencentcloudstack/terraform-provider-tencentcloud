package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/zqfan/tencentcloud-sdk-go/client"
)

func resourceTencentCloudSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSecurityGroupCreate,
		Read:   resourceTencentCloudSecurityGroupRead,
		Update: resourceTencentCloudSecurityGroupUpdate,
		Delete: resourceTencentCloudSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(2, 60),
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(2, 100),
			},

			"project_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"attach_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: func(v interface{}, s string) (i []string, errs []error) {
						value := v.(string)
						for _, prefix := range []string{"ins-", "eni-", "cdb-", "lb-"} {
							if strings.HasPrefix(value, prefix) {
								return
							}
						}

						errs = append(errs, errors.New("id is invalid"))
						return
					},
				},
			},
		},
	}
}

func resourceTencentCloudSecurityGroupCreate(d *schema.ResourceData, m interface{}) error {
	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_security_group.create")()

	ctx := context.WithValue(context.TODO(), "logId", logId)

	client := m.(*TencentCloudClient)
	vpcService := VpcService{client: client.apiV3Conn}

	nameStr := d.Get("name").(string)
	name := &nameStr

	var (
		desc      *string
		projectID *string
	)

	if descInterface, exist := d.GetOk("description"); exist {
		descStr := descInterface.(string)
		desc = &descStr
	}

	if projectIDInterface, exist := d.GetOk("project_id"); exist {
		projectIDStr := projectIDInterface.(string)
		projectID = &projectIDStr
	}

	id, err := vpcService.CreateSecurityGroup(ctx, name, desc, projectID)
	if err != nil {
		return err
	}

	// attach instances to this security group
	if idsInterface, exist := d.GetOk("attach_ids"); exist {
		ids := expandStringList(idsInterface.(*schema.Set).List())

		if err := attachIds(ctx, ids, id, client.apiV3Conn, client.commonConn); err != nil {
			return err
		}
	}

	log.Printf("[DEBUG] SgId=%s", id)
	d.SetId(id)
	return resourceTencentCloudSecurityGroupRead(d, m)
}

func resourceTencentCloudSecurityGroupRead(d *schema.ResourceData, m interface{}) error {
	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_security_group.read")()

	ctx := context.WithValue(context.TODO(), "logId", logId)

	client := m.(*TencentCloudClient).apiV3Conn

	vpcService := VpcService{client: client}

	id := d.Id()

	securityGroup, has, err := vpcService.DescribeSecurityGroup(ctx, id)
	if err != nil {
		return err
	}

	switch has {
	default:
		err := fmt.Errorf("one security_group_id read get %d security_group info", has)
		log.Printf("[CRITAL]%s %v", logId, err)

		return err

	case 0:
		d.SetId("")
		return nil

	case 1:
		_ = d.Set("name", *securityGroup.SecurityGroupName)
		_ = d.Set("description", *securityGroup.SecurityGroupDesc)
		if securityGroup.ProjectId != nil {
			projectID, _ := strconv.Atoi(*securityGroup.ProjectId)
			_ = d.Set("project_id", projectID)
		}
	}

	var attachIds []string

	// find cvm which attach this security group
	cvmDescRequest := cvm.NewDescribeInstancesRequest()
	cvmDescRequest.Limit = common.Int64Ptr(100)
	cvmDescRequest.Filters = []*cvm.Filter{
		{
			Name:   common.StringPtr("security-group-id"),
			Values: common.StringPtrs([]string{id}),
		},
	}

	cvmDescResponse, err := client.UseCvmClient().DescribeInstances(cvmDescRequest)
	if err != nil {
		log.Printf("[CRITAL]%s %v", logId, err)
		return err
	}

	for _, cvmInstance := range cvmDescResponse.Response.InstanceSet {
		attachIds = append(attachIds, *cvmInstance.InstanceId)
	}

	// find ENI which attach this security group
	eniDescRequest := vpc.NewDescribeNetworkInterfacesRequest()
	eniDescRequest.Limit = common.Uint64Ptr(100)
	eniDescRequest.Filters = []*vpc.Filter{
		{
			Name:   common.StringPtr("groups.security-group-id"),
			Values: common.StringPtrs([]string{id}),
		},
	}

	eniDescResponse, err := client.UseVpcClient().DescribeNetworkInterfaces(eniDescRequest)
	if err != nil {
		log.Printf("[CRITAL]%s %v", logId, err)
		return err
	}

	for _, eniInstance := range eniDescResponse.Response.NetworkInterfaceSet {
		attachIds = append(attachIds, *eniInstance.NetworkInterfaceId)
	}

	// find mysql which attach this security group
	mysqlDescRequest := cdb.NewDescribeDBInstancesRequest()
	mysqlDescRequest.Limit = common.Uint64Ptr(2000)
	mysqlDescRequest.SecurityGroupId = &id

	mysqlDescResponse, err := client.UseMysqlClient().DescribeDBInstances(mysqlDescRequest)
	if err != nil {
		log.Printf("[CRITAL]%s %v", logId, err)
		return err
	}

	for _, mysqlInstance := range mysqlDescResponse.Response.Items {
		attachIds = append(attachIds, *mysqlInstance.InstanceId)
	}

	// find clb which attach this security group
	clbService := ClbService{client: client}
	clbs, err := clbService.DescribeLoadBalancers(ctx, nil, &id)
	if err != nil {
		log.Printf("[CRITAL]%s %v", logId, err)
		return err
	}

	for _, clbInstance := range clbs {
		attachIds = append(attachIds, *clbInstance.LoadBalancerId)
	}

	if len(attachIds) > 0 {
		_ = d.Set("attach_ids", attachIds)
	}

	return nil
}

func resourceTencentCloudSecurityGroupUpdate(d *schema.ResourceData, m interface{}) error {
	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_security_group.update")()

	ctx := context.WithValue(context.TODO(), "logId", logId)

	commonClient := m.(*TencentCloudClient).commonConn
	client := m.(*TencentCloudClient).apiV3Conn

	vpcService := VpcService{client: client}

	id := d.Id()

	attributeUpdate := d.HasChange("name") || d.HasChange("description") || d.HasChange("attach_ids")
	var (
		newName         *string
		newDesc         *string
		removeAttachIds []string
		newAttachIds    []string
	)

	d.Partial(true)
	defer d.Partial(false)

	if d.HasChange("name") {
		d.SetPartial("name")
		newName = common.StringPtr(d.Get("name").(string))
	}

	if d.HasChange("description") {
		d.SetPartial("description")
		newDesc = common.StringPtr(d.Get("description").(string))
	}

	if d.HasChange("attach_ids") {
		d.SetPartial("attach_ids")
		o, n := d.GetChange("attach_ids")
		oldSet := o.(*schema.Set)
		newSet := n.(*schema.Set)

		removeAttachIds = expandStringList(oldSet.Difference(newSet).List())
		newAttachIds = expandStringList(newSet.Difference(oldSet).List())
	}

	if !attributeUpdate {
		return nil
	}

	// update security group itself
	if err := vpcService.ModifySecurityGroup(ctx, id, newName, newDesc); err != nil {
		return err
	}

	// update attach ids
	if len(removeAttachIds) > 0 {
		if err := unattachIds(ctx, removeAttachIds, id, client, commonClient); err != nil {
			return err
		}
	}

	if len(newAttachIds) > 0 {
		if err := attachIds(ctx, newAttachIds, id, client, commonClient); err != nil {
			return err
		}
	}

	return resourceTencentCloudSecurityGroupRead(d, m)
}

func attachIds(
	ctx context.Context,
	ids []string,
	securityGroupId string,
	client *connectivity.TencentCloudClient,
	commonClient *client.Client,
) error {
	logId := GetLogId(ctx)

	cvmIns, eniIns, mysqlIns, clbIns := splitAttachIds(ids)

	cvmClient := client.UseCvmClient()
	for _, cvmId := range cvmIns {
		descRequest := cvm.NewDescribeInstancesRequest()
		descRequest.InstanceIds = common.StringPtrs([]string{cvmId})

		descResponse, err := cvmClient.DescribeInstances(descRequest)
		if err != nil {
			log.Printf("[CRITAL]%s get cvm by id %s error, reason: %v", logId, cvmId, err)
			return err
		}

		if *descResponse.Response.TotalCount != 1 {
			err := fmt.Errorf("cvm id %s does not exist", cvmId)
			log.Printf("[CRITAL]%s %v", logId, err)
			return err
		}

		cvmInstance := descResponse.Response.InstanceSet[0]
		existSgIds := make([]string, 0, len(cvmInstance.SecurityGroupIds))
		for _, sgId := range cvmInstance.SecurityGroupIds {
			existSgIds = append(existSgIds, *sgId)
		}

		existSgIds = addSecurityGroupToExistSecurityGroups(existSgIds, securityGroupId)

		if err := bindInstanceWithSgIds(commonClient, cvmId, existSgIds); err != nil {
			log.Printf("[CRITAL]%s attach cvm %s to security group %s error, reason: %v",
				logId, cvmId, logId, err)
			return err
		}
	}

	vpcService := VpcService{client: client}
	for _, eniId := range eniIns {
		enis, err := vpcService.DescribeNetworkInterfaces(ctx, []string{eniId})
		if err != nil {
			log.Printf("[CRITAL]%s get eni by id %s error, reason: %v", logId, eniId, err)
			return err
		}

		if len(enis) != 1 {
			err := fmt.Errorf("eni id %s does not exist", eniId)
			log.Printf("[CRITAL]%s %v", logId, err)
			return err
		}

		eniInstance := enis[0]
		existSgIds := make([]string, 0, len(eniInstance.GroupSet))
		for _, sgId := range eniInstance.GroupSet {
			existSgIds = append(existSgIds, *sgId)
		}

		existSgIds = addSecurityGroupToExistSecurityGroups(existSgIds, securityGroupId)

		if err := vpcService.AttachEniToSecurityGroup(ctx, eniId, existSgIds); err != nil {
			log.Printf("[CRITAL]%s attach secondary ENI %s to security group %s error, reason: %v",
				logId, eniId, logId, err)
			return err
		}
	}

	mysqlService := MysqlService{client: client}
	for _, mysqlId := range mysqlIns {
		existSgIds, err := mysqlService.DescribeDBSecurityGroups(ctx, mysqlId)
		if err != nil {
			log.Printf("[CRITAL]%s get mysql security groups by id %s error, reason: %v", logId, mysqlId, err)
			return err
		}

		existSgIds = addSecurityGroupToExistSecurityGroups(existSgIds, securityGroupId)

		if err := mysqlService.ModifyDBInstanceSecurityGroups(ctx, mysqlId, existSgIds); err != nil {
			log.Printf("[CRITAL]%s attach mysql %s to security group %s error, reason: %v",
				logId, mysqlId, logId, err)
			return err
		}
	}

	clbService := ClbService{client: client}
	for _, clbId := range clbIns {
		clb, err := clbService.DescribeLoadBalancers(ctx, []string{clbId}, nil)
		if err != nil {
			log.Printf("[CRITAL]%s get clb by id %s error, reason: %v", logId, clbId, err)
			return err
		}

		if len(clb) != 1 {
			err := fmt.Errorf("clb id %s does not exist", clbId)
			log.Printf("[CRITAL]%s %v", logId, err)
			return err
		}

		clbInstance := clb[0]
		existSgIds := make([]string, 0, len(clbInstance.SecureGroups))
		for _, sgId := range clbInstance.SecureGroups {
			existSgIds = append(existSgIds, *sgId)
		}

		existSgIds = addSecurityGroupToExistSecurityGroups(existSgIds, securityGroupId)

		if err := clbService.ModifyLoadBalancerSecurityGroups(ctx, clbId, existSgIds); err != nil {
			log.Printf("[CRITAL]%s attach clb %s to security group %s error, reason: %v",
				logId, clbId, logId, err)
			return err
		}
	}

	return nil
}

func unattachIds(
	ctx context.Context,
	ids []string,
	securityGroupId string,
	client *connectivity.TencentCloudClient,
	commonClient *client.Client,
) error {
	logId := GetLogId(ctx)

	cvmIns, eniIns, mysqlIns, clbIns := splitAttachIds(ids)

	// unattach cvm
	for _, cvmId := range cvmIns {
		cvmDescRequest := cvm.NewDescribeInstancesRequest()
		cvmDescRequest.InstanceIds = common.StringPtrs([]string{cvmId})
		cvmDescRequest.Limit = common.Int64Ptr(1)

		cvmDescResponse, err := client.UseCvmClient().DescribeInstances(cvmDescRequest)
		if err != nil {
			log.Printf("[CRITAL]%s %v", logId, err)
			return err
		}

		if len(cvmDescResponse.Response.InstanceSet) != 1 {
			err := fmt.Errorf("give cvm id %s does not exist", cvmId)
			log.Printf("[CRITAL]%s %v", logId, err)
			return err
		}

		cvmInstance := cvmDescResponse.Response.InstanceSet[0]
		existSgIds := make([]string, 0, len(cvmInstance.SecurityGroupIds))
		for _, sgId := range cvmInstance.SecurityGroupIds {
			if *sgId == securityGroupId {
				continue
			}
			existSgIds = append(existSgIds, *sgId)
		}

		if err := bindInstanceWithSgIds(commonClient, cvmId, existSgIds); err != nil {
			log.Printf("[CRITAL]%s unattach cvm id %s security group %s error, reason: %v",
				logId, cvmId, logId, err)
			return err
		}
	}

	// unattach eni
	vpcService := VpcService{client: client}
	for _, eniId := range eniIns {
		enis, err := vpcService.DescribeNetworkInterfaces(ctx, []string{eniId})
		if err != nil {
			log.Printf("[CRITAL]%s get eni by id %s error, reason: %v", logId, eniId, err)
			return err
		}

		if len(enis) != 1 {
			err := fmt.Errorf("eni id %s does not exist", eniId)
			log.Printf("[CRITAL]%s %v", logId, err)
			return err
		}

		eniInstance := enis[0]
		existSgIds := make([]string, 0, len(eniInstance.GroupSet))
		for _, sgId := range eniInstance.GroupSet {
			if *sgId == securityGroupId {
				continue
			}

			existSgIds = append(existSgIds, *sgId)
		}

		if err := vpcService.AttachEniToSecurityGroup(ctx, eniId, existSgIds); err != nil {
			log.Printf("[CRITAL]%s unattach ENI %s security group %s error, reason: %v",
				logId, eniId, securityGroupId, err)
			return err
		}
	}

	// unattach mysql
	mysqlService := MysqlService{client: client}
	for _, mysqlId := range mysqlIns {
		existSgIds, err := mysqlService.DescribeDBSecurityGroups(ctx, mysqlId)
		if err != nil {
			log.Printf("[CRITAL]%s get mysql security groups by id %s error, reason: %v", logId, mysqlId, err)
			return err
		}

		for i := range existSgIds {
			if existSgIds[i] == securityGroupId {
				existSgIds = append(existSgIds[:i], existSgIds[i+1:]...)
				break
			}
		}

		if err := mysqlService.ModifyDBInstanceSecurityGroups(ctx, mysqlId, existSgIds); err != nil {
			log.Printf("[CRITAL]%s unattach mysql id %s security group %s error, reason: %v",
				logId, mysqlId, securityGroupId, err)
			return err
		}
	}

	// unattach clb
	clbService := ClbService{client: client}
	for _, clbId := range clbIns {
		clbs, err := clbService.DescribeLoadBalancers(ctx, []string{clbId}, nil)
		if err != nil {
			log.Printf("[CRITAL]%s get clb security groups by id %s error, reason: %v", logId, clbId, err)
			return err
		}

		if len(clbs) != 1 {
			err := fmt.Errorf("get clb by id %s but return %d clb instances", clbId, len(clbs))
			log.Printf("[CRITAL]%s %v", logId, err)
			return err
		}

		clbInstance := clbs[0]
		existSgIds := make([]string, 0, len(clbInstance.SecureGroups))
		for _, sgId := range clbInstance.SecureGroups {
			if *sgId == securityGroupId {
				continue
			}

			existSgIds = append(existSgIds, *sgId)
		}

		if err := clbService.ModifyLoadBalancerSecurityGroups(ctx, clbId, existSgIds); err != nil {
			log.Printf("[CRITAL]%s unattach clb id %s security group %s error, reason: %v",
				logId, clbId, securityGroupId, err)
			return err
		}
	}

	return nil
}

func resourceTencentCloudSecurityGroupDelete(d *schema.ResourceData, m interface{}) error {
	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_security_group.delete")()

	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()

	client := m.(*TencentCloudClient).apiV3Conn
	commonClient := m.(*TencentCloudClient).commonConn

	if idsInterface, exist := d.GetOk("attach_ids"); exist {
		ids := expandStringList(idsInterface.(*schema.Set).List())

		if err := unattachIds(ctx, ids, id, client, commonClient); err != nil {
			return err
		}
	}

	vpcService := VpcService{client: client}

	return vpcService.DeleteSecurityGroup(ctx, id)
}

func splitAttachIds(ids []string) (cvmIns, eniIns, mysqlIns, clbIns []string) {
	for _, id := range ids {
		switch {
		case strings.HasPrefix(id, "ins-"):
			cvmIns = append(cvmIns, id)

		case strings.HasPrefix(id, "eni-"):
			eniIns = append(eniIns, id)

		case strings.HasPrefix(id, "cdb-"):
			mysqlIns = append(mysqlIns, id)

		case strings.HasPrefix(id, "lb-"):
			clbIns = append(clbIns, id)
		}
	}

	return
}

// addSecurityGroupToExistSecurityGroups if id exists, will do nothing
func addSecurityGroupToExistSecurityGroups(existSecurityGroups []string, securityGroupId string) []string {
	var alreadyHasSecurityGroup bool
	for _, sgId := range existSecurityGroups {
		if sgId == securityGroupId {
			alreadyHasSecurityGroup = true
			break
		}
	}

	if !alreadyHasSecurityGroup {
		existSecurityGroups = append(existSecurityGroups, securityGroupId)
	}

	return existSecurityGroups
}
