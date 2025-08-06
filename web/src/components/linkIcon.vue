<template>
  <img
    style="height: 14px; margin-left: 5px; margin-top: 1px; cursor: pointer"
    src="@/assets/imgs/linkIcon.png"
    @click="handleClick"
  />
</template>

<script>
import { mapGetters } from "vuex"

export default {
  props: {
    type: ''
  },
  data() {
    return {
      basePath: this.$basePath,
      linkList: {}
    }
  },
  watch: {
    commonInfo:{
      handler(val) {
        const { linkList } = val.data || {}
        this.linkList = linkList || {}
      },
      deep: true
    }
  },
  computed: {
    ...mapGetters('user', ['commonInfo']),
  },
  methods: {
    handleClick() {
      console.log(this.type, '-----------------------------type')
      const docUrl = this.linkList[this.type]
      if (docUrl) window.open(docUrl)
    }
  }
}
</script>

<style lang="scss" scoped>

</style>
