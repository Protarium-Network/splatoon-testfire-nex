# Splatoon Global Testfire NEX server

This is a preservation-oriented NEX server derived from
[PretendoNetwork/splatoon](https://github.com/PretendoNetwork/splatoon), under
the AGPL-3.0 license included in this repository. It changes the retail
Splatoon configuration to the Global Testfire configuration recovered from
`Gambit.rpx`:

- Game server ID: `0x1017E300`
- Access key: `da693ee5`
- NEX SDK: `3.8.3` (`NEX_3_8_3AGM`)

It exposes separate PRUDP authentication and secure endpoints. The secure
endpoint registers Ticket Granting, Secure Connection, NAT Traversal, Match
Making, Match Making Ext, Matchmake Extension, and Splatoon Ranking through
Pretendo's maintained NEX libraries.

## Important scope

The original Global Testfire's decoded network traces are not public. This
project implements the complete publicly documented NEX service set, but an
actual playable 4v4 match still needs validation with owned clients and packet
captures. In particular, client routing/account-token discovery, NNCS NAT
discovery, BOSS/event scheduling, and the exact Testfire RMC sequence are
outside this repository. It does not route a console to this server or interact
with Nintendo's production services.

Turf War gameplay itself is PIA peer-to-peer traffic. The NEX server brokers
identity, tickets, NAT traversal, gatherings, and session metadata; it does
not centrally simulate battle packets.

## Local setup

Prerequisites:

- Go 1.25 or newer, or Docker Compose
- PostgreSQL 16 or newer
- A local `settings.json` containing NEX PID/password pairs for test clients

Copy the templates, update the reachable secure-server address and local
account details, then start the server:

```text
copy .env.example .env
copy settings.example.json settings.json
go run .
```

For a containerized PostgreSQL and server, create `.env` from `.env.example`,
create `settings.json`, set `PN_GLOBAL_TESTFIRE_SECURE_SERVER_HOST` to the
server's LAN/public address, then run:

```text
docker compose up --build
```

`PN_GLOBAL_TESTFIRE_LOCAL_MODE=1` is for an isolated preservation environment.
It uses the local account registry and intentionally does not contact Pretendo
account/friends services. Disable it and provide the gRPC settings when
integrating with a real account/friends backend.

## Verification

```text
go test ./...
go build ./...
```

See [PROTOCOL_COVERAGE.md](PROTOCOL_COVERAGE.md) for registered NEX services,
their RMC IDs, and the interoperability work still requiring a Testfire trace.

_Deployed and maintained as part of the [Protarium Network](https://github.com/Protarium-Network) Wii U online service revival project._
_Derived from Pretendo Network’s original codebase (PretendoNetwork). Copyright (C) Pretendo Network contributors._
