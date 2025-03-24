<template>
    <div class="flex flex-wrap gap-5 m-5 justify-center">
        <CommodityCard v-for="item in commodities" :commodity="item">
        </CommodityCard>
    </div>
</template>

<script setup lang="ts">
import CommodityCard from '@/components/CommodityCard.vue';
import Commodity from '@/modals/Commodity';

const commodities: Commodity[] = [];

for (let i = 0; i < 10; i++) {
    const letter = String.fromCharCode(65 + i); // A-Z
    commodities.push(
        new Commodity(
            String(i+1),
            letter,
            100 + i * 10,
            { color: ["紅色", "藍色"], size: ["S", "M", "L"] },
            10 + i
        )
    );
}


</script>

<style scoped></style>