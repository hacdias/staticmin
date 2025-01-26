<template>
  <select v-on:change="change" :value="theme">
    <option value="">{{ t("settings.themes.default") }}</option>
    <option value="light">{{ t("settings.themes.light") }}</option>
    <option value="dark">{{ t("settings.themes.dark") }}</option>
  </select>
</template>

<script setup lang="ts">
import type { SelectHTMLAttributes } from "vue";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

defineProps<{
  theme: UserTheme;
}>();

const emit = defineEmits<{
  (e: "update:theme", val: string | null): void;
}>();

const change = (event: Event) => {
  emit("update:theme", (event.target as SelectHTMLAttributes)?.value);
};
</script>
