# Changelog

## [0.5.0] - 2025-01-26

### Changed
- **BREAKING**: Migrated all unused events directly to CloudEvents format
  - Quality Assessment events (06, 07A, 07B) → CloudEvents constants
  - Reviews events (08B, 09B, 10C, 10D, 11A, 11B, 12A, 12B) → CloudEvents constants  
  - Product Ready for Publication (13) → CloudEvents constant
  - Scheduling events (14A, 14B, 14C) → CloudEvents constants
  - Price events (15A, 15B) → CloudEvents constants
  - Product Update events (16A, 16B) → CloudEvents constants
  - Product Status/Delete events (18, 19) → CloudEvents constants
  - Enrichment events (05A, 05B, 05C, 05D) → CloudEvents constants

### Added
- New CloudEvents constants for enrichment events:
  - `EnrichmentCompletedV1` = "enrichment.completed.v1"
  - `EnrichmentFailedV1` = "enrichment.failed.v1"
  - `VariantsEnrichedV1` = "variants.enriched.v1"
  - `EnrichmentRetryV1` = "enrichment.retry.v1"

### Migration Impact
This is a breaking change for any code that directly references the old string literals. Since all these events were identified as unused in production, the migration should have minimal impact. The NormalizeEventType function continues to handle both old and new formats for backward compatibility.

## [0.4.1] - 2025-01-26

### Fixed
- Fixed duplicate switch case compilation errors in NormalizeEventType
- Removed constant names from case statements, kept only string literals

## [0.4.0] - 2025-01-26

### Added
- CloudEvents v1.0 format support with backward compatibility
- New CloudEvents constants for all event types
- NormalizeEventType function for dual format support

### Changed
- All event publishing now supports both old and new formats
- Services can accept both numbered and CloudEvents formats

## Previous Versions
- See git history for changes before 0.4.0