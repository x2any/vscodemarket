<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { http } from '../api/http'
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
  } catch (e) {
    error.value = e instanceof Error ? e.message : String(e)
  }
}
</script>

<template>
  <section>
    <h2>{{ t.home.submit }}</h2>
    <VersionForm ref="formRef" @submit="onSubmit" @infer="onInfer" />

    <p v-if="error" class="err">{{ error }}</p>

    <div v-if="result" class="results">
      <DownloadLinkCard
        :title="t.home.clientCard"
        :url="result.client.downloadUrl"
        :platform="result.client.platform"
        :architecture="result.client.architecture"
      />
      <DownloadLinkCard
        v-if="result.server"
        :title="t.home.serverCard"
        :url="result.server.downloadUrl"
        :platform="result.server.platform"
        :architecture="result.server.architecture"
        :commit-hash="result.server.commitHash"
      />
    </div>
  </section>
</template>

<style scoped>
.err { color: var(--el-color-danger); }
.results { margin-top: 1rem; }
</style>