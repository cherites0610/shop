<template>
  <div>
    <Header></Header>
    <RouterView></RouterView>
  </div>
</template>

<script setup lang="ts">
import { onBeforeMount } from 'vue';
import Header from './components/Header.vue';
import { useUserStore } from './stores/userStore';
import { requests } from './axios';
import { useCommodityStore } from './stores/commodityStore';

const userStore = useUserStore();
const commodityStore = useCommodityStore();

onBeforeMount(async () => {
  // const token = window.localStorage.getItem("lineToken");
  await commodityStore.getCommodities()

  // const getLineUrl = async (): Promise<string> => {
  //   const lineUrl = await requests.get("/lineLogin")
  //   return lineUrl.data.url
  // }
  
  // try {
  //   if (!token) throw new Error("Token not found");
  //   await userStore.getUserProfile(token);
  // } catch (err: any) {
  //   window.location.href = await getLineUrl();
  // }



})


</script>

<style scoped></style>