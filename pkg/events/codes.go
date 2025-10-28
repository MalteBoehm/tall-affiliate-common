package events

// Mapping between orchestrator codes (e.g., "02A_PRODUCT_VALIDATED") and CloudEvents types.
// Keep this as the single source of truth for interop with non-Go orchestrators.

// CodeToCE maps orchestrator event codes to canonical CloudEvents types (domain.subdomain.action.v1)
var CodeToCE = map[string]string{
    // Phase 0: Scraper Intake
    "00A_SCRAPER_JOB_REQUESTED":          "scraper.job.requested.v1",

    // Phase 1/2: Discovery & Validation
    "01_PRODUCT_DETECTED":                CatalogProductDetectedV1,
    "02A_PRODUCT_VALIDATED":              CatalogProductValidatedV1,
    "02B_PRODUCT_IGNORED":                CatalogProductIgnoredV1,
    "02C_PRODUCT_REVIEW_REQUIRED":        CatalogProductReviewRequiredV1,

    // Phase 3: Enrichment Orchestration (PA-API)
    // Map STARTED/ENRICHED/COMPLETED to requested/completed for canonical CE types
    "03A_PA_API_ENRICHMENT_ORCHESTRATION_STARTED":   CatalogProductEnrichmentRequestedV1,
    "03B_PA_API_ENRICHMENT_ORCHESTRATION_ENRICHED":  CatalogProductEnrichmentCompletedV1,
    "03C_PA_API_ENRICHMENT_ORCHESTRATION_FAILED":    CatalogProductEnrichmentFailedV1,
    "03D_PA_API_ENRICHMENT_ORCHESTRATION_COMPLETED": CatalogProductEnrichmentCompletedV1,

    // Phase 4: Content Generation
    "04A_CONTENT_GENERATION_REQUESTED":    ContentGenerationRequestedV1,
    "04B_CONTENT_GENERATION_STARTED":      ContentGenerationStartedV1,
    "04C_CONTENT_GENERATED_COMPLETED":     ContentGeneratedV1,
    "04D_CONTENT_GENERATION_FAILED":       ContentGenerationFailedV1,
    "04E_CONTENT_GENERATION_RETRIED":      ContentGenerationRetriedV1,

    // Phase 5: Portal Publishing (new canonical CE types)
    "05A_PORTAL_PUBLICATION_REQUESTED":    "portal.publication.requested.v1",
    "05B_PORTAL_PUBLICATION_COMPLETED":    "portal.publication.completed.v1",
    "05C_PORTAL_PUBLICATION_FAILED":       "portal.publication.failed.v1",

    // Monitoring (MON_*) — content
    "MON_CONTENT_HEALTH_CHECK_SCHEDULED":  "monitoring.content.health_check.scheduled.v1",
    "MON_CONTENT_HEALTH_CHECK_COMPLETED":  "monitoring.content.health_check.completed.v1",
    "MON_CONTENT_HEALTH_CHECK_FAILED":     "monitoring.content.health_check.failed.v1",
    "MON_CONTENT_VISIBILITY_ALERT":        "monitoring.content.visibility.alert.v1",

    // Monitoring — price
    "MON_PRICE_UPDATE_SUCCEEDED":          PriceUpdatedV1,
    "MON_PRICE_UPDATE_FAILED":             PriceUpdateFailedV1,
    "MON_PRICE_UPDATE_SKIPPED":            "price.update.skipped.v1",
    "MON_PRICE_CHANGE_DETECTED":           "price.change.detected.v1",
    "MON_PRICE_HISTORY_RECORDED":          "price.history.recorded.v1",
    "MON_PRICE_HISTORY_RECORD_FAILED":     "price.history.record_failed.v1",
    "MON_PRICE_TREND_ANALYSIS_COMPLETED":  "price.trend_analysis.completed.v1",
    "MON_PRICE_ALERT_TRIGGERED":           "price.alert.triggered.v1",
    "MON_PRICE_ALERT_RESOLVED":            "price.alert.resolved.v1",

    // Monitoring — product
    "MON_PRODUCT_DATA_REFRESH_SCHEDULED":  "catalog.product.refresh.scheduled.v1",
    "MON_PRODUCT_DATA_REFRESH_COMPLETED":  "catalog.product.refresh.completed.v1",
    "MON_PRODUCT_DATA_REFRESH_FAILED":     "catalog.product.refresh.failed.v1",
    "MON_PRODUCT_AVAILABILITY_CHANGED":    ProductAvailabilityChangedV1,
    "MON_PRODUCT_STATUS_CHANGED":          ProductStatusChangedV1,
    "MON_PRODUCT_ARCHIVED":                "product.archived.v1",

    // Monitoring — social proof
    "MON_SOCIAL_PROOF_UPDATED":            "social_proof.updated.v1",
    "MON_SOCIAL_PROOF_UPDATE_FAILED":      "social_proof.update_failed.v1",

    // Monitoring — compliance
    "MON_COMPLIANCE_REVIEW_REQUIRED":      "compliance.review.required.v1",
    "MON_COMPLIANCE_REVIEW_PASSED":        "compliance.review.passed.v1",
    "MON_COMPLIANCE_REVIEW_FAILED":        "compliance.review.failed.v1",

    // Monitoring — supply chain
    "MON_SUPPLY_CHAIN_ALERT_RAISED":       "supply_chain.alert.raised.v1",
    "MON_SUPPLY_CHAIN_ALERT_RESOLVED":     "supply_chain.alert.resolved.v1",
}

// CEToCode is the inverse mapping for lookups
var CEToCode = func() map[string]string {
    m := make(map[string]string, len(CodeToCE))
    for k, v := range CodeToCE {
        // prefer first writer in case of duplicates
        if _, exists := m[v]; !exists {
            m[v] = k
        }
    }
    return m
}()

// GetCETypeForCode returns the CloudEvents type for the given orchestrator code
func GetCETypeForCode(code string) (string, bool) {
    v, ok := CodeToCE[code]
    return v, ok
}

// GetCodeForCEType returns the orchestrator code for the given CloudEvents type
func GetCodeForCEType(ceType string) (string, bool) {
    v, ok := CEToCode[ceType]
    return v, ok
}
