<script setup lang="ts">
import { useRouter } from 'vue-router'
import { track } from '../api/track'
import { t } from '../i18n'
import VersionForm from '../components/VersionForm.vue'
import type { Channel } from '../api/contracts'

const router = useRouter()

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
    <VersionForm @submit="onSubmit" />
  </section>
</template>