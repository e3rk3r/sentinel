package watchtower

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"time"

	"github.com/opus-domini/sentinel/internal/activity"
	"github.com/opus-domini/sentinel/internal/store"
)

// fallbackErrorMarkers and fallbackWarnMarkers are used when no patterns are
// loaded from the store (e.g. during tests with a nil store).
var (
	fallbackErrorMarkers = []string{
		"panic",
		"fatal",
		"segmentation fault",
		"traceback",
		"exception",
		"xdebug",
		"permission denied",
		"error",
		"failed",
	}
	fallbackWarnMarkers = []string{
		"warning",
		"warn",
		"deprecated",
		"timeout",
		"retry",
		"slow",
	}
)

func normalizeRuntimeCommand(current, start string) string {
	command := strings.TrimSpace(current)
	if command == "" {
		command = strings.TrimSpace(start)
	}
	if command == "-" {
		return ""
	}
	return command
}

func isShellLikeCommand(command string) bool {
	command = strings.ToLower(strings.TrimSpace(command))
	switch command {
	case "", "sh", "bash", "zsh", "fish", "tmux":
		return true
	default:
		return false
	}
}

// detectTimelineMarker matches preview text against the provided patterns,
// returning the matched pattern string, its severity, and whether a match
// was found. Patterns are evaluated in priority order (as returned by the
// store); the first match wins.
func detectTimelineMarker(preview string, patterns []store.MarkerPattern) (string, string, bool) {
	normalized := strings.ToLower(strings.TrimSpace(preview))
	if normalized == "" {
		return "", "", false
	}

	// When patterns is non-nil (even if empty), use only the configured
	// patterns. A nil slice means the cache has not been populated yet,
	// so we fall back to the hardcoded defaults.
	if patterns != nil {
		for _, p := range patterns {
			if !p.Enabled {
				continue
			}
			needle := strings.ToLower(strings.TrimSpace(p.Pattern))
			if needle == "" {
				continue
			}
			if strings.Contains(normalized, needle) {
				return p.Pattern, p.Severity, true
			}
		}
		return "", "", false
	}

	// Fallback: use hardcoded lists when patterns cache is nil.
	for _, marker := range fallbackErrorMarkers {
		if strings.Contains(normalized, marker) {
			return marker, activity.SeverityError, true
		}
	}
	for _, marker := range fallbackWarnMarkers {
		if strings.Contains(normalized, marker) {
			return marker, activity.SeverityWarn, true
		}
	}
	return "", "", false
}

// refreshMarkerCache reloads marker patterns from the store if the cache
// is stale or empty. It is safe to call concurrently.
func (s *Service) refreshMarkerCache(ctx context.Context) {
	if s == nil || s.store == nil {
		return
	}

	s.markerMu.Lock()
	defer s.markerMu.Unlock()

	if len(s.markerCache) > 0 && time.Since(s.markerCacheAt) < s.markerCacheTTL {
		return
	}

	patterns, err := s.store.ListMarkerPatterns(ctx)
	if err != nil {
		slog.Warn("watchtower marker cache refresh failed", "err", err)
		return
	}
	s.markerCache = patterns
	s.markerCacheAt = time.Now()
}

// cachedMarkerPatterns returns the currently cached marker patterns.
func (s *Service) cachedMarkerPatterns() []store.MarkerPattern {
	if s == nil {
		return nil
	}
	s.markerMu.Lock()
	defer s.markerMu.Unlock()
	return s.markerCache
}

func timelineLastLine(preview string) string {
	preview = strings.TrimSpace(preview)
	if preview == "" {
		return ""
	}
	lines := strings.Split(preview, "\n")
	last := strings.TrimSpace(lines[len(lines)-1])
	if len(last) > 240 {
		return last[:240]
	}
	return last
}

func timelineMetadataJSON(values map[string]any) json.RawMessage {
	if values == nil {
		return json.RawMessage("{}")
	}
	payload, err := json.Marshal(values)
	if err != nil {
		return json.RawMessage("{}")
	}
	return json.RawMessage(payload)
}
