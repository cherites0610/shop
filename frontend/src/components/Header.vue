<template>
    <div class="navbar bg-base-100 shadow-sm">
        <div class="flex-1">
            <a href="bots.ts" class="btn btn-ghost text-xl">bots</a>
        </div>
        <div class="flex-none">
            <ul class="menu menu-horizontal px-1">
                <li>
                    <RouterLink to="/">商品列表</RouterLink>
                </li>
                <li>
                    <RouterLink to="/">規格管理</RouterLink>
                </li>
                <li>
                    <RouterLink to="/">{{userStore.user?.displayName}}-您好</RouterLink>
                </li>
            </ul>
        </div>
    </div>
</template>

<script setup lang="ts">
import { useUserStore } from '@/stores/userStore';

const userStore = useUserStore()

</script>

<style scoped></style>