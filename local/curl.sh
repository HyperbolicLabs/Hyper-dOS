#!/usr/bin/env sh

# curl -X POST "https://dev-baron.ts.net/api/v1/namespaces/hyperdos/services/hyperai-vllm:http/proxy/v1/chat/completions" \
curl -X POST "http://localhost:8000" \
      -H "Content-Type: application/json" \
      --data-raw '{
          "messages": [{
            "role": "user",
            "content": "write a short poem about mycelium"
          }],
          "model": "TinyLlama/TinyLlama-1.1B-Chat-v1.0",
          "max_tokens": 512,
          "temperature": 0.7,
          "top_p": 0.9,
          "stream": false
      }' \
| jq -r ".choices[0].message.content"
