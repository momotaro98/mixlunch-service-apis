#!/bin/bash

export DDL_PATH="db/init/01_ddl.sql"
export DDL_REL_PATH="./init/01_ddl.sql"

if [ -z "$1" ]; then
  echo "Need arg TARGET_BRANCH like 'master', 'develop'"
  exit 1
fi

export TARGET_BRANCH=$1

diff=$(diff -q <(git show $TARGET_BRANCH:$DDL_PATH) <(cat $DDL_REL_PATH))
if [ -z "$diff" ]; then
  echo "No difference"
  exit 0
fi

schemalex "local-git://.?file=$DDL_PATH&commitish=$TARGET_BRANCH" $DDL_REL_PATH
