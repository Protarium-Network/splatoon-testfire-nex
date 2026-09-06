# Protocol coverage

The server uses Pretendo's current Go NEX stack. The values below are the
Global Testfire / `ShowDL` values recovered from `Gambit.rpx`.

| Layer | Protocol ID | Status |
| --- | ---: | --- |
| PRUDPv1 + RMC framing | — | Provided by `nex-go` |
| Ticket Granting | `0x0A` | Registered on the authentication endpoint |
| Secure Connection | `0x0B` | Registered on the secure endpoint |
| NAT Traversal | `0x03` | Registered on the secure endpoint |
| Match Making | `0x15` | Registered on the secure endpoint |
| Match Making Ext | `0x32` | Registered on the secure endpoint |
| Matchmake Extension | `0x6D` | Registered on the secure endpoint |
| Splatoon Ranking | `0x70` | Registered on the secure endpoint |

## Testfire identity

| Setting | Value |
| --- | --- |
| Game server ID | `0x1017E300` |
| NEX access key | `da693ee5` |
| NEX SDK version | `3.8.3` (`NEX_3_8_3AGM`) |

The RPX's non-Testfire fallback is deliberately not used: it is the retail JP
configuration (`0x10162B00`, `6f599f81`, NEX 3.8.15).

## What the NEX server does

It authenticates NEX sessions, grants secure-server tickets, registers clients,
coordinates NAT traversal, creates/matches gatherings, and stores match state.
Actual Turf War traffic is PIA peer-to-peer traffic between consoles, so it is
not simulated by this process.

## Evidence still needed for end-to-end compatibility

No public decoded Global Testfire packet trace establishes the exact RMC method
order, match attributes, event-time content, or all error codes. This server
logs protocol and method IDs at both endpoints to support comparison with
captures from an owned, locally routed client. It intentionally does not
include client routing, Nintendo account token emulation, or a bypass for
Nintendo's production services.

The RPX also shows that a playable Internet session depends on services outside
this NEX process: account-token/auth-server discovery, NNCS-style NAT discovery
for PIA peer connections, and BOSS-provided schedule/stage content. Those need
their own compatible preservation services and captured Testfire data; no
implementation can safely infer the missing payloads from the RPX alone.
