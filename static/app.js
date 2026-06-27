/* ── Config ──────────────────────────────────────────────────── */
// When deploying frontend to GitHub Pages, set this to your live API URL:
// e.g. 'https://my-portfolio.up.railway.app'
// Leave empty to use relative URLs (works when Go serves both).
const API_BASE = window.PORTFOLIO_API_BASE || '';

const api = (path) => `${API_BASE}/api${path}`;

/* ── Theme ───────────────────────────────────────────────────── */
(function initTheme() {
  const saved = localStorage.getItem('theme') || 'dark';
  document.documentElement.setAttribute('data-theme', saved);
})();

function toggleTheme() {
  const current = document.documentElement.getAttribute('data-theme');
  const next = current === 'dark' ? 'light' : 'dark';
  document.documentElement.setAttribute('data-theme', next);
  localStorage.setItem('theme', next);
}

/* ── Social icons (SVG inline) ───────────────────────────────── */
const ICONS = {
  github: `<svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor"><path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0 0 24 12c0-6.63-5.37-12-12-12Z"/></svg>`,
  linkedin:`<svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor"><path d="M20.447 20.452h-3.554v-5.569c0-1.328-.027-3.037-1.852-3.037-1.853 0-2.136 1.445-2.136 2.939v5.667H9.351V9h3.414v1.561h.046c.477-.9 1.637-1.85 3.37-1.85 3.601 0 4.267 2.37 4.267 5.455v6.286zM5.337 7.433a2.062 2.062 0 0 1-2.063-2.065 2.064 2.064 0 1 1 2.063 2.065zm1.782 13.019H3.555V9h3.564v11.452zM22.225 0H1.771C.792 0 0 .774 0 1.729v20.542C0 23.227.792 24 1.771 24h20.451C23.2 24 24 23.227 24 22.271V1.729C24 .774 23.2 0 22.222 0h.003z"/></svg>`,
  twitter:`<svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor"><path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z"/></svg>`,
  mail:     `<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="4" width="20" height="16" rx="2"/><path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7"/></svg>`,
  whatsapp: `<svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor"><path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51a12.8 12.8 0 0 0-.57-.01c-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 0 1-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 0 1-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 0 1 2.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0 0 12.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 0 0 5.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 0 0-3.48-8.413Z"/></svg>`,
};

/* ── Fetch helper ────────────────────────────────────────────── */
async function fetchJSON(url) {
  const res = await fetch(url);
  if (!res.ok) throw new Error(`HTTP ${res.status}`);
  return res.json();
}

/* ── Intersection Observer (reveal + skill bars) ─────────────── */
const revealObserver = new IntersectionObserver((entries) => {
  entries.forEach(entry => {
    if (!entry.isIntersecting) return;
    entry.target.classList.add('visible');
    revealObserver.unobserve(entry.target);
  });
}, { threshold: 0.12 });

function observeReveal(el) { revealObserver.observe(el); }

const barObserver = new IntersectionObserver((entries) => {
  entries.forEach(entry => {
    if (!entry.isIntersecting) return;
    const bar = entry.target;
    bar.style.width = bar.dataset.level + '%';
    barObserver.unobserve(bar);
  });
}, { threshold: 0.3 });

/* ── Render helpers ──────────────────────────────────────────── */
function clearSkeletons(el) { el.innerHTML = ''; }

/* ── Hero / About ────────────────────────────────────────────── */
async function loadAbout() {
  try {
    const d = await fetchJSON(api('/about'));

    // Hero
    document.title = `${d.name} — ${d.role}`;
    const heroName = document.getElementById('hero-name');
    heroName.innerHTML = `<span class="accent">${d.name.split(' ')[0]}</span> ${d.name.split(' ').slice(1).join(' ')}`;

    document.getElementById('hero-role').textContent = d.role;
    document.getElementById('hero-tagline').textContent = d.tagline;

    const heroSocial = document.getElementById('hero-social');
    heroSocial.innerHTML = d.social.map(s =>
      `<a href="${s.url}" target="_blank" rel="noopener" class="social-link">
        ${ICONS[s.icon] || ''} ${s.label}
      </a>`
    ).join('');

    // About section
    const photo = document.getElementById('about-photo');
    photo.src = d.photo;
    photo.alt = d.name;

    document.getElementById('about-location').textContent = `📍 ${d.location}`;
    document.getElementById('about-text').innerHTML =
      `<p>${d.bio}</p>
       ${d.available ? `<p style="margin-top:1rem"><span style="color:var(--accent-2);font-weight:600;">🟢 I'm open to new opportunities</span></p>` : ''}`;

    document.querySelectorAll('.reveal').forEach(observeReveal);
  } catch (e) {
    console.error('About load failed', e);
  }
}

/* ── Skills ──────────────────────────────────────────────────── */
async function loadSkills() {
  try {
    const groups = await fetchJSON(api('/skills'));
    const grid = document.getElementById('skills-grid');
    clearSkeletons(grid);

    groups.forEach(group => {
      const card = document.createElement('div');
      card.className = 'skill-group reveal';

      card.innerHTML = `
        <div class="skill-group-header">
          <span>${group.icon}</span>
          <span>${group.group}</span>
        </div>
        ${group.skills.map(s => `
          <div class="skill-item">
            <div class="skill-label">
              <span>${s.name}</span>
              <span>${s.level}%</span>
            </div>
            <div class="skill-bar-bg">
              <div class="skill-bar" data-level="${s.level}"></div>
            </div>
          </div>
        `).join('')}
      `;

      grid.appendChild(card);
      observeReveal(card);
      card.querySelectorAll('.skill-bar').forEach(b => barObserver.observe(b));
    });
  } catch (e) {
    console.error('Skills load failed', e);
  }
}

/* ── Projects ────────────────────────────────────────────────── */
async function loadProjects() {
  try {
    const projects = await fetchJSON(api('/projects'));
    const grid = document.getElementById('projects-grid');
    clearSkeletons(grid);

    projects.forEach(p => {
      const card = document.createElement('div');
      card.className = `project-card reveal${p.featured ? ' featured' : ''}`;

      card.innerHTML = `
        <div class="project-title">${p.title}${p.featured ? ' <span style="font-size:0.7rem;color:var(--accent-2);font-family:var(--font-code);font-weight:600;">★ FEATURED</span>' : ''}</div>
        <p class="project-desc">${p.description}</p>
        <div class="project-tags">
          ${p.tags.map(t => `<span class="tag">${t}</span>`).join('')}
        </div>
        <div class="project-actions">
          ${p.live ? `<a href="${p.live}" target="_blank" rel="noopener" class="btn btn-primary btn-sm">Live Demo ↗</a>` : ''}
          ${p.github ? `<a href="${p.github}" target="_blank" rel="noopener" class="btn btn-ghost btn-sm">${ICONS.github} GitHub</a>` : ''}
        </div>
      `;

      grid.appendChild(card);
      observeReveal(card);
    });
  } catch (e) {
    console.error('Projects load failed', e);
  }
}

/* ── Blog list ────────────────────────────────────────────────── */
async function loadBlog() {
  try {
    const posts = await fetchJSON(api('/blog'));
    const grid = document.getElementById('blog-grid');
    clearSkeletons(grid);

    if (posts.length === 0) {
      grid.innerHTML = `<p style="color:var(--text-muted)">No posts yet — check back soon.</p>`;
      return;
    }

    posts.forEach(post => {
      const card = document.createElement('div');
      card.className = 'blog-card reveal';
      card.setAttribute('role', 'button');
      card.setAttribute('tabindex', '0');
      card.setAttribute('aria-label', `Read ${post.title}`);

      card.innerHTML = `
        <div class="blog-meta">${formatDate(post.date)}</div>
        <div class="blog-card-title">${post.title}</div>
        <p class="blog-card-excerpt">${post.excerpt}</p>
        <span class="read-more">Read more →</span>
      `;

      card.addEventListener('click', () => openPost(post.slug));
      card.addEventListener('keydown', e => { if (e.key === 'Enter') openPost(post.slug); });

      grid.appendChild(card);
      observeReveal(card);
    });
  } catch (e) {
    console.error('Blog load failed', e);
  }
}

/* ── Blog post reader ─────────────────────────────────────────── */
async function openPost(slug) {
  const listEl = document.getElementById('blog-list');
  const postEl = document.getElementById('blog-post');
  const contentEl = document.getElementById('blog-post-content');

  listEl.classList.add('hidden');
  postEl.classList.remove('hidden');
  contentEl.innerHTML = '<div class="skeleton skeleton-text mb" style="height:2.5rem;width:70%"></div><div class="skeleton skeleton-text mb"></div><div class="skeleton skeleton-text" style="width:50%"></div>';

  // Scroll to blog section
  document.getElementById('blog').scrollIntoView({ behavior: 'smooth', block: 'start' });

  try {
    const post = await fetchJSON(api(`/blog/${slug}`));
    contentEl.innerHTML = `
      <div class="prose-header">
        <h1 class="prose-title">${post.title}</h1>
        <p class="prose-date">${formatDate(post.date)}</p>
      </div>
      <div class="prose-body">${post.html}</div>
    `;

    // Syntax highlight
    contentEl.querySelectorAll('pre code').forEach(block => {
      hljs.highlightElement(block);
    });

  } catch (e) {
    contentEl.innerHTML = '<p style="color:var(--text-muted)">Could not load post. Please try again.</p>';
  }
}

function closePost() {
  document.getElementById('blog-list').classList.remove('hidden');
  document.getElementById('blog-post').classList.add('hidden');
  document.getElementById('blog').scrollIntoView({ behavior: 'smooth', block: 'start' });
}

/* ── Contact form ────────────────────────────────────────────── */
async function submitContact() {
  const name    = document.getElementById('c-name').value.trim();
  const email   = document.getElementById('c-email').value.trim();
  const message = document.getElementById('c-message').value.trim();
  const status  = document.getElementById('form-status');
  const btn     = document.getElementById('send-btn');

  if (!name || !email || !message) {
    status.textContent = 'Please fill in all fields.';
    status.className = 'form-status error';
    return;
  }

  btn.disabled = true;
  btn.textContent = 'Sending…';
  status.textContent = '';
  status.className = 'form-status';

  try {
    const res = await fetch(api('/contact'), {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name, email, message }),
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.error || 'Server error');

    status.textContent = data.message || 'Message sent!';
    status.className = 'form-status success';
    document.getElementById('c-name').value = '';
    document.getElementById('c-email').value = '';
    document.getElementById('c-message').value = '';
  } catch (e) {
    status.textContent = 'Failed to send. Please try again.';
    status.className = 'form-status error';
  } finally {
    btn.disabled = false;
    btn.textContent = 'Send Message';
  }
}

/* ── Date formatter ──────────────────────────────────────────── */
function formatDate(str) {
  if (!str) return '';
  const d = new Date(str);
  return d.toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' });
}

/* ── Navbar scroll shadow ─────────────────────────────────────── */
window.addEventListener('scroll', () => {
  const nav = document.getElementById('navbar');
  nav.style.boxShadow = window.scrollY > 10
    ? '0 2px 24px rgba(0,0,0,0.3)'
    : 'none';
}, { passive: true });

/* ── Mobile hamburger ────────────────────────────────────────── */
document.getElementById('hamburger').addEventListener('click', () => {
  const drawer = document.getElementById('nav-drawer');
  drawer.classList.toggle('open');
});

document.querySelectorAll('.drawer-link').forEach(link => {
  link.addEventListener('click', () => {
    document.getElementById('nav-drawer').classList.remove('open');
  });
});

/* ── Event bindings ──────────────────────────────────────────── */
document.getElementById('theme-toggle').addEventListener('click', toggleTheme);
document.getElementById('back-btn').addEventListener('click', closePost);
document.getElementById('send-btn').addEventListener('click', submitContact);

/* ── Footer year ─────────────────────────────────────────────── */
document.getElementById('footer-year').textContent = new Date().getFullYear();

/* ── Boot ────────────────────────────────────────────────────── */
(async function init() {
  // Initial reveal pass for elements already in view
  document.querySelectorAll('.reveal').forEach(observeReveal);

  await Promise.all([
    loadAbout(),
    loadSkills(),
    loadProjects(),
    loadBlog(),
  ]);
})();

/* ── Floating buttons ────────────────────────────────────────── */
(function initFabs() {
  const fabTop = document.getElementById('fab-totop');
  if (!fabTop) return;
  window.addEventListener('scroll', () => {
    fabTop.classList.toggle('visible', window.scrollY > 400);
  }, { passive: true });

  // Sync WhatsApp URL from API
  fetch('/api/about').then(r => r.json()).then(d => {
    const wa = (d.social || []).find(s => s.icon === 'whatsapp');
    const btn = document.getElementById('fab-whatsapp');
    if (wa && wa.url && btn) btn.href = wa.url;
  }).catch(() => {});
})();

/* ── Carousel ────────────────────────────────────────────────── */
(function initCarousel() {
  let slides = [];
  let current = 0;
  let timer = null;
  const DURATION = 5000;

  async function loadCarousel() {
    try {
      slides = await fetchJSON('/api/carousel');

      // Only show enabled slides
      slides = (slides || []).filter(slide => slide.enabled !== false);

      const section = document.getElementById('carousel-section');
      if (!section) return;

      // No slides → hide the entire carousel
      if (slides.length === 0) {
        section.style.display = 'none';
        stopAuto();
        return;
      }

      section.style.display = '';

      current = 0;

      buildCarousel();
      goTo(0);
      startAuto();

    } catch (e) {
      console.warn('Carousel load failed', e);
    }
  }

  function buildCarousel() {
    const track = document.getElementById('carousel-track');
    const dots  = document.getElementById('carousel-dots');
    const bar   = document.getElementById('carousel-caption-bar');

    if (!track) return;

    // Fixed: Added closing curly brace to block below
    if (slides.length === 0) {
        document.getElementById('hero').style.display = 'none';
        return;
    }

    track.innerHTML = '';
    dots.innerHTML  = '';

    // Set CSS duration for progress bar
    bar.style.setProperty('--carousel-duration', DURATION + 'ms');

    slides.forEach((slide, i) => {
      // ── Build slide element
      const el = document.createElement('div');
      el.className = 'carousel-slide' + (i === 0 ? ' active' : '');
      el.dataset.index = i;

      if (slide.type === 'hero') {
        // Slide 1 — mirrors the hero section visually
        el.innerHTML = `
          <div style="position:absolute;inset:0;background:linear-gradient(135deg,rgba(124,58,237,0.22) 0%,rgba(6,182,212,0.12) 60%,transparent 100%)"></div>
          <div style="position:absolute;inset:0;background:var(--bg);opacity:0.82"></div>
          <div class="carousel-slide-content" style="z-index:2;position:relative">
            <div class="carousel-slide-eyebrow">Portfolio Highlights</div>
            <div class="carousel-slide-title" style="background:linear-gradient(90deg,var(--accent-1),var(--accent-2));-webkit-background-clip:text;-webkit-text-fill-color:transparent;background-clip:text">
              Welcome to My Work
            </div>
            <div class="carousel-slide-subtitle">
              Full-stack engineer crafting fast backends and pixel-perfect frontends. Scroll to explore.
            </div>
          </div>`;
      } else {
        // Image slides
        el.innerHTML = `
          <img class="carousel-slide-img" src="${slide.image || ''}" alt="${slide.title || ''}" loading="${i === 0 ? 'eager' : 'lazy'}" />
          <div class="carousel-slide-overlay"></div>
          <div class="carousel-slide-content">
            ${slide.caption ? `<div class="carousel-slide-eyebrow">${slide.caption}</div>` : ''}
            ${slide.title   ? `<div class="carousel-slide-title">${slide.title}</div>` : ''}
            ${slide.subtitle ? `<div class="carousel-slide-subtitle">${slide.subtitle}</div>` : ''}
          </div>`;
      }
      track.appendChild(el);

      // ── Dot
      const dot = document.createElement('button');
      dot.className = 'carousel-dot' + (i === 0 ? ' active' : '');
      dot.setAttribute('aria-label', `Go to slide ${i + 1}`);
      dot.addEventListener('click', () => goTo(i));
      dots.appendChild(dot);
    });

    // Counter
    const counter = document.createElement('div');
    counter.className = 'carousel-counter';
    counter.id = 'carousel-counter';
    document.getElementById('carousel-section').appendChild(counter);
    updateCounter();

    // Arrow bindings
    document.getElementById('carousel-prev').addEventListener('click', () => goTo(current - 1));
    document.getElementById('carousel-next').addEventListener('click', () => goTo(current + 1));

    // Touch / swipe
    let touchStartX = 0;
    track.addEventListener('touchstart', e => { touchStartX = e.touches[0].clientX; }, { passive: true });
    track.addEventListener('touchend',   e => {
      const dx = e.changedTouches[0].clientX - touchStartX;
      if (Math.abs(dx) > 50) goTo(dx < 0 ? current + 1 : current - 1);
    });

    // Pause on hover
    const section = document.getElementById('carousel-section');
    section.addEventListener('mouseenter', stopAuto);
    section.addEventListener('mouseleave', startAuto);
  }

  function goTo(idx) {
    const len = slides.length;
    current = ((idx % len) + len) % len;

    // Move track
    document.getElementById('carousel-track').style.transform = `translateX(-${current * 100}%)`;

    // Active class on slides (for ken-burns effect)
    document.querySelectorAll('.carousel-slide').forEach((s, i) => {
      s.classList.toggle('active', i === current);
    });

    // Dots
    document.querySelectorAll('.carousel-dot').forEach((d, i) => {
      d.classList.toggle('active', i === current);
    });

    updateCounter();

    // Restart progress bar animation
    const bar = document.getElementById('carousel-caption-bar');
    bar.style.animation = 'none';
    // Force reflow
    void bar.offsetWidth;
    bar.style.animation = '';

    // Restart auto timer
    stopAuto();
    startAuto();
  }

  function updateCounter() {
    const el = document.getElementById('carousel-counter');
    if (el) el.textContent = `${current + 1} / ${slides.length}`;
  }

  function startAuto() {
    if (timer) return;
    timer = setInterval(() => goTo(current + 1), DURATION);
  }

  function stopAuto() {
    clearInterval(timer);
    timer = null;
  }

  loadCarousel();
})();