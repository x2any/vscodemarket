<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { http } from '../api/http'
import { track } from '../api/track'
import { t } from '../i18n'
import DownloadLinkCard from '../components/DownloadLinkCard.vue'
import type { LookupResponse, ClientPayload, ServerPayload, Channel } from '../api/contracts'

const route = useRoute()

const channel = computed<Channel>(() => {
  const c = String(route.params.channel ?? 'stable')
  return c === 'insider' ? 'insider' : 'stable'
})
const version = computed(() => String(route.params.ver ?? ''))

const data = ref<LookupResponse | null>(null)
const loading = ref(false)
const error = ref<string>('')

async function load() {
  if (!version.value) return
  loading.value = true
  error.value = ''
  data.value = null
  try {
    data.value = await http.post<LookupResponse>('/versions/lookup', {
      channel: channel.value,
      version: version.value
    })
  } catch (e) {
    error.value = e instanceof Error ? e.message : String(e)
  } finally {
    loading.value = false
  }
}

onMounted(load)
watch([channel, version], load)

function onDownload(kind: 'CLIENT' | 'SERVER', c: ClientPayload | ServerPayload) {
  track({
    eventType: 'DOWNLOAD',
    targetType: kind,
    targetIdentifier: version.value,
    platform: c.platform,
    architecture: c.architecture,
    channel: channel.value
  })
}

function label(c: ClientPayload | ServerPayload): string {
  return `${c.platform} / ${c.architecture}`
}
</script>

<template>
  <section>
    <h2>{{ channel }} {{ version }}</h2>

    <p v-if="loading">Loading…</p>
    <p v-if="error" class="err">{{ error }}</p>

    <div v-if="data">
      <p v-if="data.commit" class="commit">
        {{ t.home.commitHash }}: <code>{{ data.commit.slice(0, 12) }}</code>
      </p>

      <h3>{{ t.home.clientCard }}</h3>
      <div class="grid">
        <DownloadLinkCard
          v-for="(c, i) in data.clients"
          :key="`c-${c.platform}-${c.architecture}-${i}`"
          :title="label(c)"
          :url="c.downloadUrl"
          :platform="c.platform"
          :architecture="c.architecture"
          @click="onDownload('CLIENT', c)"
        />
      </div>

      <template v-if="data.servers && data.servers.length">
        <h3>{{ t.home.serverCard }}</h3>
        <div class="grid">
          <DownloadLinkCard
            v-for="(s, i) in data.servers"
            :key="`s-${s.platform}-${s.architecture}-${i}`"
            :title="label(s)"
            :url="s.downloadUrl"
            :platform="s.platform"
            :architecture="s.architecture"
            :commit-hash="s.commitHash"
            @click="onDownload('SERVER', s)"
          />
        </div>
      </template>
      <p v-else-if="data.channel === 'stable'" class="empty">
        {{ t.versionDetail.serverUnavailable }}
      </p>
    </div>
  </section>
</template>

<style scoped>
.err { color: var(--el-color-danger); }
.commit { opacity: .85; }
.empty { opacity: .65; font-style: italic; }
.grid {
  display: grid;
  gap: .75rem;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  margin-bottom: 1.5rem;
}
h3 { margin: 1rem 0 .5rem; font-size: 1rem; }
</style>