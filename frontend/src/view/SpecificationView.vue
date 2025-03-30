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

const commodity = ref<CommodityDetailResponse>()
const commodityStore = useCommodityStore()
const route = useRoute()
const newSpecTypeName = ref(''); // 用於輸入新規格類型名稱
const loading = ref(false)

onMounted(async () => {
    const id = Number(route.params.id);
    if (isNaN(id)) {
        console.error('Invalid ID');
        return;
    }
    commodity.value = await commodityStore.getCommodityById(id);
})

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

// 刪除指定的規格類型並重新生成 commodity_specifications
const deleteSpecType = (specTypeId: number) => {
    if (!commodity.value || !commodity.value.specification_types) return;

    // 找到並移除指定的規格類型
    const specTypeIndex = commodity.value.specification_types.findIndex(
        (type) => type.spec_type_id === specTypeId
    );
    if (specTypeIndex === -1) return;

    commodity.value.specification_types.splice(specTypeIndex, 1);

    // 重新生成 commodity_specifications
    commodity.value.commodity_specifications = generateCommoditySpecifications();
};

// 生成 commodity_specifications 的默認值
const generateCommoditySpecifications = () => {
    if (!commodity.value || !commodity.value.specification_types) return [];

    const specTypes = commodity.value.specification_types;
    let combinations: any[] = [];

    if (specTypes.length === 0) return [];

    // 獲取所有規格值的 ID 和名稱
    const specValues1 = specTypes[0].specification_values;
    const specValues2 = specTypes[1] ? specTypes[1].specification_values : [{ spec_value_id: 0, spec_value: '' }];

    // 生成規格組合（笛卡爾積）
    specValues1.forEach((val1, index1) => {
        specValues2.forEach((val2, index2) => {
            const commoditySpecId = commodity.value!.commodity_specifications.length + index1 * specValues2.length + index2 + 1;
            combinations.push({
                spec_value_1_id: val1.spec_value_id,
                spec_value_1: val1.spec_value,
                spec_value_2_id: val2.spec_value_id || null, // 如果只有一個 spec_type，設置為 0
                spec_value_2: val2.spec_value || null,
                price: 0,
                stock: 0,
                picture_url: 'https://picsum.photos/200/300', // 默認圖片，可根據需要修改
            });
        });
    });
    console.log(combinations);

    return combinations;
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

    // 新增規格類型，包含一個默認規格值，不指定 spec_type_id 和 spec_value_id
    commodity.value.specification_types.push({
        spec_type_name: newSpecTypeName.value,
        specification_values: [
            { spec_value: '默認值' } // 不指定 spec_value_id
        ]
    });

    // 清空輸入框
    newSpecTypeName.value = '';

    // 重新生成 commodity_specifications
    commodity.value.commodity_specifications = generateCommoditySpecifications();
};

// 新增規格值
const addSpecValue = (specTypeId: number) => {
    if (!commodity.value || !commodity.value.specification_types) return;

    const specType = commodity.value.specification_types.find(
        (type) => type.spec_type_id === specTypeId
    );
    if (!specType) return;


    // 添加新規格值
    specType.specification_values.push({
        spec_value: '新值' // 可改為用戶輸入，例如通過額外輸入框
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

const submit = async (commodity: CommodityDetailResponse) => {
    loading.value = true;

    try {
        // 更改名字
        await requests.put(`commodities/${commodity.commodity_id}`, {
            "commodity_name": commodity.commodity_name
        });

        // 更改規格類型
        const requestSpec = commodity.specification_types.map((item) => {
            if (item.spec_type_id) {
                return {
                    "spec_type_id": item.spec_type_id,
                    "spec_type_name": item.spec_type_name,
                    "spec_type_values": item.specification_values.map((value) => value.spec_value)
                };
            } else {
                return {
                    "spec_type_name": item.spec_type_name,
                    "spec_type_values": item.specification_values.map((value) => value.spec_value)
                };
            }

        });
        await requests.put(`commodities/${commodity.commodity_id}/specification-types`, requestSpec);

        // 更改 SKU
        await Promise.all(commodity.commodity_specifications.map(async (item) => {
            if (!item.commodity_spec_id) {
                // 新增
                await requests.post(`/commodities/${commodity.commodity_id}/sku`, {
                    ...item,
                    "commodity_id": commodity.commodity_id
                });
            } else {
                // 修改
                await requests.put(`/commodities/${commodity.commodity_id}/sku/${item.commodity_spec_id}`, {
                    ...item,
                    "commodity_id": commodity.commodity_id
                });
            }
        }));

        loading.value = false;
    } catch (error) {
        console.error("提交失敗:", error);
        loading.value = false;
    }

    alert("成功")
};

</script>

<style scoped></style>