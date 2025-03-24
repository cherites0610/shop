<template>
    <div v-if="loading">
        <span class="loading loading-spinner loading-xl"></span>
    </div>
    <div v-else>
        <Error></Error>
    </div>
</template>

<script setup lang="ts">
import { requests } from '@/axios';
import Error from '@/components/Error.vue';
import { ref } from 'vue';

const loading = ref<boolean>(true);

const urlParams = new URLSearchParams(window.location.search);
const code = urlParams.get('code');
const state = urlParams.get('state');

if (code) {
    exchangeLineToken(code);
} else {
    loading.value=false
    console.error('Code parameter is missing');
}

async function exchangeLineToken(code: string) {
    const result = await requests.post("/LineAcess", {
        "code": code,
        "state": state
    }, {
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' }
    })

    let line = JSON.parse(result.data.response)
    if (line.access_token) {
        window.localStorage.setItem("lineToken", line.access_token)
        window.location.href="/"
    }
}
</script>

<style scoped></style>