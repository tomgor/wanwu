<template>
    <el-dialog
        title="创建标签"
        :visible.sync="dialogVisible"
        width="20%"
        :before-close="handleClose">
        <div>
            <el-input placeholder="搜索标签" suffix-icon="el-icon-search" @keyup.enter.native="addByEnterKey" v-model="tagName"></el-input>
            <div class="add"  @click="addTag"><span class="el-icon-plus add-icon"></span>创建标签</div>
            <div class="tag-box">
                <div 
                v-for="item in tagList" 
                class="tag_item"  
                @mouseenter="mouseEnter(item)"
                @mouseleave="mouseLeave(item)"
                @dblclick="handleDoubleClick(item)"
                >
                    <el-checkbox v-model="item.selected" v-if="!item.showIpt">{{item.tagName}}</el-checkbox>
                    <el-input v-model="item.tagName" v-if="item.showIpt" @blur="inputBlur(item)" ></el-input>
                    <span class="el-icon-close del-icon" v-if="item.showDel && !item.showIpt" @click="delTag(item)"></span>
                </div>
            </div>
        </div>
        <span slot="footer" class="dialog-footer">
            <el-button type="primary" @click="submitDialog">确 定</el-button>
        </span>
    </el-dialog>
</template>
<script>
import { delTag,tagList,createTag,editTag,bindTag } from "@/api/knowledge";
export default {
    data(){
        return{
            dialogVisible:false,
            tagList:[],
            tagName:'',
            knowledgeId:''
        }
    },
    created(){
    //    this.getList();
    },
    methods:{
        submitDialog(){
            const ids = this.tagList.filter(item => item.selected).map(item => item.tagId)
            bindTag({knowledgeId:this.knowledgeId,tagIdList:ids}).then(res =>{
                if(res.code === 0){
                    this.$emit('relodaData')
                }
            })
            this.dialogVisible = false
        },
        getList(){
            tagList({knowledgeId:this.knowledgeId,tagName:this.tagName}).then(res => {
                if(res.code === 0){
                    this.tagList = res.data.knowledgeTagList.map(item =>({
                        ...item,
                        showDel:false,
                        showIpt:false
                    }))
                }
            })
        },
        delTag(item){
            if(item.selected) return
            delTag({tagId:item.tagId}).then(res =>{
                if(res.code === 0){
                    this.getList()
                }
            })
        },
        showDiaglog(id){
            this.knowledgeId = id
            this.dialogVisible = true;
            this.getList();
        },
        handleClose(){
            this.dialogVisible = false;
        },
        mouseEnter(n){
            n.showDel = true
        },
        mouseLeave(n){
            n.showDel = false
        },
        handleDoubleClick(n){
            n.showIpt = true
        },
        inputBlur(n){
            if(n.tagId){
                this.edit_tag(n);
            }else{
                this.add_Tag(n);
            }
        },
        add_Tag(n){
            createTag({tagName:n.tagName}).then(res =>{
                if(res.code === 0){
                    n.showIpt = false;
                    this.getList();
                }
            })
        },
        edit_tag(n){
            editTag({tagId:n.tagId,tagName:n.tagName}).then(res =>{
                if(res.code === 0){
                    n.showIpt = false;
                    this.getList();
                }
            })
        },
        addTag(){
            this.tagList.unshift(
                {
                    tagName:'',
                    checked:false,
                    showDel:false,
                    showIpt:true
                }
            )
        },
        addByEnterKey(e){
            if (e.keyCode === 13){
                 this.getList();
            }
        }
    }
}
</script>
<style lang="scss" scoped>
/deep/{
    .el-dialog__body{
        padding:5px 20px!important;
    }
    .add{
        margin-top:10px;
        padding:10px 0;
        cursor: pointer;
      .add-icon{
        margin-right:5px;
      }
    }
    .tag-box{
        max-height:300px;
        overflow-y:scroll;
    }
    .tag_item{
        cursor: pointer;
        background:#F4F5FF;
        padding:5px;
        margin:10px 0;
        border-radius:4px;
        display:flex;
        justify-content:space-between;
        align-items:center;
        .del-icon{
            color: #384BF7;
            cursor: pointer;
            font-size: 16px;
        }
    }
}
</style>