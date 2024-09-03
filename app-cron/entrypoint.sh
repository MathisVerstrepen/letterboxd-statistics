#!/bin/bash

# Add environment variables to /etc/environment so they are available to cron
printenv | grep -v "no_proxy" >> /etc/environment

# Start the cron service in the background
cron -f