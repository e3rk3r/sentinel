# Tmux Workspace

![Desktop tmux sessions](assets/images/desktop-tmux-sessions.png)

Sentinel models tmux as:

- Session
- Window
- Pane

The UI and API are built to keep this hierarchy consistent with realtime updates.

## Core Capabilities

- List/create/rename/kill sessions.
- List/select/create/kill windows.
- List/select/split/kill panes.
- Attach to any session over WebSocket PTY stream.
- Rename window and pane labels.
- Session icon metadata.
- Frequent directories endpoint (`GET /api/tmux/frequent-dirs`) powers quick-pick suggestions in the session creation dialog.

![Desktop tmux fullscreen](assets/images/desktop-tmux-fullscreen.png)

## Realtime Interaction

- `/ws/tmux?session=<name>` streams the active tmux PTY.
- `/ws/events` carries projection updates for lists, unread state, and recovery status.

Mouse stability in browser terminals is enforced server-side by tmux binding patches and mouse-mode enablement.

## Optimistic UX

Frontend actions assume success first and reconcile with backend events.

Examples:

- Session create/kill updates UI immediately.
- Window create/kill/select is optimistic.
- Pane split/select/kill is optimistic.

If backend rejects the action, UI is corrected by subsequent API/event reconciliation.

## Default Naming Rules

When creating entities without custom names:

- New window: `win-<sequence>` where sequence is monotonic per session.
- New pane title: `pan-<pane-id-suffix>`.

This avoids repeated ambiguous names after tmux index reuse.

When creating a session with a name that already exists, the server auto-suffixes the name with `-1`, `-2`, ... up to `-9` to resolve the collision. The response `name` field may differ from the requested name.

## Multi-User Sessions

Sessions can be created as different OS users via the `user` field in the create payload. This allows a single Sentinel instance to manage sessions across multiple system accounts.

- Launchers support `userMode` (`session` or `fixed`) and `userValue` for per-launcher user targeting.
- The sidebar shows a user indicator when a session runs as a different user than the process user.
- Requires `[multi_user]` configuration — see Configuration Reference.

## Unread and Activity Semantics

Watchtower tracks revisions per pane and seen revisions per focus scope.

- Pane receives new output -> pane can become unread.
- Window is considered unread when any pane in it is unread.
- Session summary aggregates unread windows/panes.
- Unread activity is indicated by the session icon colour in the sidebar (amber), not by a window badge.

Seen operations happen via WS events channel (`type: "seen"`) and emit patch updates immediately.

## Sidebar Density

The sidebar adapts to 3 tiers based on available width:

- Minimal (<=250px): icon + name only.
- Compact (<=300px): badges visible.
- Full (>300px): content preview visible.
