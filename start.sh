#! /bin/bash
pkill -f 'pcs'
echo '######Stop PCS service down...'
echo '######Start PCS service again...'
nohup ./pcs &
#tail -f /export/Logs/pcs/pcs.log

