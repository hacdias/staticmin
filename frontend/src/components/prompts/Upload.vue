<template>
  <div class="card floating">
    <div class="card-title">
      <h2>{{ t("prompts.upload") }}</h2>
    </div>

    <div class="card-content">
      <p>{{ t("prompts.uploadMessage") }}</p>
    </div>

    <div class="card-action full">
      <div
        @click="uploadFile"
        @keypress.enter="uploadFile"
        class="action"
        id="focus-prompt"
        tabindex="1"
      >
        <i class="material-icons">insert_drive_file</i>
        <div class="title">{{ t("buttons.file") }}</div>
      </div>
      <div
        @click="uploadFolder"
        @keypress.enter="uploadFolder"
        class="action"
        tabindex="2"
      >
        <i class="material-icons">folder</i>
        <div class="title">{{ t("buttons.folder") }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from "vue-i18n";
import { useRoute } from "vue-router";
import { useFileStore } from "@/stores/file";
import { useLayoutStore } from "@/stores/layout";

import * as upload from "@/utils/upload";

const { t } = useI18n();
const route = useRoute();

const fileStore = useFileStore();
const layoutStore = useLayoutStore();

// TODO: this is a copy of the same function in FileListing.vue
const uploadInput = (event: Event) => {
  const files = (event.currentTarget as HTMLInputElement)?.files;
  if (files === null) return;

  const folder_upload = !!files[0].webkitRelativePath;

  const uploadFiles: UploadList = [];
  for (let i = 0; i < files.length; i++) {
    const file = files[i];
    const fullPath = folder_upload ? file.webkitRelativePath : undefined;
    uploadFiles.push({
      file,
      name: file.name,
      size: file.size,
      isDir: false,
      fullPath,
    });
  }

  const path = route.path.endsWith("/") ? route.path : route.path + "/";
  const conflict = upload.checkConflict(uploadFiles, fileStore.req!.items);

  if (conflict) {
    layoutStore.showHover({
      prompt: "replace",
      action: (event: Event) => {
        event.preventDefault();
        layoutStore.closeHovers();
        upload.handleFiles(uploadFiles, path, false);
      },
      confirm: (event: Event) => {
        event.preventDefault();
        layoutStore.closeHovers();
        upload.handleFiles(uploadFiles, path, true);
      },
    });

    return;
  }

  upload.handleFiles(uploadFiles, path);
};

const openUpload = (isFolder: boolean) => {
  const input = document.createElement("input");
  input.type = "file";
  input.multiple = true;
  input.webkitdirectory = isFolder;
  // TODO: call the function in FileListing.vue instead
  input.onchange = uploadInput;
  input.click();
};

const uploadFile = () => {
  openUpload(false);
};
const uploadFolder = () => {
  openUpload(true);
};
</script>
