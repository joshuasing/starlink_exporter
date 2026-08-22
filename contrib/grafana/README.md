# Grafana dashboards

Example dashboards for use with starlink_exporter.

## Starlink Dish (`starlink.json`)

Overview of dish health and performance: status row (up, uptime, current
obstruction, 24h obstructed time, power draw, PoP latency), actual
uplink/downlink throughput, PoP ping latency and drop ratio, obstruction
fraction, and power draw over time.

**Import**: Grafana → Dashboards → New → Import → upload `starlink.json`,
then select your Prometheus data source. Assumes the default `starlink` job
scrape interval of ~15s.
