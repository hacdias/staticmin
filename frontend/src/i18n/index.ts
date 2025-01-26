import dayjs from "dayjs";
import { createI18n } from "vue-i18n";

import("dayjs/locale/ar");
import("dayjs/locale/de");
import("dayjs/locale/el");
import("dayjs/locale/en");
import("dayjs/locale/es");
import("dayjs/locale/fr");
import("dayjs/locale/he");
import("dayjs/locale/hu");
import("dayjs/locale/is");
import("dayjs/locale/it");
import("dayjs/locale/ja");
import("dayjs/locale/ko");
import("dayjs/locale/nl-be");
import("dayjs/locale/pl");
import("dayjs/locale/pt-br");
import("dayjs/locale/pt");
import("dayjs/locale/ro");
import("dayjs/locale/ru");
import("dayjs/locale/sk");
import("dayjs/locale/sv");
import("dayjs/locale/tr");
import("dayjs/locale/uk");
import("dayjs/locale/zh-cn");
import("dayjs/locale/zh-tw");

// All i18n resources specified in the plugin `include` option can be loaded
// at once using the import syntax
import messages from "@intlify/unplugin-vue-i18n/messages";

export function detectLocale() {
  // locale is an RFC 5646 language tag
  // https://developer.mozilla.org/en-US/docs/Web/API/Navigator/language
  let locale = navigator.language.toLowerCase();
  switch (true) {
    case /^he\b/.test(locale):
      locale = "he";
      break;
    case /^hu\b/.test(locale):
      locale = "hu";
      break;
    case /^ar\b/.test(locale):
      locale = "ar";
      break;
    case /^el.*/i.test(locale):
      locale = "el";
      break;
    case /^es\b/.test(locale):
      locale = "es";
      break;
    case /^en\b/.test(locale):
      locale = "en";
      break;
    case /^is\b/.test(locale):
      locale = "is";
      break;
    case /^it\b/.test(locale):
      locale = "it";
      break;
    case /^fr\b/.test(locale):
      locale = "fr";
      break;
    case /^pt-br\b/.test(locale):
      locale = "pt-br";
      break;
    case /^pt\b/.test(locale):
      locale = "pt";
      break;
    case /^ja\b/.test(locale):
      locale = "ja";
      break;
    case /^zh-tw\b/.test(locale):
      locale = "zh-tw";
      break;
    case /^zh-cn\b/.test(locale):
    case /^zh\b/.test(locale):
      locale = "zh-cn";
      break;
    case /^de\b/.test(locale):
      locale = "de";
      break;
    case /^ro\b/.test(locale):
      locale = "ro";
      break;
    case /^ru\b/.test(locale):
      locale = "ru";
      break;
    case /^pl\b/.test(locale):
      locale = "pl";
      break;
    case /^ko\b/.test(locale):
      locale = "ko";
      break;
    case /^sk\b/.test(locale):
      locale = "sk";
      break;
    case /^tr\b/.test(locale):
      locale = "tr";
      break;
    // ua wasnt a valid locale for ukraine
    case /^uk\b/.test(locale):
      locale = "uk";
      break;
    case /^sv-se\b/.test(locale):
    case /^sv\b/.test(locale):
      locale = "sv";
      break;
    case /^nl-be\b/.test(locale):
      locale = "nl-be";
      break;
    default:
      locale = "en";
  }

  return locale;
}

// TODO: was this really necessary?
// function removeEmpty(obj: Record<string, any>): void {
//   Object.keys(obj)
//     .filter((k) => obj[k] !== null && obj[k] !== undefined && obj[k] !== "") // Remove undef. and null and empty.string.
//     .reduce(
//       (newObj, k) =>
//         typeof obj[k] === "object"
//           ? Object.assign(newObj, { [k]: removeEmpty(obj[k]) }) // Recurse.
//           : Object.assign(newObj, { [k]: obj[k] }), // Copy value.
//       {}
//     );
// }

export const rtlLanguages = ["he", "ar"];

export const i18n = createI18n({
  locale: detectLocale(),
  fallbackLocale: "en",
  messages,
  // expose i18n.global for outside components
  legacy: true,
});

export const isRtl = (locale?: string) => {
  // see below
  // @ts-expect-error incorrect type when legacy
  return rtlLanguages.includes(locale || i18n.global.locale.value);
};

export function setLocale(locale: string) {
  dayjs.locale(locale);
  // according to doc u only need .value if legacy: false but they lied
  // https://vue-i18n.intlify.dev/guide/essentials/scope.html#local-scope-1
  // @ts-expect-error incorrect type when legacy
  i18n.global.locale.value = locale;
}

export function setHtmlLocale(locale: string) {
  const html = document.documentElement;
  html.lang = locale;
  if (isRtl(locale)) html.dir = "rtl";
  else html.dir = "ltr";
}

export default i18n;
