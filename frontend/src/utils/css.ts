export default function getRule(rules: string[]) {
  for (let i = 0; i < rules.length; i++) {
    rules[i] = rules[i].toLowerCase();
  }

  let result = null;
  const find = Array.prototype.find;

  find.call(document.styleSheets, (styleSheet: CSSStyleSheet) => {
    result = find.call(styleSheet.cssRules, (cssRule: CSSRule) => {
      let found = false;

      // faster than checking instanceof for every element
      if (cssRule.constructor.name === "CSSStyleRule") {
        for (let i = 0; i < rules.length; i++) {
          if (
            (cssRule as CSSStyleRule).selectorText.toLowerCase() === rules[i]
          ) {
            found = true;
          }
        }
      }

      return found;
    });

    return result != null;
  });

  return result as CSSStyleRule | null;
}
