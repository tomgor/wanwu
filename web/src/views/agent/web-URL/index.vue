<template>
    <div class="weburl-container">
        <div class="weburl-title">
            <span class="el-icon-arrow-left goback"></span>
            <span class="weburl-title-text">发布配置</span>
        </div>
        <CommonLayout
        :showAside="true"
        :showTitle="false"
        :asideWidth="asideWidth"
        class="weburl-content"
        >
        <template #aside>
            <div v-for="item in toolList" class="toolList">
                <h3>{{item.name}}</h3>
                <p>{{item.desc}}</p>
            </div>
        </template>
        <template #main-content>

        </template>
        </CommonLayout>
    </div>
</template>
<script>
import CommonLayout from '@/components/exploreContainer.vue';
import CreateApi from './createApi.vue';
import CreateUrl from './createUrl.vue';
export default {
    components: {CommonLayout,CreateApi,CreateUrl},
    data(){
        return{
            asideWidth:'320',
            toolList:[
                {
                    name:'Web URL',
                    desc:'分享链接后，可通过网页直接访问智能体应用',
                    type:'url'
                },
                {
                    name:'API',
                    desc:'支持嵌入第三方应用系统',
                    type:'api'
                }
            ]
        }
    }
}
</script>
<style lang="scss" scoped>
    .weburl-container{
        width:100%;
        height:100%;
        padding:0 10px;
        .weburl-title{
            width:100%;
            height:60px;
            padding:0 20px;
            border-bottom: 1px solid #dbdbdb;
            display:flex;
            align-items: center;
            .goback{
                font-size: 20px;
                margin-right: 10px;
                cursor:pointer;
            }
            .weburl-title-text{
                font-size:18px;
                font-weight: bold;
            }
        }
        .weburl-content{
            width:100%;
            padding:10px;
            height:calc(100% - 60px);
                .toolList{
                    cursor:pointer;
                    border:1px solid #dbdbdb;
                    text-align:center;
                    width:90%;
                    height:80px;
                    border-radius:6px;
                    padding:10px;
                    margin:20px auto;
                    h3{
                        font-size:16px;
                    }
                    p{
                        color:#666;
                        padding-top:10px;
                    }
                }
            
        }
        
    }
</style>