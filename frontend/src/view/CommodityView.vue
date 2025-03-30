<template>
  <div v-if="commodity && selectedItem" class="flex gap-5">
    <img :src="selectedItem.picture_url" alt="該規格未上傳圖片" />
    <div class="flex flex-col gap-10">
      <span class="font-bold text-4xl">
        {{ commodity.commodity_name }}
      </span>
      <span class="font-bold text-3xl bg-gray-400">
        $ {{ selectedItem.price }}
      </span>

      <div>
        <div class="flex" v-for="(specType, index) in commodity.specification_types" :key="specType.spec_type_id">
          <span>{{ specType.spec_type_name }}:</span>
          <div v-for="value in specType.specification_values" :key="value.spec_value_id">
            <label>
              <input type="radio" :name="specType.spec_type_name" class="radio" :value="value.spec_value_id"
                v-model="selectedSpecValues[index]" @change="updateSelectedItem" />
              {{ value.spec_value }}
            </label>
          </div>
        </div>

        <div class="join">
          <button @click="num--" class="join-item btn" :disabled="num <= 1">-1</button>
          <input type="number" v-model="num" class="join-item text-xl text-center" min="1" :max="selectedItem.stock" />
          <button @click="num++" class="join-item btn" :disabled="num >= selectedItem.stock">+1</button>
          <span>庫存: {{ selectedItem.stock }}</span>
        </div>
      </div>
      <button @click="buyHandler" class="btn btn-primary">Buy!</button>
    </div>
  </div>

  <div v-else>
    <Error></Error>
  </div>
</template>

<script setup lang="ts">
import type { CommodityDetailResponse, CommoditySpecResponse } from '@/models/Commodity';
import { useCommodityStore } from '@/stores/commodityStore';
import { onMounted, ref } from 'vue';
import Error from '@/components/Error.vue';
import { useRoute } from 'vue-router';
import { requests } from '@/axios';
import { useUserStore } from '@/stores/userStore';

const route = useRoute()
const userStore = useUserStore()
// 商品數據
const commodityStore = useCommodityStore();
const commodity = ref<CommodityDetailResponse>();
// 已選擇的規格值（索引對應 specification_types）
const selectedSpecValues = ref<number[]>([]);
// 已選擇的商品規格項
const selectedItem = ref<CommoditySpecResponse | null>(null);
// 購買數量
const num = ref(1);

onMounted(async () => {
  await getCommodity()
});

const getCommodity = async () => {
  const id = Number(route.params.id);
  if (isNaN(id)) {
    console.error('Invalid ID');
    return;
  }
  commodity.value = await commodityStore.getCommodityById(id);
  if (commodity.value) {
    // 初始化選擇的規格值（預設選第一個）
    selectedSpecValues.value = commodity.value.specification_types.map(
      (specType) => specType.specification_values[0]?.spec_value_id ?? 0
    );
    updateSelectedItem(); // 初始化 selectedItem
  }
}

// 更新已選擇的規格項
function updateSelectedItem() {
  if (!commodity.value) return;

  // 根據選擇的規格值找到匹配的 commodity_specifications 項
  const selected = commodity.value.commodity_specifications.find((spec) => {
    const specValueIds = [spec.spec_value_1_id, spec.spec_value_2_id ?? 0];
    return selectedSpecValues.value.every((val, idx) => specValueIds[idx] === val);
  });

  selectedItem.value = selected || commodity.value.commodity_specifications[0] || null;
  // 重置數量，避免超過庫存
  if (selectedItem.value && num.value > selectedItem.value.stock) {
    num.value = selectedItem.value.stock;
  }
}

// 購買處理
const buyHandler = async () => {
  if (!selectedItem.value || num.value <= 0) {
    alert('請選擇商品規格並設置有效數量！');
    return;
  }
  const id = Number(route.params.id);
  await requests.post("/commodities/buy", {
    "commodity_id": id,
    "spec_type_id": selectedItem.value.commodity_spec_id,
    "user_id": userStore.user?.userId,
    "num": num.value
  })
  alert("購買成功")
  await getCommodity()

}
</script>