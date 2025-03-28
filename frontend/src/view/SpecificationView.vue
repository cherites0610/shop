<template>
    <!-- 規格選擇區 -->
    <div v-if="commodity">
        <div v-for="(value, key, index) in commodity.spec" :key="key">
            <SpecificationEdit @delete="(spec_name) => deleteHandler(spec_name)" @updateSpecName="updateSpecNameHandler" :specification="{ [key]: value }" :isFirst="index === 0">
            </SpecificationEdit>
        </div>
        <button v-if="Object.keys(commodity.spec).length < 2" @click="addHandler" class="btn">新增</button>

        <SpecificationInformationEdit :specification="commodity.spec" :specificationInf="commodity.specifications">
        </SpecificationInformationEdit>
    </div>
    <div v-else>
        <Error></Error>
    </div>

    <pre>{{ commodity }}</pre>
</template>

<script setup lang="ts">
import Error from '@/components/Error.vue';
import type Commodity from '@/modals/Commodity';
import { useCommodityStore } from '@/stores/commodityStore';
import { ref, onMounted } from 'vue';
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

const addHandler = () => {
    if (commodity.value) {
        const newKey = `規格${Object.keys(commodity.value.spec).length + 1}`;

        // 重新建立新物件
        commodity.value.spec = {
            ...commodity.value.spec,
            [newKey]: []
        };
    }
};

const deleteHandler = (name: string) => {
    if (commodity.value) {
        // 重新建立新物件，確保 Vue 偵測變化
        const newSpec = { ...commodity.value.spec };
        delete newSpec[name];
        commodity.value.spec = newSpec;

        // 直接刪除符合條件的 specifications 項目
        commodity.value.specifications = commodity.value.specifications.filter(
            (item) => !item.spec_value.includes(name)
        );
    }
};

const updateSpecNameHandler = (oldName: string, newName: string) => {
    if (commodity.value) {
        const newSpec = { ...commodity.value.spec };
        newSpec[newName] = newSpec[oldName];
        delete newSpec[oldName];
        commodity.value.spec = newSpec;

        // 更新 specifications 中的 spec_value
        commodity.value.specifications.forEach((item) => {
            const index = item.spec_value.indexOf(oldName);
            if (index !== -1) {
                item.spec_value[index] = newName;
            }
        });

        // 確保 Vue 追蹤變更
        commodity.value = { ...commodity.value };
    }
};

// 初始化載入
onMounted(async () => {
    await getCommodity();
});
</script>

<style scoped></style>