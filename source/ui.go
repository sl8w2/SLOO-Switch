package main

const pageHTML = `<!DOCTYPE html>
<html lang="ar" dir="rtl">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>SLOO Switch</title>
<link rel="icon" type="image/png" href="/icon.png">
<script>try{if(localStorage.getItem("theme")==="dark")document.documentElement.setAttribute("data-theme","dark")}catch(e){}</script>
<style>
:root{--bg:radial-gradient(1100px 560px at 85% -10%,rgba(240,138,36,.2),transparent 60%),radial-gradient(900px 520px at -10% 115%,rgba(217,164,65,.12),transparent 60%),repeating-linear-gradient(92deg,rgba(255,255,255,.014) 0 2px,rgba(0,0,0,.06) 2px 6px),linear-gradient(180deg,#1b100a,#110906);--s1:#3b2316;--s2:#2b1910;--line:#5e3c25;--text:#f6e7d2;--muted:#b89c82;--amber:#f08a24;--amber2:#ffb454;--amber3:#c96d14;--gold:#d9a441;--ok:#7bc47f;--bad:#e2574c;--chip:rgba(0,0,0,.28);--ghost:rgba(255,255,255,.05);--hdr:linear-gradient(180deg,rgba(0,0,0,.28),transparent);--ftr:rgba(0,0,0,.28);--rule:rgba(217,164,65,.22);--shadow:0 10px 26px rgba(0,0,0,.4);--shadow-h:0 18px 38px rgba(0,0,0,.5),0 0 0 1px rgba(240,138,36,.35),0 0 30px rgba(240,138,36,.16);--inset:inset 0 1px 0 rgba(255,255,255,.06);--input:rgba(0,0,0,.35);--overlay:rgba(10,5,2,.72);--toast:#24140c;--on-accent:#2b1405;--glow:rgba(240,138,36,.32)}
[data-theme=dark]{--bg:radial-gradient(900px 500px at 85% -10%,rgba(255,154,46,.17),transparent 60%),radial-gradient(700px 420px at -5% 110%,rgba(255,154,46,.07),transparent 60%),linear-gradient(180deg,#0f0f15,#07070a);--s1:#1c1c27;--s2:#14141b;--line:#2c2c3b;--text:#f3f3f7;--muted:#9d9db3;--amber:#ff9a2e;--amber2:#ffb454;--amber3:#e07d17;--gold:#ffb454;--ok:#8be08f;--bad:#ff6b5f;--chip:rgba(255,255,255,.05);--ghost:rgba(255,255,255,.06);--hdr:linear-gradient(180deg,rgba(255,255,255,.04),transparent);--ftr:rgba(255,255,255,.03);--rule:rgba(255,255,255,.08);--shadow:0 10px 30px rgba(0,0,0,.55);--shadow-h:0 18px 40px rgba(0,0,0,.65),0 0 0 1px rgba(255,154,46,.5),0 0 34px rgba(255,154,46,.18);--inset:inset 0 1px 0 rgba(255,255,255,.07);--input:rgba(255,255,255,.05);--overlay:rgba(4,4,8,.78);--toast:#1c1c27;--on-accent:#1a0d02;--glow:rgba(255,154,46,.3)}
*{box-sizing:border-box;margin:0;padding:0}
html,body{min-height:100%}
.tm,.tm *{transition:background-color .3s,border-color .3s,color .3s,box-shadow .3s!important}
body{font-family:"Segoe UI",Tahoma,"Noto Sans Arabic",Arial,sans-serif;color:var(--text);display:flex;flex-direction:column;min-height:100vh;background:var(--bg);background-attachment:fixed}
.skip{position:absolute;inset-inline-start:12px;top:-64px;z-index:60;padding:10px 16px;border-radius:10px;background:var(--amber);color:var(--on-accent);font-weight:700;text-decoration:none}
.skip:focus{top:12px}
header{display:flex;align-items:center;justify-content:space-between;padding:20px 34px;border-bottom:1px solid var(--rule);background:var(--hdr)}
.brand{display:flex;align-items:center;gap:16px}
.mark{width:52px;height:52px;flex:none;filter:drop-shadow(0 6px 14px var(--glow))}
h1{font:700 26px Georgia,"Times New Roman",serif;letter-spacing:6px;color:var(--gold);line-height:1}
.brand p{color:var(--muted);font-size:13px;margin-top:6px}
.tools{display:flex;gap:10px;align-items:center}
button{font-family:inherit;cursor:pointer;border:0;color:var(--text);transition:transform .15s,box-shadow .2s,background .2s,border-color .2s}
button:focus-visible,.card:focus-visible,input:focus-visible,a:focus-visible{outline:2px solid var(--amber2);outline-offset:2px}
.ghost{background:var(--ghost);border:1px solid var(--line);padding:9px 16px;border-radius:12px;font-size:14px;min-height:42px}
.ghost:hover{background:rgba(240,138,36,.14);border-color:var(--amber)}
.icon{width:42px;padding:0;display:grid;place-items:center}
.icon svg{width:20px;height:20px;fill:none;stroke:currentColor;stroke-width:1.9;stroke-linecap:round;stroke-linejoin:round}
.i-sun{display:none}
[data-theme=dark] .i-sun{display:block}
[data-theme=dark] .i-moon{display:none}
.primary{background:linear-gradient(145deg,var(--amber2),var(--amber) 60%,var(--amber3));color:var(--on-accent);font-weight:700;font-size:15px;padding:13px 24px;border-radius:14px;box-shadow:0 8px 24px var(--glow),inset 0 1px 0 rgba(255,255,255,.4)}
.primary:hover{transform:translateY(-2px);box-shadow:0 12px 30px var(--glow),inset 0 1px 0 rgba(255,255,255,.4)}
.actions{display:flex;gap:12px;flex-wrap:wrap}
.outline{background:transparent;border:1.5px solid var(--amber);color:var(--amber2);font-weight:700;font-size:15px;padding:12px 22px;border-radius:14px}
.outline:hover{background:rgba(240,138,36,.14);transform:translateY(-2px)}
.danger{background:linear-gradient(145deg,#f1705f,#d9483c);color:#1a0603;font-weight:700;padding:11px 22px;border-radius:12px}
main{width:100%;max-width:1020px;margin:0 auto;padding:32px 34px;flex:1;outline:0}
.bar{display:flex;align-items:center;justify-content:space-between;gap:18px;flex-wrap:wrap;margin-bottom:28px}
.status{display:flex;align-items:center;gap:10px;padding:11px 18px;border-radius:999px;background:var(--chip);border:1px solid var(--line);font-size:14px}
.dot{width:10px;height:10px;border-radius:50%;background:var(--muted);flex:none}
.status.ok .dot{background:var(--ok);box-shadow:0 0 10px var(--ok)}
.status.warn .dot{background:var(--amber);box-shadow:0 0 10px var(--amber)}
.notice{margin-bottom:22px;padding:13px 18px;border-radius:14px;background:rgba(226,87,76,.12);border:1px solid rgba(226,87,76,.45);font-size:14px}
.grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(220px,1fr));gap:20px;list-style:none}
.item{position:relative}
.card{display:block;width:100%;padding:26px 20px 20px;border-radius:20px;text-align:center;font:inherit;background:linear-gradient(160deg,var(--s1),var(--s2));border:1px solid var(--line);box-shadow:var(--shadow),var(--inset);transition:transform .2s,box-shadow .25s,border-color .25s}
.card:hover{transform:translateY(-5px);border-color:var(--amber);box-shadow:var(--shadow-h)}
.card.active{border-color:var(--gold);box-shadow:var(--shadow),0 0 0 1px var(--gold),0 0 28px rgba(217,164,65,.16)}
.av{width:72px;height:72px;border-radius:50%;margin:0 auto 14px;display:grid;place-items:center;font:700 30px Georgia,"Times New Roman",serif;color:#2b1405;box-shadow:0 6px 18px rgba(0,0,0,.45),inset 0 2px 0 rgba(255,255,255,.4)}
.g0{background:linear-gradient(145deg,#ffc069,#ee8420)}
.g1{background:linear-gradient(145deg,#f7d08a,#c8872c)}
.g2{background:linear-gradient(145deg,#ffb07a,#d4561f)}
.g3{background:linear-gradient(145deg,#e6c39a,#a8693a)}
.g4{background:linear-gradient(145deg,#ffd27a,#e0a02a)}
.g5{background:linear-gradient(145deg,#f4a261,#bc5a1c)}
.nm{display:block;font-size:18px;font-weight:700;word-break:break-word}
.sub{display:block;font-size:12px;color:var(--muted);margin-top:6px;direction:ltr;letter-spacing:.5px}
.tag{display:block;margin-top:14px;font-size:12px;color:var(--muted);min-height:18px}
.card:hover .tag.hint{color:var(--amber2)}
.badge{display:inline-flex;align-items:center;gap:6px;color:var(--gold);font-weight:700}
.badge i{width:8px;height:8px;border-radius:50%;background:var(--ok);box-shadow:0 0 8px var(--ok)}
.x{position:absolute;top:10px;inset-inline-end:10px;width:30px;height:30px;border-radius:50%;background:var(--chip);color:var(--muted);font-size:13px;opacity:0}
.item:hover .x,.item:focus-within .x{opacity:1}
.x:hover{background:var(--bad);color:#1a0603}
.ed{position:absolute;top:10px;inset-inline-start:10px;width:30px;height:30px;border-radius:50%;background:var(--chip);color:var(--muted);opacity:0;display:grid;place-items:center}
.ed svg{width:15px;height:15px;fill:none;stroke:currentColor;stroke-width:2;stroke-linecap:round;stroke-linejoin:round}
.item:hover .ed,.item:focus-within .ed{opacity:1}
.ed:hover{background:var(--amber);color:var(--on-accent)}
.box .chk{display:flex;align-items:center;gap:10px;margin:0 0 14px;color:var(--text);font-size:14px}
.box .chk input{width:18px;height:18px;padding:0;accent-color:var(--amber)}
.empty{grid-column:1/-1;text-align:center;padding:64px 20px;border:1px dashed var(--line);border-radius:22px;background:var(--chip);list-style:none}
.empty .ring2{width:78px;height:78px;border-radius:50%;margin:0 auto 18px;border:3px solid var(--gold);display:grid;place-items:center;font:700 34px Georgia,serif;color:var(--gold)}
.empty h2{font:700 20px Georgia,"Segoe UI",serif;margin-bottom:8px}
.empty p{color:var(--muted);font-size:14px}
footer{display:flex;align-items:center;justify-content:space-between;gap:14px;flex-wrap:wrap;padding:18px 34px;border-top:1px solid var(--rule);background:var(--ftr);font-size:13px;color:var(--muted)}
footer b{color:var(--gold);letter-spacing:2px;font-family:Georgia,serif}
.links{display:flex;gap:18px}
.links a{color:var(--amber2);text-decoration:none;font-weight:600;direction:ltr}
.links a:hover{color:var(--text);text-decoration:underline}
.modal,.busy{position:fixed;inset:0;display:grid;place-items:center;background:var(--overlay);backdrop-filter:blur(6px);z-index:20}
.hidden{display:none!important}
.box{width:min(440px,92vw);padding:28px;border-radius:22px;background:linear-gradient(160deg,var(--s1),var(--s2));border:1px solid var(--gold);box-shadow:0 24px 60px rgba(0,0,0,.6)}
.box h2{font:700 21px Georgia,"Segoe UI",serif;color:var(--gold);margin-bottom:16px}
.box p{color:var(--muted);font-size:14px;line-height:1.7;margin-bottom:12px}
.box label{display:block;font-size:13px;color:var(--muted);margin:12px 0 6px}
.box input{width:100%;padding:12px 14px;border-radius:12px;font:inherit;color:var(--text);background:var(--input);border:1px solid var(--line)}
.box input:focus{border-color:var(--amber)}
.val{direction:ltr;text-align:start;font-size:12px;color:var(--text);word-break:break-all;padding:10px 12px;border-radius:10px;background:var(--chip)}
.row{display:flex;gap:10px;justify-content:flex-end;margin-top:22px}
.busy{flex-direction:column;gap:20px;z-index:30}
.busy p{font-size:17px;color:var(--text)}
.ring{width:56px;height:56px;border-radius:50%;border:5px solid rgba(240,138,36,.2);border-top-color:var(--amber);animation:spin .9s linear infinite}
@keyframes spin{to{transform:rotate(360deg)}}
.toast{position:fixed;left:50%;bottom:28px;transform:translate(-50%,30px);opacity:0;pointer-events:none;padding:13px 24px;border-radius:14px;background:var(--toast);border:1px solid var(--gold);color:var(--text);font-size:14px;box-shadow:0 14px 34px rgba(0,0,0,.55);transition:opacity .25s,transform .25s;z-index:40;max-width:90vw;text-align:center}
.toast.show{opacity:1;transform:translate(-50%,0)}
.toast.err{border-color:var(--bad)}
@media(max-width:600px){header,main,footer{padding-inline:18px}h1{font-size:21px}}
@media(prefers-reduced-motion:reduce){*,*::before,*::after{animation-duration:.01ms!important;transition-duration:.01ms!important}.ring{animation-duration:2.4s!important}.card:hover,.primary:hover,.outline:hover{transform:none}}
</style>
</head>
<body>
<a class="skip" href="#main" data-i="skip">تخطّي إلى المحتوى</a>
<header>
<div class="brand"><svg class="mark" viewBox="0 0 512 512" aria-hidden="true" focusable="false"><defs><linearGradient id="mg" x1="0.1" y1="0" x2="0.9" y2="1"><stop offset="0" stop-color="#ffc56e"/><stop offset="0.55" stop-color="#f08a24"/><stop offset="1" stop-color="#b8610f"/></linearGradient><linearGradient id="ms" x1="0" y1="0" x2="0" y2="1"><stop offset="0" stop-color="#3a1d09"/><stop offset="1" stop-color="#1f0e03"/></linearGradient><linearGradient id="mh" x1="0" y1="0" x2="0" y2="1"><stop offset="0" stop-color="#fff" stop-opacity="0.38"/><stop offset="1" stop-color="#fff" stop-opacity="0"/></linearGradient><clipPath id="mc"><rect x="20" y="20" width="472" height="472" rx="118"/></clipPath></defs><rect x="20" y="20" width="472" height="472" rx="118" fill="url(#mg)"/><rect x="20" y="20" width="472" height="300" fill="url(#mh)" clip-path="url(#mc)"/><path d="M350 170C338 126 300 106 256 106C206 106 166 134 166 178C166 222 206 238 256 252C310 267 350 288 350 338C350 386 312 408 256 408C208 408 168 388 158 344" fill="none" stroke="url(#ms)" stroke-width="60" stroke-linecap="round" stroke-linejoin="round"/></svg><div><h1>SLOO</h1><p data-i="sub"></p></div></div>
<nav class="tools" data-a="tools">
<button id="theme" class="ghost icon" type="button" aria-pressed="false"><svg class="i-moon" viewBox="0 0 24 24" aria-hidden="true"><path d="M21 12.8A9 9 0 1 1 11.2 3a7 7 0 0 0 9.8 9.8z"/></svg><svg class="i-sun" viewBox="0 0 24 24" aria-hidden="true"><circle cx="12" cy="12" r="4"/><path d="M12 2v2M12 20v2M4.9 4.9l1.4 1.4M17.7 17.7l1.4 1.4M2 12h2M20 12h2M4.9 19.1l1.4-1.4M17.7 6.3l1.4-1.4"/></svg></button>
<button id="lang" class="ghost" type="button" data-a="langlabel"></button>
<button id="gear" class="ghost icon" type="button" data-a="settings"><svg viewBox="0 0 24 24" aria-hidden="true"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.7 1.7 0 0 0 .3 1.8l.1.1a2 2 0 1 1-2.8 2.8l-.1-.1a1.7 1.7 0 0 0-1.8-.3 1.7 1.7 0 0 0-1 1.5V21a2 2 0 1 1-4 0v-.1a1.7 1.7 0 0 0-1.1-1.5 1.7 1.7 0 0 0-1.8.3l-.1.1a2 2 0 1 1-2.8-2.8l.1-.1a1.7 1.7 0 0 0 .3-1.8 1.7 1.7 0 0 0-1.5-1H3a2 2 0 1 1 0-4h.1a1.7 1.7 0 0 0 1.5-1.1 1.7 1.7 0 0 0-.3-1.8l-.1-.1a2 2 0 1 1 2.8-2.8l.1.1a1.7 1.7 0 0 0 1.8.3H9a1.7 1.7 0 0 0 1-1.5V3a2 2 0 1 1 4 0v.1a1.7 1.7 0 0 0 1 1.5 1.7 1.7 0 0 0 1.8-.3l.1-.1a2 2 0 1 1 2.8 2.8l-.1.1a1.7 1.7 0 0 0-.3 1.8V9a1.7 1.7 0 0 0 1.5 1H21a2 2 0 1 1 0 4h-.1a1.7 1.7 0 0 0-1.5 1z"/></svg></button>
</nav>
</header>
<main id="main" tabindex="-1">
<div id="notice" class="notice hidden" role="alert" aria-live="assertive"></div>
<section class="bar"><div id="status" class="status" role="status" aria-live="polite"></div><div class="actions"><button id="newlogin" class="outline" type="button"></button><button id="add" class="primary" type="button"></button></div></section>
<ul id="grid" class="grid"></ul>
</main>
<footer>
<span>&copy; 2026 <b>SLOO</b> &middot; <span data-i="rights"></span> &middot; v__VERSION__</span>
<nav class="links" data-a="social"><a href="https://www.tiktok.com/@sl8w2" target="_blank" rel="noopener noreferrer">TikTok @sl8w2</a><a href="https://www.instagram.com/sl8w2" target="_blank" rel="noopener noreferrer">Instagram @sl8w2</a></nav>
</footer>
<div id="modal" class="modal hidden"><div id="box" class="box" role="dialog" aria-modal="true" aria-labelledby="dlgTitle"></div></div>
<div id="busy" class="busy hidden" role="alert" aria-live="assertive"><div class="ring"></div><p id="busyText"></p></div>
<div id="toast" class="toast" role="status" aria-live="polite"></div>
<script>
var T="__TOKEN__";
var I={
ar:{skip:"تخطّي إلى المحتوى",sub:"مبدّل حسابات Epic Games",rights:"جميع الحقوق محفوظة",add:"+ حفظ الحساب الحالي",newlogin:"تسجيل حساب جديد",rename:"تسمية",rename_title:"إعادة تسمية الحساب",new_auto_msg:"تم حفظ الحساب تلقائياً باسم «{0}». يمكنك تغيير الاسم الآن.",auto_added:"تم حفظ حساب جديد تلقائياً: {0}",autosave:"حفظ الحسابات تلقائياً عند تسجيل الدخول",renamed_ok:"تم تغيير الاسم",busy_new:"جارٍ فتح Epic لتسجيل حساب جديد...",new_title:"حفظ الحساب الحالي أولاً",new_msg:"الحساب الحالي غير محفوظ. اكتب له اسماً ليُحفظ، ثم يُفتح Epic على شاشة تسجيل الدخول.",save_and_new:"حفظ ومتابعة",new_done:"سجّل دخولك بالحساب الجديد في Epic ثم ارجع هنا",new_nolaunch:"شغّل Epic يدوياً وسجّل دخولك بالحساب الجديد",new_detected:"تم اكتشاف الحساب الجديد. اكتب له اسماً",e_unsaved_current:"احفظ الحساب الحالي أولاً.",cur_saved:"الحساب الحالي: {0}",cur_unsaved:"الحساب الحالي غير محفوظ",cur_none:"لا يوجد حساب مسجّل في Epic",exe_missing:"لم يتم العثور على Epic Games. حدّد مساره من الإعدادات.",empty_t:"لا توجد حسابات بعد",empty_d:"سجّل دخولك في Epic ثم اضغط «حفظ الحساب الحالي».",active:"نشط",hint:"اضغط للتبديل",save_title:"حفظ الحساب الحالي",name_ph:"اسم الحساب",save:"حفظ",cancel:"إلغاء",del_title:"حذف الحساب",del_msg:"حذف «{0}» من التطبيق؟ لن يتأثر حسابك في Epic.",del:"حذف",settings:"الإعدادات",tools:"الأدوات",social:"حسابات SLOO",langlabel:"اللغة",theme:"تبديل المظهر",theme_dark:"الوضع الداكن",theme_wood:"المظهر الخشبي",exe_auto:"المسار المكتشف",data_path:"مكان البيانات",exe_custom:"مسار مخصص (اختياري)",exe_ph:"C:\\Program Files (x86)\\Epic Games\\...\\EpicGamesLauncher.exe",reset:"افتراضي",none:"غير موجود",switching:"جارٍ التبديل إلى {0}...",done_launch:"تم التبديل إلى {0} — جارٍ تشغيل Epic",done_nolaunch:"تم التبديل إلى {0} — شغّل Epic يدوياً",already:"هذا هو الحساب الحالي",saved_ok:"تم حفظ {0}",deleted_ok:"تم الحذف",settings_ok:"تم حفظ الإعدادات",lost:"انقطع الاتصال بالتطبيق",e_no_login:"سجّل دخولك في Epic أولاً مع تفعيل «تذكّرني».",e_no_id:"تعذّر تحديد هوية الحساب. سجّل الدخول في Epic ثم أعد المحاولة.",e_name_taken:"الاسم «{0}» مستخدم.",e_id_taken:"هذا الحساب محفوظ مسبقاً باسم «{0}».",e_bad_name:"اكتب اسماً من 1 إلى 40 حرفاً.",e_not_found:"الحساب غير موجود.",e_snapshot_missing:"ملف الحساب المحفوظ مفقود.",e_kill_failed:"تعذّر إغلاق Epic. أغلقه يدوياً وأعد المحاولة.",e_bad_path:"المسار غير صحيح. اختر ملف EpicGamesLauncher.exe.",e_internal:"حدث خطأ غير متوقع."},
en:{skip:"Skip to content",sub:"Epic Games Account Switcher",rights:"All rights reserved",add:"+ Save current account",newlogin:"New login",rename:"Rename",rename_title:"Rename account",new_auto_msg:"The account was saved automatically as “{0}”. You can rename it now.",auto_added:"New account saved automatically: {0}",autosave:"Save accounts automatically when I log in",renamed_ok:"Renamed",busy_new:"Opening Epic for a new login...",new_title:"Save the current account first",new_msg:"The current account is not saved. Name it to save it, then Epic opens at the login screen.",save_and_new:"Save and continue",new_done:"Log in with the new account in Epic, then come back here",new_nolaunch:"Start Epic manually and log in with the new account",new_detected:"New account detected. Give it a name",e_unsaved_current:"Save the current account first.",cur_saved:"Current account: {0}",cur_unsaved:"Current account is not saved",cur_none:"No account logged in to Epic",exe_missing:"Epic Games was not found. Set its path in Settings.",empty_t:"No accounts yet",empty_d:"Log in to Epic, then press “Save current account”.",active:"Active",hint:"Click to switch",save_title:"Save current account",name_ph:"Account name",save:"Save",cancel:"Cancel",del_title:"Delete account",del_msg:"Delete “{0}” from the app? Your Epic account is not affected.",del:"Delete",settings:"Settings",tools:"Tools",social:"SLOO accounts",langlabel:"Language",theme:"Switch theme",theme_dark:"Dark mode",theme_wood:"Wood theme",exe_auto:"Detected path",data_path:"Data location",exe_custom:"Custom path (optional)",exe_ph:"C:\\Program Files (x86)\\Epic Games\\...\\EpicGamesLauncher.exe",reset:"Default",none:"Not found",switching:"Switching to {0}...",done_launch:"Switched to {0} — starting Epic",done_nolaunch:"Switched to {0} — start Epic manually",already:"This is the current account",saved_ok:"{0} saved",deleted_ok:"Deleted",settings_ok:"Settings saved",lost:"Lost connection to the app",e_no_login:"Log in to Epic first with “Remember me” enabled.",e_no_id:"Could not identify the account. Log in to Epic and retry.",e_name_taken:"The name “{0}” is taken.",e_id_taken:"This account is already saved as “{0}”.",e_bad_name:"Enter a name of 1 to 40 characters.",e_not_found:"Account not found.",e_snapshot_missing:"The saved account file is missing.",e_kill_failed:"Could not close Epic. Close it manually and retry.",e_bad_path:"Invalid path. Pick EpicGamesLauncher.exe.",e_internal:"Unexpected error."}
};
var lang="ar",theme="wood",S=null,lastJSON="",toastTimer=0,lastFocus=null,switching=false,awaitingNew=false,preKeys=[],known=null,silent=false;
try{lang=localStorage.getItem("lang")||"ar";theme=localStorage.getItem("theme")||"wood"}catch(e){}
function $(id){return document.getElementById(id)}
function t(k,a){var s=I[lang][k]||k;return a===undefined?s:s.replace("{0}",a)}
function h(tag,cls,txt){var e=document.createElement(tag);if(cls)e.className=cls;if(txt!==undefined)e.textContent=txt;return e}
function api(p,b){
return fetch("/api/"+p,{method:"POST",headers:{"Content-Type":"application/json","X-SLOO":T},body:JSON.stringify(b||{})}).then(function(r){return r.json()}).catch(function(){return{error:"lost"}})
}
function errText(r){return r.error==="lost"?t("lost"):t("e_"+r.error,r.arg)}
function toast(msg,err){
var e=$("toast");e.textContent=msg;e.className="toast show"+(err?" err":"");
clearTimeout(toastTimer);toastTimer=setTimeout(function(){e.className="toast"},3600)
}
function busy(on,txt){$("busy").className="busy"+(on?"":" hidden");if(txt)$("busyText").textContent=txt}
function modalOpen(){return $("modal").className.indexOf("hidden")<0}
function closeDialog(){$("modal").className="modal hidden";if(lastFocus&&lastFocus.focus)lastFocus.focus();lastFocus=null}
function dialog(title,nodes,okText,onOk,danger){
lastFocus=document.activeElement;
var box=$("box");box.textContent="";
var hd=h("h2","",title);hd.id="dlgTitle";box.appendChild(hd);
nodes.forEach(function(n){box.appendChild(n)});
var row=h("div","row"),c=h("button","ghost",t("cancel")),o=h("button",danger?"danger":"primary",okText);
c.type="button";o.type="button";c.onclick=closeDialog;o.onclick=function(){onOk()};
row.appendChild(c);row.appendChild(o);box.appendChild(row);
$("modal").className="modal";
var inp=box.querySelector("input");
if(inp){inp.onkeydown=function(e){if(e.key==="Enter")onOk()};inp.focus()}else if(danger)c.focus();else o.focus()
}
function hue(key){var n=0;for(var i=0;i<key.length;i++)n+=key.charCodeAt(i);return"g"+(n%6)}
function shortId(id){return id.length>14?id.slice(0,6)+"..."+id.slice(-6):id}
function render(){
if(!S)return;
var st=$("status");st.textContent="";st.appendChild(h("i","dot"));
var cur=S.current,txt,cls="status";
if(cur.loggedIn&&cur.matched){var m=S.accounts.filter(function(x){return x.key===cur.matched})[0];txt=t("cur_saved",m?m.name:"");cls+=" ok"}
else if(cur.loggedIn){txt=t("cur_unsaved");cls+=" warn"}
else txt=t("cur_none");
if(awaitingNew&&cur.loggedIn&&!modalOpen()){
awaitingNew=false;
if(!cur.matched){toast(t("new_detected"));setTimeout(askSave,150)}
else if(preKeys.indexOf(cur.matched)<0){var na=S.accounts.filter(function(x){return x.key===cur.matched})[0];if(na)setTimeout(function(){askRename(na,true)},150)}
}
st.className=cls;st.appendChild(h("span","",txt));
var nt=$("notice");
if(!S.exeFound){nt.textContent=t("exe_missing");nt.className="notice"}else nt.className="notice hidden";
var g=$("grid");g.textContent="";
if(!S.accounts.length){
var em=h("li","empty");em.appendChild(h("div","ring2","S"));em.appendChild(h("h2","",t("empty_t")));em.appendChild(h("p","",t("empty_d")));g.appendChild(em);return
}
S.accounts.forEach(function(a){
var li=h("li","item");
var c=h("button","card"+(a.active?" active":""));c.type="button";
if(a.active)c.setAttribute("aria-current","true");
c.appendChild(h("span","av "+hue(a.key),Array.from(a.name)[0].toUpperCase()));
c.appendChild(h("span","nm",a.name));c.appendChild(h("span","sub",shortId(a.id)));
var tag=h("span","tag"+(a.active?"":" hint"));
if(a.active){var b=h("span","badge");b.appendChild(h("i"));b.appendChild(h("span","",t("active")));tag.appendChild(b)}else tag.textContent=t("hint");
c.appendChild(tag);
c.onclick=function(){doSwitch(a)};
var x=h("button","x","\u2715");x.type="button";x.setAttribute("aria-label",t("del")+" "+a.name);
x.onclick=function(){askDelete(a)};
var ed=h("button","ed");ed.type="button";ed.setAttribute("aria-label",t("rename")+" "+a.name);
ed.innerHTML='<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M12 20h9"/><path d="M16.5 3.5a2.1 2.1 0 0 1 3 3L7 19l-4 1 1-4z"/></svg>';
ed.onclick=function(){askRename(a,false)};
li.appendChild(c);li.appendChild(ed);li.appendChild(x);g.appendChild(li)
})
}
function noteNew(){
var keys=S.accounts.map(function(a){return a.key});
if(known&&!silent&&!awaitingNew){S.accounts.forEach(function(a){if(known.indexOf(a.key)<0)toast(t("auto_added",a.name))})}
silent=false;
known=keys
}
function load(){
return api("state").then(function(r){
if(r.error)return;
var j=JSON.stringify(r);
if(j===lastJSON)return;
lastJSON=j;S=r;noteNew();render()
})
}
function doSwitch(a){
if(switching)return;
if(a.active&&S.epicRunning){toast(t("already"));return}
switching=true;busy(true,t("switching",a.name));
api("switch",{key:a.key}).then(function(r){
switching=false;busy(false);
if(r.error)toast(errText(r),true);else toast(t(r.launched?"done_launch":"done_nolaunch",r.name));
load();setTimeout(load,1500)
})
}
function askSave(){
var inp=h("input");inp.type="text";inp.maxLength=40;inp.placeholder=t("name_ph");inp.setAttribute("aria-label",t("name_ph"));
dialog(t("save_title"),[inp],t("save"),function(){
api("save",{name:inp.value}).then(function(r){
if(r.error){toast(errText(r),true);return}
silent=true;closeDialog();toast(t("saved_ok",r.name));load()
})
})
}
function startNew(){
busy(true,t("busy_new"));
api("newlogin").then(function(r){
busy(false);
if(r.error==="unsaved_current"){askSaveFirst();return}
if(r.error){toast(errText(r),true);return}
preKeys=S?S.accounts.map(function(a){return a.key}):[];awaitingNew=true;toast(t(r.launched?"new_done":"new_nolaunch"));
load();setTimeout(load,1500)
})
}
function askSaveFirst(){
var inp=h("input");inp.type="text";inp.maxLength=40;inp.placeholder=t("name_ph");inp.setAttribute("aria-label",t("name_ph"));
dialog(t("new_title"),[h("p","",t("new_msg")),inp],t("save_and_new"),function(){
api("save",{name:inp.value}).then(function(r){
if(r.error){toast(errText(r),true);return}
silent=true;closeDialog();load();startNew()
})
})
}
function askRename(a,fresh){
var inp=h("input");inp.type="text";inp.maxLength=40;inp.value=a.name;inp.setAttribute("aria-label",t("name_ph"));
var nodes=fresh?[h("p","",t("new_auto_msg",a.name)),inp]:[inp];
dialog(t("rename_title"),nodes,t("save"),function(){
if(inp.value.trim()===a.name){closeDialog();return}
api("rename",{key:a.key,name:inp.value}).then(function(r){
if(r.error){toast(errText(r),true);return}
closeDialog();toast(t("renamed_ok"));load()
})
});
inp.select()
}
function askDelete(a){
dialog(t("del_title"),[h("p","",t("del_msg",a.name))],t("del"),function(){
api("delete",{key:a.key}).then(function(r){
if(r.error){toast(errText(r),true);return}
closeDialog();toast(t("deleted_ok"));load()
})
},true)
}
function openSettings(){
var auto=h("div","val",S&&S.exeAuto?S.exeAuto:t("none"));
var inp=h("input");inp.type="text";inp.id="exeIn";inp.placeholder=t("exe_ph");inp.value=S?S.exePath:"";
var l1=h("label","",t("exe_auto")),l2=h("label","",t("exe_custom")),l3=h("label","",t("data_path")),dp=h("div","val",S&&S.dataPath?S.dataPath:"-");l2.setAttribute("for","exeIn");
var cb=h("input");cb.type="checkbox";cb.id="autoCb";cb.checked=!S||S.autoSave!==false;
var chk=h("label","chk");chk.setAttribute("for","autoCb");chk.appendChild(cb);chk.appendChild(h("span","",t("autosave")));
dialog(t("settings"),[chk,l1,auto,l2,inp,l3,dp],t("save"),function(){
api("autosave",{on:cb.checked}).then(function(){return api("exe",{path:inp.value})}).then(function(r){
if(r.error){toast(errText(r),true);return}
closeDialog();toast(t("settings_ok"));load()
})
});
var row=$("box").querySelector(".row"),rs=h("button","ghost",t("reset"));
rs.type="button";rs.onclick=function(){inp.value="";inp.focus()};
row.insertBefore(rs,row.firstChild)
}
function applyTheme(anim){
var d=document.documentElement;
if(anim){d.classList.add("tm");setTimeout(function(){d.classList.remove("tm")},350)}
if(theme==="dark")d.setAttribute("data-theme","dark");else d.removeAttribute("data-theme");
var b=$("theme");b.setAttribute("aria-pressed",theme==="dark"?"true":"false");
b.setAttribute("aria-label",t("theme"));b.title=t(theme==="dark"?"theme_wood":"theme_dark")
}
function applyLang(){
document.documentElement.lang=lang;document.documentElement.dir=lang==="ar"?"rtl":"ltr";
document.querySelectorAll("[data-i]").forEach(function(e){e.textContent=t(e.getAttribute("data-i"))});
document.querySelectorAll("[data-a]").forEach(function(e){e.setAttribute("aria-label",t(e.getAttribute("data-a")))});
$("add").textContent=t("add");$("newlogin").textContent=t("newlogin");$("lang").textContent=lang==="ar"?"EN":"\u0639\u0631\u0628\u064a";
$("gear").title=t("settings");applyTheme(false);render();api("lang",{lang:lang})
}
$("theme").onclick=function(){theme=theme==="dark"?"wood":"dark";try{localStorage.setItem("theme",theme)}catch(e){}applyTheme(true)};
$("lang").onclick=function(){lang=lang==="ar"?"en":"ar";try{localStorage.setItem("lang",lang)}catch(e){}applyLang()};
$("add").onclick=askSave;
$("newlogin").onclick=startNew;
$("gear").onclick=openSettings;
$("modal").onclick=function(e){if(e.target===$("modal"))closeDialog()};
document.addEventListener("keydown",function(e){
if(!modalOpen())return;
if(e.key==="Escape"){closeDialog();return}
if(e.key!=="Tab")return;
var f=Array.prototype.slice.call($("box").querySelectorAll("button,input"));
if(!f.length)return;
var first=f[0],last=f[f.length-1],cur=document.activeElement;
if(f.indexOf(cur)<0||(e.shiftKey&&cur===first)){e.preventDefault();last.focus()}
else if(!e.shiftKey&&cur===last){e.preventDefault();first.focus()}
});
document.querySelectorAll(".links a").forEach(function(a){a.onclick=function(e){e.preventDefault();api("open",{url:a.getAttribute("href")})}});
document.addEventListener("visibilitychange",function(){if(!document.hidden)load()});
window.addEventListener("pagehide",function(){navigator.sendBeacon("/api/bye?t="+T)});
applyLang();load();api("ping");
setInterval(function(){api("ping")},5000);
setInterval(function(){if(!document.hidden)load()},8000);
setInterval(function(){if(awaitingNew)load()},3000);
</script>
</body>
</html>
`
