<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import { currentLocale, setLocale, t } from './i18n'

function toggleLocale() {
  setLocale(currentLocale.value === 'zh-CN' ? 'en-US' : 'zh-CN')
}
</script>

<template>
  <div class="app">
    <header class="bar">
      <strong class="brand">{{ t.app.title }}</strong>
      <nav class="links">
        <RouterLink to="/">{{ t.nav.home }}</RouterLink>
        <RouterLink to="/search">{{ t.nav.search }}</RouterLink>
        <RouterLink to="/releases">{{ t.nav.releases }}</RouterLink>
        <RouterLink to="/trending">{{ t.nav.trending }}</RouterLink>
        <RouterLink to="/docs">{{ t.nav.docs }}</RouterLink>
      </nav>
      <button class="lang" type="button" @click="toggleLocale">
        {{ currentLocale === 'zh-CN' ? 'EN' : '中' }}
      </button>
    </header>

    <main class="content">
      <RouterView />
    </main>

    <footer class="foot">
      <small>{{ t.app.intro }}</small>
    </footer>
  </div>
</template>

<style>
:root { color-scheme: light dark; }
body { margin: 0; font-family: system-ui, -apple-system, sans-serif; }
.app { display: flex; flex-direction: column; min-height: 100vh; }
.bar {
  display: flex; align-items: center; gap: 1rem;
  padding: .75rem 1rem; border-bottom: 1px solid #e5e7eb;
}
.brand { font-size: 1rem; }
.links { display: flex; gap: 1rem; flex: 1; flex-wrap: wrap; }
.links a { text-decoration: none; color: inherit; }
.links a.router-link-active { font-weight: 600; }
.lang { padding: .25rem .5rem; }
.content { padding: 1rem; flex: 1; }
.foot { padding: 1rem; border-top: 1px solid #e5e7eb; opacity: .75; }

@media (max-width: 1023px) {
  .bar { flex-wrap: wrap; }
}
@media (max-width: 767px) {
  .links { gap: .5rem; font-size: .9rem; }
  .content { padding: .75rem; }
}
</style>
