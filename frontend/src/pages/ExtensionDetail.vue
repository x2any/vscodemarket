<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { http } from '../api/http'
import { t } from '../i18n'

interface Version {
  version: string
  publishTime: string
  enginesVscode: string
  downloadUrl: string
}

interface VersionsResp {
  extension: { publisher: string; name: string; displayName?: string }
  versions: Version[]
}

const route = useRoute()
const router = useRouter()

const publisher = computed(() => String(route.params.pub ?? ''))
const name = computed(() => String(route.params.name ?? ''))
const specificVersion = computed(() => route.params.ver ? String(route.params.ver) : '')

const versions = ref<Version[]>([])
const error = ref('')
const loading = ref(false)

async function load() {
  if (!publisher.value || !name.value) return
  loading.value = true
  error.value = ''
  versions.value = []
  try {
    const resp = await http.get<VersionsResp>(
      `/extensions/${publisher.value}/${name.value}/versions`
    )
    versions.value = resp.versions
  } catch (e) {
    error.value = e instanceof Error ? e.message : String(e)
  } finally {
    loading.value = false
  }
}

function goToVersion(v: Version) {
  router.push({
    name: 'extension-version',
    params: { pub: publisher.value, name: name.value, ver: v.version }
  })
}

onMounted(load)
</script>

<template>
  <section>
    <h2>{{ publisher }}.{{ name }}</h2>

    <p v-if="specificVersion" class="hint">
      Filtering: <code>{{ specificVersion }}</code>
    </p>

    <p v-if="loading">Loading…</p>
    <p v-if="error" class="err">{{ error }}</p>

    <table v-if="versions.length" class="grid">
      <thead>
        <tr>
          <th>Version</th>
          <th>{{ t.extension.engines }}</th>
          <th>Published</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="v in versions" :key="v.version">
          <td><code>{{ v.version }}</code></td>
          <td><el-tag>{{ v.enginesVscode || '—' }}</el-tag></td>
          <td>{{ new Date(v.publishTime).toISOString().slice(0, 10) }}</td>
          <td>
            <el-button v-if="specificVersion === v.version" type="primary" tag="a" :href="v.downloadUrl" target="_blank" rel="noopener">
              Download
            </el-button>
            <el-button v-else link @click="goToVersion(v)">View</el-button>
          </td>
        </tr>
      </tbody>
    </table>
  </section>
</template>

<style scoped>
.err { color: var(--el-color-danger); }
.grid { width: 100%; border-collapse: collapse; margin-top: 1rem; }
.grid th, .grid td { padding: .5rem; text-align: left; border-bottom: 1px solid #e5e7eb; }
.hint { opacity: .75; }
</style>