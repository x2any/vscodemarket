<script setup lang="ts">
import { t } from '../i18n'

defineProps<{
  topicKey: keyof typeof t.docs.topics
}>()
</script>

<template>
  <article class="docs-block">
    <h3>{{ t.docs.topics[topicKey] }}</h3>
    <slot />
  </article>
</template>

<style scoped>
.docs-block { margin-bottom: 2rem; }
.docs-block h3 { margin-bottom: .5rem; }
</style>