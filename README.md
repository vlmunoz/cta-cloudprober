# cta-cloudprober
Cloudprober configuration and probes for monitoring CTA

Includes the following probers:

1. Aliveness prober for eos (depends on eos-client)
1. Aliveness prober for postgres (depends on psql client)
1. Aliveness prober for quarkdb
1. Network connectivity prober
1. CTA drive status prober (depends on psql client)

To set up cloudprober as a systemd service:

1. Install go (https://go.dev/doc/install)
1. `cp cloudprober.service /etc/systemd/system/`
1. Update PGHOST, PGPORT and PGPASSWORD env variables in /etc/systemd/system/cloudprober.service to real values
1. `systemctl daemon-reload`
1. `systemctl enable cloudprober`
1. `systemctl start cloudprober`

Firewall must be open on port 9313, then dashboard can be seen at http://hostname:9313/status, and metrics can be seen at http://hostname:9313/metrics.

To rebuild probers:

1. Modify source files (end in `.go`)
1. `export GO111MODULE=on`
1. `go build <prober>.go`
