<template>
  <div>
    <HeaderBar />
    <router-view></router-view>
  </div>
</template>

<script>
import HeaderBar from "@/components/header/HeaderBar.vue";
// eslint-disable-next-line no-undef
__webpack_public_path__ = window.FileBrowser.StaticURL + "/";

export default {
  name: "app",
  components: {
    HeaderBar
  },
  mounted() {
    const loading = document.getElementById("loading");
    loading.classList.add("done");

    setTimeout(function () {
      loading.parentNode.removeChild(loading);
    }, 200);
  },
};
</script>

<style>
@import "./css/styles.css";
</style>
