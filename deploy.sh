#!/usr/bin/env bash
set -e

which ssh-agent || ( apt-get update -y && apt-get install openssh-client -y )

eval $(ssh-agent -s)
echo "$GOVKAPI_SSH_PRIVATE_KEY" | tr -d '\r' | ssh-add -

# clone GO VKAPI repo
mkdir -p ~/.ssh
touch ~/.ssh/known_hosts
ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts
git clone "$GOVKAPI_SSH_REPO_URL" || exit 1
GOVKAPI_REPO_DIR=$(echo $GOVKAPI_SSH_REPO_URL | awk -F "\/" '{print $2}' | cut -d "." -f 1)
cd "$GOVKAPI_REPO_DIR" || exit 1

git checkout master || exit 1
git config --global user.email "cicd-go-vkapi-gen@github.com"
git config --global user.name "GitHub Actions GO VK API Generator"

# create a new branch named (HOW??)
br_name=$(date +"generated-%m-%d-%Y-%H-%M-%S")
git checkout -b "$br_name" || exit 1

# copy generated output to the target repo dir
find . -type f -not -path "./.git/*" -exec rm -f '{}' \; || exit 1
# include hidden files and directories
shopt -s dotglob
cp -aRfp ../output/* ./ || exit 1
ls -la

# add all files to a new branch and commit it to repo
git status
git add --all
CHANGED=$(git status -s | wc -l)
echo CHANGED=$CHANGED
if [[ $CHANGED != 0 ]]; then
  git add . --all || exit 1
  git status
  git commit -m "Auto-generated VK API SDK build. $(date)" || exit 1
  git push -u origin "$br_name" || exit 1
else
  echo "No changes in destination code."
fi
echo "DONE"

