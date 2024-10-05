package olivia

import (
	"github.com/gorilla/websocket"
	gocache "github.com/patrickmn/go-cache"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"net/http"
	"os"
	"time"
)

// =================================================================
var cachedDataStore = map[string][]DataPacket{}

var cachedUserData = map[string]UserProfile{}

var (

	// globalNeuralNetworks is a map to hold the neural network instances
	globalNeuralNetworks map[string]Network

	// cacheInstance initializes the cache with a 5-minute lifetime
	cacheInstance = gocache.New(5*time.Minute, 5*time.Minute)
)

var websocketUpgrader = websocket.Upgrader{
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

// =================================================================
const adviceURL = "https://api.adviceslip.com/advice"
const day = time.Hour * 24
const jokeURL = "https://official-joke-api.appspot.com/random_joke"
const DontUnderstand = "don't understand"

// =================================================================
