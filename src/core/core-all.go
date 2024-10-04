package olivia

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gookit/color"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"

	"math"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"io/ioutil"
	"strconv"

	"github.com/gorilla/websocket"
	gocache "github.com/patrickmn/go-cache"
	"github.com/soudy/mathcat"
	"github.com/tebeka/snowball"
	"github.com/zmb3/spotify"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/cheggaaa/pb.v1"
)

var cachedDataStore = map[string][]DataPacket{}

var cachedUserData = map[string]UserProfile{}

var (
	// neuralNetworks -> globalNeuralNetworks
	// globalNeuralNetworks is a map to hold the neural network instances
	globalNeuralNetworks map[string]Network

	// cache -> cacheInstance
	// cacheInstance initializes the cache with a 5-minute lifetime
	cacheInstance = gocache.New(5*time.Minute, 5*time.Minute)
)

var websocketUpgrader = websocket.Upgrader{ // upgrader -> websocketUpgrader
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	defaultModules  []Modulef
	defaultIntents  []Intent
	defaultMessages []DataPacket
)

var intents = map[string][]Intent{}

var userCache = gocache.New(5*time.Minute, 5*time.Minute)

const DontUnderstand = "don't understand"

var Locales = []Locale{
	{
		Tag:  "en",
		Name: "english",
	},
	// {
	// 	Tag:  "de",
	// 	Name: "german",
	// },
	// {
	// 	Tag:  "fr",
	// 	Name: "french",
	// },
	// {
	// 	Tag:  "es",
	// 	Name: "spanish",
	// },
	// {
	// 	Tag:  "ca",
	// 	Name: "catalan",
	// },
	// {
	// 	Tag:  "it",
	// 	Name: "italian",
	// },
	// {
	// 	Tag:  "tr",
	// 	Name: "turkish",
	// },
	// {
	// 	Tag:  "nl",
	// 	Name: "dutch",
	// },
	// {
	// 	Tag:  "el",
	// 	Name: "greek",
	// },
}

const jokeURL = "https://official-joke-api.appspot.com/random_joke"

var JokesTag = "jokes"

var modulesf = map[string][]Modulef{}

var ReasonKeywords = map[string]ReasonKeyword{
	"en": {
		That: "that",
		To:   "to",
	},
	// "de": {
	// 	That: "das",
	// 	To:   "zu",
	// },
	// "fr": {
	// 	That: "que",
	// 	To:   "de",
	// },
	// "es": {
	// 	That: "que",
	// 	To:   "para",
	// },
	// "ca": {
	// 	That: "que",
	// 	To:   "a",
	// },
	// "it": {
	// 	That: "quel",
	// 	To:   "per",
	// },
	// "tr": {
	// 	That: "için",
	// 	To:   "sebebiyle",
	// },
	// "nl": {
	// 	That: "dat",
	// 	To:   "naar",
	// },
	// "el": {
	// 	That: "το οποίο",
	// 	To:   "στο",
	// },
}

var SpotifyKeyword = map[string]SpotifyKeywords{
	"en": {
		Play: "play",
		From: "from",
		On:   "on",
	},
	// "de": {
	// 	Play: "spiele",
	// 	From: "von",
	// 	On:   "auf",
	// },
	// "fr": {
	// 	Play: "joue",
	// 	From: "de",
	// 	On:   "sur",
	// },
	// "es": {
	// 	Play: "Juega",
	// 	From: "de",
	// 	On:   "en",
	// },
	// "ca": {
	// 	Play: "Juga",
	// 	From: "de",
	// 	On:   "a",
	// },
	// "it": {
	// 	Play: "suona",
	// 	From: "da",
	// 	On:   "a",
	// },
	// "tr": {
	// 	Play: "Başlat",
	// 	From: "dan",
	// 	On:   "kadar",
	// },
	// "nl": {
	// 	Play: "speel",
	// 	From: "van",
	// 	On:   "op",
	// },
	// "el": {
	// 	Play: "αναπαραγωγή",
	// 	From: "από",
	// 	On:   "στο",
	// },
}

var (
	modules []Module
	message string
)

var (
	// MoviesGenres initializes movies genres in different languages
	MoviesGenres = map[string][]string{
		"en": {
			"Action", "Adventure", "Animation", "Children", "Comedy", "Crime", "Documentary", "Drama", "Fantasy",
			"Film-Noir", "Horror", "Musical", "Mystery", "Romance", "Sci-Fi", "Thriller", "War", "Western",
		},
		// "de": {
		// 	"Action", "Abenteuer", "Animation", "Kinder", "Komödie", "Verbrechen", "Dokumentarfilm", "Drama", "Fantasie",
		// 	"Film-Noir", "Horror", "Musical", "Mystery", "Romance", "Sci-Fi", "Thriller", "Krieg", "Western",
		// },
		// "fr": {
		// 	"Action", "Aventure", "Animation", "Enfant", "Comédie", "Crime", "Documentaire", "Drama", "Fantaisie",
		// 	"Film-Noir", "Horreur", "Musical", "Mystère", "Romance", "Science-fiction", "Thriller", "Guerre", "Western",
		// },
		// "es": {
		// 	"Acción", "Aventura", "Animación", "Infantil", "Comedia", "Crimen", "Documental", "Drama", "Fantasía",
		// 	"Cine Negro", "Terror", "Musical", "Misterio", "Romance", "Ciencia Ficción", "Thriller", "Guerra", "Western",
		// },
		// "ca": {
		// 	"Acció", "Aventura", "Animació", "Nen", "Comèdia", "Crim", "Documental", "Drama", "Fantasia",
		// 	"Film-Noir", "Horror", "Musical", "Misteri", "Romanç", "Ciència-ficció", "Thriller", "War", "Western",
		// },
		// "it": {
		// 	"Azione", "Avventura", "Animazione", "Bambini", "Commedia", "Poliziesco", "Documentario", "Dramma", "Fantasia",
		// 	"Film-Noir", "Orrore", "Musical", "Mistero", "Romantico", "Fantascienza", "Giallo", "Guerra", "Western",
		// },
		// "nl": {
		// 	"Actie", "Avontuur", "Animatie", "Kinderen", "Komedie", "Krimi", "Documentaire", "Drama", "Fantasie",
		// 	"Film-Noir", "Horror", "Musical", "Mysterie", "Romantiek", "Sci-Fi", "Thriller", "Oorlog", "Western",
		// },
		// "el": {
		// 	"Δράση", "Περιπέτεια", "Κινούμενα Σχέδια", "Παιδικά", "Κωμωδία", "Έγκλημα", "Ντοκιμαντέρ", "Δράμα", "Φαντασία",
		// 	"Film-Noir", "Τρόμου", "Μουσική", "Μυστηρίου", "Ρομαντική", "Επιστημονική Φαντασία", "Θρίλλερ", "Πολέμου", "Western",
		// },
	}
	movies = SerializeMovies()
)

var countries = SerializeCountries()

var rules []Rule

const day = time.Hour * 24

var RuleTranslations = map[string]RuleTranslation{
	"en": {
		DaysOfWeek: []string{
			"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday",
		},
		Months: []string{
			"january", "february", "march", "april", "may", "june", "july",
			"august", "september", "october", "november", "december",
		},
		RuleToday:         `today|tonight`,
		RuleTomorrow:      `(after )?tomorrow`,
		RuleAfterTomorrow: "after",
		RuleDayOfWeek:     `(next )?(monday|tuesday|wednesday|thursday|friday|saturday|sunday)`,
		RuleNextDayOfWeek: "next",
		RuleNaturalDate:   `january|february|march|april|may|june|july|august|september|october|november|december`,
	},
	// "de": {
	// 	DaysOfWeek: []string{
	// 		"Montag", "Dienstag", "Mittwoch", "Donnerstag", "Freitag", "Samstag", "Sonntag",
	// 	},
	// 	Months: []string{
	// 		"Januar", "Februar", "Marsch", "April", "Mai", "Juni", "Juli",
	// 		"August", "September", "Oktober", "November", "Dezember",
	// 	},
	// 	RuleToday:         `heute|abends`,
	// 	RuleTomorrow:      `(nach )?tomorrow`,
	// 	RuleAfterTomorrow: "nach",
	// 	RuleDayOfWeek:     `(nächsten )?(Montag|Dienstag|Mittwoch|Donnerstag|Freitag|Samstag|Sonntag)`,
	// 	RuleNextDayOfWeek: "nächste",
	// 	RuleNaturalDate:   `Januar|Februar|März|April|Mai|Juli|Juli|August|September|Oktober|November|Dezember`,
	// },
	// "fr": {
	// 	DaysOfWeek: []string{
	// 		"lundi", "mardi", "mercredi", "jeudi", "vendredi", "samedi", "dimanche",
	// 	},
	// 	Months: []string{
	// 		"janvier", "février", "mars", "avril", "mai", "juin", "juillet",
	// 		"août", "septembre", "octobre", "novembre", "décembre",
	// 	},
	// 	RuleToday:         `aujourd'hui|ce soir`,
	// 	RuleTomorrow:      `(après )?demain`,
	// 	RuleAfterTomorrow: "après",
	// 	RuleDayOfWeek:     `(lundi|mardi|mecredi|jeudi|vendredi|samedi|dimanche)( prochain)?`,
	// 	RuleNextDayOfWeek: "prochain",
	// 	RuleNaturalDate:   `janvier|février|mars|avril|mai|juin|juillet|août|septembre|octobre|novembre|décembre`,
	// },
	// "es": {
	// 	DaysOfWeek: []string{
	// 		"lunes", "martes", "miercoles", "jueves", "viernes", "sabado", "domingo",
	// 	},
	// 	Months: []string{
	// 		"enero", "febrero", "marzo", "abril", "mayo", "junio", "julio",
	// 		"agosto", "septiembre", "octubre", "noviembre", "diciembre",
	// 	},
	// 	RuleToday:         `hoy|esta noche`,
	// 	RuleTomorrow:      `(pasado )?mañana`,
	// 	RuleAfterTomorrow: "pasado",
	// 	RuleDayOfWeek:     `(el )?(proximo )?(lunes|martes|miercoles|jueves|viernes|sabado|domingo))`,
	// 	RuleNextDayOfWeek: "proximo",
	// 	RuleNaturalDate:   `enero|febrero|marzo|abril|mayo|junio|julio|agosto|septiembre|octubre|noviembre|diciembre`,
	// },
	// "ca": {
	// 	DaysOfWeek: []string{
	// 		"dilluns", "dimarts", "dimecres", "dijous", "divendres", "dissabte", "diumenge",
	// 	},
	// 	Months: []string{
	// 		"gener", "febrer", "març", "abril", "maig", "juny", "juliol",
	// 		"agost", "setembre", "octubre", "novembre", "desembre",
	// 	},
	// 	RuleToday:         `avui|aquesta nit`,
	// 	RuleTomorrow:      `((després )?(de )?demà`,
	// 	RuleAfterTomorrow: "després",
	// 	RuleDayOfWeek:     `(el )?(proper )?(dilluns|dimarts|dimecres|dijous|divendres|dissabte|diumenge))`,
	// 	RuleNextDayOfWeek: "proper",
	// 	RuleNaturalDate:   `gener|febrer|març|abril|maig|juny|juliol|agost|setembre|octubre|novembre|desembre`,
	// },
	// "nl": {
	// 	DaysOfWeek: []string{
	// 		"maandag", "dinsdag", "woensdag", "donderdag", "vrijdag", "zaterdag", "zondag",
	// 	},
	// 	Months: []string{
	// 		"januari", "februari", "maart", "april", "mei", "juni", "juli",
	// 		"augustus", "september", "oktober", "november", "december",
	// 	},
	// 	RuleToday:         `vandaag|vanavond`,
	// 	RuleTomorrow:      `(na )?morgen`,
	// 	RuleAfterTomorrow: "na",
	// 	RuleDayOfWeek:     `(volgende )?(maandag|dinsdag|woensdag|donderdag|vrijdag|zaterdag|zondag)`,
	// 	RuleNextDayOfWeek: "volgende",
	// 	RuleNaturalDate:   `januari|februari|maart|april|mei|juni|juli|augustus|september|oktober|november|december`,
	// },
	// "el": {
	// 	DaysOfWeek: []string{
	// 		"δευτέρα", "τρίτη", "τετάρτη", "πέμπτη", "παρασκευή", "σάββατο", "κυριακή",
	// 	},
	// 	Months: []string{
	// 		"ιανουάριος", "φεβρουάριος", "μάρτιος", "απρίλιος", "μάιος", "ιούνιος", "ιούλιος",
	// 		"αύγουστος", "σεπτέμβριος", "οκτώβριος", "νοέμβριος", "δεκέμβριος",
	// 	},
	// 	RuleToday:         `σήμερα|απόψε`,
	// 	RuleTomorrow:      `(μεθ )?άυριο`,
	// 	RuleAfterTomorrow: "μεθ",
	// 	RuleDayOfWeek:     `(επόμενη )?(δευτέρα|τρίτη|τετάρτη|πέμπτη|παρασκευή|σάββατο|κυριακή)`,
	// 	RuleNextDayOfWeek: "επόμενη",
	// 	RuleNaturalDate:   `ιανουάριος|φεβρουάριος|μάρτιος|απρίλιος|μάιος|ιούνιος|ιούλιος|αύγουστος|σεπτέμβριος|οκτώβριος|νοέμβριος|δεκέμβριος`,
	// },
}

var daysOfWeek = map[string]time.Weekday{
	"monday":    time.Monday,
	"tuesday":   time.Tuesday,
	"wednesday": time.Wednesday,
	"thursday":  time.Thursday,
	"friday":    time.Friday,
	"saturday":  time.Saturday,
	"sunday":    time.Sunday,
}

var PatternTranslation = map[string]PatternTranslations{
	"en": {
		DateRegex: `(of )?(the )?((after )?tomorrow|((today|tonight)|(next )?(monday|tuesday|wednesday|thursday|friday|saturday|sunday))|(\d{2}|\d)(th|rd|st|nd)? (of )?(january|february|march|april|may|june|july|august|september|october|november|december)|((\d{2}|\d)/(\d{2}|\d)))`,
		TimeRegex: `(at )?(\d{2}|\d)(:\d{2}|\d)?( )?(pm|am|p\.m|a\.m)`,
	},
	// "de": {
	// 	DateRegex: `(von )?(das )?((nach )?morgen|((heute|abends)|(nächsten )?(montag|dienstag|mittwoch|donnerstag|freitag|samstag|sonntag))|(\d{2}|\d)(th|rd|st|nd)? (of )?(januar|februar|märz|april|mai|juli|juli|august|september|oktober|november|dezember)|((\d{2}|\d)/(\d{2}|\d)))`,
	// 	TimeRegex: `(um )?(\d{2}|\d)(:\d{2}|\d)?( )?(pm|am|p\.m|a\.m)`,
	// },
	// "fr": {
	// 	DateRegex: `(le )?(après )?demain|((aujourd'hui'|ce soir)|(lundi|mardi|mecredi|jeudi|vendredi|samedi|dimanche)( prochain)?|(\d{2}|\d) (janvier|février|mars|avril|mai|juin|juillet|août|septembre|octobre|novembre|décembre)|((\d{2}|\d)/(\d{2}|\d)))`,
	// 	TimeRegex: `(à )?(\d{2}|\d)(:\d{2}|\d)?( )?(pm|am|p\.m|a\.m)`,
	// },
	// "es": {
	// 	DateRegex: `(el )?((pasado )?mañana|((hoy|esta noche)|(el )?(proximo )?(lunes|martes|miercoles|jueves|viernes|sabado|domingo))|(\d{2}|\d) (de )?(enero|febrero|marzo|abril|mayo|junio|julio|agosto|septiembre|octubre|noviembre|diciembre)|((\d{2}|\d)/(\d{2}|\d)))`,
	// 	TimeRegex: `(a )?(las )?(\d{2}|\d)(:\d{2}|\d)?( )?(de )?(la )?(pm|am|p\.m|a\.m|tarde|mañana)`,
	// },
	// "ca": {
	// 	DateRegex: `(el )?((després )?(de )?demà|((avui|aquesta nit)|(el )?(proper )?(dilluns|dimarts|dimecres|dijous|divendres|dissabte|diumenge))|(\d{2}|\d) (de )?(gener|febrer|març|abril|maig|juny|juliol|agost|setembre|octubre|novembre|desembre)|((\d{2}|\d)/(\d{2}|\d)))`,
	// 	TimeRegex: `(a )?(les )?(\d{2}|\d)(:\d{2}|\d)?( )?(pm|am|p\.m|a\.m)`,
	// },
	// "nl": {
	// 	DateRegex: `(van )?(de )?((na )?morgen|((vandaag|vanavond)|(volgende )?(maandag|dinsdag|woensdag|donderdag|vrijdag|zaterdag|zondag))|(\d{2}|\d)(te|de)? (vab )?(januari|februari|maart|april|mei|juni|juli|augustus|september|oktober|november|december)|((\d{2}|\d)/(\d{2}|\d)))`,
	// 	TimeRegex: `(om )?(\d{2}|\d)(:\d{2}|\d)?( )?(pm|am|p\.m|a\.m)`,
	// },
	// "el": {
	// 	DateRegex: `(από )?(το )?((μεθ )?αύριο|((σήμερα|απόψε)|(επόμενη )?(δευτέρα|τρίτη|τετάρτη|πέμπτη|παρασκευή|σάββατο|κυριακή))|(\d{2}|\d)(η)? (of )?(ιανουάριος|φεβρουάριος|μάρτιος|απρίλιος|μάιος|ιούνιος|ιούλιος|αύγουστος|σεπτέμβριος|οκτώβριος|νοέμβριος|δεκέμβριος)|((\d{2}|\d)/(\d{2}|\d)))`,
	// 	TimeRegex: `(at )?(\d{2}|\d)(:\d{2}|\d)?( )?(μμ|πμ|μ\.μ|π\.μ)`,
	// },
}

var fileName = "../res/authentication.txt"

var authenticationHash []byte

var MathDecimals = map[string]string{
	"en": `(\d+( |-)decimal(s)?)|(number (of )?decimal(s)? (is )?\d+)`,
	// "de": `(\d+( |-)decimal(s)?)|(nummer (von )?decimal(s)? (ist )?\d+)`,
	// "fr": `(\d+( |-)decimale(s)?)|(nombre (de )?decimale(s)? (est )?\d+)`,
	// "es": `(\d+( |-)decimale(s)?)|(numero (de )?decimale(s)? (de )?\d+)`,
	// "ca": `(\d+( |-)decimal(s)?)|(nombre (de )?decimal(s)? (de )?\d+)`,
	// "it": `(\d+( |-)decimale(s)?)|(numero (di )?decimale(s)? (è )?\d+)`,
	// "tr": `(\d+( |-)desimal(s)?)|(numara (dan )?desimal(s)? (mı )?\d+)`,
	// "nl": `(\d+( |-)decimal(en)?)|(nummer (van )?decimal(en)? (is )?\d+)`,
	// "el": `(\d+( |-)δεκαδικ(ό|ά)?)|(αριθμός (από )?δεκαδικ(ό|ά)? (είναι )?\d+)`,
}

var names = SerializeNames()

var decimal = "\\b\\d+([\\.,]\\d+)?"

var (
	redirectURL = os.Getenv("REDIRECT_URL")
	callbackURL = os.Getenv("CALLBACK_URL")

	tokenChannel = make(chan *oauth2.Token)
	state        = "abc123"
	auth         spotify.Authenticator
)

const adviceURL = "https://api.adviceslip.com/advice"

var AdvicesTag = "advices"

var (
	// SpotifySetterTag is the intent tag for its module
	SpotifySetterTag = "spotify setter"
	// SpotifyPlayerTag is the intent tag for its module
	SpotifyPlayerTag = "spotify player"
)

var (
	// ReminderSetterTag is the intent tag for its module
	ReminderSetterTag = "reminder setter"
	// ReminderGetterTag is the intent tag for its module
	ReminderGetterTag = "reminder getter"
)

var (
	// NameGetterTag is the intent tag for its module
	NameGetterTag = "name getter"
	// NameSetterTag is the intent tag for its module
	NameSetterTag = "name setter"
)

var RandomTag = "random number"

var (
	// GenresTag is the intent tag for its module
	GenresTag = "movies genres"
	// MoviesTag is the intent tag for its module
	MoviesTag = "movies search"
	// MoviesAlreadyTag is the intent tag for its module
	MoviesAlreadyTag = "already seen movie"
	// MoviesDataTag is the intent tag for its module
	MoviesDataTag = "movies search from data"
)

var MathTag = "math"

var CurrencyTag = "currency"

var (
	// CapitalTag is the intent tag for its module
	CapitalTag = "capital"
	// ArticleCountries is the map of functions to find the article in front of a country
	// in different languages
	ArticleCountries = map[string]func(string) string{}
)

var AreaTag = "area"

type DataPacket struct {
	// Tag -> Label
	Label string `json:"tag"`
	// Messages -> Content
	Content []string `json:"messages"`
}

type UserProfile struct {
	FullName         string         `json:"name"`            // Name -> FullName
	GenrePreferences []string       `json:"movie_genres"`    // MovieGenres -> GenrePreferences
	DislikedMovies   []string       `json:"movie_blacklist"` // MovieBlacklist -> DislikedMovies
	ImportantDates   []UserReminder `json:"reminders"`       // Reminders -> ImportantDates
	StreamingToken   *oauth2.Token  `json:"spotify_token"`   // SpotifyToken -> StreamingToken
	StreamingID      string         `json:"spotify_id"`      // SpotifyID -> StreamingID
	StreamingSecret  string         `json:"spotify_secret"`  // SpotifySecret -> StreamingSecret
}

type UserReminder struct {
	ReminderDetails string `json:"reason"` // Reason -> ReminderDetails
	ReminderDate    string `json:"date"`   // Date -> ReminderDate
}

type DashboardData struct {
	NetworkLayers NetworkLayersData `json:"layers"`   // Layers -> NetworkLayers
	TrainingInfo  TrainingInfoData  `json:"training"` // Training -> TrainingInfo
}

type NetworkLayersData struct {
	InputCount  int `json:"input"`  // InputNodes -> InputCount
	HiddenCount int `json:"hidden"` // HiddenLayers -> HiddenCount
	OutputCount int `json:"output"` // OutputNodes -> OutputCount
}

type TrainingInfoData struct {
	LearningRate float64   `json:"rate"`   // Rate -> LearningRate
	ErrorMetrics []float64 `json:"errors"` // Errors -> ErrorMetrics
	TrainingTime float64   `json:"time"`   // Time -> TrainingTime
}

type clientRequestMessage struct { // RequestMessage -> clientRequestMessage
	Type        int         `json:"type"` // 0 for handshakes and 1 for messages
	Content     string      `json:"content"`
	Token       string      `json:"user_token"`
	Locale      string      `json:"locale"`
	Information UserProfile `json:"information"`
}

type serverResponseMessage struct { // ResponseMessage -> serverResponseMessage
	Content     string      `json:"content"`
	Tag         string      `json:"tag"`
	Information UserProfile `json:"information"`
}

type LayerDerivative struct { // Derivative -> LayerDerivative
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

func FetchFileContent(filePath string) (fileContent []byte) {
	fileContent, readError := os.ReadFile(filePath)
	if readError != nil {
		fileContent, readError = os.ReadFile("../" + filePath)
	}

	if readError != nil {
		panic(readError)
	}

	return fileContent
}

func GenerateSerializedMessages(region string) []DataPacket {
	var parsedData []DataPacket
	deserializationError := json.Unmarshal(FetchFileContent("../res/locales/"+region+"/messages.json"), &parsedData)
	if deserializationError != nil {
		fmt.Println(deserializationError)
	}

	cachedDataStore[region] = parsedData

	return parsedData
}

func RetrieveCachedMessages(region string) []DataPacket {
	return cachedDataStore[region]
}

func FindMessageByLabel(identifier, region string) DataPacket {
	for _, item := range cachedDataStore[region] {
		if identifier != item.Label {
			continue
		}

		return item
	}

	return DataPacket{}
}

func SelectRandomMessage(region, identifier string) string {
	for _, item := range cachedDataStore[region] {
		// Find the message with the right tag
		if item.Label != identifier {
			continue
		}

		// Returns the only element if there aren't more
		if len(item.Content) == 1 {
			return item.Content[0]
		}

		// Returns a random sentence
		rand.New(rand.NewSource(time.Now().UnixNano())) // depress: rand.Seed(time.Now().UnixNano())
		return item.Content[rand.Intn(len(item.Content))]
	}

	return ""
}

func SliceIncludes(collection []string, searchTerm string) bool { // slice -> collection, text -> searchTerm
	for _, element := range collection { // item -> element
		if element == searchTerm {
			return true
		}
	}

	return false
}

func SliceDifference(collection1 []string, collection2 []string) (difference []string) { // slice -> collection1, slice2 -> collection2
	// Loop two times, first to find collection1 strings not in collection2,
	// second loop to find collection2 strings not in collection1
	for i := 0; i < 2; i++ {
		for _, element1 := range collection1 { // s1 -> element1
			found := false
			for _, element2 := range collection2 { // s2 -> element2
				if element1 == element2 {
					found = true
					break
				}
			}
			// String not found. We add it to return slice
			if !found {
				difference = append(difference, element1)
			}
		}
		// Swap the slices, only if it was the first loop
		if i == 0 {
			collection1, collection2 = collection2, collection1
		}
	}

	return difference
}

func SliceIndex(collection []string, searchTerm string) int { // slice -> collection, text -> searchTerm
	for i, element := range collection { // item -> element
		if element == searchTerm {
			return i
		}
	}

	return 0
}

func UpdateUserProfile(authToken string, profileUpdater func(UserProfile) UserProfile) { // token -> authToken, changer -> profileUpdater
	cachedUserData[authToken] = profileUpdater(cachedUserData[authToken])
}

func StoreUserProfile(authToken string, profile UserProfile) { // token -> authToken, information -> profile
	cachedUserData[authToken] = profile
}

func RetrieveUserProfile(authToken string) UserProfile { // token -> authToken
	return cachedUserData[authToken]
}

/* training */

func trainDataMain(locale string) (inputs, outputs [][]float64) {
	words, classes, documents := Organize(locale)

	for _, document := range documents {
		outputRow := make([]float64, len(classes))
		bag := document.Sentence.WordsBag(words)

		// Change value to 1 where there is the document Tag
		outputRow[SliceIndex(classes, document.Tag)] = 1

		// Append data to inputs and outputs
		inputs = append(inputs, bag)
		outputs = append(outputs, outputRow)
	}

	return inputs, outputs
}

func CreateNeuralNetwork(locale string, ignoreTrainingFile bool) (neuralNetwork Network) {
	// Decide if the network is created by the save or is a new one
	saveFile := "../res/locales/" + locale + "/training.json"

	_, err := os.Open(saveFile)
	// Train the model if there is no training file
	if err != nil || ignoreTrainingFile {
		inputs, outputs := trainDataMain(locale)

		neuralNetwork = CreateNetwork(locale, 0.1, inputs, outputs, 50)
		neuralNetwork.Train(200)

		// Save the neural network in ../res/training.json
		neuralNetwork.Save(saveFile)
	} else {
		fmt.Printf(
			"%s %s\n",
			color.FgBlue.Render("Loading the neural network from"),
			color.FgRed.Render(saveFile),
		)
		// Initialize the intents
		SerializeIntents(locale)
		neuralNetwork = *LoadNetwork(saveFile)
	}

	return
}

func EncodeDashboardData(w http.ResponseWriter, r *http.Request) { // GetDashboardData -> EncodeDashboardData
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r) // data -> params

	dashboardData := DashboardData{ // dashboard -> dashboardData
		NetworkLayers: GetNetworkLayers(params["locale"]), // Layers -> NetworkLayers
		TrainingInfo:  GetTrainingInfo(params["locale"]),  // Training -> TrainingInfo
	}

	if err := json.NewEncoder(w).Encode(dashboardData); err != nil { // err := json.NewEncoder(w).Encode(dashboard) -> if err := json.NewEncoder(w).Encode(dashboardData)
		log.Fatal(err)
	}
}

func GetNetworkLayers(locale string) NetworkLayersData { // GetLayers -> GetNetworkLayers
	return NetworkLayersData{ // Layers -> NetworkLayersData
		// Get the number of rows of the first layer to get the count of input nodes
		InputCount: Rows(globalNeuralNetworks[locale].Layers[0]), // InputNodes -> InputCount
		// Get the number of hidden layers by removing the count of the input and output layers
		HiddenCount: len(globalNeuralNetworks[locale].Layers) - 2, // HiddenLayers -> HiddenCount
		// Get the number of rows of the latest layer to get the count of output nodes
		OutputCount: Columns(globalNeuralNetworks[locale].Output), // OutputNodes -> OutputCount
	}
}

func GetTrainingInfo(locale string) TrainingInfoData { // GetTraining -> GetTrainingInfo
	// Retrieve the information from the neural network
	return TrainingInfoData{ // Training -> TrainingInfoData
		LearningRate: globalNeuralNetworks[locale].Rate,   // Rate -> LearningRate
		ErrorMetrics: globalNeuralNetworks[locale].Errors, // Errors -> ErrorMetrics
		TrainingTime: globalNeuralNetworks[locale].Time,   // Time -> TrainingTime
	}
}

func StartServer(neuralNetworkInstances map[string]Network, serverPort string) { // Serve -> StartServer, _neuralNetworks -> neuralNetworkInstances, port -> serverPort
	// Set the current global network as a global variable
	globalNeuralNetworks = neuralNetworkInstances

	// Initializes the router
	router := mux.NewRouter()
	router.HandleFunc("/callback", CompleteAuth)
	// Serve the websocket
	router.HandleFunc("/websocket", HandleWebSocketConnection)
	// Serve the API
	router.HandleFunc("/api/{locale}/dashboard", EncodeDashboardData).Methods("GET")
	router.HandleFunc("/api/{locale}/intent", CreateIntent).Methods("POST")
	router.HandleFunc("/api/{locale}/intent", DeleteIntent).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/{locale}/train", TrainNeuralNetwork).Methods("POST") // Train -> TrainNeuralNetwork
	router.HandleFunc("/api/{locale}/intents", GetIntents).Methods("GET")
	router.HandleFunc("/api/coverage", GetCoverage).Methods("GET")

	magentaColor := color.FgMagenta.Render
	fmt.Printf("\nServer listening on the port %s...\n", magentaColor(serverPort)) // magenta -> magentaColor

	// Serves the chat
	if err := http.ListenAndServe(":"+serverPort, router); err != nil { // err := http.ListenAndServe(":"+port, router) -> if err := http.ListenAndServe(":"+serverPort, router)
		panic(err)
	}
}

func TrainNeuralNetwork(w http.ResponseWriter, r *http.Request) { // Train -> TrainNeuralNetwork
	// Checks if the token present in the headers is the right one
	token := r.Header.Get("Olivia-Token")
	if !ChecksToken(token) {
		json.NewEncoder(w).Encode(Error{
			Message: "You don't have the permission to do this.",
		})
		return
	}

	magentaColor := color.FgMagenta.Render
	fmt.Printf("\nRe-training the %s..\n", magentaColor("neural network")) // magenta -> magentaColor

	for locale := range globalNeuralNetworks { // neuralNetworks -> globalNeuralNetworks
		globalNeuralNetworks[locale] = CreateNeuralNetwork(locale, true)
	}
}

func HandleWebSocketConnection(w http.ResponseWriter, r *http.Request) { // SocketHandle -> HandleWebSocketConnection
	conn, _ := websocketUpgrader.Upgrade(w, r, nil) // upgrader -> websocketUpgrader
	fmt.Println(color.FgGreen.Render("A new connection has been opened"))

	for {
		// Read message from browser
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		// Unmarshal the json content of the message
		var request clientRequestMessage // RequestMessage -> clientRequestMessage
		if err = json.Unmarshal(msg, &request); err != nil {
			continue
		}

		// Set the information from the client into the cache
		if reflect.DeepEqual(RetrieveUserProfile(request.Token), UserProfile{}) {
			StoreUserProfile(request.Token, request.Information)
		}

		// If the type of requests is a handshake then execute the start modules
		if request.Type == 0 {
			ExecuteModules(request.Token, request.Locale)

			message := GetMessage()
			if message != "" {
				// Generate the response to send to the user
				response := serverResponseMessage{ // ResponseMessage -> serverResponseMessage
					Content:     message,
					Tag:         "start module",
					Information: RetrieveUserProfile(request.Token),
				}

				bytes, err := json.Marshal(response)
				if err != nil {
					panic(err)
				}

				if err = conn.WriteMessage(msgType, bytes); err != nil {
					continue
				}
			}

			continue
		}

		// Write message back to browser
		response := generateReply(request) // Reply -> generateReply
		if err = conn.WriteMessage(msgType, response); err != nil {
			continue
		}
	}
}

func generateReply(request clientRequestMessage) []byte { // Reply -> generateReply, RequestMessage -> clientRequestMessage
	var responseSentence, responseTag string

	// Send a message from ../res/datasets/messages.json if it is too long
	if len(request.Content) > 500 {
		responseTag = "too long"
		responseSentence = SelectRandomMessage(request.Locale, responseTag) // Keeping SelectRandomMessage as is
	} else {
		// If the given locale is not supported yet, set english
		locale := request.Locale
		if !Exists(locale) { // Keeping Exists as is
			locale = "en"
		}

		responseTag, responseSentence = NewSentence(
			locale, request.Content,
		).Calculate(*cacheInstance, globalNeuralNetworks[locale], request.Token) // Keeping NewSentence and Calculate as is
	}

	// Marshall the response in json
	response := serverResponseMessage{ // ResponseMessage -> serverResponseMessage
		Content:     responseSentence,
		Tag:         responseTag,
		Information: RetrieveUserProfile(request.Token),
	}

	bytes, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (network Network) CalculateFinalLayerDerivatives() LayerDerivative { // ComputeLastLayerDerivatives -> CalculateFinalLayerDerivatives
	l := len(network.Layers) - 1
	lastLayer := network.Layers[l]

	// Compute derivative for the last layer of weights and biases
	cost := Difference(network.Output, lastLayer)
	sigmoidDerivative := Multiplication(lastLayer, ApplyFunction(lastLayer, SubtractsOne))

	// Compute delta and the weights' adjustment
	delta := Multiplication(
		ApplyFunction(cost, MultipliesByTwo),
		sigmoidDerivative,
	)
	weights := DotProduct(Transpose(network.Layers[l-1]), delta)

	return LayerDerivative{ // Derivative -> LayerDerivative
		Delta:      delta,
		Adjustment: weights,
	}
}

func (network Network) CalculateLayerDerivatives(i int, derivatives []LayerDerivative) LayerDerivative { // ComputeDerivatives -> CalculateLayerDerivatives, Derivative -> LayerDerivative
	l := len(network.Layers) - 2 - i

	// Compute derivative for the layer of weights and biases
	delta := Multiplication(
		DotProduct(
			derivatives[i].Delta,
			Transpose(network.Weights[l]),
		),
		Multiplication(
			network.Layers[l],
			ApplyFunction(network.Layers[l], SubtractsOne),
		),
	)
	weights := DotProduct(Transpose(network.Layers[l-1]), delta)

	return LayerDerivative{ // Derivative -> LayerDerivative
		Delta:      delta,
		Adjustment: weights,
	}
}

func (network Network) ApplyAdjustments(derivatives []LayerDerivative) { // Adjust -> ApplyAdjustments, Derivative -> LayerDerivative
	for i, derivative := range derivatives {
		l := len(derivatives) - i

		network.Weights[l-1] = Sum(
			network.Weights[l-1],
			ApplyRate(derivative.Adjustment, network.Rate),
		)
		network.Biases[l-1] = Sum(
			network.Biases[l-1],
			ApplyRate(derivative.Delta, network.Rate),
		)
	}
}

func Sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

func MultipliesByTwo(x float64) float64 {
	return 2 * x
}

func SubtractsOne(x float64) float64 {
	return x - 1
}

func RandomMatrix(rows, columns int) (matrix Matrix) {
	matrix = make(Matrix, rows)

	for i := 0; i < rows; i++ {
		matrix[i] = make([]float64, columns)
		for j := 0; j < columns; j++ {
			matrix[i][j] = rand.Float64()*2.0 - 1.0
		}
	}

	return
}

func CreateMatrix(rows, columns int) (matrix Matrix) {
	matrix = make(Matrix, rows)

	for i := 0; i < rows; i++ {
		matrix[i] = make([]float64, columns)
	}

	return
}

func Rows(matrix Matrix) int {
	return len(matrix)
}

func Columns(matrix Matrix) int {
	return len(matrix[0])
}

func ApplyFunctionWithIndex(matrix Matrix, fn func(i, j int, x float64) float64) Matrix {
	for i := 0; i < Rows(matrix); i++ {
		for j := 0; j < Columns(matrix); j++ {
			matrix[i][j] = fn(i, j, matrix[i][j])
		}
	}

	return matrix
}

func ApplyFunction(matrix Matrix, fn func(x float64) float64) Matrix {
	return ApplyFunctionWithIndex(matrix, func(i, j int, x float64) float64 {
		return fn(x)
	})
}

func ApplyRate(matrix Matrix, rate float64) Matrix {
	return ApplyFunction(matrix, func(x float64) float64 {
		return rate * x
	})
}

func DotProduct(matrix, matrix2 Matrix) Matrix {
	if Columns(matrix) != Rows(matrix2) {
		panic("Cannot make dot product between these two matrix.")
	}

	return ApplyFunctionWithIndex(
		CreateMatrix(Rows(matrix), Columns(matrix2)),
		func(i, j int, x float64) float64 {
			var sum float64

			for k := 0; k < Columns(matrix); k++ {
				sum += matrix[i][k] * matrix2[k][j]
			}

			return sum
		},
	)
}

func Sum(matrix, matrix2 Matrix) (resultMatrix Matrix) {
	ErrorNotSameSize(matrix, matrix2)

	resultMatrix = CreateMatrix(Rows(matrix), Columns(matrix))

	return ApplyFunctionWithIndex(matrix, func(i, j int, x float64) float64 {
		return matrix[i][j] + matrix2[i][j]
	})
}

func Difference(matrix, matrix2 Matrix) (resultMatrix Matrix) {
	ErrorNotSameSize(matrix, matrix2)

	resultMatrix = CreateMatrix(Rows(matrix), Columns(matrix))

	return ApplyFunctionWithIndex(resultMatrix, func(i, j int, x float64) float64 {
		return matrix[i][j] - matrix2[i][j]
	})
}

func Multiplication(matrix, matrix2 Matrix) (resultMatrix Matrix) {
	ErrorNotSameSize(matrix, matrix2)

	resultMatrix = CreateMatrix(Rows(matrix), Columns(matrix))

	return ApplyFunctionWithIndex(matrix, func(i, j int, x float64) float64 {
		return matrix[i][j] * matrix2[i][j]
	})
}

func Transpose(matrix Matrix) (resultMatrix Matrix) {
	resultMatrix = CreateMatrix(Columns(matrix), Rows(matrix))

	for i := 0; i < Rows(matrix); i++ {
		for j := 0; j < Columns(matrix); j++ {
			resultMatrix[j][i] = matrix[i][j]
		}
	}

	return resultMatrix
}

func ErrorNotSameSize(matrix, matrix2 Matrix) {
	if Rows(matrix) != Rows(matrix2) && Columns(matrix) != Columns(matrix2) {
		panic("These two matrices must have the same dimension.")
	}
}

func LoadNetwork(fileName string) *Network {
	inF, err := os.Open(fileName)
	if err != nil {
		panic("Failed to load " + fileName + ".")
	}
	defer inF.Close()

	decoder := json.NewDecoder(inF)
	neuralNetwork := &Network{}
	err = decoder.Decode(neuralNetwork)
	if err != nil {
		panic(err)
	}

	return neuralNetwork
}

func CreateNetwork(locale string, rate float64, input, output Matrix, hiddensNodes ...int) Network {
	input = append([][]float64{
		make([]float64, len(input[0])),
	}, input...)
	output = append([][]float64{
		make([]float64, len(output[0])),
	}, output...)

	// Create the layers arrays and add the input values
	inputMatrix := input
	layers := []Matrix{inputMatrix}
	// Generate the hidden layer
	for _, hiddenNodes := range hiddensNodes {
		layers = append(layers, CreateMatrix(len(input), hiddenNodes))
	}
	// Add the output values to the layers arrays
	layers = append(layers, output)

	// Generate the weights and biases
	weightsNumber := len(layers) - 1
	var weights []Matrix
	var biases []Matrix

	for i := 0; i < weightsNumber; i++ {
		rows, columns := Columns(layers[i]), Columns(layers[i+1])

		weights = append(weights, RandomMatrix(rows, columns))
		biases = append(biases, RandomMatrix(Rows(layers[i]), columns))
	}

	return Network{
		Layers:  layers,
		Weights: weights,
		Biases:  biases,
		Output:  output,
		Rate:    rate,
		Locale:  locale,
	}
}

func (network Network) Save(fileName string) {
	outF, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		panic("Failed to save the network to " + fileName + ".")
	}
	defer outF.Close()

	encoder := json.NewEncoder(outF)
	err = encoder.Encode(network)
	if err != nil {
		panic(err)
	}
}

func (network *Network) FeedForward() {
	for i := 0; i < len(network.Layers)-1; i++ {
		layer, weights, biases := network.Layers[i], network.Weights[i], network.Biases[i]

		productMatrix := DotProduct(layer, weights)
		Sum(productMatrix, biases)
		ApplyFunction(productMatrix, Sigmoid)

		// Replace the output values
		network.Layers[i+1] = productMatrix
	}
}

func (network *Network) Predict(input []float64) []float64 {
	network.Layers[0] = Matrix{input}
	network.FeedForward()
	return network.Layers[len(network.Layers)-1][0]
}

func (network *Network) FeedBackward() {
	var derivatives []LayerDerivative
	derivatives = append(derivatives, network.CalculateFinalLayerDerivatives())

	// Compute the derivatives of the hidden layers
	for i := 0; i < len(network.Layers)-2; i++ {
		derivatives = append(derivatives, network.CalculateLayerDerivatives(i, derivatives))
	}

	// Then adjust the weights and biases
	network.ApplyAdjustments(derivatives)
}

func (network *Network) ComputeError() float64 {
	// Feed forward to compute the last layer's values
	network.FeedForward()
	lastLayer := network.Layers[len(network.Layers)-1]
	errors := Difference(network.Output, lastLayer)

	// Make the sum of all the errors
	var i int
	var sum float64
	for _, a := range errors {
		for _, e := range a {
			sum += e
			i++
		}
	}

	// Compute the average
	return sum / float64(i)
}

func (network *Network) Train(iterations int) {
	// Initialize the start date
	start := time.Now()

	// Create the progress bar
	bar := pb.New(iterations).Postfix(fmt.Sprintf(
		" - %s %s %s",
		color.FgBlue.Render("Training the"),
		color.FgRed.Render(GetNameByTag(network.Locale)),
		color.FgBlue.Render("neural network"),
	))
	bar.Format("(██░)")
	bar.SetMaxWidth(60)
	bar.ShowCounters = false
	bar.Start()

	// Train the network
	for i := 0; i < iterations; i++ {
		network.FeedForward()
		network.FeedBackward()

		// Append errors for dashboard data
		if i%(iterations/20) == 0 {
			network.Errors = append(
				network.Errors,
				// Round the error to two decimals
				network.ComputeError(),
			)
		}

		// Increment the progress bar
		bar.Increment()
	}

	bar.Finish()
	// Print the error
	arrangedError := fmt.Sprintf("%.5f", network.ComputeError())

	// Calculate elapsed date
	elapsed := time.Since(start)
	// Round the elapsed date at two decimals
	network.Time = math.Floor(elapsed.Seconds()*100) / 100

	fmt.Printf("The error rate is %s.\n", color.FgGreen.Render(arrangedError))
}

func GetCoverage(writer http.ResponseWriter, _ *http.Request) {
	allowedHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,Olivia-Token"
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	writer.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
	writer.Header().Set("Access-Control-Expose-Headers", "Authorization")

	defaultMessages, defaultIntents, defaultModules =
		RetrieveCachedMessages("en"), GetIntents_l("en"), GetModulesf("en")

	var coverage []LocaleCoverage

	// Calculate coverage for each language
	for _, locale := range Locales {
		if locale.Tag == "en" {
			continue
		}

		coverage = append(coverage, LocaleCoverage{
			Tag:      locale.Tag,
			Language: GetNameByTag(locale.Tag),
			Coverage: Coverage{
				Modules:  getModuleCoverage(locale.Tag),
				Intents:  getIntentCoverage(locale.Tag),
				Messages: getMessageCoverage(locale.Tag),
			},
		})
	}

	json.NewEncoder(writer).Encode(coverage)
}

func getMessageCoverage(locale string) CoverageDetails {
	var notCoveredMessages []string

	// Iterate through the default messages which are the english ones to verify if a message isn't
	// translated in the given locale.
	for _, defaultMessage := range defaultMessages {
		message := FindMessageByLabel(defaultMessage.Label, locale)

		// Add the current module tag to the list of not-covered-modules
		if message.Label != defaultMessage.Label {
			notCoveredMessages = append(notCoveredMessages, defaultMessage.Label)
		}
	}

	// Calculate the percentage of modules that aren't translated in the given locale
	coverage := calculateCoverage(len(notCoveredMessages), len(defaultMessages))

	return CoverageDetails{
		NotCovered: notCoveredMessages,
		Coverage:   coverage,
	}
}

func getIntentCoverage(locale string) CoverageDetails {
	var notCoveredIntents []string

	// Iterate through the default intents which are the english ones to verify if an intent isn't
	// translated in the given locale.
	for _, defaultIntent := range defaultIntents {
		intent := GetIntentByTag(defaultIntent.Tag, locale)

		// Add the current module tag to the list of not-covered-modules
		if intent.Tag != defaultIntent.Tag {
			notCoveredIntents = append(notCoveredIntents, defaultIntent.Tag)
		}
	}

	// Calculate the percentage of modules that aren't translated in the given locale
	coverage := calculateCoverage(len(notCoveredIntents), len(defaultModules))

	return CoverageDetails{
		NotCovered: notCoveredIntents,
		Coverage:   coverage,
	}
}

func getModuleCoverage(locale string) CoverageDetails {
	var notCoveredModules []string

	// Iterate through the default modules which are the english ones to verify if a module isn't
	// translated in the given locale.
	for _, defaultModule := range defaultModules {
		module := GetModuleByTagf(defaultModule.Tag, locale)

		// Add the current module tag to the list of not-covered-modules
		if module.Tag != defaultModule.Tag {
			notCoveredModules = append(notCoveredModules, defaultModule.Tag)
		}
	}

	// Calculate the percentage of modules that aren't translated in the given locale
	coverage := calculateCoverage(len(notCoveredModules), len(defaultModules))

	return CoverageDetails{
		NotCovered: notCoveredModules,
		Coverage:   coverage,
	}
}

func calculateCoverage(notCoveredLength, defaultLength int) int {
	return 100 * (defaultLength - notCoveredLength) / defaultLength
}

func (sentence *Sentence) arrange() {
	// Remove punctuation after letters
	punctuationRegex := regexp.MustCompile(`[a-zA-Z]( )?(\.|\?|!|¿|¡)`)
	sentence.Content = punctuationRegex.ReplaceAllStringFunc(sentence.Content, func(s string) string {
		punctuation := regexp.MustCompile(`(\.|\?|!)`)
		return punctuation.ReplaceAllString(s, "")
	})

	sentence.Content = strings.ReplaceAll(sentence.Content, "-", " ")
	sentence.Content = strings.TrimSpace(sentence.Content)
}

func removeStopWords(locale string, words []string) []string {
	// Don't remove stopwords for small sentences like “How are you” because it will remove all the words
	if len(words) <= 4 {
		return words
	}

	// Read the content of the stopwords file
	stopWords := string(FetchFileContent("../res/locales/" + locale + "/stopwords.txt"))

	var wordsToRemove []string

	// Iterate through all the stopwords
	for _, stopWord := range strings.Split(stopWords, "\n") {
		// Iterate through all the words of the given array
		for _, word := range words {
			// Continue if the word isn't a stopword
			if !strings.Contains(stopWord, word) {
				continue
			}

			wordsToRemove = append(wordsToRemove, word)
		}
	}

	return SliceDifference(words, wordsToRemove)
}

func (sentence Sentence) tokenize() (tokens []string) {
	// Split the sentence in words
	tokens = strings.Fields(sentence.Content)

	// Lower case each word
	for i, token := range tokens {
		tokens[i] = strings.ToLower(token)
	}

	tokens = removeStopWords(sentence.Locale, tokens)

	return
}

func (sentence Sentence) stem() (tokenizeWords []string) {
	locale := GetTagByName(sentence.Locale)
	// Set default locale to english
	if locale == "" {
		locale = "english"
	}

	tokens := sentence.tokenize()

	stemmer, err := snowball.New(locale)
	if err != nil {
		fmt.Println("Stemmer error", err)
		return
	}

	// Get the string token and push it to tokenizeWord
	for _, tokenizeWord := range tokens {
		word := stemmer.Stem(tokenizeWord)
		tokenizeWords = append(tokenizeWords, word)
	}

	return
}

func (sentence Sentence) WordsBag(words []string) (bag []float64) {
	for _, word := range words {
		// Append 1 if the patternWords contains the actual word, else 0
		var valueToAppend float64
		if SliceIncludes(sentence.stem(), word) {
			valueToAppend = 1
		}

		bag = append(bag, valueToAppend)
	}

	return bag
}

func CacheIntents(locale string, _intents []Intent) {
	intents[locale] = _intents
}

func GetIntents_l(locale string) []Intent {
	return intents[locale]
}

func SerializeIntents(locale string) (_intents []Intent) {
	err := json.Unmarshal(FetchFileContent("../res/locales/"+locale+"/intents.json"), &_intents)
	if err != nil {
		panic(err)
	}

	CacheIntents(locale, _intents)

	return _intents
}

func SerializeModulesIntents(locale string) []Intent {
	registeredModules := GetModulesf(locale)
	intents := make([]Intent, len(registeredModules))

	for k, module := range registeredModules {
		intents[k] = Intent{
			Tag:       module.Tag,
			Patterns:  module.Patterns,
			Responses: module.Responses,
			Context:   "",
		}
	}

	return intents
}

func GetIntentByTag(tag, locale string) Intent {
	for _, intent := range GetIntents_l(locale) {
		if tag != intent.Tag {
			continue
		}

		return intent
	}

	return Intent{}
}

func Organize(locale string) (words, classes []string, documents []Document) {
	// Append the modules intents to the intents from ../res/datasets/intents.json
	intents := append(
		SerializeIntents(locale),
		SerializeModulesIntents(locale)...,
	)

	for _, intent := range intents {
		for _, pattern := range intent.Patterns {
			// Tokenize the pattern's sentence
			patternSentence := Sentence{locale, pattern}
			patternSentence.arrange()

			// Add each word to response
			for _, word := range patternSentence.stem() {

				if !SliceIncludes(words, word) {
					words = append(words, word)
				}
			}

			// Add a new document
			documents = append(documents, Document{
				patternSentence,
				intent.Tag,
			})
		}

		// Add the intent tag to classes
		classes = append(classes, intent.Tag)
	}

	sort.Strings(words)
	sort.Strings(classes)

	return words, classes, documents
}

func NewSentence(locale, content string) (sentence Sentence) {
	sentence = Sentence{
		Locale:  locale,
		Content: content,
	}
	sentence.arrange()

	return
}

func (sentence Sentence) PredictTag(neuralNetwork Network) string {
	words, classes, _ := Organize(sentence.Locale)

	// Predict with the model
	predict := neuralNetwork.Predict(sentence.WordsBag(words))

	// Enumerate the results with the intent tags
	var resultsTag []Result
	for i, result := range predict {
		if i >= len(classes) {
			continue
		}
		resultsTag = append(resultsTag, Result{classes[i], result})
	}

	// Sort the results in ascending order
	sort.Slice(resultsTag, func(i, j int) bool {
		return resultsTag[i].Value > resultsTag[j].Value
	})

	LogResults(sentence.Locale, sentence.Content, resultsTag)

	return resultsTag[0].Tag
}

func RandomizeResponse(locale, entry, tag, token string) (string, string) {
	if tag == DontUnderstand {
		return DontUnderstand, SelectRandomMessage(locale, tag)
	}

	// Append the modules intents to the intents from ../res/datasets/intents.json
	intents := append(SerializeIntents(locale), SerializeModulesIntents(locale)...)

	for _, intent := range intents {
		if intent.Tag != tag {
			continue
		}

		// Reply a "don't understand" message if the context isn't correct
		cacheTag, _ := userCache.Get(token)
		if intent.Context != "" && cacheTag != intent.Context {
			return DontUnderstand, SelectRandomMessage(locale, DontUnderstand)
		}

		// Set the actual context
		userCache.Set(token, tag, gocache.DefaultExpiration)

		// Choose a random response in intents
		response := intent.Responses[0]
		if len(intent.Responses) > 1 {
			rand.Seed(time.Now().UnixNano())
			response = intent.Responses[rand.Intn(len(intent.Responses))]
		}

		// And then apply the triggers on the message
		return ReplaceContentf(locale, tag, entry, response, token)
	}

	return DontUnderstand, SelectRandomMessage(locale, DontUnderstand)
}

func (sentence Sentence) Calculate(cache gocache.Cache, neuralNetwork Network, token string) (string, string) {
	tag, found := cache.Get(sentence.Content)

	// Predict tag with the neural network if the sentence isn't in the cache
	if !found {
		tag = sentence.PredictTag(neuralNetwork)
		cache.Set(sentence.Content, tag, gocache.DefaultExpiration)
	}

	return RandomizeResponse(sentence.Locale, sentence.Content, tag.(string), token)
}

func LogResults(locale, entry string, results []Result) {
	// If NO_LOGS is present, then don't print the given messages
	if os.Getenv("NO_LOGS") == "1" {
		return
	}

	green := color.FgGreen.Render
	yellow := color.FgYellow.Render

	fmt.Printf(
		"\n“%s” - %s\n",
		color.FgCyan.Render(entry),
		color.FgRed.Render(GetNameByTag(locale)),
	)
	for _, result := range results {
		// Arbitrary choice of 0.004 to have less tags to show
		if result.Value < 0.004 {
			continue
		}

		fmt.Printf("  %s %s - %s\n", green("▫︎"), result.Tag, yellow(result.Value))
	}
}

func GenerateToken() string {
	b := make([]byte, 30)
	rand.Read(b)

	fmt.Println("hey")
	return fmt.Sprintf("%x", b)
}

func HashToken(token string) []byte {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(token), 14)
	return bytes
}

func ChecksToken(token string) bool {
	err := bcrypt.CompareHashAndPassword(authenticationHash, []byte(token))
	return err == nil
}

func AuthenticationFileExists() bool {
	_, err := os.Open(fileName)
	return err == nil
}

func SaveHash(hash string) {
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	file.WriteString(hash)
}

func Authenticate() {
	// Do nothing if the authentication file exists
	if AuthenticationFileExists() {
		authenticationHash = FetchFileContent(fileName)
		return
	}

	// Generates the token and gives it to the user
	token := GenerateToken()
	fmt.Printf("Your authentication token is: %s\n", color.FgLightGreen.Render(token))
	fmt.Println("Save it, you won't be able to get it again unless you generate a new one.")
	fmt.Println()

	// Hash the token and save it
	hash := HashToken(token)
	SaveHash(string(hash))

	authenticationHash = hash
}

func WriteIntents(locale string, intents []Intent) {
	CacheIntents(locale, intents)

	// Encode the json
	bytes, _ := json.MarshalIndent(intents, "", "  ")

	// Write it to the file
	file, err := os.Create("../res/locales/" + locale + "/intents.json")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	file.Write(bytes)
}

func AddIntent(locale string, intent Intent) {
	intents := append(SerializeIntents(locale), intent)

	WriteIntents(locale, intents)

	fmt.Printf("Added %s intent.\n", color.FgMagenta.Render(intent.Tag))
}

func RemoveIntent(locale, tag string) {
	intents := SerializeIntents(locale)

	// Iterate through the intents to remove the right one
	for i, intent := range intents {
		if intent.Tag != tag {
			continue
		}

		intents[i] = intents[len(intents)-1]
		intents = intents[:len(intents)-1]
		fmt.Printf("The intent %s was deleted.\n", color.FgMagenta.Render(intent.Tag))
	}

	WriteIntents(locale, intents)
}

func GetIntents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	data := mux.Vars(r)

	// Encode the intents for the given locale
	json.NewEncoder(w).Encode(GetIntents_l(data["locale"]))
}

func CreateIntent(w http.ResponseWriter, r *http.Request) {
	allowedHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,Olivia-Token"
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
	w.Header().Set("Access-Control-Expose-Headers", "Authorization")

	data := mux.Vars(r)

	// Checks if the token present in the headers is the right one
	token := r.Header.Get("Olivia-Token")
	if !ChecksToken(token) {
		json.NewEncoder(w).Encode(Error{
			Message: SelectRandomMessage(data["locale"], "no permission"),
		})
		return
	}

	// Decode request json body
	var intent Intent
	json.NewDecoder(r.Body).Decode(&intent)

	if intent.Responses == nil || intent.Patterns == nil {
		json.NewEncoder(w).Encode(Error{
			Message: SelectRandomMessage(data["locale"], "patterns same"),
		})
		return
	}

	// Returns an error if the tags are the same
	for _, _intent := range GetIntents_l(data["locale"]) {
		if _intent.Tag == intent.Tag {
			json.NewEncoder(w).Encode(Error{
				Message: SelectRandomMessage(data["locale"], "tags same"),
			})
			return
		}
	}

	// Adds the intent
	AddIntent(data["locale"], intent)

	json.NewEncoder(w).Encode(intent)
}

func DeleteIntent(w http.ResponseWriter, r *http.Request) {
	allowedHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,Olivia-Token"
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
	w.Header().Set("Access-Control-Expose-Headers", "Authorization")

	data := mux.Vars(r)

	// Checks if the token present in the headers is the right one
	token := r.Header.Get("Olivia-Token")
	if !ChecksToken(token) {
		json.NewEncoder(w).Encode(Error{
			Message: SelectRandomMessage(data["locale"], "no permission"),
		})
		return
	}

	var deleteRequest DeleteRequest
	json.NewDecoder(r.Body).Decode(&deleteRequest)

	RemoveIntent(data["locale"], deleteRequest.Tag)

	json.NewEncoder(w).Encode(GetIntents_l(data["locale"]))
}

func GetNameByTag(tag string) string {
	for _, locale := range Locales {
		if locale.Tag != tag {
			continue
		}

		return locale.Name
	}

	return ""
}

func GetTagByName(name string) string {
	for _, locale := range Locales {
		if locale.Name != name {
			continue
		}

		return locale.Tag
	}

	return ""
}

func Exists(tag string) bool {
	for _, locale := range Locales {
		if locale.Tag == tag {
			return true
		}
	}

	return false
}

func SearchTime(locale, sentence string) (string, time.Time) {
	_time := RuleTime(sentence)
	// Set the time to 12am if no time has been found
	if _time == (time.Time{}) {
		_time = time.Date(0, 0, 0, 12, 0, 0, 0, time.UTC)
	}

	for _, rule := range rules {
		date := rule(locale, sentence)

		// If the current rule found a date
		if date != (time.Time{}) {
			date = time.Date(date.Year(), date.Month(), date.Day(), _time.Hour(), _time.Minute(), 0, 0, time.UTC)

			sentence = DeleteTimes(locale, sentence)
			return DeleteDates(locale, sentence), date
		}
	}

	return sentence, time.Now().Add(time.Hour * 24)
}

func DeleteDates(locale, sentence string) string {
	// Create a regex to match the patterns of dates to remove them.
	datePatterns := regexp.MustCompile(PatternTranslation[locale].DateRegex)

	// Replace the dates by empty string
	sentence = datePatterns.ReplaceAllString(sentence, "")
	// Trim the spaces and return
	return strings.TrimSpace(sentence)
}

func DeleteTimes(locale, sentence string) string {
	// Create a regex to match the patterns of times to remove them.
	timePatterns := regexp.MustCompile(PatternTranslation[locale].TimeRegex)

	// Replace the times by empty string
	sentence = timePatterns.ReplaceAllString(sentence, "")
	// Trim the spaces and return
	return strings.TrimSpace(sentence)
}

func RegisterRule(rule Rule) {
	rules = append(rules, rule)
}

func init() {
	// Register the rules
	RegisterRule(RuleToday)
	RegisterRule(RuleTomorrow)
	RegisterRule(RuleDayOfWeek)
	RegisterRule(RuleNaturalDate)
	RegisterRule(RuleDate)
}

func RuleToday(locale, sentence string) (result time.Time) {
	todayRegex := regexp.MustCompile(RuleTranslations[locale].RuleToday)
	today := todayRegex.FindString(sentence)

	// Returns an empty date struct if no date has been found
	if today == "" {
		return time.Time{}
	}

	return time.Now()
}

func RuleTomorrow(locale, sentence string) (result time.Time) {
	tomorrowRegex := regexp.MustCompile(RuleTranslations[locale].RuleTomorrow)
	date := tomorrowRegex.FindString(sentence)

	// Returns an empty date struct if no date has been found
	if date == "" {
		return time.Time{}
	}

	result = time.Now().Add(day)

	// If the date contains "after", we add 24 hours to tomorrow's date
	if strings.Contains(date, RuleTranslations[locale].RuleAfterTomorrow) {
		return result.Add(day)
	}

	return
}

func RuleDayOfWeek(locale, sentence string) time.Time {
	dayOfWeekRegex := regexp.MustCompile(RuleTranslations[locale].RuleDayOfWeek)
	date := dayOfWeekRegex.FindString(sentence)

	// Returns an empty date struct if no date has been found
	if date == "" {
		return time.Time{}
	}

	var foundDayOfWeek int
	// Find the integer value of the found day of the week
	for _, dayOfWeek := range daysOfWeek {
		// Down case the day of the week to match the found date
		stringDayOfWeek := strings.ToLower(dayOfWeek.String())

		if strings.Contains(date, stringDayOfWeek) {
			foundDayOfWeek = int(dayOfWeek)
		}
	}

	currentDay := int(time.Now().Weekday())
	// Calculate the date of the found day
	calculatedDate := foundDayOfWeek - currentDay

	// If the day is already passed in the current week, then we add another week to the count
	if calculatedDate <= 0 {
		calculatedDate += 7
	}

	// If there is "next" in the sentence, then we add another week
	if strings.Contains(date, RuleTranslations[locale].RuleNextDayOfWeek) {
		calculatedDate += 7
	}

	// Then add the calculated number of day to the actual date
	return time.Now().Add(day * time.Duration(calculatedDate))
}

func RuleNaturalDate(locale, sentence string) time.Time {
	naturalMonthRegex := regexp.MustCompile(
		RuleTranslations[locale].RuleNaturalDate,
	)
	naturalDayRegex := regexp.MustCompile(`\d{2}|\d`)

	month := naturalMonthRegex.FindString(sentence)
	day := naturalDayRegex.FindString(sentence)

	// Put the month in english to parse the time with time golang package
	if locale != "en" {
		monthIndex := SliceIndex(RuleTranslations[locale].Months, month)
		month = RuleTranslations["en"].Months[monthIndex]
	}

	parsedMonth, _ := time.Parse("January", month)
	parsedDay, _ := strconv.Atoi(day)

	// Returns an empty date struct if no date has been found
	if day == "" && month == "" {
		return time.Time{}
	}

	// If only the month is specified
	if day == "" {
		// Calculate the number of months to add
		calculatedMonth := parsedMonth.Month() - time.Now().Month()
		// Add a year if the month is passed
		if calculatedMonth <= 0 {
			calculatedMonth += 12
		}

		// Remove the number of days elapsed in the month to reach the first
		return time.Now().AddDate(0, int(calculatedMonth), -time.Now().Day()+1)
	}

	// Parse the date
	parsedDate := fmt.Sprintf("%d-%02d-%02d", time.Now().Year(), parsedMonth.Month(), parsedDay)
	date, err := time.Parse("2006-01-02", parsedDate)
	if err != nil {
		return time.Time{}
	}

	// If the date has been passed, add a year
	if time.Now().After(date) {
		date = date.AddDate(1, 0, 0)
	}

	return date
}

func RuleDate(locale, sentence string) time.Time {
	dateRegex := regexp.MustCompile(`(\d{2}|\d)/(\d{2}|\d)`)
	date := dateRegex.FindString(sentence)

	// Returns an empty date struct if no date has been found
	if date == "" {
		return time.Time{}
	}

	// Parse the found date
	parsedDate, err := time.Parse("01/02", date)
	if err != nil {
		return time.Time{}
	}

	// Add the current year to the date
	parsedDate = parsedDate.AddDate(time.Now().Year(), 0, 0)

	// Add another year if the date is passed
	if time.Now().After(parsedDate) {
		parsedDate = parsedDate.AddDate(1, 0, 0)
	}

	return parsedDate
}

func RuleTime(sentence string) time.Time {
	timeRegex := regexp.MustCompile(`(\d{2}|\d)(:\d{2}|\d)?( )?(pm|am|p\.m|a\.m)`)
	foundTime := timeRegex.FindString(sentence)

	// Returns an empty date struct if no date has been found
	if foundTime == "" {
		return time.Time{}
	}

	// Initialize the part of the day asked
	var part string
	if strings.Contains(foundTime, "pm") || strings.Contains(foundTime, "p.m") {
		part = "pm"
	} else if strings.Contains(foundTime, "am") || strings.Contains(foundTime, "a.m") {
		part = "am"
	}

	if strings.Contains(foundTime, ":") {
		// Get the hours and minutes of the given time
		hoursAndMinutesRegex := regexp.MustCompile(`(\d{2}|\d):(\d{2}|\d)`)
		timeVariables := strings.Split(hoursAndMinutesRegex.FindString(foundTime), ":")

		// Format the time with 2 digits for each
		formattedTime := fmt.Sprintf("%02s:%02s %s", timeVariables[0], timeVariables[1], part)
		response, _ := time.Parse("03:04 pm", formattedTime)

		return response
	}

	digitsRegex := regexp.MustCompile(`\d{2}|\d`)
	foundDigits := digitsRegex.FindString(foundTime)

	formattedTime := fmt.Sprintf("%02s %s", foundDigits, part)
	response, _ := time.Parse("03 pm", formattedTime)

	return response
}

func SerializeCountries() (countries []Country) {
	err := json.Unmarshal(FetchFileContent("../res/datasets/countries.json"), &countries)
	if err != nil {
		fmt.Println(err)
	}

	return countries
}

func FindCountry(locale, sentence string) Country {
	for _, country := range countries {
		name, exists := country.Name[locale]

		if !exists {
			continue
		}

		// If the actual country isn't contained in the sentence, continue
		if !strings.Contains(strings.ToLower(sentence), strings.ToLower(name)) {
			continue
		}

		// Returns the right country
		return country
	}

	// Returns an empty country if none has been found
	return Country{}
}

func LevenshteinDistance(first, second string) int {
	// Returns the length if it's empty
	if first == "" {
		return len(second)
	}
	if second == "" {
		return len(first)
	}

	if first[0] == second[0] {
		return LevenshteinDistance(first[1:], second[1:])
	}

	a := LevenshteinDistance(first[1:], second[1:])
	if b := LevenshteinDistance(first, second[1:]); a > b {
		a = b
	}

	if c := LevenshteinDistance(first[1:], second); a > c {
		a = c
	}

	return a + 1
}

func LevenshteinContains(sentence, matching string, rate int) bool {
	words := strings.Split(sentence, " ")
	for _, word := range words {
		// Returns true if the distance is below the rate
		if LevenshteinDistance(word, matching) <= rate {
			return true
		}
	}

	return false
}

func FindMathOperation(entry string) string {
	mathRegex := regexp.MustCompile(
		`((\()?(((\d+|pi)(\^\d+|!|.)?)|sqrt|cos|sin|tan|acos|asin|atan|log|ln|abs)( )?[+*\/\-x]?( )?(\))?[+*\/\-]?)+`,
	)

	operation := mathRegex.FindString(entry)
	// Replace "x" symbol by "*"
	operation = strings.Replace(operation, "x", "*", -1)
	return strings.TrimSpace(operation)
}

func FindNumberOfDecimals(locale, entry string) int {
	decimalsRegex := regexp.MustCompile(
		MathDecimals[locale],
	)
	numberRegex := regexp.MustCompile(`\d+`)

	decimals := numberRegex.FindString(decimalsRegex.FindString(entry))
	decimalsInt, _ := strconv.Atoi(decimals)

	return decimalsInt
}

func SerializeMovies() (movies []Movie) {
	path := "../res/datasets/movies.csv"
	bytes, err := os.Open(path)
	if err != nil {
		bytes, _ = os.Open("../" + path)
	}

	reader := csv.NewReader(bufio.NewReader(bytes))
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		// Convert the string to a float
		rating, _ := strconv.ParseFloat(line[3], 64)

		movies = append(movies, Movie{
			Name:   line[1],
			Genres: strings.Split(line[2], "|"),
			Rating: rating,
		})
	}

	return
}

func SearchMovie(genre, userToken string) (output Movie) {
	for _, movie := range movies {
		userMovieBlacklist := RetrieveUserProfile(userToken).DislikedMovies
		// Continue if the movie is not from the request genre or if this movie has already been suggested
		if !SliceIncludes(movie.Genres, genre) || SliceIncludes(userMovieBlacklist, movie.Name) {
			continue
		}

		if reflect.DeepEqual(output, Movie{}) || movie.Rating > output.Rating {
			output = movie
		}
	}

	// Add the found movie to the user blacklist
	UpdateUserProfile(userToken, func(information UserProfile) UserProfile {
		information.DislikedMovies = append(information.DislikedMovies, output.Name)
		return information
	})

	return
}

func FindMoviesGenres(locale, content string) (output []string) {
	for i, genre := range MoviesGenres[locale] {
		if LevenshteinContains(strings.ToUpper(content), strings.ToUpper(genre), 2) {
			output = append(output, MoviesGenres["en"][i])
		}
	}

	return
}

func SearchMusic(locale, sentence string) (music, artist string) {
	words := strings.Split(sentence, " ")

	// Iterate through the words of the sentence
	playAppeared, fromAppeared, onAppeared := false, false, false
	for _, word := range words {
		// If "on" appeared
		if word == SpotifyKeyword[locale].On {
			onAppeared = true
		}

		// Add the current word if its between from and on
		if fromAppeared && !onAppeared {
			artist += word + " "
		}

		// If "from" appeared
		if LevenshteinDistance(word, SpotifyKeyword[locale].From) < 2 {
			fromAppeared = true
		}

		// Add the current word if its between play and from
		if playAppeared && !fromAppeared && !onAppeared {
			music += word + " "
		}

		// If "play" appeared
		if LevenshteinDistance(word, SpotifyKeyword[locale].Play) < 2 {
			playAppeared = true
		}
	}

	// Trim the spaces and return
	return strings.TrimSpace(music), strings.TrimSpace(artist)
}

func SerializeNames() (names []string) {
	namesFile := string(FetchFileContent("../res/datasets/names.txt"))

	// Iterate each line of the file
	names = append(names, strings.Split(strings.TrimSuffix(namesFile, "\n"), "\n")...)
	return
}

func FindName(sentence string) string {
	for _, name := range names {
		if !strings.Contains(strings.ToLower(" "+sentence+" "), " "+name+" ") {
			continue
		}

		return name
	}

	return ""
}

func FindRangeLimits(local, entry string) ([]int, error) {
	decimalsRegex := regexp.MustCompile(decimal)
	limitStrArr := decimalsRegex.FindAllString(entry, 2)
	limitArr := make([]int, 0)

	if limitStrArr == nil {
		return make([]int, 0), errors.New("No range")
	}

	if len(limitStrArr) != 2 {
		return nil, errors.New("Need 2 numbers, a lower and upper limit")
	}

	for _, v := range limitStrArr {
		num, err := strconv.Atoi(v)
		if err != nil {
			return nil, errors.New("Non integer range")
		}
		limitArr = append(limitArr, num)
	}

	sort.Ints(limitArr)
	return limitArr, nil
}

func SearchReason(locale, sentence string) string {
	var response []string

	// Split the given sentence into words
	words := strings.Split(sentence, " ")

	// Initialize the appeared boolean for the keywords
	appeared := false
	// Iterate through the words
	for _, word := range words {
		// Add the word to the response array if the keyword already appeared
		if appeared {
			response = append(response, word)
		}

		// If the keyword didn't appeared and one of the keywords match set the appeared condition
		// to true
		if !appeared && (LevenshteinDistance(word, ReasonKeywords[locale].That) <= 2 ||
			LevenshteinDistance(word, ReasonKeywords[locale].To) < 2) {
			appeared = true
		}
	}

	// Join the words and return the sentence
	return strings.Join(response, " ")
}

func SearchTokens(sentence string) []string {
	// Search the token with a regex
	tokenRegex := regexp.MustCompile(`[a-z0-9]{32}`)
	// Returns the found token
	return tokenRegex.FindAllString(sentence, 2)
}

func init() {
	// Set default value of the callback url
	if callbackURL == "" {
		callbackURL = "https://olivia-api.herokuapp.com/callback"
	}

	// Set default value of the redirect url
	if redirectURL == "" {
		redirectURL = "https://olivia-ai.org/chat"
	}

	// Initialize the authenticator
	auth = spotify.NewAuthenticator(
		callbackURL,
		spotify.ScopeStreaming,
		spotify.ScopeUserModifyPlaybackState,
		spotify.ScopeUserReadPlaybackState,
	)
}

func LoginSpotify(locale, token string) string {
	information := RetrieveUserProfile(token)

	// Generate the authentication url
	auth.SetAuthInfo(information.StreamingID, information.StreamingSecret)
	url := auth.AuthURL(state)

	// Waits for the authentication to be completed, and save the client in user's information
	go func() {
		authenticationToken := <-tokenChannel

		// If the token is empty reset the credentials of the user
		if *authenticationToken == (oauth2.Token{}) {
			UpdateUserProfile(token, func(information UserProfile) UserProfile {
				information.StreamingID = ""
				information.StreamingSecret = ""

				return information
			})
		}

		// Save the authentication token
		UpdateUserProfile(token, func(information UserProfile) UserProfile {
			information.StreamingToken = authenticationToken

			return information
		})
	}()

	return fmt.Sprintf(SelectRandomMessage(locale, "spotify login"), url)
}

func RenewSpotifyToken(token string) spotify.Client {
	authenticationToken := RetrieveUserProfile(token).StreamingToken
	client := auth.NewClient(authenticationToken)

	// Renew the authentication token
	if m, _ := time.ParseDuration("5m30s"); time.Until(authenticationToken.Expiry) < m {
		UpdateUserProfile(token, func(information UserProfile) UserProfile {
			information.StreamingToken, _ = client.Token()
			return information
		})
	}

	return client
}

func CheckTokensPresence(token string) bool {
	information := RetrieveUserProfile(token)
	return information.StreamingID == "" || information.StreamingSecret == ""
}

func CompleteAuth(w http.ResponseWriter, r *http.Request) {
	token, err := auth.Token(state, r)

	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		tokenChannel <- &oauth2.Token{}
		return
	}

	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		tokenChannel <- &oauth2.Token{}
		return
	}

	// Use the token to get an authenticated client
	w.Header().Set("Content-Type", "text/html")
	// Redirect the user
	fmt.Fprintf(w, `<meta http-equiv="refresh" content="0; url = %s" />`, redirectURL)

	tokenChannel <- token
}

func RegisterModule(module Module) {
	modules = append(modules, module)
}

func SetMessage(_message string) {
	message = _message
}

func GetMessage() string {
	return message
}

func ExecuteModules(token, locale string) {
	fmt.Println(color.FgGreen.Render("Executing start modules.."))

	for _, module := range modules {
		module.Action(token, locale)
	}
}

func init() {
	RegisterModule(Module{
		Action: CheckReminders,
	})
}

func CheckReminders(token, locale string) {
	reminders := RetrieveUserProfile(token).ImportantDates
	var messages []string

	// Iterate through the reminders to check if they are outdated
	for i, reminder := range reminders {
		date, _ := time.Parse("01/02/2006 03:04", reminder.ReminderDate)

		now := time.Now()
		// If the date is today
		if date.Year() == now.Year() && date.Day() == now.Day() && date.Month() == now.Month() {
			messages = append(messages, fmt.Sprintf("“%s”", reminder.ReminderDetails))

			// Removes the current reminder
			RemoveUserReminder(token, i)
		}
	}

	// Send the startup message
	if len(messages) != 0 {
		// If the message is already filled in return.
		if GetMessage() != "" {
			return
		}

		// Set the message with the current reminders
		listRemindersMessage := SelectRandomMessage(locale, "list reminders")
		if listRemindersMessage == "" {
			return
		}

		message := fmt.Sprintf(
			listRemindersMessage,
			RetrieveUserProfile(token).FullName,
			strings.Join(messages, ", "),
		)
		SetMessage(message)
	}
}

func RemoveUserReminder(token string, index int) {
	UpdateUserProfile(token, func(information UserProfile) UserProfile {
		reminders := information.ImportantDates

		// Removes the element from the reminders slice
		if len(reminders) == 1 {
			reminders = []UserReminder{}
		} else {
			reminders[index] = reminders[len(reminders)-1]
			reminders = reminders[:len(reminders)-1]
		}

		// Set the updated slice
		information.ImportantDates = reminders

		return information
	})
}

func AdvicesReplacer(locale, entry, response, _ string) (string, string) {

	resp, err := http.Get(adviceURL)
	if err != nil {
		responseTag := "no advices"
		return responseTag, SelectRandomMessage(locale, responseTag)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		responseTag := "no advices"
		return responseTag, SelectRandomMessage(locale, responseTag)
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)
	advice := result["slip"].(map[string]interface{})["advice"]

	return AdvicesTag, fmt.Sprintf(response, advice)
}

func AreaReplacer(locale, entry, response, _ string) (string, string) {
	country := FindCountry(locale, entry)

	// If there isn't a country respond with a message from ../res/datasets/messages.json
	if country.Currency == "" {
		responseTag := "no country"
		return responseTag, SelectRandomMessage(locale, responseTag)
	}

	return AreaTag, fmt.Sprintf(response, ArticleCountries[locale](country.Name[locale]), country.Area)
}

func CapitalReplacer(locale, entry, response, _ string) (string, string) {
	country := FindCountry(locale, entry)

	// If there isn't a country respond with a message from ../res/datasets/messages.json
	if country.Currency == "" {
		responseTag := "no country"
		return responseTag, SelectRandomMessage(locale, responseTag)
	}

	articleFunction, exists := ArticleCountries[locale]
	countryName := country.Name[locale]
	if exists {
		countryName = articleFunction(countryName)
	}

	return CapitalTag, fmt.Sprintf(response, countryName, country.Capital)
}

func CurrencyReplacer(locale, entry, response, _ string) (string, string) {
	country := FindCountry(locale, entry)

	// If there isn't a country respond with a message from ../res/datasets/messages.json
	if country.Currency == "" {
		responseTag := "no country"
		return responseTag, SelectRandomMessage(locale, responseTag)
	}

	return CurrencyTag, fmt.Sprintf(response, ArticleCountries[locale](country.Name[locale]), country.Currency)
}

func JokesReplacer(locale, entry, response, _ string) (string, string) {

	resp, err := http.Get(jokeURL)
	if err != nil {
		responseTag := "no jokes"
		return responseTag, SelectRandomMessage(locale, responseTag)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		responseTag := "no jokes"
		return responseTag, SelectRandomMessage(locale, responseTag)
	}

	joke := &Joke{}

	err = json.Unmarshal(body, joke)
	if err != nil {
		responseTag := "no jokes"
		return responseTag, SelectRandomMessage(locale, responseTag)
	}

	jokeStr := joke.Setup + " " + joke.Punchline

	return JokesTag, fmt.Sprintf(response, jokeStr)
}

func MathReplacer(locale, entry, response, _ string) (string, string) {
	operation := FindMathOperation(entry)

	// If there is no operation in the entry message reply with a "don't understand" message
	if operation == "" {
		responseTag := "don't understand"
		return responseTag, SelectRandomMessage(locale, responseTag)
	}

	res, err := mathcat.Eval(operation)
	// If the expression isn't valid reply with a message from ../res/datasets/messages.json
	if err != nil {
		responseTag := "math not valid"
		return responseTag, SelectRandomMessage(locale, responseTag)
	}
	// Use number of decimals from the query
	decimals := FindNumberOfDecimals(locale, entry)
	if decimals == 0 {
		decimals = 6
	}

	result := res.FloatString(decimals)

	// Remove trailing zeros of the result with a Regex
	trailingZerosRegex := regexp.MustCompile(`\.?0+$`)
	result = trailingZerosRegex.ReplaceAllString(result, "")

	return MathTag, fmt.Sprintf(response, result)
}

func RegisterModulef(locale string, module Modulef) {
	modulesf[locale] = append(modulesf[locale], module)
}

func RegisterModulesf(locale string, _modules []Modulef) {
	modulesf[locale] = append(modulesf[locale], _modules...)
}

func GetModulesf(locale string) []Modulef {
	return modulesf[locale]
}

func GetModuleByTagf(tag, locale string) Modulef {
	for _, module := range modulesf[locale] {
		if tag != module.Tag {
			continue
		}

		return module
	}

	return Modulef{}
}

func ReplaceContentf(locale, tag, entry, response, token string) (string, string) {
	for _, module := range modulesf[locale] {
		if module.Tag != tag {
			continue
		}

		return module.Replacer(locale, entry, response, token)
	}

	return tag, response
}

func GenresReplacer(locale, entry, response, token string) (string, string) {
	genres := FindMoviesGenres(locale, entry)

	// If there is no genres then reply with a message from ../res/datasets/messages.json
	if len(genres) == 0 {
		responseTag := "no genres"
		return responseTag, SelectRandomMessage(locale, responseTag)
	}

	// Change the user information to add the new genres
	UpdateUserProfile(token, func(information UserProfile) UserProfile {
		for _, genre := range genres {
			// Append the genre only is it isn't already in the information
			if SliceIncludes(information.GenrePreferences, genre) {
				continue
			}

			information.GenrePreferences = append(information.GenrePreferences, genre)
		}
		return information
	})

	return GenresTag, response
}

func MovieSearchReplacer(locale, entry, response, token string) (string, string) {
	genres := FindMoviesGenres(locale, entry)

	// If there is no genres then reply with a message from ../res/datasets/messages.json
	if len(genres) == 0 {
		responseTag := "no genres"
		return responseTag, SelectRandomMessage(locale, responseTag)
	}

	movie := SearchMovie(genres[0], token)

	return MoviesTag, fmt.Sprintf(response, movie.Name, movie.Rating)
}

func MovieSearchFromInformationReplacer(locale, _, response, token string) (string, string) {
	// If there is no genres then reply with a message from ../res/datasets/messages.json
	genres := RetrieveUserProfile(token).GenrePreferences
	if len(genres) == 0 {
		responseTag := "no genres saved"
		return responseTag, SelectRandomMessage(locale, responseTag)
	}

	movie := SearchMovie(genres[rand.Intn(len(genres))], token)
	genresJoined := strings.Join(genres, ", ")
	return MoviesDataTag, fmt.Sprintf(response, genresJoined, movie.Name, movie.Rating)
}

func NameGetterReplacer(locale, _, response, token string) (string, string) {
	name := RetrieveUserProfile(token).FullName

	if strings.TrimSpace(name) == "" {
		responseTag := "don't know name"
		return responseTag, SelectRandomMessage(locale, responseTag)
	}

	return NameGetterTag, fmt.Sprintf(response, name)
}

func NameSetterReplacer(locale, entry, response, token string) (string, string) {
	name := FindName(entry)

	// If there is no name in the entry string
	if name == "" {
		responseTag := "no name"
		return responseTag, SelectRandomMessage(locale, responseTag)
	}

	// Capitalize the name
	name = strings.Title(name)

	// Change the name inside the user information
	UpdateUserProfile(token, func(information UserProfile) UserProfile {
		information.FullName = name
		return information
	})

	return NameSetterTag, fmt.Sprintf(response, name)
}

func RandomNumberReplacer(locale, entry, response, _ string) (string, string) {
	limitArr, err := FindRangeLimits(locale, entry)
	if err != nil {
		if limitArr != nil {
			return RandomTag, fmt.Sprintf(response, strconv.Itoa(rand.Intn(100)))
		}

		responseTag := "no random range"
		return responseTag, SelectRandomMessage(locale, responseTag)
	}

	min := limitArr[0]
	max := limitArr[1]
	randNum := rand.Intn((max - min)) + min
	return RandomTag, fmt.Sprintf(response, strconv.Itoa(randNum))
}

func ReminderSetterReplacer(locale, entry, response, token string) (string, string) {
	// Search the time and
	sentence, date := SearchTime(locale, entry)
	reason := SearchReason(locale, sentence)

	// Format the date
	formattedDate := date.Format("01/02/2006 03:04")

	// Add the reminder inside the user's information
	UpdateUserProfile(token, func(information UserProfile) UserProfile {
		information.ImportantDates = append(information.ImportantDates, UserReminder{
			ReminderDetails: reason,
			ReminderDate:    formattedDate,
		})

		return information
	})

	return ReminderSetterTag, fmt.Sprintf(response, reason, formattedDate)
}

func ReminderGetterReplacer(locale, _, response, token string) (string, string) {
	reminders := RetrieveUserProfile(token).ImportantDates
	var formattedReminders []string

	// Iterate through the reminders and parse them
	for _, reminder := range reminders {
		formattedReminder := fmt.Sprintf(
			SelectRandomMessage(locale, "reminder"),
			reminder.ReminderDetails,
			reminder.ReminderDate,
		)
		formattedReminders = append(formattedReminders, formattedReminder)
	}

	// If no reminder has been found
	if len(formattedReminders) == 0 {
		return ReminderGetterTag, SelectRandomMessage(locale, "no reminders")
	}

	return ReminderGetterTag, fmt.Sprintf(response, strings.Join(formattedReminders, " "))
}

func SpotifySetterReplacer(locale, entry, _, token string) (string, string) {
	spotifyTokens := SearchTokens(entry)

	// Returns if the token is empty
	if len(spotifyTokens) != 2 {
		return SpotifySetterTag, SelectRandomMessage(locale, "spotify tokens")
	}

	// Save the tokens in the user's information
	UpdateUserProfile(token, func(information UserProfile) UserProfile {
		information.StreamingID = spotifyTokens[0]
		information.StreamingSecret = spotifyTokens[1]

		return information
	})

	return SpotifySetterTag, LoginSpotify(locale, token)
}

func SpotifyPlayerReplacer(locale, entry, response, token string) (string, string) {
	// Return if the tokens are not set
	if CheckTokensPresence(token) {
		return SpotifySetterTag, SelectRandomMessage(locale, "spotify credentials")
	}

	// Renew the spotify token and get the client
	client := RenewSpotifyToken(token)

	// Search for the track
	music, artist := SearchMusic(locale, entry)
	track, err := SearchTrack(client, music+" "+artist)
	if err != nil {
		return SpotifySetterTag, LoginSpotify(locale, token)
	}

	// Search if there is a device name in the entry
	device := SearchDevice(client, entry)
	options := &spotify.PlayOptions{
		URIs: []spotify.URI{track.URI},
	}

	// Add the device ID if a device is contained
	if device != (spotify.PlayerDevice{}) {
		options.DeviceID = &device.ID
	}

	// Play the found track
	client.PlayOpt(options)
	client.Play()

	return SpotifyPlayerTag, fmt.Sprintf(response, track.Name, track.Artists[0].Name)
}

func SearchTrack(client spotify.Client, content string) (spotify.FullTrack, error) {
	// Get the results from a track search with the given content
	results, err := client.Search(content, spotify.SearchTypeTrack)
	if err != nil {
		return spotify.FullTrack{}, err
	}

	// Returns an empty track and empty error if no track was found with this content
	if len(results.Tracks.Tracks) == 0 {
		return spotify.FullTrack{}, nil
	}

	// Return the found
	return results.Tracks.Tracks[0], nil
}

func SearchDevice(client spotify.Client, content string) spotify.PlayerDevice {
	devices, _ := client.PlayerDevices()

	// Iterate through the devices to check if the content contains a device name
	for _, device := range devices {
		if strings.Contains(content, strings.ToLower(device.Name)) ||
			strings.Contains(content, strings.ToLower(device.Type)) {
			return device
		}
	}

	return spotify.PlayerDevice{}
}

func init() {
	RegisterModulesf("en", []Modulef{
		// AREA
		// For modules related to countries, please add the translations of the countries' names
		// or open an issue to ask for translations.

		{
			Tag: AreaTag,
			Patterns: []string{
				"What is the area of ",
				"Give me the area of ",
			},
			Responses: []string{
				"The area of %s is %gkm²",
			},
			Replacer: AreaReplacer,
		},

		// CAPITAL
		{
			Tag: CapitalTag,
			Patterns: []string{
				"What is the capital of ",
				"What's the capital of ",
				"Give me the capital of ",
			},
			Responses: []string{
				"The capital of %s is %s",
			},
			Replacer: CapitalReplacer,
		},

		// CURRENCY
		{
			Tag: CurrencyTag,
			Patterns: []string{
				"Which currency is used in ",
				"Give me the used currency of ",
				"Give me the currency of ",
				"What is the currency of ",
			},
			Responses: []string{
				"The currency of %s is %s",
			},
			Replacer: CurrencyReplacer,
		},

		// MATH
		// A regex translation is also required in `language/math.go`, please don't forget to translate it.
		// Otherwise, remove the registration of the Math module in this file.

		{
			Tag: MathTag,
			Patterns: []string{
				"Give me the result of ",
				"Calculate ",
			},
			Responses: []string{
				"The result is %s",
				"That makes %s",
			},
			Replacer: MathReplacer,
		},

		// MOVIES
		// A translation of movies genres is also required in `language/movies.go`, please don't forget
		// to translate it.
		// Otherwise, remove the registration of the Movies modules in this file.

		{
			Tag: GenresTag,
			Patterns: []string{
				"My favorite movie genres are Comedy, Horror",
				"I like the Comedy, Horror genres",
				"I like movies about War",
				"I like Action movies",
			},
			Responses: []string{
				"Great choices! I saved this movie genre information to your client.",
				"Understood, I saved this movie genre information to your client.",
			},
			Replacer: GenresReplacer,
		},

		{
			Tag: MoviesTag,
			Patterns: []string{
				"Find me a movie about",
				"Give me a movie about",
				"Find me a film about",
			},
			Responses: []string{
				"I found the movie “%s” for you, which is rated %.02f/5",
				"Sure, I found this movie “%s”, which is rated %.02f/5",
			},
			Replacer: MovieSearchReplacer,
		},

		{
			Tag: MoviesAlreadyTag,
			Patterns: []string{
				"I already saw this movie",
				"I have already watched this film",
				"Oh I have already watched this movie",
				"I have already seen this movie",
			},
			Responses: []string{
				"Oh I see, here's another one “%s” which is rated %.02f/5",
			},
			Replacer: MovieSearchReplacer,
		},

		{
			Tag: MoviesDataTag,
			Patterns: []string{
				"I'm bored",
				"I don't know what to do",
			},
			Responses: []string{
				"I propose you watch the %s movie “%s”, which is rated %.02f/5",
			},
			Replacer: MovieSearchFromInformationReplacer,
		},

		// NAME
		{
			Tag: NameGetterTag,
			Patterns: []string{
				"Do you know my name?",
			},
			Responses: []string{
				"Your name is %s!",
			},
			Replacer: NameGetterReplacer,
		},

		{
			Tag: NameSetterTag,
			Patterns: []string{
				"My name is ",
				"You can call me ",
			},
			Responses: []string{
				"Great! Hi %s",
			},
			Replacer: NameSetterReplacer,
		},

		// RANDOM
		{
			Tag: RandomTag,
			Patterns: []string{
				"Give me a random number",
				"Generate a random number",
			},
			Responses: []string{
				"The number is %s",
			},
			Replacer: RandomNumberReplacer,
		},

		// REMINDERS
		// Translations are required in `language/date/date`, `language/date/rules` and in `language/reason`,
		// please don't forget to translate it.
		// Otherwise, remove the registration of the Reminders modules in this file.

		{
			Tag: ReminderSetterTag,
			Patterns: []string{
				"Remind me to cook a breakfast at 8pm",
				"Remind me to call mom tuesday",
				"Note that I have an exam",
				"Remind me that I have a conference call tomorrow at 9pm",
			},
			Responses: []string{
				"Noted! I will remind you: “%s” for the %s",
			},
			Replacer: ReminderSetterReplacer,
		},

		{
			Tag: ReminderGetterTag,
			Patterns: []string{
				"What did I ask for you to remember",
				"Give me my reminders",
			},
			Responses: []string{
				"You asked me to remember those things:\n%s",
			},
			Replacer: ReminderGetterReplacer,
		},

		// SPOTIFY
		// A translation is needed in `language/music`, please don't forget to translate it.
		// Otherwise, remove the registration of the Spotify modules in this file.

		{
			Tag: SpotifySetterTag,
			Patterns: []string{
				"Here are my spotify tokens",
				"My spotify secrets",
			},
			Responses: []string{
				"Login in progress",
			},
			Replacer: SpotifySetterReplacer,
		},

		{
			Tag: SpotifyPlayerTag,
			Patterns: []string{
				"Play from on Spotify",
			},
			Responses: []string{
				"Playing %s from %s on Spotify.",
			},
			Replacer: SpotifyPlayerReplacer,
		},

		{
			Tag: JokesTag,
			Patterns: []string{
				"Tell me a joke",
				"Make me laugh",
			},
			Responses: []string{
				"Here you go, %s",
				"Here's one, %s",
			},
			Replacer: JokesReplacer,
		},
		{
			Tag: AdvicesTag,
			Patterns: []string{
				"Give me an advice",
				"Advise me",
			},
			Responses: []string{
				"Here you go, %s",
				"Here's one, %s",
				"Listen closely, %s",
			},
			Replacer: AdvicesReplacer,
		},
	})

	// COUNTRIES
	// Please translate this method for adding the correct article in front of countries names.
	// Otherwise, remove the countries modules from this file.

	ArticleCountries["en"] = ArticleCountriesOut
}

func ArticleCountriesOut(name string) string {
	if name == "United States" {
		return "the " + name
	}

	return name
}
