#!/usr/bin/env bash
set -e

which ssh-agent || ( apt-get update -y && apt-get install openssh-client -y )

eval $(ssh-agent -s)
echo "$GOVKAPI_SSH_PRIVATE_KEY" | tr -d '\r' | ssh-add -

# build GO VKAPI Generator
go build -o go-vkapi-gen

# run generator and check exit code
./go-vkapi-gen || exit 1

# clone GO VKAPI repo
mkdir -p ~/.ssh
touch ~/.ssh/known_hosts
ssh-keyscan -t rsa gitlab.com >> ~/.ssh/known_hosts
git clone "$GOVKAPI_SSH_REPO_URL" || exit 1
cd "$GOVKAPI_REPO_DIR" || exit 1

git checkout master || exit 1
git config --global user.email "cicd-go-vkapi-gen@gitlab.com"
git config --global user.name "CIDI GO VK API Generator"

# create a new branch named (HOW??)
br_name=$(date +"generated-%m-%d-%Y-%H-%M-%S")
git checkout -b "$br_name" || exit 1

# copy generated output to the target repo dir
cp -Rfp ../output/* ./ || exit 1
ls -l

# add all files to a new branch and commit it to repo
git status
git add --all || exit 1
git status

git commit -m "Auto-generated VK API SDK build. $(date)" || exit 1
git push -u origin "$br_name" -o merge_request.create -o merge_request.remove_source_branch -o merge_request.title="Auto-generated VK API SDK build. $(date)" -o merge_request.target=master || exit 1

echo "DONE"
