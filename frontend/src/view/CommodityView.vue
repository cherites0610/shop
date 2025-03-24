<template>
    <div v-if="commodity" class="flex gap-5">
        <img src="https://img.daisyui.com/images/stock/photo-1559703248-dcaaec9fab78.webp">
        <div class="flex flex-col gap-10">
            <span class="font-bold text-4xl">
                {{ commodity.name }}
            </span>
            <span class="font-bold text-3xl bg-gray-400">
                $ {{ commodity.price }}
            </span>
            <div v-if="commodity.specification">
                <div class="flex" v-for="(value, key) in commodity.specification" :key="key">
                    <span>{{ key }}:</span>
                    <div v-for="item in value" :key="item">
                        <label>
                            <input type="radio" :name="item" class="radio" v-model="selectedItem.specification[key]"
                                :value="item">
                            {{ item }}
                        </label>
                    </div>
                </div>

                <div class="join">
                    <button @click="selectedItem.number--" class="join-item btn">-1</button>
                    <input type="number" v-model="selectedItem.number" class="join-item" />
                    <button @click="selectedItem.number++" class="join-item btn">+1</button>
                    <span>庫存:{{ commodity.stock }}</span>
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
import { requests } from '@/axios';

const commodityStore = useCommodityStore();
const route = useRoute();
const errorMessage = ref<string | null>(null);
const userStore = useUserStore()

const commodity = ref<Commodity | undefined>(undefined);
const selectedItem = ref<{ id: number, specification: { [key: string]: string }, number: number }>({
    id: -1,
    specification: {},
    number: 1,
});

const buyHandler = () => {
    console.log("Buy!");
    console.log(userStore.user?.userId);
    requests.post("/sendMessage",{"text":"Test","userID":userStore.user?.userId},{
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' }
    })
    
    
}

const getCommodity = async () => {
    const id = Number(route.params.id);
    if (!isNaN(id)) {
        const fetchedCommodity = await commodityStore.getCommodityById(id);
        commodity.value = fetchedCommodity;

        if (commodity.value) {
            selectedItem.value.id = commodity.value.id;

            // 自动选择每个 specification 的第一个值
            for (const key in commodity.value.specification) {
                if (commodity.value.specification[key].length > 0) {
                    selectedItem.value.specification[key] = commodity.value.specification[key][0];

                }
            }
        }
    } else {
        console.error('Invalid ID');
    }
};

watch(() => selectedItem.value.number, (newVal) => {
    if (commodity.value) {
        if (newVal < 1) {
            selectedItem.value.number = 1;  // 防止设置为负数
        } else if (newVal > commodity.value.stock) {
            selectedItem.value.number = commodity.value.stock;  // 限制最大值
        } else {
            errorMessage.value = null;  // 清除错误消息
        }
    }

});


onMounted(() => {
    getCommodity();
});
</script>

<style scoped></style>
