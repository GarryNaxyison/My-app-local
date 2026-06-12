package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPrepareWebAppHTMLForMobileShell(t *testing.T) {
	html := []byte(`<html lang="ru" data-app-shell="auto" data-app-device="desktop"><head></head>`)
	got := string(prepareWebAppHTML(html, "mobile"))
	if !strings.Contains(got, `data-app-shell="mobile" data-app-device="mobile"`) {
		t.Fatalf("mobile shell was not injected: %s", got)
	}
}

func TestPrepareWebAppHTMLDefaultsToDesktopShell(t *testing.T) {
	html := []byte(`<html lang="ru" data-app-shell="auto" data-app-device="desktop"><head></head>`)
	got := string(prepareWebAppHTML(html, "unknown"))
	if !strings.Contains(got, `data-app-shell="desktop" data-app-device="desktop"`) {
		t.Fatalf("desktop shell was not injected: %s", got)
	}
}

func TestPrepareWebAppHTMLCanMarkStandaloneLogin(t *testing.T) {
	html := []byte(`<html lang="ru" data-app-shell="auto" data-app-device="desktop"><head></head>`)
	got := string(prepareWebAppHTMLWithMode(html, "mobile", true))
	if !strings.Contains(got, `data-auth-standalone="true"`) {
		t.Fatalf("standalone auth flag was not injected: %s", got)
	}
}

func TestWebAppHTMLLoadsBuiltReactApp(t *testing.T) {
	html, err := os.ReadFile(filepath.Join("web", "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	page := string(html)
	for _, want := range []string{
		`<div id="root">`,
		`poliglot-boot`,
		`type="module"`,
		`/app/assets/`,
		`.css`,
		`.js`,
	} {
		if !strings.Contains(page, want) {
			t.Fatalf("React web app shell is missing %q", want)
		}
	}
	for _, forbidden := range []string{
		`/src/main.tsx`,
		`activation-key-notice`,
		`showActivationKeyNotice`,
	} {
		if strings.Contains(page, forbidden) {
			t.Fatalf("React web app shell still contains legacy inline marker %q", forbidden)
		}
	}
}

func TestReactFrontendKeepsGoAPIContracts(t *testing.T) {
	appSource, err := os.ReadFile(filepath.Join("web-react", "src", "App.tsx"))
	if err != nil {
		t.Fatal(err)
	}
	apiSource, err := os.ReadFile(filepath.Join("web-react", "src", "lib", "api.ts"))
	if err != nil {
		t.Fatal(err)
	}
	page := string(appSource) + "\n" + string(apiSource)
	for _, want := range []string{
		`/api/session`,
		`/api/daily/claim`,
		`/api/auth/logout`,
		`/login`,
		`/api/ai-tutor/start`,
		`/api/ai-tutor/completed`,
		`/api/ai-tutor/restart`,
		`/api/ai-tutor/answer`,
		`/api/lesson/start`,
		`/api/lesson/answer`,
		`/api/practice`,
		`/api/shadowing/start`,
		`/api/shadowing/answer`,
		`/api/words/next`,
		`/api/word-game/next`,
		`/api/spelling/start`,
		`/api/level-test/start`,
		`/api/vocabulary`,
		`/api/phrasebook`,
		`/api/mistakes`,
		`/api/mistakes/practice/start`,
		`/api/mistakes/practice/answer`,
		`/api/mistakes/clear`,
		`/api/tools/voice-text`,
		`/api/tools/image-translate`,
		`/api/tools/translator`,
		`/api/leaderboard`,
		`/api/premium/plans`,
		`/api/premium/activation-key`,
		`/api/settings`,
		`telegram_account`,
		`referral_code`,
		`premium_audio_required`,
		`premium_pronunciation_required`,
		`complete_daily_first`,
		`cleanTutorHistoryLabel`,
		`to talk about`,
	} {
		if !strings.Contains(page, want) {
			t.Fatalf("React frontend is missing Go API contract marker %q", want)
		}
	}
	for _, want := range []string{
		"Listening and pronunciation",
		"AI Tutor, listening, pronunciation",
	} {
		if !strings.Contains(page, want) {
			t.Fatalf("React premium copy is missing %q", want)
		}
	}
	for _, forbidden := range []string{
		"unlinkTelegram",
		"data-unlink-telegram",
		"/api/auth/telegram/unlink",
	} {
		if strings.Contains(page, forbidden) {
			t.Fatalf("React frontend should not expose Telegram unlink action %q", forbidden)
		}
	}
}

func TestReactFrontendContainsTutorOverflowGuards(t *testing.T) {
	appSource, err := os.ReadFile(filepath.Join("web-react", "src", "App.tsx"))
	if err != nil {
		t.Fatal(err)
	}
	cssSource, err := os.ReadFile(filepath.Join("web-react", "src", "styles", "app.css"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appSource) + "\n" + string(cssSource)
	for _, want := range []string{
		"tutor-composer-v2 textarea",
		"tutor-feedback-v2",
		"tutor-completed-lesson-v2",
		"overflow-wrap: anywhere",
		"min-width: 0",
	} {
		if !strings.Contains(source, want) {
			t.Fatalf("React tutor overflow guard is missing %q", want)
		}
	}
}

func TestReactPWAUpdatesInstalledShellFromNetwork(t *testing.T) {
	appSource, err := os.ReadFile(filepath.Join("web-react", "src", "App.tsx"))
	if err != nil {
		t.Fatal(err)
	}
	workerSource, err := os.ReadFile(filepath.Join("web-react", "public", "offline-deck-sw.js"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appSource) + "\n" + string(workerSource)
	for _, want := range []string{
		`updateViaCache: "none"`,
		`const hadController = !!navigator.serviceWorker.controller`,
		`registration.update()`,
		`SKIP_WAITING`,
		`poliglot-v2-offline-decks-20260530`,
		`fetch(request)`,
		`request.mode === "navigate"`,
	} {
		if !strings.Contains(source, want) {
			t.Fatalf("PWA update flow is missing %q", want)
		}
	}
	if strings.Contains(string(workerSource), "return cached || network") {
		t.Fatal("service worker should not serve cached app shell before checking the network")
	}
}

func TestReactFrontendIncludesPremiumDashboardSurfaces(t *testing.T) {
	source, err := os.ReadFile(filepath.Join("web-react", "src", "App.tsx"))
	if err != nil {
		t.Fatal(err)
	}
	fallingPatternSource, err := os.ReadFile(filepath.Join("web-react", "src", "components", "ui", "falling-pattern.tsx"))
	if err != nil {
		t.Fatal(err)
	}
	shaderSource, err := os.ReadFile(filepath.Join("web-react", "src", "components", "ui", "shaders-hero-section.tsx"))
	if err != nil {
		t.Fatal(err)
	}
	chatSource, err := os.ReadFile(filepath.Join("web-react", "src", "components", "ui", "chat-interface.tsx"))
	if err != nil {
		t.Fatal(err)
	}
	page := string(source) + "\n" + string(fallingPatternSource) + "\n" + string(shaderSource) + "\n" + string(chatSource)
	for _, want := range []string{
		`data-v2-shell="ios-function-ribbon"`,
		"FunctionRibbon",
		"ContextHeader",
		"ThemeBackground",
		"FallingPattern",
		"ShaderBackground",
		"ChatComponent",
		"ChoiceTrainer",
		"ToolsView",
		"PremiumView",
		"SettingsView",
		"ActionCard",
		"translateViewDetails",
		"appCopy",
	} {
		if !strings.Contains(page, want) {
			t.Fatalf("React frontend is missing premium surface marker %q", want)
		}
	}
}

func TestReactFrontendKeepsExpandedInterfaceLocales(t *testing.T) {
	source, err := os.ReadFile(filepath.Join("web-react", "src", "lib", "i18n.ts"))
	if err != nil {
		t.Fatal(err)
	}
	page := string(source)
	for _, code := range []string{"ru", "en", "es", "de", "fr", "it", "zh", "ja", "ko", "tg", "uz", "tt", "hy", "kk", "ky", "ka", "uk", "pl", "ro", "pt", "ar", "bn", "cs", "el", "hi", "hu", "id", "nl", "sv", "ta", "te", "th", "tl", "tr", "vi"} {
		if !strings.Contains(page, `"`+code+`"`) {
			t.Fatalf("React v2 i18n is missing locale %q", code)
		}
	}
	for _, want := range []string{
		"interface_language",
		"learning_language",
		"view_home_title",
		"view_settings_title",
		"auth_login_title",
	} {
		if !strings.Contains(page, want) {
			t.Fatalf("React v2 i18n is missing key %q", want)
		}
	}
}

func TestWebLanguageDTOsUseCleanVisibleNames(t *testing.T) {
	languages := webLanguageDTOs(interfaceLanguages())
	if len(languages) != 35 {
		t.Fatalf("got %d interface languages, want 35", len(languages))
	}
	want := map[string]string{
		"ru": "\u0420\u0443\u0441\u0441\u043a\u0438\u0439",
		"es": "Espa\u00f1ol",
		"zh": "\u4e2d\u6587",
		"ja": "\u65e5\u672c\u8a9e",
		"ko": "\ud55c\uad6d\uc5b4",
		"hy": "\u0540\u0561\u0575\u0565\u0580\u0565\u0576",
		"ka": "\u10e5\u10d0\u10e0\u10d7\u10e3\u10da\u10d8",
		"pt": "Portugu\u00eas",
		"ar": "\u0627\u0644\u0639\u0631\u0628\u064a\u0629",
		"bn": "\u09ac\u09be\u0982\u09b2\u09be",
		"cs": "\u010ce\u0161tina",
		"el": "\u0395\u03bb\u03bb\u03b7\u03bd\u03b9\u03ba\u03ac",
		"hi": "\u0939\u093f\u0902\u0926\u0940",
		"ta": "\u0ba4\u0bae\u0bbf\u0bb4\u0bcd",
		"te": "\u0c24\u0c46\u0c32\u0c41\u0c17\u0c41",
		"th": "\u0e20\u0e32\u0e29\u0e32\u0e44\u0e17\u0e22",
		"tr": "T\u00fcrk\u00e7e",
		"vi": "Ti\u1ebfng Vi\u1ec7t",
	}
	for _, language := range languages {
		if expected := want[language.Code]; expected != "" && language.NativeName != expected {
			t.Fatalf("%s native name: got %q want %q", language.Code, language.NativeName, expected)
		}
		for _, value := range []string{language.Name, language.NativeName, language.InterfaceName} {
			if strings.Contains(value, "РЎ") || strings.Contains(value, "Рџ") || strings.Contains(value, "Г±") || strings.Contains(value, "дё") || strings.Contains(value, "ж—") || strings.Contains(value, "н•") {
				t.Fatalf("%s still has mojibake-looking visible name %q", language.Code, value)
			}
		}
	}
}

func TestPublicSiteReactBuildSupportsLandingAndLegalPages(t *testing.T) {
	for _, name := range []string{"poliglot-ai.html", "privacy.html", "terms.html"} {
		html, err := os.ReadFile(filepath.Join("Сайт полиглота для бота", name))
		if err != nil {
			t.Fatal(err)
		}
		page := string(html)
		for _, want := range []string{
			`<div id="root"></div>`,
			`/assets/site-react/`,
			`/assets/site-i18n.js`,
			`/assets/site-phrases.js`,
			`type="module"`,
		} {
			if !strings.Contains(page, want) {
				t.Fatalf("%s React static page is missing %q", name, want)
			}
		}
	}
	source, err := os.ReadFile(filepath.Join("site-react", "src", "PublicSiteApp.tsx"))
	if err != nil {
		t.Fatal(err)
	}
	legalSource, err := os.ReadFile(filepath.Join("site-react", "src", "legacyLegalContent.ts"))
	if err != nil {
		t.Fatal(err)
	}
	i18nSource, err := os.ReadFile(filepath.Join("site-react", "public", "assets", "site-i18n.js"))
	if err != nil {
		t.Fatal(err)
	}
	page := string(source) + "\n" + string(legalSource)
	for _, want := range []string{
		"SparklesCore",
		"nav-theme-toggle",
		"landing-hero",
		"pricing-grid",
		"termsDocumentHtml",
		"privacyDocumentHtml",
		"1. Принятие условий",
		"\u041f\u041e\u041b\u0418\u0422\u0418\u041a\u0410 \u041e\u0411\u0420\u0410\u0411\u041e\u0422\u041a\u0418 \u041f\u0415\u0420\u0421\u041e\u041d\u0410\u041b\u042c\u041d\u042b\u0425 \u0414\u0410\u041d\u041d\u042b\u0425",
	} {
		if !strings.Contains(page, want) {
			t.Fatalf("public React site source is missing %q", want)
		}
	}
	localePage := string(i18nSource)
	for _, code := range []string{"ru", "en", "es", "de", "fr", "it", "zh", "ja", "ko", "tg", "uz", "tt", "hy", "kk", "ky", "ka", "uk", "pl", "ro", "pt"} {
		if !strings.Contains(localePage, `"`+code+`"`) {
			t.Fatalf("public React site is missing locale %q", code)
		}
	}
}

func TestPrivacyPolicyAssetCoversTwentyLocalesAndBotConsent(t *testing.T) {
	for _, pathParts := range [][]string{
		{"site-react", "public", "assets", "privacy-policy-i18n.js"},
		{"Сайт полиглота для бота", "assets", "privacy-policy-i18n.js"},
	} {
		asset, err := os.ReadFile(filepath.Join(pathParts...))
		if err != nil {
			t.Fatal(err)
		}
		page := string(asset)
		for _, code := range []string{"ru", "en", "es", "de", "fr", "it", "zh", "ja", "ko", "tg", "uz", "tt", "hy", "kk", "ky", "ka", "uk", "pl", "ro", "pt"} {
			if !strings.Contains(page, `"`+code+`"`) {
				t.Fatalf("%s is missing locale %q", filepath.Join(pathParts...), code)
			}
		}
		for _, want := range []string{
			"poliglotai.online",
			"@AsaselD",
			"@poliglot_ai_bot",
			"Telegram ID",
			"\u041f\u041e\u041b\u0418\u0422\u0418\u041a\u0410 \u041e\u0411\u0420\u0410\u0411\u041e\u0422\u041a\u0418 \u041f\u0415\u0420\u0421\u041e\u041d\u0410\u041b\u042c\u041d\u042b\u0425 \u0414\u0410\u041d\u041d\u042b\u0425",
		} {
			if !strings.Contains(page, want) {
				t.Fatalf("%s is missing privacy marker %q", filepath.Join(pathParts...), want)
			}
		}
	}
}

func TestWebAssetContentTypesCoverReactBuild(t *testing.T) {
	cases := map[string]string{
		"index.js":   "text/javascript; charset=utf-8",
		"index.css":  "text/css; charset=utf-8",
		"font.woff2": "font/woff2",
		"logo.svg":   "image/svg+xml",
		"data.json":  "application/json; charset=utf-8",
	}
	for name, want := range cases {
		if got := webAssetContentType(name); got != want {
			t.Fatalf("content type for %s: got %q want %q", name, got, want)
		}
	}
}
