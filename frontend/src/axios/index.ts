import axios from "axios";

export const requests = axios.create({
    baseURL: 'https://api.pt-sms.org/',
    timeout: 1000,
  });