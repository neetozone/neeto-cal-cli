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

## Quick Reference

```bash
# Auth
neetocal login                    # Log in via browser
neetocal whoami                   # Show current user
neetocal doctor                   # Health check

# Meetings
neetocal meetings list --toon
neetocal meetings show <sid> --toon
neetocal meetings create --name "Demo" --slug demo --hosts "alice@co.com" --kind one_on_one --spot zoom --duration 30
neetocal meetings delete <sid>

# Bookings
neetocal bookings list --type upcoming --toon
neetocal bookings show <id> --toon
neetocal bookings create --meeting-slug demo --name "John" --email john@example.com --slot-date 2025-01-15 --slot-start-time 10:00 --time-zone "America/New_York"
neetocal bookings update <id> --status cancelled --cancel-reason "Rescheduling"

# Available Slots
neetocal meetings slots <meeting-sid> --year 2025 --month 1 --time-zone "America/New_York" --toon

# Availabilities
neetocal availabilities list --toon
neetocal availabilities create --email alice@co.com --name "Work Hours" --time-zone "America/New_York" --json-file periods.json

# Meeting Templates
neetocal meeting-templates list --toon

# Packages
neetocal packages list --toon
neetocal packages purchases list <package-id> --toon

# Automation Rules
neetocal automation-rules list --toon
neetocal automation-rules create --event booking_confirmed --json-file rule.json

# Discount Codes
neetocal discount-codes create --code SAVE20 --kind percentage --value 20

# Discovery
neetocal commands --json          # Full command catalog
```

## Common Workflows

### Check upcoming schedule
```bash
neetocal bookings list --type upcoming --toon
```

### Book a meeting
```bash
# 1. Find available meetings
neetocal meetings list --toon

# 2. Check available slots
neetocal meetings slots <meeting-sid> --year 2025 --month 6 --day 15 --time-zone "America/New_York" --toon

# 3. Create the booking
neetocal bookings create \
  --meeting-slug <slug> \
  --name "John Doe" \
  --email john@example.com \
  --slot-date 2025-06-15 \
  --slot-start-time 10:00 \
  --time-zone "America/New_York"
```

### Cancel a booking
```bash
neetocal bookings update <id> --status cancelled --cancel-reason "No longer needed"
```

### Create a one-off meeting link
```bash
neetocal meetings one-off-link <meeting-sid>
```

## Token Efficiency

- Use `--toon` for all commands. TOON format uses 30-60% fewer tokens than JSON while preserving all fields.

```bash
neetocal meetings list --toon
neetocal bookings list --type upcoming --toon
neetocal meetings slots <sid> --year 2025 --month 6 --toon
neetocal meetings show <sid> --toon
```

## Important Notes

- Dates are in `YYYY-MM-DD` format
- Times are in `HH:MM` format (24-hour)
- Use `neetocal commands --json` to discover all available commands and flags
- For complex payloads (availability periods, automation actions), use `--json-file <path>`
- IDs can be either short IDs (sid) or UUIDs
