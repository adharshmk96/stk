#!/bin/bash

# Download the binary
curl -L -o /tmp/stk https://github.com/adharshmk96/stk/releases/download/v0.6.6/stk

# Move the binary to /bin
sudo mv /tmp/stk /bin

# Make the binary executable
sudo chmod +x /bin/stk