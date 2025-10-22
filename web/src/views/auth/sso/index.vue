<template>
  <div v-loading="true" class="loading-dev" />
</template>

<script>
import { mapActions } from "vuex";
import { getImgVerCode } from "@/api/user";
import { urlEncrypt } from "@/utils/crypto";
import { redirectUrl, rawQuery } from "@/utils/util";

const rawQuery1 = rawQuery()
export default {
  data() {
    return {
      rawQuery1
    };
  },
  beforeCreate() {
  },
  mounted() {
        console.log("ssoLogin", this.rawQuery1);

    this.doLogin();
  },
  methods: {
    ...mapActions("user", ["ssoLogin", "LoginOut"]),
    doLogin() {
      const code = this.$route.query.code || this.$route.query.key;
      if (!code) {
        return this.$message.warning("未提供参数！");
      }
      try {
        this.ssoLogin(this.rawQuery1);
      } catch (e) {
        console.log("error", e);
      }
    },
    // 工具：从 location.href 里取原始查询串，不 decode，保留参数内的 =
    
  },
};
</script>

<style lang="scss" scoped>
.loading-dev {
  width: 200px;
  height: 100px;
}
</style>
