import { defineStore } from "pinia";
// import { useAuthPreferencesStore } from "./auth-preferences";
// import { useAuthEmailStore } from "./auth-email";

export const useLayoutStore = defineStore("layout", {
  // convert to a function
  state: () => ({
    loading: false,
    show: null,
    showConfirm: null,
    showAction: null,
    showShell: false,
  }),
  getters: {
    // user and jwt getter removed, no longer needed
  },
  actions: {
    // no context as first argument, use `this` instead
    toggleShell() {
      this.showShell = !this.showShell;
    },
    showHover(value) {
      if (typeof value !== "object") {
        this.show = value;
        return;
      }

      this.show = value.prompt;
      this.showConfirm = value.confirm;
      if (value.action !== undefined) {
        this.showAction = value.action;
      }
    },
    showError() {
      this.show = "error";
    },
    showSuccess() {
      this.show = "success";
    },
    closeHovers() {
      this.show = null;
      this.showConfirm = null;
      this.showAction = null;
    },
    // easily reset state using `$reset`
    clearLayout() {
      this.$reset();
    },
  },
});