import { defineStore } from "pinia";
import { computed, ref } from "vue";
import { login as loginApi } from "@/api/cloudops";

const TOKEN_KEY = "cloudops_token";
const USER_KEY = "cloudops_user";

export const useAuthStore = defineStore("auth", () => {
  const token = ref(localStorage.getItem(TOKEN_KEY) || "");
  const username = ref(localStorage.getItem(USER_KEY) || "");

  const isLoggedIn = computed(() => !!token.value);

  async function login(user: string, password: string) {
    const res = await loginApi(user, password);
    token.value = res.token;
    username.value = res.user.username;
    localStorage.setItem(TOKEN_KEY, res.token);
    localStorage.setItem(USER_KEY, res.user.username);
  }

  function logout() {
    token.value = "";
    username.value = "";
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(USER_KEY);
  }

  return { token, username, isLoggedIn, login, logout };
});
