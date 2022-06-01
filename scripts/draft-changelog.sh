#!/bin/sh

# this script can generate draft changelog instead of painful writing it
# run with `make changelog` then copy it!

version=""
gitTag=$(git describe --tag --abbrev=0)
IFS=. read -r major minor patch <<< "$gitTag"

major=${major:1}

type=$1

if [[ -z $type ]]; then
  read -r -p "Specify semver: major, minor, patch(default) " input

  type=$input
  if [[ -z $input ]]; then
    type="patch"
  fi
fi

case $type in
v*)
  version=$1
  ;;
major)
  version="$((major + 1)).0.0"
  ;;
minor)
  version="$major.$((minor + 1)).0"
  ;;
patch)
  version="$major.$minor.$((patch+1))"
esac

diffs=$(git diff --name-only HEAD $gitTag | grep "tencentcloud/*")


resource="^tencentcloud\/resource_tc_([a-z_]+)\.go$"
data="^tencentcloud\/data_source_tc_([a-z_]+)\.go$"
service="^tencentcloud\/service_([a-z_]+)\.go$"
test="([a-z_]+)_test$"

items=""


for file in ${diffs}; do
  module=""
  fileType="resource"
  if [[ $file =~ $resource ]]; then
    module="tencentcloud_${BASH_REMATCH[1]}"
  elif [[ $file =~ $data ]]; then
    fileType="data source"
    module="tencentcloud_${BASH_REMATCH[1]}"
  elif [[ $file =~ $service ]]; then
    module="tencentcloud_${BASH_REMATCH[1]}"
  fi

  if [[ $module =~ $test ]]; then
    module=${BASH_REMATCH[1]}
  fi
  if [[ $module != "" ]]; then
    item="* $fileType \`$module\`"
    if [[ ! $items =~ "$item" ]]; then
        items="$items\n$item"
    fi
  fi

done

LANG=en_US
dateStr=$(date +"%B %d, %Y")

template="
## $version $dateStr
\n
\nFEATURES:
\nDEPRECATED:
\nENHANCEMENTS:
\nBUGFIXES:
\nCOMMON:
\n
$items
\n
"

echo $template