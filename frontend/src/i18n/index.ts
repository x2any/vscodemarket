import { ref, computed } from 'vue'
import { zhCN, type Dict } from './zh-CN'
import { enUS } from './en-US'

export type Locale = 'zh-CN' | 'en-US'

const dicts: Record<Locale, Dict> = { 'zh-CN': zhCN, 'en-US': enUS }

const stored = (typeof localStorage !== 'undefined' && localStorage.getItem('locale')) as Locale | null
const locale = ref<Locale>(stored && stored in dicts ? stored : 'zh-CN')

export const currentLocale = computed(() => locale.value)
export const t = computed<Dict>(() => dicts[locale.value])

export function setLocale(next: Locale) {
  locale.value = next
  if (typeof localStorage !== 'undefined') localStorage.setItem('locale', next)
}
