<template>
  <select name="selectLanguage" v-on:change="change" :value="locale">
    <option v-for="(language, value) in locales" :key="value" :value="value">
      {{ language }}
    </option>
  </select>
</template>

<script>
import { markRaw } from "vue";

export default {
  name: "languages",
  props: ["locale"],
  data() {
    const dataObj = {};
    const locales = {
      he: "עברית",
      hu: "Magyar",
      ar: "العربية",
      ca: "Català",
      de: "Deutsch",
      el: "Ελληνικά",
      en: "English",
      es: "Español",
      fr: "Français",
      is: "Icelandic",
      it: "Italiano",
      ja: "日本語",
      ko: "한국어",
      "nl-be": "Dutch (Belgium)",
      pl: "Polski",
      "pt-br": "Português",
      pt: "Português (Brasil)",
      ro: "Romanian",
      ru: "Русский",
      sk: "Slovenčina",
      "sv-se": "Swedish (Sweden)",
      tr: "Türkçe",
      uk: "Українська",
      "zh-cn": "中文 (简体)",
      "zh-tw": "中文 (繁體)",
    };

    // Vue3 reactivity breaks with this configuration
    // so we need to use markRaw as a workaround
    // https://github.com/vuejs/core/issues/3024
    Object.defineProperty(dataObj, "locales", {
      value: markRaw(locales),
      configurable: false,
      writable: false,
    });

    return dataObj;
  },
  methods: {
    change(event) {
      this.$emit("update:locale", event.target.value);
    },
  },
};
</script>
