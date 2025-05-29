#!/usr/bin/env sh

aider \
    --no-auto-commits \
    --architect \
    --model openai/deepseek-ai/DeepSeek-R1-0528 \
    --editor-model openai/deepseek-ai/DeepSeek-R1

# https://aider.chat/docs/config/aider_conf.html
# default aider config lives at ~/.aider.conf.yml
# because macOS engineers don't believe in the
# .config directory for some reason
#
# pass your own config file like so:
# --config path/to/config.yaml \
