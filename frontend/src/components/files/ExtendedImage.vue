<template>
  <div
    class="image-ex-container"
    ref="container"
    @touchstart="touchStart"
    @touchmove="touchMove"
    @dblclick="zoomAuto"
    @mousedown="mousedownStart"
    @mousemove="mouseMove"
    @mouseup="mouseUp"
    @wheel="wheelMove"
  >
    <img class="image-ex-img image-ex-img-center" ref="imgex" @load="onLoad" />
  </div>
</template>
<script setup lang="ts">
import { throttle } from "lodash-es";
import UTIF from "utif";
import { onBeforeUnmount, onMounted, ref, watch } from "vue";

interface IProps {
  src: string;
  moveDisabledTime?: number;
  classList?: any[];
  zoomStep?: number;
}

const props = withDefaults(defineProps<IProps>(), {
  moveDisabledTime: () => 200,
  classList: () => [],
  zoomStep: () => 0.25,
});

const scale = ref<number>(1);
const lastX = ref<number | null>(null);
const lastY = ref<number | null>(null);
const inDrag = ref<boolean>(false);
const touches = ref<number>(0);
const lastTouchDistance = ref<number | null>(0);
const moveDisabled = ref<boolean>(false);
const disabledTimer = ref<number | null>(null);
const imageLoaded = ref<boolean>(false);
const position = ref<{
  center: { x: number; y: number };
  relative: { x: number; y: number };
}>({
  center: { x: 0, y: 0 },
  relative: { x: 0, y: 0 },
});
const maxScale = ref<number>(4);
const minScale = ref<number>(0.25);

// Refs
const imgex = ref<HTMLImageElement | null>(null);
const container = ref<HTMLDivElement | null>(null);

onMounted(() => {
  if (!decodeUTIF() && imgex.value !== null) {
    imgex.value.src = props.src;
  }

  props.classList.forEach((className) =>
    container.value !== null ? container.value.classList.add(className) : ""
  );

  if (container.value === null) {
    return;
  }

  // set width and height if they are zero
  if (getComputedStyle(container.value).width === "0px") {
    container.value.style.width = "100%";
  }
  if (getComputedStyle(container.value).height === "0px") {
    container.value.style.height = "100%";
  }

  window.addEventListener("resize", onResize);
});

onBeforeUnmount(() => {
  window.removeEventListener("resize", onResize);
  document.removeEventListener("mouseup", onMouseUp);
});

watch(
  () => props.src,
  () => {
    if (!decodeUTIF() && imgex.value !== null) {
      imgex.value.src = props.src;
    }

    scale.value = 1;
    setZoom();
    setCenter();
  }
);

// Modified from UTIF.replaceIMG
const decodeUTIF = () => {
  const sufs = ["tif", "tiff", "dng", "cr2", "nef"];
  if (document?.location?.pathname === undefined) {
    return;
  }
  const suff =
    document.location.pathname.split(".")?.pop()?.toLowerCase() ?? "";

  if (sufs.indexOf(suff) == -1) return false;
  const xhr = new XMLHttpRequest();
  UTIF._xhrs.push(xhr);
  UTIF._imgs.push(imgex.value);
  xhr.open("GET", props.src);
  xhr.responseType = "arraybuffer";
  xhr.onload = UTIF._imgLoaded;
  xhr.send();
  return true;
};

const onLoad = () => {
  imageLoaded.value = true;

  if (imgex.value === null) {
    return;
  }

  imgex.value.classList.remove("image-ex-img-center");
  setCenter();
  imgex.value.classList.add("image-ex-img-ready");

  document.addEventListener("mouseup", onMouseUp);

  let realSize = imgex.value.naturalWidth;
  let displaySize = imgex.value.offsetWidth;

  // Image is in portrait orientation
  if (imgex.value.naturalHeight > imgex.value.naturalWidth) {
    realSize = imgex.value.naturalHeight;
    displaySize = imgex.value.offsetHeight;
  }

  // Scale needed to display the image on full size
  const fullScale = realSize / displaySize;

  // Full size plus additional zoom
  maxScale.value = fullScale + 4;
};

const onMouseUp = () => {
  inDrag.value = false;
};

const onResize = throttle(function () {
  if (imageLoaded.value) {
    setCenter();
    doMove(position.value.relative.x, position.value.relative.y);
  }
}, 100);

const setCenter = () => {
  if (container.value === null || imgex.value === null) {
    return;
  }

  position.value.center.x = Math.floor(
    (container.value.clientWidth - imgex.value.clientWidth) / 2
  );
  position.value.center.y = Math.floor(
    (container.value.clientHeight - imgex.value.clientHeight) / 2
  );

  imgex.value.style.left = position.value.center.x + "px";
  imgex.value.style.top = position.value.center.y + "px";
};

const mousedownStart = (event: Event) => {
  lastX.value = null;
  lastY.value = null;
  inDrag.value = true;
  event.preventDefault();
};
const mouseMove = (event: MouseEvent) => {
  if (!inDrag.value) return;
  doMove(event.movementX, event.movementY);
  event.preventDefault();
};
const mouseUp = (event: Event) => {
  inDrag.value = false;
  event.preventDefault();
};
const touchStart = (event: TouchEvent) => {
  lastX.value = null;
  lastY.value = null;
  lastTouchDistance.value = null;
  if (event.targetTouches.length < 2) {
    setTimeout(() => {
      touches.value = 0;
    }, 300);
    touches.value++;
    if (touches.value > 1) {
      zoomAuto(event);
    }
  }
  event.preventDefault();
};

const zoomAuto = (event: Event) => {
  switch (scale.value) {
    case 1:
      scale.value = 2;
      break;
    case 2:
      scale.value = 4;
      break;
    default:
    case 4:
      scale.value = 1;
      setCenter();
      break;
  }
  setZoom();
  event.preventDefault();
};

const touchMove = (event: TouchEvent) => {
  event.preventDefault();
  if (lastX.value === null) {
    lastX.value = event.targetTouches[0].pageX;
    lastY.value = event.targetTouches[0].pageY;
    return;
  }
  if (imgex.value === null) {
    return;
  }
  const step = imgex.value.width / 5;
  if (event.targetTouches.length === 2) {
    moveDisabled.value = true;
    if (disabledTimer.value) clearTimeout(disabledTimer.value);
    disabledTimer.value = window.setTimeout(
      () => (moveDisabled.value = false),
      props.moveDisabledTime
    );

    const p1 = event.targetTouches[0];
    const p2 = event.targetTouches[1];
    const touchDistance = Math.sqrt(
      Math.pow(p2.pageX - p1.pageX, 2) + Math.pow(p2.pageY - p1.pageY, 2)
    );
    if (!lastTouchDistance.value) {
      lastTouchDistance.value = touchDistance;
      return;
    }
    scale.value += (touchDistance - lastTouchDistance.value) / step;
    lastTouchDistance.value = touchDistance;
    setZoom();
  } else if (event.targetTouches.length === 1) {
    if (moveDisabled.value) return;
    const x = event.targetTouches[0].pageX - (lastX.value ?? 0);
    const y = event.targetTouches[0].pageY - (lastY.value ?? 0);
    if (Math.abs(x) >= step && Math.abs(y) >= step) return;
    lastX.value = event.targetTouches[0].pageX;
    lastY.value = event.targetTouches[0].pageY;
    doMove(x, y);
  }
};

const doMove = (x: number, y: number) => {
  if (imgex.value === null) {
    return;
  }
  const style = imgex.value.style;

  const posX = pxStringToNumber(style.left) + x;
  const posY = pxStringToNumber(style.top) + y;

  style.left = posX + "px";
  style.top = posY + "px";

  position.value.relative.x = Math.abs(position.value.center.x - posX);
  position.value.relative.y = Math.abs(position.value.center.y - posY);

  if (posX < position.value.center.x) {
    position.value.relative.x = position.value.relative.x * -1;
  }

  if (posY < position.value.center.y) {
    position.value.relative.y = position.value.relative.y * -1;
  }
};
const wheelMove = (event: WheelEvent) => {
  scale.value += -Math.sign(event.deltaY) * props.zoomStep;
  setZoom();
};
const setZoom = () => {
  scale.value = scale.value < minScale.value ? minScale.value : scale.value;
  scale.value = scale.value > maxScale.value ? maxScale.value : scale.value;
  if (imgex.value !== null)
    imgex.value.style.transform = `scale(${scale.value})`;
};
const pxStringToNumber = (style: string) => {
  return +style.replace("px", "");
};
</script>
<style>
.image-ex-container {
  margin: auto;
  overflow: hidden;
  position: relative;
}

.image-ex-img {
  position: absolute;
}

.image-ex-img-center {
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
  position: absolute;
  transition: none;
}

.image-ex-img-ready {
  left: 0;
  top: 0;
  transition: transform 0.1s ease;
}
</style>
