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
  background: #020617;
  color: #e5e7eb;
}

.app-shell {
  min-height: 100vh;
  background: radial-gradient(circle at top, #111827 0, #020617 55%);
  color: #e5e7eb;
}

.app-header {
  padding: 1.5rem 2rem;
  border-bottom: 1px solid #1f2937;
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
  color: #9ca3af;
  font-size: 0.875rem;
}

.app-nav {
  display: flex;
  gap: 0.6rem;
}

.app-nav a {
  color: #e5e7eb;
  text-decoration: none;
  padding: 0.45rem 0.9rem;
  border-radius: 9999px;
  border: 1px solid transparent;
  font-size: 0.9rem;
  transition: border-color 0.2s ease;
}

.app-nav a:hover {
  border-color: #1f2937;
}

.app-nav .active-link {
  border-color: #0ea5e9;
  color: #0ea5e9;
  background: rgba(14, 165, 233, 0.1);
}

.app-main {
  padding: 2rem clamp(1rem, 5vw, 3rem) 3rem;
  max-width: 1200px;
  margin: 0 auto;
}
</style>
