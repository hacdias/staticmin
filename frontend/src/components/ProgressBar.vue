<!-- This component taken directly from vue-simple-progress
since it didnt support Vue 3 but the component itself does
https://raw.githubusercontent.com/dzwillia/vue-simple-progress/master/src/components/Progress.vue -->
<template>
  <div>
    <div
      class="vue-simple-progress-text"
      :style="text_style"
      v-if="text.length > 0 && textPosition == 'top'"
    >
      {{ text }}
    </div>
    <div class="vue-simple-progress" :style="progress_style">
      <div
        class="vue-simple-progress-text"
        :style="text_style"
        v-if="text.length > 0 && textPosition == 'middle'"
      >
        {{ text }}
      </div>
      <div
        style="position: relative; left: -9999px"
        :style="text_style"
        v-if="text.length > 0 && textPosition == 'inside'"
      >
        {{ text }}
      </div>
      <div class="vue-simple-progress-bar" :style="bar_style">
        <div
          :style="text_style"
          v-if="text.length > 0 && textPosition == 'inside'"
        >
          {{ text }}
        </div>
      </div>
    </div>
    <div
      class="vue-simple-progress-text"
      :style="text_style"
      v-if="text.length > 0 && textPosition == 'bottom'"
    >
      {{ text }}
    </div>
  </div>
</template>

<script>
// We're leaving this untouched as you can read in the beginning
const isNumber = function (n) {
  return !isNaN(parseFloat(n)) && isFinite(n);
};

export default {
  name: "progress-bar",
  props: {
    val: {
      default: 0,
    },
    max: {
      default: 100,
    },
    size: {
      // either a number (pixel width/height) or 'tiny', 'small',
      // 'medium', 'large', 'huge', 'massive' for common sizes
      default: 3,
    },
    "bg-color": {
      type: String,
      default: "#eee",
    },
    "bar-color": {
      type: String,
      default: "#2196f3", // match .blue color to Material Design's 'Blue 500' color
    },
    "bar-transition": {
      type: String,
      default: "all 0.5s ease",
    },
    "bar-border-radius": {
      type: Number,
      default: 0,
    },
    spacing: {
      type: Number,
      default: 4,
    },
    text: {
      type: String,
      default: "",
    },
    "text-align": {
      type: String,
      default: "center", // 'left', 'right'
    },
    "text-position": {
      type: String,
      default: "bottom", // 'bottom', 'top', 'middle', 'inside'
    },
    "font-size": {
      type: Number,
      default: 13,
    },
    "text-fg-color": {
      type: String,
      default: "#222",
    },
  },
  computed: {
    pct() {
      let pct = (this.val / this.max) * 100;
      pct = pct.toFixed(2);
      return Math.min(pct, this.max);
    },
    size_px() {
      switch (this.size) {
        case "tiny":
          return 2;
        case "small":
          return 4;
        case "medium":
          return 8;
        case "large":
          return 12;
        case "big":
          return 16;
        case "huge":
          return 32;
        case "massive":
          return 64;
      }

      return isNumber(this.size) ? this.size : 32;
    },
    text_padding() {
      switch (this.size) {
        case "tiny":
        case "small":
        case "medium":
        case "large":
        case "big":
        case "huge":
        case "massive":
          return Math.min(Math.max(Math.ceil(this.size_px / 8), 3), 12);
      }

      return isNumber(this.spacing) ? this.spacing : 4;
    },
    text_font_size() {
      switch (this.size) {
        case "tiny":
        case "small":
        case "medium":
        case "large":
        case "big":
        case "huge":
        case "massive":
          return Math.min(Math.max(Math.ceil(this.size_px * 1.4), 11), 32);
      }

      return isNumber(this.fontSize) ? this.fontSize : 13;
    },
    progress_style() {
      const style = {
        background: this.bgColor,
      };

      if (this.textPosition == "middle" || this.textPosition == "inside") {
        style["position"] = "relative";
        style["min-height"] = this.size_px + "px";
        style["z-index"] = "-2";
      }

      if (this.barBorderRadius > 0) {
        style["border-radius"] = this.barBorderRadius + "px";
      }

      return style;
    },
    bar_style() {
      const style = {
        background: this.barColor,
        width: this.pct + "%",
        height: this.size_px + "px",
        transition: this.barTransition,
      };

      if (this.barBorderRadius > 0) {
        style["border-radius"] = this.barBorderRadius + "px";
      }

      if (this.textPosition == "middle" || this.textPosition == "inside") {
        style["position"] = "absolute";
        style["top"] = "0";
        style["height"] = "100%";
        (style["min-height"] = this.size_px + "px"), (style["z-index"] = "-1");
      }

      return style;
    },
    text_style() {
      const style = {
        color: this.textFgColor,
        "font-size": this.text_font_size + "px",
        "text-align": this.textAlign,
      };

      if (
        this.textPosition == "top" ||
        this.textPosition == "middle" ||
        this.textPosition == "inside"
      )
        style["padding-bottom"] = this.text_padding + "px";
      if (
        this.textPosition == "bottom" ||
        this.textPosition == "middle" ||
        this.textPosition == "inside"
      )
        style["padding-top"] = this.text_padding + "px";

      return style;
    },
  },
};
</script>
