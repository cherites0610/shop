import { requests } from "@/axios";
import type User from "@/models/User";
import { defineStore } from "pinia";
import { ref } from "vue";

export const useUserStore = defineStore("user", () => {
    const user = ref<User>();

    const getUserProfile = async (accessToken: string) => {
        const url = 'https://api.line.me/v2/profile';

        try {
            const response = await requests.get(url, {
                headers: {
                    'Authorization': `Bearer ${accessToken}`
                }
            });
            user.value = response.data;
        } catch (error: any) {
            throw(error);
        }
    }

    return {
        user,
        getUserProfile
    }
})