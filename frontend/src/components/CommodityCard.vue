<template>
    <div class="card bg-base-100 w-96 shadow-sm">
        <figure>
            <img :src="commodity.picture_url" :alt=commodity.commodity_name />
        </figure>
        <div class="card-body">
            <div class="flex flex-col gap-1">
                <h2 class=" text-3xl font-bold">{{ commodity.commodity_name }}</h2>
                <span>價格:{{ commodity.price_range.min }} - {{ commodity.price_range.max }}</span>
                <span>總庫存:{{ commodity.total_stock }}</span>
                <div class="flex gap-2 flex-row-reverse">

                    <button @click="router.push(`/commodity/${commodity.commodity_id}`)" class="btn btn-primary">
                        查看
                    </button>
                    <button @click="router.push(`/specification/${commodity.commodity_id}`)" class="btn">管理</button>
                </div>

            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import type { CommodityListResponse } from '@/models/Commodity';
import router from '@/router';
const { commodity } = defineProps<{ commodity: CommodityListResponse }>();
</script>

<style scoped></style>