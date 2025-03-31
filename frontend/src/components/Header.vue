<template>
    <div class="navbar bg-base-100 shadow-sm">
        <div class="flex-1">
            <a href="/" class="btn btn-ghost text-xl">bots</a>
        </div>
        <div class="flex-none">
            <ul class="menu menu-horizontal px-1">
                <li>
                    <RouterLink to="/">商品列表</RouterLink>
                </li>
                <li>
                    <details>
                        <summary>{{ userStore.user?.displayName }}</summary>
                        <ul class="bg-base-100 rounded-t-none p-2">
                            <li><button @click="logoutHandler">登出</button></li>
                        </ul>
                    </details>
                </li>
            </ul>
        </div>
    </div>
</template>

<script setup lang="ts">
import { useUserStore } from '@/stores/userStore';

const userStore = useUserStore()

const logoutHandler = () => {
    localStorage.removeItem("lineToken")
    window.location.reload()
}
</script>

<style scoped></style>
