#!/usr/bin/env bash
ZIPNAME="go-ai-devops.zip"
DIR="go-ai-devops"
if [ -d "$DIR" ]; then
  zip -r $ZIPNAME $DIR
  echo "Created $ZIPNAME"
else
  echo "Directory $DIR not found"
fi
