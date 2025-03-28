<template>
    <div class="flex flex-col">
        <!-- 讓 spec_name 變成可編輯的 input -->
        <input v-model="editableSpecName" @blur="updateSpecName" class="input">

        <div class="flex">
            <input class="input" v-for="(_, index) in 2" v-model="spec_value[index]">
        </div>
        <button @click="emit('delete', spec_name)" class="btn">刪除</button>
    </div>
</template>

<script setup lang="ts">
import { reactive, ref, watchEffect, defineProps, defineEmits } from 'vue';

const props = defineProps<{ specification: Record<string, any>, isFirst: boolean }>();
const emit = defineEmits(['delete', 'updateSpecName']);

// 建立響應式的 local_specification，確保能夠雙向綁定
const local_specification = reactive({ ...props.specification });

// 取得 spec_name 和 spec_value
const spec_name = ref(Object.keys(local_specification)[0] || "");
const spec_value = ref(local_specification[spec_name.value] || []);

// 用來編輯的 spec_name
const editableSpecName = ref(spec_name.value);

// 更新規格名稱
const updateSpecName = () => {
    if (editableSpecName.value !== spec_name.value) {
        const oldKey = spec_name.value;
        const value = local_specification[oldKey];

        // 先刪除舊的 key
        delete local_specification[oldKey];

        // 設定新的 key
        local_specification[editableSpecName.value] = value;

        // 更新 spec_name
        spec_name.value = editableSpecName.value;

        // 讓 Vue 重新追蹤 spec_value
        spec_value.value = local_specification[editableSpecName.value];

        // 通知父組件規格名稱的變更
        emit('updateSpecName', oldKey, editableSpecName.value);
    }
};

// 監聽 `specification`，確保 UI 同步變更
watchEffect(() => {
    const key = Object.keys(props.specification)[0] || "";
    spec_name.value = key;
    editableSpecName.value = key;
    spec_value.value = props.specification[key] || [];
});
</script>