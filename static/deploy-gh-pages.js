// ─────────────────────────────────────────────────────────────────────────────
// GitHub Pages deployment helper
//
// When your Go backend is deployed (Railway / Render), set PORTFOLIO_API_BASE
// to your live URL so the frontend JS knows where to call the API.
//
// Usage (run once before pushing to gh-pages branch):
//   node deploy-gh-pages.js https://my-portfolio.up.railway.app
//
// This patches app.js in-place with the live API URL, then you push
// the static/ folder to the gh-pages branch.
// ─────────────────────────────────────────────────────────────────────────────

const fs   = require('fs');
const path = require('path');

const liveURL = process.argv[2];
if (!liveURL) {
  console.error('Usage: node deploy-gh-pages.js <LIVE_API_URL>');
  process.exit(1);
}

const appJs = path.join(__dirname, 'app.js');
let src = fs.readFileSync(appJs, 'utf8');

// Replace the placeholder line with the real URL
src = src.replace(
  /const API_BASE = window\.PORTFOLIO_API_BASE \|\| '';/,
  `const API_BASE = '${liveURL}';`
);

fs.writeFileSync(appJs, src);
console.log(`✅  API_BASE set to: ${liveURL}`);
console.log(`    Push the static/ folder to your gh-pages branch.`);
