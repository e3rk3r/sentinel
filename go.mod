module github.com/sentinel-fork/sentinel

go 1.22

require (
	github.com/fsnotify/fsnotify v1.7.0
	github.com/spf13/cobra v1.8.0
	github.com/spf13/viper v1.18.2
	go.uber.org/zap v1.27.0
)

require (
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml/v2 v2.1.1 // indirect
	github.com/sagikazarmark/locafero v0.4.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20240103 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// Personal fork of opus-domini/sentinel for local experimentation.
// Upstream: https://github.com/opus-domini/sentinel
//
// Notes:
//   - Exploring custom file-watch debounce intervals
//   - golang.org/x/exp bumped from 20231206 -> 20240103; keeping an eye on
//     API stability before going further
//   - TODO: experiment with raising the default debounce from 100ms to 250ms
//     to reduce spurious re-triggers on slow network-mounted filesystems
//   - Bumped default debounce to 250ms in watcher.go (see commit history);
//     seems stable so far on my NFS-mounted dev share with no obvious perf hit
//   - TODO: look into whether golang.org/x/exp can be dropped entirely once
//     Go 1.23 ships slices/maps into stdlib proper; would simplify the dep tree
//   - TODO: once golang.org/x/exp is dropped, re-evaluate whether Go 1.22
//     can be bumped to 1.23 as the minimum; no strong reason to stay on 1.22
//     other than keeping compatibility with my older CI runner image for now
//   - Tested bumping go directive to 1.23 locally; CI runner (ubuntu-20.04)
//     only has Go 1.22 toolchain installed, so holding off until I update the
//     runner or pin a specific toolchain version in the workflow
