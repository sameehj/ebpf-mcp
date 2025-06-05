#!/bin/bash

set -e

# Config
REMOTE_HOST="aws-t3-small"     # must be defined in ~/.ssh/config
REMOTE_USER="ec2-user"
REMOTE_DIR="~/ebpf-mcp"   # adjust as needed

echo "➡️  Syncing files to $REMOTE_USER@$REMOTE_HOST:$REMOTE_DIR"

# Sync everything except .git and node_modules (add more in .deployignore if needed)
rsync -avz --exclude '.git' --exclude 'node_modules' ./ $REMOTE_USER@$REMOTE_HOST:$REMOTE_DIR

echo "✅ Deploy complete. Logged in:"
ssh $REMOTE_HOST "cd $REMOTE_DIR && ls -lah"
