<template>
  <video ref="videoPlayer" class="video-max video-js" controls preload="auto">
    <source />
    <track
      kind="subtitles"
      v-for="(sub, index) in subtitles"
      :key="index"
      :src="sub"
      :label="subLabel(sub)"
      :default="index === 0"
    />
    <p class="vjs-no-js">
      Sorry, your browser doesn't support embedded videos, but don't worry, you
      can <a :href="source">download it</a>
      and watch it with your favorite video player!
    </p>
  </video>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, nextTick } from "vue";
import videojs from "video.js";
import type Player from "video.js/dist/types/player";
import "videojs-mobile-ui";
import "videojs-hotkeys";
import "video.js/dist/video-js.min.css";
import "videojs-mobile-ui/dist/videojs-mobile-ui.css";

const videoPlayer = ref<HTMLElement | null>(null);
const player = ref<Player | null>(null);

const props = withDefaults(
  defineProps<{
    source: string;
    subtitles?: string[];
    options?: any;
  }>(),
  {
    options: {},
  }
);

const source = ref(props.source);
const sourceType = ref("");

nextTick(() => {
  initVideoPlayer();
});

onMounted(() => {});

onBeforeUnmount(() => {
  if (player.value) {
    player.value.dispose();
    player.value = null;
  }
});

const initVideoPlayer = async () => {
  try {
    const lang = document.documentElement.lang;
    const languagePack = await (
      languageImports[lang] || languageImports.en
    )?.();
    videojs.addLanguage("videoPlayerLocal", languagePack.default);
    sourceType.value = "";

    //
    sourceType.value = getSourceType(source.value);

    const srcOpt = { sources: { src: props.source, type: sourceType.value } };
    //Supporting localized language display.
    const langOpt = { language: "videoPlayerLocal" };
    // support for playback at different speeds.
    const playbackRatesOpt = { playbackRates: [0.5, 1, 1.5, 2, 2.5, 3] };
    const options = getOptions(
      props.options,
      langOpt,
      srcOpt,
      playbackRatesOpt
    );
    player.value = videojs(videoPlayer.value!, options, () => {});

    // TODO: need to test on mobile
    // @ts-expect-error no ts definition for mobileUi
    player.value!.mobileUi();
  } catch (error) {
    console.error("Error initializing video player:", error);
  }
};

const getOptions = (...srcOpt: any[]) => {
  const options = {
    controlBar: {
      skipButtons: {
        forward: 5,
        backward: 5,
      },
    },
    html5: {
      nativeTextTracks: false,
    },
    plugins: {
      hotkeys: {
        volumeStep: 0.1,
        seekStep: 10,
        enableModifiersForNumbers: false,
      },
    },
  };

  return videojs.obj.merge(options, ...srcOpt);
};

//  Attempting to fix the issue of being unable to play .MKV format video files
const getSourceType = (source: string) => {
  const fileExtension = source ? source.split("?")[0].split(".").pop() : "";
  if (fileExtension?.toLowerCase() === "mkv") {
    return "video/mp4";
  }
  return "";
};

const subLabel = (subUrl: string) => {
  let url: URL;
  try {
    url = new URL(subUrl);
  } catch {
    // treat it as a relative url
    // we only need this for filename
    url = new URL(subUrl, window.location.origin);
  }

  const label = decodeURIComponent(
    url.pathname
      .split("/")
      .pop()!
      .replace(/\.[^/.]+$/, "")
  );

  return label;
};

interface LanguageImports {
  [key: string]: () => Promise<any>;
}

const languageImports: LanguageImports = {
  he: () => import("video.js/dist/lang/he.json"),
  hu: () => import("video.js/dist/lang/hu.json"),
  ar: () => import("video.js/dist/lang/ar.json"),
  de: () => import("video.js/dist/lang/de.json"),
  el: () => import("video.js/dist/lang/el.json"),
  en: () => import("video.js/dist/lang/en.json"),
  es: () => import("video.js/dist/lang/es.json"),
  fr: () => import("video.js/dist/lang/fr.json"),
  it: () => import("video.js/dist/lang/it.json"),
  ja: () => import("video.js/dist/lang/ja.json"),
  ko: () => import("video.js/dist/lang/ko.json"),
  "nl-be": () => import("video.js/dist/lang/nl.json"),
  pl: () => import("video.js/dist/lang/pl.json"),
  "pt-br": () => import("video.js/dist/lang/pt-BR.json"),
  pt: () => import("video.js/dist/lang/pt-PT.json"),
  ro: () => import("video.js/dist/lang/ro.json"),
  ru: () => import("video.js/dist/lang/ru.json"),
  sk: () => import("video.js/dist/lang/sk.json"),
  tr: () => import("video.js/dist/lang/tr.json"),
  uk: () => import("video.js/dist/lang/uk.json"),
  "zh-cn": () => import("video.js/dist/lang/zh-CN.json"),
  "zh-tw": () => import("video.js/dist/lang/zh-TW.json"),
};
</script>
<style scoped>
.video-max {
  width: 100%;
  height: 100%;
}
</style>
