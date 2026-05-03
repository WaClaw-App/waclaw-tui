package i18n

import (
        "fmt"
        "sync"
        "time"
)

// Locale represents a supported UI language.
type Locale string

const (
        LocaleID Locale = "id" // casual Indonesian (default)
        LocaleEN Locale = "en" // casual English
)

var (
        currentLocale Locale = LocaleID
        mu            sync.RWMutex
)

// SetLocale changes the active display language.
func SetLocale(l Locale) {
        mu.Lock()
        currentLocale = l
        mu.Unlock()
}

// GetLocale returns the current active locale.
func GetLocale() Locale {
        mu.RLock()
        defer mu.RUnlock()
        return currentLocale
}

// T looks up a display string by key in the current locale.
// Falls back to English if key not found in current locale.
// Returns the key itself if not found in any locale.
func T(key string) string {
        mu.RLock()
        defer mu.RUnlock()

        // Try current locale first.
        if m, ok := locales[currentLocale]; ok {
                if val, ok := m[key]; ok {
                        return val
                }
        }

        // Fallback to English.
        if val, ok := en[key]; ok {
                return val
        }

        // Key not found — return the key itself as last resort.
        return key
}

// SupportedLocales returns all supported locale codes.
func SupportedLocales() []Locale {
        return []Locale{LocaleID, LocaleEN}
}

// IsValidLocale checks if a locale code is supported.
func IsValidLocale(l Locale) bool {
        for _, supported := range SupportedLocales() {
                if l == supported {
                        return true
                }
        }
        return false
}

// DayLabels returns the 7 weekday abbreviations (Monday–Sunday) for the
// current locale. This is a frontend concern — locale-aware display data
// must not be hardcoded in screen packages.
func DayLabels() []string {
        mu.RLock()
        defer mu.RUnlock()

        if currentLocale == LocaleID {
                return []string{"senin", "selasa", "rabu", "kamis", "jumat", "sabtu", "minggu"}
        }
        return []string{"mon", "tue", "wed", "thu", "fri", "sat", "sun"}
}

// FormatDate returns a locale-aware long date string (e.g. "selasa, 30 april 2024"
// for ID locale or "Tuesday, 30 April 2024" for EN locale).
func FormatDate(t time.Time) string {
        mu.RLock()
        defer mu.RUnlock()

        if currentLocale == LocaleID {
                // Indonesian: "selasa, 30 april 2024"
                dayName := DayLabels()[int(t.Weekday())-1]
                if t.Weekday() == time.Sunday {
                        dayName = DayLabels()[6]
                }
                monthLower := map[time.Month]string{
                        time.January: "januari", time.February: "februari", time.March: "maret",
                        time.April: "april", time.May: "mei", time.June: "juni",
                        time.July: "juli", time.August: "agustus", time.September: "september",
                        time.October: "oktober", time.November: "november", time.December: "desember",
                }
                m := monthLower[t.Month()]
                if m == "" {
                        m = t.Month().String()
                }
                return fmt.Sprintf("%s, %d %s %d", dayName, t.Day(), m, t.Year())
        }
        // English: Go's built-in formatting.
        return t.Format("Monday, 2 January 2006")
}

// locales maps each locale code to its string table.
// Populated by init() to avoid circular reference issues.
var locales map[Locale]map[string]string

func init() {
        locales = map[Locale]map[string]string{
                LocaleID: id,
                LocaleEN: en,
        }
}
