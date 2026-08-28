<script setup lang="ts">
import { computed, ref } from 'vue'
import { http } from '../api/http'
import { t } from '../i18n'

interface ReleaseEntry {
  channel: 'stable' | 'insider'
  version: string
  platform?: string
  architecture?: string
  downloadUrl: string
  commitHash?: string
  releaseDate?: string
}

interface ReleasesResp {
  results: ReleaseEntry[]
  page: number
  pageSize: number
  total: number
}

const channel = ref<'stable' | 'insider'>('stable')
const platform = ref<string>('')
const architecture = ref<string>('')
const page = ref(1)
const pageSize = 20

const results = ref<ReleaseEntry[]>([])
const total = ref(0)
const error = ref('')

const query = computed(() => {
  const parts = [`channel=${channel.value}`, `page=${page.value}`, `pageSize=${pageSize}`]
  if (platform.value) parts.push(`platform=${platform.value}`)
  if (architecture.value) parts.push(`architecture=${architecture.value}`)
  return parts.join('&')
})

async function load() {
  error.value = ''
  try {
    const resp = await http.get<ReleasesResp>(`/releases?${query.value}`)
    results.value = resp.results
    total.value = resp.total
  } catch (e) {
    error.value = e instanceof Error ? e.message : String(e)
  }
}

function nextPage() {
  page.value += 1
  load()
}
function prevPage() {
  if (page.value > 1) {
    page.value -= 1
    load()
  }
}

load()
</script>

<template>
  <section>
    <h2>{{ t.nav.releases }}</h2>

    <div class="filters">
      <el-radio-group v-model="channel" @change="load">
        <el-radio-button label="stable">stable</el-radio-button>
        <el-radio-button label="insider">insider</el-radio-button>
      </el-radio-group>

      <el-select v-model="platform" placeholder="Platform" clearable @change="load" style="width: 9rem">
        <el-option label="windows" value="windows" />
        <el-option label="linux" value="linux" />
        <el-option label="darwin" value="darwin" />
      </el-select>

      <el-select v-model="architecture" placeholder="Architecture" clearable @change="load" style="width: 9rem">
        <el-option label="x86_64" value="x86_64" />
        <el-option label="arm64" value="arm64" />
        <el-option label="armv7" value="armv7" />
      </el-select>
    </div>

    <p v-if="error" class="err">{{ error }}</p>

    <table v-if="results.length" class="grid">
      <thead>
        <tr>
          <th>Version</th>
          <th>{{ channel === 'stable' ? 'Platform/Arch' : '' }}</th>
          <th>Commit</th>
          <th>Date</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="r in results" :key="r.version + (r.platform || '')">
          <td><code>{{ r.version }}</code></td>
          <td>{{ r.platform ? `${r.platform}/${r.architecture}` : '—' }}</td>
          <td><code v-if="r.commitHash">{{ r.commitHash.slice(0, 8) }}</code><span v-else>—</span></td>
          <td>{{ r.releaseDate || '—' }}</td>
          <td><el-link :href="r.downloadUrl" target="_blank" rel="noopener" type="primary">link</el-link></td>
        </tr>
      </tbody>
    </table>

    <div class="pager">
      <el-button :disabled="page <= 1" @click="prevPage">‹</el-button>
      <span>page {{ page }} / {{ Math.max(1, Math.ceil(total / pageSize)) }}</span>
      <el-button :disabled="page * pageSize >= total" @click="nextPage">›</el-button>
    </div>
  </section>
</template>

<style scoped>
.filters { display: flex; gap: 1rem; flex-wrap: wrap; margin-bottom: 1rem; }
.err { color: var(--el-color-danger); }
.grid { width: 100%; border-collapse: collapse; }
.grid th, .grid td { padding: .5rem; text-align: left; border-bottom: 1px solid #e5e7eb; }
.pager { display: flex; gap: 1rem; align-items: center; margin-top: 1rem; }
</style>