<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n-light' // intentionally avoided; use our own t
import { t } from '../i18n'

const props = defineProps<{
  title: string
  url: string
  platform?: string
  architecture?: string
  commitHash?: string
}>()

const sub = computed(() => {
  const parts = [props.platform, props.architecture].filter(Boolean)
  return parts.join(' / ')
})
</script>

<template>
  <el-card class="dl-card" shadow="hover">
    <h3 class="title">{{ title }}</h3>
    <p v-if="sub" class="meta">{{ sub }}</p>
    <p v-if="commitHash" class="commit">{{ t.home.commitHash }}: <code>{{ commitHash }}</code></p>
    <el-link :href="url" target="_blank" rel="noopener" type="primary">
      {{ url }}
    </el-link>
  </el-card>
</template>

<style scoped>
.dl-card { margin-bottom: 1rem; }
.title { margin: 0 0 .5rem; font-size: 1rem; }
.meta { margin: 0 0 .25rem; opacity: .75; font-size: .9rem; }
.commit { margin: 0 0 .5rem; font-size: .85rem; }
</style>