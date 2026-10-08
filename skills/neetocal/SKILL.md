---
name: neetocal
description: >
  Manage NeetoCal calendar events, bookings, meetings, and scheduling from the command line.
  Use when the user asks about their calendar, scheduling, bookings, availabilities, or
  meetings managed through NeetoCal.
---

## Prerequisites

Run `neetocal doctor` to check authentication and connectivity.
If not authenticated, run `neetocal login`.

## Authentication & multi-subdomain

Credentials for every logged-in subdomain are stored together in
`~/.config/neetocal/auth.json`. A command that talks to the API picks which
subdomain to use by these rules:

- 0 subdomains authenticated → every credential-using command errors with
  "Not authenticated. Run 'neetocal login' to authenticate.".
- 1 subdomain authenticated → that one is the implicit default; `--subdomain`
  may be omitted.
- 2+ subdomains authenticated → **`--subdomain <name>` is required** on every
  credential-using command, including `doctor`. The error lists every
  authenticated subdomain so the agent can offer a choice.

`login` / `logout` / `whoami` have dedicated behavior:

| Command | Behavior |
|---|---|
| `neetocal login --subdomain <name>` | Adds or refreshes the entry for `<name>`. No flag → prompts for the subdomain. |
| `neetocal logout --subdomain <name>` | Removes that one entry. |
| `neetocal logout --all` | Removes every entry. |
| `neetocal logout` (no flag) | Removes the only entry if exactly one is logged in; errors if multiple. |
| `neetocal whoami` | Lists every logged-in account. Marks the entry `(default)` when exactly one. |
| `neetocal whoami --subdomain <name>` | Shows just that one. |

## Global flags (persistent on every command)

| Flag | Purpose |
|---|---|
| `--subdomain <name>` | Select which logged-in subdomain the command targets. Required when multiple are logged in. |
| `--json` | Force JSON envelope output even on a TTY. |
| `--quiet` | Emit only the raw payload — no envelope, no breadcrumbs. For action commands (create/update), emits just the resource identifier; `delete` emits `success`. Designed for scripting. |
| `--toon` | Emit TOON (Token Optimized Output Notation). Preferred for feeding list/show/slots output back to an LLM; ~30–60% fewer tokens than JSON. |

Precedence if multiple are set: `--toon` > `--quiet` > `--json` > pretty.

## Output modes & response envelope

**Pretty (default on a TTY)** — tables for arrays, key-value for objects,
breadcrumbs appended. Not intended for machine consumption.

**JSON envelope** (non-TTY, or `--json`):
```json
{
  "data": <resource body>,
  "breadcrumbs": [{ "label": "List", "command": "neetocal meetings list" }],
  "pagination": {
    "current_page_number": 1,
    "total_pages": 10,
    "total_records": 250
  }
}
```
`breadcrumbs` is omitted when empty. `pagination` is present only for list
commands.

**Quiet** (`--quiet`) — `data` contents only, no envelope. For action
commands `PrintQuiet` unwraps a single-key wrapper and prints the first of
`sid` / `id` / `name`. For `delete` it prints `success`.

**TOON** (`--toon`) — same data as JSON, re-encoded into TOON. Shape is
equivalent but whitespace/keys are compressed. Parse by re-reading keys as
you would JSON.

### Resource wrapping

API responses usually wrap the body under a resource key:
`{"meeting": { ... }}`, `{"bookings": [ ... ], "pagination": { ... }}`.
The CLI preserves the wrapper in JSON/TOON mode. In `--quiet` mode the
wrapper is unwrapped for identifier extraction.

### Pagination

List commands accept `--page` (1-indexed) and `--page-size` (max 100).
The envelope's `pagination` field always exposes:
`current_page_number`, `total_pages`, `total_records`. Agents should loop
by incrementing `--page` until `current_page_number == total_pages`.

## Discovery

The full, always-accurate command tree (including any flags added after
this skill was built) is available as JSON:

```bash
neetocal commands
```

Each catalog entry has `command`, `description`, optional `flags` (with
`name`, `type`, `default`, `description`, `required`), and `subcommands`.
Use this whenever a user asks about a flag or command not covered below.

## Command reference

Required flags are marked with `*`. All list commands also accept
`--page` / `--page-size`.

### Meetings

| Command | Positional | Flags |
|---|---|---|
| `meetings list` | — | `--host-email`, `--search` |
| `meetings show` | `<sid>` | — |
| `meetings create` | — | `--name*`, `--slug*`, `--hosts*` (csv emails), `--kind` (e.g., `one_on_one`), `--spot` (e.g., `zoom`), `--duration*` (minutes, int), `--description`, `--json-file` |
| `meetings update` | `<sid>` | `--name`, `--slug`, `--description`, `--hosts` (csv emails, optional — omit to keep current hosts), `--kind`, `--spot`, `--duration` (int), `--json-file` (partial) |
| `meetings delete` | `<sid>` | — |
| `meetings slots` | `<meeting-sid>` | `--year*` (int), `--month*` (int 1-12), `--day` (int), `--time-zone*`, `--host-email` (only that host's slots on a round-robin or multi-host meeting), `--override-availability` (every future time, including times outside availability; host or admin only) |
| `meetings one-off-link` | `<meeting-sid>` | — |
| `meetings durations list` | `<meeting-sid>` | — |
| `meetings durations show` | `<meeting-sid> <id>` | — |
| `meetings durations create` | `<meeting-sid>` | `--duration*` (int), `--is-default` (bool) |
| `meetings durations update` | `<meeting-sid> <id>` | `--duration`, `--is-default` (partial) |
| `meetings durations delete` | `<meeting-sid> <id>` | — |
| `meetings spots list` | `<meeting-sid>` | — |
| `meetings spots show` | `<meeting-sid> <id>` | — |
| `meetings spots create` | `<meeting-sid>` | `--spot*`, `--is-default`, `--phone-number`, `--location`, `--custom-text` |
| `meetings spots update` | `<meeting-sid> <id>` | same as create (partial) |
| `meetings spots delete` | `<meeting-sid> <id>` | — |
| `meetings calendar-preferences show` | `<meeting-sid> <integration>` (`google_calendar`, `outlook`, `icloud`) | — |
| `meetings calendar-preferences update` | `<meeting-sid> <integration>` | `--override-calendars`, `--event-add-calendar-ids` (csv, in order), `--override-conflict-check-calendars`, `--conflict-check-calendar-ids` (csv), `--busy-statuses` (csv, Outlook only), `--override-event-layout`, `--summary-type`, `--custom-summary`, `--body`, `--event-color` (Google only), `--json-file` (partial) |

**Response shapes**

- `meetings list` → `{"data": {"meetings": [{"sid","name","slug","kind","duration","spot","host_email",…}], "pagination":…}}`.
- `meetings show` → `{"data": {"meeting": {"sid","name","slug","description","kind","duration","spot","hosts":[…], "availability_id",…}}}`.
- `meetings create|update` → `{"data": {"meeting": {…}}}`. `--quiet` prints the `sid` only.
- `meetings delete` → empty body (204). `--quiet` prints `success`.
- `meetings slots` → `{"data": {"slots": [{"date":"YYYY-MM-DD","start_time":"HH:MM","end_time":"HH:MM"}…]}}`. With `--override-availability` each slot also has `is_available` and `already_overridden` (another overridden booking already takes it).
- `meetings one-off-link` → `{"data": {"url": "https://…"}}`.
- `meetings calendar-preferences show|update` → `{"data": {"calendar_preference": {"integration","override_calendars","event_add_calendars":[{"id","name","is_primary","is_editable"}],"override_conflict_check_calendars","conflict_check_calendars":[…],"busy_statuses","override_event_layout","summary_type","custom_summary","body","event_color","available_calendars":[…]}}}`. Pick calendar IDs from `available_calendars`: these are the host's calendars, and the list is empty unless the meeting is one-on-one. `update` changes only the flags passed. Turn a switch off with `--override-calendars=false`, and clear a list with `--event-add-calendar-ids ""`. Calendar overrides work only on one-on-one meetings, and read-only calendars (`is_editable: false`) can't receive events. `summary_type` is one of `host_and_client`, `client_and_host`, `meeting_name` or `custom` (`custom` needs `--custom-summary`). Group meetings allow only `meeting_name` and `custom`, and keep `--override-event-layout` on.

### Bookings

| Command | Positional | Flags |
|---|---|---|
| `bookings list` | — | `--host-email`, `--client-email`, `--type` (`upcoming`/`past`/`cancelled`/`incomplete`), `--sorting-key` (`created_at`/`starts_at`), `--sorting-order` (`asc`/`desc`), `--search` (min 3 chars; matches client name/email/sid, meeting name, host name/email and form answers, e.g. `acme.com`), `--meeting-sid`, `--starts-after` / `--starts-before` (ISO 8601 datetime or `YYYY-MM-DD`) |
| `bookings show` | `<id>` | — |
| `bookings create` | — | `--meeting-slug*`, `--email*`, `--name*`, `--slot-date*` (YYYY-MM-DD), `--slot-start-time*` (HH:MM), `--time-zone*`, `--preferred-meeting-spot`, `--override-availability` (book outside the meeting's availability; host or admin only), `--host-email` (host to assign on a multi-host meeting; needs `--override-availability` unless clients can choose the host) |
| `bookings update` | `<id>` | `--status` (`cancelled`/`approved`/`rejected`), `--cancel-reason`, `--rejection-reason`, `--slot-date`, `--slot-start-time`, `--time-zone`, `--reschedule-reason`, `--preferred-meeting-spot` (only with `--slot-date`/`--slot-start-time`), `--meeting-outcome-id` (pass `""` to clear the outcome), `--name`, `--email` (reschedule reuses the existing booking's client details; pass `--name`/`--email` only to override), `--override-availability` (reschedule outside the meeting's availability; host or admin only; keeps the current host) (all partial — only flags the user sets are sent) |
| `bookings payments create` | `<booking-id>` | `--payment-provider*`, `--identifier`, `--discount-code` |
| `bookings payments update` | `<booking-id> <payment-id>` | `--payment-provider*`, `--status*` (`successful`/`rejected`), `--notes` |

**Response shapes**

- `bookings list` → `{"data": {"bookings": [{"id","status","starts_at","ends_at","time_zone","host_email","client_email","meeting_name",…}], "pagination":…}}`.
- `bookings show` → `{"data": {"booking": {"id","status","starts_at","ends_at","client":{…},"host":{…},"meeting":{…},…}}}`.
- `bookings create|update` → `{"data": {"booking": {…}}}`. `--quiet` prints the booking `id`. The booking has `is_overridden`; with `--override-availability` it also has `bypassed_rules`, the availability rules that were skipped.

### Team members

| Command | Positional | Flags |
|---|---|---|
| `team-members list` | — | `--email` |
| `team-members show` | `<id>` | — |
| `team-members create` | — | `--emails*` (csv), `--organization-role`, `--invited-by`, `--send-invitation-email` (bool), `--json-file` |
| `team-members update` | `<id>` | `--email`, `--first-name`, `--last-name`, `--time-zone`, `--organization-role`, `--json-file` (partial) |
| `team-members delete` | `<id>` | — |
| `team-members slots` | — | `--emails*` (csv), `--duration*` (minutes, int), `--start-date*` (YYYY-MM-DD), `--end-date*` (YYYY-MM-DD, at most 31 days after start), `--time-zone*` |

**Response shapes**

- `team-members slots` → `{"data": {"slots": [{"starts_at":"2026-10-06T10:00:00-04:00","ends_at":"2026-10-06T10:45:00-04:00"}…]}}`. Each slot is a start time when every listed person is free for the whole duration. Not paginated.

### Availabilities

| Command | Positional | Flags |
|---|---|---|
| `availabilities list` | — | `--emails` (csv) |
| `availabilities show` | `<id>` | — |
| `availabilities create` | — | `--email*`, `--name*`, `--time-zone`, `--json-file` (contains `periods` and `overrides` — required for non-trivial cases) |
| `availabilities update` | `<id>` | `--name`, `--json-file` (partial) |

`periods` / `overrides` JSON example:
```json
{
  "periods": [
    {"wday": "monday",    "start_time": "09:00", "end_time": "17:00"},
    {"wday": "tuesday",   "start_time": "09:00", "end_time": "17:00"}
  ],
  "overrides": [
    {"date": "2026-04-20", "start_time": "00:00", "end_time": "00:00", "disabled": true}
  ]
}
```

### Meeting templates

| Command | Positional | Flags |
|---|---|---|
| `meeting-templates list` | — | `--host-email`, `--search` |
| `meeting-templates show` | `<id>` | — |
| `meeting-templates create` | — | `--name*`, `--slug*`, `--hosts` (csv), `--kind`, `--spot`, `--duration*` (int), `--json-file` |
| `meeting-templates update` | `<id>` | `--name`, `--slug`, `--json-file` (partial) |
| `meeting-templates delete` | `<id>` | — |

### Automation rules

| Command | Positional | Flags |
|---|---|---|
| `automation-rules list` | — | — |
| `automation-rules create` | — | `--event*` (e.g., `booking_confirmed`, `booking_cancelled`), `--name`, `--json-file*` (contains `meeting_ids` and `actions`) |
| `automation-rules delete` | `<id>` | — |

Typical `--json-file` payload:
```json
{
  "meeting_ids": ["mtg_abc", "mtg_def"],
  "actions": [
    {"type": "send_email", "template_id": "tmpl_123"},
    {"type": "add_to_calendar"}
  ]
}
```

### Packages & discount codes

| Command | Positional | Flags |
|---|---|---|
| `packages list` | — | — |
| `packages show` | `<id>` | — |
| `packages purchases list` | `<package-id>` | — |
| `packages purchases show` | `<package-id> <purchase-id>` | — |
| `discount-codes create` | — | `--code*`, `--kind*` (`percentage`/`fixed`), `--value*` (float), `--meeting-ids` (csv), `--expires-at` (ISO date) |

### Diagnostics & IDE setup

| Command | Purpose |
|---|---|
| `doctor` | Auth check + API reachability + version. Uses `--subdomain` when multiple are logged in. |
| `version` | Print CLI version / commit / build date. |
| `update` | Update the CLI to the latest version (auto-detects brew / shell / PowerShell install). |
| `commands` | Emit the full command/flag catalog as JSON. |
| `setup claude` | Install NeetoCal plugin into Claude Code (`plugin.json`, hooks, this SKILL.md). |
| `setup cursor` / `windsurf` / `copilot` / `gemini` / `codex` | Write NeetoCal rule files into the current project directory; re-run after an upgrade to refresh them. |

## Common workflows

### Authenticate into a new tenant without losing the current one
```bash
neetocal login --subdomain acme      # first tenant (becomes default)
neetocal meetings list               # works — only one logged in
neetocal login --subdomain beta      # add a second tenant
neetocal meetings list               # errors: --subdomain required
neetocal meetings list --subdomain acme --toon
```

### Book a slot end-to-end
```bash
# 1. Find the meeting's slug / sid
neetocal meetings list --search "demo" --toon

# 2. Pull available slots for a month (or specific day)
neetocal meetings slots mtg_abc123 --year 2026 --month 5 --day 15 \
  --time-zone "America/New_York" --toon

# 3. Create the booking using a slot from step 2
neetocal bookings create \
  --meeting-slug demo \
  --name "John Doe" \
  --email john@example.com \
  --slot-date 2026-05-15 \
  --slot-start-time 10:00 \
  --time-zone "America/New_York" --quiet
# → prints booking id on stdout
```

### Bulk-approve pending bookings
```bash
neetocal bookings list --type incomplete --quiet \
  | jq -r '.bookings[].id' \
  | xargs -I {} neetocal bookings update {} --status approved --quiet
```

### Reschedule a confirmed booking
```bash
neetocal bookings update bkg_123 \
  --slot-date 2026-05-22 \
  --slot-start-time 14:00 \
  --time-zone "America/New_York" \
  --reschedule-reason "Client conflict" --quiet
```

### Book or reschedule outside availability
Only a host of the meeting or an admin can do this. Past times are still refused, and a host can't have two overridden bookings at the same time.
```bash
# Times that can be booked with the override (taken ones show already_overridden)
neetocal meetings slots mtg_abc123 --year 2026 --month 10 --day 5 \
  --time-zone "America/New_York" --override-availability --toon

neetocal bookings create --meeting-slug product-demo \
  --name "Eve Smith" --email eve@example.com \
  --slot-date 2026-10-05 --slot-start-time 20:00 \
  --time-zone "America/New_York" --override-availability

neetocal bookings update bkg_123 --slot-date 2026-10-05 --slot-start-time 21:00 \
  --time-zone "America/New_York" --override-availability --quiet
```

### Find bookings for a domain in a date range
```bash
neetocal bookings list --search acme.com \
  --starts-after 2026-10-01 --starts-before 2026-11-01 --toon
```

### Count one meeting's bookings last month
```bash
neetocal bookings list --meeting-sid product-demo --type past \
  --starts-after 2026-09-01 --starts-before 2026-10-01 --json \
  | jq '.pagination.total_records'
```

### Find a common free slot for several people
```bash
neetocal team-members slots \
  --emails oliver@example.com,sam@example.com \
  --duration 45 --start-date 2026-10-05 --end-date 2026-10-09 \
  --time-zone "America/New_York" --toon
```

### Cancel with a reason
```bash
neetocal bookings update bkg_123 --status cancelled --cancel-reason "No longer needed" --quiet
```

### Run a report across multiple tenants
```bash
for sub in acme beta gamma; do
  echo "=== $sub ==="
  neetocal bookings list --subdomain "$sub" --type upcoming --quiet \
    | jq '.bookings | length'
done
```

### Create a one-off meeting link
```bash
neetocal meetings one-off-link mtg_abc123
```

### Send a meeting's bookings to a specific calendar
```bash
# 1. Find the host's calendar IDs
neetocal meetings calendar-preferences show mtg_abc123 google_calendar --toon

# 2. Add booking events to those calendars and check only one of them for conflicts
neetocal meetings calendar-preferences update mtg_abc123 google_calendar \
  --override-calendars --event-add-calendar-ids "<work-id>,<personal-id>" \
  --override-conflict-check-calendars --conflict-check-calendar-ids "<work-id>"
```

### Provision an availability from JSON
```bash
neetocal availabilities create \
  --email alice@co.com --name "Work Hours" --time-zone "America/New_York" \
  --json-file periods.json --quiet
```

### Pick the right output mode

| Goal | Flag |
|---|---|
| Show a list/show/slots result to the user or feed it back to the LLM | `--toon` |
| Pipe an identifier into another command or into `jq` | `--quiet` |
| Produce a machine-readable payload (with pagination metadata) | `--json` |
| Interactive TTY browsing | no flag |

## Error surface

Every command exits non-zero on failure and writes a single-line message to
stderr. Common errors the agent should expect:

- `Not authenticated. Run 'neetocal login' to authenticate.` — empty credential store.
- `Multiple subdomains authenticated (acme, beta); specify --subdomain.` — pick one.
- `Not authenticated for "foo". Authenticated subdomains: acme, beta.` — bad `--subdomain`.
- `required flag(s) "xxx" not set` (from cobra) — missing required flag.
- API errors come through with the server's message body; inspect the
  JSON envelope (or the `--quiet` payload) for `error` / `errors` / `notice`
  keys and any suggestions the API returns.

## Conventions

- Dates: `YYYY-MM-DD`. Times: `HH:MM` (24-hour). Time zones: IANA names (`America/New_York`).
- IDs: `meetings` uses short IDs (`sid`); `bookings`, `availabilities`,
  `packages`, `discount-codes`, `meeting-templates`, `automation-rules` use
  `id`. `PrintQuiet` picks whichever is present.
- For any flag or field not covered above, `neetocal commands` is
  authoritative.
