// จัดการ Token ใน localStorage
export const authService = {
  // บันทึก token และข้อมูล user
  login: (token, user) => {
    localStorage.setItem("token", token);
    localStorage.setItem("user", JSON.stringify(user));
  },

  // ลบ token และข้อมูล user
  logout: () => {
    localStorage.removeItem("token");
    localStorage.removeItem("user");
  },

  // ดึง token
  getToken: () => {
    return localStorage.getItem("token");
  },

  // ดึงข้อมูล user
  getUser: () => {
    const user = localStorage.getItem("user");
    return user ? JSON.parse(user) : null;
  },

  // เช็คว่า login แล้วหรือยัง
  isAuthenticated: () => {
    return !!localStorage.getItem("token");
  },
};

// Axios interceptor สำหรับใส่ token ใน header ทุก request
import axios from "axios";

axios.interceptors.request.use(
  (config) => {
    const token = authService.getToken();
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  },
);

// Handle 401 Unauthorized (token หมดอายุ)
axios.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      authService.logout();
      window.location.href = "/login";
    }
    return Promise.reject(error);
  },
);
