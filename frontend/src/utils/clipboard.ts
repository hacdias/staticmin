// Based on code by the following links:
// https://stackoverflow.com/a/74528564
// https://web.dev/articles/async-clipboard

interface ClipboardArgs {
  text?: string;
  data?: ClipboardItems;
}

interface ClipboardOpts {
  permission?: boolean;
}

export function copy(data: ClipboardArgs, opts?: ClipboardOpts) {
  return new Promise<void>((resolve, reject) => {
    if (
      // Clipboard API requires secure context
      window.isSecureContext &&
      typeof navigator.clipboard !== "undefined"
    ) {
      if (opts?.permission) {
        getPermission("clipboard-write")
          .then(() => writeToClipboard(data).then(resolve).catch(reject))
          .catch(reject);
      } else {
        writeToClipboard(data).then(resolve).catch(reject);
      }
    } else if (
      document.queryCommandSupported &&
      document.queryCommandSupported("copy") &&
      data.text // old method only supports text
    ) {
      const textarea = createTemporaryTextarea(data.text);
      const body = document.activeElement || document.body;
      try {
        body.appendChild(textarea);
        textarea.focus();
        textarea.select();
        document.execCommand("copy");
        resolve();
      } catch (e) {
        reject(e);
      } finally {
        body.removeChild(textarea);
      }
    } else {
      reject(
        new Error("None of copying methods are supported by this browser!")
      );
    }
  });
}

function getPermission(name: string) {
  return new Promise<void>((resolve, reject) => {
    typeof navigator.permissions !== "undefined" &&
      navigator.permissions
        // @ts-expect-error chrome specific api
        .query({ name })
        .then((permission) => {
          if (permission.state === "granted" || permission.state === "prompt") {
            resolve();
          } else {
            reject(new Error("Permission denied!"));
          }
        });
  });
}

function writeToClipboard(data: ClipboardArgs) {
  if (data.text) {
    return navigator.clipboard.writeText(data.text);
  }
  if (data.data) {
    return navigator.clipboard.write(data.data);
  }

  return new Promise<void>((resolve, reject) => {
    reject(new Error("No data was supplied!"));
  });
}

const styles = {
  fontSize: "12pt",
  position: "fixed",
  top: 0,
  left: 0,
  width: "2em",
  height: "2em",
  padding: 0,
  margin: 0,
  border: "none",
  outline: "none",
  boxShadow: "none",
  background: "transparent",
};

function createTemporaryTextarea(text: string) {
  const textarea = document.createElement("textarea");
  textarea.value = text;
  textarea.setAttribute("readonly", "");
  Object.assign(textarea.style, styles);
  return textarea;
}
