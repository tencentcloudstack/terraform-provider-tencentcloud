#!/bin/bash

# Check goimports
echo "==> Checking that code complies with fmt requirements..."
goimports_files=$(goimports -l tencentcloud)
if [[ -n ${goimports_files} ]]; then
  echo 'fmt needs running on the following files:'
  echo "${goimports_files}"
  echo "You can use the command: \`make fmt\` to reformat code."
  exit 1
fi

exit 0
