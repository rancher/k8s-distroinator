#!/bin/sh

# Run confd
confd -onetime -backend env

# Start nginx
nginx -g 'daemon off;'
