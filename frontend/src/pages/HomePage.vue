<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { track } from '../api/track'
import { t } from '../i18n'
import VersionForm from '../components/VersionForm.vue'
import type { Channel } from '../api/contracts'

const formRef = ref<InstanceType<typeof VersionForm> | null>(null)
const router = useRouter()

onMounted(() => {
  formRef.value?.pingInfer()
})

function onInfer(_payload: { userAgent: string }) {
  // Best-effort prefill; nothing to do here — the form handles inference.
}

function onSubmit(payload: { channel: Channel; version: string }) {
  track({
    eventType: 'SEARCH',
    targetType: 'CLIENT',
    targetIdentifier: payload.version,
    channel: payload.channel
  })
  router.push({
    name: 'version-detail',
    params: { channel: payload.channel, ver: payload.version }
  })
}
</script>

<template>
  <section>
    <h2>{{ t.home.submit }}</h2>
    <VersionForm ref="formRef" @submit="onSubmit" @infer="onInfer" />
  </section>
</template>