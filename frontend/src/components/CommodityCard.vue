<template>
    <div class="card bg-base-100 w-96 shadow-sm">
        <figure>
            <img src="https://random.imagecdn.app/350/150" :alt=commodity.name />
        </figure>
        <div class="card-body">
            <div class="flex flex-wrap navbar-center justify-between">
                <h2 class=" text-3xl font-bold">{{ commodity.name }}</h2>
                <button @click="router.push(`/commodity/${commodity.id}`)" class="btn btn-primary">
                    Check
                </button>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import type Commodity from '@/modals/Commodity';
import router from '@/router';
const { commodity } = defineProps<{ commodity: Commodity }>();
</script>

<style scoped></style>