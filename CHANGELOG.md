# Changelog

All notable changes to this project will be documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v0.9.0] - 2026-04-19

### Added

- Support Starlink Dishy firmware `2026.04.07.mr77639.1` ([#201])

### Fixed

- Fixed scraper race condition, deadlock and descriptor bugs ([#186], [#187])

### Changed

- Added `tar.zst` release archives alongside `tar.gz` ([#202])
- Improved build reproducibility with deterministic binary timestamps ([#202])
- Updated to Go 1.26.2 ([#200], [#199])
- Updated dependencies ([#168], [#177], [#182], [#193])

## [v0.8.0] - 2026-02-16

### Changed

- Support Starlink Dishy firmware `2026.01.31.mr72966` ([#160])
- Updated to Go 1.26 ([#159], [#151])
- Updated dependencies ([#159])

## [v0.7.4] - 2026-01-11

### Changed

- Updated dependencies ([#136], [#138], [#139])

## [v0.7.3] - 2026-01-04

### Added

- Support Starlink Dishy firmware `2025.12.07.mr69330.2` (API v42)
  ([#134])

### Changed

- Updated Go to 1.25.5 ([#123], [#124])
- Updated dependencies ([#112], [#119], [#121], [#125], [#127], [#131])

## [v0.7.1] - 2025-10-27

### Added

- Support Starlink Dishy firmware `2025.10.08.mr65265` (API v40)
  ([#97])

### Changed

- Updated Go to 1.25 ([#103])
- Updated GoReleaser configuration ([#102], [#104], [#105])

### Removed

- Removed support for `linux/arm/v7` in Docker images ([#105])

-----

_Looking for the changelog for an older version? Older releases can be found at:
https://github.com/joshuasing/starlink_exporter/releases_

[Unreleased]: https://github.com/joshuasing/starlink_exporter/compare/v0.9.0...HEAD
[v0.9.0]: https://github.com/joshuasing/starlink_exporter/compare/v0.8.0...v0.9.0
[v0.8.0]: https://github.com/joshuasing/starlink_exporter/compare/v0.7.4...v0.8.0
[v0.7.4]: https://github.com/joshuasing/starlink_exporter/compare/v0.7.3...v0.7.4
[v0.7.3]: https://github.com/joshuasing/starlink_exporter/compare/v0.7.1...v0.7.3
[v0.7.1]: https://github.com/joshuasing/starlink_exporter/releases/tag/v0.7.1

[#97]: https://github.com/joshuasing/starlink_exporter/pull/97
[#102]: https://github.com/joshuasing/starlink_exporter/pull/102
[#103]: https://github.com/joshuasing/starlink_exporter/pull/103
[#104]: https://github.com/joshuasing/starlink_exporter/pull/104
[#105]: https://github.com/joshuasing/starlink_exporter/pull/105
[#112]: https://github.com/joshuasing/starlink_exporter/pull/112
[#119]: https://github.com/joshuasing/starlink_exporter/pull/119
[#121]: https://github.com/joshuasing/starlink_exporter/pull/121
[#123]: https://github.com/joshuasing/starlink_exporter/pull/123
[#124]: https://github.com/joshuasing/starlink_exporter/pull/124
[#125]: https://github.com/joshuasing/starlink_exporter/pull/125
[#127]: https://github.com/joshuasing/starlink_exporter/pull/127
[#131]: https://github.com/joshuasing/starlink_exporter/pull/131
[#134]: https://github.com/joshuasing/starlink_exporter/pull/134
[#136]: https://github.com/joshuasing/starlink_exporter/pull/136
[#138]: https://github.com/joshuasing/starlink_exporter/pull/138
[#139]: https://github.com/joshuasing/starlink_exporter/pull/139
[#151]: https://github.com/joshuasing/starlink_exporter/pull/151
[#159]: https://github.com/joshuasing/starlink_exporter/pull/159
[#160]: https://github.com/joshuasing/starlink_exporter/pull/160
[#168]: https://github.com/joshuasing/starlink_exporter/pull/168
[#177]: https://github.com/joshuasing/starlink_exporter/pull/177
[#182]: https://github.com/joshuasing/starlink_exporter/pull/182
[#186]: https://github.com/joshuasing/starlink_exporter/pull/186
[#187]: https://github.com/joshuasing/starlink_exporter/pull/187
[#193]: https://github.com/joshuasing/starlink_exporter/pull/193
[#199]: https://github.com/joshuasing/starlink_exporter/pull/199
[#200]: https://github.com/joshuasing/starlink_exporter/pull/200
[#201]: https://github.com/joshuasing/starlink_exporter/pull/201
[#202]: https://github.com/joshuasing/starlink_exporter/pull/202
