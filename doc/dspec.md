# Test Diagnostic Redis Specification

This document specifies how test diagnostics are placed into Redis database using the wtf conventions.

Redis organizes things in terms of keys and values. Each test diagnostic has two keys and two values.

```
test:participant:counter              level:::message
test:participant:counter:~time~       [unix-time, microseconds]
```

## Diagnostic Payload
The first entry has what we will call the _base key value_ for its key. The following are the semantics behind the components of the key.
- `test` - the name of the test the diagnostic belongs to
- `participant` - the name of the instance of the application that produced the diagnostic.
- `counter` - a counter that is monotonically increasing per test-participant pair

The value of the first entry is a level that is one of
- error
- warning
- ok
and a diagnostic message that comes from the testing application.

## Diagnostic Timestamp
The second entry keeps track of the time at which a diagnostic is created from the collectors temporal frame of reference. The key for this diagnostic is the base key plus a `~time~` suffix to indicate that this is the timestamp for is representative base key. The value of the timestamp follows the same format as the output of the Redis [TIME](https://redis.io/commands/time) command, a Unix timestamp in seconds plus the number of microseconds within the Unix second.

Redis itself does not keep track of this time stamping. It is up to the language drivers to actually create and push the time stamps to the database. See the C, Python or Perl drivers for examples of how this is done.
