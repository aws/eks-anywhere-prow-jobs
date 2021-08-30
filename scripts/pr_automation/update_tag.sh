#!/usr/bin/env bash
# Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


set -e
set -o pipefail
set -x

SCRIPT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
REPO="eks-anywhere-prow-jobs"

OLD_TAG="$(cat ${SCRIPT_ROOT}/../../BUILDER_BASE_TAG_FILE)"
LATEST_TAG="$(curl https://raw.githubusercontent.com/aws/eks-distro-prow-jobs/main/BUILDER_BASE_TAG_FILE)"
if [ "$OLD_TAG" = "$LATEST_TAG" ]; then
    echo "Builder-base tag in Prowjobs is up to date!"
    exit 0
fi

echo "Builder-base tag in Prowjobs is out of date! Updating"

if [ -z "$REPO_OWNER" ]; then
    echo "No org information was provided, please set and export REPO_OWNER environment variable. \
      This is used to raise a pull request against your org after updating tags in the respective files."
    exit 1
fi
if [ $REPO_OWNER = "aws" ]; then
    ORIGIN_ORG="eks-distro-pr-bot"
    UPSTREAM_ORG="aws"
else
    ORIGIN_ORG=$REPO_OWNER
    UPSTREAM_ORG=$REPO_OWNER
fi

COMMIT_MESSAGE="[PR BOT] Update builder-base image tag in Prow jobs"
PR_TITLE="Update base image tag in Prowjobs"
PR_BODY=$(cat ${SCRIPT_ROOT}/builder_base_pr_body)
PR_BRANCH="image-tag-update"

git config --global push.default current
git remote add origin git@github.com:${ORIGIN_ORG}/${REPO}.git
git remote add upstream git@github.com:${UPSTREAM_ORG}/${REPO}.git
git checkout -b $PR_BRANCH

for FILE in $(find ./ -type f -name "*.yaml"); do
    sed -i "s,${OLD_TAG},${LATEST_TAG}," $FILE
    git add $FILE
done
sed -i "s,.*,${LATEST_TAG}," ./BUILDER_BASE_TAG_FILE
git add ./BUILDER_BASE_TAG_FILE

FILES_ADDED=$(git diff --staged --name-only)
if [ "$FILES_ADDED" = "" ]; then
    exit 0
fi

git commit -m "$COMMIT_MESSAGE"
ssh-agent bash -c 'ssh-add /secrets/ssh-secrets/ssh-privatekey; ssh -o StrictHostKeyChecking=no git@github.com; git fetch upstream; git rebase upstream/main; git push -u origin $PR_BRANCH -f'

gh auth login --with-token < /secrets/github-secrets/token

PR_EXISTS=$(gh pr list | grep -c "${PR_BRANCH}" || true)
if [ $PR_EXISTS -eq 0 ]; then
  gh pr create --title "$PR_TITLE" --body "$PR_BODY"
fi
