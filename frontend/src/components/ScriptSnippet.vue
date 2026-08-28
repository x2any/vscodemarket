<script setup lang="ts">
// Renders a run/install script block. NO wget / curl / Invoke-WebRequest / fetch
// is permitted in the snippet text — Constitution Principle III.
defineProps<{
  language?: 'sh' | 'ps1' | 'bat'
  code: string
}>()
</script>

<template>
  <pre class="snippet"><code>{{ code }}</code></pre>
</template>

<style scoped>
.snippet {
  background: #0b1020;
  color: #e6edf3;
  padding: .75rem 1rem;
  border-radius: 6px;
  overflow-x: auto;
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
  font-size: .85rem;
  line-height: 1.4;
}
</style>
