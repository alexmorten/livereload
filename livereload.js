var source = new EventSource("/livereload");
source.addEventListener("message", function (_event) {
  console.log("Reloading page");
  window.location.reload();
});

source.addEventListener("open", function (_event) {
  console.log("Livereload connection opened");
});

source.addEventListener("error", function (event) {
  console.log("error:");
  console.log(event);
});
