<template>
  <el-button @click="handleCopy" class="copy-icon">
    <i class="el-icon-document-copy"></i>
    复制
  </el-button>
</template>

<script>

export default {
  name: "copyIcon",
  props: {
    iconClass: {
      type: String,
      required: true
    }
  },
  methods: {
    async handleCopy() {
      try {
        await navigator.clipboard.writeText(this.iconClass);
        this.$message.success("复制成功");
      } catch (err) {
        this.$message.error("复制失败");
      }
    },
  },
};
</script>

<style lang="scss" scoped>
.copy-icon {
  color: #333;
  background: #f5f5f5;
  border-radius: 4px;
  padding: 5px 10px;
  cursor: pointer;
  margin-left: 10px;
}
</style>

