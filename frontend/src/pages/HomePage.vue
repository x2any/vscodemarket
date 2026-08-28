<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { http } from '../api/http'
import { track } from '../api/track'
import { t } from '../i18n'
import VersionForm from '../components/VersionForm.vue'
import DownloadLinkCard from '../components/DownloadLinkCard.vue'
import type { LookupResponse, Platform, Architecture } from '../api/contracts'

const formRef = ref<InstanceType<typeof VersionForm> | null>(null)
const result = ref<LookupResponse | null>(null)
const error = ref<string>('')

onMounted(() => {
  formRef.value?.pingInfer()
})

async function onInfer(payload: { userAgent: string }) {
  try {
    const resp = await http.post<{ platform: Platform; architecture: Architecture }>(
      '/ua/infer',
      payload
    )
    formRef.value?.fillInferred(resp.platform, resp.architecture)
  } catch {
    // UA inference is best-effort; ignore failures.
  }
}

async function onSubmit(payload: Parameters<InstanceType<typeof VersionForm>['$emit']>[1] extends infer X ? X : never) {
  error.value = ''
  result.value = null
  try {
    result.value = await http.post<LookupResponse>('/versions/lookup', payload)
    if (result.value) {
      const r = result.value
      track({
        eventType: 'SEARCH',
        targetType: r.server ? 'SERVER' : 'CLIENT',
        targetIdentifier: r.version,
        platform: r.client.platform,
        architecture: r.client.architecture,
        channel: r.channel
      })
    }
  } catch (e) {
    error.value = e instanceof Error ? e.message : String(e)
  }
}

function onDownloadClick(url: string, kind: 'CLIENT' | 'SERVER') {
  track({
    eventType: 'DOWNLOAD',
    targetType: kind,
    targetIdentifier: result.value?.version ?? '',
    platform: result.value?.client.platform,
    architecture: result.value?.client.architecture,
    channel: result.value?.channel
  })
  // Open the link in a new tab; the browser blocks window.open from async
  // handlers only in some browsers, so rely on the link itself.
  void url
}
</script>

<template>
  <section>
    <h2>{{ t.home.submit }}</h2>
    <VersionForm ref="formRef" @submit="onSubmit" @infer="onInfer" />

    <p v-if="error" class="err">{{ error }}</p>

    <p v-if="!error && result" class="hint">
      <RouterLink to="/trending">{{ t.nav.trending }} →</RouterLink>
    </p>

    <div v-if="result" class="results">
      <DownloadLinkCard
        :title="t.home.clientCard"
        :url="result.client.downloadUrl"
        :platform="result.client.platform"
        :architecture="result.client.architecture"
        @click="onDownloadClick(result.client.downloadUrl, 'CLIENT')"
      />
      <DownloadLinkCard
        v-if="result.server"
        :title="t.home.serverCard"
        :url="result.server.downloadUrl"
        :platform="result.server.platform"
        :architecture="result.server.architecture"
        :commit-hash="result.server.commitHash"
        @click="onDownloadClick(result.server.downloadUrl, 'SERVER')"
      />
    </div>
  </section>
</template>

<style scoped>
.err { color: var(--el-color-danger); }
.results { margin-top: 1rem; }
</style>