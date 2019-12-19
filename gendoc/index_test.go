package main

import (
	"reflect"
	"testing"
)

const successDoc = `
AS
  Data Source
    tencentcloud_as_scaling_configs
    tencentcloud_as_scaling_groups


  Resource
    tencentcloud_as_scaling_config
    tencentcloud_as_scaling_group


CAM
  Data Source
    tencentcloud_cam_group_memberships
    tencentcloud_cam_group_policy_attachments


  Resource
    tencentcloud_cam_role
    tencentcloud_cam_role_policy_attachment

`

func TestGetIndex(t *testing.T) {
	type args struct {
		doc string
	}
	tests := []struct {
		name    string
		args    args
		want    []Product
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				doc: successDoc,
			},
			want: []Product{
				{
					Name: "AS",
					DataSources: []string{
						"tencentcloud_as_scaling_configs",
						"tencentcloud_as_scaling_groups",
					},
					Resources: []string{
						"tencentcloud_as_scaling_config",
						"tencentcloud_as_scaling_group",
					},
				},
				{
					Name: "CAM",
					DataSources: []string{
						"tencentcloud_cam_group_memberships",
						"tencentcloud_cam_group_policy_attachments",
					},
					Resources: []string{
						"tencentcloud_cam_role",
						"tencentcloud_cam_role_policy_attachment",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetIndex(tt.args.doc)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetIndex() got = %v, want %v", got, tt.want)
			}
		})
	}
}
