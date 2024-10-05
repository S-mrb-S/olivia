package olivia

import (
	"golang.org/x/oauth2"
	"time"
)

// =================================================================
type DataPacket struct {
	Label string `json:"tag"`

	Content []string `json:"messages"`
}

type UserProfile struct {
	FullName         string         `json:"name"`
	GenrePreferences []string       `json:"movie_genres"`
	DislikedMovies   []string       `json:"movie_blacklist"`
	ImportantDates   []UserReminder `json:"reminders"`
	StreamingToken   *oauth2.Token  `json:"spotify_token"`
	StreamingID      string         `json:"spotify_id"`
	StreamingSecret  string         `json:"spotify_secret"`
}

type UserReminder struct {
	ReminderDetails string `json:"reason"`
	ReminderDate    string `json:"date"`
}

type DashboardData struct {
	NetworkLayers NetworkLayersData `json:"layers"`
	TrainingInfo  TrainingInfoData  `json:"training"`
}

type NetworkLayersData struct {
	InputCount  int `json:"input"`
	HiddenCount int `json:"hidden"`
	OutputCount int `json:"output"`
}

type TrainingInfoData struct {
	LearningRate float64   `json:"rate"`
	ErrorMetrics []float64 `json:"errors"`
	TrainingTime float64   `json:"time"`
}

type clientRequestMessage struct {
	Type        int         `json:"type"` // 0 for handshakes and 1 for messages
	Content     string      `json:"content"`
	Token       string      `json:"user_token"`
	Locale      string      `json:"locale"`
	Information UserProfile `json:"information"`
}

type serverResponseMessage struct {
	Content     string      `json:"content"`
	Tag         string      `json:"tag"`
	Information UserProfile `json:"information"`
}

type LayerDerivative struct {
	Delta      Matrix
	Adjustment Matrix
}

type Matrix [][]float64

type Network struct {
	Layers  []Matrix
	Weights []Matrix
	Biases  []Matrix
	Output  Matrix
	Rate    float64
	Errors  []float64
	Time    float64
	Locale  string
}

type LocaleCoverage struct {
	Tag      string   `json:"locale_tag"`
	Language string   `json:"language"`
	Coverage Coverage `json:"coverage"`
}

type Coverage struct {
	Modules  CoverageDetails `json:"modules"`
	Intents  CoverageDetails `json:"intents"`
	Messages CoverageDetails `json:"messages"`
}

type CoverageDetails struct {
	NotCovered []string `json:"not_covered"`
	Coverage   int      `json:"coverage"`
}

type Intent struct {
	Tag       string   `json:"tag"`
	Patterns  []string `json:"patterns"`
	Responses []string `json:"responses"`
	Context   string   `json:"context"`
}

type Document struct {
	Sentence Sentence
	Tag      string
}

type Sentence struct {
	Locale  string
	Content string
}

type Result struct {
	Tag   string
	Value float64
}

type Error struct {
	Message string `json:"message"`
}

type DeleteRequest struct {
	Tag string `json:"tag"`
}

type Locale struct {
	Tag  string
	Name string
}

type Modulef struct {
	Tag       string
	Patterns  []string
	Responses []string
	Replacer  func(string, string, string, string) (string, string)
	Context   string
}

type Joke struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	Setup     string `json:"setup"`
	Punchline string `json:"punchline"`
}

type Module struct {
	Action func(string, string)
}

type ReasonKeyword struct {
	That string
	To   string
}

type SpotifyKeywords struct {
	Play string
	From string
	On   string
}

type Movie struct {
	Name   string
	Genres []string
	Rating float64
}

type Country struct {
	Name     map[string]string `json:"name"`
	Capital  string            `json:"capital"`
	Code     string            `json:"code"`
	Area     float64           `json:"area"`
	Currency string            `json:"currency"`
}

type Rule func(string, string) time.Time

type RuleTranslation struct {
	DaysOfWeek        []string
	Months            []string
	RuleToday         string
	RuleTomorrow      string
	RuleAfterTomorrow string
	RuleDayOfWeek     string
	RuleNextDayOfWeek string
	RuleNaturalDate   string
}

type PatternTranslations struct {
	DateRegex string
	TimeRegex string
}

// =================================================================
