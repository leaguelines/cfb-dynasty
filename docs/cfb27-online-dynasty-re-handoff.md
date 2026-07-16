# CFB27 Online Dynasty RE — Agent Handoff Brief

## Mission

Reverse-engineer **how College Football 27 authenticates to EA servers** and **how clients receive live league state updates** for online dynasty. Goal is academic understanding of distributed save architecture, with a possible POC for EADP token exchange (not necessarily full Blaze session yet).

**Constraints:**

- **No live packet capture** — Javelin anti-cheat blocks the game when Wireshark/similar tools run
- **Static analysis only** via Ghidra (+ optional `strings` on the binary)
- User is a systems engineer; wants architecture insight, not game modding

---

## Artifacts & Paths

| Item | Path |
|------|------|
| Binary | `/run/media/preston/Data/College Football 27/CollegeFB27.exe` |
| Ghidra project | `/run/media/preston/Data/College Football 27/CFB27.rep` |
| Ghidra version | 12.1.2 |
| Related Go project (local save export, separate concern) | `/home/preston/Documents/Projects/golang/cfb-dynasty` |

Binary facts: **PE32+ x86-64**, ~237 MB, EA sports title sharing Madden/Origin stack.

---

## Architecture Summary (confirmed via `strings`)

Online dynasty is **not file-sync**. It's:

```
Authoritative server league instance (FranTkMode)
  + Blaze RPC commands (mutations)
  + Server-pushed EventInfo notifications (updates)
  + EADP Persistence Service (cloud save blobs/metadata)
  + Client local mirror patched from notifications
```

### Two separate network stacks

| Stack | Protocol | Purpose |
|-------|----------|---------|
| **EADP Nexus** | Protobuf + gRPC | OAuth tokens, identity |
| **Blaze** | TDF over TLS/TCP (DirtySDK) | Game service session, league RPCs, push notifications |

### Known identifiers (from strings)

- Blaze service name: `COLLEGEFB_2027_PC`
- Product slug: `collegefb.27`
- Redirector: `spring25.client.blazeredirector.ea.com`
- Nexus gRPC host: `gateway.grpc.ea.com`
- Accounts: `https://accounts.ea.com`, `https://gateway.ea.com/proxy`

---

## Phase 1: Ghidra Setup

1. Open project `CFB27` at `/run/media/preston/Data/College Football 27/`
2. Open `CollegeFB27.exe` in **CodeBrowser**
3. Wait for auto-analysis to finish
4. Enable demangler (symbols like `PcLoginStrategyConcrete::Authenticate` should resolve)

**Do NOT search in Symbol Tree first.** Start with **Defined Strings** or **Search → For Strings**.

### Ghidra workflow per string

```
Defined Strings → click result → Right-click → References → Show References to
→ double-click CODE reference → read Decompiler pane
```

---

## Phase 2: EADP Auth (Protobuf — highest ROI)

### Goal

Reconstruct token exchange flow and extract static credentials (`clientId`, `clientSecret`, `productId`) without running the game.

### Entry-point strings (search in order)

1. `FootballEadpHub::requestAuthCode`
2. `GetAuthCodeJob`
3. `GrantTokenByAuthorizationCodeRequest`
4. `GetAuthForPcBlazeServerRequest`
5. `NucleusClientSecretProd_Decoder`
6. `PcLoginStrategyConcrete::Authenticate`
7. `EadpNexusConnectGrpcV1TokenServiceService`

### Expected auth call chain

```
PcLoginStrategyConcrete::Authenticate
  → FootballEadpHub::requestAuthCode
      → GetAuthCodeJob (Origin SDK LSX → EA App IPC)
  → build GrantTokenByAuthorizationCodeRequest
  → gRPC TokenService/GrantTokenByAuthorizationCode @ gateway.grpc.ea.com
  → build GetAuthForPcBlazeServerRequest
  → gRPC AuthService/GetAuthForPcBlazeServer
  → ConnectionManagerConcrete::Login (Blaze — Phase 3)
```

### Protobuf reconstruction (no .proto files shipped)

The binary embeds full type names and debug field printers. Key messages:

**`eadp.nexus.connect.grpc.v1.GrantTokenByAuthorizationCodeRequest`** fields (from `%s  fieldName:` debug strings):

- `clientId`, `clientSecret`, `code`, `redirectUri`
- `codeVerifier`, `codeChallenge`, `codeChallengeMethod` (PKCE)
- `machineProfile` (`PcMachineProfile`)

**gRPC paths:**

- `/eadp.nexus.connect.grpc.v1.TokenService/GrantTokenByAuthorizationCode`
- `/eadp.nexus.connect.grpc.v1.AuthService/GetAuthForPcBlazeServer`
- `/eadp.nexus.connect.grpc.v1.AuthService/GetAuthForOriginPc`

**Config anchors:**

- `eadp/foundation/Hub.h HubConfig::clientId`
- `eadp/foundation/Hub.h HubConfig::productId`
- `NucleusClientSecretProd_Decoder` — obfuscated client secret decoder

### Static credential extraction task

1. Xref `NucleusClientSecretProd_Decoder` in Ghidra
2. Find encoded blob in `.data`/`.rdata`
3. Decompile decoder → reimplement in Python/Go
4. Same for `HubConfig::clientId` initialization path

### Terminal pre-work (optional, fast)

```bash
strings -a "/run/media/preston/Data/College Football 27/CollegeFB27.exe" \
  | rg '^eadp\.nexus\.connect' | sort -u > nexus_types.txt

strings -a "/run/media/preston/Data/College Football 27/CollegeFB27.exe" \
  | rg '^\%s  [a-zA-Z]+:' | sort -u > proto_debug_fields.txt
```

---

## Phase 3: Blaze Session (TDF — harder, defer until Phase 2 done)

### Goal

Understand how authenticated client opens persistent Blaze connection and receives pushes.

### Entry-point strings

1. `ConnectionManagerConcrete::Login`
2. `ConnectionManagerConcrete::onAuthenticated`
3. `ConnectionManagerConcrete::StartLogin`
4. `Blaze::Authentication::TrustedLoginRequest`
5. `Blaze::Authentication::LoginResponse`
6. `UpdateNotificationHandler::onServerNotification`
7. `Blaze::FranTkMode::Notifications::ClientNotificationMessage`

### Expected Blaze flow

```
Redirector (spring25.client.blazeredirector.ea.com, serviceName=COLLEGEFB_2027_PC)
  → QoS ping site selection
  → TLS connect via DirtySDK/ProtoSSL
  → TrustedLoginRequest (TDF-encoded)
  → LoginResponse (blazeId, session)
  → persistent connection for push notifications
```

**TDF is NOT protobuf.** Use [jacobtread/tdf](https://docs.rs/tdf) crate or [grid-leak/blaze](https://github.com/grid-leak/blaze) as reference for other EA titles.

---

## Phase 4: Online Dynasty State Updates (core interest)

### Goal

Map how server state changes reach the client.

### Architecture pattern

- **Push:** `*EventInfo` notifications over Blaze connection
- **Pull:** `EnterRequest`/`EnterResponse` and `*RefreshFormResponse` when opening screens
- **Apply:** `UCM::patchFranTkInstance`, `UpdateNotificationHandler::ClassStateUpdate`

### Key symbols

| Symbol | Role |
|--------|------|
| `Online::DynastyMode::DynastyModeManager` | Client coordinator |
| `UCM::mOnlineClient` / `UCM::mOfflineClient` | Online vs offline fork |
| `LoadOnlineLeagueInstance` | Join league |
| `RegisterUser` / `UserRegisteredEventInfo` | Subscribe to league events |
| `Blaze::DynastyMode::Career::AdvanceWeekEventInfo` | Week advance push |
| `Blaze::DynastyMode::Career::AutoSaveEventInfo` | Save checkpoint trigger |
| `Blaze::FranTkMode::Notifications::LeagueInstanceLoadedEventInfo` | Instance hydrated |
| `UpdateNotificationHandler::onServerNotification` | Central push dispatcher |
| `DynastyServer::*DataConverter::CreateCustomRecordTree` | Server record → client TDF |

### Advance-week example flow

```
Commissioner RPC (AdvanceStage)
  → server sim
  → push StartAdvanceWeekEventInfo, AdvanceWeekEventInfo, AutoSaveEventInfo
  → client onServerNotification → ClassStateUpdate → patchFranTkInstance
  → UI adapters refresh
```

---

## Phase 5: Cloud Save Listing (separate from live updates)

| Symbol | Role |
|--------|------|
| `FranchiseLauncherDataProvider::RequestOnlineSaveFiles` | List online dynasties |
| `LeagueSavesManager::GetLeagueSavesAsync` | Async save API |
| `IPersistenceService::initializeDelegate()` | Cloud persistence init |
| `ParseLeagueSaves` | Deserialize save index |
| `https://gcs.ea.com` | Blob storage endpoint |

---

## What NOT to conflate

- **Gameplay desync** (`A Desync Has Occurred`, `GameManager`) ≠ dynasty menu persistence
- **EADP push notifications** (`eadp.pushnotification`, game invites) ≠ Blaze league EventInfo
- **Local dynasty export** (existing `cfb-dynasty` Go project) ≠ online server protocol

---

## Deliverables (in priority order)

1. **Auth map:** decompiled call graph from `FootballEadpHub::requestAuthCode` through gRPC stubs
2. **Extracted credentials:** `clientId`, `clientSecret`, `productId`, `redirectUri` from static analysis
3. **Reconstructed `.proto` snippets** for `GrantTokenByAuthorizationCodeRequest` and `GetAuthForPcBlazeServerRequest`
4. **Notification dispatch map:** `onServerNotification` → event type handlers → `patchFranTkInstance`
5. **Optional:** minimal Go/Python gRPC client that exchanges a manually-obtained auth code for tokens (auth code must come from EA App at runtime — cannot be fully static)

---

## Known blockers

| Blocker | Workaround |
|---------|------------|
| Javelin blocks Wireshark on game | Static RE only; or capture EA App traffic separately |
| Auth code requires EA App IPC | One-time manual auth code paste for POC; or EA App-only capture |
| Blaze uses TDF not protobuf | Reference tdf crate / other EA title emulators |
| `clientSecret` obfuscated | Decompile `NucleusClientSecretProd_Decoder` |
| 237 MB binary, slow analysis | Start from string xrefs, don't browse randomly |

---

## First 30 minutes checklist

- [ ] Open Ghidra project, confirm analysis complete
- [ ] Search `FootballEadpHub::requestAuthCode` → xref → decompile caller
- [ ] Search `NucleusClientSecretProd_Decoder` → decompile → extract secret
- [ ] Search `GrantTokenByAuthorizationCodeRequest` → find field population + gRPC invoke
- [ ] Search `UpdateNotificationHandler::onServerNotification` → map event dispatch switch
- [ ] Document findings in `cfb-dynasty/docs/` (or agent's working notes)

---

## Legal note

Academic/research use only. User owns legitimate EA account. Do not distribute extracted secrets or build production impersonation tooling. EA ToS may restrict third-party clients even with valid credentials.
