# CFB27 Online Dynasty RE — Findings (Static Analysis)

Status: static analysis of `CollegeFB27.exe` (PE32+ x86-64, ~237 MB) via Ghidra MCP + raw
binary parsing. Companion to `cfb27-online-dynasty-re-handoff.md`.

Build id (from embedded PDB path):
`E:\dev\TnT\LocalCFB\Bin\CollegeFB\Win64\retail\CollegeFB.Main_Win64_retail.pdb`

---

## 0. Tooling reality / methodology note

The handoff's core GUI workflow ("Defined Strings → Show References To → Decompiler")
does **not** translate through the MCP bridge on this binary:

- MCP `data_list_strings` reliably **times out** (42 MB `.rdata`, 237 MB image).
- The high-value type/name/field strings have **zero** cross-references of any kind:
  no 8-byte absolute pointer, no 32-bit rip-relative `LEA`/`disp32` from any of the 3
  code sections (`.text`, `ctr`, `.text2`), and no 32-bit RVA.
- Reason: EADP's protobuf/logging is **table-driven & reflective**. Field-name and
  type-name strings live in packed metadata tables and are reached by `base + offset`
  at runtime, not by direct references. Debug output uses one big `\n`-joined format
  string per message; the `strings` tool splits it, so each `"%s  field: %s,"` line is
  a *fragment*, interleaved with groups of 6 code pointers (per-message method tables).

Consequence: string→code xref chasing is a dead end here. Productive paths are
(a) raw string/credential mining (done below) and (b) following the reflection
pointer-tables into the serialize/parse/debug functions for true decompilation (TODO).

Section offset→VA map (from `objdump -h`), image base `0x140000000`:

| section | raw file off | virtual addr | size | perm |
|---------|-------------|--------------|------|------|
| .text   | 0x00000600  | 0x140001000  | 0x09748600 | R-X |
| ctr     | 0x09748c00  | 0x14974a000  | 0x008e0600 | R-X |
| .rdata  | 0x0a029200  | 0x14a02b000  | 0x02836c00 | R-- |
| .data   | 0x0c85fe00  | 0x14c862000  | 0x003a6a00 | RW- |
| .text(2)| 0x0e51ce00  | 0x15064a000  | 0x004c0200 | R-X |

`.rdata` conversion: `VA = file_offset + 0x140001e00`.

---

## 1. Extracted credentials (Deliverable #2)

Client credential table in `.rdata` at ~`0x14a6fce38` (immediately preceding the
`NucleusClientSecret_Decoder` / `NucleusClientSecretProd_Decoder` name strings).
Values appear in **plaintext** — the "decoder" likely just selects/returns them.

| Field | Value |
|-------|-------|
| clientId | `CFB_27_PC_CLIENT` |
| clientSecret (rot 2025-05-01) | `ZmGSfNEtUOlZC0jVZ0538zaU6MI2HcqpMNLHBOBnhTieawMGK6JQoRd0l2dhacHBs_20250501203047` |
| clientSecret (rot 2026-01-13, newer) | `8eMFBK8PcwezDFWg2i4RejqXhw3zR827i1DHSCUk34HS9QafIVNymPlFfHweWoVsd_20260113170640` |
| redirectUri | `nucleus:rest` |
| authorize endpoint path | `/connect/auth` |
| identity scheme | `eadp.identity.v2` |
| productId / slug | `collegefb.27` |
| Blaze serviceName | `COLLEGEFB_2027_PC` |

Two secrets = standard Nucleus rotation (old + current). Format is
`<62-char base62>_<YYYYMMDDhhmmss rotation timestamp>`.

Adjacent entries (`MAD23PCC`, `niquu2iek…/xxxx…`, `oi5tei2ce…`, `ojaing3of…`,
`chohwie0v…`) are **EADP-SDK demo/placeholder** keys, not live CFB27 credentials.

> SENSITIVE. Academic use only; do not distribute (see handoff legal note).

---

## 2. EADP Nexus gRPC surface (Deliverables #1 / #3)

Host: `gateway.grpc.ea.com` (int: `gateway.grpc.int.ea.com`, lt: `…lt.ea.com`).
Package: `eadp.nexus.connect.grpc.v1`. Auth wrappers:
`eadp.common.v1.NexusServiceAuth`, `eadp.common.v1.NexusMethodAuth`.

### Services / methods relevant to PC login

`TokenService`:
`GrantTokenByAuthorizationCode`, `GrantTokenByRefreshToken`, `GrantTokenByExchange`,
`GrantTokenByClientCredentials`, `GrantTokenByPassword`, `GrantTokenByOtc`,
`GrantTokenByChangeGameState`, `GrantTokenByGetGameState`, `DeleteToken`, … (full list
of ~30 incl. Retrofit* migration grants captured).

`AuthService` (blaze-server ticket issuance):
`GetAuthForPcBlazeServer`, `GetAuthForPcClient`, `GetAuthForOriginPc`,
`GetAuthForCommonBlazeServer`, `GetAuthForToken`, `GetAuthForAuthV2`, plus per-platform
variants (Ps4/Ps5/Xone/Xbsx/Nx/Steam/Epic/Mobile…).

`TokenInfoService/GetTokenInfo`, `JwkService/GetJwk`.

Key method paths for PC online-dynasty login:
```
/eadp.nexus.connect.grpc.v1.TokenService/GrantTokenByAuthorizationCode
/eadp.nexus.connect.grpc.v1.AuthService/GetAuthForPcBlazeServer
/eadp.nexus.connect.grpc.v1.AuthService/GetAuthForOriginPc
```

### Reconstructed request messages (field names + wire types from debug printers)

Field numbers are NOT yet recovered (require the serialize fn / descriptor table).
Types inferred from printer format specifiers: `"%s"`→string, `%d`→int32/enum,
`%lld`→int64, `%f`→float, `[`→repeated, `%s`(msg)→nested message.

```proto
// package eadp.nexus.connect.grpc.v1;  (field numbers TBD)
message GrantTokenByAuthorizationCodeRequest {
  string clientId            = ?;
  string clientSecret        = ?;
  string code                = ?;
  string redirectUri         = ?;
  string codeVerifier        = ?;   // PKCE
  string codeChallenge       = ?;   // PKCE
  string codeChallengeMethod = ?;   // PKCE
  PcMachineProfile machineProfile = ?;
}
```
Adjacent authorize-request fields also seen: `responseType`, `releaseType`, `display`,
`locale`, `registrationSource`, `prompt`, `fid`.

Token responses reference: `accessToken`, `refreshToken`, `tokenType`, `code`,
`idToken`, `longLiveToken`, `persona`, `pidId`.

---

## 3. Auth-code (LSX/EA App) flow constants

Strings clustered together (`0x14a3be0d0`+):
`GetAuthCodeJob`, `nucleus:rest`, `/connect/auth`, `eadp.identity.v2`,
`client_id`, `response_type`, `access_t…`, plus error
`Invalid request: Client ID parameter cannot be empty.`
Log tags: `FootballEadpHub::requestAuthCode success` / `…error: %d/%s`,
`PcLoginStrategyConcrete::Authenticate - null eaUser`, `…loginRequestFailure`.

Flow (matches handoff hypothesis):
```
PcLoginStrategyConcrete::Authenticate
  -> FootballEadpHub::requestAuthCode  (GetAuthCodeJob via Origin LSX / EA App IPC;
        redirect_uri=nucleus:rest, response_type=code, client_id=CFB_27_PC_CLIENT)
  -> TokenService/GrantTokenByAuthorizationCode  (+ PKCE, clientSecret, machineProfile)
  -> AuthService/GetAuthForPcBlazeServer  (blaze ticket)
  -> ConnectionManagerConcrete::Login  (Blaze — Phase 3)
```

---

## 4. Blaze / endpoints (Phase 3 anchors)

- Redirector: `spring25.client.blazeredirector.ea.com`
  (+ `.cert`/`.dev`/`.test`, `gosca25.blazeredirector.ea.com` for trusted cert).
- serviceName `COLLEGEFB_2027_PC`.
- Accounts: `https://accounts.ea.com`, `https://accounts.grpc.ea.com`.
- Gateway proxy: `https://gateway.ea.com/proxy`.
- Cloud saves (Phase 5): `https://gcs.ea.com` (+ int/load).
- Errors telemetry: `https://collector.errors.ea.com`.

---

## 5. Open TODO

1. Recover protobuf **field numbers** by decompiling the serialize/parse fns reached
   via the reflection pointer-table groups adjacent to the field-format strings.
2. Confirm which secret is active (rotation) and how `NucleusClientSecretProd_Decoder`
   selects/returns it.
3. Map `UpdateNotificationHandler::onServerNotification` dispatch (Phase 4) — same
   table-driven approach needed.
4. Optional: minimal gRPC client for `GrantTokenByAuthorizationCode` (auth `code` must
   still come from a live EA App LSX exchange).
