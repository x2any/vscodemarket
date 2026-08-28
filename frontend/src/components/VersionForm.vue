<script setup lang="ts">
import { ref } from 'vue'
import { t } from '../i18n'
import type { Channel, Platform, Architecture } from '../api/contracts'

const emit = defineEmits<{
  (e: 'submit', v: { channel: Channel; version: string; platform?: Platform; architecture?: Architecture }): void
  (e: 'infer', v: { userAgent: string }): void
}>()

const channel = ref<Channel>('stable')
const version = ref('')
// platform / architecture are optional. When both blank the response returns
// the full platform × architecture matrix.
const platform = ref<Platform | ''>('')
const architecture = ref<Architecture | ''>('')

function onSubmit() {
  if (!version.value.trim()) return
  const payload: { channel: Channel; version: string; platform?: Platform; architecture?: Architecture } = {
    channel: channel.value,
    version: version.value.trim()
  }
  if (platform.value && architecture.value) {
    payload.platform = platform.value
    payload.architecture = architecture.value
  }
  emit('submit', payload)
}

defineExpose({
  fillInferred(p: Platform, a: Architecture) {
    platform.value = p
    architecture.value = a
  },
  pingInfer() {
    if (typeof navigator !== 'undefined') {
      emit('infer', { userAgent: navigator.userAgent })
    }
  }
})
</script>

<template>
  <el-form @submit.prevent="onSubmit" label-position="top">
    <el-form-item :label="t.home.channel">
      <el-radio-group v-model="channel">
        <el-radio-button label="stable">stable</el-radio-button>
        <el-radio-button label="insider">insider</el-radio-button>
      </el-radio-group>
    </el-form-item>

    <el-form-item :label="t.home.version">
      <el-input v-model="version" placeholder="1.94.2" />
    </el-form-item>

    <el-form-item :label="t.home.platform">
      <el-select v-model="platform" clearable placeholder="(可选)">
        <el-option label="windows" value="windows" />
        <el-option label="linux" value="linux" />
        <el-option label="darwin" value="darwin" />
      </el-select>
    </el-form-item>

    <el-form-item :label="t.home.architecture">
      <el-select v-model="architecture" clearable placeholder="(可选)">
        <el-option label="x86_64" value="x86_64" />
        <el-option label="arm64" value="arm64" />
        <el-option label="armv7" value="armv7" />
      </el-select>
    </el-form-item>

    <el-button type="primary" native-type="submit" :disabled="!version.trim()">
      {{ t.home.submit }}
    </el-button>
    <p class="hint">{{ t.home.matrixHint }}</p>
  </el-form>
</template>