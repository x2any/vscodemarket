<script setup lang="ts">
import { ref } from 'vue'
import { t } from '../i18n'
import type { Channel } from '../api/contracts'

const emit = defineEmits<{
  (e: 'submit', v: { channel: Channel; version: string }): void
}>()

const channel = ref<Channel>('stable')
const version = ref('')

function onSubmit() {
  if (!version.value.trim()) return
  emit('submit', { channel: channel.value, version: version.value.trim() })
}
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

    <el-button type="primary" native-type="submit" :disabled="!version.trim()">
      {{ t.home.submit }}
    </el-button>
  </el-form>
</template>