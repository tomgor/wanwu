<template>
  <div class="txt-box">
    <div class="txt" :loading="loading" style="width: 100%; height: 100%;white-space: pre-wrap;">{{ textContent }}</div>
  </div>
</template>
<script>
import axios from 'axios'
export default {
  name: 'DocPeview',
  components: {},
  data() {
    return {
      fileUrl: "",
      loading: true
    }
  },
  watch: {
    $route: {
      deep: true,
      handler: function (to) {
        if (to.query) {
          let fileUrl = to.query.fileUrl
          this.fileUrl = fileUrl
        }
      },
      immediate: true
    }
  },
  created() { },
  mounted() {
    this.$nextTick(() => {
      this.getPreview()
    })
  },
  methods: {
    async getPreview() {
      this.loading = true
      axios({
        method: 'get',
        responseType: "text",
        url: this.fileUrl,
        headers: {
          'Accept': 'text/plain; charset=utf-8' // 告诉服务器我们期望接受的编码格式
       }
      }).then(({ data }) => {
        this.$nextTick(() => {
          this.textContent = data
          this.loading = false
        })
      })
    },
    renderedHandler() {
      console.log('渲染完成')
    },
    errorHandler() {
      console.log('渲染失败')
    }
  }
}
</script>
<style lang="scss" scoped>
.txt-box {
  padding: 20px;
}
</style>
