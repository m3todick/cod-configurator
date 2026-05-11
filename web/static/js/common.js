// ─── Shared Utilities ────────────────────────────────────────────

const API = '/api';

// Toast notifications
const toast = (() => {
  let container;
  const init = () => {
    container = document.createElement('div');
    container.className = 'toast-container';
    document.body.appendChild(container);
  };
  const show = (msg, type = 'info', dur = 3500) => {
    if (!container) init();
    const el = document.createElement('div');
    const icons = { success: '✓', error: '✕', info: 'ℹ' };
    el.className = `toast toast-${type}`;
    el.innerHTML = `<span style="font-size:16px">${icons[type]||'ℹ'}</span><span>${msg}</span>`;
    container.appendChild(el);
    requestAnimationFrame(() => { requestAnimationFrame(() => el.classList.add('show')); });
    setTimeout(() => {
      el.classList.remove('show');
      setTimeout(() => el.remove(), 350);
    }, dur);
  };
  return { show, success: m => show(m,'success'), error: m => show(m,'error'), info: m => show(m,'info') };
})();

// HTTP helpers
const http = {
  async get(url) {
    const r = await fetch(url, { credentials: 'include' });
    if (!r.ok) { const e = await r.json().catch(()=>({error:'Ошибка сети'})); throw new Error(e.error||'Ошибка'); }
    return r.json();
  },
  async post(url, data) {
    const r = await fetch(url, { method:'POST', credentials:'include', headers:{'Content-Type':'application/json'}, body: JSON.stringify(data) });
    if (!r.ok) { const e = await r.json().catch(()=>({error:'Ошибка сети'})); throw new Error(e.error||'Ошибка'); }
    return r.json();
  },
  async put(url, data) {
    const r = await fetch(url, { method:'PUT', credentials:'include', headers:{'Content-Type':'application/json'}, body: JSON.stringify(data) });
    if (!r.ok) { const e = await r.json().catch(()=>({error:'Ошибка'})); throw new Error(e.error||'Ошибка'); }
    return r.json();
  },
  async del(url) {
    const r = await fetch(url, { method:'DELETE', credentials:'include' });
    if (!r.ok) { const e = await r.json().catch(()=>({error:'Ошибка'})); throw new Error(e.error||'Ошибка'); }
    return r.json();
  }
};

// Formatting
const fmt = {
  money: n => new Intl.NumberFormat('ru-RU', { style:'currency', currency:'RUB', maximumFractionDigits:0 }).format(n),
  date: d => new Date(d).toLocaleDateString('ru-RU', { day:'2-digit', month:'long', year:'numeric' }),
  dateTime: d => new Date(d).toLocaleString('ru-RU', { day:'2-digit', month:'2-digit', year:'numeric', hour:'2-digit', minute:'2-digit' }),
};

// Modal helper
function openModal(id) { document.getElementById(id).classList.add('open'); }
function closeModal(id) { document.getElementById(id).classList.remove('open'); }

// Navbar scroll effect
window.addEventListener('scroll', () => {
  const nb = document.querySelector('.navbar');
  if (nb) nb.classList.toggle('scrolled', window.scrollY > 10);
});

// Mobile hamburger menu
function toggleNav() {
  const nav = document.querySelector('.navbar-nav');
  const btn = document.querySelector('.nav-toggle');
  if (!nav) return;
  const isOpen = nav.classList.toggle('open');
  if (btn) btn.classList.toggle('open', isOpen);
}
// Close nav on link click or outside click
document.addEventListener('DOMContentLoaded', () => {
  document.querySelectorAll('.navbar-nav a').forEach(a => {
    a.addEventListener('click', () => {
      document.querySelector('.navbar-nav')?.classList.remove('open');
      document.querySelector('.nav-toggle')?.classList.remove('open');
    });
  });
  document.addEventListener('click', e => {
    if (!e.target.closest('.navbar-nav') && !e.target.closest('.nav-toggle')) {
      document.querySelector('.navbar-nav')?.classList.remove('open');
      document.querySelector('.nav-toggle')?.classList.remove('open');
    }
  });
});

// Category labels
const catLabels = {
  server: 'Серверы', storage: 'Хранилища', network: 'Сеть',
  cooling: 'Охлаждение', power: 'Питание', security: 'Безопасность'
};
const catIcons = {
  server: '🖥', storage: '💾', network: '🌐',
  cooling: '❄️', power: '⚡', security: '🔒'
};

// Auth state
let currentUser = null;
async function loadCurrentUser() {
  try { currentUser = await http.get(`${API}/auth/me`); return currentUser; }
  catch { currentUser = null; return null; }
}
function updateNavAuth(user) {
  const loginLink = document.getElementById('nav-login');
  const cabinetLink = document.getElementById('nav-cabinet');
  const adminLink = document.getElementById('nav-admin');
  if (!loginLink) return;
  if (user) {
    loginLink.style.display = 'none';
    if (cabinetLink) cabinetLink.style.display = '';
    if (adminLink && user.role === 'admin') adminLink.style.display = '';
  } else {
    loginLink.style.display = '';
    if (cabinetLink) cabinetLink.style.display = 'none';
    if (adminLink) adminLink.style.display = 'none';
  }
}
