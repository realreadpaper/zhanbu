package service

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"zhanbu/config"
	"zhanbu/internal/model"
)

func TestNewLiuYaoV2ServiceLoadsTakashimaCorpus(t *testing.T) {
	svc, err := NewLiuYaoV2Service("yarrow")
	require.NoError(t, err)

	hexagram, err := svc.GetHexagramByID(1)
	require.NoError(t, err)
	require.Equal(t, "乾", hexagram.Name)
	require.Equal(t, "乾为天", hexagram.FullName)
	require.Contains(t, hexagram.Judgment.Text, "元亨利贞")
	require.Len(t, hexagram.Lines, 7)
	require.Equal(t, "初九", hexagram.Lines[0].Name)
	require.Contains(t, hexagram.Lines[0].Original, "潜龙勿用")
	require.Contains(t, hexagram.Lines[0].TakashimaAnalysis, "潜龙勿用")
}

func TestLiuYaoV2ThrowAddsSearchBasedBookEvidence(t *testing.T) {
	svc, err := NewLiuYaoV2Service("coin")
	require.NoError(t, err)

	result := mustLiuYaoV2Result(t, svc, "事业是否该等待")
	result.MutableLines = []int{0}
	result.BookEvidence = svc.BuildBookEvidence(result)

	require.NotNil(t, result.BookEvidence)
	require.Contains(t, result.BookEvidence.QueryTerms, "乾为天")
	require.Contains(t, result.BookEvidence.QueryTerms, "初九")
	require.NotEmpty(t, result.BookEvidence.Snippets)

	joined := strings.Join(snippetTexts(result.BookEvidence.Snippets), "\n")
	require.Contains(t, joined, "潜龙勿用")
	require.Contains(t, joined, "第 58 页")
	require.LessOrEqual(t, len(joined), 7000)
}

func TestTakashimaEvidenceLogSummariesExposeSelectedBookParts(t *testing.T) {
	svc, err := NewLiuYaoV2Service("coin")
	require.NoError(t, err)

	result := mustLiuYaoV2Result(t, svc, "事业是否该等待")
	result.BookEvidence = svc.BuildBookEvidence(result)

	summaries := summarizeEvidenceForLog(result.BookEvidence.Snippets)
	require.NotEmpty(t, summaries)
	require.NotEmpty(t, summaries[0].Title)
	require.NotEmpty(t, summaries[0].Preview)
	require.NotEmpty(t, summaries[0].SourcePages)
	require.LessOrEqual(t, len([]rune(summaries[0].Preview)), 120)

	joined := strings.Join(logPreviews(summaries), "\n")
	require.Contains(t, joined, "潜龙勿用")
	require.Contains(t, joined, "第 58 页")
}

func TestTakashimaCorpusHasCleanCanonicalText(t *testing.T) {
	svc, err := NewLiuYaoV2Service("coin")
	require.NoError(t, err)

	hexagrams := svc.GetHexagrams()
	require.Len(t, hexagrams, 64)

	for _, hexagram := range hexagrams {
		require.NotEmpty(t, hexagram.Tuan.Text, "%s should include the 彖传/卦象总论 text", hexagram.FullName)
		require.NotEmpty(t, hexagram.Image.Text, "%s should include the 大象/象辞 text", hexagram.FullName)
		require.NotRegexp(t, `[\p{Han}] [\p{Han}]`, hexagram.Judgment.Text, "%s judgment should not contain OCR spaces between Chinese characters", hexagram.FullName)
		require.NotRegexp(t, `[\p{Han}] [\p{Han}]`, hexagram.Tuan.Text, "%s tuan should not contain OCR spaces between Chinese characters", hexagram.FullName)
		require.NotRegexp(t, `[\p{Han}] [\p{Han}]`, hexagram.Image.Text, "%s image should not contain OCR spaces between Chinese characters", hexagram.FullName)
		for _, line := range hexagram.Lines {
			require.NotRegexp(t, `^[:：；;]`, line.Commentary, "%s %s commentary should not start with a stray separator", hexagram.FullName, line.Name)
			require.NotRegexp(t, `[\p{Han}] [\p{Han}]`, line.Original, "%s %s original should not contain OCR spaces between Chinese characters", hexagram.FullName, line.Name)
			require.NotRegexp(t, `[\p{Han}] [\p{Han}]`, line.Commentary, "%s %s commentary should not contain OCR spaces between Chinese characters", hexagram.FullName, line.Name)
			require.NotContains(t, line.TakashimaAnalysis, "象传：:", "%s %s analysis should not contain duplicated 象传 separators", hexagram.FullName, line.Name)
		}
	}

	jian, err := svc.GetHexagramByID(39)
	require.NoError(t, err)
	require.Contains(t, jian.Image.Text, "山上有水，蹇，君子以反身修德")

	dayou, err := svc.GetHexagramByID(14)
	require.NoError(t, err)
	require.Contains(t, dayou.Tuan.Text, "曰大有。其德刚健而文明")
	require.NotContains(t, dayou.Tuan.Text, "曰 大")
}

func TestCleanEvidenceTextRemovesPDFArtifacts(t *testing.T) {
	raw := "《大象》曰：天地交泰，后以财 成天地之\n                     “\n宜，以左右民。\n《彖传》曰：《大有》，柔得尊位，大中而上下应之，曰 大 “\n有 。其德刚健而文明。\n ”"
	cleaned := cleanEvidenceText(raw)

	require.Contains(t, cleaned, "财成天地之宜")
	require.Contains(t, cleaned, "曰大有。其德刚健而文明")
	require.NotContains(t, cleaned, "曰 大")
	require.NotRegexp(t, `[\p{Han}] [\p{Han}]`, cleaned)
}

func TestOpenAIProviderLoadsLiuYaoV2Prompt(t *testing.T) {
	provider, err := NewOpenAIProvider(&config.AIConfig{
		APIKey:      "sk-test",
		BaseURL:     "https://example.test/v1",
		Model:       "test-model",
		MaxTokens:   1000,
		Temperature: 0.7,
	})
	require.NoError(t, err)
	require.Contains(t, provider.prompts, "liuyao_v2")
}

func logPreviews(summaries []takashimaEvidenceLogSummary) []string {
	out := make([]string, 0, len(summaries))
	for _, summary := range summaries {
		out = append(out, summary.Preview)
	}
	return out
}

func mustLiuYaoV2Result(t *testing.T, svc *LiuYaoV2Service, question string) *model.LiuYaoV2Result {
	t.Helper()
	ben, err := svc.GetHexagramByID(1)
	require.NoError(t, err)
	bian, err := svc.GetHexagramByID(44)
	require.NoError(t, err)
	return &model.LiuYaoV2Result{
		Question:     question,
		BenGua:       ben,
		BianGua:      bian,
		MutableLines: []int{0},
		Method:       "coin",
	}
}

func snippetTexts(snippets []model.TakashimaEvidenceSnippet) []string {
	out := make([]string, 0, len(snippets))
	for _, snippet := range snippets {
		out = append(out, snippet.Text)
	}
	return out
}
