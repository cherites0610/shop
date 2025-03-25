<template>
    <div v-if="commodity && selectedItem" class="flex gap-5">
        <img src="https://img.daisyui.com/images/stock/photo-1559703248-dcaaec9fab78.webp">
        <div class="flex flex-col gap-10">
            <span class="font-bold text-4xl">
                {{ commodity.name }}
            </span>
            <span class="font-bold text-3xl bg-gray-400">
                $ {{ selectedItem.price }}
            </span>

            <div>
                <div class="flex" v-for="(value, key) in commodity.spec" :key="key">
                    <span>{{ key }}:</span>
                    <div v-for="item in value" :key="item">
                        <label>
                            <input type="radio" :name="String(key)" class="radio" :value="item"
                                v-model="selectedSpecValues[key]" @change="updateSelectedItem">
                            {{ item }}
                        </label>
                    </div>
                </div>

                <div class="join">
                    <button @click="num--" class="join-item btn">-1</button>
                    <input type="number" v-model="num" class="join-item text-xl text-center" />
                    <button @click="num++" class="join-item btn">+1</button>
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
import { useRoute } from 'vue-router';
import { useCommodityStore } from '@/stores/commodityStore';
import { onMounted, ref, watch } from 'vue';
import type Commodity from '@/modals/Commodity';
import Error from '@/components/Error.vue';
import { useUserStore } from '@/stores/userStore';
import type { Specification } from '@/modals/Commodity';
import { requests } from '@/axios';

const commodityStore = useCommodityStore();
const route = useRoute();
const errorMessage = ref<string | null>(null);
const userStore = useUserStore();
const num = ref<number>(1);

// 商品數據
const commodity = ref<Commodity | undefined>(undefined);
// 當前選中的規格
const selectedItem = ref<Specification | undefined>(undefined);
// 儲存 radio 選擇的 spec_value，使用對象而不是陣列
const selectedSpecValues = ref<{ [key: string]: string }>({});

const getCommodity = async () => {
    const id = Number(route.params.id);
    if (!isNaN(id)) {
        const fetchedCommodity = await commodityStore.getCommodityById(id);
        commodity.value = fetchedCommodity;

        if (commodity.value) {
            // 初始化 selectedItem 为第一个规格
            selectedItem.value = commodity.value.specifications[0];
            // 初始化 selectedSpecValues
            selectedSpecValues.value = {};

            const specKeys = Object.keys(commodity.value.spec);
            const firstSpecValue = commodity.value.specifications[0].spec_value;

            // 为每个 spec 键找到匹配的值
            specKeys.forEach((key) => {
                const possibleValues = commodity.value!.spec[key];
                // 从 firstSpecValue 中找到属于 possibleValues 的值
                const matchedValue = firstSpecValue.find(val => possibleValues.includes(val));
                if (matchedValue) {
                    selectedSpecValues.value[key] = matchedValue;
                } else {
                    // 如果没有匹配的值，使用第一个可能值作为默认
                    selectedSpecValues.value[key] = possibleValues[0];
                }
            });
        }
    } else {
        console.error('Invalid ID');
    }
};

// 根據 selectedSpecValues 更新 selectedItem
const updateSelectedItem = () => {
    if (commodity.value) {
        const specValuesArray = Object.keys(commodity.value.spec).map(key => selectedSpecValues.value[key]);
        const matchedSpec = commodity.value.specifications.find(spec => {
            return spec.spec_value.sort().join(',') === specValuesArray.sort().join(',')
        }

        );
        if (matchedSpec) {
            selectedItem.value = matchedSpec;
            // 如果 num 超過新的庫存，調整 num
            if (num.value > matchedSpec.stock) {
                num.value = matchedSpec.stock;
            }
        } else {
            errorMessage.value = "未找到匹配的規格";
        }
    }
};

// 監聽 num 的變化
watch(() => num.value, (newVal) => {
    if (selectedItem.value) {
        if (newVal < 1) {
            num.value = 1; // 最小值為 1
        } else if (newVal > selectedItem.value.stock) {
            num.value = selectedItem.value.stock; // 限制為庫存上限
        } else {
            errorMessage.value = null; // 清除錯誤訊息
        }
    }
});

const buyHandler = () => {
    let specificationText = '';
    if(commodity.value&&selectedSpecValues.value) {
        Object.keys(commodity.value.spec).map(item => {
            specificationText += `${item}:${selectedSpecValues.value[item]},`
            
        })
    }
    console.log(specificationText);
    

    const message = `您購買了${commodity.value?.name}，規格為:${specificationText}，購買數量為:${num.value}`
    requests.post("/sendMessage", { "text": message, "userID": userStore.user?.userId }, {
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' }
    })
}

onMounted(async () => {
    await getCommodity();
});
</script>

<style scoped></style>
