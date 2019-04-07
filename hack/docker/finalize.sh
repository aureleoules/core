#!/bin/sh
set -x
set -e

# Create user for backpulse
addgroup -S backpulse
adduser -G backpulse -H -D -g 'backpulse User' backpulse -h /data/backpulse -s /bin/bash && usermod -p '*' backpulse && passwd -u backpulse
echo "export BACKPULSE_CUSTOM=${BACKPULSE_CUSTOM}" >> /etc/profile

# Final cleaning
rm /app/backpulse/docker/finalize.sh
rm /app/backpulse/docker/nsswitch.conf
