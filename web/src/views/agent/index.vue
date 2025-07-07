<template>
    <CommonLayout
    :aside-title="asideTitle" 
    :isButton="true"
    @handleBtnClick="handleBtnClick"
    :isBtnDisabled="sessionStatus === 0"
  >
    <template #aside-content>
      <div class="explore-aside-app">
        <div v-for="n in historyList " class="appList" 
        :class="['appList',{'disabled':sessionStatus === 0},{'active':n.active}]"
        @click="historyClick(n)"  
        @mouseenter="mouseEnter(n)"
        @mouseleave="mouseLeave(n)">
           <span class="appName">
                <span class="appTag"></span>
                {{n.title}}
            </span>
            <span class="el-icon-delete appDelete" v-if="n.hover || n.active" @click.stop="deleteConversation(n)"></span>
        </div>
      </div>
    </template>
    <template #main-content>
      <div class="app-content">
        <div class="app-header-api">
            <div class="header-api-box">
                <div class="header-api-url">
                    <el-tag  effect="plain" class="root-url">API根地址</el-tag>
                    {{apiURL}}
                </div>
                <el-button size="small" @click="openApiDialog" plain class="apikeyBtn" >
                    <img :src="require('@/assets/imgs/apikey.png')" />
                    API秘钥
                </el-button>
            </div>
        </div>
        <Chat :chatType="'chat'" :editForm="editForm"  ref="agentChat" @reloadList="reloadList" @setHistoryStatus="setHistoryStatus"/>
        <ApiKeyDialog ref="apiKeyDialog" :appId="editForm.assistantId" :appType="'agent'" />
      </div>
    </template>
  </CommonLayout>
</template>
<script>
import CommonLayout from '@/components/exploreContainer.vue'
import Chat from './components/chat.vue'
 import {mapGetters} from 'vuex'
 import ApiKeyDialog from './components/ApiKeyDialog.vue'
 import { getAgentInfo } from "@/api/agent";
 import {getApiKeyRoot} from "@/api/appspace";
import sseMethod from '@/mixins/sseMethod'
import { getConversationlist} from "@/api/agent";
export default {
    components: {CommonLayout,Chat,ApiKeyDialog },
    mixins: [sseMethod],
    data(){
        return {
            apiURL:'',
            asideTitle:'新建对话',
            assistantId:'',
            historyList:[],
            editForm:{
                assistantId:'',
                avatar:{},
                name:'',
                desc:'',
                prologue:''
            },
        }
    },
    computed: {
        ...mapGetters('app', ['sessionStatus'])
    },
    created(){
        if(this.$route.query.id){
            this.assistantId = this.$route.query.id
            this.editForm.assistantId = this.$route.query.id;
            this.getDetail()
            this.getList()
            this.apiKeyRootUrl()
        }
    },
    methods:{
      reloadList(val){
        this.getList(val)
      },
      getDetail(){
            getAgentInfo({assistantId:this.editForm.assistantId}).then(res =>{
                if(res.code === 0){
                    this.editForm.avatar = res.data.avatar;
                    this.editForm.name = res.data.name;
                    this.editForm.desc = res.data.desc;
                    this.editForm.prologue = res.data.prologue;
                }
            })
        },
        getList(noInit){
            getConversationlist({assistantId:this.assistantId,pageNo:1,pageSize:1000}).then(res =>{
                if(res.code === 0){
                    if(res.data.list && res.data.list.length > 0){
                        this.historyList = res.data.list.map(n =>{
                            return {...n, hover: false, active: false}
                        })
                        if (noInit) {
                            this.historyList[0].active = true  //noInit 是true时，左侧默认选中第一个,但是不要调接口刷新详情
                        } else {
                            this.historyClick[this.historyList[0]]
                        }
                    }else{
                        this.historyList = []
                    }
                }else{
                    this.historyList = []
                }
            })
        },
        setHistoryStatus(){
            this.historyList.forEach(m => {
                m.active = false
            })
        },
        historyClick(n){//切换对话
            n.hover = true;
            this.$refs['agentChat'].conversionClick(n)
        },
        deleteConversation(n){
            this.$refs['agentChat'].preDelConversation(n)
        },
        handleBtnClick(){//新建对话
            this.$refs['agentChat'].createConversion()
        },
         async doDeleteHistory(){//删除历史记录
            let res = await deleteConversationHistory({conversationId:this.conversationId})
            if(res.code === 0){
                this.$message.success(this.$t('yuanjing.deleteTips'))
            }
        },
        mouseEnter(n){
            n.hover = true;
        },
        mouseLeave(n){
            n.hover = false;
        },
        apiKeyRootUrl(){
            const data = {appId:this.editForm.assistantId,appType:'agent'}
            getApiKeyRoot(data).then(res =>{
                if(res.code === 0){
                    this.apiURL = res.data || ''
                }
            })
        },
        openApiDialog(){
            this.$refs.apiKeyDialog.showDialog();
        }
    }
}
</script>
<style lang="scss" scoped>
@import '@/style/chat.scss';
.active{
    background-color: $color_opacity !important;
    .appTag{
        background-color: #384BF7!important;
    }
}
.explore-aside-app{
    .appList:hover{
        background-color: $color_opacity !important;
    }
    .appList{
        margin:10px 20px;
        padding:10px;
        border-radius: 6px;
        margin-bottom: 6px;
        display: flex;
        gap: 8px;
        align-items: center;
        justify-content:space-between;
        cursor: pointer;
        position:relative;
        .appDelete{
            color:#384BF7;
            margin-right:-5px;
            cursor: pointer;
        }
        .appName{
            display: block;
            max-width: 130px;
            overflow: hidden;
            white-space: nowrap;
            pointer-events: none;
            text-overflow: ellipsis;
            .appTag{
                display:inline-block;
                width:8px;
                height:8px;
                border-radius:50%;
                background:#ccc;
            }
      }
    }
}
 .app-content{
    width:100%;
    height:100%;
    .app-header-api{
        width:100%;
        padding:10px;
        position:absolute;
        z-index:999;
        top:0;
        left:0;
        border-bottom:1px solid #eaeaea;
        display:flex;
        justify-content:flex-end;
        align-content:center;
        .header-api-box{
            display:flex;
            .header-api-url{
                padding: 6px 10px;
                background:#fff;
                margin:0 10px;
                border-radius:6px;
                .root-url{
                    background-color:#ECEEFE;
                    color:#384BF7;
                    border:none;
                }
            }
        }
    }
 }
</style>