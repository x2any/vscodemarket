<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { http } from '../api/http'
import { t } from '../i18n'

type TargetType = 'CLIENT' | 'SERVER' | 'EXTENSION'
type Window = '24h' | '7d' | '30d'

interface TrendingRow { targetIdentifier: string; count: number }
interface TrendingResp { targetType: TargetType; window: Window; items: TrendingRow[] }

const groups: { key: TargetType; labelKey: keyof typeof t.trending }[] = [
  { key: 'CLIENT', labelKey: 'client' },
  { key: 'SERVER', labelKey: 'server' },
  { key: 'EXTENSION', labelKey: 'extension' }
]

const windows: Window[] = ['24h', '7d', '30d']
const currentWindow = ref<Window>('24h')

const data = ref<Record<TargetType, TrendingRow[]>>({
  CLIENT: [], SERVER: [], EXTENSION: []
})

async function load() {
  await Promise.all(groups.map(async g => {
    try {
      const r = await http.get<TrendingResp>(
        `/trending?targetType=${g.key}&window=${currentWindow.value}`
      )
      data.value[g.key] = r.items
    } catch {
      data.value[g.key] = []
    }
  }))
}

onMounted(load)
watch(currentWindow, load)

function labelFor(key: keyof typeof t.trending): string {
  return t.trending[key]
}
function windowLabel(w: Window): string {
  if (w === '24h') return t.trending.window24h
  if (w === '7d') return t.trending.window7d
  return t.trending.window30d
}
</script>

<template>
  <section>
    <header class="bar">
      <h2>{{ t.nav.trending }}</h2>
      <el-radio-group v-model="currentWindow">
        <el-radio-button v-for="w in windows" :key="w" :label="w">
          {{ windowLabel(w) }}
        </el-radio-button>
      </el-radio-group>
    </header>

    <div class="grid">
      <article v-for="g in groups" :key="g.key" class="panel">
        <h3>{{ labelFor(g.labelKey) }}</h3>
        <ol v-if="data[g.key].length">
          <li v-for="(r, i) in data[g.key]" :key="r.targetIdentifier">
            <span class="rank">#{{ i + 1 }}</span>
            <code class="id">{{ r.targetIdentifier }}</code>
            <span class="cnt">{{ r.count }}</span>
          </li>
        </ol>
        <p v-else class="empty">{{ t.trending.empty }}</p>
      </article>
    </div>
  </section>
</template>

<style scoped>
.bar { display: flex; gap: 1rem; align-items: center; flex-wrap: wrap; margin-bottom: 1rem; }
.grid { display: grid; gap: 1rem; grid-template-columns: repeat(auto-fit, minmax(280px, 1fr)); }
.panel { padding: 1rem; border: 1px solid #e5e7eb; border-radius: 6px; }
.panel h3 { margin: 0 0 .5rem; }
ol { padding-left: 1rem; margin: 0; }
li { display: flex; gap: .5rem; align-items: baseline; padding: .25rem 0; }
.rank { opacity: .5; min-width: 2rem; }
.id { flex: 1; font-family: ui-monospace, monospace; }
.cnt { opacity: .75; }
.empty { opacity: .6; }
</style>