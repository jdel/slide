
# Run standalone server on websocket 4333, disable nats port 4222 for user `UAFUAGK2MEB37BBA76MRKRGUS4KLDSJTBVT2U3TZZX34SLQGE4VCFVN3`

slide server --data-dir /tmp/slide --port 0 --ws-port 4333 --user UAFUAGK2MEB37BBA76MRKRGUS4KLDSJTBVT2U3TZZX34SLQGE4VCFVN3

# Run standalone server with custom config file

slide server --config /etc/slide/slide.conf
