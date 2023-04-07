console.log("file javascript");
const promise = new Promise((resolve, reject) => {
  setTimeout(() => resolve("finished promise"), 5000);
});
(async () => {
  const out = await promise;
  console.log(out);
})();
