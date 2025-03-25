<template>
    <div class="grid grid-cols-1 gap-4">
        <div class="grid grid-cols-[1fr_1fr_1fr_1fr] gap-2 items-center justify-items-cente">
            <span class=" text-center" v-for="name in spec_name">{{ name }}</span>
        </div>
        <div class="grid grid-cols-[1fr_1fr_1fr_1fr] gap-2 items-center justify-items-center"
            v-for="combo in combinations" :key="combo.join('-')">
            <span v-for="value in combo" :key="value" class="text-center p-2">{{ value }}</span>
            <input type="number" class="text-center p-2 border rounded w-20" v-model="getInf(getIndex(combo)).price">
            <input type="number" class="text-center p-2 border rounded w-20" v-model="getInf(getIndex(combo)).stock">
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue';

const { specification, specificationInf } = defineProps<{ specification: Record<string, any>, specificationInf: any }>();

const spec_name = ref(Object.keys(specification));
spec_name.value.push("價格")
spec_name.value.push("庫存")
const combinations = computed(() => {
    const values = Object.values(specification);
    return getCombinations(values);
})

const getCombinations = (arrays: string[][]): string[][] => {
    if (!arrays.length) return [[]];
    return arrays.reduce((acc, curr) => {
        const result: string[][] = [];
        acc.forEach(a => {
            curr.forEach(b => {
                result.push([...a, b]);
            });
        });
        return result;
    }, [[]] as string[][]);
};

const getIndex = (args: string[]): number => {
    if (!specificationInf) return -1;
    const matchedItem = specificationInf.find((item: any) =>
        args.every((value) => item.spec_value.includes(value)) &&
        item.spec_value.length === args.length // 確保長度相等
    );
    return matchedItem ? matchedItem.commodity_spec_id : -1;
};

const getInf = (commodity_spec_id: number) => {
    return (
        specificationInf?.find((item: { commodity_spec_id: number; }) => item.commodity_spec_id === commodity_spec_id) ||
        { price: 'N/A', stock: 'N/A' }
    );
};

let oldVal = JSON.parse(JSON.stringify(specification));
watch(
    () => specification,
    (newVal) => {
        if (!specificationInf?.length || !newVal) return; // 防護
        const combinations = getCombinations(Object.values(oldVal || {}));
        const newValCombo = getCombinations(Object.values(newVal));

        combinations.forEach((combo, i) => {
            const index = getIndex(combo);
            const matchedItem = getInf(index)
            if (index > 0 && i < newValCombo.length && matchedItem) {
                matchedItem.spec_value = newValCombo[i];

            }
        });
        oldVal = JSON.parse(JSON.stringify(newVal));
    },
    { deep: true }
);
</script>

<style scoped></style>