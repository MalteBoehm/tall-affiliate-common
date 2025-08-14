package patterns

import (
	"strings"
)

// BrowseNode interface defines the methods that browse node implementations must provide
type BrowseNode interface {
	GetID() string
	GetDisplayName() string
	GetContextFreeName() string
	IsRoot() bool
}

// BrowseNodeData represents a single Browse Node from PA-API
type BrowseNodeData struct {
	ID              string `json:"id"`
	DisplayName     string `json:"displayName"`
	ContextFreeName string `json:"contextFreeName"`
	Ancestor        struct {
		ID              string `json:"id"`
		DisplayName     string `json:"displayName"`
		ContextFreeName string `json:"contextFreeName"`
	} `json:"ancestor"`
}

// GetID returns the browse node ID
func (b *BrowseNodeData) GetID() string {
	return b.ID
}

// GetDisplayName returns the display name
func (b *BrowseNodeData) GetDisplayName() string {
	return b.DisplayName
}

// GetContextFreeName returns the context-free name
func (b *BrowseNodeData) GetContextFreeName() string {
	return b.ContextFreeName
}

// IsRoot returns whether this is a root node (always false for PA-API browse nodes)
func (b *BrowseNodeData) IsRoot() bool {
	return false // PA-API browse nodes are typically not root nodes
}

// BrowseNodeScore represents scoring result for a Browse Node
type BrowseNodeScore struct {
	BrowseNodeID string `json:"browse_node_id"`
	Category     string `json:"category"`
	Gender       string `json:"gender"`
	Score        int    `json:"score"`
	Priority     int    `json:"priority"`
	IsSupported  bool   `json:"is_supported"`
	DisplayName  string `json:"display_name"`
}

// MultiNodeAnalysisResult represents the complete analysis result
type MultiNodeAnalysisResult struct {
	BestGender         string            `json:"best_gender"`
	BestCategory       string            `json:"best_category"`
	TotalScore         int               `json:"total_score"`
	NodeScores         []BrowseNodeScore `json:"node_scores"`
	SupportedNodes     []string          `json:"supported_nodes"`
	UnsupportedNodes   []string          `json:"unsupported_nodes"`
	ContentGender      string            `json:"content_gender"`
	ContentCategory    string            `json:"content_category"`
	CombinedConfidence float64           `json:"combined_confidence"`
}

// PatternRule represents a pattern matching rule with weight and language
type PatternRule struct {
	Patterns []string // List of patterns to match
	Weight   int      // Weight for scoring (higher = more important)
	Language string   // "de", "en", or "both"
}

// PatternMatch represents a pattern match result
type PatternMatch struct {
	Category string
	Score    int
	Matched  []string
}

// DetectCategory analyzes text to determine the product category using intelligent pattern matching
func DetectCategory(text string) string {
	// Use enhanced pattern matching with fuzzy matching and stemming
	categoryPatterns := getCategoryPatterns()

	for category, patterns := range categoryPatterns {
		if matchesPatternWithFuzzy(text, patterns) {
			return category
		}
	}

	return ""
}

// DetectGender analyzes text to determine the target gender using intelligent pattern matching
func DetectGender(text string) string {
	textLower := strings.ToLower(text)

	// Check for high-priority explicit gender markers first
	if strings.Contains(textLower, "damen") || strings.Contains(textLower, "frauen") {
		return "women"
	}
	if strings.Contains(textLower, "herren") || strings.Contains(textLower, "männer") {
		return "men"
	}

	// If no explicit markers, use pattern scoring
	genderPatterns := getGenderPatterns()
	priorityOrder := []string{"women", "men", "unisex"}
	scores := make(map[string]int)
	maxScore := 0

	for _, gender := range priorityOrder {
		if patterns, exists := genderPatterns[gender]; exists {
			score := calculateGenderScore(text, patterns)
			scores[gender] = score
			if score > maxScore {
				maxScore = score
			}
		}
	}

	// Require minimum score for confidence
	if maxScore < 5 {
		return ""
	}

	// Return the first gender in priority order that has the maximum score
	for _, gender := range priorityOrder {
		if scores[gender] == maxScore {
			return gender
		}
	}

	return ""
}

// GenerateTagsFromText generates relevant tags from text analysis using intelligent pattern matching
func GenerateTagsFromText(text string) []string {
	var tags []string
	tagPatterns := getTagPatterns()

	for tag, patterns := range tagPatterns {
		if matchesPatternWithFuzzy(text, patterns) {
			tags = append(tags, tag)
		}
	}

	// Add hierarchical tags based on detected categories
	tags = append(tags, generateHierarchicalTags(text)...)

	return tags
}

// DetermineBestGenderFromList determines the best gender from a list with priority logic
func DetermineBestGenderFromList(genders []string) string {
	// Count occurrences
	genderCount := make(map[string]int)
	for _, gender := range genders {
		if gender != "" {
			genderCount[gender]++
		}
	}

	// Priority: specific genders over unisex, women over men if both present
	if genderCount["women"] > 0 {
		return "women"
	}
	if genderCount["men"] > 0 {
		return "men"
	}
	if genderCount["unisex"] > 0 {
		return "unisex"
	}

	return ""
}

// DetermineBestCategory determines the best category from a list of candidates
func DetermineBestCategory(categories []string) string {
	if len(categories) == 0 {
		return ""
	}

	// Count occurrences of each category
	categoryCount := make(map[string]int)
	for _, category := range categories {
		if category != "" {
			categoryCount[category]++
		}
	}

	// Find the most frequent category
	var bestCategory string
	var maxCount int
	for category, count := range categoryCount {
		if count > maxCount {
			maxCount = count
			bestCategory = category
		}
	}

	return bestCategory
}

// DetectGenderWithScores returns both the detected gender and the scores for debugging
func DetectGenderWithScores(text string) (string, map[string]int) {
	genderPatterns := getGenderPatterns()

	// Calculate scores for all genders with explicit priority order
	priorityOrder := []string{"women", "men", "unisex"}
	scores := make(map[string]int)
	maxScore := 0

	for _, gender := range priorityOrder {
		if patterns, exists := genderPatterns[gender]; exists {
			score := calculateGenderScore(text, patterns)
			scores[gender] = score
			if score > maxScore {
				maxScore = score
			}
		}
	}

	// Require minimum score for confidence
	if maxScore < 5 {
		return "", scores
	}

	// Return the first gender in priority order that has the maximum score
	for _, gender := range priorityOrder {
		if scores[gender] == maxScore {
			return gender, scores
		}
	}

	return "", scores
}

// DetectGenderFromBrowseNodes analyzes Browse Nodes to determine gender using hierarchical scoring
func DetectGenderFromBrowseNodes(browseNodes []BrowseNode) (string, []BrowseNodeScore) {
	browseNodeMappings := getExtendedBrowseNodeMappings()
	scores := make([]BrowseNodeScore, 0)
	genderScores := make(map[string]int)

	for _, node := range browseNodes {
		if mapping, exists := browseNodeMappings[node.GetID()]; exists {
			score := calculatePriorityScore(mapping.Priority)

			nodeScore := BrowseNodeScore{
				BrowseNodeID: node.GetID(),
				Category:     mapping.Category,
				Gender:       mapping.Gender,
				Score:        score,
				Priority:     mapping.Priority,
				IsSupported:  mapping.IsSupported,
				DisplayName:  mapping.DisplayName,
			}
			scores = append(scores, nodeScore)

			// Only count supported nodes for gender scoring
			if mapping.IsSupported {
				genderScores[mapping.Gender] += score
			}
		}
	}

	// Determine best gender based on priority order and scores
	priorityOrder := []string{"women", "men", "unisex"}
	bestGender := ""
	maxScore := 0

	for _, gender := range priorityOrder {
		if score, exists := genderScores[gender]; exists && score > maxScore {
			maxScore = score
			bestGender = gender
		}
	}

	return bestGender, scores
}

// DetectCategoryFromBrowseNodes analyzes Browse Nodes to determine category using hierarchical scoring
func DetectCategoryFromBrowseNodes(browseNodes []BrowseNode) (string, []BrowseNodeScore) {
	browseNodeMappings := getExtendedBrowseNodeMappings()
	scores := make([]BrowseNodeScore, 0)
	categoryScores := make(map[string]int)

	for _, node := range browseNodes {
		if mapping, exists := browseNodeMappings[node.GetID()]; exists {
			score := calculatePriorityScore(mapping.Priority)

			nodeScore := BrowseNodeScore{
				BrowseNodeID: node.GetID(),
				Category:     mapping.Category,
				Gender:       mapping.Gender,
				Score:        score,
				Priority:     mapping.Priority,
				IsSupported:  mapping.IsSupported,
				DisplayName:  mapping.DisplayName,
			}
			scores = append(scores, nodeScore)

			// Only count supported nodes for category scoring
			if mapping.IsSupported {
				categoryScores[mapping.Category] += score
			}
		}
	}

	// Find category with highest score
	bestCategory := ""
	maxScore := 0

	for category, score := range categoryScores {
		if score > maxScore {
			maxScore = score
			bestCategory = category
		}
	}

	return bestCategory, scores
}

// calculatePriorityScore converts priority to hierarchical score
func calculatePriorityScore(priority int) int {
	switch priority {
	case 1:
		return 100 // Highest specificity
	case 2:
		return 80 // High specificity
	case 3:
		return 60 // Medium specificity
	case 4:
		return 40 // General category
	case 5:
		return 20 // Very general
	case 6:
		return 0 // Promotions/Sales (ignored)
	default:
		return 0 // Unknown priority
	}
}

// ===== MULTI-BROWSE-NODE ANALYZER =====

// AnalyzeMultipleBrowseNodes performs comprehensive analysis combining Browse Nodes + Content patterns
func AnalyzeMultipleBrowseNodes(browseNodes []BrowseNode, contentText string) MultiNodeAnalysisResult {
	result := MultiNodeAnalysisResult{
		NodeScores:       make([]BrowseNodeScore, 0),
		SupportedNodes:   make([]string, 0),
		UnsupportedNodes: make([]string, 0),
	}

	// 1. BROWSE NODE ANALYSIS
	browseGender, browseGenderScores := DetectGenderFromBrowseNodes(browseNodes)
	browseCategory, browseCategoryScores := DetectCategoryFromBrowseNodes(browseNodes)

	// Combine all browse node scores
	allBrowseScores := make(map[string]BrowseNodeScore)
	for _, score := range browseGenderScores {
		allBrowseScores[score.BrowseNodeID] = score
	}
	for _, score := range browseCategoryScores {
		if existing, exists := allBrowseScores[score.BrowseNodeID]; exists {
			// Keep existing score but ensure category/gender are populated
			if existing.Category == "" {
				existing.Category = score.Category
			}
			if existing.Gender == "" {
				existing.Gender = score.Gender
			}
			allBrowseScores[score.BrowseNodeID] = existing
		} else {
			allBrowseScores[score.BrowseNodeID] = score
		}
	}

	// Convert map to slice and categorize nodes
	totalBrowseScore := 0
	for _, score := range allBrowseScores {
		result.NodeScores = append(result.NodeScores, score)
		if score.IsSupported {
			result.SupportedNodes = append(result.SupportedNodes, score.BrowseNodeID)
			totalBrowseScore += score.Score
		} else {
			result.UnsupportedNodes = append(result.UnsupportedNodes, score.BrowseNodeID)
		}
	}

	// 2. CONTENT PATTERN ANALYSIS
	contentGender := DetectGender(contentText)
	contentCategory := DetectCategory(contentText)
	contentGenderScores := make(map[string]int)
	if contentGender != "" {
		_, contentGenderScores = DetectGenderWithScores(contentText)
	}

	// 3. HYBRID SCORING & CONFIDENCE CALCULATION
	result.BestGender = determineBestGenderHybrid(browseGender, contentGender, totalBrowseScore, contentGenderScores)
	result.BestCategory = determineBestCategoryHybrid(browseCategory, contentCategory, totalBrowseScore)
	result.TotalScore = totalBrowseScore
	result.ContentGender = contentGender
	result.ContentCategory = contentCategory
	result.CombinedConfidence = calculateCombinedConfidence(result, contentGenderScores)

	return result
}

// determineBestGenderHybrid combines Browse Node + Content Pattern results with intelligent weighting
func determineBestGenderHybrid(browseGender, contentGender string, browseScore int, contentScores map[string]int) string {
	// High Priority: If content explicitly mentions "Damen/Herren", trust content
	if contentGender == "women" || contentGender == "men" {
		return contentGender
	}

	// Medium Priority: Strong browse node evidence (Priority 1-2 nodes)
	if browseGender != "" && browseScore >= 80 {
		return browseGender
	}

	// Low Priority: Weak evidence - prefer content over unisex
	if contentGender != "" {
		return contentGender
	}

	// Fallback: Use browse node result
	return browseGender
}

// determineBestCategoryHybrid combines Browse Node + Content Pattern results for category
func determineBestCategoryHybrid(browseCategory, contentCategory string, browseScore int) string {
	// If content has specific category match, prefer content
	if contentCategory != "" && (contentCategory == "T-Shirts" || contentCategory == "Poloshirts" || contentCategory == "Hoodies") {
		return contentCategory
	}

	// Strong browse node evidence
	if browseCategory != "" && browseScore >= 80 {
		return browseCategory
	}

	// Fallback logic
	if contentCategory != "" {
		return contentCategory
	}

	return browseCategory
}

// calculateCombinedConfidence calculates confidence score based on evidence strength
func calculateCombinedConfidence(result MultiNodeAnalysisResult, contentScores map[string]int) float64 {
	confidence := 0.0

	// Browse Node Confidence (0-50%)
	if result.TotalScore >= 100 {
		confidence += 0.5 // Strong browse node evidence
	} else if result.TotalScore >= 40 {
		confidence += 0.3 // Medium browse node evidence
	} else if result.TotalScore > 0 {
		confidence += 0.1 // Weak browse node evidence
	}

	// Content Pattern Confidence (0-50%)
	maxContentScore := 0
	for _, score := range contentScores {
		if score > maxContentScore {
			maxContentScore = score
		}
	}

	if maxContentScore >= 10 {
		confidence += 0.5 // Strong content evidence
	} else if maxContentScore >= 5 {
		confidence += 0.3 // Medium content evidence
	} else if maxContentScore > 0 {
		confidence += 0.1 // Weak content evidence
	}

	// Consistency Bonus (0-20%)
	if result.BestGender != "" && result.ContentGender != "" && result.BestGender == result.ContentGender {
		confidence += 0.1 // Gender consistency
	}
	if result.BestCategory != "" && result.ContentCategory != "" && result.BestCategory == result.ContentCategory {
		confidence += 0.1 // Category consistency
	}

	// Cap at 100%
	if confidence > 1.0 {
		confidence = 1.0
	}

	return confidence
}

// GetUnisexBrowseNodeIDs returns all supported Unisex Browse Node IDs
func GetUnisexBrowseNodeIDs() []string {
	mappings := getExtendedBrowseNodeMappings()
	var unisexNodes []string

	for id, mapping := range mappings {
		if mapping.Gender == "unisex" && mapping.IsSupported {
			unisexNodes = append(unisexNodes, id)
		}
	}

	return unisexNodes
}

// ValidateMultiNodeResult performs validation checks on analysis result
func ValidateMultiNodeResult(result MultiNodeAnalysisResult) []string {
	var issues []string

	// Check for empty results
	if result.BestGender == "" && result.BestCategory == "" {
		issues = append(issues, "No gender or category detected")
	}

	// Check confidence level
	if result.CombinedConfidence < 0.3 {
		issues = append(issues, "Low confidence score (<30%)")
	}

	// Check for unsupported nodes dominating
	if len(result.UnsupportedNodes) > len(result.SupportedNodes) {
		issues = append(issues, "More unsupported nodes than supported nodes")
	}

	// Check for promotion nodes
	for _, score := range result.NodeScores {
		if score.Priority == 6 && score.Score > 0 {
			issues = append(issues, "Promotion nodes affecting score")
		}
	}

	return issues
}

// ===== CONVENIENCE FUNCTIONS FOR INTEGRATION =====

// AnalyzeProductForGenderAndCategory is a convenience function for the Content Generation Worker
// It takes PA-API Browse Nodes and product content, returns best gender/category with confidence
func AnalyzeProductForGenderAndCategory(browseNodes []BrowseNode, title, description string) (gender, category string, confidence float64, debugInfo MultiNodeAnalysisResult) {
	// Combine title and description for comprehensive text analysis
	combinedText := title + " " + description

	// Perform full analysis
	result := AnalyzeMultipleBrowseNodes(browseNodes, combinedText)

	return result.BestGender, result.BestCategory, result.CombinedConfidence, result
}

// ConvertPAAPIBrowseNodes converts PA-API response format to our BrowseNode format
func ConvertPAAPIBrowseNodes(paapiBrowseNodes interface{}) []BrowseNode {
	var browseNodes []BrowseNode

	// Handle different input types (slice of maps, structs, etc.)
	switch nodes := paapiBrowseNodes.(type) {
	case []map[string]interface{}:
		for _, nodeMap := range nodes {
			node := &BrowseNodeData{}
			if id, ok := nodeMap["id"].(string); ok {
				node.ID = id
			}
			if displayName, ok := nodeMap["displayName"].(string); ok {
				node.DisplayName = displayName
			}
			if contextFreeName, ok := nodeMap["contextFreeName"].(string); ok {
				node.ContextFreeName = contextFreeName
			}
			// Handle ancestor if present
			if ancestor, ok := nodeMap["ancestor"].(map[string]interface{}); ok {
				if ancestorID, ok := ancestor["id"].(string); ok {
					node.Ancestor.ID = ancestorID
				}
				if ancestorDisplayName, ok := ancestor["displayName"].(string); ok {
					node.Ancestor.DisplayName = ancestorDisplayName
				}
				if ancestorContextFreeName, ok := ancestor["contextFreeName"].(string); ok {
					node.Ancestor.ContextFreeName = ancestorContextFreeName
				}
			}
			browseNodes = append(browseNodes, node)
		}
	case []BrowseNode:
		return nodes
	}

	return browseNodes
}

// GetSupportedBrowseNodeCount returns count of supported vs unsupported browse nodes
func GetSupportedBrowseNodeCount(browseNodes []BrowseNode) (supported, unsupported int) {
	mappings := getExtendedBrowseNodeMappings()

	for _, node := range browseNodes {
		if mapping, exists := mappings[node.GetID()]; exists {
			if mapping.IsSupported {
				supported++
			} else {
				unsupported++
			}
		} else {
			unsupported++ // Unknown nodes are considered unsupported
		}
	}

	return supported, unsupported
}

// IsUnisexProduct determines if a product should be classified as unisex based on browse nodes
func IsUnisexProduct(browseNodes []BrowseNode) bool {
	mappings := getExtendedBrowseNodeMappings()

	// Check for explicit unisex browse nodes
	for _, node := range browseNodes {
		if mapping, exists := mappings[node.GetID()]; exists {
			if mapping.Gender == "unisex" && mapping.IsSupported {
				return true
			}
		}
	}

	return false
}

// ===== BROWSE NODE MAPPING =====

// ExtendedBrowseNodeMapping represents comprehensive mapping for Browse Node IDs
type ExtendedBrowseNodeMapping struct {
	ID          string   `json:"id"`
	DisplayName string   `json:"display_name"`
	Category    string   `json:"category"`
	Gender      string   `json:"gender"`   // men, women, unisex
	Priority    int      `json:"priority"` // 1=most specific, 5=most general
	Tags        []string `json:"tags"`
	Type        string   `json:"type"`         // category, gender, style, promotion
	ParentID    string   `json:"parent_id"`    // For hierarchical structure
	IsSupported bool     `json:"is_supported"` // Whether we want to save products with this Browse Node ID
}

// getExtendedBrowseNodeMappings returns the comprehensive Browse Node mapping
func getExtendedBrowseNodeMappings() map[string]ExtendedBrowseNodeMapping {
	mappings := []ExtendedBrowseNodeMapping{
		// === UNISEX KATEGORIEN (Priority 1-4) ===
		{ID: "1981507031", DisplayName: "Unisex T-Shirts mit Sprüchen", Category: "T-Shirts", Gender: "unisex", Priority: 1, Tags: []string{"t-shirts", "unisex", "sayings", "slogans", "casual"}, Type: "category", ParentID: "1981397031", IsSupported: true},
		{ID: "78689031", DisplayName: "Bekleidung", Category: "Bekleidung", Gender: "unisex", Priority: 4, Tags: []string{"bekleidung", "mode"}, Type: "category", ParentID: "", IsSupported: true},
		{ID: "203879767031", DisplayName: "Fashion Kategorie (Allgemein)", Category: "Fashion", Gender: "unisex", Priority: 4, Tags: []string{"fashion", "mode", "allgemein"}, Type: "category", ParentID: "78689031", IsSupported: true},

		// === DAMEN KATEGORIEN (Priority 1-3) ===
		{ID: "1981297031", DisplayName: "T-Shirts", Category: "T-Shirts", Gender: "women", Priority: 1, Tags: []string{"t-shirts", "damen", "basic", "casual"}, Type: "category", ParentID: "1981292031", IsSupported: true},
		{ID: "1981305031", DisplayName: "T-Shirts Kurzarm", Category: "T-Shirts", Gender: "women", Priority: 1, Tags: []string{"t-shirts", "kurzarm", "damen", "basic", "sommer"}, Type: "category", ParentID: "1981297031", IsSupported: true},
		{ID: "1981872031", DisplayName: "T-Shirts für Damen (Spezifisch)", Category: "T-Shirts", Gender: "women", Priority: 1, Tags: []string{"t-shirts", "damen", "basic", "casual", "spezifisch"}, Type: "category", ParentID: "1981297031", IsSupported: true},
		{ID: "1981294031", DisplayName: "Langarmshirts", Category: "Langarmshirts", Gender: "women", Priority: 1, Tags: []string{"langarmshirts", "damen", "basic", "casual"}, Type: "category", ParentID: "1981292031", IsSupported: true},
		{ID: "1981296031", DisplayName: "Sweatshirts", Category: "Sweatshirts", Gender: "women", Priority: 1, Tags: []string{"sweatshirts", "damen", "casual", "comfort"}, Type: "category", ParentID: "1981292031", IsSupported: true},
		{ID: "1981283031", DisplayName: "Kapuzenpullover", Category: "Kapuzenpullover", Gender: "women", Priority: 1, Tags: []string{"kapuzenpullover", "hoodies", "damen", "casual", "comfort"}, Type: "category", ParentID: "1981282031", IsSupported: true},
		{ID: "1981292031", DisplayName: "Tops, T-Shirts & Blusen", Category: "Tops", Gender: "women", Priority: 2, Tags: []string{"tops", "damen", "oberkörper"}, Type: "category", ParentID: "1981206031", IsSupported: true},
		{ID: "1981206031", DisplayName: "Damen", Category: "Bekleidung", Gender: "women", Priority: 3, Tags: []string{"damen", "bekleidung"}, Type: "gender", ParentID: "78689031", IsSupported: true},

		// === HERREN KATEGORIEN (Priority 1-3) ===
		{ID: "1981397031", DisplayName: "T-Shirts", Category: "T-Shirts", Gender: "men", Priority: 1, Tags: []string{"t-shirts", "herren", "basic", "casual"}, Type: "category", ParentID: "1981394031", IsSupported: true},
		{ID: "1981407031", DisplayName: "T-Shirts Kurzarm", Category: "T-Shirts", Gender: "men", Priority: 1, Tags: []string{"t-shirts", "kurzarm", "herren", "basic", "sommer"}, Type: "category", ParentID: "1981397031", IsSupported: true},
		{ID: "1981396031", DisplayName: "Poloshirts", Category: "Poloshirts", Gender: "men", Priority: 1, Tags: []string{"poloshirts", "herren", "sport", "casual"}, Type: "category", ParentID: "1981394031", IsSupported: true},
		{ID: "16150668031", DisplayName: "Poloshirts für Herren", Category: "Poloshirts", Gender: "men", Priority: 1, Tags: []string{"poloshirts", "herren", "golf", "sport", "casual"}, Type: "category", ParentID: "1981397031", IsSupported: true},
		{ID: "1981388031", DisplayName: "Kapuzenpullover", Category: "Kapuzenpullover", Gender: "men", Priority: 1, Tags: []string{"kapuzenpullover", "hoodies", "herren", "casual", "comfort"}, Type: "category", ParentID: "1981387031", IsSupported: true},
		{ID: "1981393031", DisplayName: "Sweatshirts", Category: "Sweatshirts", Gender: "men", Priority: 1, Tags: []string{"sweatshirts", "herren", "casual", "comfort"}, Type: "category", ParentID: "1981387031", IsSupported: true},
		{ID: "1981394031", DisplayName: "Tops, T-Shirts & Hemden", Category: "Tops", Gender: "men", Priority: 2, Tags: []string{"tops", "herren", "oberkörper"}, Type: "category", ParentID: "1981208031", IsSupported: true},
		{ID: "1981208031", DisplayName: "Herren", Category: "Bekleidung", Gender: "men", Priority: 3, Tags: []string{"herren", "bekleidung"}, Type: "gender", ParentID: "78689031", IsSupported: true},

		// === NICHT-UNTERSTÜTZTE KATEGORIEN (Priority 5-6) ===
		{ID: "201170685031", DisplayName: "Sportbekleidung", Category: "Sportbekleidung", Gender: "unisex", Priority: 4, Tags: []string{"sportbekleidung", "activewear"}, Type: "category", ParentID: "", IsSupported: false},
		{ID: "3520847031", DisplayName: "Sports-Promotions", Category: "Promotion", Gender: "unisex", Priority: 6, Tags: []string{"promotion", "sport", "sale"}, Type: "promotion", ParentID: "", IsSupported: false},
		{ID: "29843001031", DisplayName: "Das einfache weiße T-Shirt", Category: "Promotion", Gender: "unisex", Priority: 6, Tags: []string{"promotion", "t-shirts", "weiß", "basic"}, Type: "promotion", ParentID: "", IsSupported: false},
	}

	// Convert to map for fast lookup
	result := make(map[string]ExtendedBrowseNodeMapping)
	for _, mapping := range mappings {
		result[mapping.ID] = mapping
	}

	return result
}

// ===== PATTERN MATCHING FUNCTIONS =====

// getGenderPatterns returns gender detection patterns
func getGenderPatterns() map[string][]PatternRule {
	return map[string][]PatternRule{
		"women": {
			{Patterns: []string{"damen", "frauen", "women", "female"}, Weight: 10, Language: "both"},
			{Patterns: []string{"lady", "ladies", "girl", "girls"}, Weight: 8, Language: "en"},
			{Patterns: []string{"feminine", "feminin"}, Weight: 6, Language: "both"},
		},
		"men": {
			{Patterns: []string{"herren", "männer", "men", "male"}, Weight: 10, Language: "both"},
			{Patterns: []string{"guy", "guys", "gentleman"}, Weight: 8, Language: "en"},
			{Patterns: []string{"masculine", "maskulin"}, Weight: 6, Language: "both"},
		},
		"unisex": {
			{Patterns: []string{"unisex", "universal", "both"}, Weight: 10, Language: "both"},
			{Patterns: []string{"gender neutral", "geschlechtsneutral"}, Weight: 8, Language: "both"},
			{Patterns: []string{"everyone", "alle", "für alle"}, Weight: 6, Language: "both"},
		},
	}
}

// getCategoryPatterns returns category detection patterns
func getCategoryPatterns() map[string][]PatternRule {
	return map[string][]PatternRule{
		"T-Shirts": {
			{Patterns: []string{"t-shirt", "tshirt", "t shirt"}, Weight: 10, Language: "both"},
			{Patterns: []string{"kurzarm", "short sleeve"}, Weight: 6, Language: "both"},
		},
		"Poloshirts": {
			{Patterns: []string{"polo", "poloshirt", "polo shirt"}, Weight: 10, Language: "both"},
		},
		"Hoodies": {
			{Patterns: []string{"hoodie", "kapuzenpullover", "hooded"}, Weight: 10, Language: "both"},
		},
		"Sweatshirts": {
			{Patterns: []string{"sweatshirt", "sweater"}, Weight: 10, Language: "both"},
		},
	}
}

// getTagPatterns returns tag detection patterns
func getTagPatterns() map[string][]PatternRule {
	return map[string][]PatternRule{
		"casual": {
			{Patterns: []string{"casual", "everyday", "lässig"}, Weight: 5, Language: "both"},
		},
		"sport": {
			{Patterns: []string{"sport", "athletic", "fitness"}, Weight: 5, Language: "both"},
		},
		"basic": {
			{Patterns: []string{"basic", "essential", "klassisch"}, Weight: 3, Language: "both"},
		},
	}
}

// calculateGenderScore calculates score for gender patterns
func calculateGenderScore(text string, patterns []PatternRule) int {
	textLower := strings.ToLower(text)
	totalScore := 0

	for _, rule := range patterns {
		for _, pattern := range rule.Patterns {
			if strings.Contains(textLower, strings.ToLower(pattern)) {
				totalScore += rule.Weight
			}
		}
	}

	return totalScore
}

// matchesPatternWithFuzzy checks if text matches patterns with fuzzy matching
func matchesPatternWithFuzzy(text string, patterns []PatternRule) bool {
	textLower := strings.ToLower(text)

	for _, rule := range patterns {
		for _, pattern := range rule.Patterns {
			if strings.Contains(textLower, strings.ToLower(pattern)) {
				return true
			}
		}
	}

	return false
}

// generateHierarchicalTags generates hierarchical tags from text
func generateHierarchicalTags(text string) []string {
	var tags []string
	textLower := strings.ToLower(text)

	// Basic clothing tags
	if strings.Contains(textLower, "shirt") || strings.Contains(textLower, "t-shirt") {
		tags = append(tags, "shirt", "top", "oberteil")
	}
	if strings.Contains(textLower, "casual") {
		tags = append(tags, "casual", "everyday")
	}
	if strings.Contains(textLower, "sport") {
		tags = append(tags, "sport", "active")
	}

	return tags
}
