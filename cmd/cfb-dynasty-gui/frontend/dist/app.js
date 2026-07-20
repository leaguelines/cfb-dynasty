/* CFB Dynasty desktop GUI — Wails bindings + hash router */

const api = () =>
  window.go?.desktop?.App ||
  window.go?.main?.App ||
  Object.values(window.go || {}).map((p) => p?.App).find(Boolean);

function $(sel, root = document) { return root.querySelector(sel); }
function el(html) {
  const t = document.createElement("template");
  t.innerHTML = html.trim();
  return t.content.firstElementChild;
}
function esc(s) {
  return String(s ?? "")
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;");
}
function fmtSize(n) {
  if (n == null) return "";
  if (n < 1024) return n + " B";
  if (n < 1024 * 1024) return (n / 1024).toFixed(1) + " KB";
  return (n / (1024 * 1024)).toFixed(1) + " MB";
}

async function call(name, ...args) {
  const a = api();
  if (!a || typeof a[name] !== "function") {
    throw new Error("Desktop API not ready (" + name + ")");
  }
  return a[name](...args);
}

function setNav(cfg) {
  const nav = $("#nav");
  if (!cfg?.hasSave) {
    nav.innerHTML = "";
    return;
  }
  nav.innerHTML = `
    <span class="src" title="Source save">${esc(cfg.sourceName)}</span>
    <a href="#/dashboard">Dashboard</a>
    <button class="link" type="button" id="btn-close">Close save</button>
  `;
  $("#btn-close")?.addEventListener("click", async () => {
    await call("CloseSave");
    location.hash = "#/";
    refresh();
  });
}

function renderTable(header, rows, opts = {}) {
  if (!header?.length) return `<p class="muted">No columns.</p>`;
  const head = header.map((h) => `<th>${esc(h)}</th>`).join("");
  const body = (rows || []).map((row) => {
    const cells = row.map((cell, i) => {
      let v = esc(cell);
      if (opts.linkCol != null && i === opts.linkCol && cell) {
        v = `<a href="${esc(opts.linkPrefix || "")}${esc(cell)}">${esc(cell)}</a>`;
      }
      // Heuristic: player id columns
      if ((header[i] === "playerId" || header[i] === "id") && /^\d+$/.test(String(cell)) && opts.playerLinks) {
        v = `<a href="#/player/${esc(cell)}">${esc(cell)}</a>`;
      }
      if (header[i] === "firstName" || header[i] === "lastName" || header[i] === "player") {
        // leave plain unless we have an id elsewhere
      }
      return `<td>${v}</td>`;
    }).join("");
    return `<tr>${cells}</tr>`;
  }).join("");
  return `<div class="table-wrap"><table class="data"><thead><tr>${head}</tr></thead><tbody>${body || `<tr><td colspan="${header.length}" class="muted">No rows</td></tr>`}</tbody></table></div>`;
}

function logo(src, alt = "") {
  if (!src) return "";
  return `<img class="logo" src="${esc(src)}" alt="${esc(alt)}" onerror="this.style.display='none'" />`;
}

async function pageHome() {
  const cfg = await call("GetConfig");
  setNav(cfg);
  const schemaOk = cfg.schema?.valid;
  const saves = cfg.discovered || [];
  return `
    <section class="hero" style="text-align:left;max-width:900px;margin:0 auto;">
      <h1>CFB Dynasty Explorer</h1>
      <p class="lead">Point this app at your schema bundle and dynasty save to browse and export JSON or CSV. Nothing is uploaded.</p>

      <div class="setup-panel">
        <h2>1. Schema directory</h2>
        <p class="muted">Schema bundles (e.g. <code>C27_*.gz</code>) are required to decode saves. They are not shipped with this app.</p>
        <p>${schemaOk
          ? `<span class="badge" style="border-color:var(--accent-2);color:var(--accent-2)">Ready</span> ${esc(cfg.schema.dir)}`
          : `<span class="badge" style="border-color:var(--danger);color:var(--danger)">Needed</span> ${esc(cfg.schema?.message || "Not configured")}`
        }</p>
        ${cfg.schema?.bundles?.length ? `<p class="muted">Bundles: ${cfg.schema.bundles.map(esc).join(", ")}</p>` : ""}
        <div class="setup-row">
          <button class="btn" type="button" id="pick-schema">Choose schema folder…</button>
        </div>
      </div>

      <div class="setup-panel">
        <h2>2. Open a save</h2>
        <p class="muted">Default Windows location: <code>${esc(cfg.defaultSavesDir || "C:\\\\Users\\\\{you}\\\\Documents\\\\EA SPORTS CFB27\\\\saves")}</code></p>
        <div class="setup-row">
          <button class="btn" type="button" id="pick-save" ${schemaOk ? "" : "disabled"}>Browse for save…</button>
          ${cfg.hasSave ? `<a class="btn secondary" href="#/dashboard">Open dashboard</a>` : ""}
        </div>
        ${saves.length ? `
          <h3>Discovered saves</h3>
          <ul class="save-list">
            ${saves.map((s) => `
              <li>
                <div>
                  <strong>${esc(s.name)}</strong>
                  <div class="muted">${esc(s.path)} · ${fmtSize(s.size)}</div>
                </div>
                <button class="btn secondary open-save" data-path="${esc(s.path)}" ${schemaOk ? "" : "disabled"}>Open</button>
              </li>`).join("")}
          </ul>` : `<p class="muted">No saves auto-discovered. Use Browse.</p>`}
      </div>
      <div id="home-err"></div>
    </section>
  `;
}

function bindHome() {
  $("#pick-schema")?.addEventListener("click", async () => {
    try {
      await call("PickSchemaDir");
      await refresh();
    } catch (e) {
      $("#home-err").innerHTML = `<div class="err">${esc(e.message || e)}</div>`;
    }
  });
  $("#pick-save")?.addEventListener("click", async () => {
    try {
      const path = await call("PickSaveFile");
      if (!path) return;
      await call("OpenSave", path);
      location.hash = "#/dashboard";
      await refresh();
    } catch (e) {
      $("#home-err").innerHTML = `<div class="err">${esc(e.message || e)}</div>`;
    }
  });
  document.querySelectorAll(".open-save").forEach((btn) => {
    btn.addEventListener("click", async () => {
      try {
        await call("OpenSave", btn.dataset.path);
        location.hash = "#/dashboard";
        await refresh();
      } catch (e) {
        $("#home-err").innerHTML = `<div class="err">${esc(e.message || e)}</div>`;
      }
    });
  });
}

async function pageDashboard() {
  const cfg = await call("GetConfig");
  setNav(cfg);
  const dash = await call("GetDashboard");
  return `
    <section class="dashboard">
      <div class="dash-head">
        <div>
          <h1>Dashboard</h1>
          <p class="muted">Source: <strong>${esc(dash.sourceName)}</strong></p>
        </div>
        <div class="dl-group">
          <button class="btn" type="button" id="exp-json">Download all (JSON)</button>
          <button class="btn" type="button" id="exp-csv">Download all (CSV zip)</button>
        </div>
      </div>
      <div class="cards">
        ${(dash.cards || []).map((c) => `
          <div class="card ${c.count ? "" : "empty"}">
            <div class="card-top">
              <h2>${esc(c.title)}</h2>
              <span class="count">${c.count}</span>
            </div>
            ${c.count ? `
              <div class="card-actions">
                <a class="view" href="#/${esc(c.view.replace("collection:", "collection/"))}">View</a>
                <span class="dl">
                  <button class="link exp-col-json" data-name="${esc(c.name)}" type="button">JSON</button>
                  <button class="link exp-col-csv" data-name="${esc(c.name)}" type="button">CSV</button>
                </span>
              </div>` : `<div class="card-actions"><span class="muted">No data</span></div>`}
          </div>`).join("")}
      </div>
      <div id="dash-err"></div>
    </section>
  `;
}

function bindDashboard() {
  const err = (e) => { $("#dash-err").innerHTML = `<div class="err">${esc(e.message || e)}</div>`; };
  $("#exp-json")?.addEventListener("click", async () => {
    try { await call("ExportFullJSON"); } catch (e) { err(e); }
  });
  $("#exp-csv")?.addEventListener("click", async () => {
    try { await call("ExportFullCSVZip"); } catch (e) { err(e); }
  });
  document.querySelectorAll(".exp-col-json").forEach((b) => b.addEventListener("click", async () => {
    try { await call("ExportCollectionJSON", b.dataset.name); } catch (e) { err(e); }
  }));
  document.querySelectorAll(".exp-col-csv").forEach((b) => b.addEventListener("click", async () => {
    try {
      const name = b.dataset.name;
      if (name === "recruitingRankings") return;
      await call("ExportCollectionCSV", name);
    } catch (e) { err(e); }
  }));
}

async function pageCollection(name) {
  const cfg = await call("GetConfig");
  setNav(cfg);
  const t = await call("GetCollection", name);
  return `
    <section>
      <div class="dash-head">
        <div>
          <p><a href="#/dashboard">← Dashboard</a></p>
          <h1>${esc(t.title)}</h1>
        </div>
        <div class="dl-group">
          <button class="btn secondary" type="button" id="col-json">JSON</button>
          <button class="btn secondary" type="button" id="col-csv">CSV</button>
        </div>
      </div>
      ${renderTable(t.header, t.rows, { playerLinks: true })}
    </section>
  `;
}

function bindCollection(name) {
  $("#col-json")?.addEventListener("click", () => call("ExportCollectionJSON", name));
  $("#col-csv")?.addEventListener("click", () => call("ExportCollectionCSV", name));
}

async function pageSchools(query = "") {
  const cfg = await call("GetConfig");
  setNav(cfg);
  const schools = await call("GetSchools", query);
  return `
    <section>
      <p><a href="#/dashboard">← Dashboard</a></p>
      <h1>Schools</h1>
      <div class="setup-row">
        <input id="school-q" type="search" placeholder="Search schools…" value="${esc(query)}" style="flex:1;min-width:200px;padding:10px;border-radius:10px;border:1px solid var(--border);background:var(--panel);color:var(--text)" />
      </div>
      <div class="table-wrap"><table class="data"><thead>
        <tr><th></th><th>School</th><th>Conference</th><th>OVR</th><th>OFF</th><th>DEF</th><th>Prestige</th></tr>
      </thead><tbody>
        ${(schools || []).map((s) => `
          <tr>
            <td>${logo(s.SchoolLogo || s.schoolLogo, s.Name || s.name)}</td>
            <td><a href="#/schools/${s.ID ?? s.id}">${esc(s.Name || s.name)}</a></td>
            <td>${logo(s.ConferenceLogo || s.conferenceLogo)} ${esc(s.Conference || s.conference)}</td>
            <td>${esc(s.Ratings?.Overall ?? s.ratings?.Overall ?? "")}</td>
            <td>${esc(s.Ratings?.Offense ?? s.ratings?.Offense ?? "")}</td>
            <td>${esc(s.Ratings?.Defense ?? s.ratings?.Defense ?? "")}</td>
            <td>${esc(s.PrestigeStarsLabel || s.prestigeStarsLabel || s.PrestigeRank || "")}</td>
          </tr>`).join("")}
      </tbody></table></div>
    </section>
  `;
}

function bindSchools() {
  const input = $("#school-q");
  let t;
  input?.addEventListener("input", () => {
    clearTimeout(t);
    t = setTimeout(() => {
      location.hash = "#/schools?q=" + encodeURIComponent(input.value || "");
    }, 250);
  });
}

function field(obj, ...keys) {
  for (const k of keys) {
    if (obj && obj[k] != null && obj[k] !== "") return obj[k];
  }
  return "";
}

async function pageSchool(id) {
  const cfg = await call("GetConfig");
  setNav(cfg);
  const s = await call("GetSchool", id);
  const name = field(s, "Name", "name");
  return `
    <section class="stack">
      <p><a href="#/schools">← Schools</a></p>
      <div class="player-head">
        ${logo(field(s, "SchoolLogo", "schoolLogo"), name)}
        <div>
          <h1>${esc(name)}</h1>
          <p class="muted">${esc(field(s, "Conference", "conference"))}</p>
          <div class="setup-row">
            <a class="btn secondary" href="#/rosters/${id}">Roster</a>
            <a class="btn secondary" href="#/schools/${id}/class">Recruiting class</a>
          </div>
        </div>
      </div>
      <dl class="kv">
        <dt>Overall</dt><dd>${esc(s.Ratings?.Overall ?? s.ratings?.Overall ?? "")}</dd>
        <dt>Offense</dt><dd>${esc(s.Ratings?.Offense ?? s.ratings?.Offense ?? "")}</dd>
        <dt>Defense</dt><dd>${esc(s.Ratings?.Defense ?? s.ratings?.Defense ?? "")}</dd>
        <dt>Prestige</dt><dd>${esc(field(s, "PrestigeStarsLabel", "prestigeStarsLabel"))}</dd>
      </dl>
    </section>
  `;
}

async function pageSchoolClass(id) {
  const cfg = await call("GetConfig");
  setNav(cfg);
  const view = await call("GetSchoolClass", id);
  const rows = (view.Rows || view.rows || []).map((r) => [
    r.DisplayName || r.displayName || "",
    r.Position || r.position || "",
    r.StarsText || r.starsText || "",
    r.NIL || r.nil || "",
    r.NationalRank || r.nationalRank || "",
    r.Stage || r.stage || "",
    r.ID ?? r.id ?? "",
  ]);
  return `
    <section>
      <p><a href="#/schools/${id}">← ${esc(view.Name || view.name || "School")}</a></p>
      <div class="dash-head">
        <h1>Recruiting class</h1>
        <button class="btn secondary" type="button" id="class-csv">CSV</button>
      </div>
      ${renderTable(["Name", "Pos", "Stars", "NIL", "Nat#", "Stage", "playerId"], rows, { playerLinks: true })}
    </section>
  `;
}

function bindSchoolClass(id) {
  $("#class-csv")?.addEventListener("click", () => call("ExportSchoolClassCSV", id));
}

async function pageRosters() {
  const cfg = await call("GetConfig");
  setNav(cfg);
  const teams = await call("GetRosterTeams");
  return `
    <section>
      <p><a href="#/dashboard">← Dashboard</a></p>
      <h1>Rosters</h1>
      <div class="table-wrap"><table class="data"><thead>
        <tr><th></th><th>Team</th><th>Conference</th><th>Players</th><th>OVR</th></tr>
      </thead><tbody>
        ${(teams || []).map((t) => `
          <tr>
            <td>${logo(t.SchoolLogo || t.schoolLogo)}</td>
            <td><a href="#/rosters/${t.ID ?? t.id}">${esc(t.Name || t.name)}</a></td>
            <td>${esc(t.Conference || t.conference || "")}</td>
            <td>${esc(t.PlayerCount ?? t.playerCount ?? "")}</td>
            <td>${esc(t.Ratings?.Overall ?? t.ratings?.Overall ?? "")}</td>
          </tr>`).join("")}
      </tbody></table></div>
    </section>
  `;
}

async function pageRoster(id) {
  const cfg = await call("GetConfig");
  setNav(cfg);
  const r = await call("GetRoster", id);
  const players = r.Players || r.players || [];
  const positions = ["ALL", ...(r.Positions || r.positions || [])];
  const rowsFor = (pos) => players
    .filter((p) => pos === "ALL" || p.Position === pos || p.position === pos)
    .map((p) => [
      p.Jersey ?? p.jersey ?? "",
      p.FirstName || p.firstName || "",
      p.LastName || p.lastName || "",
      p.Overall ?? p.overall ?? "",
      p.StarRating || p.starRating || "",
      p.SchoolYear || p.schoolYear || "",
      p.Archetype || p.archetype || "",
      p.Position || p.position || "",
      p.ID ?? p.id,
    ]);
  return `
    <section>
      <p><a href="#/rosters">← Rosters</a></p>
      <h1>${esc(r.TeamName || r.teamName || "Roster")}</h1>
      <div class="tabs" id="pos-tabs">
        ${positions.map((p, i) => `<button type="button" class="tab ${i === 0 ? "active" : ""}" data-pos="${esc(p)}">${esc(p)}</button>`).join("")}
      </div>
      <div id="roster-table">${renderTable(
        ["jersey", "firstName", "lastName", "overall", "starRating", "class", "archetype", "position", "playerId"],
        rowsFor("ALL"),
        { playerLinks: true }
      )}</div>
    </section>
  `;
}

function bindRoster(id) {
  // re-fetch not needed; filter client-side from already rendered data is harder —
  // simplest: re-route with query. For now re-call GetRoster on tab click.
  document.querySelectorAll("#pos-tabs .tab").forEach((tab) => {
    tab.addEventListener("click", async () => {
      document.querySelectorAll("#pos-tabs .tab").forEach((t) => t.classList.remove("active"));
      tab.classList.add("active");
      const r = await call("GetRoster", id);
      const players = r.Players || r.players || [];
      const pos = tab.dataset.pos;
      const rows = players
        .filter((p) => pos === "ALL" || p.Position === pos || p.position === pos)
        .map((p) => [
          p.Jersey ?? p.jersey ?? "",
          p.FirstName || p.firstName || "",
          p.LastName || p.lastName || "",
          p.Overall ?? p.overall ?? "",
          p.StarRating || p.starRating || "",
          p.SchoolYear || p.schoolYear || "",
          p.Archetype || p.archetype || "",
          p.Position || p.position || "",
          p.ID ?? p.id,
        ]);
      $("#roster-table").innerHTML = renderTable(
        ["jersey", "firstName", "lastName", "overall", "starRating", "class", "archetype", "position", "playerId"],
        rows,
        { playerLinks: true }
      );
    });
  });
}

async function pageLeaving() {
  const cfg = await call("GetConfig");
  setNav(cfg);
  const teams = await call("GetLeavingTeams");
  return `
    <section>
      <p><a href="#/dashboard">← Dashboard</a></p>
      <h1>Leaving players</h1>
      <div class="table-wrap"><table class="data"><thead>
        <tr><th>Team</th><th>Players</th></tr>
      </thead><tbody>
        ${(teams || []).map((t) => `
          <tr>
            <td><a href="#/leaving/${t.ID ?? t.id ?? t.TeamID ?? t.teamId}">${esc(t.Name || t.name || t.TeamName || t.teamName)}</a></td>
            <td>${esc(t.PlayerCount ?? t.playerCount ?? t.Count ?? t.count ?? "")}</td>
          </tr>`).join("")}
      </tbody></table></div>
    </section>
  `;
}

async function pageLeavingTeam(id) {
  const cfg = await call("GetConfig");
  setNav(cfg);
  const v = await call("GetLeavingTeam", id);
  return `
    <section>
      <p><a href="#/leaving">← Leaving</a></p>
      <h1>${esc(v.TeamName || v.teamName)}</h1>
      ${renderTable(v.Header || v.header, v.Rows || v.rows, { playerLinks: true })}
    </section>
  `;
}

async function pageRecruits() {
  const cfg = await call("GetConfig");
  setNav(cfg);
  const t = await call("GetRecruits");
  return `
    <section>
      <div class="dash-head">
        <div>
          <p><a href="#/dashboard">← Dashboard</a></p>
          <h1>Recruits</h1>
        </div>
        <div class="dl-group">
          <button class="btn secondary" type="button" id="col-json">JSON</button>
          <button class="btn secondary" type="button" id="col-csv">CSV</button>
          <a class="btn secondary" href="#/recruiting">Recruiting board</a>
          <a class="btn secondary" href="#/recruiting-rankings">Rankings</a>
        </div>
      </div>
      ${renderTable(t.header, t.rows, { playerLinks: true })}
    </section>
  `;
}

async function pageRecruiting() {
  const cfg = await call("GetConfig");
  setNav(cfg);
  const t = await call("GetRecruiting");
  return `
    <section>
      <p><a href="#/recruits">← Recruits</a></p>
      <h1>Recruiting</h1>
      ${renderTable(t.header, t.rows, { playerLinks: true })}
    </section>
  `;
}

async function pageRecruitingRankings(conf = "") {
  const cfg = await call("GetConfig");
  setNav(cfg);
  const v = await call("GetRecruitingRankings", conf);
  const rows = (v.Rows || v.rows || []).map((r) => [
    r.Rank ?? r.rank ?? "",
    r.TeamName || r.teamName || "",
    r.Conference || r.conference || "",
    r.Commits ?? r.commits ?? "",
    r.AvgOverall || r.avgOverall || "",
    r.TotalPoints || r.totalPoints || "",
    r.TeamID ?? r.teamId ?? "",
  ]);
  const filters = v.Conferences || v.conferences || [];
  return `
    <section>
      <p><a href="#/dashboard">← Dashboard</a></p>
      <h1>Recruiting rankings</h1>
      <div class="tabs">
        <a class="tab ${!conf ? "active" : ""}" href="#/recruiting-rankings">All</a>
        ${filters.map((f) => {
          const label = f.Name || f.name || "";
          return `<a class="tab ${conf === label ? "active" : ""}" href="#/recruiting-rankings?conf=${encodeURIComponent(label)}">${esc(label)}</a>`;
        }).join("")}
      </div>
      ${renderTable(["Rank", "Team", "Conference", "Commits", "Avg OVR", "Points", "teamId"], rows)}
    </section>
  `;
}

async function pagePlayer(id) {
  const cfg = await call("GetConfig");
  setNav(cfg);
  const v = await call("GetPlayer", id);
  const p = v.player || {};
  const name = `${p.FirstName || p.firstName || ""} ${p.LastName || p.lastName || ""}`.trim();
  const ratings = p.Ratings || p.ratings || {};
  const ratingRows = Object.entries(ratings).sort((a, b) => b[1] - a[1]).map(([k, val]) => `<tr><td>${esc(k)}</td><td>${esc(val)}</td></tr>`).join("");
  const interest = (v.schoolInterest || []).map((s) => `
    <tr>
      <td>${logo(s.TeamLogo || s.teamLogo)} <a href="#/schools/${s.TeamID ?? s.teamId}/class">${esc(s.TeamName || s.teamName)}</a></td>
      <td>${esc(s.Influence ?? s.influence ?? "")}</td>
    </tr>`).join("");
  return `
    <section class="stack">
      <p><a href="#/dashboard">← Dashboard</a></p>
      <div class="player-head">
        <div>
          <h1>${esc(name || "Player")}</h1>
          <p class="muted">
            ${esc(p.Position || p.position || "")}
            · OVR ${esc(p.Overall ?? p.overall ?? "")}
            · ${esc(p.StarRating || p.starRating || "")}
            ${v.isRecruit ? " · Recruit" : ""}
            ${v.teamName ? ` · <a href="#/rosters/${v.teamId}">${esc(v.teamName)}</a>` : ""}
          </p>
        </div>
      </div>
      <dl class="kv">
        <dt>Class</dt><dd>${esc(p.SchoolYear || p.schoolYear || "")}</dd>
        <dt>Archetype</dt><dd>${esc(p.ArchetypeLabel || p.archetypeLabel || p.Archetype || p.archetype || "")}</dd>
        <dt>Dev trait</dt><dd>${esc(p.DevTrait || p.devTrait || "")}</dd>
        <dt>Hometown</dt><dd>${esc([p.HomeTown || p.homeTown, p.HomeState || p.homeState].filter(Boolean).join(", "))}</dd>
        <dt>Height/Weight</dt><dd>${esc(p.Height ?? p.height ?? "")} / ${esc(p.Weight ?? p.weight ?? "")}</dd>
      </dl>
      ${interest ? `<h2>School interest</h2><div class="table-wrap"><table class="data"><thead><tr><th>School</th><th>Influence</th></tr></thead><tbody>${interest}</tbody></table></div>` : ""}
      <h2>Ratings</h2>
      <div class="table-wrap"><table class="data"><thead><tr><th>Attribute</th><th>Value</th></tr></thead><tbody>${ratingRows || `<tr><td colspan="2" class="muted">None</td></tr>`}</tbody></table></div>
    </section>
  `;
}

async function pageGame(id) {
  const cfg = await call("GetConfig");
  setNav(cfg);
  const g = await call("GetGame", id);
  return `
    <section class="stack">
      <p><a href="#/collection/games">← Games</a></p>
      <h1>${esc(g.AwayTeam || g.awayTeam || "Away")} @ ${esc(g.HomeTeam || g.homeTeam || "Home")}</h1>
      <p class="muted">Week ${esc(g.Week ?? g.week ?? "")} · ${esc(g.WeekType || g.weekType || "")} · ${esc(g.Status || g.status || "")}</p>
      <p style="font-size:1.6rem;font-weight:700">${esc(g.AwayScore ?? g.awayScore ?? "-")} – ${esc(g.HomeScore ?? g.homeScore ?? "-")}</p>
      <dl class="kv">
        <dt>Stadium</dt><dd>${esc(g.StadiumName || g.stadiumName || "")}</dd>
        <dt>Weather</dt><dd>${esc(g.Weather || g.weather || "")}</dd>
        <dt>Network</dt><dd>${esc(g.BroadcastNetwork || g.broadcastNetwork || "")}</dd>
      </dl>
    </section>
  `;
}

async function pageRecords(period = "career") {
  const cfg = await call("GetConfig");
  setNav(cfg);
  const tabs = await call("GetRecordPeriods");
  const t = await call("GetRecords", period, "", "");
  return `
    <section>
      <p><a href="#/dashboard">← Dashboard</a></p>
      <h1>Record book</h1>
      <div class="tabs">
        ${(tabs || []).map((tab) => `
          <a class="tab ${(tab.Period || tab.period) === period ? "active" : ""}" href="#/records/${esc(tab.Period || tab.period)}">
            ${esc(tab.Title || tab.title)} (${esc(tab.Count ?? tab.count ?? 0)})
          </a>`).join("")}
      </div>
      <div class="dl-group" style="margin-bottom:12px">
        <button class="btn secondary" type="button" id="col-csv">CSV</button>
      </div>
      ${renderTable(t.header, t.rows)}
    </section>
  `;
}

function bindRecords() {
  $("#col-csv")?.addEventListener("click", () => call("ExportCollectionCSV", "recordBook"));
}

function parseRoute() {
  const raw = (location.hash || "#/").replace(/^#/, "") || "/";
  const [pathPart, queryPart] = raw.split("?");
  const path = pathPart.replace(/\/+$/, "") || "/";
  const params = new URLSearchParams(queryPart || "");
  const parts = path.split("/").filter(Boolean);
  return { path, parts, params };
}

let bindAfter = null;

async function render() {
  const app = $("#app");
  bindAfter = null;
  try {
    const { parts, params } = parseRoute();
    let html;
    if (parts.length === 0) {
      html = await pageHome();
      bindAfter = bindHome;
    } else if (parts[0] === "dashboard") {
      html = await pageDashboard();
      bindAfter = bindDashboard;
    } else if (parts[0] === "collection" && parts[1]) {
      html = await pageCollection(parts[1]);
      bindAfter = () => bindCollection(parts[1]);
    } else if (parts[0] === "schools" && parts[1] && parts[2] === "class") {
      const id = Number(parts[1]);
      html = await pageSchoolClass(id);
      bindAfter = () => bindSchoolClass(id);
    } else if (parts[0] === "schools" && parts[1]) {
      html = await pageSchool(Number(parts[1]));
    } else if (parts[0] === "schools") {
      html = await pageSchools(params.get("q") || "");
      bindAfter = bindSchools;
    } else if (parts[0] === "rosters" && parts[1]) {
      const id = Number(parts[1]);
      html = await pageRoster(id);
      bindAfter = () => bindRoster(id);
    } else if (parts[0] === "rosters") {
      html = await pageRosters();
    } else if (parts[0] === "leaving" && parts[1]) {
      html = await pageLeavingTeam(Number(parts[1]));
    } else if (parts[0] === "leaving") {
      html = await pageLeaving();
    } else if (parts[0] === "recruits") {
      html = await pageRecruits();
      bindAfter = () => bindCollection("recruits");
    } else if (parts[0] === "recruiting-rankings") {
      html = await pageRecruitingRankings(params.get("conf") || "");
    } else if (parts[0] === "recruiting") {
      html = await pageRecruiting();
    } else if (parts[0] === "player" && parts[1]) {
      html = await pagePlayer(Number(parts[1]));
    } else if (parts[0] === "game" && parts[1]) {
      html = await pageGame(Number(parts[1]));
    } else if (parts[0] === "records") {
      html = await pageRecords(parts[1] || "career");
      bindAfter = bindRecords;
    } else {
      html = `<div class="err">Unknown page</div><p><a href="#/">Home</a></p>`;
    }
    app.innerHTML = html;
    bindAfter?.();
  } catch (e) {
    const msg = e?.message || String(e);
    if (/no save loaded/i.test(msg)) {
      location.hash = "#/";
      return refresh();
    }
    app.innerHTML = `<div class="err">${esc(msg)}</div><p><a href="#/">Back home</a></p>`;
    try {
      const cfg = await call("GetConfig");
      setNav(cfg);
    } catch (_) { /* ignore */ }
  }
}

async function refresh() {
  await render();
}

window.addEventListener("hashchange", () => { refresh(); });
window.addEventListener("DOMContentLoaded", () => {
  // Wails injects bindings shortly after load
  const tryStart = async () => {
    if (!api()) {
      setTimeout(tryStart, 50);
      return;
    }
    await refresh();
  };
  tryStart();
});
