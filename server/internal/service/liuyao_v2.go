package service

import (
	"embed"
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"sort"
	"strings"
	"time"
	"zhanbu/internal/model"

	"github.com/rs/zerolog/log"
)

//go:embed data/takashima_hexagrams.json data/takashima_interpretation_rules.json
var takashimaHexagramsJSON embed.FS

var (
	takashimaQuoteOnlyLineRE = regexp.MustCompile(`(?m)^\s*[“”"']+\s*$\n?`)
	takashimaDaYouOCRRE      = regexp.MustCompile(`曰\s*大\s*[“”"']?\s*有\s*。`)
	takashimaCJKSpaceRE      = regexp.MustCompile(`([\p{Han}])\s+([\p{Han}])`)
	takashimaCJKPunctRE      = regexp.MustCompile(`([\p{Han}])\s+([，。！？；：、])`)
	takashimaPunctCJKRE      = regexp.MustCompile(`([，。！？；：、])\s+([\p{Han}])`)
)

// timeNow returns current timestamp string.
func timeNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// lineSymbol returns the symbol for a line value.
func lineSymbol(value int) string {
	if value == 7 || value == 9 {
		return "⚊" // Yang
	}
	return "⚋" // Yin
}

// determineLineFromCount determines line result from yang count.
func determineLineFromCount(yangCount int) model.LineResult {
	switch yangCount {
	case 3: // All yang = Old Yang (9)
		return model.LineResult{Value: 9, Type: "old_yang", Mutable: true, Symbol: "⚊"}
	case 2: // Two yang = Young Yin (8)
		return model.LineResult{Value: 8, Type: "young_yin", Mutable: false, Symbol: "⚋"}
	case 1: // One yang = Young Yang (7)
		return model.LineResult{Value: 7, Type: "young_yang", Mutable: false, Symbol: "⚊"}
	default: // No yang = Old Yin (6)
		return model.LineResult{Value: 6, Type: "old_yin", Mutable: true, Symbol: "⚋"}
	}
}

// LiuYaoV2Service handles Takashima-style LiuYao divination.
type LiuYaoV2Service struct {
	hexagrams map[int]*model.TakashimaHexagram
	binaryIdx map[string]*model.TakashimaHexagram
	trigrams  []model.Trigram
	rng       *rand.Rand
	method    string // yarrow, coin, or both
}

// NewLiuYaoV2Service creates a new LiuYaoV2Service.
func NewLiuYaoV2Service(method string) (*LiuYaoV2Service, error) {
	svc := &LiuYaoV2Service{
		hexagrams: make(map[int]*model.TakashimaHexagram),
		binaryIdx: make(map[string]*model.TakashimaHexagram),
		rng:       rand.New(rand.NewSource(rand.Int63())),
		method:    method,
	}

	if err := svc.loadHexagrams(); err != nil {
		return nil, fmt.Errorf("failed to load hexagrams: %w", err)
	}

	if err := svc.loadTrigrams(); err != nil {
		return nil, fmt.Errorf("failed to load trigrams: %w", err)
	}

	return svc, nil
}

// loadHexagrams loads hexagram data from embedded JSON file.
func (s *LiuYaoV2Service) loadHexagrams() error {
	data, err := takashimaHexagramsJSON.ReadFile("data/takashima_hexagrams.json")
	if err != nil {
		return fmt.Errorf("failed to read takashima_hexagrams.json: %w", err)
	}

	var hexagrams []*model.TakashimaHexagram
	if err := json.Unmarshal(data, &hexagrams); err != nil {
		return fmt.Errorf("failed to unmarshal hexagrams: %w", err)
	}

	for _, h := range hexagrams {
		if h.NameShort == "" {
			h.NameShort = h.Name
		}
		s.hexagrams[h.ID] = h
		s.binaryIdx[h.Binary] = h
	}

	return nil
}

// loadTrigrams loads trigram data from embedded JSON file.
func (s *LiuYaoV2Service) loadTrigrams() error {
	data, err := trigramsJSON.ReadFile("data/trigrams.json")
	if err != nil {
		// Use fallback
		s.trigrams = []model.Trigram{
			{ID: 7, Name: "乾", Symbol: "☰", Binary: 7, Nature: "天", Element: "金", YinYang: "阳"},
			{ID: 6, Name: "兑", Symbol: "☱", Binary: 6, Nature: "泽", Element: "金", YinYang: "阴"},
			{ID: 5, Name: "离", Symbol: "☲", Binary: 5, Nature: "火", Element: "火", YinYang: "阴"},
			{ID: 4, Name: "震", Symbol: "☳", Binary: 4, Nature: "雷", Element: "木", YinYang: "阳"},
			{ID: 3, Name: "巽", Symbol: "☴", Binary: 3, Nature: "风", Element: "木", YinYang: "阴"},
			{ID: 2, Name: "坎", Symbol: "☵", Binary: 2, Nature: "水", Element: "水", YinYang: "阳"},
			{ID: 1, Name: "艮", Symbol: "☶", Binary: 1, Nature: "山", Element: "土", YinYang: "阳"},
			{ID: 0, Name: "坤", Symbol: "☷", Binary: 0, Nature: "地", Element: "土", YinYang: "阴"},
		}
		return nil
	}

	var trigrams []model.Trigram
	if err := json.Unmarshal(data, &trigrams); err != nil {
		return fmt.Errorf("failed to unmarshal trigrams: %w", err)
	}

	s.trigrams = trigrams
	return nil
}

// Throw performs divination using the configured method.
func (s *LiuYaoV2Service) Throw(question string, method string) (*model.LiuYaoV2Result, error) {
	// Use specified method or default
	throwMethod := method
	if throwMethod == "" {
		throwMethod = s.method
	}

	var lines [6]model.LineResult

	switch throwMethod {
	case "yarrow":
		// 蓍草揲蓍法
		for i := 0; i < 6; i++ {
			lines[i] = s.throwYarrow()
		}
	case "coin":
		// 铜钱法
		for i := 0; i < 6; i++ {
			lines[i] = s.throwCoin()
		}
	default:
		// 默认使用揲蓍法
		for i := 0; i < 6; i++ {
			lines[i] = s.throwYarrow()
		}
	}

	// Compute binary string (top-to-bottom reading)
	binary := s.linesToBinary(lines)

	// Lookup hexagram
	benGua, ok := s.binaryIdx[binary]
	if !ok {
		return nil, fmt.Errorf("hexagram not found for binary: %s", binary)
	}

	// Find mutable lines
	var mutableLines []int
	for i, line := range lines {
		if line.Mutable {
			mutableLines = append(mutableLines, i)
		}
	}

	// Compute bian gua if there are mutable lines
	var bianGua *model.TakashimaHexagram
	if len(mutableLines) > 0 {
		var flippedLines [6]model.LineResult
		for i, line := range lines {
			if line.Mutable {
				flippedLines[i] = s.flipLine(line)
			} else {
				flippedLines[i] = line
			}
		}
		bianBinary := s.linesToBinary(flippedLines)
		bianGua = s.binaryIdx[bianBinary]
	}

	result := &model.LiuYaoV2Result{
		Question:     question,
		Lines:        lines,
		BenGua:       benGua,
		BianGua:      bianGua,
		MutableLines: mutableLines,
		Method:       throwMethod,
		Timestamp:    timeNow(),
	}
	result.BookEvidence = s.BuildBookEvidence(result)
	return result, nil
}

// throwYarrow simulates yarrow stalk divination (揲蓍法).
// The yarrow stalk method produces probabilities:
// - Old Yang (9): 1/16 = 6.25%
// - Young Yin (8): 7/16 = 43.75%
// - Young Yang (7): 5/16 = 31.25%
// - Old Yin (6): 3/16 = 18.75%
func (s *LiuYaoV2Service) throwYarrow() model.LineResult {
	// Simulate the yarrow stalk method probabilities
	// Total: 16 parts
	r := s.rng.Intn(16)

	var value int
	var lineType string
	var mutable bool

	switch {
	case r == 0: // 1/16 = Old Yang
		value = 9
		lineType = "old_yang"
		mutable = true
	case r >= 1 && r <= 7: // 7/16 = Young Yin
		value = 8
		lineType = "young_yin"
		mutable = false
	case r >= 8 && r <= 12: // 5/16 = Young Yang
		value = 7
		lineType = "young_yang"
		mutable = false
	default: // 3/16 = Old Yin
		value = 6
		lineType = "old_yin"
		mutable = true
	}

	return model.LineResult{
		Value:   value,
		Type:    lineType,
		Mutable: mutable,
		Symbol:  lineSymbol(value),
	}
}

// throwCoin simulates coin divination (铜钱法).
// The coin method produces equal probabilities:
// - Old Yang (9): 1/8 = 12.5%
// - Young Yin (8): 3/8 = 37.5%
// - Young Yang (7): 3/8 = 37.5%
// - Old Yin (6): 1/8 = 12.5%
func (s *LiuYaoV2Service) throwCoin() model.LineResult {
	// Simulate 3 coins
	yangCount := 0
	for i := 0; i < 3; i++ {
		if s.rng.Intn(2) == 1 {
			yangCount++
		}
	}

	return determineLineFromCount(yangCount)
}

// linesToBinary converts 6 lines to a binary string.
func (s *LiuYaoV2Service) linesToBinary(lines [6]model.LineResult) string {
	binary := make([]byte, 6)
	for i, line := range lines {
		if line.Value == 7 || line.Value == 9 {
			binary[5-i] = '1' // Yang
		} else {
			binary[5-i] = '0' // Yin
		}
	}
	return string(binary)
}

// flipLine flips a mutable line for bian gua calculation.
func (s *LiuYaoV2Service) flipLine(line model.LineResult) model.LineResult {
	if line.Value == 9 { // Old Yang -> Young Yin
		return model.LineResult{Value: 8, Type: "young_yin", Mutable: false, Symbol: "⚋"}
	} else if line.Value == 6 { // Old Yin -> Young Yang
		return model.LineResult{Value: 7, Type: "young_yang", Mutable: false, Symbol: "⚊"}
	}
	return line
}

// GetHexagrams returns all hexagrams.
func (s *LiuYaoV2Service) GetHexagrams() []*model.TakashimaHexagram {
	result := make([]*model.TakashimaHexagram, 0, len(s.hexagrams))
	for _, h := range s.hexagrams {
		result = append(result, h)
	}
	return result
}

// GetHexagramByID returns a hexagram by ID.
func (s *LiuYaoV2Service) GetHexagramByID(id int) (*model.TakashimaHexagram, error) {
	h, ok := s.hexagrams[id]
	if !ok {
		return nil, fmt.Errorf("hexagram not found: %d", id)
	}
	return h, nil
}

// GetMethod returns the configured divination method.
func (s *LiuYaoV2Service) GetMethod() string {
	return s.method
}

// BuildBookEvidence searches the extracted Takashima corpus for compact,
// reading-specific excerpts. The AI prompt receives these snippets instead of
// the full book.
func (s *LiuYaoV2Service) BuildBookEvidence(result *model.LiuYaoV2Result) *model.TakashimaBookEvidence {
	if result == nil || result.BenGua == nil {
		return &model.TakashimaBookEvidence{}
	}

	terms := evidenceQueryTerms(result)
	snippets := make([]model.TakashimaEvidenceSnippet, 0, 8)
	snippets = append(snippets, searchHexagramEvidence(result.BenGua, terms, result.MutableLines, "本卦")...)
	if result.BianGua != nil {
		snippets = append(snippets, searchHexagramEvidence(result.BianGua, terms, nil, "变卦")...)
	}

	sort.SliceStable(snippets, func(i, j int) bool {
		return snippets[i].Score > snippets[j].Score
	})
	snippets = limitEvidenceChars(snippets, 5200)

	evidence := &model.TakashimaBookEvidence{
		QueryTerms:  terms,
		Snippets:    snippets,
		MethodRules: loadTakashimaMethodRules(3),
	}
	logTakashimaEvidenceSelection(result, evidence)
	return evidence
}

type takashimaEvidenceLogSummary struct {
	Kind        string `json:"kind"`
	Title       string `json:"title"`
	SourcePages []int  `json:"source_pages"`
	Score       int    `json:"score"`
	Preview     string `json:"preview"`
}

func logTakashimaEvidenceSelection(result *model.LiuYaoV2Result, evidence *model.TakashimaBookEvidence) {
	if result == nil || result.BenGua == nil || evidence == nil {
		return
	}

	bianGua := ""
	if result.BianGua != nil {
		bianGua = result.BianGua.FullName
	}

	log.Info().
		Str("component", "liuyao_v2").
		Str("ben_gua", result.BenGua.FullName).
		Str("bian_gua", bianGua).
		Ints("mutable_lines", result.MutableLines).
		Strs("query_terms", evidence.QueryTerms).
		Interface("selected_book_evidence", summarizeEvidenceForLog(evidence.Snippets)).
		Interface("selected_method_rules", summarizeEvidenceForLog(evidence.MethodRules)).
		Msg("selected Takashima book evidence for AI prompt")
}

func summarizeEvidenceForLog(snippets []model.TakashimaEvidenceSnippet) []takashimaEvidenceLogSummary {
	summaries := make([]takashimaEvidenceLogSummary, 0, len(snippets))
	for _, snippet := range snippets {
		summaries = append(summaries, takashimaEvidenceLogSummary{
			Kind:        snippet.Kind,
			Title:       snippet.Title,
			SourcePages: snippet.SourcePages,
			Score:       snippet.Score,
			Preview:     truncateLogPreview(cleanEvidenceText(snippet.Text), 120),
		})
	}
	return summaries
}

func truncateLogPreview(text string, max int) string {
	runes := []rune(text)
	if len(runes) <= max {
		return text
	}
	if max <= 1 {
		return "…"
	}
	return string(runes[:max-1]) + "…"
}

func evidenceQueryTerms(result *model.LiuYaoV2Result) []string {
	seen := map[string]bool{}
	add := func(term string, out *[]string) {
		term = strings.TrimSpace(term)
		if term == "" || seen[term] {
			return
		}
		seen[term] = true
		*out = append(*out, term)
	}

	terms := make([]string, 0, 12)
	add(result.BenGua.FullName, &terms)
	add(result.BenGua.Name, &terms)
	for _, idx := range result.MutableLines {
		line := lineByPosition(result.BenGua, idx+1)
		if line == nil {
			continue
		}
		add(line.Name, &terms)
		add(firstClause(line.Original), &terms)
	}
	if result.BianGua != nil {
		add(result.BianGua.FullName, &terms)
		add(result.BianGua.Name, &terms)
	}
	return terms
}

func searchHexagramEvidence(hexagram *model.TakashimaHexagram, terms []string, mutableLines []int, role string) []model.TakashimaEvidenceSnippet {
	if hexagram == nil {
		return nil
	}

	snippets := []model.TakashimaEvidenceSnippet{
		newSnippet("hexagram", role+" "+hexagram.FullName+" 卦辞", hexagram.Judgment.Text, []int{hexagram.Judgment.SourcePage}, scoreText(hexagram.Judgment.Text, terms)+20),
		newSnippet("hexagram", role+" "+hexagram.FullName+" 彖象", strings.Join(nonEmpty(hexagram.Tuan.Text, hexagram.Image.Text), "\n"), hexagram.Source.Pages, scoreText(hexagram.Tuan.Text+hexagram.Image.Text, terms)+10),
	}

	if len(mutableLines) == 0 && role == "本卦" {
		snippets = append(snippets, newSnippet("hexagram_raw", role+" "+hexagram.FullName+" 通论", hexagram.RawText, hexagram.Source.Pages, scoreText(hexagram.RawText, terms)))
	}

	if role == "本卦" {
		for _, idx := range mutableLines {
			line := lineByPosition(hexagram, idx+1)
			if line == nil {
				continue
			}
			text := strings.Join(nonEmpty(
				line.Name+"："+line.Original,
				"象传："+line.Commentary,
				line.TakashimaAnalysis,
			), "\n")
			snippets = append(snippets, newSnippet("line", role+" "+hexagram.FullName+" "+line.Name, text, line.SourcePages, scoreText(text, terms)+50))
		}
	}

	out := snippets[:0]
	for _, snippet := range snippets {
		if strings.TrimSpace(snippet.Text) != "" {
			out = append(out, snippet)
		}
	}
	return out
}

func newSnippet(kind string, title string, text string, pages []int, score int) model.TakashimaEvidenceSnippet {
	pages = compactPages(pages)
	return model.TakashimaEvidenceSnippet{
		Kind:        kind,
		Title:       title,
		Text:        sourcePrefix(pages) + truncateRunes(cleanEvidenceText(text), 1600),
		SourcePages: pages,
		Score:       score,
	}
}

func lineByPosition(hexagram *model.TakashimaHexagram, position int) *model.TakashimaLine {
	for i := range hexagram.Lines {
		if hexagram.Lines[i].Position == position {
			return &hexagram.Lines[i]
		}
	}
	return nil
}

func scoreText(text string, terms []string) int {
	score := 0
	for _, term := range terms {
		if term != "" && strings.Contains(text, term) {
			score += len([]rune(term))
		}
	}
	return score
}

func limitEvidenceChars(snippets []model.TakashimaEvidenceSnippet, maxChars int) []model.TakashimaEvidenceSnippet {
	out := make([]model.TakashimaEvidenceSnippet, 0, len(snippets))
	used := 0
	for _, snippet := range snippets {
		length := len([]rune(snippet.Text))
		if used+length > maxChars && len(out) > 0 {
			continue
		}
		out = append(out, snippet)
		used += length
	}
	return out
}

func loadTakashimaMethodRules(limit int) []model.TakashimaEvidenceSnippet {
	data, err := takashimaHexagramsJSON.ReadFile("data/takashima_interpretation_rules.json")
	if err != nil {
		return nil
	}
	var rules struct {
		InterpretationWorkflow []struct {
			Name        string `json:"name"`
			Rule        string `json:"rule"`
			SourcePages []int  `json:"source_pages"`
		} `json:"interpretation_workflow"`
	}
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil
	}
	if limit > len(rules.InterpretationWorkflow) {
		limit = len(rules.InterpretationWorkflow)
	}
	out := make([]model.TakashimaEvidenceSnippet, 0, limit)
	for i := 0; i < limit; i++ {
		rule := rules.InterpretationWorkflow[i]
		out = append(out, newSnippet("method_rule", rule.Name, rule.Rule, rule.SourcePages, 1))
	}
	return out
}

func sourcePrefix(pages []int) string {
	if len(pages) == 0 || pages[0] == 0 {
		return ""
	}
	if len(pages) == 1 {
		return fmt.Sprintf("【来源：第 %d 页】\n", pages[0])
	}
	return fmt.Sprintf("【来源：第 %d-%d 页】\n", pages[0], pages[len(pages)-1])
}

func compactPages(pages []int) []int {
	seen := map[int]bool{}
	out := make([]int, 0, len(pages))
	for _, page := range pages {
		if page <= 0 || seen[page] {
			continue
		}
		seen[page] = true
		out = append(out, page)
	}
	sort.Ints(out)
	return out
}

func firstClause(text string) string {
	text = strings.TrimSpace(text)
	for _, sep := range []string{"。", "，", ",", "；", ";", "\n"} {
		if idx := strings.Index(text, sep); idx > 0 {
			return strings.TrimSpace(text[:idx])
		}
	}
	return text
}

func nonEmpty(values ...string) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			out = append(out, value)
		}
	}
	return out
}

func cleanEvidenceText(text string) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	text = takashimaQuoteOnlyLineRE.ReplaceAllString(text, "")
	text = takashimaDaYouOCRRE.ReplaceAllString(text, "曰大有。")
	for {
		cleaned := takashimaCJKSpaceRE.ReplaceAllString(text, "$1$2")
		if cleaned == text {
			break
		}
		text = cleaned
	}
	text = takashimaCJKPunctRE.ReplaceAllString(text, "$1$2")
	text = takashimaPunctCJKRE.ReplaceAllString(text, "$1$2")
	lines := strings.FieldsFunc(text, func(r rune) bool { return r == '\n' || r == '\t' })
	return strings.Join(lines, "\n")
}

func truncateRunes(text string, max int) string {
	runes := []rune(text)
	if len(runes) <= max {
		return text
	}
	return string(runes[:max]) + "……"
}
