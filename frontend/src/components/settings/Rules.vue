<template>
  <form class="rules small">
    <div v-for="(rule, index) in rules" :key="index">
      <input type="checkbox" v-model="rule.regex" /><label>Regex</label>
      <input type="checkbox" v-model="rule.allow" /><label>Allow</label>

      <input
        @keypress.enter.prevent
        type="text"
        v-if="rule.regex"
        v-model="rule.regexp.raw"
        :placeholder="$t('settings.insertRegex')"
      />
      <input
        @keypress.enter.prevent
        type="text"
        v-else
        v-model="rule.path"
        :placeholder="$t('settings.insertPath')"
      />

      <button class="button button--red" @click="remove($event, index)">
        -
      </button>
    </div>

    <div>
      <button class="button" @click="create" default="false">
        {{ $t("buttons.new") }}
      </button>
    </div>
  </form>
</template>

<script>
export default {
  name: "rules-textarea",
  props: ["rules"],
  methods: {
    remove(event, index) {
      event.preventDefault();
      const rules = [...this.rules];
      rules.splice(index, 1);
      this.$emit("update:rules", [...rules]);
    },
    create(event) {
      event.preventDefault();

      this.$emit("update:rules", [
        ...this.rules,
        {
          allow: true,
          path: "",
          regex: false,
          regexp: {
            raw: "",
          },
        },
      ]);
    },
  },
};
</script>
