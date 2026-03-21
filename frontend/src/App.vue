<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
  <div class="app-shell">
    <header class="app-header">
      <div>
        <h1>F1 Data Dashboard</h1>
        <p>Explore OpenF1 data exposed by the Go backend.</p>
      </div>
      <nav class="app-nav">
        <RouterLink to="/" exact-active-class="active-link">Home</RouterLink>
        <RouterLink to="/drivers" exact-active-class="active-link">Drivers</RouterLink>
        <RouterLink to="/cars" exact-active-class="active-link">Cars</RouterLink>
      </nav>
    </header>

    <main class="app-main">
      <RouterView />
    </main>
  </div>
</template>

<style>
:root {
  color-scheme: dark;
  --bg-0: #0b0b0d;
  --bg-1: #111114;
  --bg-2: #151519;
  --panel: #141418;
  --panel-soft: #101013;
  --border: #242428;
  --border-strong: #2f2f36;
  --text: #f5f5f5;
  --muted: #b0b0b8;
  --accent: #e10600;
  --accent-2: #ff2a1a;
  --accent-soft: rgba(225, 6, 0, 0.14);
}

body {
  margin: 0;
  font-family:
    'Space Grotesk',
    system-ui,
    -apple-system,
    BlinkMacSystemFont,
    'Segoe UI',
    sans-serif;
  background: var(--bg-0);
  color: var(--text);
}

.app-shell {
  min-height: 100vh;
  background: radial-gradient(circle at top, #1b0a0a 0, var(--bg-0) 55%);
  color: var(--text);
}

.app-header {
  padding: 1.5rem 2rem;
  border-bottom: 1px solid var(--border);
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
}

.app-header h1 {
  margin: 0;
  font-size: 1.5rem;
}

.app-header p {
  margin: 0.25rem 0 0;
  color: var(--muted);
  font-size: 0.875rem;
}

.app-nav {
  display: flex;
  gap: 0.6rem;
}

.app-nav a {
  color: var(--text);
  text-decoration: none;
  padding: 0.45rem 0.9rem;
  border-radius: 9999px;
  border: 1px solid transparent;
  font-size: 0.9rem;
  transition: border-color 0.2s ease;
}

.app-nav a:hover {
  border-color: var(--border);
}

.app-nav .active-link {
  border-color: var(--accent);
  color: var(--accent);
  background: var(--accent-soft);
}

.app-main {
  padding: 2rem clamp(1rem, 5vw, 3rem) 3rem;
  max-width: 1200px;
  margin: 0 auto;
}
</style>
