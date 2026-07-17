package generator

import "strings"

// Theme metadata used by the builder UI.
type ThemeInfo struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	Swatches string `json:"swatches"` // comma-separated hex chips for the picker
}

// ThemeList describes the four designs for the builder UI.
var ThemeList = []ThemeInfo{
	{ID: "ivory", Label: "Ivory", Swatches: "#FBF7EF,#2A2520,#8A6D3B"},
	{ID: "bazaar", Label: "Bazaar", Swatches: "#FFB703,#111111,#E5397B"},
	{ID: "mint", Label: "Mint", Swatches: "#F5FAF7,#0CA678,#12291F"},
	{ID: "nightshade", Label: "Nightshade", Swatches: "#0F0D13,#D4AF63,#F4EDDA"},
}

// resetCSS is the only CSS shared by all four themes: a box-model reset,
// smooth scrolling and the floating WhatsApp button's geometry.
const resetCSS = `
*,*::before,*::after{margin:0;padding:0;box-sizing:border-box}
html{scroll-behavior:smooth;-webkit-text-size-adjust:100%}
body{-webkit-font-smoothing:antialiased}
a{color:inherit;text-decoration:none}
ul{list-style:none}
svg{display:inline-block;vertical-align:middle}
table{border-collapse:collapse;width:100%}
section{scroll-margin-top:84px}
.wa-float{position:fixed;right:18px;bottom:18px;z-index:50;display:flex;align-items:center;justify-content:center;width:58px;height:58px;border-radius:50%;transition:transform .15s ease}
.wa-float:hover{transform:scale(1.07)}
.wa-float svg{width:28px;height:28px}
.nav-inner,.hero-inner{width:100%}
@media (prefers-reduced-motion:reduce){html{scroll-behavior:auto}*,*::before,*::after{transition:none!important;animation:none!important}}
`

// __PATTERN__ placeholders are replaced with the category glyph data-URI.

const ivoryCSS = `
body.t-ivory{background:#FBF7EF;color:#2A2520;font-family:"Iowan Old Style","Palatino Linotype",Palatino,"Book Antiqua",Georgia,serif;font-size:17px;line-height:1.65}
.t-ivory .nav{position:sticky;top:0;z-index:40;background:rgba(251,247,239,.94);backdrop-filter:blur(6px);border-bottom:1px solid #D8CFBE}
.t-ivory .nav-inner{max-width:860px;margin:0 auto;padding:15px 24px;display:flex;justify-content:space-between;align-items:baseline;gap:16px;flex-wrap:wrap}
.t-ivory .brand{font-size:15px;letter-spacing:.18em;text-transform:uppercase}
.t-ivory .navlinks{display:flex;gap:22px;font-size:12px;letter-spacing:.14em;text-transform:uppercase;flex-wrap:wrap}
.t-ivory .navlinks a{border-bottom:1px solid transparent;padding-bottom:2px;transition:border-color .2s}
.t-ivory .navlinks a:hover{border-color:#2A2520}
.t-ivory .hero{text-align:center;padding:88px 24px 64px;border-bottom:4px double #C9BFA9;background-image:url("__PATTERN__")}
.t-ivory .emblem-wrap{width:88px;height:88px;margin:0 auto 28px;border:1px solid #B9AD94;outline:1px solid #B9AD94;outline-offset:5px;border-radius:50%;display:flex;align-items:center;justify-content:center;color:#8A6D3B;background:#FBF7EF}
.t-ivory .kicker{font-size:12px;letter-spacing:.28em;text-transform:uppercase;color:#8A6D3B;margin-bottom:18px}
.t-ivory h1{font-size:clamp(34px,6vw,54px);font-weight:500;letter-spacing:.01em;margin-bottom:14px}
.t-ivory .tagline{font-style:italic;font-size:20px;color:#5C5344;max-width:34ch;margin:0 auto 28px}
.t-ivory .orn{display:none}
.t-ivory .badge-row{display:flex;justify-content:center;align-items:center;gap:14px;margin-bottom:32px;flex-wrap:wrap}
.t-ivory .badge{font-size:11px;letter-spacing:.16em;text-transform:uppercase;padding:7px 14px;border:1px solid #B9AD94}
.t-ivory .badge.open{color:#3F6B4A;border-color:#3F6B4A}
.t-ivory .badge.closed{color:#96442E;border-color:#96442E}
.t-ivory .rating-chip{display:inline-flex;align-items:center;gap:6px;font-size:13px;color:#8A6D3B}
.t-ivory .cta-row{display:flex;justify-content:center;gap:14px;flex-wrap:wrap}
.t-ivory .btn{display:inline-flex;align-items:center;gap:9px;padding:13px 26px;border:1px solid #2A2520;font-size:13px;letter-spacing:.14em;text-transform:uppercase;transition:background .25s,color .25s}
.t-ivory .btn:hover{background:#2A2520;color:#FBF7EF}
.t-ivory .btn-wa{background:#2A2520;color:#FBF7EF}
.t-ivory .btn-wa:hover{background:transparent;color:#2A2520}
.t-ivory .sec{padding:72px 24px}
.t-ivory .sec .wrap{max-width:660px;margin:0 auto}
.t-ivory .sec-title{font-size:13px;font-weight:500;letter-spacing:.3em;text-transform:uppercase;color:#8A6D3B;display:flex;align-items:center;justify-content:center;gap:18px;margin-bottom:42px;text-align:center}
.t-ivory .sec-title::before,.t-ivory .sec-title::after{content:"";height:1px;background:#D8CFBE;flex:1;max-width:120px}
.t-ivory .svc-item{display:flex;align-items:baseline;gap:10px;padding:13px 0;font-size:18px}
.t-ivory .leader{flex:1;border-bottom:1px dotted #B9AD94;transform:translateY(-4px)}
.t-ivory .svc-price{font-variant-numeric:tabular-nums;color:#8A6D3B}
.t-ivory .sec-alt{background:#F4EDE0;border-top:1px solid #E4DAC6;border-bottom:1px solid #E4DAC6}
.t-ivory .about-text{font-size:19px;line-height:1.85;text-align:center;font-style:italic;color:#4A4237}
.t-ivory .hours-note{text-align:center;margin-bottom:24px}
.t-ivory .hours td{padding:11px 4px;border-bottom:1px solid #E4DAC6}
.t-ivory .hours td.time{text-align:right;font-variant-numeric:tabular-nums}
.t-ivory .hours tr.today td{color:#8A6D3B;font-weight:700}
.t-ivory .hours .day-closed{color:#96442E;font-style:italic}
.t-ivory .review{margin-bottom:38px;text-align:center}
.t-ivory .review:last-child{margin-bottom:0}
.t-ivory .stars{color:#B98A2E;display:flex;gap:3px;justify-content:center;margin-bottom:12px}
.t-ivory blockquote{font-size:19px;font-style:italic;line-height:1.7;quotes:none}
.t-ivory figcaption{margin-top:12px;font-size:12px;letter-spacing:.2em;text-transform:uppercase;color:#8A6D3B}
.t-ivory .visit-block{text-align:center}
.t-ivory address{font-style:normal;font-size:18px;line-height:1.8;margin-bottom:10px}
.t-ivory .visit-phone{font-size:16px;color:#5C5344;margin-bottom:28px}
.t-ivory .footer{border-top:1px solid #D8CFBE;padding:40px 24px;text-align:center;font-size:12px;letter-spacing:.18em;text-transform:uppercase;color:#8A7B5F}
.t-ivory .footer .made{margin-top:8px;opacity:.65}
.t-ivory .wa-float{background:#2A2520;color:#FBF7EF;box-shadow:0 8px 24px rgba(42,37,32,.35)}
`

const bazaarCSS = `
body.t-bazaar{background:#FFFDF7;color:#111;font-family:-apple-system,"Segoe UI",system-ui,Roboto,"Helvetica Neue",Arial,sans-serif;font-size:17px;line-height:1.55}
.t-bazaar .nav{position:sticky;top:0;z-index:40;background:#111;color:#FFFDF7}
.t-bazaar .nav-inner{max-width:1080px;margin:0 auto;padding:12px 20px;display:flex;justify-content:space-between;align-items:center;gap:14px;flex-wrap:wrap}
.t-bazaar .brand{font-weight:900;text-transform:uppercase;letter-spacing:.04em;font-size:17px}
.t-bazaar .navlinks{display:flex;gap:4px;flex-wrap:wrap}
.t-bazaar .navlinks a{font-weight:800;text-transform:uppercase;font-size:12px;padding:6px 10px;border:2px solid transparent;transition:border-color .15s,color .15s}
.t-bazaar .navlinks a:hover{border-color:#FFB703;color:#FFB703}
.t-bazaar .hero{background:#FFB703 url("__PATTERN__");border-bottom:3px solid #111;padding:72px 20px 64px}
.t-bazaar .hero-inner{max-width:1080px;margin:0 auto}
.t-bazaar .emblem-wrap{width:84px;height:84px;background:#FFFDF7;border:3px solid #111;box-shadow:6px 6px 0 #111;display:flex;align-items:center;justify-content:center;color:#111;margin-bottom:28px;transform:rotate(-3deg)}
.t-bazaar .kicker{display:inline-block;background:#111;color:#FFB703;font-weight:800;text-transform:uppercase;font-size:12px;letter-spacing:.08em;padding:6px 12px;margin-bottom:18px;transform:rotate(-1deg)}
.t-bazaar h1{font-weight:900;text-transform:uppercase;font-size:clamp(40px,8vw,76px);line-height:.98;letter-spacing:-.02em;margin-bottom:16px}
.t-bazaar .tagline{font-size:20px;font-weight:600;max-width:44ch;margin-bottom:28px}
.t-bazaar .orn{display:none}
.t-bazaar .badge-row{display:flex;align-items:center;gap:14px;margin-bottom:30px;flex-wrap:wrap}
.t-bazaar .badge{background:#FFFDF7;border:3px solid #111;border-radius:999px;padding:7px 16px;font-weight:800;text-transform:uppercase;font-size:13px}
.t-bazaar .badge.open{background:#ADE25D}
.t-bazaar .badge.closed{background:#FF8FB3}
.t-bazaar .rating-chip{display:inline-flex;align-items:center;gap:6px;background:#FFFDF7;border:3px solid #111;border-radius:999px;padding:7px 16px;font-weight:800;box-shadow:4px 4px 0 #111}
.t-bazaar .cta-row{display:flex;gap:16px;flex-wrap:wrap}
.t-bazaar .btn{display:inline-flex;align-items:center;gap:9px;background:#111;color:#FFFDF7;font-weight:800;text-transform:uppercase;font-size:14px;padding:14px 24px;border:3px solid #111;box-shadow:6px 6px 0 #111;transition:transform .12s,box-shadow .12s}
.t-bazaar .btn:hover{transform:translate(3px,3px);box-shadow:2px 2px 0 #111}
.t-bazaar .btn-wa{background:#26B95F;color:#fff}
.t-bazaar .btn-call{background:#FFFDF7;color:#111}
.t-bazaar .sec{padding:72px 20px;border-bottom:3px solid #111}
.t-bazaar .wrap{max-width:1080px;margin:0 auto}
.t-bazaar .sec-title{font-weight:900;text-transform:uppercase;font-size:clamp(26px,4vw,40px);display:inline-block;background:#111;color:#FFFDF7;padding:8px 18px;transform:rotate(-1deg);margin-bottom:40px}
.t-bazaar .svc{display:grid;grid-template-columns:repeat(auto-fill,minmax(250px,1fr));gap:22px}
.t-bazaar .svc-item{border:3px solid #111;background:#FFFDF7;box-shadow:6px 6px 0 #111;padding:18px;display:flex;flex-direction:column;gap:12px;transition:transform .15s}
.t-bazaar .svc-item:nth-child(odd){transform:rotate(.6deg)}
.t-bazaar .svc-item:nth-child(even){transform:rotate(-.6deg)}
.t-bazaar .svc-item:hover{transform:rotate(0)}
.t-bazaar .svc-name{font-weight:800;font-size:17px}
.t-bazaar .leader{display:none}
.t-bazaar .svc-price{align-self:flex-start;background:#FFB703;border:2px solid #111;padding:4px 10px;font-weight:900;font-size:16px}
.t-bazaar .sec-alt{background:#E5397B;color:#FFFDF7}
.t-bazaar .sec-alt .sec-title{background:#FFFDF7;color:#111}
.t-bazaar .about-text{font-size:22px;font-weight:700;line-height:1.6;max-width:52ch}
.t-bazaar .hours-note{margin-bottom:22px}
.t-bazaar .hours-card{border:3px solid #111;box-shadow:6px 6px 0 #111;background:#fff;max-width:560px;padding:6px 18px}
.t-bazaar .hours td{padding:12px 6px;border-bottom:2px solid #111;font-weight:700}
.t-bazaar .hours tr:last-child td{border-bottom:none}
.t-bazaar .hours td.time{text-align:right}
.t-bazaar .hours tr.today td{background:#FFB703}
.t-bazaar .hours .day-closed{color:#E5397B}
.t-bazaar #reviews{background:#53B7E8}
.t-bazaar .reviews{display:grid;grid-template-columns:repeat(auto-fill,minmax(270px,1fr));gap:22px}
.t-bazaar .review{background:#FFFDF7;border:3px solid #111;box-shadow:6px 6px 0 #111;padding:20px}
.t-bazaar .stars{color:#111;display:flex;gap:2px;margin-bottom:10px}
.t-bazaar blockquote{font-weight:600;line-height:1.55}
.t-bazaar figcaption{margin-top:12px;font-weight:900;text-transform:uppercase;font-size:12px}
.t-bazaar address{font-style:normal;font-size:20px;font-weight:700;line-height:1.6;margin-bottom:8px}
.t-bazaar .visit-phone{font-weight:700;margin-bottom:26px}
.t-bazaar .footer{background:#111;color:#FFFDF7;padding:34px 20px;text-align:center;font-weight:700;text-transform:uppercase;font-size:13px}
.t-bazaar .footer .made{margin-top:6px;color:#FFB703}
.t-bazaar .wa-float{background:#26B95F;color:#fff;border:3px solid #111;box-shadow:5px 5px 0 #111}
`

const mintCSS = `
body.t-mint{background:#F5FAF7;color:#173B2C;font-family:-apple-system,"Segoe UI",system-ui,Roboto,"Helvetica Neue",Arial,sans-serif;font-size:16.5px;line-height:1.6}
.t-mint .nav{position:sticky;top:12px;z-index:40;padding:0 16px}
.t-mint .nav-inner{max-width:880px;margin:0 auto;background:rgba(255,255,255,.92);backdrop-filter:blur(8px);border-radius:999px;box-shadow:0 8px 24px rgba(23,59,44,.1);padding:10px 22px;display:flex;justify-content:space-between;align-items:center;gap:12px;flex-wrap:wrap}
.t-mint .brand{font-weight:800;font-size:16px;color:#0B7A57}
.t-mint .navlinks{display:flex;gap:6px;flex-wrap:wrap}
.t-mint .navlinks a{font-size:13.5px;font-weight:600;padding:7px 12px;border-radius:999px;color:#3E6355;transition:background .2s,color .2s}
.t-mint .navlinks a:hover{background:#DDF3E9;color:#0B7A57}
.t-mint .hero{padding:88px 20px 72px;text-align:center;background:radial-gradient(620px 320px at 12% 0%,#D8F3E6 0%,transparent 70%),radial-gradient(520px 340px at 92% 14%,#E3F6EC 0%,transparent 70%),url("__PATTERN__")}
.t-mint .emblem-wrap{width:84px;height:84px;margin:0 auto 24px;background:linear-gradient(145deg,#12BB84,#0B8A61);color:#fff;border-radius:26px;display:flex;align-items:center;justify-content:center;box-shadow:0 14px 30px rgba(12,166,120,.35)}
.t-mint .kicker{display:inline-flex;background:#DDF3E9;color:#0B7A57;font-weight:700;font-size:13px;border-radius:999px;padding:7px 16px;margin-bottom:20px}
.t-mint h1{font-weight:800;font-size:clamp(34px,6vw,56px);letter-spacing:-.02em;color:#12291F;margin-bottom:14px}
.t-mint .tagline{font-size:19px;color:#3E6355;max-width:40ch;margin:0 auto 28px}
.t-mint .orn{display:none}
.t-mint .badge-row{display:flex;justify-content:center;align-items:center;gap:12px;margin-bottom:30px;flex-wrap:wrap}
.t-mint .badge{border-radius:999px;font-weight:700;font-size:13.5px;padding:8px 16px}
.t-mint .badge.open{background:#D3F5E4;color:#0B7A57}
.t-mint .badge.closed{background:#FCE8E4;color:#B5432F}
.t-mint .rating-chip{display:inline-flex;align-items:center;gap:6px;background:#fff;border-radius:999px;padding:8px 16px;box-shadow:0 6px 18px rgba(23,59,44,.1);font-weight:700;color:#173B2C}
.t-mint .rating-chip svg{color:#F5A623}
.t-mint .cta-row{display:flex;justify-content:center;gap:14px;flex-wrap:wrap}
.t-mint .btn{border-radius:999px;padding:14px 28px;font-weight:700;display:inline-flex;gap:9px;align-items:center;transition:transform .15s,box-shadow .15s}
.t-mint .btn:hover{transform:translateY(-2px)}
.t-mint .btn-wa{background:#0CA678;color:#fff;box-shadow:0 12px 26px rgba(12,166,120,.35)}
.t-mint .btn-call{background:#fff;color:#0B7A57;box-shadow:inset 0 0 0 2px #BFE8D6}
.t-mint .sec{padding:72px 20px}
.t-mint .wrap{max-width:960px;margin:0 auto}
.t-mint .sec-title{text-align:center;font-weight:800;font-size:clamp(24px,3.5vw,34px);letter-spacing:-.01em;color:#12291F;margin-bottom:46px}
.t-mint .sec-title::after{content:"";display:block;width:44px;height:4px;border-radius:2px;background:#0CA678;margin:14px auto 0}
.t-mint .svc{display:grid;grid-template-columns:repeat(auto-fill,minmax(230px,1fr));gap:18px}
.t-mint .svc-item{background:#fff;border-radius:20px;padding:22px;box-shadow:0 10px 26px rgba(23,59,44,.07);display:flex;flex-direction:column;gap:8px;transition:transform .18s,box-shadow .18s}
.t-mint .svc-item:hover{transform:translateY(-4px);box-shadow:0 16px 34px rgba(23,59,44,.13)}
.t-mint .svc-name{font-weight:700;font-size:16.5px;color:#12291F}
.t-mint .leader{display:none}
.t-mint .svc-price{color:#0B7A57;font-weight:800;font-size:20px}
.t-mint #about .wrap{max-width:760px}
.t-mint .about-text{background:#fff;border-radius:24px;padding:36px;box-shadow:0 10px 26px rgba(23,59,44,.07);font-size:18px;line-height:1.75;color:#28503F}
.t-mint .hours-note{text-align:center;margin-bottom:24px}
.t-mint .hours-card{background:#fff;border-radius:24px;box-shadow:0 10px 26px rgba(23,59,44,.07);padding:14px 26px;max-width:560px;margin:0 auto}
.t-mint .hours td{padding:12px 4px;border-bottom:1px solid #EAF4EE;font-weight:600}
.t-mint .hours tr:last-child td{border-bottom:none}
.t-mint .hours td.time{text-align:right;color:#3E6355}
.t-mint .hours tr.today td{color:#0B7A57;font-weight:800}
.t-mint .hours .day-closed{color:#B5432F}
.t-mint .reviews{display:grid;grid-template-columns:repeat(auto-fill,minmax(260px,1fr));gap:18px}
.t-mint .review{background:#fff;border-radius:20px;padding:24px;box-shadow:0 10px 26px rgba(23,59,44,.07)}
.t-mint .stars{color:#F5A623;display:flex;gap:2px;margin-bottom:12px}
.t-mint blockquote{line-height:1.65;color:#28503F}
.t-mint figcaption{margin-top:14px;font-weight:700;font-size:13.5px;color:#0B7A57}
.t-mint .visit-block{background:#fff;border-radius:24px;box-shadow:0 10px 26px rgba(23,59,44,.07);padding:36px;text-align:center;max-width:640px;margin:0 auto}
.t-mint address{font-style:normal;font-size:17px;line-height:1.7;margin-bottom:8px}
.t-mint .visit-phone{color:#3E6355;font-weight:600;margin-bottom:24px}
.t-mint .footer{padding:44px 20px;text-align:center;color:#3E6355;font-size:14px}
.t-mint .footer .made{margin-top:6px;font-weight:700;color:#0B7A57}
.t-mint .wa-float{background:#0CA678;color:#fff;box-shadow:0 14px 30px rgba(12,166,120,.4)}
`

const nightshadeCSS = `
body.t-nightshade{background:#0F0D13;color:#E9E2D2;font-family:"Iowan Old Style","Palatino Linotype",Palatino,Georgia,serif;font-size:17px;line-height:1.7}
.t-nightshade .nav{position:sticky;top:0;z-index:40;background:rgba(15,13,19,.92);backdrop-filter:blur(8px);border-bottom:1px solid rgba(212,175,99,.25)}
.t-nightshade .nav-inner{max-width:900px;margin:0 auto;padding:16px 24px;display:flex;justify-content:space-between;align-items:baseline;gap:16px;flex-wrap:wrap}
.t-nightshade .brand{font-size:15px;letter-spacing:.22em;text-transform:uppercase;color:#D4AF63}
.t-nightshade .navlinks{display:flex;gap:20px;flex-wrap:wrap}
.t-nightshade .navlinks a{font-family:-apple-system,system-ui,sans-serif;font-size:11.5px;letter-spacing:.18em;text-transform:uppercase;color:#B8AE97;padding-bottom:3px;border-bottom:1px solid transparent;transition:color .2s,border-color .2s}
.t-nightshade .navlinks a:hover{color:#D4AF63;border-color:#D4AF63}
.t-nightshade .hero{text-align:center;padding:96px 24px 72px;background:radial-gradient(700px 380px at 50% -10%,rgba(212,175,99,.16),transparent 65%),url("__PATTERN__")}
.t-nightshade .emblem-wrap{width:92px;height:92px;margin:0 auto 28px;border:1px solid rgba(212,175,99,.65);outline:1px solid rgba(212,175,99,.25);outline-offset:6px;border-radius:50%;display:flex;align-items:center;justify-content:center;color:#D4AF63}
.t-nightshade .kicker{font-family:-apple-system,system-ui,sans-serif;font-size:11.5px;letter-spacing:.34em;text-transform:uppercase;color:#D4AF63;margin-bottom:20px}
.t-nightshade h1{font-size:clamp(36px,6vw,58px);font-weight:400;letter-spacing:.02em;color:#F4EDDA;margin-bottom:16px}
.t-nightshade .tagline{font-style:italic;font-size:20px;color:#B8AE97;margin:0 auto;max-width:38ch}
.t-nightshade .orn{display:flex;align-items:center;gap:14px;justify-content:center;margin:28px auto 30px;max-width:280px}
.t-nightshade .orn::before,.t-nightshade .orn::after{content:"";flex:1;height:1px}
.t-nightshade .orn::before{background:linear-gradient(90deg,transparent,rgba(212,175,99,.75))}
.t-nightshade .orn::after{background:linear-gradient(90deg,rgba(212,175,99,.75),transparent)}
.t-nightshade .orn span{width:7px;height:7px;transform:rotate(45deg);background:#D4AF63}
.t-nightshade .badge-row{display:flex;justify-content:center;align-items:center;gap:16px;margin-bottom:32px;flex-wrap:wrap}
.t-nightshade .badge{font-family:-apple-system,system-ui,sans-serif;font-size:11px;letter-spacing:.2em;text-transform:uppercase;padding:8px 16px;border:1px solid rgba(212,175,99,.4)}
.t-nightshade .badge.open{color:#9FD8A8;border-color:rgba(159,216,168,.55)}
.t-nightshade .badge.closed{color:#D89A8C;border-color:rgba(216,154,140,.55)}
.t-nightshade .rating-chip{display:inline-flex;align-items:center;gap:6px;color:#D4AF63;font-family:-apple-system,system-ui,sans-serif;font-size:13px;letter-spacing:.06em}
.t-nightshade .cta-row{display:flex;justify-content:center;gap:16px;flex-wrap:wrap}
.t-nightshade .btn{font-family:-apple-system,system-ui,sans-serif;font-size:12px;letter-spacing:.18em;text-transform:uppercase;padding:14px 28px;border:1px solid #D4AF63;color:#D4AF63;display:inline-flex;gap:9px;align-items:center;transition:background .25s,color .25s}
.t-nightshade .btn:hover{background:#D4AF63;color:#14111A}
.t-nightshade .btn-wa{background:#D4AF63;color:#14111A}
.t-nightshade .btn-wa:hover{background:transparent;color:#D4AF63}
.t-nightshade .sec{padding:80px 24px}
.t-nightshade .wrap{max-width:640px;margin:0 auto}
.t-nightshade .sec-title{text-align:center;font-size:30px;font-weight:400;color:#F4EDDA;margin-bottom:46px}
.t-nightshade .sec-title::after{content:"";display:block;width:56px;height:1px;background:#D4AF63;margin:18px auto 0}
.t-nightshade .svc-item{display:flex;align-items:baseline;gap:14px;padding:18px 0;border-bottom:1px solid rgba(212,175,99,.18);font-size:18px}
.t-nightshade .leader{flex:1}
.t-nightshade .svc-price{font-family:-apple-system,system-ui,sans-serif;color:#D4AF63;font-weight:600;letter-spacing:.04em;font-variant-numeric:tabular-nums}
.t-nightshade .sec-alt{background:#14111A;border-top:1px solid rgba(212,175,99,.15);border-bottom:1px solid rgba(212,175,99,.15)}
.t-nightshade .about-text{font-size:19px;line-height:1.9;text-align:center;color:#CBC2AD;font-style:italic}
.t-nightshade .hours-note{text-align:center;margin-bottom:24px}
.t-nightshade .hours-card{border:1px solid rgba(212,175,99,.25);background:#14111A;padding:10px 26px}
.t-nightshade .hours td{padding:12px 4px;border-bottom:1px solid rgba(212,175,99,.12)}
.t-nightshade .hours tr:last-child td{border-bottom:none}
.t-nightshade .hours td.time{text-align:right;font-family:-apple-system,system-ui,sans-serif;font-size:14.5px;letter-spacing:.04em}
.t-nightshade .hours tr.today td{color:#D4AF63}
.t-nightshade .hours .day-closed{color:#D89A8C;font-style:italic}
.t-nightshade .review{border:1px solid rgba(212,175,99,.2);background:#14111A;padding:28px;margin-bottom:22px}
.t-nightshade .review:last-child{margin-bottom:0}
.t-nightshade .stars{color:#D4AF63;display:flex;gap:3px;margin-bottom:14px}
.t-nightshade blockquote{font-style:italic;font-size:18px;line-height:1.75}
.t-nightshade figcaption{margin-top:14px;font-family:-apple-system,system-ui,sans-serif;font-size:11px;letter-spacing:.22em;text-transform:uppercase;color:#D4AF63}
.t-nightshade .visit-block{text-align:center}
.t-nightshade address{font-style:normal;font-size:18px;line-height:1.85;margin-bottom:10px}
.t-nightshade .visit-phone{font-family:-apple-system,system-ui,sans-serif;font-size:14px;letter-spacing:.06em;color:#B8AE97;margin-bottom:30px}
.t-nightshade .footer{border-top:1px solid rgba(212,175,99,.2);padding:44px 24px;text-align:center;font-family:-apple-system,system-ui,sans-serif;font-size:11px;letter-spacing:.22em;text-transform:uppercase;color:#857D6A}
.t-nightshade .footer .made{margin-top:8px;color:#D4AF63;opacity:.8}
.t-nightshade .wa-float{background:#D4AF63;color:#14111A;box-shadow:0 10px 30px rgba(0,0,0,.55)}
`

// patternColors picks the glyph pattern ink per theme.
var patternColors = map[string][2]string{
	"ivory":      {"#8A6D3B", "0.07"},
	"bazaar":     {"#111111", "0.08"},
	"mint":       {"#0CA678", "0.06"},
	"nightshade": {"#D4AF63", "0.05"},
}

// themeCSS assembles the final stylesheet for a theme + category.
func themeCSS(theme, category string) string {
	var css string
	switch theme {
	case "bazaar":
		css = bazaarCSS
	case "mint":
		css = mintCSS
	case "nightshade":
		css = nightshadeCSS
	default:
		css = ivoryCSS
	}
	pc := patternColors[theme]
	if pc[0] == "" {
		pc = patternColors["ivory"]
	}
	pattern := PatternDataURI(category, pc[0], pc[1])
	return resetCSS + strings.ReplaceAll(css, "__PATTERN__", pattern)
}
