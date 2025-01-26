<template>
  <div>
    <header-bar showMenu showLogo>
      <search />
      <title />
      <action
        class="search-button"
        icon="search"
        :label="t('buttons.search')"
        @action="openSearch()"
      />

      <template #actions>
        <template v-if="!isMobile">
          <action
            v-if="headerButtons.share"
            icon="share"
            :label="t('buttons.share')"
            show="share"
          />
          <action
            v-if="headerButtons.rename"
            icon="mode_edit"
            :label="t('buttons.rename')"
            show="rename"
          />
          <action
            v-if="headerButtons.copy"
            id="copy-button"
            icon="content_copy"
            :label="t('buttons.copyFile')"
            show="copy"
          />
          <action
            v-if="headerButtons.move"
            id="move-button"
            icon="forward"
            :label="t('buttons.moveFile')"
            show="move"
          />
          <action
            v-if="headerButtons.delete"
            id="delete-button"
            icon="delete"
            :label="t('buttons.delete')"
            show="delete"
          />
        </template>

        <action
          v-if="headerButtons.shell"
          icon="code"
          :label="t('buttons.shell')"
          @action="layoutStore.toggleShell"
        />
        <action
          :icon="viewIcon"
          :label="t('buttons.switchView')"
          @action="switchView"
        />
        <action
          v-if="headerButtons.download"
          icon="file_download"
          :label="t('buttons.download')"
          @action="download"
          :counter="fileStore.selectedCount"
        />
        <action
          v-if="headerButtons.upload"
          icon="file_upload"
          id="upload-button"
          :label="t('buttons.upload')"
          @action="uploadFunc"
        />
        <action icon="info" :label="t('buttons.info')" show="info" />
        <action
          icon="check_circle"
          :label="t('buttons.selectMultiple')"
          @action="toggleMultipleSelection"
        />
      </template>
    </header-bar>

    <div v-if="isMobile" id="file-selection">
      <span v-if="fileStore.selectedCount > 0">
        {{ t("prompts.filesSelected", fileStore.selectedCount) }}
      </span>
      <action
        v-if="headerButtons.share"
        icon="share"
        :label="t('buttons.share')"
        show="share"
      />
      <action
        v-if="headerButtons.rename"
        icon="mode_edit"
        :label="t('buttons.rename')"
        show="rename"
      />
      <action
        v-if="headerButtons.copy"
        icon="content_copy"
        :label="t('buttons.copyFile')"
        show="copy"
      />
      <action
        v-if="headerButtons.move"
        icon="forward"
        :label="t('buttons.moveFile')"
        show="move"
      />
      <action
        v-if="headerButtons.delete"
        icon="delete"
        :label="t('buttons.delete')"
        show="delete"
      />
    </div>

    <div v-if="layoutStore.loading">
      <h2 class="message delayed">
        <div class="spinner">
          <div class="bounce1"></div>
          <div class="bounce2"></div>
          <div class="bounce3"></div>
        </div>
        <span>{{ t("files.loading") }}</span>
      </h2>
    </div>
    <template v-else>
      <div
        v-if="
          (fileStore.req?.numDirs ?? 0) + (fileStore.req?.numFiles ?? 0) == 0
        "
      >
        <h2 class="message">
          <i class="material-icons">sentiment_dissatisfied</i>
          <span>{{ t("files.lonely") }}</span>
        </h2>
        <input
          style="display: none"
          type="file"
          id="upload-input"
          @change="uploadInput($event)"
          multiple
        />
        <input
          style="display: none"
          type="file"
          id="upload-folder-input"
          @change="uploadInput($event)"
          webkitdirectory
          multiple
        />
      </div>
      <div
        v-else
        id="listing"
        ref="listing"
        class="file-icons"
        :class="authStore.user?.viewMode ?? ''"
      >
        <div>
          <div class="item header">
            <div></div>
            <div>
              <p
                :class="{ active: nameSorted }"
                class="name"
                role="button"
                tabindex="0"
                @click="sort('name')"
                :title="t('files.sortByName')"
                :aria-label="t('files.sortByName')"
              >
                <span>{{ t("files.name") }}</span>
                <i class="material-icons">{{ nameIcon }}</i>
              </p>

              <p
                :class="{ active: sizeSorted }"
                class="size"
                role="button"
                tabindex="0"
                @click="sort('size')"
                :title="t('files.sortBySize')"
                :aria-label="t('files.sortBySize')"
              >
                <span>{{ t("files.size") }}</span>
                <i class="material-icons">{{ sizeIcon }}</i>
              </p>
              <p
                :class="{ active: modifiedSorted }"
                class="modified"
                role="button"
                tabindex="0"
                @click="sort('modified')"
                :title="t('files.sortByLastModified')"
                :aria-label="t('files.sortByLastModified')"
              >
                <span>{{ t("files.lastModified") }}</span>
                <i class="material-icons">{{ modifiedIcon }}</i>
              </p>
            </div>
          </div>
        </div>

        <h2 v-if="fileStore.req?.numDirs ?? false">
          {{ t("files.folders") }}
        </h2>
        <div v-if="fileStore.req?.numDirs ?? false">
          <item
            v-for="item in dirs"
            :key="base64(item.name)"
            v-bind:index="item.index"
            v-bind:name="item.name"
            v-bind:isDir="item.isDir"
            v-bind:url="item.url"
            v-bind:modified="item.modified"
            v-bind:type="item.type"
            v-bind:size="item.size"
            v-bind:path="item.path"
          >
          </item>
        </div>

        <h2 v-if="fileStore.req?.numFiles ?? false">{{ t("files.files") }}</h2>
        <div v-if="fileStore.req?.numFiles ?? false">
          <item
            v-for="item in files"
            :key="base64(item.name)"
            v-bind:index="item.index"
            v-bind:name="item.name"
            v-bind:isDir="item.isDir"
            v-bind:url="item.url"
            v-bind:modified="item.modified"
            v-bind:type="item.type"
            v-bind:size="item.size"
            v-bind:path="item.path"
          >
          </item>
        </div>

        <input
          style="display: none"
          type="file"
          id="upload-input"
          @change="uploadInput($event)"
          multiple
        />
        <input
          style="display: none"
          type="file"
          id="upload-folder-input"
          @change="uploadInput($event)"
          webkitdirectory
          multiple
        />

        <div :class="{ active: fileStore.multiple }" id="multiple-selection">
          <p>{{ t("files.multipleSelectionEnabled") }}</p>
          <div
            @click="() => (fileStore.multiple = false)"
            tabindex="0"
            role="button"
            :title="t('buttons.clear')"
            :aria-label="t('buttons.clear')"
            class="action"
          >
            <i class="material-icons">clear</i>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from "@/stores/auth";
import { useClipboardStore } from "@/stores/clipboard";
import { useFileStore } from "@/stores/file";
import { useLayoutStore } from "@/stores/layout";

import { users, files as api } from "@/api";
import { enableExec } from "@/utils/constants";
import * as upload from "@/utils/upload";
import css from "@/utils/css";
import { throttle } from "lodash-es";
import { Base64 } from "js-base64";

import HeaderBar from "@/components/header/HeaderBar.vue";
import Action from "@/components/header/Action.vue";
import Search from "@/components/Search.vue";
import Item from "@/components/files/ListingItem.vue";
import {
  computed,
  inject,
  nextTick,
  onBeforeUnmount,
  onMounted,
  ref,
  watch,
} from "vue";
import { useRoute } from "vue-router";
import { useI18n } from "vue-i18n";
import { storeToRefs } from "pinia";

const showLimit = ref<number>(50);
const columnWidth = ref<number>(280);
const dragCounter = ref<number>(0);
const width = ref<number>(window.innerWidth);
const itemWeight = ref<number>(0);

const $showError = inject<IToastError>("$showError")!;

const clipboardStore = useClipboardStore();
const authStore = useAuthStore();
const fileStore = useFileStore();
const layoutStore = useLayoutStore();

const { req } = storeToRefs(fileStore);

const route = useRoute();

const { t } = useI18n();

const listing = ref<HTMLElement | null>(null);

const nameSorted = computed(() =>
  fileStore.req ? fileStore.req.sorting.by === "name" : false
);

const sizeSorted = computed(() =>
  fileStore.req ? fileStore.req.sorting.by === "size" : false
);

const modifiedSorted = computed(() =>
  fileStore.req ? fileStore.req.sorting.by === "modified" : false
);

const ascOrdered = computed(() =>
  fileStore.req ? fileStore.req.sorting.asc : false
);

const dirs = computed(() => items.value.dirs.slice(0, showLimit.value));

const items = computed(() => {
  const dirs: any[] = [];
  const files: any[] = [];

  fileStore.req?.items.forEach((item) => {
    if (item.isDir) {
      dirs.push(item);
    } else {
      files.push(item);
    }
  });

  return { dirs, files };
});

const files = computed((): Resource[] => {
  let _showLimit = showLimit.value - items.value.dirs.length;

  if (_showLimit < 0) _showLimit = 0;

  return items.value.files.slice(0, _showLimit);
});

const nameIcon = computed(() => {
  if (nameSorted.value && !ascOrdered.value) {
    return "arrow_upward";
  }

  return "arrow_downward";
});

const sizeIcon = computed(() => {
  if (sizeSorted.value && ascOrdered.value) {
    return "arrow_downward";
  }

  return "arrow_upward";
});

const modifiedIcon = computed(() => {
  if (modifiedSorted.value && ascOrdered.value) {
    return "arrow_downward";
  }

  return "arrow_upward";
});

const viewIcon = computed(() => {
  const icons = {
    list: "view_module",
    mosaic: "grid_view",
    "mosaic gallery": "view_list",
  };
  return authStore.user === null
    ? icons["list"]
    : icons[authStore.user.viewMode];
});

const headerButtons = computed(() => {
  return {
    upload: authStore.user?.perm.create,
    download: authStore.user?.perm.download,
    shell: authStore.user?.perm.execute && enableExec,
    delete: fileStore.selectedCount > 0 && authStore.user?.perm.delete,
    rename: fileStore.selectedCount === 1 && authStore.user?.perm.rename,
    share: fileStore.selectedCount === 1 && authStore.user?.perm.share,
    move: fileStore.selectedCount > 0 && authStore.user?.perm.rename,
    copy: fileStore.selectedCount > 0 && authStore.user?.perm.create,
  };
});

const isMobile = computed(() => {
  return width.value <= 736;
});

watch(req, () => {
  // Reset the show value
  if (
    window.sessionStorage.getItem("listFrozen") !== "true" &&
    window.sessionStorage.getItem("modified") !== "true"
  ) {
    showLimit.value = 50;

    nextTick(() => {
      // Ensures that the listing is displayed
      // How much every listing item affects the window height
      setItemWeight();

      // Fill and fit the window with listing items
      fillWindow(true);
    });
  }
  if (req.value?.isDir) {
    window.sessionStorage.setItem("listFrozen", "false");
    window.sessionStorage.setItem("modified", "false");
  }
});

onMounted(() => {
  // Check the columns size for the first time.
  colunmsResize();

  // How much every listing item affects the window height
  setItemWeight();

  // Fill and fit the window with listing items
  fillWindow(true);

  // Add the needed event listeners to the window and document.
  window.addEventListener("keydown", keyEvent);
  window.addEventListener("scroll", scrollEvent);
  window.addEventListener("resize", windowsResize);

  if (!authStore.user?.perm.create) return;
  document.addEventListener("dragover", preventDefault);
  document.addEventListener("dragenter", dragEnter);
  document.addEventListener("dragleave", dragLeave);
  document.addEventListener("drop", drop);
});

onBeforeUnmount(() => {
  // Remove event listeners before destroying this page.
  window.removeEventListener("keydown", keyEvent);
  window.removeEventListener("scroll", scrollEvent);
  window.removeEventListener("resize", windowsResize);

  if (authStore.user && !authStore.user?.perm.create) return;
  document.removeEventListener("dragover", preventDefault);
  document.removeEventListener("dragenter", dragEnter);
  document.removeEventListener("dragleave", dragLeave);
  document.removeEventListener("drop", drop);
});

const base64 = (name: string) => Base64.encodeURI(name);

const keyEvent = (event: KeyboardEvent) => {
  // No prompts are shown
  if (layoutStore.currentPrompt !== null) {
    return;
  }

  if (event.key === "Escape") {
    // Reset files selection.
    fileStore.selected = [];
  }

  if (event.key === "Delete") {
    if (!authStore.user?.perm.delete || fileStore.selectedCount == 0) return;

    // Show delete prompt.
    layoutStore.showHover("delete");
  }

  if (event.key === "F2") {
    if (!authStore.user?.perm.rename || fileStore.selectedCount !== 1) return;

    // Show rename prompt.
    layoutStore.showHover("rename");
  }

  // Ctrl is pressed
  if (!event.ctrlKey && !event.metaKey) {
    return;
  }

  switch (event.key) {
    case "f":
      event.preventDefault();
      layoutStore.showHover("search");
      break;
    case "c":
    case "x":
      copyCut(event);
      break;
    case "v":
      paste(event);
      break;
    case "a":
      event.preventDefault();
      for (const file of items.value.files) {
        if (fileStore.selected.indexOf(file.index) === -1) {
          fileStore.selected.push(file.index);
        }
      }
      for (const dir of items.value.dirs) {
        if (fileStore.selected.indexOf(dir.index) === -1) {
          fileStore.selected.push(dir.index);
        }
      }
      break;
    case "s":
      event.preventDefault();
      document.getElementById("download-button")?.click();
      break;
  }
};

const preventDefault = (event: Event) => {
  // Wrapper around prevent default.
  event.preventDefault();
};

const copyCut = (event: Event | KeyboardEvent): void => {
  if ((event.target as HTMLElement).tagName?.toLowerCase() === "input") return;

  if (fileStore.req === null) return;

  const items = [];

  for (const i of fileStore.selected) {
    items.push({
      from: fileStore.req.items[i].url,
      name: fileStore.req.items[i].name,
    });
  }

  if (items.length === 0) {
    return;
  }

  clipboardStore.$patch({
    key: (event as KeyboardEvent).key,
    items,
    path: route.path,
  });
};

const paste = (event: Event) => {
  if ((event.target as HTMLElement).tagName?.toLowerCase() === "input") return;

  // TODO router location should it be
  const items: any[] = [];

  for (const item of clipboardStore.items) {
    const from = item.from.endsWith("/") ? item.from.slice(0, -1) : item.from;
    const to = route.path + encodeURIComponent(item.name);
    items.push({ from, to, name: item.name });
  }

  if (items.length === 0) {
    return;
  }

  let action = (overwrite: boolean, rename: boolean) => {
    api
      .copy(items, overwrite, rename)
      .then(() => {
        fileStore.reload = true;
      })
      .catch($showError);
  };

  if (clipboardStore.key === "x") {
    action = (overwrite, rename) => {
      api
        .move(items, overwrite, rename)
        .then(() => {
          clipboardStore.resetClipboard();
          fileStore.reload = true;
        })
        .catch($showError);
    };
  }

  if (clipboardStore.path == route.path) {
    action(false, true);

    return;
  }

  const conflict = upload.checkConflict(items, fileStore.req!.items);

  let overwrite = false;
  let rename = false;

  if (conflict) {
    layoutStore.showHover({
      prompt: "replace-rename",
      confirm: (event: Event, option: string) => {
        overwrite = option == "overwrite";
        rename = option == "rename";

        event.preventDefault();
        layoutStore.closeHovers();
        action(overwrite, rename);
      },
    });

    return;
  }

  action(overwrite, rename);
};

const colunmsResize = () => {
  // Update the columns size based on the window width.
  const items_ = css(["#listing.mosaic .item", ".mosaic#listing .item"]);
  if (items_ === null) return;

  let columns = Math.floor(
    (document.querySelector("main")?.offsetWidth ?? 0) / columnWidth.value
  );
  if (columns === 0) columns = 1;
  items_.style.width = `calc(${100 / columns}% - 1em)`;
};

const scrollEvent = throttle(() => {
  const totalItems =
    (fileStore.req?.numDirs ?? 0) + (fileStore.req?.numFiles ?? 0);

  // All items are displayed
  if (showLimit.value >= totalItems) return;

  const currentPos = window.innerHeight + window.scrollY;

  // Trigger at the 75% of the window height
  const triggerPos = document.body.offsetHeight - window.innerHeight * 0.25;

  if (currentPos > triggerPos) {
    // Quantity of items needed to fill 2x of the window height
    const showQuantity = Math.ceil((window.innerHeight * 2) / itemWeight.value);

    // Increase the number of displayed items
    showLimit.value += showQuantity;
  }
}, 100);

const dragEnter = () => {
  dragCounter.value++;

  // When the user starts dragging an item, put every
  // file on the listing with 50% opacity.
  const items = document.getElementsByClassName("item");

  Array.from(items).forEach((file: Element) => {
    (file as HTMLElement).style.opacity = "0.5";
  });
};

const dragLeave = () => {
  dragCounter.value--;

  if (dragCounter.value == 0) {
    resetOpacity();
  }
};

const drop = async (event: DragEvent) => {
  event.preventDefault();
  dragCounter.value = 0;
  resetOpacity();

  const dt = event.dataTransfer;
  let el: HTMLElement | null = event.target as HTMLElement;

  if (fileStore.req === null || dt === null || dt.files.length <= 0) return;

  for (let i = 0; i < 5; i++) {
    if (el !== null && !el.classList.contains("item")) {
      el = el.parentElement;
    }
  }

  const files: UploadList = (await upload.scanFiles(dt)) as UploadList;
  let items = fileStore.req.items;
  let path = route.path.endsWith("/") ? route.path : route.path + "/";

  if (
    el !== null &&
    el.classList.contains("item") &&
    el.dataset.dir === "true"
  ) {
    // Get url from ListingItem instance
    // TODO: Don't know what is happening here
    path = el.__vue__.url;

    try {
      items = (await api.fetch(path)).items;
    } catch (error: any) {
      $showError(error);
    }
  }

  const conflict = upload.checkConflict(files, items);

  if (conflict) {
    layoutStore.showHover({
      prompt: "replace",
      action: (event: Event) => {
        event.preventDefault();
        layoutStore.closeHovers();
        upload.handleFiles(files, path, false);
      },
      confirm: (event: Event) => {
        event.preventDefault();
        layoutStore.closeHovers();
        upload.handleFiles(files, path, true);
      },
    });

    return;
  }

  upload.handleFiles(files, path);
};

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

const resetOpacity = () => {
  const items = document.getElementsByClassName("item");

  Array.from(items).forEach((file: Element) => {
    (file as HTMLElement).style.opacity = "1";
  });
};

const sort = async (by: string) => {
  let asc = false;

  if (by === "name") {
    if (nameIcon.value === "arrow_upward") {
      asc = true;
    }
  } else if (by === "size") {
    if (sizeIcon.value === "arrow_upward") {
      asc = true;
    }
  } else if (by === "modified") {
    if (modifiedIcon.value === "arrow_upward") {
      asc = true;
    }
  }

  try {
    if (authStore.user?.id) {
      await users.update({ id: authStore.user?.id, sorting: { by, asc } }, [
        "sorting",
      ]);
    }
  } catch (e: any) {
    $showError(e);
  }

  fileStore.reload = true;
};

const openSearch = () => {
  layoutStore.showHover("search");
};

const toggleMultipleSelection = () => {
  fileStore.toggleMultiple();
  layoutStore.closeHovers();
};

const windowsResize = throttle(() => {
  colunmsResize();
  width.value = window.innerWidth;

  // Listing element is not displayed
  if (listing.value == null) return;

  // How much every listing item affects the window height
  setItemWeight();

  // Fill but not fit the window
  fillWindow();
}, 100);

const download = () => {
  if (fileStore.req === null) return;

  if (
    fileStore.selectedCount === 1 &&
    !fileStore.req.items[fileStore.selected[0]].isDir
  ) {
    api.download(null, fileStore.req.items[fileStore.selected[0]].url);
    return;
  }

  layoutStore.showHover({
    prompt: "download",
    confirm: (format: any) => {
      layoutStore.closeHovers();

      const files = [];

      if (fileStore.selectedCount > 0 && fileStore.req !== null) {
        for (const i of fileStore.selected) {
          files.push(fileStore.req.items[i].url);
        }
      } else {
        files.push(route.path);
      }

      api.download(format, ...files);
    },
  });
};

const switchView = async () => {
  layoutStore.closeHovers();

  const modes = {
    list: "mosaic",
    mosaic: "mosaic gallery",
    "mosaic gallery": "list",
  };

  const data = {
    id: authStore.user?.id,
    viewMode: (modes[authStore.user?.viewMode ?? "list"] ||
      "list") as ViewModeType,
  };

  users.update(data, ["viewMode"]).catch($showError);

  authStore.updateUser(data);

  setItemWeight();
  fillWindow();
};

const uploadFunc = () => {
  if (
    typeof window.DataTransferItem !== "undefined" &&
    typeof DataTransferItem.prototype.webkitGetAsEntry !== "undefined"
  ) {
    layoutStore.showHover("upload");
  } else {
    document.getElementById("upload-input")?.click();
  }
};

const setItemWeight = () => {
  // Listing element is not displayed
  if (listing.value === null || fileStore.req === null) return;

  let itemQuantity = fileStore.req.numDirs + fileStore.req.numFiles;
  if (itemQuantity > showLimit.value) itemQuantity = showLimit.value;

  // How much every listing item affects the window height
  itemWeight.value = listing.value.offsetHeight / itemQuantity;
};

const fillWindow = (fit = false) => {
  if (fileStore.req === null) return;

  const totalItems = fileStore.req.numDirs + fileStore.req.numFiles;

  // More items are displayed than the total
  if (showLimit.value >= totalItems && !fit) return;

  const windowHeight = window.innerHeight;

  // Quantity of items needed to fill 2x of the window height
  const showQuantity = Math.ceil(
    (windowHeight + windowHeight * 2) / itemWeight.value
  );

  // Less items to display than current
  if (showLimit.value > showQuantity && !fit) return;

  // Set the number of displayed items
  showLimit.value = showQuantity > totalItems ? totalItems : showQuantity;
};
</script>
