<template>
    <div class="flex flex-wrap gap-5 m-5 justify-center">
        <CommodityCard v-for="item in commodityStore.commodities" :key="item.commodity_id" :commodity="item">
        </CommodityCard>
    </div>
</template>

<script setup lang="ts">
import CommodityCard from '@/components/CommodityCard.vue';
import { useCommodityStore } from '@/stores/commodityStore';

const commodityStore = useCommodityStore();

</script>

<style scoped>

</style>
