<template>
    <div>
        <el-dialog
            title="新增工具"
            :visible.sync="dialogVisible"
            width="40%"
            :before-close="handleClose">
            <div class="tool-typ">
                <div class="toolbtn">
                    <div v-for="(item,index) in toolList" :key="index" @click="clickTool(item,index)" :class="[{'active':toolIndex === index}]">
                        {{item.name}}
                    </div>
                </div>
                <el-input v-model="toolName" placeholder="搜索工具" class="tool-input" suffix-icon="el-icon-search" @keyup.enter.native="searchTool" clearable></el-input>
            </div>
            <div class="toolContent">
                <template v-for="(items, type) in contentMap">
                    <div 
                    v-if="activeValue === type"
                    v-for="(item,i) in items"
                    :key="item[type + 'Id'] || item.id"
                    class="toolContent_item"
                    >
                    <span>{{ item.apiName || item.name }}</span>
                    <el-checkbox v-model="item.checked" @change="openTool($event,item,type)" :disabled="item.checked"></el-checkbox>
                    </div>
                </template>
            </div>
            <span slot="footer" class="dialog-footer">
                <el-button type="primary" @click="submit">确 定</el-button>
            </span>
        </el-dialog>
    </div>
</template>
<script>
import { getList } from '@/api/workflow.js';
import { addWorkFlowInfo, addMcp } from "@/api/agent";
import { getExplorationFlowList,readWorkFlow} from "@/api/workflow";
import { Base64 } from "js-base64";
export default {
    props:['assistantId'],
    data(){
        return {
            toolName:'',
            dialogVisible:false,
            toolIndex:0,
            activeValue:'mcp',
            actionInfos:[],
            workFlowInfos:[],
            mcpInfos:[],
            toolList:[
                // {
                //     value:'auto',
                //     name:'自定义'
                // },
                {
                    value:'mcp',
                    name:'MCP'
                },
                {
                    value:'workflow',
                    name:'工作流'
                }
            ]
        }
    },
    computed:{
         contentMap() {
            return {
            auto: this.actionInfos,
            mcp: this.mcpInfos,
            workflow: this.workFlowInfos
            }
        }
    },
    created(){
        this.getMcpSelect('');
        this.getWorkflowList('');
    },
    methods:{
        openTool(e,item,type){
            if(!e) return;
            item.checked = !item.checked
            if(type === 'workflow'){
                this.addWorkFlow(item)
            }else if(type === 'mcp'){
                this.addMcpItem(item)
            }
        },
        addMcpItem(n){
            addMcp({assistantId:this.assistantId,mcpId:n.mcpId}).then(res =>{
                if(res.code === 0){
                    n.checked = true;
                    this.$message.success('工具添加成功');
                    this.$emit('updateDetail');
                }
            }).catch(() =>{

            })
        },
        addWorkFlow(n){
            let params = { workflowID: n.appId};
            readWorkFlow(params).then(res => {
                if(res.code === 0){
                    this.doCreateWorkFlow(n,n.appId, res.data.base64OpenAPISchema);
                }
            })
        },
        async doCreateWorkFlow(n,workFlowId, schema){
            let params = {
                assistantId: this.assistantId,
                schema: Base64.decode(schema),
                workFlowId,
                apiAuth: {
                type: "none",
                },
            };
            let res = await addWorkFlowInfo(params);
            if (res.code === 0) {
                n.checked = true;
                this.$message.success(this.$t('agent.addPluginTips'));
                this.$emit('updateDetail');
            }
        },
        searchTool(){
            if(this.activeValue === 'auto'){
                console.log('自定义')
            }else if(this.activeValue === 'mcp'){
                this.getMcpSelect(this.toolName)
            }else{
                this.getWorkflowList(this.toolName)
            }
        },
        getMcpSelect(name){
            getList({name}).then(res => {
                if(res.code === 0){
                     this.mcpInfos = (res.data.list || []).map(item => 
                        Object.assign({}, item, { checked: false })
                    );
                }
               
            }).catch(err => {

            })
            },
            getWorkflowList(name) {
                getExplorationFlowList({name,appType:'workflow',searchType:'all'}).then(res =>{
                    if (res.code === 0) {
                        this.workFlowInfos = (res.data.list || []).map(item => 
                            Object.assign({}, item, { checked: false })
                        );
                    }
                })
        },
        showDialog(row){
            this.dialogVisible = true;
            this.setMcp(row.mcpInfos);
            this.setWorkflow(row.workFlowInfos);
        },
        setMcp(data){
           this.mcpInfos = this.mcpInfos.map(m => ({
            ...m,
            checked: data.includes(m.mcpId)
            }));
        },
        setWorkflow(data){
            this.workFlowInfos = this.workFlowInfos.map(m => ({
            ...m,
            checked: data.includes(m.appId)
            }));
        },
        handleClose(){
            this.toolIndex = -1;
            this.activeValue = '';
            this.dialogVisible = false;
        },
        clickTool(item,i){
            this.toolIndex = i;
            this.activeValue = item.value;
        },
        submit(){
            this.$emit('selectTool',this.activeValue);
            this.dialogVisible = false;
        }
    }
}
</script>
<style lang="scss" scoped>
/deep/{
    .el-dialog__body{
        padding:10px 20px;
    }
}
.tool-typ{
    display:flex;
    justify-content:space-between;
    padding:10px 0;
    border-bottom: 1px solid #dbdbdb;
    .toolbtn{
        display:flex;
        justify-content:flex-start;
        gap:20px;
        div{
            text-align: center;
            padding:5px 20px;
            border-radius:6px;
            border:1px solid #ddd;
            cursor: pointer;
        }
    }
    .tool-input{
        width:200px;
    }
}
.toolContent{
    padding:10px 0;
    max-height:300px;
    overflow-y:auto;
    .toolContent_item{
        padding:15px 20px;
        border:1px solid #dbdbdb;
        border-radius:6px;
        margin-bottom:10px;
        cursor: pointer;
        display: flex;
        justify-content:space-between;
    }
    .toolContent_item:hover{
        background:#f4f5ff;
    }
}
.active{
    border:1px solid #384BF7 !important;
    color: #fff;
    background:#384BF7;
}
</style>