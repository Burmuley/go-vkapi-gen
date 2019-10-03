#!/usr/bin/env bash

which ssh-agent || ( apt-get update -y && apt-get install openssh-client -y )

eval $(ssh-agent -s)
echo "$GOVKAPI_SSH_PRIVATE_KEY" | tr -d '\r' | ssh-add -

# build GO VKAPI Generator
go build -o go-vkapi-gen

# run generator and check exit code
if [[ "$(./go-vkapi-gen)" -ne 0 ]]; then
  echo "ERROR RUNNING GO VKAPI GENERATOR"
  exit 1
fi

# clone GO VKAPI repo
git clone "$GOVKAPI_SSH_REPO_URL"
cd "$GOVKAPI_REPO_DIR" || exit 1

git config --global user.email "cicd-go-vkapi-gen@gitlab.com"
git config --global user.name "CIDI GO VK API Generator"

# create a new branch named (HOW??)
br_name=$(date +"generated-%m-%d-%Y-%H-%M-%S")
git checkout -b "$br_name"

# copy generated output to the target repo dir
cp -R ../output/ .

# TODO: add all files to a new branch and commit it to repo
git status
