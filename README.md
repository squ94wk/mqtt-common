# mqtt-common

This module offers types and functions that implement structures used in the mqtt (version 5) protocol.

## Potential v1.0.0
A version `1.0.0` would have to include:
* Full MQTT 5.0 compliance
* Refined API, proven to work well in broker & client environment
* Return structured errors
* Use appropriate error codes & messages
* `pprof`-ed: avoiding unnecessary allocations
* move allocations to the user
* fuzz-tested for resiliency

## Possible features
* Lazy reads for big payloads (application messages/will messages) taking io.Writer to avoid mem-bloat
