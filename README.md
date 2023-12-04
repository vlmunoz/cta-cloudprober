# cta-cloudprober
Cloudprober configuration and probes for monitoring CTA

To set up cloudprober as a systemd service:

1. cp cloudprober.service to /etc/systemd/system/
1. Update PGHOST and PGPASSWORD env variables in /etc/systemd/system/cloudprober.service to real values
1. systemctl daemon-reload
1. systemctl enable cloudprober
1. systemctl start cloudprober
