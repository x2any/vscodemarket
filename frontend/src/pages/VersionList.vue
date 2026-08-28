<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { http } from '../api/http'
import { t } from '../i18n'

interface ReleaseRow {
  channel: 'stable' | 'insider'
  version: string
  downloadUrl: string
}

interface ReleasesResp {
  results: ReleaseRow[]
  page: number
  pageSize: number
  total: number
}

const router = useRouter()
const channel = ref<'stable' | 'insider'>('stable')
const page = ref(1)
const pageSize = 30

const results = ref<ReleaseRow[]>([])
const total = ref(0)
const error = ref('')

async function load() {
  error.value = ''
  try {
    const resp = await http.get<ReleasesResp>(
      `/releases?channel=${channel.value}&page=${page.value}&pageSize=${pageSize}`
    )
    results.value = resp.results
    total.value = resp.total
  } catch (e) {
    error.value = e instanceof Error ? e.message : String(e)
  }
}

onMounted(load)
watch([channel, page], load)

function open(r: ReleaseRow) {
  router.push({
    name: 'version-detail',
    params: { channel: r.channel, ver: r.version }
  })
}

function nextPage() { page.value += 1 }
function prevPage() { if (page.value > 1) page.value -= 1 }
</script>

<template>
  <section>
    <h2>{{ t.nav.releases }}</h2>

    <div class="filters">
      <el-radio-group v-model="channel" @change="load">
        <el-radio-button label="stable">stable</el-radio-button>
        <el-radio-button label="insider">insider</el-radio-button>
      </el-radio-group>
    </div>

    <p v-if="error" class="err">{{ error }}</p>

    <ul v-if="results.length" class="list">
      <li v-for="r in results" :key="r.channel + r.version">
        <a href="#" @click.prevent="open(r)">
          <code class="ver">{{ r.version }}</code>
          <small class="channel">{{ r.channel }}</small>
        </a>
      </li>
    </ul>

    <div class="pager">
      <el-button :disabled="page <= 1" @click="prevPage">‹</el-button>
      <span>page {{ page }} / {{ Math.max(1, Math.ceil(total / pageSize)) }}</span>
      <el-button :disabled="page * pageSize >= total" @click="nextPage">›</el-button>
    </div>
  </section>
</template>

<style scoped>
.filters { margin-bottom: 1rem; }
.err { color: var(--el-color-danger); }
.list { list-style: none; padding: 0; margin: 0; }
.list li { padding: .5rem 0; border-bottom: 1px solid #e5e7eb; }
.list a { display: flex; gap: 1rem; align-items: baseline; text-decoration: none; color: inherit; }
.ver { font-family: ui-monospace, monospace; font-size: .95rem; }
.channel { opacity: .6; text-transform: uppercase; font-size: .75rem; letter-spacing: .05em; }
.pager { display: flex; gap: 1rem; align-items: center; margin-top: 1rem; }
</style>