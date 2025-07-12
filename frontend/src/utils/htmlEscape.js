function htmlUnescape(str) {
      const doc = new DOMParser().parseFromString(str, "text/html");
      return doc.documentElement.textContent;
}

export default htmlUnescape;
