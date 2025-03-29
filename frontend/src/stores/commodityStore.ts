// stores/commodityStore.ts
import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { requests } from '@/axios';
import type { CommodityListResponse,CommodityDetailResponse } from "../models/Commodity";
import type { AxiosResponse } from 'axios';

export const useCommodityStore = defineStore('commodity', () => {
  const commodities = ref<CommodityListResponse[]>([]); // 用来保存商品数据
  const specifications = ref<Map<number, any>>(new Map()); // 用来保存商品规格

  // 获取商品列表
  const getCommodities = async () => {
    try {
      const response = await requests.get('/commodities');
      commodities.value = response.data;
    } catch (error) {
      console.error('Error fetching commodities:', error);
    }
  };

  // 获取单个商品
  const getCommodityById = async (id: number) => {
    try {
      const response: AxiosResponse<CommodityDetailResponse> = await requests.get(`/commodities/${id}`);
      return response.data;
    } catch (error) {
      console.error(`Error fetching commodity with id ${id}:`, error);
    }
  };


  return {
    commodities,
    specifications,
    getCommodities,
    getCommodityById,
  };
});
