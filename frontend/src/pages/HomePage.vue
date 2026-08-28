<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { http } from '../api/http'
import { track } from '../api/track'
import { t } from '../i18n'
import VersionForm from '../components/VersionForm.vue'
import DownloadLinkCard from '../components/DownloadLinkCard.vue'
import type { LookupResponse, Platform, Architecture, ClientPayload, ServerPayload } from '../api/contracts'

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

async function onSubmit(payload: { channel: 'stable' | 'insider'; version: string; platform?: Platform; architecture?: Architecture }) {
  error.value = ''
  result.value = null
  try {
    result.value = await http.post<LookupResponse>('/versions/lookup', payload)
    if (result.value) {
      const r = result.value
      track({
        eventType: 'SEARCH',
        targetType: r.servers && r.servers.length > 0 ? 'SERVER' : 'CLIENT',
        targetIdentifier: r.version,
        platform: payload.platform,
        architecture: payload.architecture,
        channel: r.channel
      })
    }
  } catch (e) {
    error.value = e instanceof Error ? e.message : String(e)
  }
}

function onDownloadClick(_url: string, _kind: 'CLIENT' | 'SERVER') {
  if (!result.value) return
  track({
    eventType: 'DOWNLOAD',
    targetType: _kind,
    targetIdentifier: result.value.version,
    channel: result.value.channel
  })
}

function labelFor(c: ClientPayload | ServerPayload): string {
  return `${ c.platform } / ${ c.architecture }`
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
      <h3>{{ t.home.clientCard }}</h3>
      <div class="grid">
        <DownloadLinkCard
          v-for="(c, i) in result.clients"
          :key="`c-${c.platform}-${c.architecture}-${i}`"
          :title="labelFor(c)"
          :url="c.downloadUrl"
          :platform="c.platform"
          :architecture="c.architecture"
          @click="onDownloadClick(c.downloadUrl, 'CLIENT')"
        />
      </div>

      <template v-if="result.servers && result.servers.length">
        <h3>{{ t.home.serverCard }} <small v-if="result.commit">({{ t.home.commitHash }}: <code>{{ result.commit.slice(0, 12) }}</code>)</small></h3>
        <div class="grid">
          <DownloadLinkCard
            v-for="(s, i) in result.servers"
            :key="`s-${s.platform}-${s.architecture}-${i}`"
            :title="labelFor(s)"
            :url="s.downloadUrl"
            :platform="s.platform"
            :architecture="s.architecture"
            :commit-hash="s.commitHash"
            @click="onDownloadClick(s.downloadUrl, 'SERVER')"
          />
        </div>
      </template>
    </div>
  </section>
</template>

<style scoped>
.err { color: var(--el-color-danger); }
.hint { opacity: .75; margin: 1rem 0; }
.results { margin-top: 1rem; }
.results h3 { margin: 1rem 0 .5rem; font-size: 1rem; }
.grid {
  display: grid;
  gap: .75rem;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
}
</style>