import { createRouter, createWebHistory } from "vue-router";
import { useAuthStore } from "@/stores/auth";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/login",
      name: "login",
      component: () => import("@/views/LoginView.vue"),
      meta: { public: true },
    },
    {
      path: "/",
      component: () => import("@/layouts/MainLayout.vue"),
      redirect: "/nodes",
      children: [
        { path: "nodes", name: "nodes", component: () => import("@/views/NodesView.vue") },
        { path: "pods", name: "pods", component: () => import("@/views/PodsView.vue") },
        {
          path: "deployments",
          name: "deployments",
          component: () => import("@/views/DeploymentsView.vue"),
        },
      ],
    },
  ],
});

router.beforeEach((to) => {
  const auth = useAuthStore();
  if (!to.meta.public && !auth.isLoggedIn) {
    return { name: "login", query: { redirect: to.fullPath } };
  }
  if (to.name === "login" && auth.isLoggedIn) {
    return { name: "nodes" };
  }
  return true;
});

export default router;
