package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsWorkflowResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsWorkflow,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mps_workflow.workflow", "id")),
			},
			{
				ResourceName:      "tencentcloud_mps_workflow.workflow",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsWorkflow = `

resource "tencentcloud_mps_workflow" "workflow" {
  workflow_name = &lt;nil&gt;
  trigger {
		type = "CosFileUpload"
		cos_file_upload_trigger {
			bucket = "TopRankVideo-125xxx88"
			region = "ap-chongqing"
			dir = "/movie/201907/"
			formats = 
		}

  }
  output_storage {
		type = "COS"
		cos_output_storage {
			bucket = "TopRankVideo-125xxx88"
			region = "ap-chongqing"
		}

  }
  output_dir = "/movie/201907/"
  media_process_task {
		transcode_task_set {
			definition = &lt;nil&gt;
			raw_parameter {
				container = &lt;nil&gt;
				remove_video = 0
				remove_audio = 0
				video_template {
					codec = &lt;nil&gt;
					fps = &lt;nil&gt;
					bitrate = &lt;nil&gt;
					resolution_adaptive = "open"
					width = 0
					height = 0
					gop = &lt;nil&gt;
					fill_type = "black"
					vcrf = &lt;nil&gt;
				}
				audio_template {
					codec = &lt;nil&gt;
					bitrate = &lt;nil&gt;
					sample_rate = &lt;nil&gt;
					audio_channel = 2
				}
				t_e_h_d_config {
					type = &lt;nil&gt;
					max_video_bitrate = &lt;nil&gt;
				}
			}
			override_parameter {
				container = &lt;nil&gt;
				remove_video = &lt;nil&gt;
				remove_audio = &lt;nil&gt;
				video_template {
					codec = &lt;nil&gt;
					fps = &lt;nil&gt;
					bitrate = &lt;nil&gt;
					resolution_adaptive = &lt;nil&gt;
					width = &lt;nil&gt;
					height = &lt;nil&gt;
					gop = &lt;nil&gt;
					fill_type = &lt;nil&gt;
					vcrf = &lt;nil&gt;
					content_adapt_stream = 0
				}
				audio_template {
					codec = &lt;nil&gt;
					bitrate = &lt;nil&gt;
					sample_rate = &lt;nil&gt;
					audio_channel = &lt;nil&gt;
					stream_selects = &lt;nil&gt;
				}
				t_e_h_d_config {
					type = &lt;nil&gt;
					max_video_bitrate = &lt;nil&gt;
				}
				subtitle_template {
					path = &lt;nil&gt;
					stream_index = &lt;nil&gt;
					font_type = "hei.ttf"
					font_size = &lt;nil&gt;
					font_color = "0xFFFFFF"
					font_alpha = 
				}
			}
			watermark_set {
				definition = &lt;nil&gt;
				raw_parameter {
					type = &lt;nil&gt;
					coordinate_origin = "TopLeft"
					x_pos = "0px"
					y_pos = "0px"
					image_template {
						image_content {
							type = "COS"
							cos_input_info {
								bucket = "TopRankVideo-125xxx88"
								region = "ap-chongqing"
								object = "/movie/201907/WildAnimal.mov"
							}
							url_input_info {
								url = &lt;nil&gt;
							}
						}
						width = "10%"
						height = "0px"
						repeat_type = &lt;nil&gt;
					}
				}
				text_content = &lt;nil&gt;
				svg_content = &lt;nil&gt;
				start_time_offset = &lt;nil&gt;
				end_time_offset = &lt;nil&gt;
			}
			mosaic_set {
				coordinate_origin = "TopLeft"
				x_pos = "0px"
				y_pos = "0px"
				width = "10%"
				height = "10%"
				start_time_offset = &lt;nil&gt;
				end_time_offset = &lt;nil&gt;
			}
			start_time_offset = &lt;nil&gt;
			end_time_offset = &lt;nil&gt;
			output_storage {
				type = "COS"
				cos_output_storage {
					bucket = "TopRankVideo-125xxx88"
					region = "ap-chongqinq"
				}
			}
			output_object_path = &lt;nil&gt;
			segment_object_name = &lt;nil&gt;
			object_number_format {
				initial_value = 0
				increment = 1
				min_length = 1
				place_holder = "0"
			}
			head_tail_parameter {
				head_set {
					type = "COS"
					cos_input_info {
						bucket = "TopRankVideo-125xxx88"
						region = "ap-chongqing"
						object = "/movie/201907/WildAnimal.mov"
					}
					url_input_info {
						url = &lt;nil&gt;
					}
				}
				tail_set {
					type = "COS"
					cos_input_info {
						bucket = "TopRankVideo-125xxx88"
						region = "ap-chongqing"
						object = "/movie/201907/WildAnimal.mov"
					}
					url_input_info {
						url = &lt;nil&gt;
					}
				}
			}
		}
		animated_graphic_task_set {
			definition = &lt;nil&gt;
			start_time_offset = &lt;nil&gt;
			end_time_offset = &lt;nil&gt;
			output_storage {
				type = "COS"
				cos_output_storage {
					bucket = "TopRankVideo-125xxx88"
					region = "ap-chongqinq"
				}
			}
			output_object_path = &lt;nil&gt;
		}
		snapshot_by_time_offset_task_set {
			definition = &lt;nil&gt;
			ext_time_offset_set = &lt;nil&gt;
			time_offset_set = &lt;nil&gt;
			watermark_set {
				definition = &lt;nil&gt;
				raw_parameter {
					type = &lt;nil&gt;
					coordinate_origin = "TopLeft"
					x_pos = "0px"
					y_pos = "0px"
					image_template {
						image_content {
							type = "COS"
							cos_input_info {
								bucket = "TopRankVideo-125xxx88"
								region = "ap-chongqing"
								object = "/movie/201907/WildAnimal.mov"
							}
							url_input_info {
								url = &lt;nil&gt;
							}
						}
						width = "10%"
						height = "0px"
						repeat_type = &lt;nil&gt;
					}
				}
				text_content = &lt;nil&gt;
				svg_content = &lt;nil&gt;
				start_time_offset = &lt;nil&gt;
				end_time_offset = &lt;nil&gt;
			}
			output_storage {
				type = "COS"
				cos_output_storage {
					bucket = "TopRankVideo-125xxx88"
					region = "ap-chongqinq"
				}
			}
			output_object_path = &lt;nil&gt;
			object_number_format {
				initial_value = 0
				increment = 1
				min_length = 1
				place_holder = "0"
			}
		}
		sample_snapshot_task_set {
			definition = &lt;nil&gt;
			watermark_set {
				definition = &lt;nil&gt;
				raw_parameter {
					type = &lt;nil&gt;
					coordinate_origin = "TopLeft"
					x_pos = "0px"
					y_pos = "0px"
					image_template {
						image_content {
							type = "COS"
							cos_input_info {
								bucket = "TopRankVideo-125xxx88"
								region = "ap-chongqing"
								object = "/movie/201907/WildAnimal.mov"
							}
							url_input_info {
								url = &lt;nil&gt;
							}
						}
						width = "10%"
						height = "0px"
						repeat_type = "repeat"
					}
				}
				text_content = &lt;nil&gt;
				svg_content = &lt;nil&gt;
				start_time_offset = &lt;nil&gt;
				end_time_offset = &lt;nil&gt;
			}
			output_storage {
				type = "COS"
				cos_output_storage {
					bucket = "TopRankVideo-125xxx88"
					region = "ap-chongqinq"
				}
			}
			output_object_path = &lt;nil&gt;
			object_number_format {
				initial_value = 0
				increment = 1
				min_length = 1
				place_holder = "0"
			}
		}
		image_sprite_task_set {
			definition = &lt;nil&gt;
			output_storage {
				type = "COS"
				cos_output_storage {
					bucket = "TopRankVideo-125xxx88"
					region = "ap-chongqinq"
				}
			}
			output_object_path = &lt;nil&gt;
			web_vtt_object_name = &lt;nil&gt;
			object_number_format {
				initial_value = 0
				increment = 1
				min_length = 1
				place_holder = "0"
			}
		}
		adaptive_dynamic_streaming_task_set {
			definition = &lt;nil&gt;
			watermark_set {
				definition = &lt;nil&gt;
				raw_parameter {
					type = &lt;nil&gt;
					coordinate_origin = "TopLeft"
					x_pos = "0px"
					y_pos = "0px"
					image_template {
						image_content {
							type = "COS"
							cos_input_info {
								bucket = "TopRankVideo-125xxx88"
								region = "ap-chongqing"
								object = "/movie/201907/WildAnimal.mov"
							}
							url_input_info {
								url = &lt;nil&gt;
							}
						}
						width = "10%"
						height = "0px"
						repeat_type = "repeat"
					}
				}
				text_content = &lt;nil&gt;
				svg_content = &lt;nil&gt;
				start_time_offset = &lt;nil&gt;
				end_time_offset = &lt;nil&gt;
			}
			output_storage {
				type = "COS"
				cos_output_storage {
					bucket = "TopRankVideo-125xxx88"
					region = "ap-chongqinq"
				}
			}
			output_object_path = &lt;nil&gt;
			sub_stream_object_name = &lt;nil&gt;
			segment_object_name = &lt;nil&gt;
		}

  }
  ai_content_review_task {
		definition = &lt;nil&gt;

  }
  ai_analysis_task {
		definition = &lt;nil&gt;
		extended_parameter = &lt;nil&gt;

  }
  ai_recognition_task {
		definition = &lt;nil&gt;

  }
  task_notify_config {
		cmq_model = &lt;nil&gt;
		cmq_region = &lt;nil&gt;
		topic_name = &lt;nil&gt;
		queue_name = &lt;nil&gt;
		notify_mode = &lt;nil&gt;
		notify_type = &lt;nil&gt;
		notify_url = &lt;nil&gt;

  }
  task_priority = &lt;nil&gt;
}

`
