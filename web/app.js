/* ZatpatSite builder — vanilla JS, no dependencies. */
(function () {
  "use strict";

  var $ = function (id) { return document.getElementById(id); };
  var api = function (path, opts) {
    return fetch("/api/v1" + path, opts).then(function (r) {
      if (!r.ok) {
        return r.json().catch(function () { return {}; }).then(function (b) {
          throw new Error(b.error || ("Request failed (" + r.status + ")"));
        });
      }
      return r;
    });
  };

  // Monday-first display order; indices match Go/JS Date.getDay() (0=Sunday).
  var DAY_ORDER = [1, 2, 3, 4, 5, 6, 0];
  var DAY_SHORT = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"];

  var CATEGORY_LABELS = {
    salon: "Salon", restaurant: "Restaurant", kirana: "Kirana store",
    coaching: "Coaching classes", clinic: "Clinic", boutique: "Boutique",
    gym: "Gym", bakery: "Bakery"
  };

  // ---------- state ----------
  var state = {
    id: null,          // null = unsaved new site
    theme: "ivory",
    lang: "en",
    rating: 0,
    reviews: [],
    createdAt: null,
    lastHTML: ""
  };
  var meta = { themes: [], categories: [] };

  function defaultHours() {
    var h = [];
    for (var i = 0; i < 7; i++) h.push({ closed: false, open: "09:30", close: "20:30" });
    h[0] = { closed: false, open: "10:00", close: "14:00" };
    return h;
  }

  // ---------- form <-> payload ----------
  function readHours() {
    var hrs = [];
    for (var d = 0; d < 7; d++) {
      hrs.push({
        closed: $("h-closed-" + d).checked,
        open: $("h-open-" + d).value,
        close: $("h-close-" + d).value
      });
    }
    return hrs;
  }

  function readServices() {
    var out = [];
    document.querySelectorAll(".svc-row").forEach(function (row) {
      var name = row.querySelector(".svc-name-in").value.trim();
      var price = parseInt(row.querySelector(".svc-price-in").value, 10);
      if (name) out.push({ name: name, price: isNaN(price) || price < 0 ? 0 : price });
    });
    return out;
  }

  function payload() {
    return {
      name: $("f-name").value.trim(),
      category: $("f-category").value,
      city: $("f-city").value.trim(),
      address: $("f-address").value.trim(),
      phone: $("f-phone").value.trim(),
      whatsapp: $("f-whatsapp").value.trim(),
      mapsUrl: $("f-maps").value.trim(),
      hours: readHours(),
      services: readServices(),
      about: $("f-about").value.trim(),
      theme: state.theme,
      lang: state.lang,
      rating: state.rating,
      reviews: state.reviews
    };
  }

  function fillForm(s) {
    $("f-name").value = s.name || "";
    $("f-category").value = s.category || "salon";
    $("f-city").value = s.city || "";
    $("f-address").value = s.address || "";
    $("f-phone").value = s.phone || "";
    $("f-whatsapp").value = s.whatsapp || "";
    $("f-maps").value = s.mapsUrl || "";
    $("f-about").value = s.about || "";
    setHours(s.hours || defaultHours());
    setServices(s.services || []);
    state.theme = s.theme || "ivory";
    state.lang = s.lang || "en";
    state.rating = s.rating || 0;
    state.reviews = s.reviews || [];
    paintThemeChips();
    paintLang();
    paintFormTitle();
  }

  // ---------- hours grid ----------
  function buildHoursGrid() {
    var grid = $("hours-grid");
    grid.innerHTML = "";
    DAY_ORDER.forEach(function (d) {
      var row = document.createElement("div");
      row.className = "hrow";
      row.id = "hrow-" + d;
      row.innerHTML =
        '<span class="dname">' + DAY_SHORT[d] + "</span>" +
        '<label class="closed-toggle"><input type="checkbox" id="h-closed-' + d + '"> Closed</label>' +
        '<input type="time" id="h-open-' + d + '" value="09:30">' +
        '<span class="dash">–</span>' +
        '<input type="time" id="h-close-' + d + '" value="20:30">';
      grid.appendChild(row);
      $("h-closed-" + d).addEventListener("change", function () {
        row.classList.toggle("is-closed", this.checked);
        schedulePreview();
      });
      $("h-open-" + d).addEventListener("input", schedulePreview);
      $("h-close-" + d).addEventListener("input", schedulePreview);
    });
  }

  function setHours(hrs) {
    for (var d = 0; d < 7; d++) {
      var h = hrs[d] || { closed: false, open: "09:30", close: "20:30" };
      $("h-closed-" + d).checked = !!h.closed;
      $("h-open-" + d).value = h.open || "09:30";
      $("h-close-" + d).value = h.close || "20:30";
      $("hrow-" + d).classList.toggle("is-closed", !!h.closed);
    }
  }

  // ---------- services rows ----------
  function addServiceRow(name, price, focus) {
    var row = document.createElement("div");
    row.className = "svc-row";
    row.innerHTML =
      '<input type="text" class="svc-name-in" placeholder="e.g. Haircut (Men)">' +
      '<span class="price-wrap"><input type="number" class="svc-price-in" min="0" step="1" placeholder="250"></span>' +
      '<button type="button" class="rm" aria-label="Remove service">✕</button>';
    row.querySelector(".svc-name-in").value = name || "";
    if (price || price === 0) row.querySelector(".svc-price-in").value = price;
    row.querySelector(".rm").addEventListener("click", function () {
      row.remove();
      paintSvcEmpty();
      schedulePreview();
    });
    row.querySelectorAll("input").forEach(function (i) {
      i.addEventListener("input", schedulePreview);
    });
    $("svc-rows").appendChild(row);
    paintSvcEmpty();
    if (focus) row.querySelector(".svc-name-in").focus();
  }

  function setServices(list) {
    $("svc-rows").innerHTML = "";
    (list || []).forEach(function (s) { addServiceRow(s.name, s.price, false); });
    paintSvcEmpty();
  }

  function paintSvcEmpty() {
    $("svc-empty").hidden = document.querySelectorAll(".svc-row").length > 0;
  }

  // ---------- theme chips / language ----------
  function paintThemeChips() {
    var box = $("theme-chips");
    box.innerHTML = "";
    meta.themes.forEach(function (t) {
      var b = document.createElement("button");
      b.type = "button";
      b.className = "theme-chip" + (t.id === state.theme ? " active" : "");
      b.setAttribute("role", "radio");
      b.setAttribute("aria-checked", t.id === state.theme ? "true" : "false");
      var sw = t.swatches.split(",").map(function (c) {
        return '<span class="sw" style="background:' + c + '"></span>';
      }).join("");
      b.innerHTML = '<span class="swatches">' + sw + "</span>" + t.label;
      b.addEventListener("click", function () {
        state.theme = t.id;
        paintThemeChips();
        schedulePreview();
      });
      box.appendChild(b);
    });
  }

  function paintLang() {
    document.querySelectorAll("#lang-toggle .seg-btn").forEach(function (b) {
      var on = b.getAttribute("data-lang") === state.lang;
      b.classList.toggle("active", on);
      b.setAttribute("aria-checked", on ? "true" : "false");
    });
  }

  function paintFormTitle() {
    $("form-title").textContent = state.id ? "Edit website" : "Build your website";
    $("btn-save").textContent = state.id ? "Update site" : "Save site";
  }

  function slugify(name) {
    var s = (name || "").toLowerCase().replace(/[^a-z0-9]+/g, "-").replace(/^-+|-+$/g, "");
    return s || "your-shop";
  }

  // ---------- live preview ----------
  var previewTimer = null;
  function schedulePreview() {
    clearTimeout(previewTimer);
    previewTimer = setTimeout(renderPreview, 320);
  }

  function renderPreview() {
    var p = payload();
    $("url-pill").textContent = "zatpat.site/" + slugify(p.name);
    if (!p.name) p.name = "Your Business Name";
    api("/render", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(p)
    }).then(function (r) { return r.text(); }).then(function (html) {
      state.lastHTML = html;
      $("preview").srcdoc = html;
    }).catch(function (e) {
      toast(e.message, true);
    });
  }

  // ---------- actions ----------
  function saveSite() {
    var p = payload();
    if (!p.name) { toast("Add a business name first", true); $("f-name").focus(); return; }
    var req = state.id
      ? api("/sites/" + state.id, { method: "PUT", headers: { "Content-Type": "application/json" }, body: JSON.stringify(p) })
      : api("/sites", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify(p) });
    req.then(function (r) { return r.json(); }).then(function (s) {
      var isNew = !state.id;
      state.id = s.id;
      paintFormTitle();
      refreshCount();
      toast(isNew ? "Saved — “" + s.name + "” is in My sites" : "Updated “" + s.name + "”");
    }).catch(function (e) { toast(e.message, true); });
  }

  function downloadSite() {
    if (state.id) {
      window.location.href = "/api/v1/sites/" + state.id + "/download";
      return;
    }
    if (!state.lastHTML) { toast("Nothing to download yet", true); return; }
    var blob = new Blob([state.lastHTML], { type: "text/html" });
    var a = document.createElement("a");
    a.href = URL.createObjectURL(blob);
    a.download = "index.html";
    document.body.appendChild(a);
    a.click();
    a.remove();
    setTimeout(function () { URL.revokeObjectURL(a.href); }, 3000);
    toast("index.html downloaded");
  }

  function fetchFromGoogle() {
    var name = $("f-name").value.trim();
    if (!name) { toast("Type the business name first", true); $("f-name").focus(); return; }
    var btn = $("btn-fetch");
    btn.disabled = true;
    api("/gbp/fetch", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        name: name,
        city: $("f-city").value.trim(),
        category: $("f-category").value,
        lang: state.lang
      })
    }).then(function (r) { return r.json(); }).then(function (res) {
      var p = res.profile;
      $("f-category").value = p.category;
      $("f-address").value = p.address;
      $("f-phone").value = p.phone;
      $("f-whatsapp").value = p.whatsapp;
      $("f-maps").value = p.mapsUrl;
      if (!$("f-about").value.trim()) $("f-about").value = res.about || "";
      setHours(p.hours);
      setServices(p.services);
      state.rating = p.rating;
      state.reviews = p.reviews;
      renderPreview();
      toast("Profile fetched — " + p.rating.toFixed(1) + " ★, " + p.reviews.length + " reviews");
    }).catch(function (e) {
      toast(e.message, true);
    }).finally(function () { btn.disabled = false; });
  }

  function newSite() {
    state.id = null;
    state.rating = 0;
    state.reviews = [];
    fillForm({ theme: state.theme, lang: state.lang, hours: defaultHours(), services: [] });
    addServiceRow("", "", false);
    addServiceRow("", "", false);
    renderPreview();
    $("f-name").focus();
  }

  // ---------- my sites drawer ----------
  function openDrawer() {
    $("drawer").hidden = false;
    $("drawer-scrim").hidden = false;
    loadSites();
  }
  function closeDrawer() {
    $("drawer").hidden = true;
    $("drawer-scrim").hidden = true;
  }

  function loadSites() {
    api("/sites").then(function (r) { return r.json(); }).then(function (res) {
      var list = $("drawer-list");
      list.innerHTML = "";
      var sites = res.sites || [];
      $("sites-count").textContent = sites.length;
      $("drawer-empty").hidden = sites.length > 0;
      sites.forEach(function (s) {
        var card = document.createElement("div");
        card.className = "site-card";
        var meta = (CATEGORY_LABELS[s.category] || s.category) +
          (s.city ? " · " + s.city : "") +
          " · updated " + new Date(s.updatedAt).toLocaleDateString("en-IN", { day: "numeric", month: "short" });
        card.innerHTML =
          '<div class="sc-top"><span class="sc-name"></span><span class="sc-theme">' + s.theme + "</span></div>" +
          '<p class="sc-meta"></p>' +
          '<div class="sc-actions">' +
          '<button type="button" class="btn btn-primary sc-edit">Edit</button>' +
          '<a class="btn btn-ghost-dark" href="/api/v1/sites/' + s.id + '/preview" target="_blank" rel="noopener">Open</a>' +
          '<a class="btn btn-ghost-dark" href="/api/v1/sites/' + s.id + '/download">Download</a>' +
          '<button type="button" class="btn btn-danger sc-del">Delete</button>' +
          "</div>";
        card.querySelector(".sc-name").textContent = s.name;
        card.querySelector(".sc-meta").textContent = meta;
        card.querySelector(".sc-edit").addEventListener("click", function () {
          api("/sites/" + s.id).then(function (r) { return r.json(); }).then(function (full) {
            state.id = full.id;
            fillForm(full);
            renderPreview();
            closeDrawer();
            toast("Editing “" + full.name + "”");
          }).catch(function (e) { toast(e.message, true); });
        });
        card.querySelector(".sc-del").addEventListener("click", function () {
          if (!confirm("Delete “" + s.name + "”? This cannot be undone.")) return;
          api("/sites/" + s.id, { method: "DELETE" }).then(function () {
            if (state.id === s.id) { state.id = null; paintFormTitle(); }
            loadSites();
            toast("Deleted “" + s.name + "”");
          }).catch(function (e) { toast(e.message, true); });
        });
        list.appendChild(card);
      });
    }).catch(function (e) { toast(e.message, true); });
  }

  function refreshCount() {
    api("/sites").then(function (r) { return r.json(); }).then(function (res) {
      $("sites-count").textContent = (res.sites || []).length;
    }).catch(function () {});
  }

  // ---------- toast ----------
  var toastTimer = null;
  function toast(msg, isErr) {
    var t = $("toast");
    t.textContent = msg;
    t.classList.toggle("err", !!isErr);
    t.hidden = false;
    clearTimeout(toastTimer);
    toastTimer = setTimeout(function () { t.hidden = true; }, 2800);
  }

  // ---------- boot ----------
  function boot() {
    buildHoursGrid();

    api("/meta").then(function (r) { return r.json(); }).then(function (m) {
      meta.themes = m.themes;
      meta.categories = m.categories;
      var sel = $("f-category");
      m.categories.forEach(function (c) {
        var o = document.createElement("option");
        o.value = c;
        o.textContent = CATEGORY_LABELS[c] || c;
        sel.appendChild(o);
      });
      paintThemeChips();
      newSite();
      refreshCount();
    }).catch(function (e) { toast(e.message, true); });

    ["f-name", "f-city", "f-address", "f-phone", "f-whatsapp", "f-maps", "f-about"].forEach(function (id) {
      $(id).addEventListener("input", schedulePreview);
    });
    $("f-category").addEventListener("change", schedulePreview);

    $("btn-add-svc").addEventListener("click", function () { addServiceRow("", "", true); });
    $("btn-copy-mon").addEventListener("click", function () {
      var mon = { closed: $("h-closed-1").checked, open: $("h-open-1").value, close: $("h-close-1").value };
      for (var d = 0; d < 7; d++) {
        if (d === 1) continue;
        $("h-closed-" + d).checked = mon.closed;
        $("h-open-" + d).value = mon.open;
        $("h-close-" + d).value = mon.close;
        $("hrow-" + d).classList.toggle("is-closed", mon.closed);
      }
      schedulePreview();
      toast("Monday hours copied to all days");
    });

    $("btn-fetch").addEventListener("click", fetchFromGoogle);
    $("btn-save").addEventListener("click", saveSite);
    $("btn-download").addEventListener("click", downloadSite);
    $("btn-new").addEventListener("click", newSite);
    $("btn-my-sites").addEventListener("click", openDrawer);
    $("btn-close-drawer").addEventListener("click", closeDrawer);
    $("drawer-scrim").addEventListener("click", closeDrawer);
    document.addEventListener("keydown", function (e) {
      if (e.key === "Escape") closeDrawer();
    });

    $("lang-toggle").addEventListener("click", function (e) {
      var b = e.target.closest(".seg-btn");
      if (!b) return;
      state.lang = b.getAttribute("data-lang");
      paintLang();
      schedulePreview();
    });
    $("device-toggle").addEventListener("click", function (e) {
      var b = e.target.closest(".seg-btn");
      if (!b) return;
      document.querySelectorAll("#device-toggle .seg-btn").forEach(function (x) {
        x.classList.toggle("active", x === b);
      });
      $("iframe-wrap").classList.toggle("mobile", b.getAttribute("data-device") === "mobile");
    });

    $("site-form").addEventListener("submit", function (e) { e.preventDefault(); saveSite(); });
  }

  boot();
})();
