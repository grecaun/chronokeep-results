# Chronokeep Results

An API for timed event results. The Chronokeep desktop app uploads results to it as a race
goes on, and the Chronokeep website reads them back, including over a live stream while a
race is running.

## API

Most requests are authorized with an API key, sent as `Authorization: Bearer <key>`.
Account management uses the token from `POST /account/login` instead.

| Group | Routes |
|---|---|
| Results | `/results`, `/results/all`, `/results/multi`, `/results/finish`, `/results/bib`, `/results/add`, `/results/delete` |
| Live results | `GET /results/stream` (server-sent events) |
| Events and years | `/event`, `/event/all`, `/event/my`, `/event-year`, `/event-year/event`, and `add`, `update`, `delete` for each |
| Participants | `/participants`, `/participants/add`, `/participants/delete` |
| Distances, segments, bib chips | `/distances`, `/segments`, `/bibchips`, and `add`, `delete` for each |
| SMS | `/sms`, `/sms/add`, `/sms/remove`, `/twilio`, `/blocked/phones/...`, `/blocked/emails/...` |
| Registration | `/r/event`, `/r/event-year`, `/r/participants` and their `add`, `update`, `delete` routes |
| Accounts and keys | `/account/...`, `/key/...` |
| Health | `/health` |

### Live results stream

`GET /results/stream?slug=<event>&year=<year>&key=<key>` keeps a
[server-sent events](https://developer.mozilla.org/docs/Web/API/Server-sent_events) stream
open for an event year. The key goes in the query because browsers' `EventSource` can't send
headers. It sends:

- `ready` once connected, with the year being watched (useful when `year` is left out),
- `results` whenever that year's results are added or deleted, folding bursts into one event
  every two seconds,
- `ping` every 30 seconds, so clients can tell the stream is still working and reconnect if
  pings stop.

Streams are kept in memory, so they only hear about results uploaded to the same server.

## Configuration

Settings come from environment variables, or from a `.env` file in the working directory.

| Variable | Required | Description |
|---|---|---|
| `DB_CONNECTOR` | | `mysql` (default), `postgres` or `sqlite3`. |
| `DB_NAME` | Yes | Database name, or the file path for SQLite. |
| `DB_HOST` | Yes | Database host. |
| `DB_PORT` | Yes | Database port. |
| `DB_USER` | Yes | Database user. |
| `DB_PASSWORD` | | Database password. |
| `PORT` | | Port to listen on. Defaults to `8181`. |
| `VERSION` | | `production` for production. Anything else runs in development mode, which logs request and response bodies and doesn't require SSL to Postgres. |
| `AUTOTLS` | | `enabled` to serve HTTPS with a Let's Encrypt certificate. |
| `DOMAIN` | With `AUTOTLS` | The domain the certificate is for. |
| `SECRET_KEY` | Yes | Signs account tokens. At least 20 characters. |
| `REFRESH_KEY` | Yes | Signs refresh tokens. At least 20 characters. |
| `ADMIN_EMAIL`, `ADMIN_NAME`, `ADMIN_PASS` | First run | Used to create the first admin account when there are no accounts. |
| `TWILIO_AUTH_TOKEN`, `TWILIO_RESPONSE_WEBHOOK_URL` | For SMS | Twilio settings for text message alerts. |
| `BAN_ALLOWLIST` | | IPs or CIDRs, separated by commas, that are never banned for probing (see below). |

## Running

```sh
go build
./results
```

The database's tables are created on first run and upgraded automatically when a new
version needs it.

## Deployment

On the server the program lives in `/chronokeep/results/` and runs as the
`chronokeep-results` systemd service, writing its logs to `/chronokeep/results/logs/`.
`update.sh` installs the latest GitHub release there and restarts the service. Releases only
contain the program and `version.txt`.

Root's crontab:

- runs `update.sh` every hour,
- restarts the service at 00:01 each day, so each day's log is its own file,
- reloads the fail2ban jail just after, so it reads the new log file,
- syncs the logs and database backups to Google Drive each morning,
- archives the previous month's logs on the 1st with `zip-logs.sh`.

## Security

### Timeouts

Connections that send nothing, or stop reading, are closed: request headers must arrive within
10 seconds and a request within 60, a response must be written within 60 seconds, and idle
kept-alive connections close after 120. Live streams stay open on purpose, but each write to
one must finish within 15 seconds. On Linux, a connection whose data goes unacknowledged for
60 seconds is dropped, which ends streams to clients that vanished without closing them.

### Probe bans

Scanners regularly ask for files this API never serves, looking for something to exploit:
`.env` files, `.php` scripts, `.git` folders, WordPress and phpMyAdmin pages and so on. The first
such request from an IP gets a 404 and bans that IP for 24 hours; while banned, it gets a 403
for everything. Bans are checked before routing and request logging, and kept in memory, so a
restart clears them.

- The IP comes from the connection itself, not `X-Forwarded-For`, so it can't be faked to get
  another address banned.
- Loopback, private and link-local addresses are never banned, and neither is anything in
  `BAN_ALLOWLIST`. Add race venues and timing networks there: one infected device on a venue's
  shared Wi-Fi could otherwise get the venue's address banned, timers included.
- Each ban is logged once, as `Banned <ip> for 24h0m0s after it asked for "<path>".`

To block these IPs at the server's firewall as well, for longer and on every port, see
[deploy/fail2ban](deploy/fail2ban/README.md).

## Testing

```sh
go test ./...
```

The tests need MySQL and Postgres databases named `results_test` on `database.lan`, with user
and password `results_test`, and use a local `results_test.sqlite` file for SQLite. Every
database test fails if a connection is left checked out of the pool.

## License

GNU Affero General Public License v3. See [LICENSE](LICENSE).
