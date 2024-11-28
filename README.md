# Liveness / Readiness Example

This application is a quick and dirty little application for showing an example of application that goes unresponsive because its saturated or busy.

If you aren't careful and set your liveness probe too agressive it'll kill the pod that is busy working.  Similiarly if you set your readiness probe too permissive you could end up with an application unable to accept work still having traffic sent to it.

## Liveness Probe
A liveness probe sole job is to check and see if an application is stuck and needs restarted.  Doesn't need to be aggressive.

## Readiness Probe
A readiness probe should be a bit more aggressive.  If a pod is saturated and unable to respond as fast as you've budgeted for.. You shouldn't be sending traffic to it.

As soon as it fails it goes to the unready status.  This removes the endpoint in kubernetes so it is no longer sent traffic.

# Usage

The yaml file can be deployed.  By default its not bound to a specific host its just listening on /.  You'll likely want to uncomment the host field for your own testing.

Calling /block?time=4s will cause all calls to /health to timeout for 4 seconds.
Calling / will tell you when its blocked until.

Internally its using a super niave lock mechanism and the /health is just trying to get a lock and immediately releasing it once it gets it.
