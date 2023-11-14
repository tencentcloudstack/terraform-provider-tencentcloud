package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsScheduleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsSchedule,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "id")),
			},
			{
				ResourceName:      "tencentcloud_mps_schedule.schedule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsSchedule = `

resource "tencentcloud_mps_schedule" "schedule" {
  schedule_name = ""
  trigger {
		type = ""
		cos_file_upload_trigger {
			bucket = ""
			region = ""
			dir = ""
			formats = 
		}
		aws_s3_file_upload_trigger {
			s3_bucket = ""
			s3_region = ""
			dir = ""
			formats = 
			s3_secret_id = ""
			s3_secret_key = ""
			aws_s_q_s {
				s_q_s_region = ""
				s_q_s_queue_name = ""
				s3_secret_id = ""
				s3_secret_key = ""
			}
		}

  }
  activities {
		activity_type = ""
		reardrive_index = 
		activity_para {
			transcode_task {
				definition = 
				raw_parameter {
					container = ""
					remove_video = 
					remove_audio = 
					video_template {
						codec = ""
						fps = 
						bitrate = 
						resolution_adaptive = ""
						width = 
						height = 
						gop = 
						fill_type = ""
						vcrf = 
					}
					audio_template {
						codec = ""
						bitrate = 
						sample_rate = 
						audio_channel = 
					}
					t_e_h_d_config {
						type = ""
						max_video_bitrate = 
					}
				}
				override_parameter {
					container = ""
					remove_video = 
					remove_audio = 
					video_template {
						codec = ""
						fps = 
						bitrate = 
						resolution_adaptive = ""
						width = 
						height = 
						gop = 
						fill_type = ""
						vcrf = 
						content_adapt_stream = 
					}
					audio_template {
						codec = ""
						bitrate = 
						sample_rate = 
						audio_channel = 
						stream_selects = 
					}
					t_e_h_d_config {
						type = ""
						max_video_bitrate = 
					}
					subtitle_template {
						path = ""
						stream_index = 
						font_type = ""
						font_size = ""
						font_color = ""
						font_alpha = 
					}
					addon_audio_stream {
						type = ""
						cos_input_info {
							bucket = ""
							region = ""
							object = ""
						}
						url_input_info {
							url = ""
						}
						s3_input_info {
							s3_bucket = ""
							s3_region = ""
							s3_object = ""
							s3_secret_id = ""
							s3_secret_key = ""
						}
					}
					std_ext_info = ""
					add_on_subtitles {
						type = ""
						subtitle {
							type = ""
							cos_input_info {
								bucket = ""
								region = ""
								object = ""
							}
							url_input_info {
								url = ""
							}
							s3_input_info {
								s3_bucket = ""
								s3_region = ""
								s3_object = ""
								s3_secret_id = ""
								s3_secret_key = ""
							}
						}
					}
				}
				watermark_set {
					definition = 
					raw_parameter {
						type = ""
						coordinate_origin = ""
						x_pos = ""
						y_pos = ""
						image_template {
							image_content {
								type = ""
								cos_input_info {
									bucket = ""
									region = ""
									object = ""
								}
								url_input_info {
									url = ""
								}
								s3_input_info {
									s3_bucket = ""
									s3_region = ""
									s3_object = ""
									s3_secret_id = ""
									s3_secret_key = ""
								}
							}
							width = ""
							height = ""
							repeat_type = ""
						}
					}
					text_content = ""
					svg_content = ""
					start_time_offset = 
					end_time_offset = 
				}
				mosaic_set {
					coordinate_origin = ""
					x_pos = ""
					y_pos = ""
					width = ""
					height = ""
					start_time_offset = 
					end_time_offset = 
				}
				start_time_offset = 
				end_time_offset = 
				output_storage {
					type = ""
					cos_output_storage {
						bucket = ""
						region = ""
					}
					s3_output_storage {
						s3_bucket = ""
						s3_region = ""
						s3_secret_id = ""
						s3_secret_key = ""
					}
				}
				output_object_path = ""
				segment_object_name = ""
				object_number_format {
					initial_value = 
					increment = 
					min_length = 
					place_holder = ""
				}
				head_tail_parameter {
					head_set {
						type = ""
						cos_input_info {
							bucket = ""
							region = ""
							object = ""
						}
						url_input_info {
							url = ""
						}
						s3_input_info {
							s3_bucket = ""
							s3_region = ""
							s3_object = ""
							s3_secret_id = ""
							s3_secret_key = ""
						}
					}
					tail_set {
						type = ""
						cos_input_info {
							bucket = ""
							region = ""
							object = ""
						}
						url_input_info {
							url = ""
						}
						s3_input_info {
							s3_bucket = ""
							s3_region = ""
							s3_object = ""
							s3_secret_id = ""
							s3_secret_key = ""
						}
					}
				}
			}
			animated_graphic_task {
				definition = 
				start_time_offset = 
				end_time_offset = 
				output_storage {
					type = ""
					cos_output_storage {
						bucket = ""
						region = ""
					}
					s3_output_storage {
						s3_bucket = ""
						s3_region = ""
						s3_secret_id = ""
						s3_secret_key = ""
					}
				}
				output_object_path = ""
			}
			snapshot_by_time_offset_task {
				definition = 
				ext_time_offset_set = 
				time_offset_set = 
				watermark_set {
					definition = 
					raw_parameter {
						type = ""
						coordinate_origin = ""
						x_pos = ""
						y_pos = ""
						image_template {
							image_content {
								type = ""
								cos_input_info {
									bucket = ""
									region = ""
									object = ""
								}
								url_input_info {
									url = ""
								}
								s3_input_info {
									s3_bucket = ""
									s3_region = ""
									s3_object = ""
									s3_secret_id = ""
									s3_secret_key = ""
								}
							}
							width = ""
							height = ""
							repeat_type = ""
						}
					}
					text_content = ""
					svg_content = ""
					start_time_offset = 
					end_time_offset = 
				}
				output_storage {
					type = ""
					cos_output_storage {
						bucket = ""
						region = ""
					}
					s3_output_storage {
						s3_bucket = ""
						s3_region = ""
						s3_secret_id = ""
						s3_secret_key = ""
					}
				}
				output_object_path = ""
				object_number_format {
					initial_value = 
					increment = 
					min_length = 
					place_holder = ""
				}
			}
			sample_snapshot_task {
				definition = 
				watermark_set {
					definition = 
					raw_parameter {
						type = ""
						coordinate_origin = ""
						x_pos = ""
						y_pos = ""
						image_template {
							image_content {
								type = ""
								cos_input_info {
									bucket = ""
									region = ""
									object = ""
								}
								url_input_info {
									url = ""
								}
								s3_input_info {
									s3_bucket = ""
									s3_region = ""
									s3_object = ""
									s3_secret_id = ""
									s3_secret_key = ""
								}
							}
							width = ""
							height = ""
							repeat_type = ""
						}
					}
					text_content = ""
					svg_content = ""
					start_time_offset = 
					end_time_offset = 
				}
				output_storage {
					type = ""
					cos_output_storage {
						bucket = ""
						region = ""
					}
					s3_output_storage {
						s3_bucket = ""
						s3_region = ""
						s3_secret_id = ""
						s3_secret_key = ""
					}
				}
				output_object_path = ""
				object_number_format {
					initial_value = 
					increment = 
					min_length = 
					place_holder = ""
				}
			}
			image_sprite_task {
				definition = 
				output_storage {
					type = ""
					cos_output_storage {
						bucket = ""
						region = ""
					}
					s3_output_storage {
						s3_bucket = ""
						s3_region = ""
						s3_secret_id = ""
						s3_secret_key = ""
					}
				}
				output_object_path = ""
				web_vtt_object_name = ""
				object_number_format {
					initial_value = 
					increment = 
					min_length = 
					place_holder = ""
				}
			}
			adaptive_dynamic_streaming_task {
				definition = 
				watermark_set {
					definition = 
					raw_parameter {
						type = ""
						coordinate_origin = ""
						x_pos = ""
						y_pos = ""
						image_template {
							image_content {
								type = ""
								cos_input_info {
									bucket = ""
									region = ""
									object = ""
								}
								url_input_info {
									url = ""
								}
								s3_input_info {
									s3_bucket = ""
									s3_region = ""
									s3_object = ""
									s3_secret_id = ""
									s3_secret_key = ""
								}
							}
							width = ""
							height = ""
							repeat_type = ""
						}
					}
					text_content = ""
					svg_content = ""
					start_time_offset = 
					end_time_offset = 
				}
				output_storage {
					type = ""
					cos_output_storage {
						bucket = ""
						region = ""
					}
					s3_output_storage {
						s3_bucket = ""
						s3_region = ""
						s3_secret_id = ""
						s3_secret_key = ""
					}
				}
				output_object_path = ""
				sub_stream_object_name = ""
				segment_object_name = ""
				add_on_subtitles {
					type = ""
					subtitle {
						type = ""
						cos_input_info {
							bucket = ""
							region = ""
							object = ""
						}
						url_input_info {
							url = ""
						}
						s3_input_info {
							s3_bucket = ""
							s3_region = ""
							s3_object = ""
							s3_secret_id = ""
							s3_secret_key = ""
						}
					}
				}
			}
			ai_content_review_task {
				definition = 
			}
			ai_analysis_task {
				definition = 
				extended_parameter = ""
			}
			ai_recognition_task {
				definition = 
			}
		}

  }
  output_storage {
		type = ""
		cos_output_storage {
			bucket = ""
			region = ""
		}
		s3_output_storage {
			s3_bucket = ""
			s3_region = ""
			s3_secret_id = ""
			s3_secret_key = ""
		}

  }
  output_dir = ""
  task_notify_config {
		cmq_model = ""
		cmq_region = ""
		topic_name = ""
		queue_name = ""
		notify_mode = ""
		notify_type = ""
		notify_url = ""
		aws_s_q_s {
			s_q_s_region = ""
			s_q_s_queue_name = ""
			s3_secret_id = ""
			s3_secret_key = ""
		}

  }
}

`
