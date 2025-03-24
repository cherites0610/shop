<template>
    <div class="flex flex-wrap gap-5 m-5 justify-center">
        <CommodityCard v-for="item in commodityStore.commodities" :key="item.id" :commodity="item">
        </CommodityCard>
    </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import CommodityCard from '@/components/CommodityCard.vue';
import { useCommodityStore } from '@/stores/commodityStore';

const commodityStore = useCommodityStore();

onMounted( async () => {
    await commodityStore.getCommodities()
});

const commodities = commodityStore.commodities;
</script>

<style scoped>

</style>
