<template>
    <div class="weburl-container">
        <div class="weburl-title">
            <span class="el-icon-arrow-left goback"></span>
            发布配置
        </div>
        <CommonLayout
        :showAside="true"
        >
            <aside-content>
                123
            </aside-content>
            <main-content>
                456
            </main-content>
        </CommonLayout>
    </div>
</template>
<script>
import CommonLayout from '@/components/exploreContainer.vue'
export default {
    components: {CommonLayout},
    data(){
        return{
        }
    }
}
</script>
<style lang="scss" scoped>
    .weburl-container{
        width:100%;
        height:100%;
        border:1px solid red;
        padding:20px;
        .weburl-title{
            width:100%;
            height:60px;
            padding:0 20px;
            border-bottom: 1px solid #dbdbdb;
            .goback{
                font-size: 18px;
                margin-right: 10px;
            }
        }
        
    }
</style>