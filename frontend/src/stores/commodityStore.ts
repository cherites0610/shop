// stores/commodityStore.ts
import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { requests } from '@/axios';
import type Commodity from '@/modals/Commodity';

export const useCommodityStore = defineStore('commodity', () => {
  const commodities = ref<Commodity[]>([]); // 用来保存商品数据
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
      return commodities.value.find(item => item.id===id)
    } catch (error) {
      console.error(`Error fetching commodity with id ${id}:`, error);
    }
  };

  // 创建商品
  const createCommodity = async (commodity: Commodity) => {
    try {
      const response = await requests.post('/commodities', commodity);
      commodities.value.push(response.data);
    } catch (error) {
      console.error('Error creating commodity:', error);
    }
  };

  // 更新商品
  const updateCommodity = async (id: number, commodity: Commodity) => {
    try {
      await requests.put(`/commodities/${id}`, commodity);
      const index = commodities.value.findIndex((c) => c.id === id);
      if (index !== -1) {
        commodities.value[index] = { ...commodities.value[index], ...commodity };
      }
    } catch (error) {
      console.error(`Error updating commodity with id ${id}:`, error);
    }
  };

  // 删除商品
  const deleteCommodity = async (id: number) => {
    try {
      await requests.delete(`/commodities/${id}`);
      commodities.value = commodities.value.filter((c) => c.id !== id);
    } catch (error) {
      console.error(`Error deleting commodity with id ${id}:`, error);
    }
  };

  // 存储商品规格
  const setSpecification = (id: number, specification: any) => {
    specifications.value.set(id, specification);
  };

  // 获取商品规格
  const getSpecification = (id: number) => {
    return specifications.value.get(id);
  };

  return {
    commodities,
    specifications,
    getCommodities,
    getCommodityById,
    createCommodity,
    updateCommodity,
    deleteCommodity,
    setSpecification,
    getSpecification,
  };
});
