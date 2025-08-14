<template>
  <div class="wp">
    <iframe id="inlineFrameExample" :src="workflowUrl"></iframe>
  </div>
</template>

<script>

export default {
  data(){
    return{
      workflowUrl: process.env.NODE_ENV === 'development'
        ? 'http://192.168.35.247:8081/work_flow?workflow_id=7537585355608817664&space_id=7535746948309975040'
        : window.location.origin + this.$basePath + '/work_flow?workflow_id=7535747020951126016&space_id=7535746948309975040'
    }
  },
  mounted() {
    window.addEventListener("message", (e) => {
      console.log('message=========>',e.data);
    });
  },
}
</script>

<style lang="scss" scoped>
.wp {
  width: 100%;
  height: 100%;
  background: #f4f5f6;
  iframe {
    width: 100%;
    height: 100%;
    border: none;
  }
}
</style>
