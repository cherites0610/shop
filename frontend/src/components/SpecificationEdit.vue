<template>
    <div class="flex flex-col">
        <span>{{ spec_name }}</span>
        <div class="flex">
            <input class="input" v-for="(_, index) in spec_value" v-model="spec_value[index]">
        </div>
    </div>
</template>

<script setup lang="ts">
import { reactive, defineProps } from 'vue';
const { specification } = defineProps<{ specification: Record<string, any> }>();
// 使用 reactive 建立響應式物件
const local_specification = reactive({ ...specification });

// 從響應式物件中提取屬性
const spec_name = Object.keys(local_specification)[0];
const spec_value = local_specification[spec_name];

</script>

<style scoped></style>