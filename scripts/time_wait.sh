while sleep 1; do netstat -n | grep -i 8080 | grep -i time_wait | wc -l; done
