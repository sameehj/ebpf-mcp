#!/bin/bash
set -e

MAP_NAME="test_map"
MAP_PATH="/sys/fs/bpf/$MAP_NAME"

echo "🔧 Creating eBPF test map..."
sudo bpftool map create "$MAP_PATH" type hash key 4 value 4 entries 128 name "$MAP_NAME"

echo "📥 Inserting test key-value pair..."
sudo bpftool map update name "$MAP_NAME" key 0x01 0x00 0x00 0x00 value 0x2a 0x00 0x00 0x00

echo "🧪 Calling map_dump via JSON-RPC..."
curl -s -X POST http://localhost:8080/rpc \
  -H "Content-Type: application/json" \
  -d "{
    \"jsonrpc\": \"2.0\",
    \"method\": \"tools/call\",
    \"params\": {
      \"tool\": \"map_dump\",
      \"input\": {
        \"map_name\": \"$MAP_NAME\",
        \"max_entries\": 10
      }
    },
    \"id\": 2
  }" | jq .

echo "🧹 Cleaning up test map..."
sudo rm "$MAP_PATH"

echo "✅ Done."
