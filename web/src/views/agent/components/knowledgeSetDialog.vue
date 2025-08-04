<template>
    <div>
        <el-dialog
        title="召回参数配置"
        :visible.sync="dialogVisible"
        width="50%"
        :before-close="handleClose">
        <span v-if="dialogVisible">
           <searchConfig ref='searchConfig' @sendConfigInfo="sendConfigInfo" :setType="'agent'" :config="knowledgeConfig"/>
        </span>
        <span slot="footer" class="dialog-footer">
            <el-button @click="dialogVisible = false">取 消</el-button>
            <el-button type="primary" @click="submit">确 定</el-button>
        </span>
        </el-dialog>
    </div>
</template>
<script>
import searchConfig from '@/components/searchConfig.vue';
export default {
    components:{
        searchConfig
    },
    data(){
        return{
            dialogVisible:false,
            configInfo:{},
            knowledgeConfig:{}
        }
    },
    methods:{
        sendConfigInfo(data){
            this.configInfo = { ...data.knowledgeMatchParams };
        },
        showDialog(row){
            this.dialogVisible = true;
            this.knowledgeConfig = row || {};
        },
        handleClose(){
            this.dialogVisible = false;
        },
        submit(){
            this.dialogVisible = false;
            this.$emit('setKnowledgeSet',this.configInfo)
        }
    }
}
</script>
<style lang="scss" scoped>
/deep/{
    .el-input-number--small{
        line-height: 28px!important;
    }
}
.question{
    cursor: pointer;
    color:#ccc;
}
</style>