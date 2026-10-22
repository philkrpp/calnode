# Calnode

Go/SQLite application with an embedded SvelteKit admin. Public booking pages and the embed widget are separate from the admin SPA. Read [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) for system contracts; use the topical docs below when changing those areas.

- Frontend assets are `go:embed`-ed: after changing `frontend/`, rebuild it and the Go binary before checking the running app. Use the scripts in `frontend/package.json`.
- In the admin, extend existing shadcn-svelte components. Destructive actions use `ConfirmDialog`, never `window.confirm`/`alert`. After editing UI components, theme, or `app.css`, run `pnpm test:visual` from `frontend/`; unit tests miss the computed-style regressions described in `frontend/TESTING.md`.
- Calendar behavior exists separately in `book.html`, `manage.html`, and `embed.js`. Keep all three aligned; shared styles belong in `booking.css`. The embed is standalone and does not load `booking-logic.js`; mirror only the helpers it needs. Verify desktop and mobile on all affected surfaces. See Architecture §8.
- Availability, booking, and the chat assistant must use the deterministic slot and booking cores. The assistant may translate language into constraints; it must not perform time math or see raw calendar data. See Architecture §20.
- LiveKit authority comes from `authorizeHost` validation (durable host or verified current access token); a room token alone does not grant host actions. Finalize recordings on stop/end, enforce screen-share permissions on the server, and merge room metadata instead of overwriting flags. See Architecture §22 and [`docs/VIDEO.md`](docs/VIDEO.md).
- Keep `slot_interval_minutes` (start cadence) distinct from `duration_minutes`; validate effective changes and ensure defaults produce a hostable location. For mail, preserve `BuildMailer` as the transport-selection point, bounded SMTP dials, and the iCalendar `REQUEST` content type. See Architecture §§8–9, 12.
- Go owns locale-aware dates, pluralization, and completed sentence fragments. A new locale starts with one canonical JSON file and `go test ./internal/i18n/`; do not market machine-drafted translations as reviewed. See Architecture §23.
- For domain/base URL behavior, read [`docs/DOMAINS.md`](docs/DOMAINS.md); for Zoom team setup, read [`docs/ZOOM.md`](docs/ZOOM.md).
