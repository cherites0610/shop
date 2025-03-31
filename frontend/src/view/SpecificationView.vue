<template>
  <div v-if="commodity">
    <!-- 商品名稱和規格編輯區 -->
    <div>
      <span>商品名稱: <input class="input" v-model="commodity.commodity_name"></span>
      <div v-for="spec_type of commodity.specification_types" :key="spec_type.spec_type_id">
        <span>規格名字</span><input class="input" v-model="spec_type.spec_type_name">
        <div v-for="spec_value of spec_type.specification_values" :key="spec_value.spec_value_id"
          class="spec-value-row">
          <input class="input" v-model="spec_value.spec_value" @input="updateSpecifications(spec_value)">
          <button class="btn" v-if="spec_type.specification_values.length > 1"
            @click="deleteSpecValue(spec_type.spec_type_id, spec_value.spec_value_id)">刪除規格值</button>
        </div>
        <button class="btn" v-if="spec_type.specification_values.length < 2"
          @click="addSpecValue(spec_type.spec_type_id)">新增規格值</button>
        <button @click="deleteSpecType(spec_type.spec_type_id)" v-if="commodity.specification_types.length > 1"
          class="btn">刪除規格類型</button>
      </div>
      <div v-if="commodity.specification_types.length < 2">
        <input class="input" v-model="newSpecTypeName" placeholder="輸入新規格名稱">
        <button @click="createSpecType" class="btn">新增規格類型</button>
      </div>
    </div>

    <!-- 編輯區：展示規格組合並允許編輯價格和庫存 -->
    <div>
      <span>編輯區</span>
      <div v-for="spec in commodity.commodity_specifications" :key="spec.commodity_spec_id">
        {{ getSpecValueName(spec.spec_value_1_id) }}
        {{ spec.spec_value_2_id ? getSpecValueName(spec.spec_value_2_id) : "" }}:
        價格 <input class="input" type="number" v-model="spec.price">
        庫存 <input class="input" type="number" v-model="spec.stock">
      </div>
    </div>

    <button @click="submit(commodity)" class="btn">送出</button>
  </div>

  <pre>{{ commodity }}</pre>
</template>

<script setup lang="ts">
import { requests } from '@/axios';
import type { CommodityDetailResponse } from '@/models/Commodity';
import { useCommodityStore } from '@/stores/commodityStore';
import { onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';

const commodity = ref<CommodityDetailResponse>();
const commodityStore = useCommodityStore();
const route = useRoute();
const newSpecTypeName = ref(''); // 用於輸入新規格類型名稱
const loading = ref(false);

// 用於生成臨時 ID 的計數器
const tempIdCounter = ref(0);

// 生成唯一的臨時 ID（負數避免與後端 ID 衝突）
const generateTempId = () => {
  tempIdCounter.value -= 1; // 從 -1 開始遞減
  return tempIdCounter.value;
};

onMounted(async () => {
  const id = Number(route.params.id);
  if (isNaN(id)) {
    console.error('Invalid ID');
    return;
  }
  commodity.value = await commodityStore.getCommodityById(id);
});

// 根據 spec_value_id 查找對應的 spec_value
const getSpecValueName = (specValueId: number) => {
  if (commodity.value) {
    for (const specType of commodity.value.specification_types) {
      const specValue = specType.specification_values.find(
        (val) => val.spec_value_id === specValueId
      );
      if (specValue) return specValue.spec_value;
    }
  }
  return '未知';
};

// 生成 commodity_specifications 的默認值，並保留舊有資料
const generateCommoditySpecifications = () => {
  if (!commodity.value || !commodity.value.specification_types) return [];

  const specTypes = commodity.value.specification_types;
  const oldSpecs = commodity.value.commodity_specifications || []; // 舊有的規格組合
  let newCombinations: any[] = [];

  if (specTypes.length === 0) return [];

  // 獲取所有規格值的 ID 和名稱
  const specValues1 = specTypes[0].specification_values;
  const specValues2 = specTypes[1] ? specTypes[1].specification_values : [{ spec_value_id: 0, spec_value: '' }];

  // 生成規格組合（笛卡爾積）
  specValues1.forEach((val1) => {
    specValues2.forEach((val2) => {
      // 查找舊資料中是否已有此組合
      const existingSpec = oldSpecs.find(
        (spec) =>
          spec.spec_value_1_id === val1.spec_value_id &&
          (spec.spec_value_2_id === val2.spec_value_id || (!spec.spec_value_2_id && !val2.spec_value_id))
      );

      if (existingSpec) {
        // 如果舊組合存在，保留其屬性（price, stock 等）
        newCombinations.push({
          ...existingSpec,
          spec_value_1: val1.spec_value, // 更新名稱以保持一致
          spec_value_2: val2.spec_value || null,
        });
      } else {
        // 如果是新組合，初始化價格和庫存
        newCombinations.push({
          spec_value_1_id: val1.spec_value_id,
          spec_value_1: val1.spec_value,
          spec_value_2_id: val2.spec_value_id || null,
          spec_value_2: val2.spec_value || null,
          price: 1,
          stock: 100,
          picture_url: 'https://picsum.photos/200/300',
        });
      }
    });
  });

  return newCombinations;
};

// 更新 commodity_specifications 中的 spec_value_1 和 spec_value_2
const updateSpecifications = (updatedSpecValue: { spec_value_id: number; spec_value: string }) => {
  if (!commodity.value || !commodity.value.commodity_specifications) return;
  commodity.value.commodity_specifications.forEach((spec) => {
    if (spec.spec_value_1_id === updatedSpecValue.spec_value_id) {
      spec.spec_value_1 = updatedSpecValue.spec_value;
    }
    if (spec.spec_value_2_id === updatedSpecValue.spec_value_id) {
      spec.spec_value_2 = updatedSpecValue.spec_value;
    }
  });
};

// 新增規格類型
const createSpecType = () => {
  if (!commodity.value || !newSpecTypeName.value.trim()) {
    alert('請輸入規格名稱！');
    return;
  }

  // 生成臨時 spec_type_id
  const tempSpecTypeId = generateTempId();

  // 新增規格類型，包含一個默認規格值
  commodity.value.specification_types.push({
    spec_type_id: tempSpecTypeId, // 臨時 ID
    spec_type_name: newSpecTypeName.value,
    specification_values: [
      { spec_value_id: generateTempId(), spec_value: '默認值' }, // 臨時 spec_value_id
    ],
  });

  // 清空輸入框
  newSpecTypeName.value = '';

  // 重新生成 commodity_specifications
  commodity.value.commodity_specifications = generateCommoditySpecifications();
};

// 刪除指定的規格類型並重新生成 commodity_specifications
const deleteSpecType = (specTypeId: number) => {
  if (!commodity.value || !commodity.value.specification_types) return;

  const specTypeIndex = commodity.value.specification_types.findIndex(
    (type) => type.spec_type_id === specTypeId
  );
  if (specTypeIndex === -1) return;

  commodity.value.specification_types.splice(specTypeIndex, 1);
  commodity.value.commodity_specifications = generateCommoditySpecifications();
};

// 新增規格值
const addSpecValue = (specTypeId: number) => {
  if (!commodity.value || !commodity.value.specification_types) return;

  const specType = commodity.value.specification_types.find(
    (type) => type.spec_type_id === specTypeId
  );
  if (!specType) return;

  // 添加新規格值並分配臨時 ID
  specType.specification_values.push({
    spec_value_id: generateTempId(), // 臨時 ID
    spec_value: '新值',
  });

  // 重新生成 commodity_specifications
  commodity.value.commodity_specifications = generateCommoditySpecifications();
};

// 刪除規格值
const deleteSpecValue = (specTypeId: number, specValueId: number) => {
  if (!commodity.value || !commodity.value.specification_types) return;

  const specType = commodity.value.specification_types.find(
    (type) => type.spec_type_id === specTypeId
  );
  if (!specType) return;

  const specValueIndex = specType.specification_values.findIndex(
    (val) => val.spec_value_id === specValueId
  );
  if (specValueIndex === -1) return;

  specType.specification_values.splice(specValueIndex, 1);
  commodity.value.commodity_specifications = generateCommoditySpecifications();
};

// 剔除臨時 ID 的輔助函數
const removeTempIds = (data: CommodityDetailResponse): CommodityDetailResponse => {
  const cleanedData = JSON.parse(JSON.stringify(data)); // 深拷貝，避免修改原始資料

  // 處理 specification_types
  cleanedData.specification_types.forEach((specType: any) => {
    if (specType.spec_type_id < 0) {
      delete specType.spec_type_id; // 移除臨時 spec_type_id
    }
    specType.specification_values.forEach((specValue: any) => {
      if (specValue.spec_value_id < 0) {
        delete specValue.spec_value_id; // 移除臨時 spec_value_id
      }
    });
  });

  // 處理 commodity_specifications
  cleanedData.commodity_specifications.forEach((spec: any) => {
    if (spec.spec_value_1_id < 0) {
      delete spec.spec_value_1_id;
    }
    if (spec.spec_value_2_id && spec.spec_value_2_id < 0) {
      delete spec.spec_value_2_id;
    }
    // 如果 commodity_spec_id 也需要前端生成並移除臨時 ID，可以在此處理
    if (spec.commodity_spec_id < 0) {
      delete spec.commodity_spec_id;
    }
  });

  return cleanedData;
};

// 提交資料到後端
const submit = async (commodity: CommodityDetailResponse) => {
  loading.value = true;

  try {
    // 剔除臨時 ID 後提交
    const cleanedCommodity = removeTempIds(commodity);
    await requests.put("commodities", cleanedCommodity);
    alert("成功");
  } catch (error) {
    console.error('提交失敗:', error);
    alert("提交失敗，請稍後再試");
  } finally {
    loading.value = false;
    await commodityStore.getCommodities();
  }
};
</script>

<style scoped></style>
