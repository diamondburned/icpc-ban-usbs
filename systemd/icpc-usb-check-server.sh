#!/usr/bin/env bash
competition_status=$(wget -qO - localhost:8080/competition-status)
echo "competition is $competition_status" 1>&2
# Only succeed if the competition status is finished.
[[ $competition_status == finished ]]
