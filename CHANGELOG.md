# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2025-05-16

### Changed
- Modified `Time.MarshalJSON()` to safely serialize time values as JSON. Added support for environment variable `LUNO_TIME_LEGACY_FORMAT` to control the output format:
  - When `LUNO_TIME_LEGACY_FORMAT=true`: Returns the original string format (with proper JSON quoting)
  - Default behavior (no env var or any other value): Returns millisecond timestamp values

### Migration Guide
If you rely on the string format previously returned by `Time.MarshalJSON()`, you have three options:

1. **Option 1**: Set the `LUNO_TIME_LEGACY_FORMAT=true` environment variable in your application
2. **Option 2**: Update your client code to handle numeric millisecond timestamp values
3. **Option 3**: Use the provided migration helper in `_examples/time_migration/main.go`

Note that the original implementation could produce invalid JSON when embedded in larger JSON structures, so the new approach is recommended for better JSON compatibility.
