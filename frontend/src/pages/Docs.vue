<script setup lang="ts">
import { t } from '../i18n'
import DocsBlock from '../components/DocsBlock.vue'
import ScriptSnippet from '../components/ScriptSnippet.vue'

const topics = ['clientInstall', 'serverInstall', 'serverStart', 'offlineImport', 'faq'] as const
</script>

<template>
  <section>
    <h2>{{ t.docs.title }}</h2>
    <p class="intro">{{ t.docs.intro }}</p>

    <DocsBlock
      v-for="key in topics"
      :key="key"
      :topic-key="key"
    >
      <p>{{ t.docs.topics[key].body }}</p>
      <ScriptSnippet :code="t.docs.topics[key].snippet" />
    </DocsBlock>
  </section>
</template>

<style scoped>
.intro { opacity: .8; margin-bottom: 1.5rem; }
</style>