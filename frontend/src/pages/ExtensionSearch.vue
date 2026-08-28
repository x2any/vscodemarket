<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { http } from '../api/http'
import { t } from '../i18n'

interface SearchResult {
  publisher: string
  name: string
  displayName: string
  latestVersion: string
}

const router = useRouter()
const q = ref('')
const results = ref<SearchResult[]>([])
const total = ref(0)
const error = ref('')

async function submit() {
  if (!q.value.trim()) return
  error.value = ''
  results.value = []
  total.value = 0
  try {
    const resp = await http.get<{ results: SearchResult[]; total: number }>(
      `/extensions/search?q=${encodeURIComponent(q.value.trim())}`
    )
    results.value = resp.results
    total.value = resp.total
    // Single match → jump straight to detail page (per spec US4 acceptance).
    if (resp.results.length === 1) {
      const only = resp.results[0]
      router.push({ name: 'extension-detail', params: { pub: only.publisher, name: only.name } })
    }
  } catch (e) {
    error.value = e instanceof Error ? e.message : String(e)
  }
}

function openDetail(r: SearchResult) {
  router.push({ name: 'extension-detail', params: { pub: r.publisher, name: r.name } })
}
</script>

<template>
  <section>
    <el-form @submit.prevent="submit" label-position="top">
      <el-form-item>
        <el-input
          v-model="q"
          :placeholder="t.extension.searchPlaceholder"
          clearable
        />
      </el-form-item>
      <el-button type="primary" native-type="submit" :disabled="!q.trim()">
        {{ t.nav.search }}
      </el-button>
    </el-form>

    <p v-if="error" class="err">{{ error }}</p>

    <p v-if="!error && q && results.length === 0">{{ t.extension.noResults }}</p>

    <ul v-if="results.length > 1" class="list">
      <li v-for="r in results" :key="r.publisher + '.' + r.name" class="row">
        <a href="#" @click.prevent="openDetail(r)">
          <strong>{{ r.displayName }}</strong>
          <span class="id">{{ r.publisher }}.{{ r.name }}</span>
          <span class="ver">v{{ r.latestVersion }}</span>
        </a>
      </li>
    </ul>

    <p v-if="total" class="total">{{ total }} {{ t.nav.search }}</p>
  </section>
</template>

<style scoped>
.err { color: var(--el-color-danger); }
.list { list-style: none; padding: 0; margin-top: 1rem; }
.row { padding: .5rem 0; border-bottom: 1px solid #e5e7eb; }
.row a { display: flex; gap: 1rem; align-items: baseline; text-decoration: none; color: inherit; }
.id { opacity: .6; font-family: ui-monospace, monospace; }
.ver { margin-left: auto; opacity: .75; }
.total { margin-top: 1rem; opacity: .6; }
</style>