<template>
    <!-- 規格選擇區 -->
    <div v-if="commodity">
        <div v-for="(value, key) in commodity.spec" :key="key">
            <SpecificationEdit :specification="{ [key]: value }"></SpecificationEdit>
        </div>


        <SpecificationInformationEdit :specification="commodity.spec" :specificationInf="commodity.specifications">
        </SpecificationInformationEdit>
    </div>
    <div v-else>
        <Error></Error>
    </div>
</template>

<script setup lang="ts">
import Error from '@/components/Error.vue';
import type Commodity from '@/modals/Commodity';
import { useCommodityStore } from '@/stores/commodityStore';
import { ref, onMounted, watch } from 'vue';
import { useRoute } from 'vue-router';
import SpecificationEdit from '@/components/SpecificationEdit.vue';
import SpecificationInformationEdit from '@/components/SpecificationInformationEdit.vue';

const commodityStore = useCommodityStore();
const commodity = ref<Commodity>();
const route = useRoute();

// 獲取商品數據
const getCommodity = async () => {
    try {
        commodity.value = await commodityStore.getCommodityById(Number(route.params?.id));
    } catch (error) {
        console.error('Failed to load commodity:', error);
    }
};

// 初始化載入
onMounted(async () => {
    await getCommodity();
});

</script>

<style scoped></style>