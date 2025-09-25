<template>
  <el-button @click="handleCopy" class="copy-icon">
    <i class="el-icon-document-copy"></i>
    复制
  </el-button>
</template>

<script>
export default {
  name: "CopyIcon",
  props: {
    iconClass: {
      type: String,
      required: true
    }
  },
  methods: {
    async handleCopy() {
      try {
        const text = this.iconClass;

        // 优先使用现代 Clipboard API
        if (navigator.clipboard && window.isSecureContext) {
          await navigator.clipboard.writeText(text);
        } else {
          // 降级方案：创建 input 并使用 execCommand
          const input = document.createElement('input');
          input.value = text;
          input.setAttribute('readonly', '');
          input.style.cssText = 'position: absolute; left: -9999px;';
          document.body.appendChild(input);
          input.select();
          document.execCommand('copy');
          document.body.removeChild(input);
        }

        this.$message.success('已复制到剪贴板');
      } catch (err) {
        console.error('复制失败:', err);
        this.$message.error('复制失败，请手动复制');
      }
    }
  }
}
</script>

<style scoped>
.copy-icon {
  margin-left: 10px;
}
</style>