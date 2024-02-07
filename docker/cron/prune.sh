#!/bin/bash
# Function to print the timestamped message
function log_timestamp() {
    echo "$(date +"%Y-%m-%d %H:%M:%S") $1"
}

log_timestamp "Stopping Bor node ..."
supervisorctl stop bor

# Wait for the bor process to stop
while supervisorctl status bor | grep -q "RUNNING"; do
    sleep 5
done

log_timestamp "Executing snapshot prune-state command..."
supervisorctl start snapshot_prune

while supervisorctl status snapshot_prune | grep -q "RUNNING"; do
    sleep 5
done


log_timestamp "Starting Bor node..."
supervisorctl start bor

log_timestamp "Bor node started successfully."

