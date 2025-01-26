<template>
  <div class="card floating">
    <div class="card-content">
      <p>
        {{ $t("prompts.discardEditorChanges") }}
      </p>
    </div>
    <div class="card-action">
      <button
        class="button button--flat button--grey"
        @click="closeHovers"
        :aria-label="$t('buttons.cancel')"
        :title="$t('buttons.cancel')"
        tabindex="2"
      >
        {{ $t("buttons.cancel") }}
      </button>
      <button
        id="focus-prompt"
        @click="submit"
        class="button button--flat button--red"
        :aria-label="$t('buttons.discardChanges')"
        :title="$t('buttons.discardChanges')"
        tabindex="1"
      >
        {{ $t("buttons.discardChanges") }}
      </button>
    </div>
  </div>
</template>

<script>
import { mapActions } from "pinia";
import url from "@/utils/url";
import { useLayoutStore } from "@/stores/layout";
import { useFileStore } from "@/stores/file";

export default {
  name: "discardEditorChanges",
  methods: {
    ...mapActions(useLayoutStore, ["closeHovers"]),
    ...mapActions(useFileStore, ["updateRequest"]),
    submit: async function () {
      this.updateRequest(null);

      const uri = url.removeLastDir(this.$route.path) + "/";
      this.$router.push({ path: uri });
    },
  },
};
</script>
