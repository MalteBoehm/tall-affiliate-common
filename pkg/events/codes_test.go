package events

import "testing"

func TestCodeToCE_Bidirectional_ScraperCore(t *testing.T) {
    cases := map[string]string{
        "00A_SCRAPER_JOB_REQUESTED":   "scraper.job.requested.v1",
        "01_PRODUCT_DETECTED":         CatalogProductDetectedV1,
        "02A_PRODUCT_VALIDATED":       CatalogProductValidatedV1,
        "02B_PRODUCT_IGNORED":         CatalogProductIgnoredV1,
        "02C_PRODUCT_REVIEW_REQUIRED": CatalogProductReviewRequiredV1,
        // Content
        "04A_CONTENT_GENERATION_REQUESTED": ContentGenerationRequestedV1,
        "04B_CONTENT_GENERATION_STARTED":   ContentGenerationStartedV1,
        "04C_CONTENT_GENERATED_COMPLETED":  ContentGeneratedV1,
        "04D_CONTENT_GENERATION_FAILED":    ContentGenerationFailedV1,
        "04E_CONTENT_GENERATION_RETRIED":   ContentGenerationRetriedV1,
        // Refresh trigger (MON)
        "MON_PRODUCT_DATA_REFRESH_SCHEDULED": "catalog.product.refresh.scheduled.v1",
    }

    for code, wantCE := range cases {
        if got, ok := GetCETypeForCode(code); !ok {
            t.Fatalf("missing mapping for %s", code)
        } else if got != wantCE {
            t.Fatalf("mapping mismatch for %s: got %s want %s", code, got, wantCE)
        }
        if back, ok := GetCodeForCEType(wantCE); !ok {
            t.Fatalf("inverse missing for %s", wantCE)
        } else if back != code {
            // Some codes share same CE type (03B/03D) – skip strict check in those cases
            if !(code == "03B_PA_API_ENRICHMENT_ORCHESTRATION_ENRICHED" || code == "03D_PA_API_ENRICHMENT_ORCHESTRATION_COMPLETED") {
                t.Fatalf("inverse mismatch: ce=%s back=%s code=%s", wantCE, back, code)
            }
        }
    }
}
