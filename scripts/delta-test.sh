#!/bin/bash

range_sha="${BASE_SHA}"

# changed files
changed_files="$(git diff --name-only "${range_sha}" | grep -E '\.go$')"
echo "Changed Files: ${changed_files}"

# changed test files
changed_test_files="$(echo "${changed_files}" | grep -E "_test\.go$")"
echo "Changed test files: ${changed_test_files}"

# test files for changed service
changed_service_files="$(echo "${changed_files}" | grep -E "^tencentcloud/(services/[^/]+/)?service_")"
changed_services_test_files=""
for service_file in ${changed_service_files}; do
  service_file_dir="$(dirname "${service_file}")"
  service_funcs="$(git diff "${range_sha}" "${service_file}" | grep -E "^.func " | awk -F ")" '{print $2}' | awk -F "(" '{print $1}' | tr -d ' ')"
  for func in ${service_funcs}; do
    files_using_func="$(grep -r --with-filename "${func}" "$service_file_dir" | awk -F ":" '{print $1}' | grep -E "^tencentcloud/(services/[^/]+/)?(resource_tc|data_source_tc)" | sort | uniq)"
    if [[ "${files_using_func}" == "" ]]; then
      continue
    fi
    files_using_func_test_files="$(echo "${files_using_func}" | awk -F "." '{print $1"_test.go"}')"
    changed_services_test_files="${changed_services_test_files} ${files_using_func_test_files}"
  done
done
echo "Test files for changed service functions: ${changed_services_test_files}"

# test files for changed resource and datasource
changed_sources_test_files="$(echo "${changed_files}" | grep -E "^tencentcloud/(services/[^/]+/)?(resource_tc|data_source_tc)" | grep -Ev "_test.go" | awk -F "." '{print $1"_test.go"}')"
echo "Test files for changed resource and datasource: ${changed_sources_test_files}"

# all need run test files
need_run_test_files="${changed_test_files} ${changed_services_test_files} ${changed_sources_test_files}"
need_run_test_files="$(echo "${need_run_test_files}" | xargs -n1 | sort | uniq)"
echo "All need run test files: ${need_run_test_files}"

# run test
for test_file in ${need_run_test_files}; do
  test_file_dir="$(dirname "$test_file")"
  test_casts="$(grep -E "func TestAcc.+\(" "${test_file}" | awk -F "(" '{print $1}' | awk '{print $2}' | grep -v "NeedFix")"
  printf "[%s]\n%s\n" "${test_file}" "${test_casts}"

  for test_cast in ${test_casts}; do
    go_test_cmd="go test -v -run ${test_cast} -timeout=0 ./${test_file_dir}/"
    $go_test_cmd
    if ! $go_test_cmd; then
      printf "[GO TEST FAILED] %s\n" "${go_test_cmd}"
      exit 1
    fi
  done
done
